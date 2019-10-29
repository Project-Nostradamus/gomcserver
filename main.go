package main

import (
	"fmt"
	"net"
	"os"
)

const (
	connHost = "localhost"
	connPort = "25565"
	connType = "tcp"
)

var IDCounter = 0
var playerMap = make(map[int]Player)

func main() {
	fmt.Println("Opening port of type", connType, "at", connHost+":"+connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		IDCounter++
		playerMap[IDCounter] = Player{
			uuid: [16]byte{},
			name: "",
			conn: conn,
		}
	}
}

type Player struct {
	uuid [16]byte
	name string
	conn net.Conn
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 100)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	fmt.Println(buf)
}
