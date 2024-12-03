package udp_handler

import (
	"fmt"
	"net"
	"os"
	"time"
)

func SetConnState(conn string, state int){
	connection_states[conn] = state
}

var connection_states map[string]int = make(map[string]int)
// Map of connection ID to map of sequence numbers to packets
var server_data_states map[string]map[int]Packet = make(map[string]map[int]Packet)
// Map of connection ID to last sequence number
var last_sequence_number map[string]uint32 = make(map[string]uint32)






/*
key1: connection (source_ip+sourceport+dest_ip+dest_port
key2: sequence number
var server_data_states map[map[Packet]int]string = make(map[map[Packet]int]string )

key: connection 
var last_sequence_number map[int]string = make(map[int]string)

 
func ListenServer(channel chan [] string,con *net.UDPConn){
	“taskId”,"name","frequencia","threshold",”client_ip”,"dest_ip",”duration”,”packet_count"
	while (string []a <- channel)
		if (a.length == 8)
		//came from server
			string connstate = con.sourceaddr + dest_ip
			int sequence = last_sequence_number[connstate]
			state = connection_states[connstate]
			Packet send = getTaskPacket(a[4],a[1],a[3],sequence)
			if (connstate == 3 || connstate == 4){
				sendUDPPacket(con,send,client_ip)
				last_sequence_number[connstate]++
				server_data_states[connstate][sequence] = send
			}
			else { //assuming case 0
				packet := &Packet{
						Type:           RegisterPacket,  
						SequenceNumber: 1,              
						AckNumber:      packet.SequenceNumber + 1,
						Flags: Flags{
							SYN: true,
							ACK: false,
							RET: false,
						},
						Data: AgentRegistration{      // Add this
							AgentID: "server-001",    // Use appropriate ID
							IPv4:    net.ParseIP("127.0.0.1"), // Use appropriate IP
						},
					}
				last_sequence_number[connstate] = 1
				connection_state[connstate] = 1;
				server_data_states[connstate][1] = packet
				sendUDPPacket(con,packet,client_ip)
			}


	//if ()
	//connection_data[address+ip] = data;
	//connection_state[address+ip] = 1;
	//send packet to receiver with SYN
}
*/


/*func getTaskPacket(string client_ip,string metrica,string frequencia,string threshold,int sequence) (Packet *p){
	return &Packet{
						Type:           TaskPacket,  
						SequenceNumber: sequence,              
						AckNumber:      1,
						Flags: Flags{
							SYN: false,
							ACK: false,
							RET: false,
						},
						Data: TaskRecord {
							Name:metrica,
							Value:0,
							ReportFreq:int(frequencia)
							CriticalValues:[threshold]
						}
}*/


/*
func sendUDPPacket(con *net.UDPConn, Packet *p, string destination){
	udp_address,error := net.ResolveUDPAddr("udp",destination)
	serialized, _ := p.Serialize()
	print("serialized.length\n")
	con.WriteToUDP(serialized, addr)
}
*/


	

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

		//fmt.Printf("Received %d bytes from %s\n", n, addr.String())

		packet,_ := Deserialize(buf[:n]);
		//packet.Print()

		// Create connection id
        sourceIP := source_address.IP.String()
        sourcePort := source_address.Port
		
        destIP := dest_address.IP.String()
        destPort := dest_address.Port
        connID := fmt.Sprintf("%s:%d:%s:%d", sourceIP, sourcePort, destIP, destPort)
		destination_ipport := fmt.Sprintf("%s:%d",destIP,destPort)
		print("Connection ID:\n")
		print(connID)

        // Get current connection state
        state, exists := connection_states[connID]
        if !exists {
            state = 0
        }
		print("\nReceived packet.\n")
		switch state {
			/// Initialize connection //////////////////////////////////////////////////////
			case 0: // No connection
				print("No connection on packet received.\n")
				if packet.Flags.SYN && !packet.Flags.ACK {
					// Send SYN+ACK
					response := &Packet{
						Type:           RegisterPacket,  
						SequenceNumber: 2,              
						AckNumber:      1,
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
					 /* print("PRINT")
					packet.Print()*/
					serialized, _ := response.Serialize()
					/* print("serialized.length\n")
					print("ADDRESS sent:")
					print(addr.String())
					print(len(serialized)) */
					connection.WriteToUDP(serialized, addr)
					connection_states[connID] = 2
					print("Sent Syn + ACK\n")
				} else {
					fmt.Printf("Expected SYN on connection state 0: %s\n", connID)
				}
			
			case 1: // Sender: sent SYN
				print("Sender: Received packet after sending SYN\n")
				if packet.Flags.SYN && packet.Flags.ACK {
					// Send ACK
					response := &Packet{
						Type:           RegisterPacket,  
						SequenceNumber: 3,              
						AckNumber:      2,
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
					sendUDPPacket_(connection,response,addr)
					connection_states[connID] = 3
					print("Sender: Sent ACK\n")			
					time.Sleep(1*time.Second)	


 					if packetMap, exists := server_data_states[connID]; exists {
						// Check if sequence number exists in the connection's packet map
						if send, exists := packetMap[4]; exists {
							fmt.Printf("Found packet with sequence 4 for connection %s\n", connID)
							send.Print()
							
							sendUDPPacket_(connection, &send, addr)
							fmt.Println("Successfully sent packet with sequence 4\n")
						} else {
							fmt.Printf("No packet with sequence 4 found for connection %s\n", connID)
						}
					} else {
						fmt.Printf("No packets found for connection %s\n", connID)
					}
				}

			case 2: // Receiver: sent SYN + ACK
				print("Receiver: received packet after sending SYN + ACK\n")
				if packet.Flags.ACK && !packet.Flags.SYN {
					connection_states[connID] = 4
					print("Receiver: received packet after sending SYN + ACK\n")
				}

			/// Send Data //////////////////////////////////////////////////////
			case 3: // Sender: Sent ACK/Packet
   				//this should be an ACK
				//remove timeout, packet from buffer 

			case 4: // Receiver: established connection 
				print("Receiver: received packet after sending SYN + ACK\n")
				/*
				server
				"client_id",”task_id”,"tipo","metrica","valor",”client_ip”,"dest_ip"				
				client
				“clientId”,“taskId”,"name","frequencia","threshold","dest_ip",”duration”,”packet_count”
				*/
				print(packet.Type.String())
				sequence := packet.SequenceNumber;
				var a [] string
				 if (packet.Type == TaskPacket){
/* 					“clientId”,“taskId”,"name","frequencia","threshold","dest_ip",”duration”,”packet_count”
 */
					a = make([] string, 8,8)
					a[0] = "0"
					a[1] = "1"
					 r := packet.Data.([]TaskRecord)[0]
					a[2] = r.Name
					a[3] = fmt.Sprint(r.ReportFreq)
					a[4] = r.CriticalValues[0]
					a[5] = destination_ipport 
					a[6] = "10"
					a[7] = "10"
				} else {
/* 					"client_id",”task_id”,"tipo","metrica","valor",”client_ip”,"dest_ip"				
 */
					if(packet.Type == RegisterPacket){
						a = make([]string,7,7)
						a[0] = "0"
						a[1] = "1"
						r := packet.Data.(TaskRecord)
						a[2] = r.Name
						a[3] = fmt.Sprint(r.ReportFreq)
						a[4] = r.CriticalValues[0]
						a[5] = destination_ipport 
						a[6] = "10"
						a[7] = "10"
					}
					if(packet.Type == ReportPacket){
						reports := packet.Data.([]ReportRecord)
						length := len(reports)
						a = make([]string, length,length)
			
						// Process each report if needed~
					/* 	print("\nReaceived packet w/ data: Length of ")
						print(length) */
						for i, report := range reports {
							// Example: convert each report to string or process it
							a[i] = fmt.Sprintf("Report %d: %v", i, report)
							//print(a[i])
						}
					}
				}
					response := &Packet{
						Type:           RegisterPacket,  
						SequenceNumber: sequence+1,
						AckNumber:      sequence,
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
					sendUDPPacket_(connection,response,addr)
				channel <- a
			//// Finalize connection //////////////////////////////////////////////////////
			// sends FIN
			// receives FIN, sends FIN+ACK
			case 5: // sent FIN, receives FIN + ACK, sends ACK
			case 6: // sent FIN + ACK, receives ACK   				
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


