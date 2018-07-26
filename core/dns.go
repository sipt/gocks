package core

import (
	"strings"
	"regexp"
	"github.com/miekg/dns"
	"net"
	"sync"
)

var (
	// ipv6
	ipv6Regexp, _ = regexp.Compile("^((([0-9A-Fa-f]{1,4}:){7}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){1,7}:)|(([0-9A-Fa-f]{1,4}:){6}:[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){5}(:[0-9A-Fa-f]{1,4}){1,2})|(([0-9A-Fa-f]{1,4}:){4}(:[0-9A-Fa-f]{1,4}){1,3})|(([0-9A-Fa-f]{1,4}:){3}(:[0-9A-Fa-f]{1,4}){1,4})|(([0-9A-Fa-f]{1,4}:){2}(:[0-9A-Fa-f]{1,4}){1,5})|([0-9A-Fa-f]{1,4}:(:[0-9A-Fa-f]{1,4}){1,6})|(:(:[0-9A-Fa-f]{1,4}){1,7})|(([0-9A-Fa-f]{1,4}:){6}(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])(\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])){3})|(([0-9A-Fa-f]{1,4}:){5}:(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])(\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])){3})|(([0-9A-Fa-f]{1,4}:){4}(:[0-9A-Fa-f]{1,4}){0,1}:(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])(\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])){3})|(([0-9A-Fa-f]{1,4}:){3}(:[0-9A-Fa-f]{1,4}){0,2}:(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])(\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])){3})|(([0-9A-Fa-f]{1,4}:){2}(:[0-9A-Fa-f]{1,4}){0,3}:(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])(\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])){3})|([0-9A-Fa-f]{1,4}:(:[0-9A-Fa-f]{1,4}){0,4}:(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])(\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])){3})|(:(:[0-9A-Fa-f]{1,4}){0,5}:(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])(\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])){3}))$")
	// ipv4
	ipv4Regexp, _ = regexp.Compile(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`)
)

var dnsServers []string
var localDNS []*DNS

func ResetDNSServers(servers []string) {
	dnsServers = servers
}
func ResetLocalDNS(hosts []*DNS) {
	localDNS = hosts
}

type DNS struct {
	Host string
	Addr string
}

func LookupIP(req *Request) {
	bytes := []byte(req.host)
	if ipv4Regexp.Match(bytes) {
		req.addrType = AddrTypeIPv4
	} else if ipv6Regexp.Match(bytes) {
		req.addrType = AddrTypeIPv6
	} else {
		req.addrType = AddrTypeDomain
	}
	for _, v := range localDNS {
		if v.Host[0] == '*' && strings.HasSuffix(req.host, v.Host[1:]) {
			req.addr = v.Addr
			return
		} else if v.Host == req.host {
			req.addr = v.Addr
			return
		}
	}
	var err error
	req.ips, err = resolve(req.host)
	if err != nil {
		Logger.Error("DNS resolve failed: " + err.Error())
		req.ips = nil
	}
}

var lock sync.Mutex

func resolve(host string) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), dns.TypeA)
	m.RecursionDesired = true

	ips := make([]string, 0)
	for _, v := range dnsServers {
		lock.Lock()
		r, _, _ := c.Exchange(m, net.JoinHostPort(v, "53"))
		lock.Unlock()
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
		break
	}

	// Stuff must be in the answer section
	Logger.Debug("[DNS] resolve: ", host, " -> ", ips)
	//return ips, nil
	return []string{}, nil
}
