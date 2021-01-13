package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

func main() {
	client := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
		api.WithHTTPClient(&http.Client{
			Timeout: 5 * time.Second,
		}),
	)

	info := models.AccountInformation{}
	err := client.Account.GetAccountInformation(&info)
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
}
