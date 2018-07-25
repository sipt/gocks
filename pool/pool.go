package pool

import (
	"sync"
)

var bufPool *sync.Pool

func init() {
	bufPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
}

func Get() []byte {
	return bufPool.Get().([]byte)
}

func Put(buf []byte) {
	bufPool.Put(buf)
}
