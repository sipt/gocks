package core

import (
	"net"
	"github.com/sipt/gocks/util"
	"bytes"
	"fmt"
)

func NewConn(network string, conn net.Conn) *Conn {
	id, _ := util.IW.NextId()
	return &Conn{
		Conn:    conn,
		ID:      id,
		Network: network,
		Values:  make(map[string]interface{}),
	}
}

type Entity struct {
	Key   string
	Value interface{}
}

type Conn struct {
	net.Conn
	ID      int64
	Network string
	Values  map[string]interface{}
}

func (c *Conn) PutValue(key string, value interface{}) {
	c.Values[key] = value
}

func (c *Conn) Record() string {
	buf := bytes.NewBufferString(fmt.Sprint(c.ID, "   "))
	for _, v := range RecordFormat {
		buf.WriteString(fmt.Sprint(c.Values[v], "   "))
	}
	return buf.String()
}
