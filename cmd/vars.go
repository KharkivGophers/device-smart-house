package main

import (
	"sync"

	"github.com/KharkivGophers/device-smart-house/config"
	"github.com/KharkivGophers/device-smart-house/models"
)

var (
	cBot    chan models.FridgeGenerData
	cTop    chan models.FridgeGenerData
	reqChan chan models.Request
	wg      sync.WaitGroup
	conf    *config.DevConfig

	//for config's listenner
	hostConf     = "0.0.0.0"
	portConf     = "3000"
	connTypeConf = "tcp"
)
