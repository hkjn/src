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
	var logger klog.Logger
	logger = klog.NewLogfmtLogger(klog.NewSyncWriter(os.Stderr))
	logger = klog.With(logger, "time", klog.DefaultTimestampUTC, "caller", klog.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte("API secret"),
	}
	ctx, _ := context.WithCancel(context.Background())
	// use second return value for cancelling request when shutting down the app

	api := binance.NewAPIService(
		"https://www.binance.com",
		"API key",
		hmacSigner,
		logger,
		ctx,
	)
	b := binance.NewBinance(api)
	//obr := nil
	ob, err := b.OrderBook(binance.OrderBookRequest{Symbol: "ETHBTC"})
	if err != nil {
		log.Fatalf("Failed to fetch order book: %v\n", err)
	}
	for i, bid := range ob.Bids {
		fmt.Printf("bid %d (%T): %v\n", i, bid, bid)
	}
	// fmt.Printf("binance orders: %v\n", ob)
}
