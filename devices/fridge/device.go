package device

import (
	"encoding/json"
	"math/rand"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vpakhuchyi/device-smart-house/models"
)

//FridgeGenerData struct for data transfer
type FridgeGenerData struct {
	Time int64
	Data float32
}

func send(r models.Request, conn net.Conn) {
	log.Println("DataTransfer: Send: intro")
	var resp models.Response
	r.Time = time.Now().UnixNano()

	err := json.NewEncoder(conn).Encode(&r)
	if err != nil {
		log.Errorln(err)
	}

	err = json.NewDecoder(conn).Decode(&resp)
	if err != nil {
		log.Errorln(err)
	}

	//info for debugging
	log.Println("DataTransfer: [done]")
	log.Warningln("Response: ", resp)
}

//DataTransfer func sends reuest as JSON to the centre
func DataTransfer(config *models.DevConfig, reqChan chan models.Request, conn net.Conn) {

	for {
		select {
		case r := <-reqChan:
			go send(r, conn)

		}
	}
}

func constructReq(mTop map[int64]float32, mBot map[int64]float32) models.Request {
	var fridgeData models.FridgeData
	fridgeData.TempCam2 = mBot
	fridgeData.TempCam1 = mTop

	req := models.Request{
		Action: "update",
		Time:   time.Now().UnixNano(),
		Meta: models.Metadata{
			Type: "fridge",
			Name: "hladik0e31",
			MAC:  "00-15-E9-2B-99-3C"},
		Data: fridgeData,
	}

	return req
}

//DataCollector func gathers data from DataGenerator
//and sends completed request's structures to the ReqChan channel
func DataCollector(ticker *time.Ticker, cBot <-chan FridgeGenerData, cTop <-chan FridgeGenerData, ReqChan chan models.Request, stopInner chan struct{}) {
	log.Warningln("collector intro")
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
func DataGenerator(ticker *time.Ticker, cBot chan<- FridgeGenerData, cTop chan<- FridgeGenerData, stopInner chan struct{}) {
	log.Warningln("gener intro")

	for {
		select {
		case <-ticker.C:
			cTop <- FridgeGenerData{Time: makeTimestamp(), Data: rand.Float32() * 10}
			cBot <- FridgeGenerData{Time: makeTimestamp(), Data: (rand.Float32() * 10) - 8}
		case <-stopInner:
			log.Warningln("DataGenerator - stopInner-case")
			return
		}

	}
}

func RunDataCollector(config *models.DevConfig, cBot <-chan FridgeGenerData, cTop <-chan FridgeGenerData, ReqChan chan models.Request, stop <-chan struct{}, start <-chan struct{}) {
	duration := config.GetSendFreq()
	stopInner := make(chan struct{})
	ticker := time.NewTicker(time.Duration(duration) * time.Second)

	go DataCollector(ticker, cBot, cTop, ReqChan, stopInner)

	for {
		if duration != config.GetSendFreq() {
			ticker.Stop()
			log.Warningln(config.GetSendFreq())
			ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Second)
			log.Warningln("new ticker")
			log.Warningln(ticker)
		}
		select {
		case <-stop:
			ticker.Stop()
			close(stopInner)
			log.Warningln("RunDataCollector - stop-case")
		case <-start:
			ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Second)
			log.Warningln("RunDataCollector - start-case")
			continue
		}
	}
}

func RunDataGenerator(config *models.DevConfig, cBot chan<- FridgeGenerData, cTop chan<- FridgeGenerData, stop chan struct{}, start chan struct{}) {
	duration := config.GetCollectFreq()
	ticker := time.NewTicker(time.Duration(duration) * time.Second)
	stopInner := make(chan struct{})

	go DataGenerator(ticker, cBot, cTop, stopInner)

	for {
		if duration != config.GetCollectFreq() {
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(config.GetCollectFreq()) * time.Second)
		}
		select {

		case <-stop:
			ticker.Stop()
			close(stopInner)
			log.Warningln("RunDataGenerator - stop-case")
		case <-start:
			log.Warningln("1 start ")
			ticker = time.NewTicker(time.Duration(config.GetCollectFreq()) * time.Second)
			log.Warningln("2 start ")
			log.Warningln("RunDataGenerator - start-case")
			continue
		}
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
