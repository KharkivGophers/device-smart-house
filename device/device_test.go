package device

import (
	"errors"
	"net"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	connTypeConf = "tcp"
	hostConf     = "localhost"
	portConf     = "3000"

	mTop = make(map[int64]float32)
	mBot = make(map[int64]float32)
)

func init() {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
}
func TestGetDial(t *testing.T) {
	ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
	defer ln.Close()
	Convey("TCP connection should be estabilished", t, func() {
		conn := getDial(connTypeConf, hostConf, portConf)
		So(conn, ShouldNotBeNil)
	})
}

func TestCheckError(t *testing.T) {

	Convey("CheckError should return error's value", t, func() {
		exErr := errors.New("Produce error")
		err := checkError("Error message", exErr)
		So(err.Error(), ShouldEqual, exErr.Error())
	})
}

// func TestJSONTrensfer(t *testing.T) {
// 	var response models.Response

// 	res := models.Response{Status: 200, Descr: "Data has been delivered successfully"}
// 	req := models.Request{Action: "update", Time: 1496741392463499334, Meta: models.Metadata{Type: "fridge", Name: "LG", MAC: "00-00-00-00-00-00"}}
// 	Convey("JSON response should be the same", t, func() {
// 		conn := getDial(connTypeOut, hostOut, portOut)
// 		json.NewEncoder(*conn).Encode(req)

// 		json.NewDecoder(*conn).Decode(&response)

// 		So(response.Descr, ShouldEqual, res.Descr)
// 	})
// }
