package core

import (
	"net"
	"fmt"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"errors"
)

const EmptyString = ""

var ErrSelectFailed = errors.New("select failed")

var CurrentServerGroup *ServerGroup

type ServerGroup struct {
	Servers    []interface{}
	Name       string
	SelectType string
	Selector   ISelector
}

type Server struct {
	Name     string
	Host     string
	Port     string
	Password string
	Method   string
}

var serverBucket map[string]interface{}

func ResetServers(groups []*ServerGroup) error {
	bucket := make(map[string]interface{}, len(groups))
	var ok bool
	var s *Server
	for i, v := range groups {
		if _, ok = bucket[v.Name]; ok {
			return errors.New("duplication of server name")
		}
		bucket[v.Name] = groups[i]

		for _, j := range v.Servers {
			s, ok = j.(*Server)
			if ok {
				bucket[s.Name] = j
			}
		}
	}
	serverBucket = bucket
	if len(groups) > 0 {
		CurrentServerGroup = groups[0]
	}
	return nil
}

func getServer(name string) *Server {
	if serverBucket == nil {
		return nil
	}
	v, ok := serverBucket[name]
	if !ok {
		return nil
	}
	switch v.(type) {
	case *Server:
		return v.(*Server)
	case *ServerGroup:
		group := v.(*ServerGroup)
		s, _ := group.Selector.Get(EmptyString)
		return s
	}
	return nil
}

func HandSelect(groupName, serverName string) error {
	if serverBucket == nil {
		return ErrSelectFailed
	}
	v, ok := serverBucket[groupName]
	if !ok {
		return ErrSelectFailed
	}
	switch v.(type) {
	case *Server:
		return ErrSelectFailed
	case *ServerGroup:
		group := v.(*ServerGroup)
		_, err := group.Selector.Get(serverName)
		if err != nil {
			return ErrSelectFailed
		}
		return nil
	}
	return ErrSelectFailed
}

func (s *Server) Conn(host string) (net.Conn, error) {
	rawAddr, err := ss.RawAddr(host)
	if err != nil {
		panic("Error getting raw address.")
	}
	cipher, err := ss.NewCipher(s.Method, s.Password)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return nil, err
	}
	return ss.DialWithRawAddr(rawAddr, net.JoinHostPort(s.Host, s.Port), cipher.Copy())
}
