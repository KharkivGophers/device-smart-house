package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/vpakhuchyi/device-smart-house/devices/fridge"
	"github.com/vpakhuchyi/device-smart-house/models"
)

var (
	cBot    chan float32
	cTop    chan float32
	reqChan chan models.Request
	wg      sync.WaitGroup

	//BreakerVar singletone; gives accsess to On or Off device
	breakerVar = models.GetBreaker()
)

func init() {
	cTop = make(chan float32, 2)
	cBot = make(chan float32, 2)
	reqChan = make(chan models.Request)

	//Device must be TurnedON from the beginning
	breakerVar.SetTurned(true)
}

func main() {
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Received a request from centre")
		})
	}()

	// swap our device state every few secs - for debugging
	// t := time.NewTicker(time.Second * 5)
	// go func() {
	// 	for range t.C {
	// 		swap(BreakerVar)
	// 		log.Warnln("Status has changed")
	// 	}
	// }()

	go device.DataGenerator(breakerVar, cBot, cTop)

	go device.DataCollector(breakerVar, cBot, cTop, reqChan)

	go device.DataTransfer(breakerVar, reqChan)
	http.ListenAndServe(":8000", nil)
}

// swap func change state of device (TurnON or TurnOFF) - for debugging
// func swap(br *models.Breaker) {
// 	if br.GetTurned() == true {
// 		br.SetTurned(false)
// 	} else {
// 		br.SetTurned(true)
// 	}
// }
