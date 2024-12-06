package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	listenAddr = flag.String("l", ":5001", "listen address")
	remoteAddr = flag.String("r", "10.0.0.10:5001", "remote address")
	verbose    = flag.Bool("v", false, "verbose on link")
)

func run(a, b net.Conn) {
	defer a.Close()
	defer b.Close()

	go func() {
		defer a.Close()
		defer b.Close()
		var buf [2048]byte
		for {
			n, err := io.CopyBuffer(a, b, buf[:])
			if err != nil || n == 0 {
				return
			}
		}
	}()

	var buf [2048]byte
	for {
		n, err := io.CopyBuffer(b, a, buf[:])
		if err != nil || n == 0 {
			return
		}
	}
}

func main() {
	flag.Parse()

	if *verbose {
		log.Println("linking: " + *listenAddr + " -> " + *remoteAddr)
	}
	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		lss, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		if *verbose {
			log.Println("Accepted", lss.RemoteAddr())
		}
		conn, err := net.Dial("tcp", *remoteAddr)
		if err != nil {
			log.Print(err)
			continue
		}
		go run(lss, conn)
	}
}
