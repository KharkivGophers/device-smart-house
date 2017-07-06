package models

type FridgeRequest struct {
	Action string     `json:"action"`
	Time   int64      `json:"time"`
	Meta   Metadata   `json:"meta"`
	Data   FridgeData `json:"data"`
}

type FridgeConfig struct {
	TurnedOn    bool   `json:"turnedOn"`
	CollectFreq int64  `json:"collectFreq"`
	SendFreq    int64  `json:"sendFreq"`
	MAC         string `json:"mac"`
}

func (c FridgeConfig) IsEmpty() bool {
	if c.CollectFreq == 0 && c.SendFreq == 0 && c.MAC == "" && c.TurnedOn == false {
		return true
	}
	return false
}

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
	ReqChan chan FridgeRequest
}