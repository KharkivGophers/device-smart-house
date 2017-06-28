package models

import "sync"

type WasherData struct {
	Turnovers map[int64]int64 `json:"turnovers"`
	WaterTemp map[int64]float32
}

type GenerateWasherData struct {
	Time int64
	TurnoversData int64
	WaterTempData float32
}

type CollectWasherData struct {
	TurnoversStorage chan GenerateWasherData
	TemperatureStorage chan GenerateWasherData
	RequestStorage chan Request
	Wg      sync.WaitGroup
}