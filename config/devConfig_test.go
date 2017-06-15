package config

import (
	"encoding/json"
	"net"
	"os"
	"testing"

	"github.com/device-smart-house/models"

	"errors"

	"fmt"

	"time"

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
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}

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

func TestInit(t *testing.T) {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
	devCfg := models.Config{
		TurnedOn:    true,
		CollectFreq: 1000,
		SendFreq:    5000}

	connTypeConf := "tcp"
	hostConf := "localhost"
	portConf := "3000"

	Convey("Init should receive config", t, func() {
		ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
		defer ln.Close()
		go func() {
			conn, _ := ln.Accept()
			fmt.Println(os.Args)
			json.NewEncoder(conn).Encode(devCfg)
		}()

		Init(connTypeConf, hostConf, portConf)
		cfg := GetConfig()
		So(cfg.GetSendFreq(), ShouldEqual, 5000)
		So(cfg.GetCollectFreq(), ShouldEqual, 1000)
		So(cfg.GetTurned(), ShouldEqual, true)
	})
}

func TestListenConfig(t *testing.T) {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
	firstSubChan := make(chan struct{})
	secondSubChan := make(chan struct{})
	server, client := net.Pipe()

	defer client.Close()
	defer server.Close()

	cfg := models.Config{
		TurnedOn:    true,
		CollectFreq: 1000,
		SendFreq:    5000}

	Convey("PublishConfigfig should notificate all subs", t, func() {

		go func() {
			json.NewEncoder(server).Encode(cfg)
		}()

		devConfig := GetConfig()

		devConfig.AddSubIntoPool("firstSub", firstSubChan)
		devConfig.AddSubIntoPool("secondSub", secondSubChan)

		go func() {
			json.NewEncoder(client).Encode(cfg)
		}()
		go listenConfig(devConfig, client)
		_, a := <-firstSubChan
		_, b := <-secondSubChan

		So(a, ShouldEqual, true)
		So(b, ShouldEqual, true)
		So(devConfig.GetSendFreq(), ShouldEqual, 5000)
		So(devConfig.GetCollectFreq(), ShouldEqual, 1000)
		So(devConfig.GetTurned(), ShouldEqual, true)

	})
}

func TestPublishConfig(t *testing.T) {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
	firstSubChan := make(chan struct{})
	secondSubChan := make(chan struct{})
	server, client := net.Pipe()

	cfg := models.Config{
		TurnedOn:    true,
		CollectFreq: 1000,
		SendFreq:    5000}

	defer client.Close()
	defer server.Close()
	Convey("PublishConfigfig should notificate all subs", t, func() {
		go func() {
			json.NewEncoder(server).Encode(cfg)

		}()
		defer server.Close()
		devConfig := GetConfig()

		devConfig.AddSubIntoPool("firstSub", firstSubChan)
		devConfig.AddSubIntoPool("secondSub", secondSubChan)

		go listenConfig(devConfig, client)
		go func() {
			json.NewEncoder(client).Encode(cfg)

		}()

		_, a := <-firstSubChan
		_, b := <-secondSubChan
		time.Sleep(time.Millisecond * 2)
		client.Close()

		So(a, ShouldEqual, true)
		So(b, ShouldEqual, true)
	})
}
