package main

import (
	"sync"
)

type SyncMap struct {
	state map[string]interface{}
	mutex *sync.Mutex
}

func NewSyncMap() *SyncMap {
	s := &SyncMap{
		state: make(map[string]interface{}),
		mutex: &sync.Mutex{},
	}
	return s
}

func (s *SyncMap) Get(key string) interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.state[key]
}

func (s *SyncMap) Put(key string, value interface{}) {
	s.mutex.Lock()
	s.state[key] = value
	s.mutex.Unlock()
}

func (s *SyncMap) Delete(key string) {
	s.mutex.Lock()
	delete(s.state, key)
	s.mutex.Unlock()
}

func (s *SyncMap) Keys() chan string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	keys := make(chan string, len(s.state))
	for key := range s.state {
		keys <- key
	}
	close(keys)
	return keys
}

func (s *SyncMap) Values() chan interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	values := make(chan interface{}, len(s.state))
	for _, value := range s.state {
		values <- value
	}
	close(values)
	return values
}
