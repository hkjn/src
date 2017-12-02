package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/beldur/kraken-go-api-client"
)

type (
	balance struct {
		ticker string
		value  float32
	}
	response struct {
		ticker       *krakenapi.TickerResponse
		balance      *krakenapi.BalanceResponse
		openOrders   *krakenapi.OpenOrdersResponse
		closedOrders *krakenapi.ClosedOrdersResponse
	}
)

// descBalance returns list of strings describing the balances.
//
// Note that there's many currencies and altcoins we don't check.
func descBalance(r *krakenapi.BalanceResponse) []balance {
	result := []balance{}
	if r.XXBT > 0.00001 {
		result = append(result, balance{
			ticker: "mBTC",
			value:  r.XXBT * 1000,
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

func query(api *krakenapi.KrakenApi) (resp *response, err error) {
	maxRetries := 10
	prefix := ""
	resp = &response{}
	for i := 0; i < maxRetries; i++ {
		var ticker *krakenapi.TickerResponse
		if i > 0 {
			prefix = fmt.Sprintf("[retry %d] ", i)
		}
		fmt.Printf("%sFetching ticker for BTC / EUR..\n", prefix)
		ticker, err = api.Ticker(krakenapi.XXBTZEUR)
		if err == nil {
			resp.ticker = ticker
			break
		}
		fmt.Printf("Failed to fetch ticker, retrying: %v\n", err)
		time.Sleep(time.Second * time.Duration(i))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ticker too many times: %v", err)
	}

	for i := 0; i < maxRetries; i++ {
		var balance *krakenapi.BalanceResponse
		if i > 0 {
			prefix = fmt.Sprintf("[retry %d] ", i)
		}
		fmt.Printf("%sFetching all account balances..\n", prefix)
		balance, err = api.Balance()
		if err == nil {
			resp.balance = balance
			break
		}
		fmt.Printf("Failed to fetch balance, retrying: %v\n", err)
		time.Sleep(time.Second * time.Duration(i))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account balances too many times: %v", err)
	}

	for i := 0; i < maxRetries; i++ {
		var openOrders *krakenapi.OpenOrdersResponse
		if i > 0 {
			prefix = fmt.Sprintf("[retry %d] ", i)
		}
		fmt.Printf("%sFetching open orders..\n", prefix)
		openOrders, err = api.OpenOrders(map[string]string{})
		if err == nil {
			resp.openOrders = openOrders
			break
		}
		fmt.Printf("Failed to fetch open orders, retrying: %v\n", err)
		time.Sleep(time.Second * time.Duration(i))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch open orders too many times: %v", err)
	}

	for i := 0; i < maxRetries; i++ {
		var closedOrders *krakenapi.ClosedOrdersResponse
		if i > 0 {
			prefix = fmt.Sprintf("[retry %d] ", i)
		}
		fmt.Printf("%sFetching closed orders..\n", prefix)
		closedOrders, err = api.ClosedOrders(map[string]string{})
		if err == nil {
			resp.closedOrders = closedOrders
			break
		}
		fmt.Printf("Failed to fetch closed orders, retrying: %v\n", err)
		time.Sleep(time.Second * time.Duration(i))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch closed orders too many times: %v", err)
	}

	return resp, nil
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

	resp, err := query(api)
	if err != nil {
		log.Fatalf("Failed to query Kraken: %v\n", err)
	}
	fmt.Printf("Opening BTC price: %+v\n", resp.ticker.XXBTZEUR.OpeningPrice)

	fmt.Println("Balances:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	for _, b := range descBalance(resp.balance) {
		fmt.Fprintf(w, "%f\t %s\n", b.value, b.ticker)
	}
	w.Flush()

	i := 0
	orders := make([]krakenapi.Order, len(resp.openOrders.Open), len(resp.openOrders.Open))
	for _, o := range resp.openOrders.Open {
		orders[i] = o
		i += 1
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OpenT.After(orders[j].OpenT)
	})
	fmt.Println("")
	if len(resp.openOrders.Open) == 0 {
		fmt.Println("No open orders.")
	} else {
		fmt.Println("Open orders:")
	}
	for i, o := range orders {
		if i >= 10 {
			break
		}
		fmt.Printf("%v: %s of %s, %s %v order @ %s BTC (or whatever)\n",
			o.OpenT, o.Description.Type, o.Description.AssetPair, o.Status, o.Description.OrderType, o.Description.PrimaryPrice)
	}

	fmt.Println("")
	if len(resp.closedOrders.Closed) == 0 {
		fmt.Println("No closed orders.")
	} else {
		fmt.Println("Closed orders:")
	}

	i = 0
	orders = make([]krakenapi.Order, len(resp.closedOrders.Closed), len(resp.closedOrders.Closed))
	for _, o := range resp.closedOrders.Closed {
		orders[i] = o
		i += 1
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CloseT.After(orders[j].CloseT)
	})
	for i, o := range orders {
		if i >= 5 {
			break
		}
		fmt.Printf("%v: %s of %s, status %s for reason %s, cost %v EUR (or whatever), volume executed %v BTC (or whatever)\n",
			o.CloseT, o.Description.Type, o.Description.AssetPair, o.Status, o.Reason, o.Cost, o.VolumeExecuted)
	}
}
