package util

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type IPInfo struct {
	Code int `json:"code"`
	Data IP  `json:"data`
}

type IP struct {
	IP        string `json:"ip"`
	Country   string `json:"country"`
	Area      string `json:"area"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Isp       string `json:"isp"`
	CountryID string `json:"country_id"`
	AreaID    string `json:"area_id"`
	RegionID  string `json:"region_id"`
	CityID    string `json:"city_id"`
	IspID     string `json:"isp_id"`
}

func tabaoAPI(ip string) (*IPInfo, error) {
	resp, err := http.Get(fmt.Sprintf("http://ip.taobao.com/service/getIpInfo.php?ip=%s", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func WatchIP(addr string) (*IP, error) {
	reply, err := tabaoAPI(addr)
	if err != nil {
		return nil, err
	}
	if reply != nil {
		return &(reply.Data), nil
	}
	return nil, nil
}
