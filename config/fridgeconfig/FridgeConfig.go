package fridgeconfig

import (
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/error"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"sync"
	"net"
)

type DevFridgeConfig struct {
	sync.Mutex
	turned      bool
	collectFreq int64
	sendFreq    int64
	subsPool    map[string]chan struct{}
}

func (d *DevFridgeConfig) AddSubIntoPool(key string, value chan struct{}) {
	d.Mutex.Lock()
	d.subsPool[key] = value
	d.Mutex.Unlock()
}

func (d *DevFridgeConfig) RemoveSubFromPool(key string) {
	d.Mutex.Lock()
	delete(d.subsPool, key)
	d.Mutex.Unlock()
}

func listenConfig(devConfig *DevFridgeConfig, conn net.Conn) {
	var resp models.Response
	var config models.FridgeConfig

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

func publishConfig(d *DevFridgeConfig) {
	for _, v := range d.subsPool {
		v <- struct{}{}
	}
}

func (d *DevFridgeConfig) updateConfig(c models.FridgeConfig) {
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

func (fridge *DevFridgeConfig) RequestFridgeConfig(connType string, host string, port string, c *models.Control, args []string) {

	conn, err := net.Dial(connType, host+":"+port)
	for err != nil {
		log.Error("Can't connect to the server: " + host + ":" + port)
		panic("No center found!")
	}

	var response models.FridgeConfig
	var request models.FridgeRequest

	request = models.FridgeRequest{
		Action: "config",
		Meta: models.Metadata{
			Type: args[0],
			Name: args[1],
			MAC:  args[2]},
	}

	err = json.NewEncoder(conn).Encode(request)
	error.CheckError("askConfig(): Encode JSON", err)

	err = json.NewDecoder(conn).Decode(&response)
	error.CheckError("askConfig(): Decode JSON", err)

	if err != nil && response.IsEmpty() {
		panic("Connection has been closed by center")
	}

	fridge.updateConfig(response)

	go func() {
		for {
			defer func() {
				if r := recover(); r != nil {
					c.Close()
					log.Error("Initialization Failed")
				}
			} ()
			listenConfig(fridge, conn)
		}
	}()
}
