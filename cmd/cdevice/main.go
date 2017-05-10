package main

import (
	"net"
	"sync"

	"encoding/json"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vpakhuchyi/device-smart-house/models"
)

//Constants fr dialup setup
const (
	HOST = "localhost"
	PORT = "3030"
	TYPE = "tcp"
)

var wg sync.WaitGroup

func main() {
	ln, _ := net.Listen(TYPE, HOST+":"+PORT)

	wg.Add(1)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Errorln(err)
			}
			go handleRequest(conn)
		}
	}()

	time.Sleep(time.Second * 10)
	SendRequest("localhost:8080", true, 2, 4)
	wg.Wait()
}

func handleRequest(conn net.Conn) {
	var req models.Request
	var resp models.Response

	err := json.NewDecoder(conn).Decode(&req)
	if err != nil {
		log.Errorln(err)
	}

	resp.Descr = req.Action + ":completed"
	resp.Status = 200

	err = json.NewEncoder(conn).Encode(&resp)
	if err != nil {
		log.Errorln(err, " Encode")
	}

	log.Println("response: ", resp)
	log.Println("handleRequest: [done]")

}

func SendRequest(host string, turned bool, collFreq int, sendFreq int) {
	var req models.ConfigRequest
	var resp models.Response

	req.Turned = turned
	req.CollectFreq = collFreq
	req.SendFreq = sendFreq

	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Errorln(err)
	}

	err = json.NewEncoder(conn).Encode(&req)
	if err != nil {
		log.Errorln(err, " Encode")
	}

	err = json.NewDecoder(conn).Decode(&resp)
	if err != nil {
		log.Errorln(err)
	}

	log.Println("response: ", resp)
	log.Println("handleRequest: [done]")

}
