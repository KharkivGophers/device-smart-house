package connectionconfig

import (
	"net"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/error"
	"os"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
)

func AskConfig(conn net.Conn) models.FridgeConfig {
	args := os.Args[1:]
	log.Warningln("Type:"+"["+args[0]+"];", "Name:"+"["+args[1]+"];", "MAC:"+"["+args[2]+"]")
	if len(args) < 3 {
		panic("Incorrect devices's information")
	}

	var req models.FridgeRequest
	var resp models.FridgeConfig
	req = models.FridgeRequest{
		Action: "config",
		Meta: models.Metadata{
			Type: args[0],
			Name: args[1],
			MAC:  args[2]},
	}
	err := json.NewEncoder(conn).Encode(req)
	error.CheckError("askConfig(): Encode JSON", err)

	err = json.NewDecoder(conn).Decode(&resp)
	error.CheckError("askConfig(): Decode JSON", err)

	if err != nil && resp.IsEmpty() {
		panic("Connection has been closed by center")
	}

	return resp
}