package fridge

import (
	"encoding/json"
	"time"
	"reflect"
	"testing"
	"os"
	"github.com/KharkivGophers/device-smart-house/models"
	"net"
	"github.com/smartystreets/goconvey/convey"
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

	convey.Convey("DataTransfer should receive req from chan and transfer it to the server", t, func() {
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
				t.Fail()
			}
			err = json.NewDecoder(server).Decode(&req)
			if err != nil {
				t.Fail()
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
			}
		} ()
		go DataTransfer(testConfig, ch, control)

		ch <- exReq

		time.Sleep(time.Millisecond * 10)
		b := reflect.DeepEqual(req.Data, exReq.Data)
		convey.So(req.Action, convey.ShouldEqual, exReq.Action)
		//Compare struct
		convey.So(b, convey.ShouldEqual, true)
		convey.So(req.Meta.MAC, convey.ShouldEqual, exReq.Meta.MAC)
		convey.So(req.Meta.Name, convey.ShouldEqual, exReq.Meta.Name)
		convey.So(req.Meta.Type, convey.ShouldEqual, exReq.Meta.Type)
	})
}