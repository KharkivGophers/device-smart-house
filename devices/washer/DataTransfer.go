package washer

import (
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/error"
	"github.com/KharkivGophers/device-smart-house/config/washerconfig"
	log "github.com/Sirupsen/logrus"
	"os"
	"net"
	"time"
	"encoding/json"

)

// Connection
func GetEnvCenter(key string) string {
	host := os.Getenv(key)
	if len(host) == 0 {
		return "127.0.0.1"
	}
	return host
}

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

	error.CheckError("send(): JSON Decode: ", err)
	if err != nil {
		panic("No response found")
	}

	log.Infoln("Request:")
	log.Infoln("send(): Response from center: ", resp)
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

