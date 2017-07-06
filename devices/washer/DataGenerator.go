package washer

import (
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
	"time"
)

//DataGenerator generates pseudo-random data that represents devices's behavior
func DataGenerator(ticker *time.Ticker, cBot chan<- models.FridgeGenerData, cTop chan<- models.FridgeGenerData,
	stopInner chan struct{}) {
}

//RunDataGenerator setups DataGenerator
func RunDataGenerator(config *washerconfig.DevWasherConfig, cBot chan<- models.GenerateWasherData,
	cTop chan<- models.GenerateWasherData, c *models.Control) {
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}