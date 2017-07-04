package main

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/tcp/connectionupdate"
	"github.com/KharkivGophers/device-smart-house/devices/fridge"
	"sync"
	log "github.com/Sirupsen/logrus"
)

func main() {
	var Wg sync.WaitGroup

	collectData := models.CollectFridgeData{
		CTop: make(chan models.FridgeGenerData, 100), // First Cam
		CBot: make(chan models.FridgeGenerData, 100), // Second Cam
		ReqChan: make(chan models.Request),
	}

	configConnParams := models.ConfigConnParams{
		ConnTypeConf: "tcp",
		HostConf:     connectionupdate.GetEnvCenter("CENTER_PORT_3000_TCP_ADDR"),
		PortConf:     "3000",
	}

	conf := config.NewConfig()
	defer func() {
		if r := recover(); r != nil {
		}
	} ()

	control := &models.Control{make(chan struct{})}
	conf.Init(configConnParams.ConnTypeConf, configConnParams.HostConf, configConnParams.PortConf, &Wg, control)

	go fridge.RunDataGenerator(conf, collectData.CBot, collectData.CTop, &Wg, control)
	go fridge.RunDataCollector(conf, collectData.CBot, collectData.CTop, collectData.ReqChan, &Wg, control)
	go fridge.DataTransfer(conf, collectData.ReqChan, &Wg, control)

	Wg.Wait()
	log.Info("DONE")
}