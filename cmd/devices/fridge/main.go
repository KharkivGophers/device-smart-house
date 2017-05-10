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
	cBot    chan device.FridgeGenerData
	cTop    chan device.FridgeGenerData
	reqChan chan *models.Request
	wg      sync.WaitGroup

	//BreakerVar singletone; gives accsess to On or Off device
	config = models.GetConfig()
)

//Constants fr dialup setup
const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func init() {
	cTop = make(chan device.FridgeGenerData, 2)
	cBot = make(chan device.FridgeGenerData, 2)
	reqChan = make(chan *models.Request)

	//Device must be TurnedON from the beginning
	config.SetTurned(true)
	config.SetCollectFreq(1)
	config.SetSendFreq(5)

}

func main() {
	//Listens for request from centre (it may contain config file)
	ln, _ := net.Listen(TYPE, HOST+":"+PORT)
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

	go device.DataGenerator(config, cBot, cTop)
	go device.DataCollector(config, cBot, cTop, reqChan)
	go device.DataTransfer(config, reqChan)
	wg.Wait()
}

func handleRequest(conn net.Conn, config *models.DevConfig) {
	var resp models.Response
	var req models.ConfigRequest
	err := json.NewDecoder(conn).Decode(&req)
	if err != nil {
		log.Errorln(err)
	}

	config.SetTurned(req.Turned)
	config.SetCollectFreq(req.CollectFreq)
	config.SetSendFreq(req.SendFreq)

	resp.Descr = "New config accepted"
	resp.Status = 200

	err = json.NewEncoder(conn).Encode(&resp)
	if err != nil {
		log.Errorln(err, " Encode")
	}

	log.Println("response: ", resp)
	log.Println("handleRequest: [done]")
}
