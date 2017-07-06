package washerconfig

import (
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/error"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"sync"
	"net"
)

type DevWasherConfig struct {
	sync.Mutex
	TurnedOn    	bool
	WashTime		int64
	WashTurnovers 	int64
	RinseTime		int64
	RinseTurnovers	int64
	SpinTime		int64
	SpinTurnovers	int64
	subsPool    map[string]chan struct{}
}

func (d *DevWasherConfig) AddSubIntoPool(key string, value chan struct{}) {
	d.Mutex.Lock()
	d.subsPool[key] = value
	d.Mutex.Unlock()
}

func (d *DevWasherConfig) RemoveSubFromPool(key string) {
	d.Mutex.Lock()
	delete(d.subsPool, key)
	d.Mutex.Unlock()
}

func listenWasherConfig(devConfig *DevWasherConfig, conn net.Conn) {
	var resp models.Response
	var config models.WasherConfig

	err := json.NewDecoder(conn).Decode(&config)
	if err != nil {
		panic("No config found!")
	}
	error.CheckError("listenConfig(): Decode JSON", err)

	resp.Descr = "Config have been received"
	devConfig.updateWasherConfig(config)
	go publishWasherConfig(devConfig)

	err = json.NewEncoder(conn).Encode(&resp)
	error.CheckError("listenConfig(): Encode JSON", err)
}

func publishWasherConfig(d *DevWasherConfig) {
	for _, v := range d.subsPool {
		v <- struct{}{}
	}
}

func (d *DevWasherConfig) updateWasherConfig(c models.WasherConfig) {
	d.TurnedOn = c.TurnedOn
	log.Warningln("TurnedOn: ", d.TurnedOn)

	d.WashTime = c.WashTime
	log.Warningln("WashTime: ", d.WashTime)

	d.WashTurnovers = c.WashTurnovers
	log.Warningln("WashTurnovers: ", d.WashTurnovers)

	d.RinseTime = c.RinseTime
	log.Warningln("RinseTime: ", d.RinseTime)

	d.RinseTurnovers = c.RinseTurnovers
	log.Warningln("RinseTurnovers: ", d.RinseTurnovers)

	d.SpinTime = c.SpinTime
	log.Warningln("SpinTime: ", d.SpinTime)

	d.SpinTurnovers = c.SpinTurnovers
	log.Warningln("SpinTurnovers: ", d.SpinTurnovers)

	switch d.TurnedOn {
	case false:
		log.Warningln("ON PAUSE")
	case true:
		log.Warningln("WORKING")
	}
}

func (washer *DevWasherConfig) RequestWasherConfig(connType string, host string, port string, c *models.Control, args []string) {

	conn, err := net.Dial(connType, host+":"+port)
	for err != nil {
		log.Error("Can't connect to the server: " + host + ":" + port)
		panic("No center found!")
	}

	var response models.WasherConfig
	var request models.WasherRequest

	request = models.WasherRequest{
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

	washer.updateWasherConfig(response)

	go func() {
		for {
			defer func() {
				if r := recover(); r != nil {
					c.Close()
					log.Error("Initialization Failed")
				}
			} ()
			listenWasherConfig(washer, conn)
		}
	}()
}
