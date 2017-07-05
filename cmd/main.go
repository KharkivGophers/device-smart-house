package main

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/tcp/connectionupdate"
	"github.com/KharkivGophers/device-smart-house/devices/fridge"
	log "github.com/Sirupsen/logrus"
)

func main() {

	collectData := models.CollectFridgeData{
		CTop: make(chan models.FridgeGenerData, 100), // First Camera
		CBot: make(chan models.FridgeGenerData, 100), // Second Camera
		ReqChan: make(chan models.Request),
	}

	configConnParams := models.ConfigConnParams{
		ConnTypeConf: "tcp",
		HostConf:     connectionupdate.GetEnvCenter("CENTER_PORT_3000_TCP_ADDR"),
		PortConf:     "3000",
	}

	conf := config.NewConfig()
	defer func() {
		if r := recover(); r != nil {}} ()

	control := &models.Control{make(chan struct{})}
	conf.Init(configConnParams.ConnTypeConf, configConnParams.HostConf, configConnParams.PortConf, control)

	go fridge.RunDataGenerator(conf, collectData.CBot, collectData.CTop, control)
	go fridge.RunDataCollector(conf, collectData.CBot, collectData.CTop, collectData.ReqChan, control)
	go fridge.DataTransfer(conf, collectData.ReqChan, control)

	control.Wait()
	log.Info("DONE")
}