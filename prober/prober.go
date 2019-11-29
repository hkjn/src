// Package prober provides black-box monitoring mechanisms.
//
// To use, define Probe() and Alert() on a type, then pass it to NewProbe:
//   struct FooProber{ someState int }
//
//   // Probe "Foo". E.g. do a network call and compare it to what
//   // was expected.
//   func (p FooProber) Probe() Result {
//     // Returning FailedWith(err) indicates that the probe failed.
//     // Returning Passed() indicates that the probe succeeded.
//   }
//   // Send an alert. Called if the probe fails "too often".
//   //
//   // By passing in FailurePenalty() and/or SuccessReward() options to NewProbe(),
//   // the adjustments to the state when probe fails or passes can be modified.
//   func (p FooProber) Alert(name, desc string, badness int, records Records) error {
//   }
//   ...
//
//   // Create the probe.
//   p := NewProbe(FooProber{1}, "FooProber", "Probes the Foo")
//
//   // Run the probe. This call blocks forever, so you may
//   // want to do this in a goroutine — you could e.g. register a web
//   // handler to show the contents of p.Records() here.
//   go p.Run()
package prober

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	DefaultInterval = time.Minute // default duration to pause between Probe() runs
	// TODO(hkjn): Replace this with a global "max 1 alerts per X
	// minutes" setting? If we have 1000 probes, we still would get 1000
	// alerts every 15 min..
	MaxAlertFrequency     = time.Minute * 15      // never call Alert() more often than this
	logDir                = os.TempDir()          // logging directory
	logName               = "prober.outcomes.log" // name of logging f1ile
	alertThreshold        = flag.Int("alert_threshold", 200, "level of 'badness' before alerting")
	alertsDisabled        = flag.Bool("no_alerts", false, "disables alerts when probes fail too often")
	disabledProbes        = make(selectedProbes)
	onlyProbes            = make(selectedProbes)
	defaultFailurePenalty = 10 // default increment of `badness` on failed probe run
	defaultSuccessReward  = 1  // default decrement of `badness` on successful probe run
	onceOpen              sync.Once
	logFile               *os.File
	bufferSize            = 200 // maximum number of results per prober to keep
	parseFlags            = sync.Once{}
	results               = [2]string{"Pass", "Fail"}
)

const (
	Pass ResultCode = iota
	Fail
)

type (
	// Result describes the outcome of a single probe.
	Result struct {
		Code    ResultCode
		Error   error
		Info    string // Optional extra information
		InfoUrl string // Optional URL to further information
	}

	// ResultCode describes pass/fail outcomes for probes.
	ResultCode int

	// Record is the result of a single probe run.
	Record struct {
		Timestamp  time.Time `yaml:"-"`
		TimeMillis string    // same as Timestamp but makes it into the YAML logs
		Result     Result    // the result of the probe run
	}

	// Records is a grouping of probe records that implements sort.Interface.
	Records []Record

	// AlertFn is function that is called when a Prober is alerting.
	AlertFn func(name, desc string, badness int, records Records) error

	// Prober is a mechanism that can probe some target(s).
	Prober interface {
		Probe() Result                                               // probe target(s) once
		Alert(name, desc string, badness int, records Records) error // send alert
	}

	// Option is a setting for an individual prober.
	Option func(*Probe)

	// selectedProbes is a set of probes that flags specify should be enabled/disabled.
	selectedProbes map[string]bool

	// Probe is a stateful representation of repeated probe runs.
	Probe struct {
		Prober                      // underlying prober mechanism
		Name, Desc    string        // name, description of the probe
		Interval      time.Duration // how often to probe
		Disabled      bool          // whether this probe is disabled
		SilencedUntil SilenceTime   // the earliest time this probe can alert
		// If `badness` reaches alert threshold, an alert email is sent and
		// the value resets to 0.
		badness        int
		failurePenalty int          // how much to increment `badness` on failure
		successReward  int          // how much to decrement `badness` on success
		reportFn       func(Result) // function to call to report probe results
		t              timeT
		alerting       bool         // whether this probe is currently alerting
		lastAlert      time.Time    // time of last alert sent, if any
		alertLock      sync.RWMutex // protects reads and writes to alerting state
		records        Records      // historical records of probe runs
		recordsLock    sync.RWMutex // protects reads and writes to stateful records
	}
	Probes []*Probe
	// SilenceTime represents a Time until which the probe is
	// silenced. It exists to provide a custom String() method.
	SilenceTime struct{ time.Time }

	// timeT represents time-dependent functionality.
	timeT interface {
		Now() time.Time
		Sleep(time.Duration)
	}
)

// realTime implements timeT for actual time.
type realTime struct{}

func (realTime) Now() time.Time        { return time.Now() }
func (realTime) Sleep(d time.Duration) { time.Sleep(d) }

// String returns the English name of the result.
func (r ResultCode) String() string { return results[r] }

// String returns a human-readable representation of the Result.
func (r Result) String() string {
	parts := []string{
		fmt.Sprintf("Code: %q", r.Code),
	}
	if r.Error != nil {
		parts = append(parts, fmt.Sprintf("Error: %q", r.Error))
	}
	if r.Info != "" {
		parts = append(parts, fmt.Sprintf("Info: %q", r.Info))
	}
	if r.InfoUrl != "" {
		parts = append(parts, fmt.Sprintf("InfoUrl: %q", r.InfoUrl))
	}
	return fmt.Sprintf("Result{%s}", strings.Join(parts, ", "))
}

// Passed returns whether the probe result indicates a pass.
func (r Result) Passed() bool { return r.Code == Pass }

func (r1 Result) Equal(r2 Result) bool {
	if r1.Code != r2.Code {
		return false
	}
	if r1.Info != r2.Info {
		return false
	}
	equalError := func(err1, err2 error) bool {
		if err1 == nil {
			return err2 == nil
		}
		if err2 == nil {
			return err1 == nil
		}
		return err1.Error() == err2.Error()
	}
	return equalError(r1.Error, r2.Error)
}

// FailedWith returns a Result representing failure with given error.
func FailedWith(err error) Result {
	return Result{
		Code:  Fail,
		Error: err,
		Info:  fmt.Sprintf("The probe failed with %q", err.Error()),
	}
}

// FailedWith returns a Result representing failure with given error and extra information.
func FailedWithInfo(err error, info, infoUrl string) Result {
	return Result{
		Code:    Fail,
		Error:   err,
		Info:    info,
		InfoUrl: infoUrl,
	}
}

// Passed returns a Result representing pass.
func Passed() Result { return Result{Code: Pass} }

// PasseWith returns a Result representing pass, with extra info.
func PassedWith(info, url string) Result {
	return Result{
		Code:    Pass,
		Info:    info,
		InfoUrl: url,
	}
}

// String returns the flag's value.
func (d *selectedProbes) String() string {
	s := ""
	i := 0
	for p, _ := range *d {
		if i > 0 {
			s += ","
		}
		s += p
	}
	return s
}

// Get is part of the flag.Getter interface. It always returns nil for
// this flag type since the struct is not exported.
func (d *selectedProbes) Get() interface{} { return nil }

// Syntax: -disable_probes=FooProbe,BarProbe
func (d *selectedProbes) Set(value string) error {
	vals := strings.Split(value, ",")
	m := *d
	for _, p := range vals {
		m[p] = true
	}
	return nil
}

// NewProbe returns a new probe from given prober implementation.
func NewProbe(p Prober, name, desc string, options ...Option) *Probe {
	parseFlags.Do(func() {
		if !flag.Parsed() {
			flag.Parse()
		}
	})
	probe := &Probe{
		Prober:         p,
		Name:           name,
		Desc:           desc,
		Interval:       DefaultInterval,
		badness:        0,
		failurePenalty: defaultFailurePenalty,
		successReward:  defaultSuccessReward,
		records:        Records{},
		t:              realTime{},
		alertLock:      sync.RWMutex{},
	}
	for _, opt := range options {
		opt(probe)
	}
	return probe
}

// Interval sets the interval for the prober.
func Interval(interval time.Duration) func(*Probe) {
	return func(p *Probe) {
		p.Interval = interval
	}
}

// Report sets the function to call to report probe results.
func Report(fn func(Result)) func(*Probe) {
	return func(p *Probe) {
		p.reportFn = fn
	}
}

// FailurePenalty sets the amount `badness` is incremented on failure for the prober.
func FailurePenalty(penalty int) func(*Probe) {
	return func(p *Probe) {
		p.failurePenalty = penalty
	}
}

// SuccessReward sets the amount `badness` is decremented on success for the prober.
func SuccessReward(reward int) func(*Probe) {
	return func(p *Probe) {
		p.successReward = reward
	}
}

// Run repeatedly runs the probe, blocking forever.
func (p *Probe) Run() {
	log.Printf("[%s] Starting..\n", p.Name)

	if !enabledInFlags(p.Name) {
		p.Disabled = true
		log.Printf("[%s] is disabled, will now exit", p.Name)
		return
	}

	for {
		wait := p.runProbe()
		p.t.Sleep(wait)
	}
}

// String returns a human-readable representation of the Probe.
func (p *Probe) String() string {
	parts := []string{
		fmt.Sprintf("Name: %q", p.Name),
		fmt.Sprintf("Desc: %q", p.Desc),
		fmt.Sprintf("Records: %s", p.Records()),
	}
	if p.Badness() != 0 {
		parts = append(parts, fmt.Sprintf("Badness: %d", p.Badness()))
	}
	if p.Interval != DefaultInterval {
		parts = append(parts, fmt.Sprintf("Interval: %v", p.Interval))
	}
	if p.IsAlerting() {
		parts = append(parts, fmt.Sprintf("alerting: true"))
	}
	lastAlert := p.getLastAlert()
	if !lastAlert.IsZero() {
		parts = append(parts, fmt.Sprintf("lastAlert: %v", lastAlert))
	}
	if p.Disabled {
		parts = append(parts, fmt.Sprintf("Disabled: true"))
	}
	if !p.SilencedUntil.Equal(time.Time{}) {
		parts = append(parts, fmt.Sprintf("SilencedUntil: %v", p.SilencedUntil))
	}
	if p.failurePenalty != defaultFailurePenalty {
		parts = append(parts, fmt.Sprintf("failurePenalty: %v", p.failurePenalty))
	}
	return fmt.Sprintf("&Probe{%s}", strings.Join(parts, ", "))
}

// enabledInFlags returns true if this probe is enabled via -only_probes or -disabled_probes flags.
func enabledInFlags(name string) bool {
	if len(onlyProbes) > 0 {
		if _, ok := onlyProbes[name]; ok {
			// We only want specific probes, but this probe is one of them.
			return true
		}
		return false
	}

	if _, ok := disabledProbes[name]; ok {
		// This probe is explicitly disabled.
		return false
	}
	return true
}

// runProbe runs the probe once, returning the amount of time to wait
// before the next runProbe() run is due.
func (p *Probe) runProbe() time.Duration {
	c := make(chan Result, 1)
	start := p.t.Now()
	go func() {
		log.Printf("[%s] Probing..\n", p.Name)
		c <- p.Probe()
	}()
	select {
	case r := <-c:
		// We got a result of some sort from the prober.
		p.handleResult(r)
		wait := p.Interval - p.t.Now().Sub(start)
		log.Printf("[%s] needs to sleep %v more here\n", p.Name, wait)
		return wait
	case <-time.After(p.Interval):
		// Probe didn't finish in time for us to run the next one, report as failure.
		log.Printf("[%s] Timed out\n", p.Name)
		timeoutFail := FailedWith(
			fmt.Errorf("%s timed out (with probe interval %1.1f sec)",
				p.Name,
				p.Interval.Seconds()))
		p.handleResult(timeoutFail)
		return time.Duration(0)
	}
}

// Records returns the historical records of probe runs.
func (p *Probe) Records() Records {
	p.recordsLock.RLock()
	defer p.recordsLock.RUnlock()
	return p.records
}

// add appends the record to the buffer for the probe, keeping it within bufferSize.
func (p *Probe) addRecord(r Record) {
	p.recordsLock.Lock()
	p.records = append(p.records, r)
	if len(p.records) >= bufferSize {
		over := len(p.records) - bufferSize
		log.Printf("[%s] buffer is over %d, reslicing it\n", p.Name, bufferSize)
		p.records = p.records[over:]
	}
	p.recordsLock.Unlock()
	log.Printf("[%s] buffer is now %d elements\n", p.Name, len(p.Records()))
}

// Silenced returns true if the probe is currently silenced.
func (p *Probe) Silenced() bool {
	return p.SilencedUntil.After(p.t.Now())
}

// Silence silences the Probe until specified time.
func (p *Probe) Silence(until time.Time) {
	p.SilencedUntil = SilenceTime{until}
	log.Printf("[%s] is now silenced until %v\n", p.Name, until)
}

// String returns a human-readable description of the time until which a probe is silenced.
func (t SilenceTime) String() string {
	return fmt.Sprintf("%s (%f hrs more)", t.Format(time.RFC822), t.Sub(time.Now()).Hours())
}

// Equal returns true if the probes are equal.
func (p1 *Probe) Equal(p2 *Probe) bool {
	if p2 == nil {
		return false
	}
	if p1.Name != p2.Name {
		return false
	}
	if p1.Desc != p2.Desc {
		return false
	}
	if !p1.Records().Equal(p2.Records()) {
		return false
	}
	if p1.Badness() != p2.Badness() {
		return false
	}
	if p1.Interval != p2.Interval {
		return false
	}
	if p1.IsAlerting() != p2.IsAlerting() {
		return false
	}
	if !p1.getLastAlert().Equal(p2.getLastAlert()) {
		return false
	}
	if p1.Disabled != p2.Disabled {
		return false
	}
	if !p1.SilencedUntil.Equal(p2.SilencedUntil.Time) {
		return false
	}
	if p1.failurePenalty != p2.failurePenalty {
		return false
	}
	return true
}

// Equal returns true if the Records are equal.
func (rs1 Records) Equal(rs2 Records) bool {
	if len(rs1) != len(rs2) {
		return false
	}
	for i, r1 := range rs1 {
		r2 := rs2[i]
		if !r1.Equal(r2) {
			return false
		}
	}
	return true
}

// Implement sort.Interface for Records. The sort order is chronological.
func (rs Records) Len() int           { return len(rs) }
func (rs Records) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs Records) Less(i, j int) bool { return rs[i].Timestamp.Before(rs[j].Timestamp) }

func (rs Records) String() string {
	s := make([]string, len(rs))
	for i, r := range rs {
		s[i] = r.String()
	}

	return strings.Join(s, ", ")
}

// RecentFailures returns only recent probe failures among the records.
func (pr Records) RecentFailures() Records {
	failures := make(Records, 0)
	for _, r := range pr {
		if !r.Result.Passed() && !r.Timestamp.Before(time.Now().Add(-time.Hour)) {
			failures = append(failures, r)
		}
	}
	sort.Sort(sort.Reverse(failures))
	return failures
}

func (r Record) String() string {
	return fmt.Sprintf(
		"Record{Timestamp: %v, TimeMillis: %q, Result: %s}",
		r.Timestamp,
		r.TimeMillis,
		r.Result)
}

// Ago describes the duration since the record occured.
func (r Record) Ago() string {
	d := time.Since(r.Timestamp)
	if d < time.Minute {
		return fmt.Sprintf("%0.1f sec ago", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%0.1f min ago", d.Minutes())
	} else if d < time.Hour*24 {
		return fmt.Sprintf("%0.1f hrs ago", d.Hours())
	} else {
		return fmt.Sprintf("%0.1f days ago", d.Hours()/24.0)
	}
}

// Marshal returns the record in YAML form.
func (r Record) marshal() []byte {
	b, err := yaml.Marshal(r)
	if err != nil {
		log.Printf("failed to marshal record %+v: %v", r, err)
	}
	return b
}

// Equal returns true if the Record objects are equal.
func (r1 Record) Equal(r2 Record) bool {
	if !r1.Timestamp.Equal(r2.Timestamp) {
		return false
	}
	if r1.TimeMillis != r2.TimeMillis {
		return false
	}
	if !r1.Result.Equal(r2.Result) {
		return false
	}
	return true
}

// openLog opens the log file.
func openLog() {
	logPath := filepath.Join(logDir, logName)
	log.Printf("Using YAML log file %q\n", logPath)
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Printf("failed to open %q: %v\n", logPath, err)
	}
	logFile = f
}

// handleResult handles a return value from a Probe() run.
func (p *Probe) handleResult(r Result) {
	if p.reportFn != nil {
		// Call custom report function, if specified.
		p.reportFn(r)
	}
	b := p.Badness()
	if r.Passed() {
		b -= p.successReward
		if b < 0 {
			b = 0
		}
		log.Printf("[%s] Pass, badness is now %d.\n", p.Name, b)
	} else {
		b += p.failurePenalty
		log.Printf("[%s] Failed while probing, badness is now %d: %v\n", p.Name, b, r.Error)
	}
	p.setBadness(b)
	p.logResult(r)

	if p.Silenced() {
		log.Printf("[%s] is silenced until %v, will not alert, resetting badness to 0\n", p.Name, p.SilencedUntil)
		p.setBadness(0)
	}

	p.setIsAlerting(p.Badness() >= *alertThreshold)
	if !p.IsAlerting() {
		return
	}
	if *alertsDisabled {
		log.Printf("[%s] would now be alerting, but alerts are disabled\n", p.Name)
		return
	}

	lastAlert := p.getLastAlert()
	if time.Since(lastAlert) < MaxAlertFrequency {
		log.Printf("[%s] will not alert, since last alert was sent %v back\n", p.Name, time.Since(lastAlert))
		return
	}

	log.Printf("[%s] is alerting\n", p.Name)
	// Send alert notification in goroutine to not block further
	// probing.
	// TODO: There is a race condition here, if email sending takes long
	// enough for further Probe() runs to finish, which would queue up
	// several duplicate alert emails. This shouldn't often happen, but
	// technically should be bounded by a timeout to prevent the
	// possibility.
	go p.sendAlert()
}

// setIsAlerting changes the alerting status of the probe.
func (p *Probe) setIsAlerting(alerting bool) {
	p.alertLock.Lock()
	p.alerting = alerting
	p.alertLock.Unlock()
}

// IsAlerting returns true if the Probe is currently alerting.
func (p *Probe) IsAlerting() bool {
	p.alertLock.RLock()
	defer p.alertLock.RUnlock()
	return p.alerting
}

// Badness returns the current `badness` value.
func (p *Probe) Badness() int {
	p.alertLock.RLock()
	defer p.alertLock.RUnlock()
	return p.badness
}

// setBadness sets the `badness` to specified value.
func (p *Probe) setBadness(b int) {
	p.alertLock.Lock()
	p.badness = b
	p.alertLock.Unlock()
}

// setLastAlert sets the last alert fired to given time.
func (p *Probe) setLastAlert(t time.Time) {
	p.alertLock.Lock()
	p.lastAlert = t
	p.alertLock.Unlock()
}

// getLastAlert returns the time the Probe last alerted.
//
// getLastAlert returns the zero value for time.Time if the Probe has never alerted.
func (p *Probe) getLastAlert() time.Time {
	p.alertLock.RLock()
	defer p.alertLock.RUnlock()
	return p.lastAlert
}

// sendAlert calls the Alert() implementation and handles the outcome.
func (p *Probe) sendAlert() {
	err := p.Alert(p.Name, p.Desc, p.Badness(), p.Records())
	if err != nil {
		log.Printf("[%s] Failed to alert: %v", p.Name, err)
		// Note: We don't reset badness here; next cycle we'll keep
		// trying to send the alert.
	} else {
		log.Printf("[%s] Called Alert(), resetting badness to 0\n", p.Name)
		p.setLastAlert(p.t.Now())
		p.setBadness(0)
	}
}

// logResult logs the result of a probe run.
func (p *Probe) logResult(res Result) {
	onceOpen.Do(openLog)
	now := p.t.Now()
	rec := Record{
		Timestamp:  now,
		TimeMillis: now.Format(time.StampMilli),
		Result:     res,
	}

	p.addRecord(rec)
	_, err := logFile.Write(rec.marshal())
	if err != nil {
		log.Printf("failed to write record to log: %v", err)
	}
}

// Silenced returns the currently silenced probes, if any.
func (ps Probes) Silenced() Probes {
	var silenced Probes
	for _, p := range ps {
		if p.Silenced() {
			silenced = append(silenced, p)
		}
	}
	return silenced
}

// Equal returns true if both Probes are equal.
func (ps1 Probes) Equal(ps2 Probes) bool {
	if len(ps1) != len(ps2) {
		return false
	}
	for i, p1 := range ps1 {
		if !ps2[i].Equal(p1) {
			return false
		}
	}
	return true
}

// Implement sort.Interface for Probes.
func (ps Probes) Len() int { return len(ps) }

// Less returns true if probe i should sort before probe j.
//
// Less is implemented to give the natural order that's likely to be
// most useful when ordering Probes, i.e. with the ones requiring
// attention first. Since the default sort order is ascending, this
// means that "lower values" will correspond to probes in worse state.
func (ps Probes) Less(i, j int) bool {
	if ps[i].Disabled != ps[j].Disabled {
		// Disabled probes sort after (higher value than) non-disabled ones.
		return ps[j].Disabled
	}
	if !ps[i].SilencedUntil.Equal(ps[j].SilencedUntil.Time) {
		// Probes that are silenced for a longer time sort after (higher
		// value than) ones that are silenced shorter. (Possibly not
		// silenced at all, but that depends on the current time.)
		return ps[i].SilencedUntil.Before(ps[j].SilencedUntil.Time)
	}
	b1, b2 := ps[i].Badness(), ps[j].Badness()
	if b1 != b2 {
		// Probes with higher `badness` sort before ones with lower `badness`.
		return b1 > b2
	}
	a1, a2 := ps[i].IsAlerting(), ps[j].IsAlerting()
	if a1 != a2 {
		// Alerting probes sort before (lower value than) non-alerting ones.
		return a1
	}
	l1, l2 := ps[i].getLastAlert(), ps[j].getLastAlert()
	if !l1.Equal(l2) {
		// Probes that alerted longer ago sort after ones that alerted
		// more recently.
		return l1.After(l2)
	}
	r1, r2 := ps[i].Records(), ps[j].Records()
	if len(r1) != len(r2) {
		// Probes with shorter history sort after those with longer
		// history.
		return len(r1) > len(r2)
	}
	// Tie-breaker: Sort by name.
	if ps[i].Name != ps[j].Name {
		return ps[i].Name < ps[j].Name
	}
	// Tie-breaker #2: Sort by desc.
	if ps[i].Desc != ps[j].Desc {
		return ps[i].Desc < ps[j].Desc
	}
	// We have no way of comparing.
	return true
}
func (ps Probes) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

func init() {
	flag.Var(&disabledProbes, "disabled_probes", "comma-separated list of probes to disable")
	flag.Var(&onlyProbes, "only_probes", "comma-separated list of the only probes to enable")
}
