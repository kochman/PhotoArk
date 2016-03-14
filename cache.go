package main

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Cache struct {
	sizeLimit float64 // in megabytes
	cacheDir  string
	fillFunc  func(string) []byte
	cacheHits *SyncMap
}

func NewCache(cacheDir string, sizeLimit float64, fill func(string) []byte) *Cache {
	c := &Cache{
		sizeLimit: sizeLimit,
		cacheDir:  cacheDir,
		cacheHits: NewSyncMap(),
		fillFunc:  fill,
	}

	// create cache folder
	err := os.MkdirAll(c.cacheDir, os.FileMode(0764))
	if err != nil {
		log.Print(err)
	}

	// restore cacheHits from file
	cacheHitsPath := filepath.Join(c.cacheDir, "cacheHits.gob")

	// register time.Time with gob
	gob.Register(time.Now())

	if _, err = os.Stat(cacheHitsPath); !os.IsNotExist(err) {
		val, err := ioutil.ReadFile(cacheHitsPath)
		if err != nil {
			log.Print(err)
		}
		buf := bytes.NewBuffer(val)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&c.cacheHits)
		if err != nil {
			log.Print(err)
		}
	}

	// every ten seconds, persist updated cacheHits
	go func() {
		for {
			time.Sleep(10 * time.Second)
			cacheHitsPath := filepath.Join(c.cacheDir, "cacheHits.gob")

			buf := &bytes.Buffer{}
			enc := gob.NewEncoder(buf)
			err := enc.Encode(c.cacheHits)
			if err != nil {
				log.Print(err)
			}

			err = ioutil.WriteFile(cacheHitsPath, buf.Bytes(), os.FileMode(0664))
			if err != nil {
				log.Print(err)
			}
		}

	}()

	return c
}

func (c *Cache) Get(key string) []byte {
	path := filepath.Join(c.cacheDir, key)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// file doesn't exist; create it
		val := c.fillFunc(key)

		// set hit counter
		c.cacheHits.Put(key, time.Now())

		// write val to disk
		parentDir, _ := filepath.Split(path)
		err = os.MkdirAll(parentDir, os.FileMode(0764))
		if err != nil {
			log.Print(err)
		}
		err := ioutil.WriteFile(path, val, os.FileMode(0664))
		if err != nil {
			log.Print(err)
		}
		return val
	} else if err != nil {
		log.Print(err)
		return nil
	} else {
		// file exists in cache; return contents

		// set hit counter
		c.cacheHits.Put(key, time.Now())

		// read file contents and return
		val, err := ioutil.ReadFile(path)
		if err != nil {
			log.Print(err)
		}
		return val
	}
}

func (c *Cache) Prune() {
	for c.Size() > c.sizeLimit {
		// prune cache on disk

		// find least-recently-used item
		if c.cacheHits.Len() > 0 {
			var oldestKey string
			var oldestTime time.Time
			found := false
			for key := range c.cacheHits.Keys() {
				time := c.cacheHits.Get(key).(time.Time)
				if !found || time.Before(oldestTime) {
					found = true
					oldestKey = key
					oldestTime = time
				}
			}

			// remove oldest
			path := filepath.Join(c.cacheDir, oldestKey)
			err := os.Remove(path)
			if os.IsNotExist(err) {
				log.Print(err)
				c.cacheHits.Delete(oldestKey)
			} else if err != nil {
				log.Print(err)
			} else {
				c.cacheHits.Delete(oldestKey)
			}
		}
	}
}

func (c *Cache) PeriodicallyPrune(pause time.Duration) {
	for {
		time.Sleep(pause)
		c.Prune()
	}
}

func (c *Cache) Size() float64 {
	var dirSizeBytes int64
	pruneCacheWalk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if !info.IsDir() {
			dirSizeBytes += info.Size()
		}
		return nil
	}

	filepath.Walk(c.cacheDir, pruneCacheWalk)
	dirSizeMegabytes := float64(dirSizeBytes) / 1024 / 1024

	return dirSizeMegabytes
}
