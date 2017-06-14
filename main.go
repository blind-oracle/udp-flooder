package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

var (
	Dst         = *flag.String("dst", "", "IP:Port where to send UDP packets")
	PayloadPath = *flag.String("payload", "", "Path to file with UDP payload")
	Threads     = *flag.Int("threads", 10, "Number of threads to start")
	Count       = *flag.Int("count", 1000, "Number of packets to send per thread")

	DstHost net.IP
	DstPort int
	Payload []byte
)

func UDPThread() {
	var (
		conn *net.UDPConn
		err  error
	)

	if conn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0}); err != nil {
		log.Fatalf("%s", err)
	}

	for i := 0; i < Count; i++ {
		if _, err = conn.WriteToUDP(Payload, &net.UDPAddr{IP: DstHost, Port: DstPort}); err != nil {
			log.Fatalf("Unable to send UDP: %s", err)
		}
	}
}

func main() {
	var err error
	flag.Parse()

	if Dst == "" {
		log.Fatal("Specify dst")
	}

	if PayloadPath == "" {
		log.Fatal("Specify payload")
	}

	t := strings.Split(Dst, ":")
	if len(t) != 2 {
		log.Fatal("Dst must be host:port")
	}

	if DstHost = net.ParseIP(t[0]); DstHost == nil {
		log.Fatal("Unable to parse dst host")
	}

	if DstPort, err = strconv.Atoi(t[1]); err != nil {
		log.Fatal("Unable to parse dst port")
	}

	if Payload, err = ioutil.ReadFile(PayloadPath); err != nil {
		log.Fatalf("Unable to read payload: %s", err)
	}

	log.Printf("Starting %d threads, sending %d packets from each thread to %s", Threads, Count, Dst)

	for i := 0; i < Threads; i++ {

	}
}
