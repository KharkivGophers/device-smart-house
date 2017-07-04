package config

import (
	"encoding/json"
	"net"
	"sync"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/error"
	"github.com/KharkivGophers/device-smart-house/tcp/connectionconfig"
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

func listenConfig(devConfig *DevConfig, conn net.Conn) {
	var resp models.Response
	var config models.Config

	err := json.NewDecoder(conn).Decode(&config)
	if err != nil {
		panic("No config found!")
	}
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

func (dc *DevConfig) Init(connType string, host string, port string, wg *sync.WaitGroup, c *models.Control) {
	conn, err := net.Dial(connType, host+":"+port)
	wg.Add(1)
	for err != nil {
		log.Error("Can't connect to the server: " + host + ":" + port)
		panic("No center found!")
	}

	dc.updateConfig(connectionconfig.AskConfig(conn))
	go func() {
		for {
			defer func() {
				if r := recover(); r != nil {
					c.Close()
					wg.Done()
					log.Error("Initialization Failed")
				}
			} ()
			listenConfig(dc, conn)
		}
	}()
}
