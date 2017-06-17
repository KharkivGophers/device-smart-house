package config

import (
	"encoding/json"
	"net"
	"os"
	"testing"

	"github.com/device-smart-house/models"

	"errors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetTurned(t *testing.T) {
	Convey("Should set valid value", t, func() {
		cfg := GetConfig()
		cfg.SetTurned(false)
		So(cfg.GetTurned(), ShouldEqual, false)
	})
}

func TestGetTurned(t *testing.T) {

	Convey("Should get valid value", t, func() {
		cfg := GetConfig()
		cfg.SetTurned(false)
		So(cfg.GetTurned(), ShouldEqual, false)
	})
}

func TestSetCollectFreq(t *testing.T) {

	Convey("Should set valid value", t, func() {
		cfg := GetConfig()
		cfg.SetCollectFreq(1000)
		So(cfg.GetCollectFreq(), ShouldEqual, 1000)
	})
}

func TestGetCollectFreq(t *testing.T) {

	Convey("Should get valid value", t, func() {
		cfg := GetConfig()
		cfg.SetCollectFreq(1000)
		So(cfg.GetCollectFreq(), ShouldEqual, 1000)
	})
}

func TestSetSendFreq(t *testing.T) {

	Convey("Should set valid value", t, func() {
		cfg := GetConfig()
		cfg.SetSendFreq(1000)
		So(cfg.GetSendFreq(), ShouldEqual, 1000)
	})
}

func TestGetSendFreq(t *testing.T) {

	Convey("Should get valid value", t, func() {
		cfg := GetConfig()
		cfg.SetSendFreq(1000)
		So(cfg.GetSendFreq(), ShouldEqual, 1000)
	})
}

func TestAddSubIntoPool(t *testing.T) {
	ch := make(chan struct{})
	key := "19-29"

	Convey("AddSubIntoPool should add chan into the pool", t, func() {
		cfg := GetConfig()
		cfg.AddSubIntoPool(key, ch)
		So(cfg.subsPool[key], ShouldEqual, ch)
	})
}

func TestRemoveSubFromPool(t *testing.T) {
	ch := make(chan struct{})
	key := "19-29"

	Convey("RemoveSubFromPool should remove chan from the pool", t, func() {
		cfg := GetConfig()
		cfg.AddSubIntoPool(key, ch)

		cfg.RemoveSubFromPool(key)
		So(cfg.subsPool[key], ShouldEqual, nil)
	})
}

func TestUpdateConfig(t *testing.T) {
	maskOsArgs()

	exCfg := models.Config{
		TurnedOn:    true,
		SendFreq:    100,
		CollectFreq: 50}

	Convey("UpdateConfig should update struct by new struct's values", t, func() {
		cfg := GetConfig()
		cfg.updateConfig(exCfg)
		So(cfg.GetTurned(), ShouldEqual, exCfg.TurnedOn)
		So(cfg.GetCollectFreq(), ShouldEqual, exCfg.CollectFreq)
		So(cfg.GetSendFreq(), ShouldEqual, exCfg.SendFreq)
	})
}

func TestCheckError(t *testing.T) {

	Convey("CheckError should return error's value", t, func() {
		exErr := errors.New("Produce error")
		err := checkError("Error message", exErr)
		So(err.Error(), ShouldEqual, exErr.Error())
	})
}

func TestListenConfig(t *testing.T) {
	maskOsArgs()

	cfg := models.Config{
		TurnedOn:    true,
		CollectFreq: 1000,
		SendFreq:    5000}

	connTypeConf := "tcp"
	hostConf := "localhost"
	portConf := "3000"

	Convey("ListenConfig should receive a configuration", t, func() {

		ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
		go func() {
			defer ln.Close()
			server, err := ln.Accept()
			if err != nil {
				t.Fail()
			}
			err = json.NewEncoder(server).Encode(cfg)
			if err != nil {
				t.Fail()
			}
		}()

		client, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			t.Fail()
		}

		devConfig := GetConfig()

		listenConfig(devConfig, client)

		So(devConfig.GetSendFreq(), ShouldEqual, 5000)
		So(devConfig.GetCollectFreq(), ShouldEqual, 1000)
		So(devConfig.GetTurned(), ShouldEqual, true)
	})
}

// func TestPublishConfig(t *testing.T) {
// 	maskOsArgs()

// 	connTypeConf := "tcp"
// 	hostConf := "localhost"
// 	portConf := "3001"

// 	firstSubChan := make(chan struct{})
// 	secondSubChan := make(chan struct{})

// 	cfg := models.Config{
// 		TurnedOn:    true,
// 		CollectFreq: 1000,
// 		SendFreq:    5000}

// 	Convey("PublishConfigfig should notificate all subs", t, func() {

// 		ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
// 		go func() {
// 			defer ln.Close()
// 			server, err := ln.Accept()
// 			if err != nil {
// 				t.Fail()
// 			}
// 			err = json.NewEncoder(server).Encode(cfg)
// 			if err != nil {
// 				t.Fail()
// 			}
// 		}()

// 		client, err := net.Dial("tcp", ln.Addr().String())
// 		if err != nil {
// 			t.Fail()
// 		}

// 		devConfig := GetConfig()

// 		go listenConfig(devConfig, client)

// 		_, a := <-firstSubChan
// 		_, b := <-secondSubChan

// 		So(a, ShouldEqual, true)
// 		So(b, ShouldEqual, true)

// 	})
// }

func TestInit(t *testing.T) {
	maskOsArgs()
	devCfg := models.Config{
		TurnedOn:    true,
		CollectFreq: 1000,
		SendFreq:    5000}

	connTypeConf := "tcp"
	hostConf := "localhost"
	portConf := "3000"

	Convey("Init should receive config", t, func() {
		ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
		go func() {
			defer ln.Close()
			server, err := ln.Accept()
			if err != nil {
				t.Fail()
			}
			err = json.NewEncoder(server).Encode(devCfg)
			if err != nil {
				t.Fail()
			}
		}()

		Init(connTypeConf, hostConf, portConf)
		cfg := GetConfig()
		So(cfg.GetSendFreq(), ShouldEqual, 5000)
		So(cfg.GetCollectFreq(), ShouldEqual, 1000)
		So(cfg.GetTurned(), ShouldEqual, true)
	})
}

func maskOsArgs() {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
}
