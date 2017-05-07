package main

import (
	"sync"

	"github.com/vpakhuchyi/device-smart-house/devices/fridge"
	"github.com/vpakhuchyi/device-smart-house/models"
)

var (
	cBot    chan float32
	cTop    chan float32
	reqChan chan *models.Request
	wg      sync.WaitGroup

	//BreakerVar singletone; gives accsess to On or Off device
	breakerVar = models.GetBreaker()
)

//Constants fr dialup setup
const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func init() {
	cTop = make(chan float32, 2)
	cBot = make(chan float32, 2)
	reqChan = make(chan *models.Request)

	//Device must be TurnedON from the beginning
	breakerVar.SetTurned(true)
}

func main() {
	//Listens for request from centre (it may contain config file)
	// ln, _ := net.Listen(TYPE, HOST+":"+PORT)

	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		log.Errorln(err)
	// 	}
	// }

	// swap our device state every few secs - for debugging
	// t := time.NewTicker(time.Second * 5)
	// go func() {
	// 	for range t.C {
	// 		swap(BreakerVar)
	// 		log.Warnln("Status has changed")
	// 	}
	// }()

	wg.Add(1)
	go device.DataGenerator(breakerVar, cBot, cTop)
	wg.Add(1)
	go device.DataCollector(breakerVar, cBot, cTop, reqChan)
	wg.Add(1)
	go device.DataTransfer(breakerVar, reqChan)
	wg.Wait()
}

// swap func change state of device (TurnON or TurnOFF) - for debugging
// func swap(br *models.Breaker) {
// 	if br.GetTurned() == true {
// 		br.SetTurned(false)
// 	} else {
// 		br.SetTurned(true)
// 	}
// }
