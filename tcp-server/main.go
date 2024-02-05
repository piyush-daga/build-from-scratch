package main

import (
	"log"
	"net"
	"time"
)

func doWork(conn net.Conn) {
	log.Println("Starting work on the new connection.")

	// Setting up a connection timeout
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	data := make([]byte, 100)
	_, err := conn.Read(data)
	logError(err)

	// Do some work
	time.Sleep(5 * time.Second)

	// Write back the response
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello Hello!!\r\n"))

	conn.Close()
}

// TCP Server
func main() {
	// Listen to a port
	listener, err := net.Listen("tcp", ":9000")
	logError(err)

	for {
		log.Println("Listening for new connections !!")

		// Now we need to start accepting connections
		conn, err := listener.Accept()
		logError(err)

		log.Println("Accepted a new connection.")

		go doWork(conn)
	}
}

func logError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
