package main

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/device"
	"github.com/KharkivGophers/device-smart-house/models"
)

func main() {

	collectData := models.CollectData{
		CTop: make(chan models.FridgeGenerData, 100), // First Cam
		CBot: make(chan models.FridgeGenerData, 100), // Second Cam
		ReqChan: make(chan models.Request),
	}

	configConnParams := models.ConfigConnParams{
		ConnTypeConf: "tcp",
		HostConf: device.GetEnvCenter("CENTER_PORT_3000_TCP_ADDR"),
		PortConf: "3000",
	}

	var conf *config.DevConfig
	conf = config.GetConfig()

	config.Init(configConnParams.ConnTypeConf, configConnParams.HostConf, configConnParams.PortConf)
	collectData.Wg.Add(1)
	go device.RunDataGenerator(conf, collectData.CBot, collectData.CTop, &collectData.Wg)
	go device.RunDataCollector(conf, collectData.CBot, collectData.CTop, collectData.ReqChan, &collectData.Wg)
	go device.DataTransfer(conf, collectData.ReqChan)
	collectData.Wg.Wait()
}
