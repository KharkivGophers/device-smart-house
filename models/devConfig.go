package models

import "sync"

type DevConfig struct {
	sync.Mutex
	turned      bool
	collectFreq int
	sendFreq    int
}

var config *DevConfig
var once sync.Once

func GetConfig() *DevConfig {
	once.Do(func() {
		config = &DevConfig{}
	})
	return config
}

func (br *DevConfig) SetTurned(b bool) {
	br.Mutex.Lock()
	br.turned = b
	br.Mutex.Unlock()
}

func (br *DevConfig) GetTurned() bool {
	br.Mutex.Lock()
	defer br.Mutex.Unlock()
	return br.turned
}

func (br *DevConfig) GetCollectFreq() int {
	br.Mutex.Lock()
	defer br.Mutex.Unlock()
	return br.collectFreq
}

func (br *DevConfig) GetSendFreq() int {
	br.Mutex.Lock()
	defer br.Mutex.Unlock()
	return br.sendFreq
}

func (br *DevConfig) SetCollectFreq(b int) {
	br.Mutex.Lock()
	br.collectFreq = b
	br.Mutex.Unlock()

}

func (br *DevConfig) SetSendFreq(b int) {
	br.Mutex.Lock()
	br.sendFreq = b
	br.Mutex.Unlock()

}
