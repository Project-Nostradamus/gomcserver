package main

import (
	"encoding/json"
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
	listener, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
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

		handleRequest(conn)
	}
}

type Player struct {
	uuid [16]byte
	name string
	conn net.Conn
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 100)
	sendbuf := make([]byte, 100)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println(buf)

	//JSON!!!!!!
	serverinfo :=
		Info{
			version: Version{
				name:     "1.8.9",
				protocol: 47,
			},
			players: Players{
				max:    3,
				online: 0,
				sample: Sample{
					name: "Bob",
					id:   "4566e69f-c907-48ee-8d71-d7ba5aa00d20",
				},
			},
			description: Description{
				text: "Hi",
			},
			favicon: nil,
		}

	sendbuf, _ = json.Marshal(serverinfo)
	conn.Write(sendbuf)
}

type Info struct {
	version     Version
	players     Players
	description Description
	favicon     string
}

type Version struct {
	name     string
	protocol int
}

type Players struct {
	max    int
	online int
	sample Sample
}

type Sample struct {
	name string
	id   string
}

type Description struct {
	text string
}
