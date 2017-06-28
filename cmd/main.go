package main

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/connection"
	"github.com/KharkivGophers/device-smart-house/fridge"
)

func main() {
	collectData := models.CollectFridgeData{
		CTop: make(chan models.FridgeGenerData, 100), // First Cam
		CBot: make(chan models.FridgeGenerData, 100), // Second Cam
		ReqChan: make(chan models.Request),
	}

	configConnParams := models.ConfigConnParams{
		ConnTypeConf: "tcp",
		HostConf: connection.GetEnvCenter("CENTER_PORT_3000_TCP_ADDR"),
		PortConf: "3000",
	}

	var conf *config.DevConfig
	conf = config.GetConfig()

	config.Init(configConnParams.ConnTypeConf, configConnParams.HostConf, configConnParams.PortConf)
	collectData.Wg.Add(1)
	go fridge.RunDataGenerator(conf, collectData.CBot, collectData.CTop, &collectData.Wg)
	go fridge.RunDataCollector(conf, collectData.CBot, collectData.CTop, collectData.ReqChan, &collectData.Wg)
	go fridge.DataTransfer(conf, collectData.ReqChan)
	collectData.Wg.Wait()
}
