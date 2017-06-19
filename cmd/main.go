package main

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/device"
	"github.com/KharkivGophers/device-smart-house/models"
)

func init() {
	cTop = make(chan models.FridgeGenerData, 100)
	cBot = make(chan models.FridgeGenerData, 100)
	reqChan = make(chan models.Request)
	conf = config.GetConfig()
}

func main() {
	config.Init(connTypeConf, hostConf, portConf)
	wg.Add(1)
	go device.RunDataGenerator(conf, cBot, cTop, &wg)
	go device.RunDataCollector(conf, cBot, cTop, reqChan, &wg)
	go device.DataTransfer(conf, reqChan)
	wg.Wait()

}
