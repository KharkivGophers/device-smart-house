package main

import (
	"github.com/KharkivGophers/device-smart-house/devices/washer"
	"github.com/KharkivGophers/device-smart-house/config/fridgeconfig"
	"github.com/KharkivGophers/device-smart-house/devices/fridge"
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
	"github.com/KharkivGophers/device-smart-house/models"
)

func startWasher(connType string, host string, port string, control *models.Control, args []string) {
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
	nextStep := make(chan struct{})
	firstStep := make(chan struct{})

	go washer.DataTransfer(washerConfig, collectWasherData.RequestStorage, control)

	go func() {
		for {
			select {
			case <-firstStep:
				go washerConfig.SendWasherRequests(connType, host, port, control, args, nextStep)
				<-nextStep
				go washer.RunDataGenerator(washerConfig, collectWasherData.TurnoversStorage, collectWasherData.TemperatureStorage, control, firstStep)
				go washer.RunDataCollector(washerConfig, collectWasherData.TurnoversStorage, collectWasherData.TemperatureStorage, collectWasherData.RequestStorage)
			}
		}
	}()
	firstStep <- struct{}{}
}

func startFridge(connType string, host string, port string, control *models.Control, args []string) {
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
	fridgeConfig.RequestFridgeConfig(connType, host, port, control, args)

	go fridge.RunDataGenerator(fridgeConfig, collectFridgeData.CBot, collectFridgeData.CTop, control)
	go fridge.RunDataCollector(fridgeConfig, collectFridgeData.CBot, collectFridgeData.CTop, collectFridgeData.ReqChan, control)
	go fridge.DataTransfer(fridgeConfig, collectFridgeData.ReqChan, control)
}