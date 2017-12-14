// fetch_checksums.go is a tool that generates .sha512 files under checksums/ based on config.json.
package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

// checksumSecret downloads and checksums specified secret, and appends it to the checksum file.
func checksumSecret(url, checksumfile, secretfile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer checkClose(resp.Body, &err)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from GET %q, want 200 OK, got %s", url, resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	digest := sha512.Sum512(b)

	f, err := os.OpenFile(checksumfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer checkClose(f, &err)
	log.Printf("Appending checksum %x[..] to %q..\n", digest[:5], checksumfile)
	line := fmt.Sprintf("%x  %s\n", digest, secretfile)
	if _, err := f.Write([]byte(line)); err != nil {
		return err
	}
	_, err = io.Copy(f, resp.Body)
	return err
}

// fetchChecksums downloads specified URL of checksum file.
func fetchChecksums(url, filename string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer checkClose(resp.Body, &err)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from GET %q, want 200 OK, got %s", url, resp.Status)
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer checkClose(f, &err)
	log.Printf("Saving checksum file to %q..\n", filename)
	_, err = io.Copy(f, resp.Body)
	return err
}

// downloadChecksums downloads the checksum files.
func downloadChecksums(conf ignite.Config, sshash string) error {
	fetched := map[string]bool{}
	for node, nc := range conf.NodeConfigs {
		log.Printf("Fetching checksums for node %q..\n", node)
		for _, pv := range nc.ProjectVersions {
			// TODO: Also need to handle secrets, like decenter.world.pem for "decenter.world"..
			// fetch from secret service directly?
			if pv.Name == ignite.ProjectName("bitcoin") {
				// TODO: Instead of special-casing "core" (bitcoin) project, which has
				// no checksums since there's no binaries to download, maybe start
				// checksumming / versioning systemd unit (.service, .mount) and
				// dropins (.conf) within the project?
				log.Printf("Skipping bitcoin, no binaries to download..\n")
				continue
			}
			url := ignite.GetChecksumURL(pv)
			filename := fmt.Sprintf("checksums/%s_%s.sha512", pv.Name, pv.Version)
			if !fetched[url] {
				log.Printf("Fetching %q..\n", url)
				if err := fetchChecksums(url, filename); err != nil {
					return err
				}
				fetched[url] = true
			}
			secrets, err := conf.ProjectConfigs.GetSecrets(pv.Name)
			if err != nil {
				return err
			}
			for _, secret := range secrets {
				url := secret.GetURL(secretservice.BaseDomain, sshash, pv)
				if !fetched[url] {
					log.Printf("Fetching and checksumming secret %q..\n", secret.Name)
					if err := checksumSecret(url, filename, secret.Name); err != nil {
						return err
					}
					fetched[url] = true
				}
			}
		}
	}
	return nil
}

func main() {
	conf, err := ignite.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read node config: %v\n", err)
	}

	sshash, err := secretservice.GetHash()
	if err != nil {
		log.Fatalf("Unable to fetch secret service hash: %v\n", err)
	}
	log.Printf("Read %d node configs..\n", len(conf.NodeConfigs))
	if err := downloadChecksums(*conf, sshash); err != nil {
		log.Fatalf("Failed to download checksums: %v\n", err)
	}
}
