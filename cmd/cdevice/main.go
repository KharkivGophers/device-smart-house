package main

import (
	"net"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/vpakhuchyi/device-smart-house/models"
)

//Constants fr dialup setup
const (
	HOST = "localhost"
	PORT = "3030"
	TYPE = "tcp"
)

func main() {
	ln, _ := net.Listen(TYPE, HOST+":"+PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Errorln(err)
		}
		go handleRequest(conn)
	}

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
