package models

type Request struct {
	Action string     `json:"action"`
	Time   int64      `json:"time"`
	Meta   Metadata   `json:"meta"`
	Data   FridgeData `json:"data"`
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

type Response struct {
	Descr string `json:"descr"`
}

type Control struct {
	Controller chan struct{}
}

func (c *Control) Close() {
	select {
	case <- c.Controller:
	default:
		close(c.Controller)
	}
}