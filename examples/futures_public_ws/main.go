package main

import (
	"os"

	"github.com/aopoltorzhicky/go_kraken/futureswebsocket"
)

func main() {
	ws := futureswebsocket.New(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET"))
	ws.ConnectAndProcessUpdates()
}
