package udp_handler

import (
	"fmt"
	"net"
	//"os"
	"time"
)

func sendUDPPackets(con *net.UDPConn, p *Packet, destination string) {
	//print("SendUDPPacket called. Destination: ",destination,"\n")
	udpAddr, err := net.ResolveUDPAddr("udp", destination)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return
	}

	sendUDPPackets_(con,p,udpAddr)
}

func sendUDPPackets_(con *net.UDPConn, p *Packet,udpAddr *net.UDPAddr) {

	serialized, err := p.Serialize()
	if err != nil {
		fmt.Printf("Error serializing packet: %v\n", err)
		return
	}

	_, err = con.WriteToUDP(serialized, udpAddr)
	if err != nil {
		fmt.Printf("Error sending UDP packet: %v\n", err)
		return
	}

}

var waitTime int = 7
var maxRetries int = 5

func sendWithRetransmission(con *net.UDPConn, p *Packet, destination string, con_id string, sequence_no int) bool{
	//print("[SEND] Convert address ",destination,".\n")
    udpAddr, err := net.ResolveUDPAddr("udp", destination)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return false
	}
	return sendWithRetransmission_(con,p,udpAddr,con_id,sequence_no)
}

 var retransmission bool = true
func sendWithRetransmission_(con *net.UDPConn, p *Packet, udpAddr *net.UDPAddr, conID string, sequenceNo int) bool {
    // Send packet immediately
    fmt.Printf("[SEND] Initial transmission for connection %s, sequence %d\n", conID, sequenceNo)
    sendUDPPackets_(con, p, udpAddr)
    //connection_states[conID] = 3
    // Set up retransmission
	sent := true
	if (retransmission){
    go func() {
        tries := 0
        for {
            // Check if packet still needs retransmission
            if _, exists := server_data_states[conID][sequenceNo]; !exists {
                fmt.Printf("[SUCCESS] Packet acknowledged for connection %s, sequence %d\n", conID, sequenceNo)
                sent = true
				return
            }
            if tries >= maxRetries {
                fmt.Printf("[FAILED] Max retries reached for connection %s, sequence %d after %d attempts\n", 
                conID, sequenceNo, tries)
				sent = false
				return
            }
            time.Sleep(time.Duration(waitTime) * time.Second)
            tries++
            fmt.Printf("[RETRY] Attempt %d/%d for connection %s, sequence %d\n", 
                tries, maxRetries, conID, sequenceNo)
            sendUDPPackets_(con, p, udpAddr)
        }
    }()
}
return sent
}

 