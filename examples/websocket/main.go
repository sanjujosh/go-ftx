package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

var symbols []string = []string{"BTC-PERP", "BTC/USD"}

func main() {

	done, sigs := make(chan struct{}), make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := api.New()
	client.Stream.SetDebugMode(true)

	wstickers, err := client.Stream.SubscribeToTickers(ctx, symbols...)
	if err != nil {
		log.Fatalln(err)
	}

	wsmarkets, err := client.Stream.SubscribeToMarkets(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	wstrades, err := client.Stream.SubscribeToTrades(ctx, symbols...)
	if err != nil {
		log.Fatalln(err)
	}

	wsbooks, err := client.Stream.SubscribeToOrderBooks(ctx, symbols...)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		time.Sleep(15 * time.Second)
		done <- struct{}{}
	}()

	for {

		select {

		case <-done:
			log.Println("Exiting")
			return

		case <-sigs:
			log.Println("Exiting")
			return

		case event, ok := <-wstickers.EventC:
			if ok {
				ticker, ok := event.(*models.TickerResponse)
				if ok && ticker != nil {
					fmt.Printf("Ticker %+v\n", *ticker)
				}
			}

		case event, ok := <-wsmarkets.EventC:
			if ok {
				markets, err := api.MapToMarketData(event)
				if err == nil {
					fmt.Println("Market")
					for c, m := range markets {
						fmt.Printf("%s: %+v\n", c, *m)
					}
				}
			}

		case event, ok := <-wstrades.EventC:
			if ok {
				trades, ok := event.(*models.TradesResponse)
				if ok && trades != nil {
					fmt.Println("Trades")
					for _, t := range trades.Trades {
						fmt.Printf("Trade: %+v\n", t)
					}
				}
			}

		case event, ok := <-wsbooks.EventC:
			if ok {
				book, ok := event.(*models.OrderBookResponse)
				if ok && book != nil {
					fmt.Printf("Book Response %s: %+v\n", book.Symbol, *book)
				}
			}

		default:
			time.Sleep(time.Millisecond)
		}
	}
}
