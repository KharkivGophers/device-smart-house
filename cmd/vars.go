package main

import (
	"sync"

	"github.com/vpakhuchyi/device-smart-house/config"
	"github.com/vpakhuchyi/device-smart-house/models"
)

var (
	cBot    chan models.FridgeGenerData
	cTop    chan models.FridgeGenerData
	reqChan chan models.Request
	stop    chan struct{}
	start   chan struct{}
	wg      sync.WaitGroup
	conf    *config.DevConfig

	//for config's listenner
	hostConf     = "192.168.104.60"
	portConf     = "3000"
	connTypeConf = "tcp"
)
