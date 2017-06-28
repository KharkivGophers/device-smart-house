package fridge

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/TCPConnection"
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

	conn := TCPConnection.GetDial(transferConnParams.ConnTypeOut, transferConnParams.HostOut, transferConnParams.PortOut)
	var requestsCounter int
	for {
		select {
		case r := <-reqChan:
			go TCPConnection.Send(r, conn, &requestsCounter)
		}
	}
}

