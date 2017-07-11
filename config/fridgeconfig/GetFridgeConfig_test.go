package fridgeconfig

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetTurned(t *testing.T) {

	convey.Convey("Should get valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetTurned(false)
		convey.So(testConfig.GetTurned(), convey.ShouldEqual, false)
	})
}

func TestGetCollectFreq(t *testing.T) {

	convey.Convey("Should get valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetCollectFreq(1000)
		convey.So(testConfig.GetCollectFreq(), convey.ShouldEqual, 1000)
	})
}


func TestGetSendFreq(t *testing.T) {
	testConfig := NewFridgeConfig()
	convey.Convey("Should get valid value", t, func() {
		testConfig.SetSendFreq(1000)
		convey.So(testConfig.GetSendFreq(), convey.ShouldEqual, 1000)
	})
}