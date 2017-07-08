package models

type WasherRequest struct {
	Action string     `json:"action"`
	Time   int64      `json:"time"`
	Meta   Metadata   `json:"meta"`
	Data   WasherData `json:"data"`
}

type WasherConfig struct {
	MAC				string	`json:"mac"`
	Temperature		float32 `json:"temperature"`
	WashTime		int64	`json:"washTime"`
	WashTurnovers 	int64	`json:"washTurnovers"`
	RinseTime		int64	`json:"rinseTime"`
	RinseTurnovers	int64	`json:"rinseTurnovers"`
	SpinTime		int64	`json:"spinTime"`
	SpinTurnovers	int64	`json:"spinTurnovers"`
}

func (washer WasherConfig) IsEmpty() bool {
	if washer.WashTime == 0 && washer.WashTurnovers == 0 {
		return true
	}
	return false
}

type WasherData struct {
	Turnovers map[int64]int64 	`json:"turnovers"`
	WaterTemp map[int64]float32	`json:"waterTemp"`
}

type GenerateWasherData struct {
	Time int64
	Turnovers int64
	WaterTemp float32
}

type CollectWasherData struct {
	TurnoversStorage chan GenerateWasherData
	TemperatureStorage chan GenerateWasherData
	RequestStorage chan WasherRequest
}