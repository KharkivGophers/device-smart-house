package config


func GetConfig() *DevConfig {
	once.Do(func() {
		config = &DevConfig{}
		config.subsPool = make(map[string]chan struct{})
	})
	return config
}

func (d *DevConfig) GetTurned() bool {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.turned
}

func (d *DevConfig) GetCollectFreq() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.collectFreq
}

func (d *DevConfig) GetSendFreq() int64 {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	return d.sendFreq
}