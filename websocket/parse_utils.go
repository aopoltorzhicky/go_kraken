package websocket

import (
	"log"
	"strconv"
	"time"
)

type ParseFactory interface {
	Parse(data interface{}, pair string) (interface{}, error)
}

func valToFloat64(value interface{}) float64 {
	if v, ok := value.(string); ok {
		result, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Printf("Can't parse float %#v", value)
			return .0
		}
		return result
	}
	return .0
}

func valToTime(data interface{}) time.Time {
	ret := valToFloat64(data) * 1e9
	return time.Unix(0, int64(ret))
}

func parseLevel(data []interface{}) Level {
	return Level{
		Price:          valToFloat64(data[0].(string)),
		Volume:         valToFloat64(data[2].(string)),
		WholeLotVolume: data[1].(float64),
	}
}

func parseValues(data []interface{}) Values {
	var today, last24h float64
	switch data[0].(type) {
	case string:
		today = valToFloat64(data[0].(string))
		last24h = valToFloat64(data[1].(string))
	case float64:
		today = data[0].(float64)
		last24h = data[1].(float64)
	}

	return Values{
		Today:  today,
		Last24: last24h,
	}
}

func parseSide(data string) string {
	if data == "s" {
		return Sell
	} else if data == "b" {
		return Buy
	}
	log.Printf("Unknown side: %s", data)
	return ""
}

func parseOrderType(data string) string {
	if data == "l" {
		return Limit
	} else if data == "m" {
		return Market
	}
	log.Printf("Unknown order type: %s", data)
	return ""
}
