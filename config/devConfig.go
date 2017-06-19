package config

import (
	"encoding/json"
	"net"
	"sync"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/KharkivGophers/device-smart-house/models"
)

type DevConfig struct {
	sync.Mutex
	turned      bool
	collectFreq int64
	sendFreq    int64
	subsPool    map[string]chan struct{}
}

var config *DevConfig
var once sync.Once

func GetConfig() *DevConfig {
	once.Do(func() {
		config = &DevConfig{}
		config.subsPool = make(map[string]chan struct{})
	})
	return config
}

func (d *DevConfig) SetTurned(b bool) {
	d.Mutex.Lock()
	d.turned = b
	defer d.Mutex.Unlock()
}

func (d *DevConfig) GetTurned() bool {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.turned
}

func (d *DevConfig) GetCollectFreq() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.collectFreq
}

func (d *DevConfig) GetSendFreq() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.sendFreq
}

func (d *DevConfig) SetCollectFreq(b int64) {
	d.Mutex.Lock()
	d.collectFreq = b
	d.Mutex.Unlock()

}

func (d *DevConfig) SetSendFreq(b int64) {
	d.Mutex.Lock()
	d.sendFreq = b
	d.Mutex.Unlock()

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

func askConfig(conn net.Conn) models.Config {
	args := os.Args[1:]
	log.Warningln("Type:"+"["+args[0]+"];", "Name:"+"["+args[1]+"];", "MAC:"+"["+args[2]+"]")
	if len(args) < 3 {
		panic("Incorrect device's information")
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
	checkError("askConfig(): Encode JSON", err)

	err = json.NewDecoder(conn).Decode(&resp)
	checkError("askConfig(): Decode JSON", err)

	if err != nil && resp.IsEmpty() {
		panic("Connection has been closed by center")
	}

	return resp

}

func listenConfig(devConfig *DevConfig, conn net.Conn) {
	var resp models.Response
	var config models.Config

	err := json.NewDecoder(conn).Decode(&config)
	checkError("listenConfig(): Decode JSON", err)

	resp.Descr = "Config have been received"

	devConfig.updateConfig(config)
	go publishConfig(devConfig)

	err = json.NewEncoder(conn).Encode(&resp)
	checkError("listenConfig(): Encode JSON", err)

}

func publishConfig(d *DevConfig) {
	for _, v := range d.subsPool {
		v <- struct{}{}
	}
}

func Init(connType string, host string, port string) {
	// var times int
	config := GetConfig()
	conn, err := net.Dial(connType, host+":"+port)
	for err != nil {
		panic("Can't connect to the server")
	}

	config.updateConfig(askConfig(conn))
	go func() {
		for {
			listenConfig(config, conn)
		}
	}()

}

func checkError(desc string, err error) error {
	if err != nil {
		log.Errorln(desc, err)
		return err
	}
	return nil
}
