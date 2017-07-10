package washerconfig

import (
	"encoding/json"
	"github.com/KharkivGophers/device-smart-house/error"
	"github.com/KharkivGophers/device-smart-house/models"
	log "github.com/Sirupsen/logrus"
	"net"
	"sync"
	"time"
)

type DevWasherConfig struct {
	sync.Mutex
	Temperature    int64
	WashTime       int64
	WashTurnovers  int64
	RinseTime      int64
	RinseTurnovers int64
	SpinTime       int64
	SpinTurnovers  int64
	subsPool       map[string]chan struct{}
}

func (d *DevWasherConfig) updateWasherConfig(c models.WasherConfig) {
	d.Temperature = c.Temperature

	d.WashTime = c.WashTime

	d.WashTurnovers = c.WashTurnovers

	d.RinseTime = c.RinseTime

	d.RinseTurnovers = c.RinseTurnovers

	d.SpinTime = c.SpinTime

	d.SpinTurnovers = c.SpinTurnovers

	log.Warn("New Configuration:")
	log.Warn("Temperature: ", d.Temperature, "; WashTime: ", d.WashTime, "; WashTurnovers: ", d.WashTurnovers,
		"; RinseTime: ", d.RinseTime, "; RinseTurnovers: ", d.RinseTurnovers, "; SpinTime: ", d.SpinTime,
		"; SpinTurnovers: ", d.SpinTurnovers)
}

func (washer *DevWasherConfig) RequestWasherConfig(connType string, host string, port string, args []string) models.WasherConfig {
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

	log.Println("Request:", request)

	err = json.NewDecoder(conn).Decode(&response)
	error.CheckError("askConfig(): Decode JSON", err)

	return response
}

func (washer *DevWasherConfig) SendWasherRequests(connType string, host string, port string, c *models.Control, args []string, nextStep chan struct{}) {

	ticker := time.NewTicker(time.Second)
	response := washer.RequestWasherConfig(connType, host, port, args)
	log.Println("Response:", response)

	for {
		select {
		case <-ticker.C:
			switch response.IsEmpty() {
			case true:
				log.Println("Response:", response)
				washer.RequestWasherConfig(connType, host, port, args)
			default:
				washer.updateWasherConfig(response)
				ticker.Stop()
				nextStep<-struct{}{}
				return
			}
		}
	}
}

//func (d *DevWasherConfig) AddSubIntoPool(key string, value chan struct{}) {
//	d.Mutex.Lock()
//	d.subsPool[key] = value
//	d.Mutex.Unlock()
//}
//
//func (d *DevWasherConfig) RemoveSubFromPool(key string) {
//	d.Mutex.Lock()
//	delete(d.subsPool, key)
//	d.Mutex.Unlock()
//}
//
//func publishWasherConfig(d *DevWasherConfig) {
//	for _, v := range d.subsPool {
//		v <- struct{}{}
//	}
//}