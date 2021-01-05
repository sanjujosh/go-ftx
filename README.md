# ftx
FTX exchange golang library

Forked from https://github.com/grishinsana/ftx

### Install
```shell script
go get github.com/uscott/go-ftx
```

### Usage

> See examples directory and test cases for more examples

### TODO
- Private Streams (working on it)
- Orders (mostly done)
- Wallet
- Converts
- Funding Payments
- Leveraged Tokens
- Options
- SRM Staking

#### REST
```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/uscott/go-ftx"
)

func main() {
	client := ftx.New(
		ftx.WithAuth("API-KEY", "API-SECRET"),
		ftx.WithHTTPClient(&http.Client{
			Timeout: 5 * time.Second,
		}),
	)

	info, err := client.Account.GetAccountInformation()
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
}
```

#### WebSocket
```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uscott/go-ftx"
)

func main() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    ctx, cancel := context.WithCancel(context.Background())

    client := ftx.New()
    client.Stream.SetDebugMode(true)

    data, err := client.Stream.SubscribeToTickers(ctx, "ETH/BTC")
    if err != nil {
        log.Fatalf("%+v", err)
    }

    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case msg, ok := <-data:
                if !ok {
                    return
                }
                log.Printf("%+v\n", msg)
            }
        }
    }()

    <-sigs
    cancel()
    time.Sleep(time.Second)
}
```

### Websocket Debug Mode
If needed, it is possible to set debug mode to look at error and system messages in stream methods
```go
    client := ftx.New()
    client.Stream.SetDebugMode(true)
```

### No Logged In Error
"Not logged in" errors usually come from wrong signatures.

FTX released an article on how to authenticate https://blog.ftx.com/blog/api-authentication/

If you have unauthorized error to private methods, then you need to use SetServerTimeDiff()
```go
ftx := New()
ftx.SetServerTimeDiff()
```
