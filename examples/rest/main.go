package main

import (
	"log"

	"github.com/aopoltorzhicky/go_kraken/rest"
)

func main() {
	api := rest.New("", "")
	data, err1 := api.AddOrder("XXBTZUSD", "buy", "limit", 0.1, map[string]interface{}{
		"price":    1000.00000,
		"leverage": "2",
	})
	if err1 != nil {
		log.Fatalln(err1)
	}
	log.Println(data)
	orderID := data.TransactionIds[0]
	resp, err2 := api.Cancel(orderID)
	if err2 != nil {
		log.Fatalln(err2)
	}
	log.Println(resp)
}
