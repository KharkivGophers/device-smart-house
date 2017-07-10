package main

import (
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/tcp/connectionupdate"
	log "github.com/Sirupsen/logrus"
)

func main() {
	configConnParams := models.ConfigConnParams{
		ConnTypeConf: "tcp",
		HostConf:     connectionupdate.GetEnvCenter("CENTER_PORT_3000_TCP_ADDR"),
		PortConf:     "3000",
	}

	newDevice := config.CreateDevice()
	newDeviceType := newDevice[0]
	control := &models.Control{make(chan struct{})}

	switch newDeviceType {
	case "washer":
		startWasher(configConnParams.ConnTypeConf, configConnParams.HostConf, configConnParams.PortConf, control, newDevice)
	case "fridge":
		startFridge(configConnParams.ConnTypeConf, configConnParams.HostConf, configConnParams.PortConf, control, newDevice)
	}

	control.Wait()
	log.Info("Device has been terminated due to the center's issue")
}