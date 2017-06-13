package device

import (
	"encoding/json"
	"testing"

	"github.com/device-smart-house/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mTop = make(map[int64]float32)
	mBot = make(map[int64]float32)
)

func TestGetDial(t *testing.T) {
	Convey("TCP connection should be estabilished", t, func() {
		conn := getDial(connTypeOut, hostOut, portOut)
		So(conn, ShouldNotBeNil)
	})
}

func TestJSONTrensfer(t *testing.T) {
	var response models.Response

	res := models.Response{Status: 200, Descr: "Data has been delivered successfully"}
	req := models.Request{Action: "update", Time: 1496741392463499334, Meta: models.Metadata{Type: "fridge", Name: "LG", MAC: "00-00-00-00-00-00"}}
	Convey("JSON response should be the same", t, func() {
		conn := getDial(connTypeOut, hostOut, portOut)
		json.NewEncoder(*conn).Encode(req)

		json.NewDecoder(*conn).Decode(&response)

		So(response.Descr, ShouldEqual, res.Descr)
	})
}
