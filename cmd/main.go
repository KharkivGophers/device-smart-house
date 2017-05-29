package main

import (
	"github.com/vpakhuchyi/device-smart-house/config"
	"github.com/vpakhuchyi/device-smart-house/device"
	"github.com/vpakhuchyi/device-smart-house/models"
)

func init() {
	cTop = make(chan models.FridgeGenerData, 10)
	cBot = make(chan models.FridgeGenerData, 10)
	reqChan = make(chan models.Request)
	stop = make(chan struct{})
	start = make(chan struct{})
	conf = config.GetConfig()
}

func main() {
	config.Init(connTypeConf, hostConf, portConf)
	// go func() {
	// 	for {
	// 		animation()
	// 	}
	// }()
	wg.Add(3)
	go device.RunDataGenerator(conf, cBot, cTop, stop, start)
	go device.RunDataCollector(conf, cBot, cTop, reqChan, stop, start)
	go device.DataTransfer(conf, reqChan)
	wg.Wait()
}
