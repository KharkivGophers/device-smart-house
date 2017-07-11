package fridgeconfig

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
)
func TestSetTurned(t *testing.T) {
	convey.Convey("Should set valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetTurned(false)
		convey.So(testConfig.GetTurned(), convey.ShouldEqual, false)
	})
}

func TestSetCollectFreq(t *testing.T) {

	convey.Convey("Should set valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetCollectFreq(1000)
		convey.So(testConfig.GetCollectFreq(), convey.ShouldEqual, 1000)
	})
}

func TestSetSendFreq(t *testing.T) {

	convey.Convey("Should set valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetSendFreq(1000)
		convey.So(testConfig.GetSendFreq(), convey.ShouldEqual, 1000)
	})
}