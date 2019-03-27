package kraken_ws

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
		Time:      valToTime(body[0]),
		EndTime:   valToTime(body[1]),
		Open:      valToFloat64(body[2]),
		High:      valToFloat64(body[3]),
		Low:       valToFloat64(body[4]),
		Close:     valToFloat64(body[5]),
		VolumeWAP: valToFloat64(body[6]),
		Volume:    valToFloat64(body[7]),
		Count:     body[8].(int),
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
		Ask:  valToFloat64(body[0]),
		Bid:  valToFloat64(body[1]),
		Time: valToTime(body[2]),
		Pair: pair,
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
