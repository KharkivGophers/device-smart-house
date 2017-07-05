package fridge

import (
	"testing"
	"time"
	"reflect"
	"github.com/KharkivGophers/device-smart-house/models"
	. "github.com/smartystreets/goconvey/convey"

)

func TestDataGenerator(t *testing.T) {
	ticker := time.NewTicker(time.Millisecond)
	top := make(chan models.FridgeGenerData)
	bot := make(chan models.FridgeGenerData)
	stopInner := make(chan struct{})

	Convey("DataGenerator should produce structs with data", t, func() {
		var fromTop, fromBot models.FridgeGenerData
		var okTop, okBot bool

		go DataGenerator(ticker, bot, top, stopInner)
		fromTop, okTop = <-top
		fromBot, okBot = <-bot

		time.Sleep(time.Millisecond * 10)

		So(okTop, ShouldEqual, true)
		So(okBot, ShouldEqual, true)
		So(fromBot.Data, ShouldNotEqual, 0)
		So(fromTop.Data, ShouldNotEqual, 0)
		So(reflect.TypeOf(fromBot.Data).String(), ShouldEqual, "float32")
		So(reflect.TypeOf(fromTop.Data).String(), ShouldEqual, "float32")
	})
}

func TestMakeTimeStamp(t *testing.T) {
	Convey("MakeTimeStamp should return timestamp as int64", t, func() {
		time := makeTimestamp()
		So(reflect.TypeOf(time).String(), ShouldEqual, "int64")
		So(time, ShouldNotBeEmpty)
		So(time, ShouldNotEqual, 0)
	})
}