// Package ignite deals with Ignite JSON configs.
//
// TODO: Fix issue with missing /etc/ssl .pem files in hkjninfra project, seems like call to ProjectConfigs.GetSecrets() got lost..
package ignite

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"sort"
	"strings"
)

type (
	fileVerification struct {
		Hash string `json:"hash,omitempty"`
	}
	fileContents struct {
		Source       string           `json:"source"`
		Verification fileVerification `json:"verification"`
	}
	file struct {
		Filesystem string            `json:"filesystem"`
		Path       string            `json:"path"`
		Contents   fileContents      `json:"contents"`
		Mode       int               `json:"mode"`
		User       map[string]string `json:"user"`
		Group      map[string]string `json:"group"`
	}
	storage struct {
		Filesystem []string `json:"filesystem"`
		Files      []file   `json:"files"`
	}
	systemdDropin struct {
		Name     string `json:"name"`
		Contents string `json:"contents"`
	}
	systemdUnit struct {
		Enable   bool            `json:"enable"`
		Name     string          `json:"name"`
		Contents string          `json:"contents,omitempty"`
		Dropins  []systemdDropin `json:"dropins,omitempty"`
	}
	systemd struct {
		Units    []systemdUnit     `json:"units"`
		Passwd   map[string]string `json:"passwd"`
		Networkd map[string]string `json:"networkd"`
	}
	ignition struct {
		Version string            `json:"version"`
		Config  map[string]string `json:"config"`
	}
	ignitionConfig struct {
		Ignition ignition `json:"ignition"`
		Storage  storage  `json:"storage"`
		Systemd  systemd  `json:"systemd"`
	}
	// binary to fetch on a node
	binary struct {
		// url to fetch binary from, e.g. "https://github.com/hkjn/hkjninfra/releases/download/1.1.7/tserver_x86_64"
		url string
		// checksum of the file, e.g. "sha512-123[...]"
		checksum string
		// path on the remote node for the binary, e.g. "/opt/bin/tserver"
		path string
	}
	Version string
	// nodeName is the name of a node, e.g. "core".
	nodeName string
	// node is a single instance.
	node struct {
		// name is the name of the node.
		name nodeName
		// binaries are the files to install on the node.
		binaries []binary
		// systemdUnits are the systemd units to use for the node.
		systemdUnits []systemdUnit
	}
	nodes map[nodeName]node
	// ProjectName is the name of a project.
	ProjectName string
	checksums   map[string]string
	// ProjectVersion is the name and version of a project.
	ProjectVersion struct {
		// name is the name of a project the node should run node, e.g. "hkjninfra".
		Name ProjectName `json:"name"`
		// version is the version of the project that should run on the node, e.g. "1.0.1".
		Version Version `json:"version"`
	}
	// NodeConfig is the configuration of a single node
	NodeConfig struct {
		// sshash is the secretservice hash to use
		sshash string
		// ProjectVersions is the names of all the projects the node should run
		ProjectVersions []ProjectVersion `json:"projects"`
		// checksums holds each version of each project's checksums.
		checksums map[ProjectVersion]checksums
		// arch is the CPU architecture the node runs, e.g. "x86_64"
		Arch string `json:"arch"`
	}

	NodeFile struct {
		Path        string `json:"path"`
		Name        string `json:"name"`
		ChecksumKey string `json:"checksum_key"`
	}
	NodeFiles  []NodeFile
	Secret     NodeFile
	Secrets    []Secret
	DropinName struct {
		Unit, Dropin string
	}
	// TODO: Unify with checksums type above
	checksumlines []string
	Checksums     map[ProjectVersion]checksumlines

	// projectConfig is the full configuration for a project.
	projectConfig struct {
		Units   []string     `json:"units"`
		Dropins []DropinName `json:"dropins"`
		Files   NodeFiles    `json:"files"`
		Secrets NodeFiles    `json:"secrets"`
	}
	// ProjectConfigs represents all the project configurations.
	ProjectConfigs map[ProjectName]projectConfig
	// NodeConfigs is the configuration of all nodes.
	NodeConfigs map[nodeName]NodeConfig
	Config      struct {
		ProjectConfigs ProjectConfigs `json:"project_configs"`
		NodeConfigs    NodeConfigs    `json:"nodes"`
	}
)

// sharedFiles are the shared files for each node.
var sharedFiles = []file{
	{
		Filesystem: "root",
		Path:       "/etc/coreos/update.conf",
		Contents: fileContents{
			Source:       "data:,GROUP%3Dbeta%0AREBOOT_STRATEGY%3D%22etcd-lock%22",
			Verification: fileVerification{},
		},
		Mode:  420,
		User:  map[string]string{},
		Group: map[string]string{},
	},
}

func (b binary) toFile() file {
	return file{
		Filesystem: "root",
		Path:       b.path,
		Contents: fileContents{
			Source: b.url,
			Verification: fileVerification{
				Hash: fmt.Sprintf("sha512-%s", b.checksum),
			},
		},
		Mode:  493,
		User:  map[string]string{},
		Group: map[string]string{},
	}
}

// newSystemdUnit reads systemd unit from file name.
func newSystemdUnit(unitFile string) (*systemdUnit, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf(
		"units/%s",
		unitFile,
	))
	if err != nil {
		return nil, err
	}
	return &systemdUnit{
		Enable:   true,
		Name:     unitFile,
		Contents: string(b),
	}, nil
}

// Load returns the systemd units.
func (dn DropinName) Load() (*systemdUnit, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("units/%s", dn.Dropin))
	if err != nil {
		return nil, err
	}
	return &systemdUnit{
		Name: dn.Unit,
		Dropins: []systemdDropin{
			{
				Name:     dn.Dropin,
				Contents: string(b),
			},
		},
	}, nil
}

// getFiles returns the files to put on the node.
func (n node) getFiles() []file {
	result := make(
		[]file,
		len(n.binaries)+len(sharedFiles),
		len(n.binaries)+len(sharedFiles),
	)
	for i, sharedFile := range sharedFiles {
		result[i] = sharedFile
	}
	for j, bin := range n.binaries {
		result[j+len(sharedFiles)] = bin.toFile()
	}
	return result
}

// String returns a human-readable description of the node.
func (n node) String() string {
	return fmt.Sprintf(
		"%q (%d binaries, %d systemd units)",
		n.name,
		len(n.binaries),
		len(n.systemdUnits),
	)
}

// Write writes the Ignition config to disk.
func (n node) Write() error {
	bp := "bootstrap"
	_, err := os.Stat(bp)
	if os.IsNotExist(err) {
		if mkerr := os.Mkdir(bp, 755); mkerr != nil {
			u, _ := user.Current()
			return fmt.Errorf("failed to create dir %q as %s:%s: %v", bp, u.Uid, u.Gid, mkerr)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat %q: %v", bp, err)
	}
	f, err := os.Create(fmt.Sprintf("%s/%s.json", bp, n.name))
	if err != nil {
		return err
	}
	defer f.Close()

	conf := n.getIgnitionConfig()
	return json.NewEncoder(f).Encode(&conf)
}

// getIgnitionConfig returns the ignition config for the nod.
func (n node) getIgnitionConfig() ignitionConfig {
	return ignitionConfig{
		Ignition: ignition{
			Version: "2.0.0",
			Config:  map[string]string{},
		},
		Storage: storage{
			Filesystem: []string{},
			Files:      n.getFiles(),
		},
		Systemd: systemd{
			Units:    n.systemdUnits,
			Passwd:   map[string]string{},
			Networkd: map[string]string{},
		},
	}
}

// newProject returns the systemd units created from config.
func (conf projectConfig) getSystemdUnits() ([]systemdUnit, error) {
	units := []systemdUnit{}
	for _, unitFile := range conf.Units {
		unit, err := newSystemdUnit(unitFile)
		if err != nil {
			return nil, err
		}
		units = append(units, *unit)
	}
	for _, d := range conf.Dropins {
		dropin, err := d.Load()
		if err != nil {
			return nil, err
		}
		units = append(units, *dropin)
	}
	return units, nil
}

// getChecksums returns the checksums for the project version.
func (pv ProjectVersion) getChecksums() (checksums, error) {
	checksumFile := fmt.Sprintf(
		"checksums/%s_%s.sha512",
		pv.Name,
		pv.Version,
	)
	checksumData, err := ioutil.ReadFile(checksumFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read checksums for %q version %q: %v", pv.Name, pv.Version, err)
	}
	result := checksums{}
	for _, line := range strings.Split(string(checksumData), "\n") {
		if len(line) == 0 {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line in checksum file %s: %q", checksumFile, line)
		}
		result[parts[1]] = parts[0]
	}
	return result, nil
}

// GetChecksumURL returns the URL to fetch the checksums for the project.
func (pv ProjectVersion) GetChecksumURL() string {
	return fmt.Sprintf(
		"https://github.com/hkjn/%s/releases/download/%s/SHA512SUMS",
		pv.Name,
		pv.Version,
	)
}

// Names returns the names of the project configs in sorted order.
func (conf ProjectConfigs) Names() []ProjectName {
	names := make([]ProjectName, len(conf), len(conf))
	i := 0
	for name, _ := range conf {
		names[i] = name
		i += 1
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

// String returns a human-readable description of the project configs.
func (conf ProjectConfigs) String() string {
	desc := make([]string, len(conf), len(conf))
	for i, name := range conf.Names() {
		desc[i] = fmt.Sprintf(
			"%s: %s",
			name,
			conf[name],
		)
	}
	return fmt.Sprintf(
		"ProjectConfigs{%s}",
		strings.Join(desc, ", "),
	)
}

// GetSecrets returns any secrets in the project.
func (conf ProjectConfigs) GetSecrets(projectName ProjectName) (Secrets, error) {
	p, exists := conf[projectName]
	if !exists {
		return nil, fmt.Errorf("no project %q", projectName)
	}
	result := make(Secrets, len(p.Secrets), len(p.Secrets))
	for i, s := range p.Secrets {
		result[i] = Secret(s)
	}
	return result, nil
}

// String returns a human-readable description of the NodeFiles.
func (nf NodeFiles) String() string {
	if len(nf) == 0 {
		return "[]"
	}
	files := make([]string, len(nf), len(nf))
	for i, f := range nf {
		if f.ChecksumKey != "" {
			files[i] = fmt.Sprintf(
				"NodeFile{Name: %s, ChecksumKey: %s, Path: %s}}",
				f.Name,
				f.ChecksumKey,
				f.Path,
			)
		} else {
			files[i] = fmt.Sprintf(
				"NodeFile{Name: %s, Path: %s}}",
				f.Name,
				f.Path,
			)
		}
	}
	return fmt.Sprintf(
		"NodeFiles{%s}",
		strings.Join(files, ", "),
	)
}

// GetURL returns the URL to fetch the secret.
func (s Secret) GetURL(secretServiceDomain, sshash string, pv ProjectVersion) string {
	return fmt.Sprintf(
		"https://%s/%s/files/%s/%s/certs/%s",
		secretServiceDomain,
		sshash,
		pv.Name,
		pv.Version,
		s.Name,
	)
}

// String returns a human-readable description of the NodeConfig.
func (nc NodeConfig) String() string {
	return fmt.Sprintf(
		"NodeConfig{Arch: %s}",
		nc.Arch,
	)
}

// String returns a human-readable description of the project.
func (conf projectConfig) String() string {
	if len(conf.Secrets) > 0 {
		return fmt.Sprintf(
			"project{Units: %s, Secrets: %s}",
			strings.Join(conf.Units, ", "),
			conf.Secrets.String(),
		)
	} else {
		return fmt.Sprintf(
			"project{Units: %s}",
			strings.Join(conf.Units, ", "),
		)
	}
}

// getBinaries returns the binaries for this project and version, given configs.
func (pv ProjectVersion) getBinaries(conf ProjectConfigs, checksums checksums) ([]binary, error) {
	pc, exists := conf[pv.Name]
	if !exists {
		return nil, fmt.Errorf("bug: no such project %q", pv.Name)
	}
	result := []binary{}
	for _, file := range pc.Files {
		key := file.ChecksumKey
		if key == "" {
			key = file.Name
		}
		checksum, exists := checksums[key]
		if !exists {
			return nil, fmt.Errorf("missing checksum for key %q; all checksums \"%v\"", key, checksums)
		}
		result = append(result, binary{
			url: fmt.Sprintf(
				"https://github.com/hkjn/%s/releases/download/%s/%s",
				pv.Name,
				pv.Version,
				file.Name,
			),
			checksum: checksum,
			path:     file.Path,
		})
	}

	return result, nil
}

// getUnits returns the systemd units for the specific projects.
func (conf ProjectConfigs) getSystemdUnits(pversions []ProjectVersion) ([]systemdUnit, error) {
	result := []systemdUnit{}
	for _, pv := range pversions {
		pc, exists := conf[pv.Name]
		if !exists {
			return nil, fmt.Errorf("bug: no such project %q", pv.Name)
		}
		units, err := pc.getSystemdUnits()
		if err != nil {
			return nil, err
		}
		result = append(result, units...)
	}
	return result, nil
}

// String returns a human-readable description of the config.
func (conf Config) String() string {
	return fmt.Sprintf(
		"Config{%s, %s}",
		conf.ProjectConfigs,
		conf.NodeConfigs,
	)
}

// createNodes returns the nodes created from configs.
func (nconf NodeConfig) createNode(name nodeName, pconf ProjectConfigs) (*node, error) {
	bins := []binary{}
	for _, pv := range nconf.ProjectVersions {
		newbins, err := pv.getBinaries(pconf, nconf.checksums[pv])
		if err != nil {
			return nil, err
		}
		bins = append(bins, newbins...)
	}
	units, err := pconf.getSystemdUnits(nconf.ProjectVersions)
	if err != nil {
		return nil, err
	}
	return &node{
		name:         name,
		binaries:     bins,
		systemdUnits: units,
	}, nil
}

// getNodes returns the nodes created from the config.
func (conf Config) getNodes() (nodes, error) {
	result := nodes{}
	for name, nc := range conf.NodeConfigs {
		log.Printf("Generating config for node %q..\n", name)
		n, err := nc.createNode(name, conf.ProjectConfigs)
		if err != nil {
			return nil, err
		}
		result[name] = *n
		log.Printf("Generated config %v\n", result[name])
	}
	return result, nil
}

// checkClose closes specified closer and sets err to the result.
func checkClose(c io.Closer, err *error) {
	cerr := c.Close()
	if *err == nil {
		*err = cerr
	}
}

// checksum returns the sha512 checksum of file at specified url.
func (s Secret) checksum(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer checkClose(resp.Body, &err)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code from GET %q, want 200 OK, got %s", url, resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	digest := sha512.Sum512(b)
	return fmt.Sprintf("%x  %s\n", digest, s.Name), nil
}

// getSecretChecksums returns the checksums of secrets in given combination of node and project version.
func (nc NodeConfig) getSecretChecksums(sshash, ssbasedomain string, pconfs ProjectConfigs) (Checksums, error) {
	result := Checksums{}
	fetched := map[string]bool{}
	for _, pv := range nc.ProjectVersions {
		// TODO: Also need to handle secrets, like decenter.world.pem for "decenter.world"..
		// fetch from secret service directly?
		if pv.Name == ProjectName("bitcoin") {
			// TODO: Instead of special-casing "core" (bitcoin) project, which has
			// no checksums since there's no binaries to download, maybe start
			// checksumming / versioning systemd unit (.service, .mount) and
			// dropins (.conf) within the project?
			log.Printf("Skipping bitcoin, no binaries to download..\n")
			continue
		}
		secrets, err := pconfs.GetSecrets(pv.Name)
		if err != nil {
			return nil, err
		}
		for _, secret := range secrets {
			url := secret.GetURL(ssbasedomain, sshash, pv)
			if fetched[url] {
				continue
			}
			log.Printf("Fetching and checksumming secret %q..\n", secret.Name)
			line, err := secret.checksum(url)
			if err != nil {
				return nil, err
			}
			lines, exist := result[pv]
			if !exist {
				result[pv] = checksumlines([]string{line})
			} else {
				result[pv] = append(lines, line)
			}
			fetched[url] = true
		}
	}
	return result, nil
}

// GetChecksums returns the checksums specified by the config.
func (conf *Config) GetChecksums(sshash, ssbasedomain string) (Checksums, error) {
	result := Checksums{}
	// TODO: Should include non-secret checksums here too, or otherwise make sure that fetch_checksums.go will include those.
	for node, nc := range conf.NodeConfigs {
		log.Printf("Fetching checksums for node %q..\n", node)
		newchecksums, err := nc.getSecretChecksums(sshash, ssbasedomain, conf.ProjectConfigs)
		if err != nil {
			return nil, err
		}
		for k, v := range newchecksums {
			_, exists := result[k]
			if !exists {
				result[k] = v
			}
		}
	}
	return result, nil
}

// ReadConfig returns the node/project configs, read from disk.
func ReadConfig() (*Config, error) {
	conf := Config{}
	f, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		return nil, err
	}

	for nn, nc := range conf.NodeConfigs {
		nc := nc
		nc.checksums = map[ProjectVersion]checksums{}
		for _, pv := range nc.ProjectVersions {
			checksums, err := pv.getChecksums()
			if err != nil {
				return nil, err
			}
			nc.checksums[pv] = checksums
		}
		conf.NodeConfigs[nn] = nc
	}
	// TODO: Should probably include systemd units / files and secrets in *Config returned here..
	return &conf, nil
}

// CreateNodes returns nodes created from the configs.
func CreateNodes() (nodes, error) {
	conf, err := ReadConfig()
	if err != nil {
		return nil, err
	}
	log.Printf("Read config: %+v\n", conf)
	for _, nc := range conf.NodeConfigs {
		if len(nc.checksums) == 0 {
			return nil, fmt.Errorf("bug: missing checksums")
		}
	}
	result, err := conf.getNodes()
	if err != nil {
		log.Fatalf("Failed to get node: %v\n", err)
	}
	return result, nil
}
