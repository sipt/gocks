package core

import (
	"net"
	"net/http"
	"github.com/sipt/gocks/core/plugin"
)

const (
	HttpProxyHttpDumpOpen  = 1
	HttpProxyHttpDumpOff   = 0
	HttpProxyHttpsDumpOpen = 2
	HttpProxyHttpsDumpOff  = 1
)

var HttpProxyWorkType httpProxyWorkType = 0

type httpProxyWorkType int

func StartHttpDump() {
	HttpProxyWorkType |= HttpProxyHttpDumpOpen
}
func StopHttpDump() {
	HttpProxyWorkType &= HttpProxyHttpDumpOff
}
func StartHttpsDump() {
	HttpProxyWorkType |= HttpProxyHttpsDumpOpen
}
func StopHttpsDump() {
	HttpProxyWorkType &= HttpProxyHttpsDumpOff
}

type Plugin func(lc, sc net.Conn) error

func HttpProxy(lc, sc net.Conn, req *http.Request) error {
	if HttpProxyWorkType&HttpProxyHttpDumpOpen == HttpProxyHttpDumpOpen {
		lc.Close()
		sc.Close()
	} else {
		plugin.DefaultPlugin(lc, sc)
		req.Write(sc)
	}
	return nil
}
func HttpsProxy(lc, sc net.Conn, req *http.Request) error {
	lc.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if HttpProxyWorkType&HttpProxyHttpsDumpOpen == HttpProxyHttpsDumpOpen {
		lc.Close()
		sc.Close()
	} else {
		plugin.DefaultPlugin(lc, sc)
	}
	return nil
}
