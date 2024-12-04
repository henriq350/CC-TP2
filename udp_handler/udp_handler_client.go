package udp_handler

import (
	"fmt"
	"net"
	//"os"
	//"time"
)


func ListenClient(channel chan []string, con *net.UDPConn) {/* 
	"client_id",”task_id”,“tipo”,"metrica","valor",”client_ip”,"dest_ip" */
	print("Listen Client started.\n")
	for {
		a := <-channel
		print("cleint:received array\n")
		print("Length: ", len(a),"\n")
		if len(a) == 7 {
			//client_id := a[0]
			//task_id := a[1]
			tipo := a[2]
			metrica := a[3]
			valor := a[4]
			client_ip := a[5]
			dest_ip := a[6]

			// Get local address details
			localAddr := con.LocalAddr().(*net.UDPAddr)
			localIP := localAddr.IP.String()
			localPort := localAddr.Port

			// Parse destination address for port
			destAddr, err := net.ResolveUDPAddr("udp", dest_ip)
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
			print("connsstate:", connstate,"\n")
			var sequence uint32 
			sequence = last_sequence_number[connstate]
			state := connection_states[connstate]

			var send *Packet

			if(tipo == "Report"){
				send = getReportPacket(client_ip,metrica,valor,sequence)
			}else if (tipo == "Register"){
				send = getRegisterPacket(client_ip,"0x01",sequence)
			}
			if state == 3 || state == 4 {
				print("connection found.\n")
				sendUDPPacket(con, send, client_ip)
				last_sequence_number[connstate]++
				
				if _, exists := server_data_states[connstate]; !exists {
					server_data_states[connstate] = make(map[int]Packet)
				}
				server_data_states[connstate][int(sequence)] = *send
				
			} else { // assuming case 0
				print("connection .\n")
				send.Print()
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
						IPv4:    "127.0.0.1",
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
				
				sendUDPPacket(con, packet, dest_ip)
			}
		}
	}
}




//Report Packet
func getReportPacket(client_ip string , name string , value string ,sequence uint32) *Packet {
	tr := ReportRecord{
		Name:           name,
		Value:          value,
	}
	return &Packet{
		Type:           ReportPacket,
		SequenceNumber: sequence,
		AckNumber:      1,
		Flags: Flags{
			SYN: false,
			ACK: false,
			RET: false,
		},
		Data: []ReportRecord{tr},
	}
}
//Register Packet
func getRegisterPacket(client_ip string , name string , sequence uint32) *Packet {
	
	tr := AgentRegistration{
		AgentID:           name,
		IPv4:          client_ip,
	}
	return &Packet{
		Type:           ReportPacket,
		SequenceNumber: sequence,
		AckNumber:      1,
		Flags: Flags{
			SYN: false,
			ACK: false,
			RET: false,
		},
		Data: []AgentRegistration{tr},
	}
}

/*
// AgentRegistration represents the registration data
type AgentRegistration struct {
	AgentID string
	IPv4    net.IP
}
*/
//Terminate