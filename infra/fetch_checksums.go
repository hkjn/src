// fetch_checksums.go is a tool that generates .sha512 files under checksums/ based on config.json.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"hkjn.me/src/infra/ignite"
	"hkjn.me/src/infra/secretservice"
)

// checkClose closes specified closer and sets err to the result.
func checkClose(c io.Closer, err *error) {
	cerr := c.Close()
	if *err == nil {
		*err = cerr
	}
}

func main() {
	sshash, err := secretservice.GetHash()
	if err != nil {
		log.Fatalf("Unable to fetch secret service hash: %v\n", err)
	}

	conf, err := ignite.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read node config: %v\n", err)
	}

	log.Printf("Read %d node configs..\n", len(conf.NodeConfigs))
	checksums, err := conf.GetChecksums(sshash)
	if err != nil {
		log.Fatalf("Failed to download checksums: %v\n", err)
	}
	basedir := "checksums"
	for pv, checksumlines := range checksums {
		var err error
		filename := filepath.Join(basedir, fmt.Sprintf("%s_%s.sha512", pv.Name, pv.Version))
		log.Printf("Creating %s with %d checksums for %v\n", filename, len(checksumlines), pv)
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Failed to open checksums file: %v\n", err)
		}
		defer checkClose(f, &err)
		for _, line := range checksumlines {
			fmt.Printf("Appending checksum %x[..] to %q..\n", line[:5], filename)
			f.Write([]byte(line))
		}
	}
}
