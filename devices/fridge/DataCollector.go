package fridge

import (
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/KharkivGophers/device-smart-house/models"
	"os"
	"github.com/KharkivGophers/device-smart-house/config/fridgeconfig"
)

//DataCollector gathers data from DataGenerator
//and sends completed request's structures to the ReqChan channel
func DataCollector(ticker *time.Ticker, cBot <-chan models.FridgeGenerData, cTop <-chan models.FridgeGenerData,
	ReqChan chan models.FridgeRequest, stopInner chan struct{}) {
	var mTop = make(map[int64]float32)
	var mBot = make(map[int64]float32)

	for {
		select {
		case <-stopInner:
			log.Println("DataCollector(): wg.Done()")
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
func RunDataCollector(config *fridgeconfig.DevFridgeConfig, cBot <-chan models.FridgeGenerData,
	cTop <-chan models.FridgeGenerData, ReqChan chan models.FridgeRequest, c *models.Control) {

	duration := config.GetSendFreq()
	stopInner := make(chan struct{})
	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)

	configChanged := make(chan struct{})
	config.AddSubIntoPool("DataCollector", configChanged)

	//wg.Add(1)
	if config.GetTurned() {
		go DataCollector(ticker, cBot, cTop, ReqChan, stopInner)
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
					ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Millisecond)
					go DataCollector(ticker, cBot, cTop, ReqChan, stopInner)
					log.Println("DataCollector() has been started")
				default:
					close(stopInner)
					stopInner = make(chan struct{})
					ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Millisecond)
					go DataCollector(ticker, cBot, cTop, ReqChan, stopInner)
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
		case <- c.Controller:
			log.Error("Data Collector Failed")
			return
		}
	}
}

func constructReq(mTop map[int64]float32, mBot map[int64]float32) models.FridgeRequest {
	var fridgeData models.FridgeData
	args := os.Args[1:]

	fridgeData.TempCam2 = mBot
	fridgeData.TempCam1 = mTop

	req := models.FridgeRequest{
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
