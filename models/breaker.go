package models

import "sync"

type Breaker struct {
	sync.Mutex
	Turned bool
}

var BreakerVar *Breaker
var once sync.Once

func GetBreaker() *Breaker {

	once.Do(func() {
		BreakerVar = &Breaker{}
	})
	return BreakerVar
}

func (br *Breaker) SetTurned(b bool) {
	br.Mutex.Lock()
	br.Turned = b
	br.Mutex.Unlock()

}

func (br *Breaker) GetTurned() bool {
	br.Mutex.Lock()
	defer br.Mutex.Unlock()
	return br.Turned
}
