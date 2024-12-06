package udp_handler

import (
	"fmt"
	"net"
	"os"
	//"time"
)

func SetConnState(conn string, state int){
	connection_states[conn] = state
}

var connection_states map[string]int = make(map[string]int)
// Map of connection ID to map of sequence numbers to packets
var server_data_states map[string]map[int]Packet = make(map[string]map[int]Packet)
// Map of connection ID to last sequence number
var last_sequence_number map[string]uint32 = make(map[string]uint32)



func ListenUdp(type_ string, address string, con *net.UDPConn, channel chan []string) {
	fmt.Printf("[ListenUDP] Starting UDP listener. Type: %s, Address: %s\n", type_, address)
	
	var connection *net.UDPConn
	var dest_address *net.UDPAddr
	var source_address *net.UDPAddr
 
	// Set up connection and addresses
	if address != "" {
		fmt.Printf("[ListenUDP] Setting up new UDP connection on %s\n", address)
		udp_address, error := net.ResolveUDPAddr("udp", address)
		dest_address = udp_address
		if error != nil {
			fmt.Printf("[ListenUDP] Error resolving address: %v\n", error)
			os.Exit(1)
		}
		source_address = udp_address
		connection_, error := net.ListenUDP("udp", udp_address)
		connection = connection_
		if error != nil {
			fmt.Printf("[ListenUDP] Error creating UDP connection: %v\n", error)
			os.Exit(1)
		}
	} else {
		fmt.Println("[ListenUDP] Using provided connection")
		connection = con
		sourceAddr := con.LocalAddr().(*net.UDPAddr)
		source_address = sourceAddr
	}
 
	fmt.Printf("[ListenUDP] Listener configured. Source address: %v\n", source_address)
 
	for {
		fmt.Println("\n[ListenUDP] Waiting for incoming packets...")
		
		// Read incoming packet
		buf := make([]byte, 4096)
		n, addr, err := connection.ReadFromUDP(buf)
		dest_address = addr
		
		if err != nil {
			fmt.Printf("[ListenUDP] Error reading UDP packet: %v\n", err)
			return
		}
 
		if n == 0 {
			fmt.Println("[ListenUDP] Received empty packet, skipping")
			continue
		}
 
		fmt.Printf("[ListenUDP] Received %d bytes from %v\n", n, addr)
 
		// Deserialize packet
		packet, err := Deserialize(buf[:n])
		if err != nil {
			fmt.Printf("[ListenUDP] Deserialization error: %v\n", err)
			continue
		}
 
		// Create connection identifiers
		sourceIP := source_address.IP.String()
		sourcePort := source_address.Port
		destIP := dest_address.IP.String()
		destPort := dest_address.Port
		connID := fmt.Sprintf("%s:%d:%s:%d", sourceIP, sourcePort, destIP, destPort)
		destination_ipport := fmt.Sprintf("%s:%d", destIP, destPort)
 
		fmt.Printf("[ListenUDP] Connection details:\n")
		fmt.Printf("  Source: %s:%d\n", sourceIP, sourcePort)
		fmt.Printf("  Destination: %s:%d\n", destIP, destPort)
		fmt.Printf("  Connection ID: %s\n", connID)
 
		// Get connection state
		state, exists := connection_states[connID]
		if !exists {
			state = 0
			fmt.Println("[ListenUDP] New connection detected, initializing state to 0")
		} else {
			fmt.Printf("[ListenUDP] Existing connection, current state: %d\n", state)
		}
 
		fmt.Printf("[ListenUDP] Received packet details:\n")
		packet.Print()
 
		// State machine
		switch state {
		case 0: // No connection
			fmt.Println("[ListenUDP] State 0: Processing new connection")
			if packet.Flags.SYN && !packet.Flags.ACK {
				fmt.Println("[ListenUDP] Received SYN, sending SYN+ACK")
				response := &Packet{
					Type:           RegisterPacket,
					SequenceNumber: 2,
					AckNumber:      1,
					Flags: Flags{
						SYN: true,
						ACK: true,
						RET: false,
					},
					Data: AgentRegistration{
						AgentID:  "server-001",
						IPv4:     "127.0.0.1",
						ClientID: "1",
					},
				}
				sendUDPPackets_(connection, response, addr)
				connection_states[connID] = 2
				fmt.Println("[ListenUDP] Sent SYN+ACK, moved to state 2")
			} else {
				fmt.Printf("[ListenUDP] Error: Expected SYN flag on connection %s\n", connID)
			}
 
		case 1: // Sender: sent SYN
			fmt.Println("[ListenUDP] State 1: Processing SYN+ACK response")
			sequence := 3
			if packet.Flags.SYN && packet.Flags.ACK {
				fmt.Println("[ListenUDP] Received SYN+ACK, sending ACK")
				response := &Packet{
					Type:           RegisterPacket,
					SequenceNumber: uint32(sequence),
					AckNumber:      2,
					Flags: Flags{
						SYN: false,
						ACK: true,
						RET: false,
					},
					Data: AgentRegistration{
						AgentID:  "server-001",
						IPv4:     "127.0.0.1",
						ClientID: "1",
					},
				}
				sendWithRetransmission_(connection, response, addr, connID, sequence)
				connection_states[connID] = 3
				fmt.Println("[ListenUDP] Sent ACK, moved to state 3")
 
				// Check for pending data packets
				if packetMap, exists := server_data_states[connID]; exists {
					if send, exists := packetMap[4]; exists {
						fmt.Printf("[ListenUDP] Found pending packet (seq 4) for %s\n", connID)
						sendWithRetransmission_(connection, &send, addr, connID, sequence)
						fmt.Println("[ListenUDP] Sent pending packet")
					}
				}
			}
 
		case 2: // Receiver: sent SYN + ACK
			fmt.Println("[ListenUDP] State 2: Waiting for ACK")
			if packet.Flags.ACK && !packet.Flags.SYN {
				connection_states[connID] = 4
				last_sequence_number[connID] = 4
				fmt.Println("[ListenUDP] Received ACK, moved to state 4")
			}
 
		/* case 3: // Sender: Sent ACK/Packet
			fmt.Println("[ListenUDP] State 3: Processing ACK")
			if packet.Flags.ACK {
				ack := packet.AckNumber
				delete(server_data_states[connID], int(ack))
				fmt.Printf("[ListenUDP] Received ACK for sequence %d, removed from buffer\n", ack)
				//connection_states[connID] = 4
			} */
 
		case 3,4: // Receiver: established connection
			if packet.Flags.ACK == true {
				fmt.Println("[ListenUDP] State 3: Processing ACK")
				ack := packet.AckNumber
				delete(server_data_states[connID], int(ack))
				fmt.Printf("[ListenUDP] Received ACK for sequence %d, removed from buffer\n", ack)
				//connection_states[connID] = 4
			} else {
				fmt.Printf("[ListenUDP] State 4: Processing data packet type: %s\n", packet.Type)
				sequence := packet.SequenceNumber
				var a []string
				// Process packet based on type
				if packet.Type == TaskPacket {
					a = make([] string, 8,8)
					r := packet.Data.([]TaskRecord)[0]
					a[0] = r.ClientID
					a[1] = r.TaskID
					a[2] = r.Name
					a[3] = fmt.Sprint(r.ReportFreq)
					a[4] = r.CriticalValues[0]
					a[5] = r.DestinationIp
					a[6] = fmt.Sprint(r.Duration)    
					a[7] = fmt.Sprint(r.PacketCount)				// ... rest of TaskPacket processing ...
				} else if packet.Type == RegisterPacket {
					fmt.Println("[ListenUDP] Processing Register packet")
						a = make([]string,7,7)
						r:= packet.Data.(AgentRegistration)
						a[0] = r.ClientID
						a[1] = ""
						a[2] = "Register"
						a[3] = ""
						a[4] = ""
						a[5] = r.IPv4 
						a[6] = ""
				} else if packet.Type == ReportPacket {
					fmt.Println("[ListenUDP] Processing Report packet")
					reports := packet.Data.([]ReportRecord)						
					r := reports[0]
					a = make([]string,7,7)
					a[0] = r.ClientID
					a[1] = r.TaskID
					a[2] = "Report"
					a[3] = r.Name
					a[4] = r.Value
					a[5] = destination_ipport
					a[6] = r.DestinationIp
				}
				last_sequence_number[connID] = uint32(sequence+2)
				// Send ACK
				response := &Packet{
					Type:           RegisterPacket,
					SequenceNumber: sequence + 1,
					AckNumber:      sequence,
					Flags: Flags{
						SYN: false,
						ACK: true,
						RET: false,
					},
					Data: AgentRegistration{
						AgentID:  "server-001",
						IPv4:     "127.0.0.1",
						ClientID: "1",
					},
				}
				sendUDPPackets_(connection, response, addr)
				fmt.Printf("[ListenUDP] Sent ACK for sequence %d\n", sequence)
				
				channel <- a
				fmt.Println("[ListenUDP] Sent processed data to channel")
			}
		}
	}
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