package main

import (
	"log"

	"github.com/aopoltorzhicky/go_kraken/rest"
)

func main() {
	api := rest.New("pZai5d5lPFm4uE+raQlSHCRRqNE4MoPU/k1C8nu0/0bZZ9+iItHsiWoc", "iNYlYHYf5j4tE8EvSV+HN6VsJyX3dhkRlw/6+NLLHzDZ+9aA8j/o3ze9gpXDXyEZVnz88vtwsqF8fPZWv82egw==")
	data, err1 := api.GetTradeBalance("ZUSD")
	if err1 != nil {
		log.Fatalln(err1)
	}
	log.Println(data)
}
