package device

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vpakhuchyi/device-smart-house/models"
)

//DataTransfer func sends reuest as JSON to the centre
func DataTransfer(br *models.Breaker, reqChan chan models.Request) {
	var Resp models.Response
	for {
		if br.GetTurned() == true {
			select {
			case r := <-reqChan:
				go func() {
					r.Time = time.Now().UnixNano()
					b := new(bytes.Buffer)

					err := json.NewEncoder(b).Encode(&r)
					if err != nil {
						log.Panic(err)
					}

					//Preparing a request
					req, err := http.NewRequest("POST", "http://localhost:8080", b)
					if err != nil {
						log.Panic(err)
					}
					req.Header.Set("Content-Type", "application/json")

					//Sending the request to the centre
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						log.Panic(err)
					}
					defer resp.Body.Close()

					//Decoding response from centre
					err = json.NewDecoder(resp.Body).Decode(&Resp)
					if err != nil {
						log.Panic(err)
					}

					//info for debugging
					log.Println("DataTransfer: [done]")
					log.Warningln("Response: ", Resp)
				}()

			}
		}
	}
}

//DataCollector func gathers data from DataGenerator
//and sends completed request's structures to the ReqChan channel
func DataCollector(br *models.Breaker, cBot <-chan float32, cTop <-chan float32, ReqChan chan models.Request) {

	var mTop = make(map[int64]float32)
	var mBot = make(map[int64]float32)
	var fridgeData models.FridgeData

	ticker := time.NewTicker(time.Second * 5)

	for {
		if br.GetTurned() == true {

			go func() {
				for v := range cTop {
					mTop[time.Now().UnixNano()] = v
				}
			}()

			go func() {

				for z := range cBot {
					mBot[time.Now().UnixNano()] = z
				}
			}()

			select {
			case <-ticker.C:
				time.Sleep(time.Millisecond * 1)
				fridgeData.TempCam2 = mBot
				fridgeData.TempCam1 = mTop

				ReqChan <- models.Request{
					Action: "update",
					Time:   time.Now().UnixNano(),
					Meta: models.Metadata{
						Type: "fridge",
						Name: "hladik0e31",
						MAC:  "00-15-E9-2B-99-3C"},
					Data: fridgeData,
				}

				log.Println("DataCollector: [done]")

				//for debugg
				log.Println("TempCam1: ", fridgeData.TempCam1)
				log.Println("TempCam2: ", fridgeData.TempCam2)

				//Cleaning temp maps
				mTop = make(map[int64]float32)
				mBot = make(map[int64]float32)
			}

		}
	}
}

func DataGenerator(br *models.Breaker, cBot chan<- float32, cTop chan<- float32) {
	ticker := time.NewTicker(time.Second * 1)

	for {
		if br.GetTurned() == true {
			select {
			case <-ticker.C:
				cTop <- rand.Float32() * 10
				cBot <- (rand.Float32() * 10) - 8
			}
		}
	}

}
