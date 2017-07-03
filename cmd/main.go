package main

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/tcp/connectionupdate"
	"github.com/KharkivGophers/device-smart-house/devices/fridge"
	"sync"
	"log"
)

//
//func (c *Closer) Close() {
//	select {
//	case <- c.Control:
//	default:
//		close(c.Control)
//	}
//}

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
			log.Print(r)
		}
	} ()
	conf.Init(configConnParams.ConnTypeConf, configConnParams.HostConf, configConnParams.PortConf, &Wg)

	Wg.Add(1)
	go fridge.RunDataGenerator(conf, collectData.CBot, collectData.CTop, &Wg)
	go fridge.RunDataCollector(conf, collectData.CBot, collectData.CTop, collectData.ReqChan, &Wg)
	go fridge.DataTransfer(conf, collectData.ReqChan, &Wg)

	Wg.Wait()
	log.Println("DONE")
}