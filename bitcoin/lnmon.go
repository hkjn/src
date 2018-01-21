// lnmon.go is a simple wrapper around c-lightning's lightning-cli for monitoring state.
package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

type cli struct{}

func (c cli) exec(cmd string) (string, error) {
	cmd := exec.Command("lightning-cli", cmd)
	out := bytes.Buffer{}
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (c cli) GetInfo() (string, error) {
	return c.exec("getinfo")
}

func (c cli) GetPeers() (string, error) {
}

func (c cli) GetNodes() (string, error) {
}

//{ "id" : "02501223d013c9c8da04fa7d482570cb38416f6faade7bd86f03413915557e0ff7", "port" : 19735, "address" :
//       [   { "type" : "ipv4", "address" : "163.172.162.18", "port" : 9735 } ], "version" : "v0.5.2-2016-11-21-1599-gf298c6b0", "blockheight" : 1259783 }

func main() {
	log.Printf("Running lightning-cli getinfo..\n")
	c := cli{}
	info, err := c.GetInfo()
	if err != nil {
		log.Fatalf("Failed to get info: %v\n", err)
	}
	// lightning-cli getinfo
	// lightning-cli getpeers
	// lightning-cli getnodes
	// lightning-cli getchannels

	// lightning-cli listfunds
	// lightning-cli listinvoice
	// lightning-cli listpayments
}
