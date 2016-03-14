package main

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestSyncMapPutGet(t *testing.T) {
	RegisterTestingT(t)

	syncMap := NewSyncMap()
	syncMap.Put("hello", "world")
	Expect(syncMap.Get("hello")).To(Equal("world"))
	syncMap.Put("hi", "world")
	Expect(syncMap.Get("hi")).To(Equal("world"))
	syncMap.Put("hello", "there")
	Expect(syncMap.Get("hello")).To(Equal("there"))
}

func TestSyncMapPutGetRace(t *testing.T) {
	RegisterTestingT(t)

	syncMap := NewSyncMap()

	go func() {
		time.Sleep(time.Millisecond)
		syncMap.Put("hello", "world")
	}()
	go syncMap.Put("hi", "world")
	go syncMap.Put("hello", "there")

	time.Sleep(time.Second)

	Expect(syncMap.Get("hi")).To(Equal("world"))
	Expect(syncMap.Get("hello")).To(Equal("world"))
}

func TestSyncMapDelete(t *testing.T) {
	RegisterTestingT(t)

	syncMap := NewSyncMap()
	syncMap.Put("hello", "world")
	syncMap.Put("hi", "there")

	keys := make([]string, 0, 1)
	for key := range syncMap.Keys() {
		keys = append(keys, key)
	}
	Expect(keys).To(ConsistOf("hi", "hello"))

	syncMap.Delete("hello")
	keys = make([]string, 0, 1)
	for key := range syncMap.Keys() {
		keys = append(keys, key)
	}
	Expect(keys).To(ConsistOf("hi"))

	syncMap.Delete("hi")
	keys = make([]string, 0, 1)
	for key := range syncMap.Keys() {
		keys = append(keys, key)
	}
	Expect(keys).To(BeEmpty())
}

func TestSyncMapKeys(t *testing.T) {
	RegisterTestingT(t)

	syncMap := NewSyncMap()
	syncMap.Put("hello", "world")

	keys := make([]string, 0, 1)
	for key := range syncMap.Keys() {
		keys = append(keys, key)
	}
	Expect(keys).To(ConsistOf("hello"))

	syncMap.Put("hi", "there")
	keys = make([]string, 0, 2)
	for key := range syncMap.Keys() {
		keys = append(keys, key)
	}
	Expect(keys).To(ConsistOf("hi", "hello"))
	syncMap.Put("hi", "again")
	keys = make([]string, 0, 2)
	for key := range syncMap.Keys() {
		keys = append(keys, key)
	}
	Expect(keys).To(ConsistOf("hi", "hello"))
}

func TestSyncMapValues(t *testing.T) {
	RegisterTestingT(t)

	syncMap := NewSyncMap()
	syncMap.Put("hello", "world")

	values := make([]interface{}, 0, 1)
	for value := range syncMap.Values() {
		values = append(values, value)
	}
	Expect(values).To(ConsistOf("world"))

	syncMap.Put("hi", "there")
	values = make([]interface{}, 0, 2)
	for value := range syncMap.Values() {
		values = append(values, value)
	}
	Expect(values).To(ConsistOf("there", "world"))
	syncMap.Put("hi", "again")
	values = make([]interface{}, 0, 2)
	for value := range syncMap.Values() {
		values = append(values, value)
	}
	Expect(values).To(ConsistOf("world", "again"))
}

func TestSyncMapLen(t *testing.T) {
	RegisterTestingT(t)

	syncMap := NewSyncMap()
	Expect(syncMap.Len()).To(Equal(0))
	syncMap.Put("hello", "world")
	Expect(syncMap.Len()).To(Equal(1))
	syncMap.Put("hi", "there")
	Expect(syncMap.Len()).To(Equal(2))
	syncMap.Put("hello", "there")
	Expect(syncMap.Len()).To(Equal(2))
	syncMap.Delete("hello")
	Expect(syncMap.Len()).To(Equal(1))
	syncMap.Delete("hi")
	Expect(syncMap.Len()).To(Equal(0))
}
