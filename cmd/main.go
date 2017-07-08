package main

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/config/fridgeconfig"
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
	"github.com/KharkivGophers/device-smart-house/devices/fridge"
	"github.com/KharkivGophers/device-smart-house/devices/washer"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/tcp/connectionupdate"
	log "github.com/Sirupsen/logrus"
)

func main() {
	configConnParams := models.ConfigConnParams{
		ConnTypeConf: "tcp",
		HostConf:     connectionupdate.GetEnvCenter("CENTER_PORT_3000_TCP_ADDR"),
		PortConf:     "3000",
	}

	newDevice := config.CreateDevice()
	newDeviceType := newDevice[0]
	control := &models.Control{make(chan struct{})}

	switch newDeviceType {
	case "washer":
		washerConfig := washerconfig.NewWasherConfig()

		collectWasherData := models.CollectWasherData{
			TurnoversStorage:   make(chan models.GenerateWasherData, 100),
			TemperatureStorage: make(chan models.GenerateWasherData, 100),
			RequestStorage:     make(chan models.WasherRequest),
		}
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		washerConfig.SendWasherRequests(configConnParams.ConnTypeConf, configConnParams.HostConf, configConnParams.PortConf, control, newDevice)

		go washer.RunDataGenerator(washerConfig, collectWasherData.TurnoversStorage, collectWasherData.TemperatureStorage, control)
		go washer.RunDataCollector(washerConfig, collectWasherData.TurnoversStorage, collectWasherData.TemperatureStorage, collectWasherData.RequestStorage, control)
		go washer.DataTransfer(washerConfig, collectWasherData.RequestStorage, control)

	default:
		fridgeConfig := fridgeconfig.NewFridgeConfig()

		collectFridgeData := models.CollectFridgeData{
			CTop:    make(chan models.FridgeGenerData, 100), // First Camera
			CBot:    make(chan models.FridgeGenerData, 100), // Second Camera
			ReqChan: make(chan models.FridgeRequest),
		}

		defer func() {
			if r := recover(); r != nil {
			}
		}()
		fridgeConfig.RequestFridgeConfig(configConnParams.ConnTypeConf, configConnParams.HostConf, configConnParams.PortConf, control, newDevice)

		go fridge.RunDataGenerator(fridgeConfig, collectFridgeData.CBot, collectFridgeData.CTop, control)
		go fridge.RunDataCollector(fridgeConfig, collectFridgeData.CBot, collectFridgeData.CTop, collectFridgeData.ReqChan, control)
		go fridge.DataTransfer(fridgeConfig, collectFridgeData.ReqChan, control)
	}

	control.Wait()
	log.Info("Device has been terminated due to the center's issue")
}
