package main

import (
	"net"
	"net/http"
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"github.com/sipt/gocks/core/plugin"
	"time"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

func main() {
	c, err := net.Dial("udp", "")
	shadowsocks.NewConn()
}
func test(){
	conn, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		c, err := conn.Accept()
		if err != nil {
			panic(err)
		}
		go func(){
			in := bufio.NewReader(c)
			req, err := http.ReadRequest(in)
			if err != nil {
				fmt.Println(err)
				return
			}
			ips, _ := resolve(req.URL.Hostname())
			fmt.Printf("connect to:[%s] by [%s]\n", req.URL.Hostname(), ips[0])
			co, err := net.Dial("tcp", net.JoinHostPort(ips[0],"443"))
			if err != nil {
				fmt.Println(err)
				return
			}
			c.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
			plugin.DefaultPlugin(c, co)
			t := time.NewTicker(10*time.Second)
			<- t.C
			c.Close()
			co.Close()
		}()
	}
}
var dnsServers = []string{"223.5.5.5", "8.8.8.8", "114.114.114.114"}
func resolve(host string) ([]string, error) {
	c := new(dns.Client)
	c.Net = "tcp"
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), dns.TypeA)
	m.RecursionDesired = true

	ips := make([]string, 0)
	for _, v := range dnsServers {
		r, _, _ := c.Exchange(m, net.JoinHostPort(v, "53"))
		if r == nil || r.Rcode != dns.RcodeSuccess {
			continue
		}
		var a *dns.A
		var ok bool
		for _, v := range r.Answer {
			a, ok = v.(*dns.A)
			if ok {
				ips = append(ips, a.A.String())
			}
		}
		if len(ips) > 0 {
			fmt.Println("DNS[",v,"] resolve: ", host, " -> ", ips)
			break
		}
	}

	// Stuff must be in the answer section
	return ips, nil
	//return []string{}, nil
}
