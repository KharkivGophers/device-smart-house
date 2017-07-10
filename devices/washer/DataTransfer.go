package washer

import (
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/error"
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
	log "github.com/Sirupsen/logrus"
	"net"
	"time"
	"encoding/json"

)

// Connection
func GetDial(connType string, host string, port string) net.Conn {
	var times int
	conn, err := net.Dial(connType, host+":"+port)

	for err != nil {
		if times >= 5 {
			panic("Can't connect to the server: send")
		}
		time.Sleep(time.Second)
		conn, err = net.Dial(connType, host+":"+port)
		error.CheckError("getDial()", err)
		times++
		log.Warningln("Reconnect times: ", times)
	}
	return conn
}

func Send(r models.WasherRequest, conn net.Conn) {
	var resp models.Response
	r.Time = time.Now().UnixNano()

	err := json.NewEncoder(conn).Encode(r)

	if err != nil {
		panic("Nothing to encode")
	}
	error.CheckError("send(): JSON Encode: ", err)

	err = json.NewDecoder(conn).Decode(&resp)
	if err != nil {
		panic("No response found")
	}
	error.CheckError("send(): JSON Decode: ", err)

	log.Infoln("Data was sent; Response from center: ", resp)
}

//DataTransfer func sends request as JSON to the centre
func DataTransfer(config *washerconfig.DevWasherConfig, requestStorage chan models.WasherRequest, c *models.Control) {

	// for data transfer
	transferConnParams := models.TransferConnParams{
		// HostOut: GetEnvCenter("CENTER_PORT_3030_TCP_ADDR"),
		HostOut: "0.0.0.0",
		PortOut: "3030",
		ConnTypeOut: "tcp",
	}

	defer func() {
		if a := recover(); a != nil {
			log.Error(a)
			c.Close()
		}
	} ()
	conn := GetDial(transferConnParams.ConnTypeOut, transferConnParams.HostOut, transferConnParams.PortOut)

	for {
		select {
		case r := <-requestStorage:
			go func() {
				defer func() {
					if a := recover(); a != nil {
						log.Error(a)
						c.Close()
					}
				} ()
				Send(r, conn)
			}()
		case <- c.Controller:
			log.Error("Data Transfer Failed")
			return
		}
	}
}

