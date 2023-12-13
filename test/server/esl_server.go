package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:18021")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		log.Fatalf("accept: %v", err)
	}
	defer conn.Close()
	s := bufio.NewScanner(conn)
	for s.Scan() {
		// t := s.Text()
	}
}
