package main

import "net"

func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}

	conn.Write([]byte("hello"))
}
