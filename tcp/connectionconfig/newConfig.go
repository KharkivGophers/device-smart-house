package connectionconfig

import (
	"net"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/error"
	"os"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
)

func InputConfig(conn net.Conn) models.Config {
	args := os.Args[1:]

	log.Warningln("Type:"+"["+args[0]+"];", "Name:"+"["+args[1]+"];", "MAC:"+"["+args[2]+"]")

	if len(args) < 3 {
		panic("Incorrect devices's information")
	}

	switch args[0] {
	case "washer":
		return AskWasherConfig(conn, args)
	default:
		return AskFridgeConfig(conn, args)
	}
}

func AskWasherConfig(conn net.Conn, args []string) models.Config {
	req := models.Request{
		Action: "config",
		Meta: models.Metadata{
			Type: args[0],
			Name: args[1],
			MAC:  args[2]},
	}

	var resp models.Config

	err := json.NewEncoder(conn).Encode(req)
	error.CheckError("askConfig(): Encode JSON", err)

	err = json.NewDecoder(conn).Decode(&resp)
	error.CheckError("askConfig(): Decode JSON", err)

	if err != nil && resp.IsEmpty() {
		panic("Connection has been closed by center")
	}

	return resp
}

func AskFridgeConfig(conn net.Conn, args []string) models.Config {
	req := models.Request{
		Action: "config",
		Meta: models.Metadata{
			Type: args[0],
			Name: args[1],
			MAC:  args[2]},
	}

	var resp models.Config

	err := json.NewEncoder(conn).Encode(req)
	error.CheckError("askConfig(): Encode JSON", err)

	err = json.NewDecoder(conn).Decode(&resp)
	error.CheckError("askConfig(): Decode JSON", err)

	if err != nil && resp.IsEmpty() {
		panic("Connection has been closed by center")
	}

	return resp
}