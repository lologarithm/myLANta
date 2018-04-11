package main

import (
	"encoding/binary"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	exit := make(chan int, 10)
	log.Printf("Launching Server.")
	network := RunServer(exit)

	go func() {
		for {
			time.Sleep(time.Second)
			network.outgoing <- &Message{
				Raw: []byte{1, 2, 3},
			}
		}
	}()

	buf := make([]byte, 1)
	os.Stdin.Read(buf)
	exit <- 1
	exit <- 1
	log.Printf("goodbye")
}

type Client struct {
	Addr  *net.UDPAddr
	ID    int16
	Alive bool
}

type Message struct {
	Raw    []byte
	Target int16
	Kind   byte
	Length uint32
	Data   interface{}
}

type MsgKind byte

const (
	MsgKindUnknown MsgKind = iota
	MsgKindClients
	MsgKindFiles
)

func decode(m *Message) *Message {
	dcd := &Message{
		Raw:    m.Raw,
		Target: m.Target,
		Kind:   m.Raw[0],
		Length: binary.LittleEndian.Uint32(m.Raw[1:3]),
	}

	return dcd
}
