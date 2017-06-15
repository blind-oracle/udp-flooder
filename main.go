package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	Dst         string
	PayloadPath string
	Threads     int
	Count       int
	DelayStr    string

	DstHost net.IP
	DstPort int
	Delay   time.Duration
	Payload []byte
	wg      sync.WaitGroup
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

		if Delay > 0 {
			time.Sleep(Delay)
		}
	}

	wg.Done()
}

func main() {
	var err error

	flag.StringVar(&Dst, "dst", "", "IP:Port where to send UDP packets")
	flag.StringVar(&PayloadPath, "payload", "", "Path to file with UDP payload")
	flag.StringVar(&DelayStr, "delay", "0s", "Time between consequent packets")
	flag.IntVar(&Threads, "threads", 10, "Number of threads to start")
	flag.IntVar(&Count, "count", 1000, "Number of packets to send per thread")
	flag.Parse()

	if Dst == "" {
		log.Fatal("Specify dst")
	}

	if PayloadPath == "" {
		log.Fatal("Specify payload")
	}

	if Delay, err = time.ParseDuration(DelayStr); err != nil {
		log.Fatal("Unable to parse %s as duration", DelayStr)
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
		wg.Add(1)
		go UDPThread()
	}

	wg.Wait()
}
