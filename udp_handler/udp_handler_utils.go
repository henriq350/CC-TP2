package udp_handler

import (
	"fmt"
	"net"
	//"os"
)

func sendUDPPacket(con *net.UDPConn, p *Packet, destination string) {
	udpAddr, err := net.ResolveUDPAddr("udp", destination)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return
	}

	serialized, err := p.Serialize()
	if err != nil {
		fmt.Printf("Error serializing packet: %v\n", err)
		return
	}

	_, err = con.WriteToUDP(serialized, udpAddr)
	if err != nil {
		fmt.Printf("Error sending UDP packet: %v\n", err)
		return
	}
}

func sendUDPPacket_(con *net.UDPConn, p *Packet,udpAddr *net.UDPAddr) {

	serialized, err := p.Serialize()
	if err != nil {
		fmt.Printf("Error serializing packet: %v\n", err)
		return
	}

	_, err = con.WriteToUDP(serialized, udpAddr)
	if err != nil {
		fmt.Printf("Error sending UDP packet: %v\n", err)
		return
	}
}
