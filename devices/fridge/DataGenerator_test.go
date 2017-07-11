package fridge

import (
	"testing"
	"time"
	"reflect"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/smartystreets/goconvey/convey"

)

func TestDataGenerator(t *testing.T) {
	ticker := time.NewTicker(time.Millisecond)
	top := make(chan models.FridgeGenerData)
	bot := make(chan models.FridgeGenerData)
	stopInner := make(chan struct{})

	convey.Convey("DataGenerator should produce structs with data", t, func() {
		var fromTop, fromBot models.FridgeGenerData
		var okTop, okBot bool

		go DataGenerator(ticker, bot, top, stopInner)
		fromTop, okTop = <-top
		fromBot, okBot = <-bot

		time.Sleep(time.Millisecond * 10)

		convey.So(okTop, convey.ShouldEqual, true)
		convey.So(okBot, convey.ShouldEqual, true)
		convey.So(fromBot.Data, convey.ShouldNotEqual, 0)
		convey.So(fromTop.Data, convey.ShouldNotEqual, 0)
		convey.So(reflect.TypeOf(fromBot.Data).String(), convey.ShouldEqual, "float32")
		convey.So(reflect.TypeOf(fromTop.Data).String(), convey.ShouldEqual, "float32")
	})
}

func TestMakeTimeStamp(t *testing.T) {
	convey.Convey("MakeTimeStamp should return timestamp as int64", t, func() {
		time := makeTimestamp()
		convey.So(reflect.TypeOf(time).String(), convey.ShouldEqual, "int64")
		convey.So(time, convey.ShouldNotBeEmpty)
		convey.So(time, convey.ShouldNotEqual, 0)
	})
}