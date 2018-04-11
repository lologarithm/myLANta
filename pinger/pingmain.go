package main

import (
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "239.0.0.0:9999")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second * 2)
		conn.Write([]byte("hello"))
	}
}
