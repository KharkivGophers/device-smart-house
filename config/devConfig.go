package config

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vpakhuchyi/device-smart-house/models"
)

type DevConfig struct {
	sync.Mutex
	turned      bool
	collectFreq int
	sendFreq    int
}

var config *DevConfig
var once sync.Once

func GetConfig() *DevConfig {
	once.Do(func() {
		config = &DevConfig{}
	})
	return config
}

func (d *DevConfig) SetTurned(b bool) {
	d.Mutex.Lock()
	d.turned = b
	d.Mutex.Unlock()
}

func (d *DevConfig) GetTurned() bool {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.turned
}

func (d *DevConfig) GetCollectFreq() int {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.collectFreq
}

func (d *DevConfig) GetSendFreq() int {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.sendFreq
}

func (d *DevConfig) SetCollectFreq(b int) {
	d.Mutex.Lock()
	d.collectFreq = b
	d.Mutex.Unlock()

}

func (d *DevConfig) SetSendFreq(b int) {
	d.Mutex.Lock()
	d.sendFreq = b
	d.Mutex.Unlock()

}

func (d *DevConfig) UpdateConfig(c models.Config) {
	d.turned = c.Turned
	d.sendFreq = c.SendFreq
	d.collectFreq = c.CollectFreq
}

func AskConfig(conn *net.Conn) models.Config {
	var req models.Request
	var resp models.Config

	req = models.Request{
		Action: "config",
		Meta: models.Metadata{
			Type: "fridge",
			Name: "hladik0e31",
			MAC:  "00-15-E9-2B-99-3C"},
	}

	err := json.NewEncoder(*conn).Encode(&req)
	if err != nil {
		log.Errorln("Encode JSON", err)
	}

	err = json.NewDecoder(*conn).Decode(&resp)
	if err != nil {
		log.Errorln("Decode JSON", err)
	}
	return resp
}

func listenConfig(conn *net.Conn) models.Config {
	var req models.Config
	var resp models.Response

	err := json.NewDecoder(*conn).Decode(&req)
	if err != nil {
		log.Errorln("Decode JSON", err)
	}

	resp.Status = 200
	resp.Descr = "Config has been accepted"

	err = json.NewEncoder(*conn).Encode(&resp)
	if err != nil {
		log.Errorln("Encode JSON", err)
	}
	return req
}

func ConfigInit(connType string, host string, port string) {
	config := GetConfig()
	var reconnect *time.Ticker
	defer reconnect.Stop()

	conn, err := net.Dial(connType, host+":"+port)
	for err != nil {
		log.Errorln(err)
		reconnect = time.NewTicker(time.Second * 1)
		for range reconnect.C {
			conn, _ = net.Dial(connType, host+":"+port)
		}
	}

	config.UpdateConfig(AskConfig(&conn))

	go func(conn *net.Conn, conf *DevConfig) {
		for {
			config.UpdateConfig(listenConfig(conn))
			log.Warningln("Listens for config in loop")
		}
	}(&conn, config)
}
