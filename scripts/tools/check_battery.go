// Simple tool to show battery status from SysFS.
package main

import (
	"log"
	"hkjn.me/power"
)

func main() {
	bat, err := power.Get()
	if err != nil {
		log.Fatalf("failed to get battery info: %v\n", err)
	}
	for i, b := range bat {
		log.Printf("[Battery %d]: %+v\n", i, b)
	}
}
