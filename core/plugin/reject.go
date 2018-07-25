package plugin

import (
	"time"
	"net"
	"errors"
)

var EmptyConn, RejectConn net.Conn

func init() {
	EmptyConn = &emptyConn{}
	RejectConn = &rejectConn{EmptyConn}
}

var RejectError = errors.New("reject this conn")

type emptyConn struct{}

func (c *emptyConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (c *emptyConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (c *emptyConn) Close() error                       { return nil }
func (c *emptyConn) LocalAddr() net.Addr                { return nil }
func (c *emptyConn) RemoteAddr() net.Addr               { return nil }
func (c *emptyConn) SetDeadline(t time.Time) error      { return nil }
func (c *emptyConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *emptyConn) SetWriteDeadline(t time.Time) error { return nil }

type rejectConn struct {
	net.Conn
}

func (c *rejectConn) Read(b []byte) (n int, err error) {
	return 0, RejectError
}
