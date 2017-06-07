package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/vpakhuchyi/device-smart-house/config"
	"github.com/vpakhuchyi/device-smart-house/device"
	"github.com/vpakhuchyi/device-smart-house/models"
)

func init() {
	cTop = make(chan models.FridgeGenerData, 100)
	cBot = make(chan models.FridgeGenerData, 100)
	reqChan = make(chan models.Request)
	conf = config.GetConfig()
}

func main() {
	runtime.SetBlockProfileRate(10)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	config.Init(connTypeConf, hostConf, portConf)
	wg.Add(1)
	go device.RunDataGenerator(conf, cBot, cTop, &wg)
	go device.RunDataCollector(conf, cBot, cTop, reqChan, &wg)
	go device.DataTransfer(conf, reqChan, &wg)
	wg.Wait()
}

func catchPanic() {

}
