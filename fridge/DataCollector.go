package fridge

import (
	"time"
	"sync"
	"log"
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
	"os"
)

//DataCollector gathers data from DataGenerator
//and sends completed request's structures to the ReqChan channel
func DataCollector(ticker *time.Ticker, cBot <-chan models.FridgeGenerData, cTop <-chan models.FridgeGenerData,
	ReqChan chan models.Request, stopInner chan struct{}, wg *sync.WaitGroup) {

	var mTop = make(map[int64]float32)
	var mBot = make(map[int64]float32)

	for {
		select {
		case <-stopInner:

			log.Println("DataCollector(): wg.Done()")
			wg.Done()
			return
		case tv := <-cTop:
			mTop[tv.Time] = tv.Data
		case bv := <-cBot:
			mBot[bv.Time] = bv.Data
		case <-ticker.C:
			ReqChan <- constructReq(mTop, mBot)

			//Cleaning temp maps
			mTop = make(map[int64]float32)
			mBot = make(map[int64]float32)
		}

	}
}

//RunDataCollector setups DataCollector
func RunDataCollector(config *config.DevConfig, cBot <-chan models.FridgeGenerData,
	cTop <-chan models.FridgeGenerData, ReqChan chan models.Request, wg *sync.WaitGroup) {
	duration := config.GetSendFreq()
	stopInner := make(chan struct{})
	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)

	configChanged := make(chan struct{})
	config.AddSubIntoPool("DataCollector", configChanged)

	wg.Add(1)
	if config.GetTurned() {
		go DataCollector(ticker, cBot, cTop, ReqChan, stopInner, wg)
	}

	for {
		select {
		case <-configChanged:
			state := config.GetTurned()
			switch state {
			case true:
				select {
				case <-stopInner:
					wg.Add(1)
					stopInner = make(chan struct{})
					ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Millisecond)
					go DataCollector(ticker, cBot, cTop, ReqChan, stopInner, wg)
					log.Println("DataCollector() has been started")
				default:
					close(stopInner)
					stopInner = make(chan struct{})
					wg.Add(1)
					ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Millisecond)
					go DataCollector(ticker, cBot, cTop, ReqChan, stopInner, wg)
					log.Println("DataCollector() has been started")
				}
			case false:
				select {
				case <-stopInner:
					ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Millisecond)
				default:
					close(stopInner)
					log.Println("DataCollector() has been killed")
				}
			}
		}
	}
}

func constructReq(mTop map[int64]float32, mBot map[int64]float32) models.Request {
	var fridgeData models.FridgeData
	args := os.Args[1:]

	fridgeData.TempCam2 = mBot
	fridgeData.TempCam1 = mTop

	req := models.Request{
		Action: "update",
		Time:   time.Now().UnixNano(),
		Meta: models.Metadata{
			Type: args[0],
			Name: args[1],
			MAC:  args[2]},
		Data: fridgeData,
	}
	return req
}
