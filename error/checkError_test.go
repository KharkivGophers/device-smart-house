package error

import (
	"testing"
	"errors"
	"github.com/smartystreets/goconvey/convey"
)

func TestCheckError(t *testing.T) {
	exErr := errors.New("Produce error")
	convey.Convey("CheckError should return error's value", t, func() {
		err := CheckError("Error message", exErr)
		convey.So(err.Error(), convey.ShouldEqual, exErr.Error())
	})
}