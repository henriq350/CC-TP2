package udp_handler

import (
    "fmt"
    "net"
    "strconv"
)

func ListenServer(channel chan []string, con *net.UDPConn) {
    for {
        // taskId, name, frequency, threshold, client_id, dest_ip, duration, packet_count
        a := <-channel
        if len(a) == 8 {
            taskID := a[0]
            name := a[1]
            frequency := a[2]
            threshold := a[3]
            clientID := a[4]
            destIP := a[5]
            duration := a[6]
            packetCount := a[7]

            localAddr := con.LocalAddr().(*net.UDPAddr)
            localIP := localAddr.IP.String()
            localPort := localAddr.Port

            destAddr, err := net.ResolveUDPAddr("udp", clientID)
            if err != nil {
                fmt.Printf("Error resolving destination address: %v\n", err)
                fmt.Println(clientID)
                continue
            }

            connstate := fmt.Sprintf("%s:%d:%s:%d", localIP, localPort, destAddr.IP.String(), destAddr.Port)
            sequence := last_sequence_number[connstate]
            state := connection_states[connstate]

            send := getTaskPacket(clientID, taskID, name, frequency, threshold, duration, packetCount, destIP, sequence)

            if state == 3 || state == 4 {
                sendUDPPacket(con, send, clientID)
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
                
                sendUDPPacket(con, packet, clientID)
            }
        }
    }
}

func getTaskPacket(clientID, taskID, name, frequency, threshold, duration, packetCount, destIP string, sequence uint32) *Packet {
    freq, _ := strconv.ParseUint(frequency, 10, 32)
    dur, _ := strconv.ParseUint(duration, 10, 32)
    pCount, _ := strconv.ParseUint(packetCount, 10, 32)
    thresh, _ := strconv.ParseFloat(threshold, 64)

    tr := TaskRecord{
        TaskID:         taskID,
        Name:           name,
        Value:          "0",
        DestinationIp:  destIP,      // Added destinationIp field
        Threshold:      thresh,
        Duration:       uint32(dur),
        PacketCount:    uint32(pCount),
        Frequency:      uint32(freq),
        ReportFreq:     uint32(freq),
        CriticalValues: []string{threshold},
        ClientID:       clientID,
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