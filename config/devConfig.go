package config

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"os"

	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/vpakhuchyi/device-smart-house/models"
)

type DevConfig struct {
	sync.Mutex
	turned      bool
	collectFreq int64
	sendFreq    int64
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

func (d *DevConfig) updateConfig(c models.Config) {
	log.Warningln(c)
	d.turned = c.TurnedOn
	log.Warningln("c.TurnedOn", c.TurnedOn)
	d.sendFreq = c.SendFreq
	log.Warningln("c.SendFreq", c.SendFreq, reflect.TypeOf(c.SendFreq))
	d.collectFreq = c.CollectFreq
	log.Warningln("c.CollectFreq", c.CollectFreq, reflect.TypeOf(c.CollectFreq))
	log.Println("Config updated")
}

func askConfig(conn *net.Conn) models.Config {
	args := os.Args[1:]

	var req models.Request
	var resp models.Config
	req = models.Request{
		Action: "config",
		Meta: models.Metadata{
			Type: args[0],
			Name: args[1],
			MAC:  args[2]},
	}
	log.Warningln(req)
	err := json.NewEncoder(*conn).Encode(&req)
	checkError("askConfig Encode JSON", err)

	err = json.NewDecoder(*conn).Decode(&resp)
	checkError("askConfig Decode JSON", err)
	return resp
}

func listenConfigf(devConfig *DevConfig, conn *net.Conn, sendFreqChan chan int64,
	collectFreqChan chan int64, turnedOnChan chan bool) {

	for {
		var config interface{}
		err := json.NewDecoder(*conn).Decode(&config)
		checkError("receiveConfig Decode JSON", err)
		log.Infoln(config)
		for k, v := range config.(map[string]interface{}) {
			log.Infoln("for range k, v:", k, v)
			switch k {
			case "sendFreq":
				sendFreqChan <- int64(v.(float64))
			case "collectFreq":
				collectFreqChan <- int64(v.(float64))
			case "turnedOn":
				turnedOnChan <- v.(bool)
			default:
				log.Println("default case in switch: listenConfig")

			}
		}
		log.Println("listenConfigf: config have been received")

		// resp.Descr = "Config have been received"
		// resp.Status = 200
		// err = json.NewEncoder(*conn).Encode(&resp)
		// checkError("receiveConfig Encode JSON", err)

	}
}
func Init(connType string, host string, port string, sendFreqChan chan int64, collectFreqChan chan int64, turnedOnChan chan bool) {
	config := GetConfig()
	var reconnect *time.Ticker

	conn, err := net.Dial(connType, host+":"+port)
	for err != nil {
		log.Errorln(err)
		reconnect = time.NewTicker(time.Second * 1)
		for range reconnect.C {
			conn, _ = net.Dial(connType, host+":"+port)
		}
	}

	config.updateConfig(askConfig(&conn))
	go listenConfigf(config, &conn, sendFreqChan, collectFreqChan, turnedOnChan)
}

func checkError(desc string, err error) error {
	if err != nil {
		log.Errorln(desc, err)
		return err
	}
	return nil
}
