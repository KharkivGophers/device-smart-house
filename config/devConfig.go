package config

import (
	"encoding/json"
	"net"
	"sync"
	"os"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/error"
	log "github.com/Sirupsen/logrus"
)

type DevConfig struct {
	sync.Mutex
	turned      bool
	collectFreq int64
	sendFreq    int64
	subsPool    map[string]chan struct{}
}

func (d *DevConfig) AddSubIntoPool(key string, value chan struct{}) {
	d.Mutex.Lock()
	d.subsPool[key] = value
	d.Mutex.Unlock()
}

func (d *DevConfig) RemoveSubFromPool(key string) {
	d.Mutex.Lock()
	delete(d.subsPool, key)
	d.Mutex.Unlock()
}

func askConfig(conn net.Conn) models.Config {
	args := os.Args[1:]
	log.Warningln("Type:"+"["+args[0]+"];", "Name:"+"["+args[1]+"];", "MAC:"+"["+args[2]+"]")
	if len(args) < 3 {
		panic("Incorrect devices's information")
	}

	var req models.Request
	var resp models.Config
	req = models.Request{
		Action: "config",
		Meta: models.Metadata{
			Type: args[0],
			Name: args[1],
			MAC:  args[2]},
	}
	err := json.NewEncoder(conn).Encode(req)
	error.CheckError("askConfig(): Encode JSON", err)

	err = json.NewDecoder(conn).Decode(&resp)
	error.CheckError("askConfig(): Decode JSON", err)

	if err != nil && resp.IsEmpty() {
		panic("Connection has been closed by center")
	}

	return resp
}

func listenConfig(devConfig *DevConfig, conn net.Conn) {
	var resp models.Response
	var config models.Config

	err := json.NewDecoder(conn).Decode(&config)
	error.CheckError("listenConfig(): Decode JSON", err)

	resp.Descr = "Config have been received"

	devConfig.updateConfig(config)
	go publishConfig(devConfig)

	err = json.NewEncoder(conn).Encode(&resp)
	error.CheckError("listenConfig(): Encode JSON", err)

}

func publishConfig(d *DevConfig) {
	for _, v := range d.subsPool {
		v <- struct{}{}
	}
}

func (d *DevConfig) updateConfig(c models.Config) {
	d.turned = c.TurnedOn
	d.sendFreq = c.SendFreq
	log.Warningln("SendFreq: ", d.sendFreq)
	d.collectFreq = c.CollectFreq
	log.Warningln("CollectFreq: ", d.collectFreq)

	switch d.turned {
	case false:
		log.Warningln("ON PAUSE")
	case true:
		log.Warningln("WORKING")
	}
}

func (dc *DevConfig) Init(connType string, host string, port string) {
	conn, err := net.Dial(connType, host+":"+port)
	for err != nil {
		panic("Can't connect to the server: " + host + ":" + port)
	}

	dc.updateConfig(askConfig(conn))
	go func() {
		for {
			listenConfig(dc, conn)
		}
	}()
}
