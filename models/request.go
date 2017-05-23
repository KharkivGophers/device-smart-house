package models

type Request struct {
	Action string      `json:"action"`
	Time   int64       `json:"time"`
	Meta   Metadata    `json:"meta"`
	Data   interface{} `json:"data"`
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
	Turned      bool `json:"turned"`
	CollectFreq int  `json:"collectFreq"`
	SendFreq    int  `json:"sendFreq"`
}
