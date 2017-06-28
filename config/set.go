package config

func (d *DevConfig) SetTurned(b bool) {
	d.Mutex.Lock()
	d.turned = b
	defer d.Mutex.Unlock()
}

func (d *DevConfig) SetCollectFreq(b int64) {
	d.Mutex.Lock()
	d.collectFreq = b
	d.Mutex.Unlock()

}

func (d *DevConfig) SetSendFreq(b int64) {
	d.Mutex.Lock()
	d.sendFreq = b
	d.Mutex.Unlock()

}