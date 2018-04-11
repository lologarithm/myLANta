package main

import (
	"fmt"
	"log"
	"net"
)

const (
	bport string = "9999"
)

type Network struct {
	conn        *net.UDPConn
	connections []*Client
	connLookup  map[string]int
	outgoing    chan *Message
}

func RunServer(exit chan int) {
	network := &Network{
		connections: []*Client{},
		connLookup:  map[string]int{},
		outgoing:    make(chan *Message, 100),
	}
	go runBroadcastListener(network, exit)
}

func runBroadcastListener(s *Network, exit chan int) {
	udpAddr, err := net.ResolveUDPAddr("udp", "239.0.0.0:"+bport)
	if err != nil {
		log.Printf("addr: %#v", udpAddr)
		panic(err)
	}
	fmt.Printf("Now listening broadcasts on port: %s\n", bport)

	s.conn, err = net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}
	incoming := make(chan *Message, 100)
	go s.listen(incoming)

	alive := true
	for alive {
		select {
		case msg := <-incoming:
			log.Printf("Got a message baby: %#v", msg)
		case msg := <-s.outgoing:
			if msg.Target >= len(s.connections) {
				break
			}
			conn := s.connections[msg.Target].Addr
			if n, err := s.conn.WriteToUDP(msg.Msg, conn); err != nil {
				fmt.Println("Error: ", err, " Bytes Written: ", n)
			}
		case <-exit:
			alive = false
			break
		default:
		}
	}
	fmt.Println("Killing Socket Server")
	s.conn.Close()
}

func (s *Network) listen(incoming chan *Message) {
	buf := make([]byte, 2048)
	for {
		n, ipaddr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		if n == 0 {
			continue
		}
		addr := ipaddr.String()
		connidx, ok := s.connLookup[addr]
		if !ok {
			connidx = len(s.connections)
			s.connLookup[addr] = connidx
			s.connections = append(s.connections, &Client{Addr: ipaddr, ID: connidx, Alive: true})
		}
		incoming <- &Message{
			Msg:    buf[:n],
			Target: connidx,
		}
	}
}

func (s *Network) handleMessages() {
	buf := make([]byte, 2048)
	for {
		n, ipaddr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		addr := ipaddr.String()
		connidx, ok := s.connLookup[addr]
		if !ok {
			connidx = len(s.connections)
			s.connLookup[addr] = connidx
			s.connections[len(s.connections)] = &Client{Addr: ipaddr, ID: connidx, Alive: true}
		}
		if n == 0 {
			continue
		}
		log.Printf("Message: %#v", buf[:n])
	}
}

type Message struct {
	Msg    []byte
	Target int
}

func (s *Network) sendMessages() {

}
