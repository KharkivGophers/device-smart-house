package models

import "time"

type Metadata struct {
	Type string `json:"type"`
	Name string `json:"name"`
	MAC  string `json:"mac"`
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

func (c *Control) Wait() {
	<- c.Controller
	<-time.NewTimer(6).C
}