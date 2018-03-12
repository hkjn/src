// gomon is a web tool that handles monitoring and alerting.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"

	"hkjn.me/src/dashboard"
)

func main() {
	flag.Parse()
	var conf dashboard.Config
	err := envconfig.Process("dashboard", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	if conf.BindAddr == "" {
		conf.BindAddr = ":8080"
	}
	fmt.Printf("gomon initializing, listening on %s..\n", conf.BindAddr)

	err = http.ListenAndServe(
		conf.BindAddr,
		dashboard.Start(conf),
	)
	log.Fatal(err.Error())
}
