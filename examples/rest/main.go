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
	log.Println(data)
}
