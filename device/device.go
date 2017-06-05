package device

import (
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vpakhuchyi/device-smart-house/config"
	"github.com/vpakhuchyi/device-smart-house/models"
)

//for data transfer
var (
	hostOut     = "localhost"
	portOut     = "3030"
	connTypeOut = "tcp"
)

//DataTransfer func sends reuest as JSON to the centre
func DataTransfer(config *config.DevConfig, reqChan chan models.Request) {
	conn := getDial(connTypeOut, hostOut, portOut)

	for {
		select {
		case r := <-reqChan:
			go send(r, conn)
		}
	}
}

//DataCollector func gathers data from DataGenerator
//and sends completed request's structures to the ReqChan channel
func DataCollector(ticker *time.Ticker, cBot <-chan models.FridgeGenerData, cTop <-chan models.FridgeGenerData, ReqChan chan models.Request, stopInner chan struct{}) {

	var mTop = make(map[int64]float32)
	var mBot = make(map[int64]float32)

	for {
		select {
		case <-stopInner:
			log.Warningln("DataCollector - stopInner-case")
			return
		case tv := <-cTop:
			mTop[tv.Time] = tv.Data
		case bv := <-cBot:
			mBot[bv.Time] = bv.Data
		case <-ticker.C:
			ReqChan <- constructReq(mTop, mBot)

			//Cleaning temp maps
			mTop = make(map[int64]float32)
			mBot = make(map[int64]float32)
		}

	}
}

//DataGenerator func generates pseudo-random data that represents device's behavior
func DataGenerator(ticker *time.Ticker, cBot chan<- models.FridgeGenerData, cTop chan<- models.FridgeGenerData, stopInner chan struct{}) {

	for {
		select {
		case <-ticker.C:
			cTop <- models.FridgeGenerData{Time: makeTimestamp(), Data: rand.Float32() * 10}
			cBot <- models.FridgeGenerData{Time: makeTimestamp(), Data: (rand.Float32() * 10) - 8}

		case <-stopInner:
			log.Warningln("DataGenerator - stopInner-case")
			return
		}

	}
}

func RunDataCollector(config *config.DevConfig, cBot <-chan models.FridgeGenerData,
	cTop <-chan models.FridgeGenerData, ReqChan chan models.Request) {
	duration := config.GetSendFreq()
	stopInner := make(chan struct{})
	log.Warningln("RunDataCollector: duration", duration)
	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)

	configChanged := make(chan struct{})
	config.AddSubIntoPool("DataCollector", configChanged)

	go DataCollector(ticker, cBot, cTop, ReqChan, stopInner)

	for {
		select {
		case <-configChanged:
			state := config.GetTurned()
			switch state {
			case true:
				close(stopInner)
				ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Millisecond)
				go DataCollector(ticker, cBot, cTop, ReqChan, stopInner)
				log.Warningln("go DataCollector has been started after signal")
			case false:
				close(stopInner)
				log.Warningln("turnedOn: off signal")
			}
		}
	}
}

func RunDataGenerator(config *config.DevConfig, cBot chan<- models.FridgeGenerData, cTop chan<- models.FridgeGenerData) {
	duration := config.GetCollectFreq()
	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)
	stopInner := make(chan struct{})

	configChanged := make(chan struct{})
	config.AddSubIntoPool("DataGenerator", configChanged)

	go DataGenerator(ticker, cBot, cTop, stopInner)

	for {

		select {
		case <-configChanged:
			state := config.GetTurned()
			switch state {
			case true:
				close(stopInner)
				log.Warningln("close stopInner before new ticker")
				ticker = time.NewTicker(time.Duration(config.GetCollectFreq()) * time.Millisecond)
				go DataGenerator(ticker, cBot, cTop, stopInner)
				log.Warningln("go DataGenerator has been started after signal")
			case false:
				close(stopInner)
				log.Warningln("turnedOn: off signal")
			}
		}
	}

}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func getDial(connType string, host string, port string) *net.Conn {
	// var reconnect *time.Ticker
	conn, err := net.Dial(connType, host+":"+port)
	checkError("getDial error", err)
	// for err != nil {
	// 	reconnect = time.NewTicker(time.Millisecond * 1)
	// 	for range reconnect.C {
	// 		conn, err = net.Dial(connType, host+":"+port)
	// 	}
	// }

	return &conn
}

func send(r models.Request, conn *net.Conn) {
	var resp models.Response
	r.Time = time.Now().UnixNano()

	err := json.NewEncoder(*conn).Encode(&r)
	checkError("send: JSON Enc", err)

	err = json.NewDecoder(*conn).Decode(&resp)
	checkError("send: JSON Dec", err)

	log.Warningln("Response: ", resp)
}

func constructReq(mTop map[int64]float32, mBot map[int64]float32) models.Request {
	var fridgeData models.FridgeData
	args := os.Args[1:]
	fridgeData.TempCam2 = mBot
	fridgeData.TempCam1 = mTop

	req := models.Request{
		Action: "update",
		Time:   time.Now().UnixNano(),
		Meta: models.Metadata{
			Type: args[0],
			Name: args[1],
			MAC:  args[2]},
		Data: fridgeData,
	}

	return req
}

func checkError(desc string, err error) error {
	if err != nil {
		log.Errorln(desc, err)
		return err
	}
	return nil
}
