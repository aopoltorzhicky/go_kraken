package futureswebsocket

import (
	"encoding/json"
	"log"
)

func (k *Kraken) handleBookMessage(msg []byte) error {
	var event EventType
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	switch event.Event {
	case SUBSCRIBED:
		log.Printf("%s to %s", event.Event, event.Feed)
		return nil
	case UNSUBSCRIBED:
		log.Printf("%s from %s", event.Event, event.Feed)
		return nil
	case INFO:
		return nil
	}

	switch event.Feed {
	case BOOK:
		var bookUpdate BookUpdateEvent
		if err := json.Unmarshal(msg, &bookUpdate); err != nil {
			return err
		}
		log.Printf("%s", bookUpdate)
	case BOOK_SNAPSHOT:
		var bookSnapshot BookSnapshotEvent
		if err := json.Unmarshal(msg, &bookSnapshot); err != nil {
			return err
		}
		log.Printf("%s", bookSnapshot)
	default:
		log.Printf("unknown event: %s", msg)
	}
	return nil
}
