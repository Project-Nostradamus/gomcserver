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
	holdbuf := make([]byte, 100)
	sendbuf := make([]byte, 196)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received Packet: ", buf)

	//Server Ping & Pong
	if buf[0] != 254 && buf[1] == 1 {
		conn.Write(buf)
		fmt.Println("Pong")
		return
	}

	//Server Info
	if buf[1] == 0 {
		serverinfo :=
			Info{
				Version: Version{
					Name:     "1.8.9",
					Protocol: 47,
				},
				Players: Players{
					Max:    100,
					Online: 5,
					Sample: []Sample{
						Sample{
							Name: "Bob",
							ID:   "4566e69f-c907-48ee-8d71-d7ba5aa00d20",
						},
					},
				},
				Description: Description{
					Text: "Hi",
				},
				Favicon: "",
			}
		holdbuf, _ = json.Marshal(serverinfo)
		sendbuf = encodeVarInt(uint32(len(holdbuf) + 2))
		sendbuf = append(sendbuf, 0)
		sendbuf = append(sendbuf, holdbuf...)

		fmt.Println("Raw Json: ", string(holdbuf))
		fmt.Println("Sent Packet: ", sendbuf)
		fmt.Println("Length of Sent Packet: ", len(sendbuf)+2)
		conn.Write(sendbuf)
	}
}

func encodeVarInt(v uint32) (vi []byte) {
	num := uint32(v)
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		vi = append(vi, byte(b))
		if num == 0 {
			break
		}
	}
	return
}

type Info struct {
	Version     Version     `json:"version"`
	Players     Players     `json:"players"`
	Description Description `json:"description"`
	Favicon     string      `json:"favicon"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []Sample `json:"sample"`
}

type Sample struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Description struct {
	Text string `json:"text"`
}
