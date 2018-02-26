package engine

import (
	"time"

	"github.com/hargrave81/binaryfix/queue"
)

// Trade represents a trade object at a given time for a given currency
type Trade struct {
	Value        float64
	TradeBuy     bool
	TradeSell    bool
	NoTradeBuy   bool
	NoTradeSell  bool
	RSI          float64
	RS           float64
	Gained       float64
	Lost         float64
	AverageGain  float64
	AverageLoss  float64
	Curreny      string
	InstanceDate time.Time
}

////////////////////////////////////// Queue Helpers

// TradeQueueEngine is the engine holder for function access
type TradeQueueEngine struct {
	MaxTradeCount int
	tradeQueue    map[string]queue.Queue
	AverageGain   float64
	AverageLoss   float64
}

// CreateTradeQueue creates a tradequeueengine
func CreateTradeQueue() *TradeQueueEngine {
	TradeQueue := &TradeQueueEngine{}
	TradeQueue.tradeQueue = make(map[string]queue.Queue)
	TradeQueue.MaxTradeCount = 10
	TradeQueue.AverageGain = 8 // T1
	TradeQueue.AverageLoss = 8 // T2
	return TradeQueue
}

// UpdateTrade will update a given currency's trade rates
func (t *TradeQueueEngine) UpdateTrade(currency string, value float64) {
	// create currency queue if it does not exist
	if _, ok := t.tradeQueue[currency]; !ok {
		genQueue := queue.NewQueue(t.MaxTradeCount)
		t.tradeQueue[currency] = genQueue
	}

	// calculate the new trade
	newTrade := t.calculateRSI(currency, value, t.LastTrade(currency))

	// push the new trade into the queue
	t.tradeQueue[currency].Push(queue.Node{Value: &newTrade})

	// throw away the last entry in the queue
	if t.tradeQueue[currency].Count() > t.MaxTradeCount {
		t.tradeQueue[currency].Pop()
	}
}

// LastTrade Return the last Trade in the queue
func (t *TradeQueueEngine) LastTrade(currency string) *Trade {
	if v, ok := t.tradeQueue[currency]; ok {
		if v.Count() > 0 {
			return v.Peek().Value.(*Trade)
		}
	}
	return &Trade{}
}

// Calculte the RSI
func (t *TradeQueueEngine) calculateRSI(currency string, value float64, lastTrade *Trade) *Trade {
	newTrade := &Trade{}
	newTrade.Value = value
	newTrade.Curreny = currency
	if value > lastTrade.Value && lastTrade.Value > 0 {
		newTrade.Gained = newTrade.Value - lastTrade.Value
	} else if value < lastTrade.Value && lastTrade.Value > 0 {
		newTrade.Lost = lastTrade.Value - newTrade.Value
	}

	if t.tradeQueue[currency].Count() == t.MaxTradeCount {
		// we have all the trades we allow in the queue
		newTrade.AverageGain, newTrade.AverageLoss = t.Average(currency, 999999) // just use the max rows rather than a custom row count here
	}
	// there is other logic that i have removed in this step
	if newTrade.AverageLoss > 0 {
		newTrade.RS = newTrade.AverageGain / newTrade.AverageLoss
	} else {
		newTrade.RS = 0
	}

	newTrade.RSI = 100 - (100 / (newTrade.RS + 1))

	if newTrade.RSI > t.AverageGain && lastTrade.TradeBuy || newTrade.RSI > t.AverageGain && lastTrade.NoTradeBuy {
		newTrade.TradeBuy = false
		newTrade.NoTradeBuy = true
	} else if newTrade.RSI > t.AverageGain && newTrade.RSI > 0 && !lastTrade.NoTradeBuy {
		newTrade.TradeBuy = true
	}

	if newTrade.RSI < t.AverageLoss && lastTrade.TradeSell || newTrade.RSI < t.AverageLoss && lastTrade.NoTradeSell {
		newTrade.TradeSell = false
		newTrade.NoTradeSell = true
	} else if newTrade.RSI < t.AverageLoss && newTrade.RSI > 0 && !lastTrade.NoTradeSell {
		newTrade.TradeSell = true
	}
	return newTrade
}

// Average Returns the average for a given currency type for both Gain and Loss
func (t *TradeQueueEngine) Average(currency string, rowCount int) (float64, float64) {
	totalUp := 0.0
	totalDown := 0.0
	if rowCount <= 0 {
		rowCount = 999999
	}
	rows := 0
	if _, ok := t.tradeQueue[currency]; ok {
		for _, v := range t.tradeQueue[currency].Slice() {
			totalUp += (v.Value.(*Trade)).Gained
			totalDown += (v.Value.(*Trade)).Lost
			rows++
			rowCount--
			if rowCount == 0 {
				break
			}
		}
	}
	return totalUp / float64(rows), totalDown / float64(rows)
}

// GetTrades will return both tradebuy and tradesell
func (t *TradeQueueEngine) GetTrades() ([]*Trade, []*Trade) {
	var result1 []*Trade
	var result2 []*Trade
	for _, v := range t.tradeQueue {
		value := v.Peek().Value
		if value != nil {
			if value.(*Trade).TradeBuy {
				result1 = append(result1, value.(*Trade))
			} else if value.(Trade).TradeSell {
				result2 = append(result2, value.(*Trade))
			}
		}
	}
	return result1, result2
}
