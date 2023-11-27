package main

import (
	"log"
	"os"

	"github.com/bknigge/go_kraken/rest"
)

func main() {
	api := rest.New(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET"))
	data, err := api.Ticker("XXBTZUSD")
	if err != nil {
		log.Panicln(err)
		return
	}
	for _, ticker := range data {
		log.Printf("ask %s", ticker.Ask.Price)
	}
}
