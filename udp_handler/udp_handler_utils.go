package udp_handler

import (
	"fmt"
	"net"
	//"os"
	"time"
)

func sendUDPPackets(con *net.UDPConn, p *Packet, destination string) {
	print("SendUDPPacket called. Destination: ",destination,"\n")
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

func sendWithRetransmission(con *net.UDPConn, p *Packet, destination string, con_id string, sequence_no int) {
    udpAddr, err := net.ResolveUDPAddr("udp", destination)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return
	}
	sendWithRetransmission_(con,p,udpAddr,con_id,sequence_no)
}

 var retransmission bool = false
func sendWithRetransmission_(con *net.UDPConn, p *Packet, udpAddr *net.UDPAddr, conID string, sequenceNo int) {
    // Send packet immediately
    fmt.Printf("[SEND] Initial transmission for connection %s, sequence %d\n", conID, sequenceNo)
    sendUDPPackets_(con, p, udpAddr)
    
    // Set up retransmission
	if (retransmission){

    go func() {
        tries := 0
        for {
            // Check if packet still needs retransmission
            if _, exists := server_data_states[conID][sequenceNo]; !exists {
                fmt.Printf("[SUCCESS] Packet acknowledged for connection %s, sequence %d\n", conID, sequenceNo)
                return
            }

            if tries >= maxRetries {
                fmt.Printf("[FAILED] Max retries reached for connection %s, sequence %d after %d attempts\n", 
                    conID, sequenceNo, tries)
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
}

 