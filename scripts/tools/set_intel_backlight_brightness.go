// set_intel_backlight_brightness allows light levels to be set.

package main

import (
	"log"
	"os"
	"strconv"

	"github.com/hkjn/brightness"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v [brightness between 0-852]\n", os.Args[0])
	}
	lvl, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid value: %v\n", lvl)
	}
	err = brightness.Set(lvl)
	if err != nil {
		log.Fatalf("Failed to set value: %v\n", err)
	}
}
