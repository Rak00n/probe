package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func runTCPServer(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	for {
		userInput, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return
		}
		fmt.Println(userInput)
		sendData(userInput)
	}
}

func sendDataOverTCP(data string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", relayTo)
	if err != nil {
		println("Resolve failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(data))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
	conn.Close()
}
