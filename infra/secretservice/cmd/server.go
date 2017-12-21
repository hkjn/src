// secretservice is a fileserver serving up secrets over HTTPS.
//
package main

import (
	"log"

	"hkjn.me/src/infra/secretservice"
)

func main() {
	log.Fatal(secretservice.Serve())
}
