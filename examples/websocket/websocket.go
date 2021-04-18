package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uscott/go-clog"
	"github.com/uscott/go-ftx/api"
)

var symbols []string = []string{"BTC-PERP", "BTC/USD"}

const N int = 5

func main() {

	done, sigs := make(chan struct{}), make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := api.New()
	client.Logger.SetLevel(clog.DebugLevel)

	tickersC, err := client.Stream.SubscribeToTickers(ctx, symbols...)
	if err != nil {
		log.Fatalln(err)
	}
	cnttickers := 0

	marketsC, err := client.Stream.SubscribeToMarkets(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	cntmarkets := 0

	tradesC, err := client.Stream.SubscribeToTrades(ctx, symbols...)
	if err != nil {
		log.Fatalln(err)
	}
	cnttrades := 0

	booksC, err := client.Stream.SubscribeToOrderBooks(ctx, symbols...)
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
			log.Println("Exiting - time limit")
			return

		case <-sigs:
			log.Println("Exiting")
			return

		case ticker := <-tickersC:
			if ticker != nil && cnttickers < N {
				fmt.Printf("Ticker %+v\n", *ticker)
				cnttickers++
				if cnttickers == N {
					fmt.Println("Seen enough tickers")
				}
			}
		case market := <-marketsC:
			if market != nil && cntmarkets < N {
				fmt.Printf("Market: %+v\n", *market)
				cntmarkets++
				if cntmarkets == N {
					fmt.Println("Seen enough markets")
				}
			}
		case trade := <-tradesC:
			if trade != nil && cnttrades < N {
				fmt.Printf("Trade: %+v\n", *trade)
				cnttrades++
				if cnttrades == N {
					fmt.Println("Seen enough trades")
				}
			}
		case book := <-booksC:
			if book != nil && cntbooks < N {
				fmt.Printf("Book Response %s: %+v\n", book.Symbol, *book)
				cntbooks++
				if cntbooks == N {
					fmt.Println("Seen enough orderbooks")
				}
			}
		default:
			time.Sleep(time.Millisecond)
		}

		go func() {
			if cnttickers >= N && cntmarkets >= N && cnttrades >= N && cntbooks >= N {
				fmt.Println("Exiting - seen enough")
				os.Exit(0)
			}
			time.Sleep(time.Millisecond)
		}()
	}
}
