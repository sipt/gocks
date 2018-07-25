package plugin

import (
	"net"
	"github.com/sipt/gocks/pool"
)

//default implement plugin
func DefaultPlugin(lc, sc net.Conn) {
	go send(lc, sc, false)
	go send(sc, lc, true)
}

func send(from, to net.Conn, toc bool) {
	buf := pool.Get()
	var n int
	var err error
	for {
		n, err = from.Read(buf)
		if err != nil {
			// error log
			from.Close()
			to.Close()
			break
		}
		n, err = to.Write(buf[:n])
		if err != nil {
			// error log
			from.Close()
			to.Close()
			break
		}
	}
	pool.Put(buf)
}
