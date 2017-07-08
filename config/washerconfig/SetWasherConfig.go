package washerconfig

func (d *DevWasherConfig) SetTemperature(b float32) {
	d.Mutex.Lock()
	d.Temperature = b
	d.Mutex.Unlock()
}

func (d *DevWasherConfig) SetWashTime(b int64) {
	d.Mutex.Lock()
	d.WashTime = b
	d.Mutex.Unlock()
}

func (d *DevWasherConfig) SetWashTurnovers(b int64) {
	d.Mutex.Lock()
	d.WashTurnovers = b
	d.Mutex.Unlock()
}

func (d *DevWasherConfig) SetRinseTime(b int64) {
	d.Mutex.Lock()
	d.RinseTime = b
	d.Mutex.Unlock()
}

func (d *DevWasherConfig) SetRinseTurnovers(b int64) {
	d.Mutex.Lock()
	d.RinseTurnovers = b
	d.Mutex.Unlock()
}

func (d *DevWasherConfig) SetSpinTime(b int64) {
	d.Mutex.Lock()
	d.SpinTime = b
	d.Mutex.Unlock()
}

func (d *DevWasherConfig) SetSpinTurnovers(b int64) {
	d.Mutex.Lock()
	d.SpinTurnovers = b
	d.Mutex.Unlock()
}