package config

import (
	"testing"

	"github.com/device-smart-house/models"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	connTypeConf = "localhost"
	hostConf     = "3000"
	portConf     = "tcp"
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

// func TestAskConfig(t *testing.T) {
// 	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
// 	conn, _ := net.Dial(connTypeConf, hostConf+":"+portConf)
// 	req := models.Config{
// 		CollectFreq: 1000,
// 		SendFreq:    5000,
// 		TurnedOn:    true,
// 	}

// 	Convey("AskConfig should send request to the server", t, func() {

// 		go func() {
// 			ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
// 			for {
// 				c, _ := ln.Accept()
// 				a := models.Config{
// 					CollectFreq: 1000,
// 					SendFreq:    5000,
// 					TurnedOn:    true,
// 				}
// 				json.NewEncoder(c).Encode(a)
// 			}
// 		}()

// 		cfg := askConfig(conn)
// 		So(cfg.TurnedOn, ShouldEqual, req.TurnedOn)
// 		So(cfg.SendFreq, ShouldEqual, req.SendFreq)
// 		So(cfg.CollectFreq, ShouldEqual, req.CollectFreq)
// 	})
// }
