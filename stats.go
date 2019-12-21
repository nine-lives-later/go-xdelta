package xdelta

import (
	"fmt"
	lib "github.com/konsorten/go-xdelta/xdelta-lib"
	"sync"
	"time"
)

type Stats struct {
	states map[lib.XdeltaState]time.Duration
	lock   sync.Mutex
}

func newStats() *Stats {
	return &Stats{
		states: make(map[lib.XdeltaState]time.Duration),
	}
}

func (s *Stats) DumpToStdout() {
	s.lock.Lock()

	fmt.Println("STATS:")

	for k, v := range s.states {
		fmt.Printf("  State %10v lastet %v\n", k, v)
	}

	s.lock.Unlock()
}

func (s *Stats) addStateTime(state lib.XdeltaState, duration time.Duration) {
	s.lock.Lock()

	prev, _ := s.states[state]

	s.states[state] = prev + duration

	s.lock.Unlock()
}
