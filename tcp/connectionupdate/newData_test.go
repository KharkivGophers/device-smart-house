package connectionupdate

import (
	"net"
	"testing"
	"time"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"encoding/json"
	"github.com/KharkivGophers/device-smart-house/models"
)

func TestGetDial(t *testing.T) {
	connTypeConf := "tcp"
	hostConf := "0.0.0.0"
	portConf := "3000"

	Convey("tcp tcp should be established", t, func() {
		ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
		conn := GetDial(connTypeConf, hostConf, portConf)
		time.Sleep(time.Millisecond * 100)
		defer ln.Close()
		defer conn.Close()
		So(conn, ShouldNotBeNil)
	})
}

func TestSend(t *testing.T) {
	os.Args = []string{"cmd", "fridge", "LG", "00-00-00-00-00-00"}
	var req models.Request
	var resp models.Response

	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	exReq := models.Request{
		Action: "update",
		Meta: models.Metadata{
			Type: os.Args[1],
			Name: os.Args[2],
			MAC:  os.Args[3]},
	}

	resp = models.Response{Descr: "Struct has been received"}
	Convey("Send should send JSON to the server", t, func() {
		go Send(exReq, client) // request counter is missing

		json.NewDecoder(server).Decode(&req)
		json.NewEncoder(server).Encode(resp)

		So(req.Action, ShouldEqual, exReq.Action)
		So(req.Meta.MAC, ShouldEqual, exReq.Meta.MAC)
		So(req.Meta.Name, ShouldEqual, exReq.Meta.Name)
		So(req.Meta.Type, ShouldEqual, exReq.Meta.Type)
	})
}

