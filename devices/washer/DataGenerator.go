package washer

import (
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
	"github.com/KharkivGophers/device-smart-house/models"
	"log"
	"math/rand"
	"time"
)

// DataGenerator generates pseudo-random numbers
func DataGenerator(stage string, ticker *time.Ticker, maxTemperature int64, minTurnovers int64, maxTurnovers int64, turnOversStorage chan<- models.GenerateWasherData,
	waterTempStorage chan<- models.GenerateWasherData) {

	log.Println(stage, "started!")
	for {
		select {
		case <-ticker.C:
			turnOversStorage <- models.GenerateWasherData{Time: makeTimestamp(), Turnovers: int64(rand.Intn(int(maxTurnovers)))}
			waterTempStorage <- models.GenerateWasherData{Time: makeTimestamp(), WaterTemp: rand.Float32() * float32(maxTemperature)}
		}
	}
}

func RunDataGenerator(config *washerconfig.DevWasherConfig, turnOversStorage chan<- models.GenerateWasherData,
	waterTempStorage chan<- models.GenerateWasherData, c *models.Control, firstStep chan struct{}) {

	maxTemperature := config.GetTemperature()
	// Run wash
	washTime := config.GetWashTime()
	maxWashTurnovers := config.GetWashTurnovers()
	minWashTurnovers := 200
	stageWash := "Wash"
	ticker := time.NewTicker(time.Second * 3)
	timer := time.NewTimer(time.Second * time.Duration(washTime))
	go DataGenerator(stageWash, ticker, int64(maxTemperature), int64(minWashTurnovers), maxWashTurnovers, turnOversStorage, waterTempStorage)
	<-timer.C
	ticker.Stop()
	log.Println(stageWash, "finished!")

	// Run rinse
	rinseTime := config.GetRinseTime()
	maxRinseTurnovers := config.GetRinseTurnovers()
	minRinseTurnovers := maxRinseTurnovers - 100
	stageRinse := "Rinse"
	ticker = time.NewTicker(time.Second * 3)
	timer = time.NewTimer(time.Second * time.Duration(rinseTime))
	go DataGenerator(stageRinse, ticker, int64(maxTemperature), int64(minRinseTurnovers), maxRinseTurnovers, turnOversStorage, waterTempStorage)
	<-timer.C
	ticker.Stop()
	log.Println(stageRinse, "finished!")

	// Run spin
	spinTime := config.GetSpinTime()
	maxSpinTurnovers := config.GetSpinTurnovers()
	minSpinTurnovers := maxSpinTurnovers - 100
	stageSpin := "Spin"
	ticker = time.NewTicker(time.Second * 3)
	timer = time.NewTimer(time.Second * time.Duration(spinTime))
	go DataGenerator(stageSpin, ticker, int64(maxTemperature), int64(minSpinTurnovers), maxSpinTurnovers, turnOversStorage, waterTempStorage)
	<-timer.C
	ticker.Stop()
	log.Println(stageSpin, "finished!")

	firstStep <- struct{}{}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
