package main

import (
	"bytes"
	"encoding/gob"
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

func (s *SyncMap) Len() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.state)
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

func (s *SyncMap) GobEncode() (data []byte, err error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if err := enc.Encode(s.state); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *SyncMap) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if err := dec.Decode(&s.state); err != nil {
		return err
	}
	return nil
}
