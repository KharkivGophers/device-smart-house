package washerconfig


func NewWasherConfig() *DevWasherConfig {
	conf := &DevWasherConfig{}
	conf.subsPool = make(map[string]chan struct{})

	return conf
}

func (d *DevWasherConfig) GetTemperature() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.Temperature
}

func (d *DevWasherConfig) GetWashTime() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.WashTime
}

func (d *DevWasherConfig) GetWashTurnovers() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.WashTurnovers
}

func (d *DevWasherConfig) GetRinseTime() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.RinseTime
}


func (d *DevWasherConfig) GetRinseTurnovers() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.RinseTurnovers
}


func (d *DevWasherConfig) GetSpinTime() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.SpinTime
}

func (d *DevWasherConfig) GetSpinTurnovers() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.SpinTurnovers
}
