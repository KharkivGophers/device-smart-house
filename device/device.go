package device

import (
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/device-smart-house/config"
	"github.com/device-smart-house/models"
)

var i int

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
func DataCollector(ticker *time.Ticker, cBot <-chan models.FridgeGenerData, cTop <-chan models.FridgeGenerData,
	ReqChan chan models.Request, stopInner chan struct{}, wg *sync.WaitGroup) {

	var mTop = make(map[int64]float32)
	var mBot = make(map[int64]float32)

	for {
		select {
		case <-stopInner:

			log.Println("DataCollector(): wg.Done()")
			wg.Done()
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
func DataGenerator(ticker *time.Ticker, cBot chan<- models.FridgeGenerData, cTop chan<- models.FridgeGenerData,
	stopInner chan struct{}, wg *sync.WaitGroup) {

	for {
		select {
		case <-ticker.C:
			cTop <- models.FridgeGenerData{Time: makeTimestamp(), Data: rand.Float32() * 10}
			cBot <- models.FridgeGenerData{Time: makeTimestamp(), Data: (rand.Float32() * 10) - 8}

		case <-stopInner:

			log.Println("DataGenerator(): wg.Done()")
			wg.Done()
			return
		}

	}
}

//RunDataCollector setups DataCollector
func RunDataCollector(config *config.DevConfig, cBot <-chan models.FridgeGenerData,
	cTop <-chan models.FridgeGenerData, ReqChan chan models.Request, wg *sync.WaitGroup) {
	duration := config.GetSendFreq()
	stopInner := make(chan struct{})
	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)

	configChanged := make(chan struct{})
	config.AddSubIntoPool("DataCollector", configChanged)

	wg.Add(1)
	if config.GetTurned() {
		go DataCollector(ticker, cBot, cTop, ReqChan, stopInner, wg)
	}

	for {
		select {
		case <-configChanged:
			state := config.GetTurned()
			switch state {
			case true:
				select {
				case <-stopInner:
					wg.Add(1)
					stopInner = make(chan struct{})
					ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Millisecond)
					go DataCollector(ticker, cBot, cTop, ReqChan, stopInner, wg)
					log.Println("DataCollector() has been started")
				default:
					close(stopInner)
					stopInner = make(chan struct{})
					wg.Add(1)
					ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Millisecond)
					go DataCollector(ticker, cBot, cTop, ReqChan, stopInner, wg)
					log.Println("DataCollector() has been started")
				}
			case false:
				select {
				case <-stopInner:
					ticker = time.NewTicker(time.Duration(config.GetSendFreq()) * time.Millisecond)
				default:
					close(stopInner)
					log.Println("DataCollector() hase been killed")
				}
			}
		}
	}

}

//RunDataGenerator setups DataGenerator
func RunDataGenerator(config *config.DevConfig, cBot chan<- models.FridgeGenerData,
	cTop chan<- models.FridgeGenerData, wg *sync.WaitGroup) {
	duration := config.GetCollectFreq()
	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)
	stopInner := make(chan struct{})

	configChanged := make(chan struct{})
	config.AddSubIntoPool("DataGenerator", configChanged)

	wg.Add(1)
	if config.GetTurned() {
		go DataGenerator(ticker, cBot, cTop, stopInner, wg)
	}

	for {
		select {
		case <-configChanged:
			state := config.GetTurned()
			switch state {
			case true:
				select {
				case <-stopInner:
					wg.Add(1)
					stopInner = make(chan struct{})
					ticker = time.NewTicker(time.Duration(config.GetCollectFreq()) * time.Millisecond)
					go DataGenerator(ticker, cBot, cTop, stopInner, wg)
					log.Println("DataGenerator() has been started")
				default:
					close(stopInner)
					stopInner = make(chan struct{})
					wg.Add(1)
					ticker = time.NewTicker(time.Duration(config.GetCollectFreq()) * time.Millisecond)
					go DataGenerator(ticker, cBot, cTop, stopInner, wg)
					log.Println("DataGenerator() has been started")
				}
			case false:
				select {
				case <-stopInner:
					ticker = time.NewTicker(time.Duration(config.GetCollectFreq()) * time.Millisecond)
				default:
					close(stopInner)
					log.Println("DataGenerator() has been killed")
				}
			}
		}
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func getDial(connType string, host string, port string) net.Conn {
	var times int
	conn, err := net.Dial(connType, host+":"+port)

	for err != nil {
		if times >= 5 {
			panic("Can't connect to the server: send")
		}
		time.Sleep(time.Second)
		conn, err = net.Dial(connType, host+":"+port)
		checkError("getDial()", err)
		times++
		log.Warningln("Recennect times: ", times)
	}
	return conn
}

func send(r models.Request, conn net.Conn) {
	var resp models.Response
	r.Time = time.Now().UnixNano()

	err := json.NewEncoder(conn).Encode(r)
	checkError("send(): JSON Encode: ", err)

	err = json.NewDecoder(conn).Decode(&resp)
	checkError("send(): JSON Decode: ", err)
	i++
	log.Infoln("Request number:", i)
	log.Infoln("send(): Response from center: ", resp)
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
