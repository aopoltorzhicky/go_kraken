package futureswebsocket

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Kraken struct {
	Key    string
	Secret string
	Conn   *websocket.Conn
	Msg    chan Update
}

// New - constructor of Kraken object
func New(key string, secret string) *Kraken {
	if key == "" || secret == "" {
		log.Print("[WARNING] You are not set api key and secret!")
	}
	return &Kraken{
		Key:    key,
		Secret: secret,
		Msg:    make(chan Update, 1024),
	}
}

func (k *Kraken) SubscribeToBooks(productIds []string) error {
	subscribeMsg, err := json.Marshal(SubscribeBook{
		Event:      SUBSCRIBE,
		Feed:       BOOK,
		ProductIds: productIds,
	})
	if err != nil {
		return err
	}

	err = k.Conn.WriteMessage(websocket.TextMessage, subscribeMsg)
	if err != nil {
		log.Println("write:", err)
		return err
	}

	return nil
}

func (k *Kraken) Connect() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: "futures.kraken.com", Path: "/ws/v1"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	k.Conn = conn
	done := make(chan struct{})

	// Read messages from server
	go func() {
		defer close(done)
		for {
			_, message, err := k.Conn.ReadMessage()
			if err != nil {
				log.Print("read:", err)
				return
			}

			if err := k.handleMessage(message); err != nil {
				log.Print("read:", err)
				return
			}
		}
	}()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			err = k.Conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("write:", err)
				return err
			}
		case <-interrupt:
			log.Println("interrupted")

			// Graceful shutdown by sending a close frame and then waiting for the server to close the connection.
			err := k.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func (k *Kraken) Listen() <-chan Update {
	return k.Msg
}

// SignChallenge signs the given challenge using the provided API secret.
func (k *Kraken) SignChallenge(challenge string) (string, error) {
	// Step 1: Hash the challenge with SHA-256
	sha256Hasher := sha256.New()
	sha256Hasher.Write([]byte(challenge))
	hashedChallenge := sha256Hasher.Sum(nil)

	// Step 2: Base64-decode the API secret
	decodedSecret, err := base64.StdEncoding.DecodeString(k.Secret)
	if err != nil {
		return "", err
	}

	// Step 3: Use the result of step 2 to hash the result of step 1 with HMAC-SHA-512
	hmac512 := hmac.New(sha512.New, decodedSecret)
	hmac512.Write(hashedChallenge)
	hmacHash := hmac512.Sum(nil)

	// Step 4: Base64-encode the result of step 3
	signedOutput := base64.StdEncoding.EncodeToString(hmacHash)

	return signedOutput, nil
}
