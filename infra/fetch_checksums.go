// fetch_checksums.go is a tool that generates .sha512 files under checksums/ based on config.json.
package main

import (
	"log"

	"hkjn.me/src/infra/ignite"
	"hkjn.me/src/infra/secretservice"
)

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
	if err := conf.DownloadChecksums("checksums", sshash, secretservice.BaseDomain); err != nil {
		log.Fatalf("Failed to download checksums: %v\n", err)
	}
}
