package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beldur/kraken-go-api-client"
)

func main() {
	key := os.Getenv("KRAKEN_KEY")
	secret := os.Getenv("KRAKEN_SECRET")
	log.Printf("Using %d char key, %d char secret\n", len(key), len(secret))
	api := krakenapi.New(key, secret)

	// There are also some strongly typed methods available
	ticker, err := api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		log.Fatal(err)
	}

	balance, err := api.Balance()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Opening BTC price: %+v\n", ticker.XXBTZEUR.OpeningPrice)
	fmt.Printf("Balances: %+v\n", balance)
}
