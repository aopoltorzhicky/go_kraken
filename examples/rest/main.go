package main

import (
	"log"

	"../../rest"
)

func main() {
	api := rest.New("api-key", "secret")
	data, err1 := api.GetWebSocketsToken()
	if err1 != nil {
		log.Fatalln(err1)
	}
	log.Println(data)
}
