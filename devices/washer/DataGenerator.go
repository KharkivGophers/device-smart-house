package washer

import (
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
	"time"
	"log"
	"math/rand"
)

func DataGenerator(stage string, ticker *time.Ticker, turnOversStorage chan<- models.GenerateWasherData,
	waterTempStorage chan<- models.GenerateWasherData) {

	log.Println(stage, "started!")
	for {
		select {
		case <-ticker.C:
			log.Println("Yo", stage)
			turnOversStorage <- models.GenerateWasherData{Time:makeTimestamp(), Turnovers: rand.Int63n(100)}
			waterTempStorage <- models.GenerateWasherData{Time:makeTimestamp(), WaterTemp: rand.Float32() * 10}
		}
	}
	log.Println(stage, "finished!")
}

func RunDataGenerator(config *washerconfig.DevWasherConfig, turnOversStorage chan<- models.GenerateWasherData,
	waterTempStorage chan<- models.GenerateWasherData, c *models.Control) {

	// Run wash
	washTime := config.GetWashTime()
	stageWash := "Wash"
	ticker := time.NewTicker(time.Second * 3)
	timer := time.NewTimer(time.Second * time.Duration(washTime))
	go DataGenerator(stageWash ,ticker, turnOversStorage, waterTempStorage)
	<-timer.C
	ticker.Stop()

	// Run rinse
	rinseTime := config.GetRinseTime()
	stageRinse := "Rinse"
	ticker = time.NewTicker(time.Second * 3)
	timer = time.NewTimer(time.Second * time.Duration(rinseTime))
	go DataGenerator(stageRinse, ticker, turnOversStorage, waterTempStorage)
	<-timer.C
	ticker.Stop()

	// Run spin
	spinTime := config.GetSpinTime()
	stageSpin := "Spin"
	ticker = time.NewTicker(time.Second * 3)
	timer = time.NewTimer(time.Second * time.Duration(spinTime))
	go DataGenerator(stageSpin, ticker, turnOversStorage, waterTempStorage)
	<-timer.C
	ticker.Stop()
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

