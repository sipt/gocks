package core

import (
	"net"
	"bufio"
	"net/http"
	"github.com/sipt/gocks/core/plugin"
)

type Request struct {
	req      *http.Request
	host     string
	port     string
	addr     string
	ips      []string
	addrType int
}

//====================================
//  Connection
//====================================
const (
	PolicyReject = "REJECT"
	PolicyDirect = "DIRECT"
	PolicyNone   = "NONE"
)

func NewConnByReq(lc *Conn, req *http.Request) (net.Conn, error) {
	request := &Request{
		req:  req,
		host: req.URL.Hostname(),
		port: req.URL.Port(),
	}

	//dns
	LookupIP(request)
	if len(request.ips) == 0 {
		request.addr = request.host
	} else {
		request.addr = request.ips[0]
	}
	lc.PutValue(NameRemote, request.addr)

	policy := connFilter(request)
	lc.PutValue(NamePolicy, policy)
	if policy == PolicyReject { //reject
		return plugin.RejectConn, nil
	}
	switch req.URL.Scheme {
	case "http":
		if request.port == "" {
			request.port = "80"
		}
	case "https":
		//https
		if request.port == "" {
			request.port = "443"
		}
	case "socks":
	}
	var (
		conn net.Conn
		err  error
	)
	request.addr = net.JoinHostPort(request.addr, request.port)
	s := getServer(policy)
	lc.PutValue(NamePolicy, policy)
	lc.PutValue(NameProxy, s.Name)
	conn, err = s.Conn(request.addr)
	return conn, err
}

const (
	ProxyModelDirect = 1
	ProxyModelProxy  = 2
	ProxyModelAutoX  = 3
)

var ProxyModel = ProxyModelAutoX

func TCPProxy(lc *Conn) error {
	in := bufio.NewReader(lc)
	req, err := http.ReadRequest(in)
	if err != nil {
		return err
	}

	//
	if req.URL.Scheme == "http" {
	} else if req.URL.Scheme == "socks5" {
	} else if req.Method == http.MethodConnect {
		req.URL.Scheme = "https"
	}

	sc, err := NewConnByReq(lc, req)
	if err != nil {
		return err
	}
	lc.PutValue(NameHost, req.URL.Host)
	lc.PutValue(NameURL, req.URL.String())
	lc.PutValue(NameScheme, req.URL.Scheme)
	switch req.URL.Scheme {
	case "http":
		HttpProxy(lc, sc, req)
	case "https":
		//https
		HttpsProxy(lc, sc, req)
	case "socks":
		plugin.DefaultPlugin(lc, sc)
	}
	Logger.Info(lc.Record())
	return nil
}
