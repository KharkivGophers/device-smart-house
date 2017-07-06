package fridge

import (
	"encoding/json"
	"time"
	"reflect"
	"testing"
	"os"
	"github.com/KharkivGophers/device-smart-house/models"
	"net"
	. "github.com/smartystreets/goconvey/convey"
	log "github.com/Sirupsen/logrus"
	"github.com/KharkivGophers/device-smart-house/config/fridgeconfig"
)

//how to change conn configs?
func TestDataTransfer(t *testing.T) {
	maskOsArgs()
	connTypeOut := "tcp"
	hostOut := "localhost"
	portOut := "3030"

	bot := make(map[int64]float32)
	top := make(map[int64]float32)

	bot[1] = 1.01
	bot[2] = 2.02
	bot[3] = 3.03

	top[1] = 10.01
	top[2] = 20.01
	top[3] = 30.01

	var req models.FridgeRequest
	exReq := models.FridgeRequest{
		Action: "update",
		Meta: models.Metadata{
			Type: os.Args[1],
			Name: os.Args[2],
			MAC:  os.Args[3]},
		Data: models.FridgeData{
			TempCam1: top,
			TempCam2: bot},
	}

	testConfig := fridgeconfig.NewFridgeConfig()
	ch := make(chan models.FridgeRequest)

	Convey("DataTransfer should receive req from chan and transfer it to the server", t, func() {
		ln, err := net.Listen(connTypeOut, hostOut+":"+portOut)
		if err != nil {
			//t.Fail()
			panic("DataTransfer() Listen: error")
		}

		control := &models.Control{make(chan struct{})}
		go func() {
			defer ln.Close()
			server, err := ln.Accept()
			if err != nil {
				//t.Fail()
				panic("DataTransfer() Accept: invalid connection")
			}
			err = json.NewDecoder(server).Decode(&req)
			if err != nil {
				//t.Fail()
				panic("DataTransfer() Decode: invalid data to decode")
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
			}
		} ()
		go DataTransfer(testConfig, ch, control)

		ch <- exReq

		// TODO Ask Viktor what the following code means
		//need to refactor DataTransfer (can't wait for it)
		time.Sleep(time.Millisecond * 10)
		b := reflect.DeepEqual(req.Data, exReq.Data)
		So(req.Action, ShouldEqual, exReq.Action)
		//Compare struct
		So(b, ShouldEqual, true)
		So(req.Meta.MAC, ShouldEqual, exReq.Meta.MAC)
		So(req.Meta.Name, ShouldEqual, exReq.Meta.Name)
		So(req.Meta.Type, ShouldEqual, exReq.Meta.Type)
	})
}