package connectionupdate

import (
	"net"
	"testing"
	"time"
	"github.com/smartystreets/goconvey/convey"
	"os"
	"encoding/json"
	"github.com/KharkivGophers/device-smart-house/models"
	log "github.com/Sirupsen/logrus"
)

func TestGetDial(t *testing.T) {
	connTypeConf := "tcp"
	hostConf := "0.0.0.0"
	portConf := "3000"

	convey.Convey("tcp tcp should be established", t, func() {
		ln, _ := net.Listen(connTypeConf, hostConf+":"+portConf)
		conn := GetDial(connTypeConf, hostConf, portConf)
		time.Sleep(time.Millisecond * 100)
		defer ln.Close()
		defer conn.Close()
		convey.So(conn, convey.ShouldNotBeNil)
	})
}

func TestSend(t *testing.T) {
	os.Args = []string{"cmd", "fridgeconfig", "LG", "00-00-00-00-00-00"}
	var req models.FridgeRequest
	var resp models.Response

	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	exReq := models.FridgeRequest{
		Action: "update",
		Meta: models.Metadata{
			Type: os.Args[1],
			Name: os.Args[2],
			MAC:  os.Args[3]},
	}

	resp = models.Response{Descr: "Struct has been received"}
	convey.Convey("Send should send JSON to the server", t, func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
			}
		} ()
		go Send(exReq, client) // request counter is missing

		json.NewDecoder(server).Decode(&req)
		json.NewEncoder(server).Encode(resp)

		convey.So(req.Action, convey.ShouldEqual, exReq.Action)
		convey.So(req.Meta.MAC, convey.ShouldEqual, exReq.Meta.MAC)
		convey.So(req.Meta.Name, convey.ShouldEqual, exReq.Meta.Name)
		convey.So(req.Meta.Type, convey.ShouldEqual, exReq.Meta.Type)
	})
}