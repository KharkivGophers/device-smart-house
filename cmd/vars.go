package main

import (
	"sync"

	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
	"github.com/KharkivGophers/device-smart-house/device"
)

var (
	cBot    chan models.FridgeGenerData
	cTop    chan models.FridgeGenerData
	reqChan chan models.Request
	wg      sync.WaitGroup
	conf    *config.DevConfig

	//for config's listener
	hostConf     = device.GetEnvCenter("CENTER_PORT_3000_TCP_ADDR")
	portConf     = "3000"
	connTypeConf = "tcp"
)
