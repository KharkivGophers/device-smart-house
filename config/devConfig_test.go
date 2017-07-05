package config

import (
	"encoding/json"
	"net"
	"os"
	"testing"
	"github.com/KharkivGophers/device-smart-house/models"
	. "github.com/smartystreets/goconvey/convey"
	log "github.com/Sirupsen/logrus"
)

func TestAddSubIntoPool(t *testing.T) {
	ch := make(chan struct{})
	key := "19-29"

	Convey("AddSubIntoPool should add chan into the pool", t, func() {
		testConfig := NewConfig()
		testConfig.AddSubIntoPool(key, ch)
		So(testConfig.subsPool[key], ShouldEqual, ch)
	})
}

func TestRemoveSubFromPool(t *testing.T) {
	ch := make(chan struct{})
	key := "19-29"

	Convey("RemoveSubFromPool should remove chan from the pool", t, func() {
		testConfig := NewConfig()
		testConfig.AddSubIntoPool(key, ch)

		testConfig.RemoveSubFromPool(key)
		So(testConfig.subsPool[key], ShouldEqual, nil)
	})
}

func TestUpdateConfig(t *testing.T) {
	maskOsArgs()

	exCfg := models.Config{
		TurnedOn:    true,
		SendFreq:    100,
		CollectFreq: 50}

	Convey("UpdateConfig should update struct by new struct's values", t, func() {
		testConfig := NewConfig()
		testConfig.updateConfig(exCfg)
		So(testConfig.GetTurned(), ShouldEqual, exCfg.TurnedOn)
		So(testConfig.GetCollectFreq(), ShouldEqual, exCfg.CollectFreq)
		So(testConfig.GetSendFreq(), ShouldEqual, exCfg.SendFreq)
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
				//t.Fail()
				panic("ListenConfig() Accept: No Connection")
			}
			err = json.NewEncoder(server).Encode(cfg)
			if err != nil {
				//t.Fail()
				panic("ListenConfig() Encode: invalid data to encode!")
			}
		}()

		client, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			//t.Fail()
			panic("ListenConfig() Dial: invalid address!")
		}
		testConfig := NewConfig()

		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
			}
		}()
		listenConfig(testConfig, client)

		So(testConfig.GetSendFreq(), ShouldEqual, 5000)
		So(testConfig.GetCollectFreq(), ShouldEqual, 1000)
		So(testConfig.GetTurned(), ShouldEqual, true)
	})
}

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
		control := &models.Control{make(chan struct{})}
		ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
		go func() {
			defer ln.Close()
			server, err := ln.Accept()
			if err != nil {
				//t.Fail()
				panic("Init() Accept: invalid connection!")
			}
			err = json.NewEncoder(server).Encode(devCfg)
			if err != nil {
				//t.Fail()
				panic("Init() Encode: invalid data to encode!")
			}
		}()
		testConfig := NewConfig()

		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
			}} ()
		testConfig.Init(connTypeConf, hostConf, portConf, control)

		So(testConfig.GetSendFreq(), ShouldEqual, 5000)
		So(testConfig.GetCollectFreq(), ShouldEqual, 1000)
		So(testConfig.GetTurned(), ShouldEqual, true)
	})
}

func maskOsArgs() {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
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