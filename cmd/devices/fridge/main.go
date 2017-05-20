package main

import (
	"encoding/json"
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/vpakhuchyi/device-smart-house/devices/fridge"
	"github.com/vpakhuchyi/device-smart-house/models"
)

var (
	cBot                 chan device.FridgeGenerData
	cTop                 chan device.FridgeGenerData
	reqChan              chan models.Request
	stop                 chan struct{}
	start                chan struct{}
	wg                   sync.WaitGroup
	HandleRequestCounter int

	//BreakerVar singletone; gives accsess to On or Off device
	config = models.GetConfig()
)

//Constants for dialup setup
var (
	hostIn     = "localhost"
	portIn     = "3000"
	connTypeIn = "tcp"

	hostOut     = "localhost"
	portOut     = "3030"
	connTypeOut = "tcp"
	//
)

func init() {
	cTop = make(chan device.FridgeGenerData, 10)
	cBot = make(chan device.FridgeGenerData, 10)
	reqChan = make(chan models.Request)
	stop = make(chan struct{})
	start = make(chan struct{})

	//Device must be TurnedON from the beginning
	config.SetTurned(true)
	config.SetCollectFreq(1)
	config.SetSendFreq(5)
}

func main() {
	//Listens for request from centre
	ln, err := net.Listen(connTypeIn, hostIn+":"+portIn)
	if err != nil {
		log.Errorln("TCP Config")
		log.Errorln(err)
	}
	wg.Add(1)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Errorln(err)
			}
			go handleRequest(conn, config)
		}
	}()

	//----TCP for sending requests to the center
	conn, err := net.Dial(connTypeOut, hostOut+":"+portOut)
	if err != nil {
		log.Errorln("TCP Send")
		log.Errorln(err)

	}

	go device.RunDataGenerator(config, cBot, cTop, stop, start)
	go device.RunDataCollector(config, cBot, cTop, reqChan, stop, start)
	//maybe without go?
	go device.DataTransfer(config, reqChan, conn)
	wg.Wait()

}

func handleRequest(conn net.Conn, config *models.DevConfig) {
	log.Warningln("hadleRequest intro")
	HandleRequestCounter++
	log.Warningln(HandleRequestCounter)

	var resp models.Response
	var req models.ConfigRequest

	err := json.NewDecoder(conn).Decode(&req)
	if err != nil {
		log.Errorln(err)
	}

	log.Warningln(req.Turned)
	log.Warningln(req)
	if req.Turned == false {
		stop <- struct{}{}
	}

	if req.Turned == true {
		start <- struct{}{}
	}

	config.SetTurned(req.Turned)
	config.SetCollectFreq(req.CollectFreq)
	config.SetSendFreq(req.SendFreq)
	log.Warningln("hadleRequest setuped DevConfig")

	resp.Descr = "New config accepted"
	resp.Status = 200

	err = json.NewEncoder(conn).Encode(&resp)
	if err != nil {
		log.Errorln(err, " Encode")
	}

	log.Println("response: ", resp)
	log.Println("handleRequest: [done]")
	log.Warningln("hadleRequest out")
}
