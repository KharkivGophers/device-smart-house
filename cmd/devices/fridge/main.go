package main

import (
	"sync"

	"github.com/vpakhuchyi/device-smart-house/config"
	"github.com/vpakhuchyi/device-smart-house/devices/fridge"
	"github.com/vpakhuchyi/device-smart-house/models"
)

var (
	cBot    chan device.FridgeGenerData
	cTop    chan device.FridgeGenerData
	reqChan chan models.Request
	stop    chan struct{}
	start   chan struct{}
	wg      sync.WaitGroup
	conf    *config.DevConfig
)

var (
	hostConf     = "localhost"
	portConf     = "3000"
	connTypeConf = "tcp"
)

func init() {
	cTop = make(chan device.FridgeGenerData, 10)
	cBot = make(chan device.FridgeGenerData, 10)
	reqChan = make(chan models.Request)
	stop = make(chan struct{})
	start = make(chan struct{})

	conf = config.GetConfig()
}

func main() {
	config.ConfigInit(connTypeConf, hostConf, portConf)

	wg.Add(1)
	go device.RunDataGenerator(conf, cBot, cTop, stop, start)
	go device.RunDataCollector(conf, cBot, cTop, reqChan, stop, start)
	go device.DataTransfer(conf, reqChan)
	wg.Wait()

}
