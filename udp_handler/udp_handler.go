package udp_handler

import (
	"fmt"
	"net"
	"os"
)

var connection_states map[string]int = make(map[string]int)

func SetConnState(conn string, state int){
	connection_states[conn] = state
}


/* var connection_data map[string]bytes = make(map[string]bytes) 
 */


func ListenServer(address string, channel chan [] string){
	//on receive: (task, name, critical calues, frequency,ip)
	//connection_data[address+ip] = data;
	//connection_state[address+ip] = 1;
	//send packet to receiver with SYN
}



	

func ListenUdp(type_ string, address string, con *net.UDPConn , channel chan [] string ){
	var connection *net.UDPConn
	var dest_address *net.UDPAddr
	var source_address *net.UDPAddr
	var a [] string
	a = make([]string,3,3)
	channel <- a
	if (address != "" ){
		udp_address,error := net.ResolveUDPAddr("udp",address)
		dest_address = udp_address
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}
		source_address = udp_address
		connection_, error := net.ListenUDP("udp", udp_address)
		connection = connection_
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}
	}else{
		connection = con
		sourceAddr := con.LocalAddr().(*net.UDPAddr)
		source_address = sourceAddr
		/* dest_address_ := con.RemoteAddr().(*net.UDPAddr)
		dest_address = dest_address_ */
	}
	for {
		fmt.Println("loop")
		//p := read_udp_packet(connection)
		buf := make([]byte, 4096)
		n, addr, err := connection.ReadFromUDP(buf)
		dest_address = addr
		if err != nil {
			fmt.Println(err)
			return
		}

		// Add check for zero bytes
		if n == 0 {
			fmt.Println("Received empty packet, skipping")
			continue
		}

		fmt.Printf("Received %d bytes from %s\n", n, addr.String())

		packet,_ := Deserialize(buf[:n]);
		packet.Print()

		// Create connection id
        sourceIP := source_address.IP.String()
        sourcePort := source_address.Port
		
        destIP := dest_address.IP.String()
        destPort := dest_address.Port
        connID := fmt.Sprintf("%s:%d:%s:%d", sourceIP, sourcePort, destIP, destPort)
		print("Connection ID:\n")
		print(connID)

        // Get current connection state
        state, exists := connection_states[connID]
        if !exists {
            state = 0
        }
		print("Received packet.")
		switch state {
			/// Initialize connection //////////////////////////////////////////////////////
			case 0: // No connection
				print("No connection on packet received.")
				if packet.Flags.SYN && !packet.Flags.ACK {
					// Send SYN+ACK
					response := &Packet{
						Type:           RegisterPacket,  
						SequenceNumber: 1,              
						AckNumber:      packet.SequenceNumber + 1,
						Flags: Flags{
							SYN: true,
							ACK: true,
							RET: false,
						},
						Data: AgentRegistration{      // Add this
							AgentID: "server-001",    // Use appropriate ID
							IPv4:    net.ParseIP("127.0.0.1"), // Use appropriate IP
						},
					}
					print("PRINT")
					packet.Print()
					serialized, _ := response.Serialize()
					print("serialized.length\n")
					print("ADDRESS sent:")
					print(addr.String())
					print(len(serialized))
					connection.WriteToUDP(serialized, addr)
					connection_states[connID] = 2
					print("Sent Syn + ACK")
				} else {
					fmt.Printf("Expected SYN on connection state 0: %s\n", connID)
				}
			
			case 1: // Sender: sent SYN
				print("Sender: Received packet after sending SYN")
				if packet.Flags.SYN && packet.Flags.ACK {
					// Send ACK
					response := &Packet{
						Type:           RegisterPacket,  
						SequenceNumber: 1,              
						AckNumber:      packet.SequenceNumber + 1,
						Flags: Flags{
							SYN: false,
							ACK: true,
							RET: false,
						},
						Data: AgentRegistration{      // Add this
							AgentID: "server-001",    // Use appropriate ID
							IPv4:    net.ParseIP("127.0.0.1"), // Use appropriate IP
						},
					}
					serialized, _ := response.Serialize()
					connection.WriteToUDP(serialized, addr)
					connection_states[connID] = 3
					print("Sender: Sent ACK")
				}

			case 2: // Receiver: sent SYN + ACK
				print("Receiver: received packet after sending SYN + ACK")
				if packet.Flags.ACK && !packet.Flags.SYN {
					connection_states[connID] = 4
					print("Receiver: received packet after sending SYN + ACK")
				}

			/// Send Data //////////////////////////////////////////////////////
			case 3: // Sender: Sent ACK/Packet
   				//this should be an ACK
				//remove timeout, packet from buffer 
			case 4: // Receiver: established connection 
				print("Receiver: received packet after sending SYN + ACK\n")
				var a [] string
				print(packet.Type.String())
				if(packet.Type == ReportPacket){
					reports := packet.Data.([]ReportRecord)
					length := len(reports)
					a = make([]string, length,length)
        
					// Process each report if needed~
					print("\nReaceived packet w/ data: Length of ")
					print(length)
					for i, report := range reports {
						// Example: convert each report to string or process it
						a[i] = fmt.Sprintf("Report %d: %v", i, report)
						print(a[i])
					}
				}
				channel <- a
				//channel chan [] string 
				//deserialize packet
				//packet
				//if server 
					//write to channel
				// if client
					//add task 
			//// Finalize connection //////////////////////////////////////////////////////
			//sender sends FIN
			//receiver receives FIN, sends FIN+ACK
			case 5: //sender sent FIN, receives FIN + ACK, sends ACK
			case 6: //receiver sent FIN + ACK, receives ACK   				
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


