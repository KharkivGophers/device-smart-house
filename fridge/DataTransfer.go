package fridge

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/connection"
	"github.com/KharkivGophers/device-smart-house/models"
)

//DataTransfer func sends request as JSON to the centre
func DataTransfer(config *config.DevConfig, reqChan chan models.Request) {

	// for data transfer
	transferConnParams := models.TransferConnParams{
		// HostOut: GetEnvCenter("CENTER_PORT_3030_TCP_ADDR"),
		HostOut: "0.0.0.0",
		PortOut: "3030",
		ConnTypeOut: "tcp",
	}

	conn := connection.GetDial(transferConnParams.ConnTypeOut, transferConnParams.HostOut, transferConnParams.PortOut)
	var requestsCounter int
	for {
		select {
		case r := <-reqChan:
			go connection.Send(r, conn, &requestsCounter)
		}
	}
}