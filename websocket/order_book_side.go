package websocket

import (
	"bytes"
	"sort"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
)

type orderBookLevel struct {
	Price  decimal.Decimal
	Volume decimal.Decimal
}

type byPrice []orderBookLevel

func (a byPrice) Len() int           { return len(a) }
func (a byPrice) Less(i, j int) bool { return a[i].Price.Cmp(a[j].Price) == -1 }
func (a byPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func newOrderBookLevels(m map[string]orderBookLevel, asc bool) []orderBookLevel {
	result := make([]orderBookLevel, 0)

	for _, value := range m {
		result = append(result, value)
	}

	if asc {
		sort.Sort(byPrice(result))
	} else {
		sort.Sort(sort.Reverse(byPrice(result)))
	}

	return result
}

// OrderBookSide -
type OrderBookSide struct {
	m               map[string]orderBookLevel
	sorted          []orderBookLevel
	depth           int
	pricePrecision  int32
	volumePrecision int32
	isAsk           bool

	mx *sync.RWMutex
}

func newOrderBookSide(depth, pricePrecision, volumePrecision int, isAsk bool) *OrderBookSide {
	return &OrderBookSide{
		m:               make(map[string]orderBookLevel),
		sorted:          make([]orderBookLevel, 0),
		depth:           depth,
		pricePrecision:  int32(pricePrecision),
		volumePrecision: int32(volumePrecision),
		isAsk:           isAsk,
		mx:              new(sync.RWMutex),
	}
}

func (o *OrderBookSide) applyUpdate(upd OrderBookItem) error {
	flValue, err := upd.Volume.Float64()
	if err != nil {
		return err
	}

	price := decimal.RequireFromString(upd.Price.String())
	key := price.StringFixed(o.pricePrecision)

	o.mx.Lock()
	if flValue == 0 {
		delete(o.m, key)
	} else {
		o.m[key] = orderBookLevel{
			Price:  price,
			Volume: decimal.RequireFromString(upd.Volume.String()),
		}
	}
	o.mx.Unlock()
	return nil
}

func (o *OrderBookSide) applyUpdates(updates []OrderBookItem) error {
	for i := range updates {
		if err := o.applyUpdate(updates[i]); err != nil {
			return err
		}
	}

	o.mx.Lock()
	levels := newOrderBookLevels(o.m, o.isAsk)
	for _, level := range levels[o.depth:] {
		delete(o.m, level.Price.StringFixed(o.pricePrecision))
	}
	o.sorted = levels[:o.depth]
	o.mx.Unlock()

	return nil
}

// Get - receives volume by price. If not exists returns false
func (o *OrderBookSide) Get(price decimal.Decimal) (decimal.Decimal, bool) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	key := price.StringFixed(o.pricePrecision)
	level, ok := o.m[key]
	if !ok {
		return decimal.Zero, ok
	}
	return level.Volume, ok
}

// Range - ranges by order book side from best price to depth
func (o *OrderBookSide) Range(handler func(price, volume decimal.Decimal) error) error {
	o.mx.RLock()
	defer o.mx.RUnlock()

	for i := range o.sorted {
		if err := handler(o.sorted[i].Price, o.sorted[i].Volume); err != nil {
			return err
		}
	}
	return nil
}

// Best - returns best price and volume at this price. If order book is not initialized it returns Zero
func (o *OrderBookSide) Best() (decimal.Decimal, decimal.Decimal) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	if len(o.sorted) == 0 {
		return decimal.Zero, decimal.Zero
	}
	return o.sorted[0].Price, o.sorted[0].Volume
}

func (o *OrderBookSide) checksum() []byte {
	o.mx.RLock()
	defer o.mx.RUnlock()

	var str bytes.Buffer
	for _, level := range o.sorted {
		price := level.Price.StringFixed(o.pricePrecision)
		price = strings.Replace(price, ".", "", 1)
		price = strings.TrimLeft(price, "0")
		str.WriteString(price)

		volume := level.Volume.StringFixed(o.volumePrecision)
		volume = strings.Replace(volume, ".", "", 1)
		volume = strings.TrimLeft(volume, "0")
		str.WriteString(volume)
	}
	return str.Bytes()
}

// String -
func (o *OrderBookSide) String() string {
	o.mx.RLock()
	defer o.mx.RUnlock()

	var str strings.Builder
	for i := range o.sorted {
		str.WriteByte('\t')
		if o.pricePrecision > 0 {
			str.WriteString(o.sorted[i].Price.StringFixed(o.pricePrecision))
		} else {
			str.WriteString(o.sorted[i].Price.String())
		}
		str.WriteString(" [ ")
		if o.volumePrecision > 0 {
			str.WriteString(o.sorted[i].Volume.StringFixed(o.volumePrecision))
		} else {
			str.WriteString(o.sorted[i].Volume.String())
		}
		str.WriteString(" ]\r\n")
	}
	return str.String()
}
