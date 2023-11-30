package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aopoltorzhicky/go_kraken/futureswebsocket"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	kraken := futureswebsocket.New(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET"))
	go kraken.Connect()
	time.Sleep(time.Second)

	if err := kraken.SubscribeToBooks([]string{"PI_XBTUSD"}); err != nil {
		log.Fatalf("SubscribeBook error: %s", err.Error())
	}

	for {
		select {
		case <-signals:
			log.Print("Stopping...")
			return
		case update := <-kraken.Listen():
			switch update.Feed {
			case futureswebsocket.BOOK_SNAPSHOT:
				log.Printf("Book snapshot received for %s", update.ProductId)

			case futureswebsocket.BOOK:
				log.Printf("Book update received for %s", update.ProductId)

			}
		}
	}
}
