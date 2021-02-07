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

const N int = 5

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
	cnttickers := 0

	wsmarkets, err := client.Stream.SubscribeToMarkets(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	cntmarkets := 0

	wstrades, err := client.Stream.SubscribeToTrades(ctx, symbols...)
	if err != nil {
		log.Fatalln(err)
	}
	cnttrades := 0

	wsbooks, err := client.Stream.SubscribeToOrderBooks(ctx, symbols...)
	if err != nil {
		log.Fatalln(err)
	}
	cntbooks := 0

	go func() {
		time.Sleep(time.Minute)
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
				if ok && ticker != nil && cnttickers < N {
					fmt.Printf("Ticker %+v\n", *ticker)
					cnttickers++
					if cnttickers == N {
						fmt.Println("Seen enough tickers")
					}
				}
			}

		case event, ok := <-wsmarkets.EventC:
			if ok {
				markets, err := api.MapToMarketData(event)
				if err == nil && cntmarkets < N {
					fmt.Println("Market")
					for c, m := range markets {
						fmt.Printf("%s: %+v\n", c, *m)
					}
					cntmarkets++
					if cntmarkets == N {
						fmt.Println("Seen enough markets")
					}
				}
			}

		case event, ok := <-wstrades.EventC:
			if ok {
				trades, ok := event.(*models.TradesResponse)
				if ok && trades != nil && cnttrades < N {
					fmt.Println("Trades")
					for _, t := range trades.Trades {
						fmt.Printf("Trade: %+v\n", t)
					}
					cnttrades++
					if cnttrades == N {
						fmt.Println("Seen enough trades")
					}
				}
			}

		case event, ok := <-wsbooks.EventC:
			if ok {
				book, ok := event.(*models.OrderBookResponse)
				if ok && book != nil && cntbooks < N {
					fmt.Printf("Book Response %s: %+v\n", book.Symbol, *book)
					cntbooks++
					if cntbooks == N {
						fmt.Println("Seen enough orderbooks")
					}
				}
			}

		default:
			time.Sleep(time.Millisecond)
		}

		go func() {
			if cnttickers >= N && cntmarkets >= N && cnttrades >= N && cntbooks >= N {
				done <- struct{}{}
			}
			time.Sleep(time.Millisecond)
		}()

	}
}
