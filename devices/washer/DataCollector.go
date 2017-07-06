package washer

import (
	"time"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
)

//DataCollector gathers data from DataGenerator
//and sends completed request's structures to the ReqChan channel
func DataCollector(ticker *time.Ticker, cBot <-chan models.GenerateWasherData, cTop <-chan models.GenerateWasherData,
	ReqChan chan models.WasherRequest, stopInner chan struct{}) {
}

//RunDataCollector setups DataCollector
func RunDataCollector(config *washerconfig.DevWasherConfig, cBot <-chan models.GenerateWasherData,
	cTop <-chan models.GenerateWasherData, ReqChan chan models.WasherRequest, c *models.Control) {
}

func constructReq(mTop map[int64]float32, mBot map[int64]float32) models.WasherRequest {
	return models.WasherRequest{}
}