package main

import (
	"sync"

	"github.com/KharkivGophers/device-smart-house/device"
	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
)

var (
	cBot    chan models.FridgeGenerData
	cTop    chan models.FridgeGenerData
	reqChan chan models.Request
	wg      sync.WaitGroup
	conf    *config.DevConfig

	//for config's listener
	hostConf     = device.GetEnvCenter("CENTER_TCP_ADDR")
	portConf     = "3000"
	connTypeConf = "tcp"
)