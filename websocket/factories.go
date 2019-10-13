package websocket

import (
	"fmt"
)

type tickerFactory struct{}

func newTickerFactory() *tickerFactory {
	return &tickerFactory{}
}

func (f *tickerFactory) Parse(data interface{}, pair string) (interface{}, error) {
	result := TickerUpdate{
		Pair: pair,
	}
	body, ok := data.(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("Can't parse data %#v", data)
	}

	for k, v := range body {
		value := v.([]interface{})
		switch k {
		case "a":
			result.Ask = parseLevel(value)
		case "b":
			result.Bid = parseLevel(value)
		case "c":
			result.Close = parseValues(value)
		case "v":
			result.Volume = parseValues(value)
		case "o":
			result.Open = parseValues(value)
		case "h":
			result.High = parseValues(value)
		case "l":
			result.Low = parseValues(value)
		case "t":
			result.TradeVolume = parseValues(value)
		case "p":
			result.VolumeAveragePrice = parseValues(value)
		}
	}
	return result, nil
}

type candlesFactory struct{}

func newCandlesFactory() *candlesFactory {
	return &candlesFactory{}
}

func (f *candlesFactory) Parse(data interface{}, pair string) (interface{}, error) {
	body, ok := data.([]interface{})
	if !ok {
		return CandleUpdate{Pair: pair}, fmt.Errorf("Can't parse data %#v", data)
	}
	return CandleUpdate{
		Time:      valToFloat64(body[0]),
		EndTime:   valToFloat64(body[1]),
		Open:      valToFloat64(body[2]),
		High:      valToFloat64(body[3]),
		Low:       valToFloat64(body[4]),
		Close:     valToFloat64(body[5]),
		VolumeWAP: valToFloat64(body[6]),
		Volume:    valToFloat64(body[7]),
		Count:     int(body[8].(float64)),
		Pair:      pair,
	}, nil
}

type tradesFactory struct{}

func newTradesFactory() *tradesFactory {
	return &tradesFactory{}
}

func (f *tradesFactory) Parse(data interface{}, pair string) (interface{}, error) {
	result := []TradeUpdate{}
	body, ok := data.([]interface{})
	if !ok {
		return result, fmt.Errorf("Can't parse data %#v", data)
	}
	for _, item := range body {
		entity := item.([]interface{})
		trade := TradeUpdate{
			Price:     valToFloat64(entity[0]),
			Volume:    valToFloat64(entity[1]),
			Time:      valToFloat64(entity[2]),
			Side:      parseSide(entity[3].(string)),
			OrderType: parseOrderType(entity[4].(string)),
			Misc:      entity[5].(string),
			Pair:      pair,
		}
		result = append(result, trade)
	}
	return result, nil
}

type spreadFactory struct{}

func newSpreadFactory() *spreadFactory {
	return &spreadFactory{}
}

func (f *spreadFactory) Parse(data interface{}, pair string) (interface{}, error) {
	body, ok := data.([]interface{})
	if !ok {
		return SpreadUpdate{Pair: pair}, fmt.Errorf("Can't parse data %#v", data)
	}
	return SpreadUpdate{
		Bid:       valToFloat64(body[0]),
		Ask:       valToFloat64(body[1]),
		Time:      valToFloat64(body[2]),
		BidVolume: valToFloat64(body[3]),
		AskVolume: valToFloat64(body[4]),
		Pair:      pair,
	}, nil
}

type bookFactory struct{}

func newBookFactory() *bookFactory {
	return &bookFactory{}
}

func (f *bookFactory) Parse(data interface{}, pair string) (interface{}, error) {
	result := OrderBookUpdate{
		Pair: pair,
	}
	body, ok := data.(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("Can't parse data %#v", data)
	}

	for k, v := range body {
		items := make([]OrderBookItem, 0)
		updates := v.([]interface{})
		for _, item := range updates {
			entity := item.([]interface{})
			orderBookItem := OrderBookItem{
				Price:  valToFloat64(entity[0]),
				Volume: valToFloat64(entity[1]),
				Time:   valToFloat64(entity[2]),
			}
			orderBookItem.Republish = (len(entity) == 4 && entity[3] == "r")
			items = append(items, orderBookItem)
		}

		switch k {
		case "as":
			result.IsSnapshot = true
			result.Asks = items
		case "a":
			result.Asks = items
		case "bs":
			result.IsSnapshot = true
			result.Bids = items
		case "b":
			result.Bids = items
		}
	}
	return result, nil
}

type ownTradesFactory struct{}

func newOwnTradesFactory() *ownTradesFactory {
	return &ownTradesFactory{}
}

func (f *ownTradesFactory) Parse(data interface{}, pair string) (interface{}, error) {
	upd := OwnTradesUpdate{
		ChannelName: ChanOwnTrades,
		Trades:      make(map[string]OwnTrade),
	}
	body, ok := data.([]interface{})
	if !ok {
		return upd, fmt.Errorf("Can't parse data %#v", data)
	}

	if len(body) != 2 {
		return upd, fmt.Errorf("Can't parse data %#v", data)
	}

	for key, value := range body[0].(map[string]map[string]interface{}) {
		upd.Trades[key] = OwnTrade{
			Cost:      valToFloat64(value["cost"]),
			Fee:       valToFloat64(value["fee"]),
			Margin:    valToFloat64(value["margin"]),
			OrderID:   value["ordertxid"].(string),
			OrderType: value["ordertype"].(string),
			Pair:      value["pair"].(string),
			PosTxID:   value["postxid"].(string),
			Price:     valToFloat64(value["price"]),
			Time:      valToFloat64(value["time"]),
			Type:      value["type"].(string),
			Vol:       valToFloat64(value["vol"]),
		}
	}

	return upd, nil
}

type openOrdersFactory struct{}

func newOpenOrdersFactory() *openOrdersFactory {
	return &openOrdersFactory{}
}

func (f *openOrdersFactory) Parse(data interface{}, pair string) (interface{}, error) {
	upd := OpenOrdersUpdate{
		ChannelName: ChanOpenOrders,
		Order:       make(map[string]OpenOrder),
	}
	body, ok := data.([]interface{})
	if !ok {
		return upd, fmt.Errorf("Can't parse data %#v", data)
	}

	if len(body) != 2 {
		return upd, fmt.Errorf("Can't parse data %#v", data)
	}

	for key, value := range body[0].(map[string]map[string]interface{}) {
		upd.Order[key] = OpenOrder{
			Cost:       valToFloat64(value["cost"]),
			Fee:        valToFloat64(value["fee"]),
			LimitPrice: valToFloat64(value["limitprice"]),
			Misc:       value["misc"].(string),
			Oflags:     value["oflags"].(string),
			OpenTime:   valToFloat64(value["opentm"]),
			StartTime:  valToFloat64(value["starttm"]),
			ExpireTime: valToFloat64(value["expiretm"]),
			Price:      valToFloat64(value["price"]),
			Refid:      value["refid"].(string),
			Status:     value["status"].(string),
			StopPrice:  valToFloat64(value["stopprice"]),
			UserRef:    int(value["userref"].(float64)),
			Vol:        valToFloat64(value["vol"]),
			VolExec:    valToFloat64(value["vol_exec"]),

			Descr: OpenOrderDescr{
				Close:     value["close"].(string),
				Leverage:  value["leverage"].(string),
				Order:     value["order"].(string),
				Ordertype: value["ordertype"].(string),
				Pair:      value["pair"].(string),
				Price:     valToFloat64(value["price"]),
				Price2:    valToFloat64(value["price2"]),
				Type:      value["type"].(string),
			},
		}
	}

	return upd, nil
}
