package device

import (
	"errors"
	"net"
	"os"
	"testing"

	"reflect"

	"encoding/json"

	"sync"

	"time"

	"fmt"

	"github.com/device-smart-house/config"
	"github.com/device-smart-house/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDial(t *testing.T) {
	connTypeConf := "tcp"
	hostConf := "localhost"
	portConf := "3000"

	Convey("TCP connection should be estabilished", t, func() {
		ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
		defer ln.Close()
		conn := getDial(connTypeConf, hostConf, portConf)
		defer conn.Close()
		So(conn, ShouldNotBeNil)
	})
}

func TestCheckError(t *testing.T) {
	exErr := errors.New("Produce error")
	Convey("CheckError should return error's value", t, func() {
		err := checkError("Error message", exErr)
		So(err.Error(), ShouldEqual, exErr.Error())
	})
}

func TestConstructReq(t *testing.T) {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
	var exReq models.Request
	bot := make(map[int64]float32)
	top := make(map[int64]float32)

	bot[1] = 1.01
	bot[2] = 2.02
	bot[3] = 3.03

	top[1] = 10.01
	top[2] = 20.01
	top[3] = 30.01

	exReq = models.Request{
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

func TestSend(t *testing.T) {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
	var req models.Request
	var resp models.Response

	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	exReq := models.Request{
		Action: "update",
		Meta: models.Metadata{
			Type: os.Args[1],
			Name: os.Args[2],
			MAC:  os.Args[3]},
	}

	resp = models.Response{Status: 200, Descr: "Struct has been received"}
	Convey("Send should send JSON to the server", t, func() {
		go send(exReq, client)

		json.NewDecoder(server).Decode(&req)
		json.NewEncoder(server).Encode(resp)

		So(req.Action, ShouldEqual, exReq.Action)
		So(req.Meta.MAC, ShouldEqual, exReq.Meta.MAC)
		So(req.Meta.Name, ShouldEqual, exReq.Meta.Name)
		So(req.Meta.Type, ShouldEqual, exReq.Meta.Type)
	})
}

func TestMakeTimeStamp(t *testing.T) {
	Convey("MakeTimeStamp should return timestamp as int64", t, func() {
		time := makeTimestamp()
		So(reflect.TypeOf(time).String(), ShouldEqual, "int64")
		So(time, ShouldNotBeEmpty)
		So(time, ShouldNotEqual, 0)
	})
}

//how to change conn configs?
func TestDataTransfer(t *testing.T) {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
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
	var wg sync.WaitGroup
	cfg := config.GetConfig()
	ch := make(chan models.Request)

	Convey("DataTransfer should receive req from chan and transfer it to the server", t, func() {
		ln, _ := net.Listen(connTypeOut, hostOut+":"+portOut)
		defer ln.Close()

		wg.Add(1)
		go func() {
			conn, _ := ln.Accept()
			err := json.NewDecoder(conn).Decode(&req)
			checkError("decode", err)
			wg.Done()
		}()

		go DataTransfer(cfg, ch, &wg)

		ch <- exReq

		wg.Wait()

		b := reflect.DeepEqual(req.Data, exReq.Data)
		So(req.Action, ShouldEqual, exReq.Action)
		//Compare struct
		So(b, ShouldEqual, true)
		So(req.Meta.MAC, ShouldEqual, exReq.Meta.MAC)
		So(req.Meta.Name, ShouldEqual, exReq.Meta.Name)
		So(req.Meta.Type, ShouldEqual, exReq.Meta.Type)
	})
}

func TestDataGenerator(t *testing.T) {
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Millisecond)
	top := make(chan models.FridgeGenerData)
	bot := make(chan models.FridgeGenerData)
	stopInner := make(chan struct{})

	Convey("DataGenerator should produce structs with data", t, func() {
		wg.Add(1)
		var fromTop, fromBot models.FridgeGenerData
		var okTop, okBot bool

		go DataGenerator(ticker, bot, top, stopInner, &wg)
		fromTop, okTop = <-top
		fromBot, okBot = <-bot

		time.Sleep(time.Millisecond * 10)

		fmt.Println(fromBot)
		fmt.Println(fromTop)
		So(okTop, ShouldEqual, true)
		So(okBot, ShouldEqual, true)
		So(fromBot.Data, ShouldNotEqual, 0)
		So(fromTop.Data, ShouldNotEqual, 0)
		So(reflect.TypeOf(fromBot.Data).String(), ShouldEqual, "float32")
		So(reflect.TypeOf(fromTop.Data).String(), ShouldEqual, "float32")
	})
}
