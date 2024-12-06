package udp_handler

import (
    "fmt"
    "net"
)

/* func ListenClient(server_ip string, channel chan []string, con *net.UDPConn) {
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
                last_sequence_number[connstate]++
                
                if _, exists := server_data_states[connstate]; !exists {
                    server_data_states[connstate] = make(map[int]Packet)
                }
                server_data_states[connstate][int(sequence)] = *send
                //sendUDPPacket(con, send, clientIP)
                sendWithRetransmission(con,send,clientIP,connstate,int(sequence))
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
                
                sendWithRetransmission(con,packet,clientIP,connstate,1)
                //sendUDPPacket(con, packet, destIP)
            }
        }
    }
}
 */

 func ListenClient(server_ip string, channel chan []string, con *net.UDPConn) {
    fmt.Println("[ListenClient] Started, waiting for channel messages...")
    
    for {
        // Wait for and receive message from channel
        a := <-channel
        fmt.Println("[ListenClient] Received message from channel")
 
        if len(a) == 7 {
            // Parse message fields
            clientID := a[0]
            taskID := a[1]
            tipo := a[2]
            metrica := a[3] 
            valor := a[4]
            destIP := a[6]
 
            fmt.Printf("[ListenClient] Message details:\n  Type: %s\n  Client ID: %s\n  Task ID: %s\n", 
                tipo, clientID, taskID)
 
            // Get local connection info
            localAddr := con.LocalAddr().(*net.UDPAddr)
            localIP := localAddr.IP.String()
            localPort := localAddr.Port
            fmt.Printf("[ListenClient] Local address: %s:%d\n", localIP, localPort)
 
            // Resolve server address
            destAddr, err := net.ResolveUDPAddr("udp", server_ip)
            if err != nil {
                fmt.Printf("[ListenClient] Error resolving destination address: %v\n", err)
                continue
            }
            fmt.Printf("[ListenClient] Destination address resolved: %s\n", destAddr.String())
 
            // Create connection identifier
            connstate := fmt.Sprintf("%s:%d:%s:%d", localIP, localPort, destAddr.IP.String(), destAddr.Port)
            sequence := last_sequence_number[connstate]
            state := connection_states[connstate]
            clientIP := fmt.Sprintf("%s:%d", localIP, localPort)
            
            fmt.Printf("[ListenClient] Connection state:\n  ID: %s\n  State: %d\n  Sequence: %d\n", 
                connstate, state, sequence)
 
            // Create packet based on type
            var send *Packet
            switch tipo {
            case "Report":
                send = getReportPacket(clientID, taskID, metrica, valor, destIP, sequence)
                fmt.Println("[ListenClient] Created Report packet")
            case "Register":
                send = getRegisterPacket(clientID, clientIP, "", sequence)
                fmt.Println("[ListenClient] Created Register packet")
            }
 
            if state == 3 || state == 4 {
                fmt.Println("[ListenClient] Connection established, sending data directly")
                last_sequence_number[connstate]++
                
                if _, exists := server_data_states[connstate]; !exists {
                    server_data_states[connstate] = make(map[int]Packet)
                    fmt.Println("[ListenClient] Initialized new server data state map")
                }
                
                server_data_states[connstate][int(sequence)] = *send
                fmt.Printf("[ListenClient] Stored packet with sequence %d\n", sequence)
                
                sendWithRetransmission(con, send, server_ip, connstate, int(sequence))
                fmt.Printf("[ListenClient] Sent packet with retransmission, sequence %d\n", sequence)
            } else {
                fmt.Println("[ListenClient] Connection not established, starting handshake")
                packet := &Packet{
                    Type:           RegisterPacket,
                    SequenceNumber: 1,
                    AckNumber:      1/* send.SequenceNumber + 1 */,
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
                fmt.Println("[ListenClient] Set initial sequence number and connection state")
                
                if _, exists := server_data_states[connstate]; !exists {
                    server_data_states[connstate] = make(map[int]Packet)
                    fmt.Println("[ListenClient] Initialized new server data state map")
                }
                
                server_data_states[connstate][1] = *packet
                server_data_states[connstate][4] = *send
                fmt.Println("[ListenClient] Stored SYN packet and data packet")
                
                sendWithRetransmission(con, packet, server_ip, connstate, 1)
                fmt.Println("[ListenClient] Sent SYN packet with retransmission")
            }
        }
    }
 }


func getReportPacket(clientID, taskID, name, value, destIP string, sequence uint32) *Packet {
    print("Sending report packet. Sequence: ",sequence, ".\n")
    return &Packet{
        Type:           ReportPacket,
        SequenceNumber: sequence,
        AckNumber:      60000,
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
    print("Sending register packet. Sequence: ",sequence, ".\n")
    return &Packet{
        Type:           RegisterPacket,
        SequenceNumber: sequence,
        AckNumber:      60000,
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