package main

import (
	"log"

	"github.com/aopoltorzhicky/go_kraken/rest"
)

func main() {
	api := rest.New("", "")
	data, err := api.GetSpread("ADAETH", 0)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(data)
}
