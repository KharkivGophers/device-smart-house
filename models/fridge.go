package models

import "sync"

type FridgeData struct {
	TempCam1 map[int64]float32 `json:"tempCam1"`
	TempCam2 map[int64]float32 `json:"tempCam2"`
}

type FridgeGenerData struct {
	Time int64
	Data float32
}

type CollectFridgeData struct {
	CBot chan FridgeGenerData
	CTop chan FridgeGenerData
	ReqChan chan Request
	Wg      sync.WaitGroup // TODO to main
}