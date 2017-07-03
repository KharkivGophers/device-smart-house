package fridge

import (
	"encoding/json"
	"time"
	"reflect"
	"testing"
	"os"
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
	"net"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
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

	var req models.Request
	exReq := models.Request{
		Action: "update",
		Meta: models.Metadata{
			Type: os.Args[1],
			Name: os.Args[2],
			MAC:  os.Args[3]},
		Data: models.FridgeData{
			TempCam1: top,
			TempCam2: bot},
	}

	testConfig := config.NewConfig()
	ch := make(chan models.Request)

	Convey("DataTransfer should receive req from chan and transfer it to the server", t, func() {
		ln, err := net.Listen(connTypeOut, hostOut+":"+portOut)
		if err != nil {
			t.Fail()
		}
		var Wg sync.WaitGroup
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
		go DataTransfer(testConfig, ch, &Wg)

		ch <- exReq
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