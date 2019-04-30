// Logging helpers
package autosite

import (
	"net/http"

	"github.com/golang/glog"
)

// LoggerFunc returns a logger from a http request.
type LoggerFunc func(*http.Request) Logger

// Logger specifies logging functions.
//
// The methods are chosen to match the logging methods from
// appengine.Context, without needing to depend on appengine.
type Logger interface {
	// Debugf formats its arguments according to the format, analogous to fmt.Printf,
	// and records the text as a log message at Debug level.
	Debugf(format string, args ...interface{})

	// Infof is like Debugf, but at Info level.
	Infof(format string, args ...interface{})

	// Warningf is like Debugf, but at Warning level.
	Warningf(format string, args ...interface{})

	// Errorf is like Debugf, but at Error level.
	Errorf(format string, args ...interface{})

	// Criticalf is like Debugf, but at Critical level.
	Criticalf(format string, args ...interface{})
}

// Glogger implements Logger using package glog.
//
// Note that Glogger should not be used on appengine, since attempting
// to write to disk causes a panic.
type Glogger struct{}

// Debugf formats its arguments according to the format, analogous to fmt.Printf,
// and records the text as a log message at Debug level.
func (Glogger) Debugf(format string, args ...interface{}) {
	glog.V(1).Infof(format, args...)
}

// Infof is like Debugf, but at Info level.
func (Glogger) Infof(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

// Warningf is like Debugf, but at Warning level.
func (Glogger) Warningf(format string, args ...interface{}) {
	glog.Warningf(format, args...)
}

// Errorf is like Debugf, but at Error level.
func (Glogger) Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args...)
}

// Criticalf is like Debugf, but at Critical level.
func (Glogger) Criticalf(format string, args ...interface{}) {
	glog.Fatalf(format, args...)
}
