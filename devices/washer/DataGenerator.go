package washer

import (
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
	"github.com/KharkivGophers/device-smart-house/models"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"time"
)

// DataGenerator generates pseudo-random numbers
func DataGenerator(stage string, ticker *time.Ticker, maxTemperature int64, maxTurnovers int64, turnOversStorage chan<- models.GenerateWasherData,
	waterTempStorage chan<- models.GenerateWasherData) {

	log.Info(stage, " -------- S T A R T E D!")
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
	stageWash := "W A S H"
	ticker := time.NewTicker(time.Second * 3)
	timer := time.NewTimer(time.Second * time.Duration(washTime))
	go DataGenerator(stageWash, ticker, int64(maxTemperature), maxWashTurnovers, turnOversStorage, waterTempStorage)
	<-timer.C
	ticker.Stop()
	log.Info(stageWash, " -------- F I N I S H E D!")

	// Run rinse
	rinseTime := config.GetRinseTime()
	maxRinseTurnovers := config.GetRinseTurnovers()
	stageRinse := "R I N S E"
	ticker = time.NewTicker(time.Second * 3)
	timer = time.NewTimer(time.Second * time.Duration(rinseTime))
	go DataGenerator(stageRinse, ticker, int64(maxTemperature), maxRinseTurnovers, turnOversStorage, waterTempStorage)
	<-timer.C
	ticker.Stop()
	log.Info(stageRinse, " -------- F I N I S H E D!")

	// Run spin
	spinTime := config.GetSpinTime()
	maxSpinTurnovers := config.GetSpinTurnovers()
	stageSpin := "S P I N"
	ticker = time.NewTicker(time.Second * 3)
	timer = time.NewTimer(time.Second * time.Duration(spinTime))
	go DataGenerator(stageSpin, ticker, int64(maxTemperature), maxSpinTurnovers, turnOversStorage, waterTempStorage)
	<-timer.C
	ticker.Stop()
	log.Info(stageSpin, " -------- F I N I S H E D!")
	log.Warn("W A S H I N G --- M A C H I N E --- F I N I S H E D")

	firstStep <- struct{}{}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}