package cNetTask

import (
	uh "ccproject/udp_handler"
	"fmt"
)

func HandleUDP(udpAddr string, taskChannel chan []string) {

	// listner udp
	go uh.ListenUDP(udpAddr, taskChannel)

	for packet := range taskChannel{
		go handleUDPMessage(packet, taskChannel)
	}

}

func handleUDPMessage(packet []string, taskChannel chan <- []string) {

	fmt.Println("Received packet: ", packet)
	taskChannel <- packet
}