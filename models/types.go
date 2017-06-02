package models

type Request struct {
	Action string      `json:"action"`
	Time   int64       `json:"time"`
	Meta   Metadata    `json:"meta"`
	Data   interface{} `json:"data"`
}

type Response struct {
	Status int    `json:"status"`
	Descr  string `json:"descr"`
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

type FridgeGenerData struct {
	Time int64
	Data float32
}
