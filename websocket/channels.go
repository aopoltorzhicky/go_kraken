package websocket

import (
	"encoding/json"
	"strings"
)

func (k *Kraken) handleChannel(data []byte) error {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	channel := strings.Split(msg.ChannelName, "-")[0]
	switch channel {
	case ChanTicker:
		var ticker TickerUpdate
		if err := json.Unmarshal(msg.Data, &ticker); err != nil {
			return err
		}
		k.msg <- msg.toUpdate(ticker)
	case ChanCandles:
		var candle Candle
		if err := json.Unmarshal(msg.Data, &candle); err != nil {
			return err
		}
		k.msg <- msg.toUpdate(candle)
	case ChanTrades:
		var trades []Trade
		if err := json.Unmarshal(msg.Data, &trades); err != nil {
			return err
		}
		k.msg <- msg.toUpdate(trades)
	case ChanSpread:
		var spread Spread
		if err := json.Unmarshal(msg.Data, &spread); err != nil {
			return err
		}
		k.msg <- msg.toUpdate(spread)
	case ChanBook:
		var update OrderBookUpdate
		if err := json.Unmarshal(msg.Data, &update); err != nil {
			return err
		}
		k.msg <- msg.toUpdate(update)
	case ChanOwnTrades:
		var update OwnTradesUpdate
		if err := json.Unmarshal(msg.Data, &update); err != nil {
			return err
		}
		k.msg <- msg.toUpdate(update)
	case ChanOpenOrders:
		var update OpenOrdersUpdate
		if err := json.Unmarshal(msg.Data, &update); err != nil {
			return err
		}
		k.msg <- msg.toUpdate(update)
	}

	return nil
}
