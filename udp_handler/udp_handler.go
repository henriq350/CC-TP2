package udp_handler

import (
	"fmt"
	"net"
	"os"
)

func ListenUdp(address string){
	udp_address,error := net.ResolveUDPAddr("udp",address);
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	connection, error := net.ListenUDP("udp", udp_address)

	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	for {
		fmt.Println("loop")
		read_udp_packet(connection)
	}

	
}


func read_udp_packet(conn *net.UDPConn){
	var buf [512]byte
		_, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("UDP packet received");
		buf2 := [512]byte{}

		packet,_ := Deserialize(buf2[:]);
		packet.Print();
}


