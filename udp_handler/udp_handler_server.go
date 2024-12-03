package udp_handler

import (
	"fmt"
	"net"
	//"os"
	//"time"
)


func ListenServer(channel chan []string, con *net.UDPConn) {
	for {
		a := <-channel
		if len(a) == 8 {
			//taskId := a[0]
			name := a[1]
			frequency := a[2]
			threshold := a[3]
			client_ip := a[4]
			//dest_ip := a[5]

			// Get local address details
			localAddr := con.LocalAddr().(*net.UDPAddr)
			localIP := localAddr.IP.String()
			localPort := localAddr.Port

			// Parse destination address for port
			destAddr, err := net.ResolveUDPAddr("udp", client_ip)
			if err != nil {
				fmt.Printf("Error resolving destination address: %v\n", err)
				print(client_ip)
				print("\n")
				continue
			}
			destIP := destAddr.IP.String()
			destPort := destAddr.Port

			// Create connection state identifier with both IPs and ports
			connstate := fmt.Sprintf("%s:%d:%s:%d", localIP, localPort, destIP, destPort)
			var sequence uint32 
			sequence = last_sequence_number[connstate]
			state := connection_states[connstate]

			send := getTaskPacket(client_ip, name, frequency, threshold, sequence)

			if state == 3 || state == 4 {
				sendUDPPacket(con, send, client_ip)
				last_sequence_number[connstate]++
				
				if _, exists := server_data_states[connstate]; !exists {
					server_data_states[connstate] = make(map[int]Packet)
				}
				server_data_states[connstate][int(sequence)] = *send
				
			} else { // assuming case 0
				packet := &Packet{
					Type:           RegisterPacket,
					SequenceNumber: 1,
					AckNumber:      send.SequenceNumber + 1,
					Flags: Flags{
						SYN: true,
						ACK: false,
						RET: false,
					},
					Data: AgentRegistration{
						AgentID: "server-001",
						IPv4:    net.ParseIP("127.0.0.1"),
					},
				}
				last_sequence_number[connstate] = 1
				connection_states[connstate] = 1
				
				if _, exists := server_data_states[connstate]; !exists {
					server_data_states[connstate] = make(map[int]Packet)
				}
				server_data_states[connstate][1] = *packet

				//update seq 4 with packet to send
				server_data_states[connstate][4] = *send
				
				sendUDPPacket(con, packet, client_ip)
			}
		}
	}
}

func getTaskPacket(client_ip, metrica, frequencia, threshold string, sequence uint32) *Packet {
	freq := 0
	// Convert frequency string to int (add error handling as needed)
	fmt.Sscanf(frequencia, "%d", &freq)
	tr := TaskRecord{
		Name:           metrica,
		Value:          "0",
		ReportFreq:     uint32(freq),
		CriticalValues: []string{threshold},
	}
	return &Packet{
		Type:           TaskPacket,
		SequenceNumber: sequence,
		AckNumber:      1,
		Flags: Flags{
			SYN: false,
			ACK: false,
			RET: false,
		},
		Data: []TaskRecord{tr},
	}
}