package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/binance-exchange/go-binance"
	klog "github.com/go-kit/kit/log"
)

func main() {
	key := os.Getenv("BINANCE_KEY")
	secret := os.Getenv("BINANCE_SECRET")
	if len(key) == 0 {
		log.Fatalf("No BINANCE_KEY specified.")
	}
	if len(secret) == 0 {
		log.Fatalf("No BINANCE_SECRET specified.")
	}
	var logger klog.Logger
	logger = klog.NewLogfmtLogger(klog.NewSyncWriter(os.Stderr))
	logger = klog.With(logger, "time", klog.DefaultTimestampUTC, "caller", klog.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(secret),
	}
	ctx, _ := context.WithCancel(context.Background())
	// use second return value for cancelling request when shutting down the app

	api := binance.NewAPIService(
		"https://www.binance.com",
		key,
		hmacSigner,
		logger,
		ctx,
	)
	b := binance.NewBinance(api)

	orders, err := b.AllOrders(binance.AllOrdersRequest{Symbol: "BATBTC"})

	//ob, err := b.OrderBook(binance.OrderBookRequest{Symbol: "ETHBTC"})
	if err != nil {
		log.Fatalf("Failed to fetch orders: %v\n", err)
	}
	//for i, bid := range ob.Bids {
	//	fmt.Printf("bid %d: %v\n", i, bid)
	//}
	fmt.Printf("binance orders: %v\n", orders)
}
