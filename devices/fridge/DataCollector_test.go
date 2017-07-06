package fridge

import (
	"os"
	"time"
	"reflect"
	"github.com/KharkivGophers/device-smart-house/models"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDataCollector(t *testing.T) {
	maskOsArgs()
	var req models.FridgeRequest
	ticker := time.NewTicker(time.Millisecond)
	top := make(chan models.FridgeGenerData)
	bot := make(chan models.FridgeGenerData)
	reqChan := make(chan models.FridgeRequest)
	stopInner := make(chan struct{})

	botMap := make(map[int64]float32)
	topMap := make(map[int64]float32)

	topMap[0] = 1.01

	botMap[0] = 10.01

	exReq := models.FridgeRequest{
		Action: "update",
		Meta: models.Metadata{
			Type: os.Args[1],
			Name: os.Args[2],
			MAC:  os.Args[3]},
		Data: models.FridgeData{
			TempCam1: topMap,
			TempCam2: botMap},
	}

	Convey("DataGenerator should produce structs with data", t, func() {

		go DataCollector(ticker, bot, top, reqChan, stopInner)
		top <- models.FridgeGenerData{Data: 1.01}
		bot <- models.FridgeGenerData{Data: 10.01}

		time.Sleep(time.Millisecond * 10)

		req = <-reqChan

		// TODO Ask Viktor what the following comments mean
		//we have to refactor DataCollector: need to control WG
		// close(stopInner)
		//Compare struct's data
		b := reflect.DeepEqual(req.Data, exReq.Data)
		So(b, ShouldEqual, true)
	})
}

func TestConstructReq(t *testing.T) {
	os.Args = []string{"cmd", "fridgeconfig", "LG", "00-00-00-00-00-00"}
	var exReq models.FridgeRequest
	bot := make(map[int64]float32)
	top := make(map[int64]float32)

	bot[1] = 1.01
	bot[2] = 2.02
	bot[3] = 3.03

	top[1] = 10.01
	top[2] = 20.01
	top[3] = 30.01

	exReq = models.FridgeRequest{
		Action: "update",
		Meta: models.Metadata{
			Type: os.Args[1],
			Name: os.Args[2],
			MAC:  os.Args[3]},
		Data: models.FridgeData{TempCam1: top, TempCam2: bot},
	}
	Convey("ConstructReq should produce Request struct with received data", t, func() {
		req := constructReq(top, bot)
		b := reflect.DeepEqual(req.Data, exReq.Data)
		So(req.Action, ShouldEqual, exReq.Action)
		//Compare struct
		So(b, ShouldEqual, true)
		So(req.Meta.MAC, ShouldEqual, exReq.Meta.MAC)
		So(req.Meta.Name, ShouldEqual, exReq.Meta.Name)
		So(req.Meta.Type, ShouldEqual, exReq.Meta.Type)
	})
}

func maskOsArgs() {
	os.Args = []string{"cmd", "fridgeconfig", "LG", "00-00-00-00-00-00"}
}

