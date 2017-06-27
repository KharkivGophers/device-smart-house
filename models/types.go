package models

import (
	"sync"
)

type Request struct {
	Action string     `json:"action"`
	Time   int64      `json:"time"`
	Meta   Metadata   `json:"meta"`
	Data   FridgeData `json:"data"`
}

type Response struct {
	Descr string `json:"descr"`
}

type FridgeData struct {
	TempCam1 map[int64]float32 `json:"tempCam1"`
	TempCam2 map[int64]float32 `json:"tempCam2"`
}

type Metadata struct {
	Type string `json:"type"`
	Name string `json:"name"`
	MAC  string `json:"mac"`
}

type Config struct {
	TurnedOn    bool   `json:"turnedOn"`
	CollectFreq int64  `json:"collectFreq"`
	SendFreq    int64  `json:"sendFreq"`
	MAC         string `json:"mac"`
}

func (c Config) IsEmpty() bool {
	if c.CollectFreq == 0 && c.SendFreq == 0 && c.MAC == "" && c.TurnedOn == false {
		return true
	}
	return false
}

type FridgeGenerData struct {
	Time int64
	Data float32
}

type ConfigConnParams struct {
	ConnTypeConf string
	HostConf string
	PortConf string

}

type TransferConnParams struct {
	HostOut string
	PortOut string
	ConnTypeOut string
}

type CollectData struct {
	CBot chan FridgeGenerData
	CTop chan FridgeGenerData
	ReqChan chan Request
	Wg      sync.WaitGroup
}