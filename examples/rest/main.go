package main

import (
	"log"

	"github.com/aopoltorzhicky/go_kraken/rest"
)

func main() {
	api := rest.New("api_key", "secret")
	data, err := api.GetAccountBalances()
	if err != nil {
		log.Panicln(err)
		return
	}
	log.Println(data)
}
