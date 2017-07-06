package fridgeconfig

func NewFridgeConfig() *DevFridgeConfig {
	conf := &DevFridgeConfig{}
	conf.subsPool = make(map[string]chan struct{})

	return conf
}

func (d *DevFridgeConfig) GetTurned() bool {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.turned
}

func (d *DevFridgeConfig) GetCollectFreq() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.collectFreq
}

func (d *DevFridgeConfig) GetSendFreq() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.sendFreq
}
