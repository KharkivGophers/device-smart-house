package fridge

import (
	"math/rand"
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/config/fridgeconfig"
)

//DataGenerator generates pseudo-random data that represents devices's behavior
func DataGenerator(ticker *time.Ticker, cBot chan<- models.FridgeGenerData, cTop chan<- models.FridgeGenerData,
	stopInner chan struct{}) {
	for {
		select {
		case <-ticker.C:
			cTop <- models.FridgeGenerData{Time: makeTimestamp(), Data: rand.Float32() * 10}
			cBot <- models.FridgeGenerData{Time: makeTimestamp(), Data: (rand.Float32() * 10) - 8}

		case <-stopInner:
			log.Println("DataGenerator(): wg.Done()")
			return
		}

	}
}

//RunDataGenerator setups DataGenerator
func RunDataGenerator(config *fridgeconfig.DevFridgeConfig, cBot chan<- models.FridgeGenerData,
	cTop chan<- models.FridgeGenerData, c *models.Control) {

	duration := config.GetCollectFreq()
	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)
	stopInner := make(chan struct{})

	configChanged := make(chan struct{})
	config.AddSubIntoPool("DataGenerator", configChanged)

	if config.GetTurned() {
		go DataGenerator(ticker, cBot, cTop, stopInner)
	}

	for {
		select {
		case <-configChanged:
			state := config.GetTurned()
			switch state {
			case true:
				select {
				case <-stopInner:
					stopInner = make(chan struct{})
					ticker = time.NewTicker(time.Duration(config.GetCollectFreq()) * time.Millisecond)
					go DataGenerator(ticker, cBot, cTop, stopInner)
					log.Println("DataGenerator() has been started")
				default:
					close(stopInner)
					ticker.Stop()
					stopInner = make(chan struct{})
					ticker = time.NewTicker(time.Duration(config.GetCollectFreq()) * time.Millisecond)
					go DataGenerator(ticker, cBot, cTop, stopInner)
					log.Println("DataGenerator() has been started")
				}
			case false:
				select {
				case <-stopInner:
					ticker = time.NewTicker(time.Duration(config.GetCollectFreq()) * time.Millisecond)
				default:
					close(stopInner)
					log.Println("DataGenerator() has been killed")
				}
			}
		case <- c.Controller:
			log.Error("Data Generator Failed")
			return
		}
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}