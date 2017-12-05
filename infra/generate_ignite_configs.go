// ignite.go generates Ignite JSON configs.
//
// TODO: could version the systemd units as well.
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
	log.Printf("Read %d character secret service hash.\n", len(sshash))
	conf, err := ignite.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to create project configs: %v\n", err)
	}

	ns, err := conf.CreateNodes()
	if err != nil {
		log.Fatalf("Failed to create nodes: %v\n", err)
	}
	for _, n := range ns {
		log.Printf("Writing Ignition config for %v..\n", n)
		err := n.Write()
		if err != nil {
			log.Fatalf("Failed to write node config: %v\n", err)
		}
	}
}
