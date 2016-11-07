package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	listen      = flag.String("listen", "localhost:8080", "listen address")
	destination = flag.String("dst", "localhost:80", "destination address")
)

func main() {
	flag.Parse()
	log.Println("listening on", *listen)
	log.Println("relay to", *destination)
	ln, err := net.Listen("tcp", *listen)
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	dst, err := net.Dial("tcp", *destination)
	if err != nil {
		log.Print(err)
		return
	}
	dstName := dst.RemoteAddr()
	srcName := conn.RemoteAddr()
	go func() {
		io.Copy(W{dst, fmt.Sprintf("%s -> %s", srcName, dstName)}, conn)
		dst.Close()
	}()
	go func() {
		io.Copy(W{conn, fmt.Sprintf("%s -> %s", dstName, srcName)}, dst)
		conn.Close()
	}()
}

type W struct {
	io.Writer
	name string
}

func (w W) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	fmt.Printf("%s: %q\n", w.name, p[:n])
	fmt.Printf("	%v\n", p[:n])
	return n, err
}
