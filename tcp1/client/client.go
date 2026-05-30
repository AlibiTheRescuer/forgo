package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	var m string
	if len(os.Args) > 1 {
		m = os.Args[1]
	}
	fmt.Fprintf(conn, "%s", "connection message from client: "+m)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(buf[:n]))
	conn.Close()
}
