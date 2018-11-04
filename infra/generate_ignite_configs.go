// generate_ignite_config.go generates Ignite JSON configs.
//
// TODO: could version the systemd units as well.
package main

import (
	"log"
	"os/user"

	"hkjn.me/src/infra/ignite"
	"hkjn.me/src/infra/secretservice"
)

func main() {
	sshash, err := secretservice.GetHash()
	if err != nil {
		log.Fatalf("Unable to fetch secret service hash: %v\n", err)
	}
	log.Printf("Read %d character secret service hash.\n", len(sshash))

	ns, err := ignite.CreateNodes()
	if err != nil {
		log.Fatalf("Failed to create nodes: %v\n", err)
	}
	for _, n := range ns {
		log.Printf("Writing Ignition config for %v..\n", n)
		err := n.Write()
		if err != nil {
			u, uerr := user.Current()
			if uerr != nil {
				log.Fatalf("Failed to write node config, also failed to find current user (%v): %v\n", uerr, err)
			}
			log.Fatalf("Failed to write node config as user %v:%v: %v\n", u.Uid, u.Gid, err)
		}
	}
}
