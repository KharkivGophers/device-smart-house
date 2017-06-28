package config

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTurned(t *testing.T) {

	Convey("Should get valid value", t, func() {
		cfg := GetConfig()
		cfg.SetTurned(false)
		So(cfg.GetTurned(), ShouldEqual, false)
	})
}

func TestGetCollectFreq(t *testing.T) {

	Convey("Should get valid value", t, func() {
		cfg := GetConfig()
		cfg.SetCollectFreq(1000)
		So(cfg.GetCollectFreq(), ShouldEqual, 1000)
	})
}


func TestGetSendFreq(t *testing.T) {

	Convey("Should get valid value", t, func() {
		cfg := GetConfig()
		cfg.SetSendFreq(1000)
		So(cfg.GetSendFreq(), ShouldEqual, 1000)
	})
}