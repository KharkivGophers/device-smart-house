package fridgeconfig

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTurned(t *testing.T) {

	Convey("Should get valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetTurned(false)
		So(testConfig.GetTurned(), ShouldEqual, false)
	})
}

func TestGetCollectFreq(t *testing.T) {

	Convey("Should get valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetCollectFreq(1000)
		So(testConfig.GetCollectFreq(), ShouldEqual, 1000)
	})
}


func TestGetSendFreq(t *testing.T) {
	testConfig := NewFridgeConfig()
	Convey("Should get valid value", t, func() {
		testConfig.SetSendFreq(1000)
		So(testConfig.GetSendFreq(), ShouldEqual, 1000)
	})
}