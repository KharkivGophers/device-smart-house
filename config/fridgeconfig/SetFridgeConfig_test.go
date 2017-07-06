package fridgeconfig

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)
func TestSetTurned(t *testing.T) {
	Convey("Should set valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetTurned(false)
		So(testConfig.GetTurned(), ShouldEqual, false)
	})
}

func TestSetCollectFreq(t *testing.T) {

	Convey("Should set valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetCollectFreq(1000)
		So(testConfig.GetCollectFreq(), ShouldEqual, 1000)
	})
}

func TestSetSendFreq(t *testing.T) {

	Convey("Should set valid value", t, func() {
		testConfig := NewFridgeConfig()
		testConfig.SetSendFreq(1000)
		So(testConfig.GetSendFreq(), ShouldEqual, 1000)
	})
}