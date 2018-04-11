package main

import (
	"log"
	"net"
	"os"
)

func main() {
	exit := make(chan int, 10)
	log.Printf("Launching Server.")
	go RunServer(exit)
	buf := make([]byte, 1)
	os.Stdin.Read(buf)
	exit <- 1
	exit <- 1
	log.Printf("goodbye")
}

type Me struct {
	Peers []*Client
}

type Client struct {
	Addr  *net.UDPAddr
	ID    int
	Alive bool
}
