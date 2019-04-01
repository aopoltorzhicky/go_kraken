package main

import (
	"log"

	"github.com/aopoltorzhicky/go_kraken/rest"
)

func main() {
	api := rest.New("", "")
	data, err := api.Candles("ADAETH", 0, 1)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(data)
}
