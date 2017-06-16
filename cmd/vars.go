package main

import (
	"sync"

	"github.com/device-smart-house/config"
	"github.com/device-smart-house/models"
)

var (
	cBot    chan models.FridgeGenerData
	cTop    chan models.FridgeGenerData
	reqChan chan models.Request
	wg      sync.WaitGroup
	conf    *config.DevConfig

	//for config's listenner
	hostConf     = "192.168.104.76"
	portConf     = "3000"
	connTypeConf = "tcp"
)
