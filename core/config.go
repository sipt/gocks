package core

import (
	"os"
	"github.com/go-ini/ini"
	"strings"
)

func Config(path string) (port, socksPort string) {
	port, socksPort = "8890", "8891"
	cfg, err := ini.Load(path)
	if err != nil {
		Logger.Error("Fail to read file: ", err)
		os.Exit(1)
	}
	// General
	{
		general, err := cfg.GetSection("General")
		if err != nil {
			Logger.Error("Fail to read [General]: ", err)
			os.Exit(1)
		}
		if general.HasKey("dns-server") {
			key, err := general.GetKey("dns-server")
			if err != nil {
				Logger.Error("Fail to read [General]-[dns-server]: ", err)
				os.Exit(1)
			}
			dnsServers := key.Strings(",")
			ResetDNSServers(dnsServers)
		} else {
			// default dns-server
			ResetDNSServers([]string{"114.114.114.114", "223.5.5.5"})
		}
		if general.HasKey("port") {
			key, err := general.GetKey("port")
			if err != nil {
				Logger.Error("Fail to read [General]-[port]: ", err)
				os.Exit(1)
			}
			port = key.String()
		}
		if general.HasKey("socks-port") {
			key, err := general.GetKey("socks-port")
			if err != nil {
				Logger.Error("Fail to read [General]-[socks-port]: ", err)
				os.Exit(1)
			}
			socksPort = key.String()
		}
	}
	Logger.Info("load [General] success")
	// Proxy
	var servers map[string]interface{}
	{
		proxy, err := cfg.GetSection("Proxy")
		if err != nil {
			Logger.Error("Fail to read [Proxy]: ", err)
			os.Exit(1)
		}
		keys := proxy.Keys()
		servers = make(map[string]interface{})
		servers[PolicyDirect] = &Server{
			Name:PolicyDirect,
		}
		servers[PolicyReject] = &Server{
			Name:PolicyReject,
		}
		var (
			vs []string
			ok bool
		)
		for _, k := range keys {
			vs = k.Strings(",")
			if _, ok = servers[k.Name()]; ok {
				Logger.Error(" duplication of server name: ", k.Name())
				os.Exit(1)
			}
			servers[k.Name()] = &Server{
				Name:     k.Name(),
				Host:     vs[0],
				Port:     vs[1],
				Method:   vs[2],
				Password: vs[3],
			}
		}
	}
	Logger.Info("load [Proxy] success")
	// Proxy Group
	{
		proxy, err := cfg.GetSection("Proxy Group")
		if err != nil {
			Logger.Error("Fail to read [Proxy Group]: ", err)
			os.Exit(1)
		}
		keys := proxy.Keys()
		groups := make([]*ServerGroup, len(keys))
		var (
			vs    []string
			group *ServerGroup
		)
		for i, k := range keys {
			group = &ServerGroup{
				Name: k.Name(),
			}
			groups[i] = group
			servers[group.Name] = group
		}
		for i, k := range keys {
			group = groups[i]
			vs = k.Strings(",")
			group.SelectType = vs[0]
			group.Servers = make([]interface{}, len(vs)-1)
			if !CheckSelectorType(vs[0]) {
				Logger.Error("Not support group select type:", vs[0])
				os.Exit(1)
			}
			vs = vs[1:]
			var (
				s  interface{}
				ok bool
			)

			for i, v := range vs {
				s, ok = servers[v]
				if !ok {
					Logger.Error("not exist [Proxy Group]: ", v)
					os.Exit(1)
				}
				group.Servers[i] = s
			}
		}
		err = ResetServers(groups)
		if err != nil {
			Logger.Error(err)
			os.Exit(1)
		}
	}
	Logger.Info("load [Proxy Group] success")

	//Host
	{
		host, err := cfg.GetSection("Host")
		if err != nil {
			Logger.Error("Fail to read [Host]: ", err)
			os.Exit(1)
		}
		keys := host.Keys()

		var hosts = make([]*DNS, len(keys))
		for i, k := range keys {
			hosts[i] = &DNS{
				Host: k.Name(),
				Addr: k.String(),
			}
		}
		ResetLocalDNS(hosts)
	}
	Logger.Info("load [Host] success")

	//Rule
	{
		rule, err := cfg.GetSection("Rule")
		if err != nil {
			Logger.Error("Fail to read [Rule]: ", err.Error())
			os.Exit(1)
		}
		lines := strings.Split(rule.Key("rules").String(), ";")
		var (
			rules = make([]*Rule, len(lines))
			geoIP = make([]*Rule, 0, 8)
			final *Rule
			r     *Rule
			i     = 0
			vs    []string
		)
		for _, l := range lines {
			vs = strings.Split(l, ",")
			if len(vs) < 3 {
				Logger.Error("Fail to read [Rule]: ", l)
				os.Exit(1)
			}
			r = &Rule{
				Type:   vs[0],
				Value:  vs[1],
				Policy: vs[2],
			}
			if len(vs) >= 4 {
				r.Comment = vs[3]
			}

			switch  r.Type {
			case RuleGeoIP:
				geoIP = append(geoIP, r)
			case RuleFinal:
				if final != nil {
					Logger.Error("rule [FINAL] exist")
					os.Exit(1)
				}
				final = r
			default:
				rules[i] = r
				i ++
			}
		}
		rules = rules[:i]
		rules = append(rules, geoIP...)
		rules = append(rules, final)
		RuleReset(rules)
	}
	Logger.Info("load [Rule] success")
	return
}
