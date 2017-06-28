package config

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)
func TestSetTurned(t *testing.T) {
	Convey("Should set valid value", t, func() {
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

func TestSetSendFreq(t *testing.T) {

	Convey("Should set valid value", t, func() {
		cfg := GetConfig()
		cfg.SetSendFreq(1000)
		So(cfg.GetSendFreq(), ShouldEqual, 1000)
	})
}