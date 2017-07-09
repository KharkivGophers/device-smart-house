package washer

import (
	"time"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
	"os"
	"log"
)

//DataCollector gathers data from DataGenerator
//and sends completed request's structures to the ReqChan channel
func DataCollector(ticker *time.Ticker, turnOversStorage <-chan models.GenerateWasherData, waterTempStorage <-chan models.GenerateWasherData,
	RequestStorage chan models.WasherRequest) {

	var requestturnOversStorage = make(map[int64]int64)
	var requestwaterTempStorage = make(map[int64]float32)

	for {
		select {
		case tv := <-waterTempStorage:
			requestwaterTempStorage[tv.Time] = tv.WaterTemp
		case bv := <-turnOversStorage:
			requestturnOversStorage[bv.Time] = bv.Turnovers
		case <-ticker.C:
			log.Print("Data Collector is working")
			RequestStorage <-constructReq(requestturnOversStorage, requestwaterTempStorage)
		}
	}
}

//RunDataCollector setups DataCollector
func RunDataCollector(config *washerconfig.DevWasherConfig, turnOversStorage <-chan models.GenerateWasherData,
	waterTempStorage <-chan models.GenerateWasherData, RequestStorage chan models.WasherRequest, c *models.Control) {
	washTime := config.GetWashTime()
	rinseTime := config.GetRinseTime()
	spinTime := config.GetSpinTime()
	ticker := time.NewTicker(time.Second * 5)

	timer := time.NewTimer(time.Second * time.Duration(washTime + rinseTime + spinTime))
	go DataCollector(ticker, turnOversStorage, waterTempStorage, RequestStorage)
	<-timer.C
	ticker.Stop()
	log.Println("Washing Machine finished!")
}

func constructReq(turnOversStorage map[int64]int64, waterTempStorage map[int64]float32) models.WasherRequest {

	var washerData models.WasherData
	args := os.Args[1:]

	washerData.WaterTemp = waterTempStorage
	washerData.Turnovers = turnOversStorage

	request := models.WasherRequest{
		Action: "update",
		Time: time.Now().UnixNano(),
		Meta: models.Metadata{
			Type: args[0],
			Name: args[1],
			MAC:  args[2]},
		Data: washerData,
	}

	return request
}