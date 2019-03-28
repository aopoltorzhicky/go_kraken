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
	return time.Unix(0, int64(ret)).UTC()
}

func parseLevel(data []interface{}) Level {
	wholeLot := 0
	if v, ok := data[1].(int); ok {
		wholeLot = v
	}
	return Level{
		Price:          valToFloat64(data[0]),
		Volume:         valToFloat64(data[2]),
		WholeLotVolume: wholeLot,
	}
}

func parseValues(data []interface{}) Values {
	switch data[0].(type) {
	case string:
		return Values{
			Today:  valToFloat64(data[0]),
			Last24: valToFloat64(data[1]),
		}
	case int:
		last24h := 0
		if v, ok := data[1].(int); ok {
			last24h = v
		}
		return Values{
			Today:  data[0].(int),
			Last24: last24h,
		}

	case float64:
		last24h := 0.0
		if v, ok := data[1].(float64); ok {
			last24h = v
		}
		return Values{
			Today:  data[0].(float64),
			Last24: last24h,
		}
	default:
		return Values{
			Today:  0,
			Last24: 0,
		}
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
