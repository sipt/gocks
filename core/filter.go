package core

import (
	"strings"
	"github.com/sipt/gocks/util"
)

//====================================
//  Filter
//====================================
const (
	RuleDomainSuffix  = "DOMAIN-SUFFIX"
	RuleDomain        = "DOMAIN"
	RuleDomainKeyword = "DOMAIN-KEYWORD"
	RuleGeoIP         = "GEOIP"
	RuleFinal         = "FINAL"
)

type Rule struct {
	Type    string
	Value   string
	Policy  string
	Comment string
}

var rules []*Rule

func RuleReset(rs []*Rule) {
	rules = rs
}
func RuleAppend(r *Rule) {
	rules = append(rules, r)
}

func connFilter(request *Request) string {
LOOP:
	for i, v := range rules {
		switch v.Type {
		case RuleDomainSuffix:
			if strings.HasSuffix(request.host, v.Value) {
				return v.Policy
			}
		case RuleDomain:
			if request.host == v.Value {
				return v.Policy
			}
		case RuleDomainKeyword:
			if strings.Index(request.host, v.Value) >= 0 {
				return v.Policy
			}
		case RuleGeoIP:
			ipInfo, err := util.WatchIP(request.addr)
			if err != nil {
				Logger.Error("ip get CountryID error, ip:", request.addr, ", err:", err)
			}
			if err != nil || ipInfo == nil || ipInfo.CountryID == "" {
				if rules[len(rules)-1].Type == RuleFinal {
					Logger.Info(v)
					return rules[len(rules)-1].Policy
				}
				return PolicyNone
			}
			rs := rules[i:]
			for _, v := range rs {
				switch v.Type {
				case RuleGeoIP:
					if ipInfo.CountryID == v.Value {
						return v.Policy
					}
				case RuleFinal:
					return v.Policy
				}
			}
			break LOOP
		}
	}
	return PolicyNone
}
