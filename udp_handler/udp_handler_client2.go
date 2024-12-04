package udp_handler

import (
    "fmt"
    "net"
)

func ListenClient(server_ip string, channel chan []string, con *net.UDPConn) {
    fmt.Println("Listen Client started.")
    for {
        a := <-channel
        if len(a) == 7 {
            clientID := a[0]
            taskID := a[1]
            tipo := a[2]
            metrica := a[3]
            valor := a[4]
            //clientIP := a[5]
            destIP := a[6]

            localAddr := con.LocalAddr().(*net.UDPAddr)
            localIP := localAddr.IP.String()
            localPort := localAddr.Port

            destAddr, err := net.ResolveUDPAddr("udp", server_ip)
            if err != nil {
                fmt.Printf("Error resolving destination address: %v\n", err)
                continue
            }

            connstate := fmt.Sprintf("%s:%d:%s:%d", localIP, localPort, destAddr.IP.String(), destAddr.Port)
            sequence := last_sequence_number[connstate]
            state := connection_states[connstate]
            clientIP := fmt.Sprintf("%s:%d", localIP, localPort)

            var send *Packet
            switch tipo {
            case "Report":
                send = getReportPacket(clientID, taskID, metrica, valor, destIP, sequence)
            case "Register":
                send = getRegisterPacket(clientID, clientIP, metrica, sequence)
            }

            if state == 3 || state == 4 {
                sendUDPPacket(con, send, clientIP)
                last_sequence_number[connstate]++
                
                if _, exists := server_data_states[connstate]; !exists {
                    server_data_states[connstate] = make(map[int]Packet)
                }
                server_data_states[connstate][int(sequence)] = *send
                
            } else {
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
                        AgentID:  "server-001",
                        IPv4:     "127.0.0.1",
                        ClientID: clientID,
                    },
                }
                
                last_sequence_number[connstate] = 1
                connection_states[connstate] = 1
                
                if _, exists := server_data_states[connstate]; !exists {
                    server_data_states[connstate] = make(map[int]Packet)
                }
                server_data_states[connstate][1] = *packet
                server_data_states[connstate][4] = *send
                
                sendUDPPacket(con, packet, destIP)
            }
        }
    }
}

func getReportPacket(clientID, taskID, name, value, destIP string, sequence uint32) *Packet {
    return &Packet{
        Type:           ReportPacket,
        SequenceNumber: sequence,
        AckNumber:      1,
        Flags: Flags{
            SYN: false,
            ACK: false,
            RET: false,
        },
        Data: []ReportRecord{{
            TaskID:        taskID,
            Name:         name,
            Value:        value,
            DestinationIp: destIP,  // Added destinationIp
            ClientID:     clientID,
        }},
    }
}

func getRegisterPacket(clientID, clientIP, name string, sequence uint32) *Packet {
    return &Packet{
        Type:           RegisterPacket,
        SequenceNumber: sequence,
        AckNumber:      1,
        Flags: Flags{
            SYN: false,
            ACK: false,
            RET: false,
        },
        Data: AgentRegistration{
            AgentID:  name,
            IPv4:     clientIP,
            ClientID: clientID,
        },
    }
}