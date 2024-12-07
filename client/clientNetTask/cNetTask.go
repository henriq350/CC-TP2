package cNetTask

import (
	uh "ccproj/udp_handler"
	"fmt"
	"net"
	"os"
)

func HandleUDP(udpAddr, udpServerIP string, taskChannel chan []string, sendChannel chan []string, term chan bool) {

	// listner udp
	udp_address, error := net.ResolveUDPAddr("udp", udpAddr)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	connection, error := net.ListenUDP("udp", udp_address)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

    go uh.ListenUdp("","",connection ,taskChannel, term)
    go uh.ListenClient(udpServerIP,sendChannel, connection)

	select {}
}

