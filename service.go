package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hargrave81/binaryfix/engine"
)

// CurrencyService is the service that will loop every X seconds and process the currency rates and if we should trade
type CurrencyService struct {
	Interval int
}

//Start will start the service
func (cs *CurrencyService) Start(interval int) error {
	cs.Interval = interval
	mainloop := make(chan bool, 1)

	// TradeQueue holds the last 10 trades for a given currency
	TradeQueue := engine.CreateTradeQueue()

	var bubble error
	go func() {
		for {
			forever := make(chan bool, 1)
			timeout := make(chan bool, 1)
			go func() {
				d := time.Duration(cs.Interval) * time.Second
				time.Sleep(d)
				timeout <- true
			}()
			go func() {
				log.Println("Processing Currency Index ...")

				select {
				case <-timeout:
					{
						// DO WORK
						err := performService(TradeQueue)
						if err != nil {
							bubble = fmt.Errorf("%s: %s", "Failed to process currency", err)
							mainloop <- true
							return
						}
						// bubble out just incase the interval has changed
						forever <- true
					}
				}
			}()
			<-forever
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-mainloop
	return bubble
}

func performService(tq *engine.TradeQueueEngine) error {
	// get the current rates
	result := engine.GetStocks()
	// loop through all the rates and update the stats
	for k, r := range result {
		tq.UpdateTrade(k, r)
	}
	// now lets see if we need to do something
	tradeBuys, tradeSells := tq.GetTrades()
	if len(tradeBuys) > 0 {

	}
	if len(tradeSells) > 0 {

	}
	return nil
}
