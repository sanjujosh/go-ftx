# go-ftx
FTX exchange golang library

Forked from https://github.com/uscott/go-ftx

This is a full implementation of the FTX REST and Websocket API.

Use at your own risk.

### Install
```shell script
go get github.com/uscott/go-ftx
```

### Usage

> See examples directory and test cases for more examples

#### REST
```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sanjujosh/go-ftx/api"
	"github.com/sanjujosh/go-ftx/models"
)

func main() {
	client := api.New(
		api.WithAuth("API-KEY", "API-SECRET"),
		api.WithHTTPClient(&http.Client{
			Timeout: 5 * time.Second,
		}),
	)

	info := models.AccountInformation{}
	err := client.Account.GetAccountInformation(&info)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account info: %+v\n", info)
}
```

#### WebSocket

Refer to examples/websocket/websocket.go

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sanjujosh/go-clog"
	"github.com/sanjujosh/go-ftx/api"
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
```

### Websocket Debug Mode

The client now uses package go-clog which is a minor extension of https://github.com/sirupsen/logrus for logging.

Debug messages or others can be logged via the logger.

```go
    client := ftx.New()
    client.Logger.SetLevel(clog.DebugLevel) // Exactly like logrus.DebugLevel
```

### No Logged In Error
"Not logged in" errors usually come from wrong signatures.

FTX released an article on how to authenticate https://blog.ftx.com/blog/api-authentication/

If you have unauthorized error to private methods, then you need to use SetServerTimeDiff()
```go
ftx := New()
ftx.SetServerTimeDiff()
```
