package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/beldur/kraken-go-api-client"
)

type (
	balance struct {
		ticker string
		value  float32
	}
)

// descBalance returns list of strings describing the balances.
//
// Note that there's many currencies and altcoins we don't check.
func descBalance(r *krakenapi.BalanceResponse) []balance {
	result := []balance{}
	if r.XXBT > 0.00001 {
		result = append(result, balance{
			ticker: "BTC",
			value:  r.XXBT,
		})
	}
	if r.ZEUR > 0.00001 {
		result = append(result, balance{
			ticker: "EUR",
			value:  r.ZEUR,
		})
	}
	if r.BCH > 0.00001 {
		result = append(result, balance{
			ticker: "BCH",
			value:  r.BCH,
		})
	}
	if r.XXMR > 0.00001 {
		result = append(result, balance{
			ticker: "XMR",
			value:  r.XXMR,
		})
	}
	return result
}

func query(api *krakenapi.KrakenApi) (ticker *krakenapi.TickerResponse, balance *krakenapi.BalanceResponse, err error) {
	maxRetries := 5
	prefix := ""
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			prefix = fmt.Sprintf("[retry %d] ", i)
		}
		fmt.Printf("%sFetching ticker for BTC / EUR..\n", prefix)
		ticker, err = api.Ticker(krakenapi.XXBTZEUR)
		if err == nil {
			break
		}
		fmt.Printf("Failed to fetch ticker, retrying: %v\n", err)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch ticker too many times: %v", err)
	}
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			prefix = fmt.Sprintf("[retry %d]", i)
		}
		fmt.Printf("%sFetching all account balances..\n", prefix)
		balance, err = api.Balance()
		if err == nil {
			break
		}
		fmt.Printf("Failed to fetch balance, retrying: %v\n", err)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch account balances too many times: %v", err)
	}
	return ticker, balance, nil
}

func main() {
	key := os.Getenv("KRAKEN_KEY")
	secret := os.Getenv("KRAKEN_SECRET")
	if len(key) == 0 {
		log.Fatalf("No KRAKEN_KEY specified.")
	}
	if len(secret) == 0 {
		log.Fatalf("No KRAKEN_SECRET specified.")
	}
	api := krakenapi.New(key, secret)

	ticker, balance, err := query(api)
	if err != nil {
		log.Fatalf("Failed to query Kraken: %v\n", err)
	}
	fmt.Printf("Opening BTC price: %+v\n", ticker.XXBTZEUR.OpeningPrice)

	fmt.Println("Balances:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	for _, b := range descBalance(balance) {
		fmt.Fprintf(w, "%f\t %s\n", b.value, b.ticker)
	}
	w.Flush()
}
