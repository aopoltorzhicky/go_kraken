package main

import (
	"log"
	"os"

	"github.com/aopoltorzhicky/go_kraken/rest"
)

func main() {
	api := rest.New(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET"))
	data, err := api.Ticker("XXBTZUSD")
	if err != nil {
		log.Panicln(err)
		return
	}
	for name, ticker := range data {
		log.Printf("%s ask %s", name, ticker.Ask.Price)
	}
}
