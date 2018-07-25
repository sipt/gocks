package storage

import "sync"

type Storage struct {
	s map[string]interface{}
	sync.RWMutex
}

var s *Storage

func init() {
	s = &Storage{
		s: make(map[string]interface{}),
	}
}

func Get(key string) interface{} {
	return s.s[key]
}

func Put(key string, value interface{}) {
	s.s[key] = value
}
