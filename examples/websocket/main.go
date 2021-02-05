package main

import (
	"context"
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

	if err := subscribeToTickers(ctx, client, symbols...); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err := subscribeToMarkets(ctx, client); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err := subscribeToTrades(ctx, client, symbols...); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err := subscribeToOrderBooks(ctx, client, symbols...); err != nil {
		log.Println(err)
		os.Exit(1)
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
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func subscribeToTickers(ctx context.Context, client *api.Client, symbols ...string) (err error) {

	wssub, err := client.Stream.SubscribeToTickers(ctx, symbols...)
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-wssub.EventC:
				if ok {
					ticker, ok := event.(*models.TickerResponse)
					if ok && ticker != nil {
						log.Printf("%+v\n", *ticker)
					}
				}
			}
		}
	}()

	return
}

func subscribeToMarkets(ctx context.Context, client *api.Client) (err error) {

	wssub, err := client.Stream.SubscribeToMarkets(ctx)
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-wssub.EventC:
				if ok {
					markets, err := api.MapToMarketData(event)
					if err != nil {
						return
					}
					for c, m := range markets {
						log.Printf("%s: %+v\n", c, *m)
					}
				}
			}
		}
	}()

	return
}

func subscribeToTrades(ctx context.Context, client *api.Client, symbols ...string) (err error) {

	wssub, err := client.Stream.SubscribeToTrades(ctx, symbols...)
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-wssub.EventC:
				if ok {
					trades, ok := event.(*models.TradesResponse)
					if ok && trades != nil {
						for _, t := range trades.Trades {
							log.Printf("Trade: %+v\n", t)
						}
					}
				}
			}
		}
	}()

	return
}

func subscribeToOrderBooks(ctx context.Context, client *api.Client, symbols ...string) (err error) {

	wssub, err := client.Stream.SubscribeToOrderBooks(ctx, symbols...)
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-wssub.EventC:
				if ok {
					book, ok := event.(*models.OrderBookResponse)
					if ok && book != nil {
						log.Printf("Book Response %s: %+v\n", book.Symbol, *book)
					}
				}
			}
		}
	}()

	return
}
