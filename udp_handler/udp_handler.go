package udp_handler

import (
	"fmt"
	"net"
	"os"
)

var connection_states map[string]int = make(map[string]int)

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
		//p := read_udp_packet(connection)
		buf := make([]byte, 4096)
		n, addr, err := connection.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			return
		}


		fmt.Printf("Received %d bytes from %s\n", n, addr.String())

		packet,_ := Deserialize(buf[:n]);

		// Create connection id
        sourceIP := addr.IP.String()
        sourcePort := addr.Port
        destIP := udp_address.IP.String()
        destPort := udp_address.Port
        connID := fmt.Sprintf("%s:%d:%s:%d", sourceIP, sourcePort, destIP, destPort)

        // Get current connection state
        state, exists := connection_states[connID]
        if !exists {
            state = 0
        }

		switch state {
			case 0: // No connection
				if packet.Flags.SYN && !packet.Flags.ACK {
					// Send SYN+ACK
					f := Flags{SYN:true,ACK:true,RET:false}
					response := Packet{Flags:f}
					serialized, _ := response.Serialize()
					connection.WriteToUDP(serialized, addr)
					connection_states[connID] = 2
				} else {
					fmt.Printf("Expected SYN on connection state 0: %s\n", connID)
				}
			
			case 1: // Sender: sent SYN
				if packet.Flags.SYN && packet.Flags.ACK {
					// Send ACK
					f := Flags{SYN: false, ACK: true, RET: false}
					response := Packet{Flags: f}
					serialized, _ := response.Serialize()
					connection.WriteToUDP(serialized, addr)
					connection_states[connID] = 3
				}

			case 2: // Receiver: sent SYN + ACK
				if packet.Flags.ACK && !packet.Flags.SYN {
					connection_states[connID] = 4
				}

			/* case 3: // Sender: sent ACK
   				//sendData(connection, addr, packet)

			case 4: // Receiver: received ACK 
   				//receiveData(connection, addr, packet) */
}
			
		}

		

		

		//identify connection (source ip + source port + dest ip + dest port)

		/* sourceIP := addr.IP.String()
		sourcePort := addr.Port
 */
		//connection := fmt.Sprintf("%s:%s:%d",address,sourceIP, sourcePort);
	
	
}

/* 
func read_udp_packet(conn *net.UDPConn) Packet{
	//var buf [512]byte
	buf := make([]byte, 4096)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			return
		}


		fmt.Printf("Received %d bytes from %s\n", n, addr.String())

		packet,_ := Deserialize(buf[:n]);
		return *packet;
} */


