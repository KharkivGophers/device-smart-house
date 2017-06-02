package main

import (
	"sync"

	"github.com/vpakhuchyi/device-smart-house/config"
	"github.com/vpakhuchyi/device-smart-house/models"
)

var (
	cBot            chan models.FridgeGenerData
	cTop            chan models.FridgeGenerData
	reqChan         chan models.Request
	confMap         map[string]interface{}
	sendFreqChan    chan int64
	collectFreqChan chan int64
	turnedOnChan    chan bool
	stop            chan struct{}
	start           chan struct{}
	wg              sync.WaitGroup
	conf            *config.DevConfig

	//for config's listenner
	hostConf     = "localhost"
	portConf     = "3000"
	connTypeConf = "tcp"
)
