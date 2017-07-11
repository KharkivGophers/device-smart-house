package fridgeconfig

import (
	"encoding/json"
	"net"
	"os"
	"testing"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/smartystreets/goconvey/convey"
	log "github.com/Sirupsen/logrus"
)

func TestAddSubIntoPool(t *testing.T) {
	ch := make(chan struct{})
	key := "19-29"

	convey.Convey("AddSubIntoPool should add chan into the pool", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.AddSubIntoPool(key, ch)
		convey.So(testConfig.subsPool[key], convey.ShouldEqual, ch)
	})
}

func TestRemoveSubFromPool(t *testing.T) {
	ch := make(chan struct{})
	key := "19-29"

	convey.Convey("RemoveSubFromPool should remove chan from the pool", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.AddSubIntoPool(key, ch)

		testConfig.RemoveSubFromPool(key)
		convey.So(testConfig.subsPool[key], convey.ShouldEqual, nil)
	})
}

func TestUpdateConfig(t *testing.T) {
	maskOsArgs()

	exCfg := models.FridgeConfig{
		TurnedOn:    true,
		SendFreq:    100,
		CollectFreq: 50}

	convey.Convey("UpdateConfig should update struct by new struct's values", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.updateConfig(exCfg)
		convey.So(testConfig.GetTurned(), convey.ShouldEqual, exCfg.TurnedOn)
		convey.So(testConfig.GetCollectFreq(), convey.ShouldEqual, exCfg.CollectFreq)
		convey.So(testConfig.GetSendFreq(), convey.ShouldEqual, exCfg.SendFreq)
	})
}


func TestListenConfig(t *testing.T) {
	maskOsArgs()

	cfg := models.FridgeConfig{
		TurnedOn:    true,
		CollectFreq: 1000,
		SendFreq:    5000}

	connTypeConf := "tcp"
	hostConf := "localhost"
	portConf := "3000"

	convey.Convey("ListenConfig should receive a configuration", t, func() {

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
		testConfig := NewFridgeConfig()

		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
			}
		}()
		listenConfig(testConfig, client)

		convey.So(testConfig.GetSendFreq(), convey.ShouldEqual, 5000)
		convey.So(testConfig.GetCollectFreq(), convey.ShouldEqual, 1000)
		convey.So(testConfig.GetTurned(), convey.ShouldEqual, true)
	})
}

func TestInit(t *testing.T) {
	maskOsArgs()
	devCfg := models.FridgeConfig{
		TurnedOn:    true,
		CollectFreq: 1000,
		SendFreq:    5000}

	connTypeConf := "tcp"
	hostConf := "localhost"
	portConf := "3000"

	convey.Convey("Init should receive config", t, func() {
		control := &models.Control{make(chan struct{})}
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
		testConfig := NewFridgeConfig()

		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
			}} ()
		testConfig.RequestFridgeConfig(connTypeConf, hostConf, portConf, control, maskOsArgs())

		convey.So(testConfig.GetSendFreq(), convey.ShouldEqual, 5000)
		convey.So(testConfig.GetCollectFreq(), convey.ShouldEqual, 1000)
		convey.So(testConfig.GetTurned(), convey.ShouldEqual, true)
	})
}

func maskOsArgs() []string {
	os.Args = []string{"cmd", "fridgeconfig", "LG", "00-00-00-00-00-00"}
	return os.Args
}
