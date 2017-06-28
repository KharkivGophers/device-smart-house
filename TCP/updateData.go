package TCP

import (
	"os"
	"net"
	"time"
	"encoding/json"
	"github.com/KharkivGophers/device-smart-house/models"
	log "github.com/Sirupsen/logrus"
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
		checkError("getDial()", err)
		times++
		log.Warningln("Reconnect times: ", times)
	}
	return conn
}

func Send(r models.Request, conn net.Conn, requestsCounter *int) {
	var resp models.Response
	r.Time = time.Now().UnixNano()

	err := json.NewEncoder(conn).Encode(r)
	checkError("send(): JSON Encode: ", err)

	err = json.NewDecoder(conn).Decode(&resp)
	checkError("send(): JSON Decode: ", err)
	*requestsCounter++
	log.Infoln("Request number:", *requestsCounter)
	log.Infoln("send(): Response from center: ", resp)
}

func checkError(desc string, err error) error {
	if err != nil {
		log.Errorln(desc, err)
		return err
	}
	return nil
}