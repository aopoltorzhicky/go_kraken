package futureswebsocket

import (
	"encoding/json"
	"log"
)

func (k *Kraken) handleMessage(msg []byte) error {
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
		return k.handleBookUpdate(msg)
	case BOOK_SNAPSHOT:
		return k.handleBookSnapshotUpdate(msg)
	}

	return nil
}

func (k *Kraken) handleBookUpdate(msg []byte) error {
	var update Update
	if err := json.Unmarshal(msg, &update); err != nil {
		return err
	}

	var bookUpdateEvent BookUpdateEvent
	if err := json.Unmarshal(msg, &bookUpdateEvent); err != nil {
		return err
	}
	update.Data = bookUpdateEvent
	k.Msg <- update

	return nil
}

func (k *Kraken) handleBookSnapshotUpdate(msg []byte) error {
	var update Update
	if err := json.Unmarshal(msg, &update); err != nil {
		return err
	}

	var bookSnapshotEvent BookSnapshotEvent
	if err := json.Unmarshal(msg, &bookSnapshotEvent); err != nil {
		return err
	}
	update.Data = bookSnapshotEvent
	k.Msg <- update

	return nil
}
