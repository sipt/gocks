package main

import (
	"net"
	"fmt"
	"github.com/sipt/gocks/core"
	_ "github.com/sipt/gocks/core/selector"
)

func main() {
	startProxy()
}

func startProxy() {
	port, _ := core.Config("~/Documents/sipt.ini")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	fmt.Println("start listen:", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go func() {
			//defer func() {
			//	if r := recover(); r != nil {
			//		core.Logger.Error("[UNKNOWN] ", r)
			//	}
			//}()
			err = core.TCPProxy(core.NewConn("tcp", conn))
			if err != nil {
				core.Logger.Error("[UNKNOWN] ", err)
			}
		}()
	}
}
