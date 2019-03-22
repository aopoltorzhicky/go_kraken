package websocket

import (
	"fmt"
)

type TickerFactory struct{}

func newTickerFactory() *TickerFactory {
	return &TickerFactory{}
}

func (f *TickerFactory) Parse(data interface{}, pair string) (interface{}, error) {
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

type CandlesFactory struct{}

func newCandlesFactory() *CandlesFactory {
	return &CandlesFactory{}
}

func (f *CandlesFactory) Parse(data interface{}, pair string) (interface{}, error) {
	body, ok := data.([]interface{})
	if !ok {
		return CandleUpdate{}, fmt.Errorf("Can't parse data %#v", data)
	}
	return CandleUpdate{
		Time:      valToTime(body[0]),
		EndTime:   valToFloat64(body[1]),
		Open:      valToFloat64(body[2]),
		High:      valToFloat64(body[3]),
		Low:       valToFloat64(body[4]),
		Close:     valToFloat64(body[5]),
		VolumeWAP: valToFloat64(body[6]),
		Volume:    valToFloat64(body[7]),
		Count:     int64(body[8].(float64)),
		Pair:      pair,
	}, nil
}

type TradesFactory struct{}

func newTradesFactory() *TradesFactory {
	return &TradesFactory{}
}

func (f *TradesFactory) Parse(data interface{}, pair string) (interface{}, error) {
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
			Time:      valToTime(entity[2]),
			Side:      parseSide(entity[3].(string)),
			OrderType: parseOrderType(entity[4].(string)),
			Misc:      entity[5].(string),
			Pair:      pair,
		}
		result = append(result, trade)
	}
	return result, nil
}

type SpreadFactory struct{}

func newSpreadFactory() *SpreadFactory {
	return &SpreadFactory{}
}

func (f *SpreadFactory) Parse(data interface{}, pair string) (interface{}, error) {
	body, ok := data.([]interface{})
	if !ok {
		return SpreadUpdate{}, fmt.Errorf("Can't parse data %#v", data)
	}
	return SpreadUpdate{
		Ask:  valToFloat64(body[0]),
		Bid:  valToFloat64(body[1]),
		Time: valToTime(body[2]),
		Pair: pair,
	}, nil
}

type BookFactory struct{}

func newBookFactory() *BookFactory {
	return &BookFactory{}
}

func (f *BookFactory) Parse(data interface{}, pair string) (interface{}, error) {
	result := OrderBookUpdate{
		Pair: pair,
	}
	body, ok := data.(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("Can't parse data %#v", data)
	}
	for k, v := range body {
		var items []OrderBookItem
		updates := v.([]interface{})
		for _, item := range updates {
			entity := item.([]interface{})
			orderBookItem := OrderBookItem{
				Price:  valToFloat64(entity[0]),
				Volume: valToFloat64(entity[1]),
				Time:   valToTime(entity[2]),
			}
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
