package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	udp_address, error := net.ResolveUDPAddr("udp", "127.0.0.1:8008")
	if error != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	connection, error := net.ListenUDP("udp", udp_address)
	conn.WriteToUDP([]byte("Hello UDP Client\n"), addr)
}
