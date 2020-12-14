package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	goftx "github.com/uscott/go-ftx"
)

func main() {
	client := goftx.New(
		goftx.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
		goftx.WithHTTPClient(&http.Client{
			Timeout: 5 * time.Second,
		}),
	)

	info, err := client.Account.GetAccountInformation()
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
}
