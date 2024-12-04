package main

import (
	"ccproj/udp_handler"
	"fmt"
	"net"
	"os"
	"time"
)

func createRegistrationPacket() *udp_handler.Packet {
	return &udp_handler.Packet{
		Type:           udp_handler.RegisterPacket,
		SequenceNumber: 1,
		AckNumber:      0,
		Flags: udp_handler.Flags{
			SYN: true,
			ACK: false,
			RET: false,
		},
		Data: udp_handler.AgentRegistration{
			AgentID: "agent-001",
			IPv4:    net.ParseIP("192.168.1.100"),
		},
	}
}

func createTaskPacket() *udp_handler.Packet {
	tasks := []udp_handler.TaskRecord{
		{
			Name:       "cpu_usage",
			Value:      "monitor",
			ReportFreq: 30,
			CriticalValues: []string{
				"80",  // CPU usage threshold
				"90",  // Critical CPU threshold
			},
		},
		{
			Name:       "network_latency",
			Value:      "ping 8.8.8.8",
			ReportFreq: 60,
			CriticalValues: []string{
				"100", // Warning latency (ms)
				"200", // Critical latency (ms)
			},
		},
	}

	return &udp_handler.Packet{
		Type:           udp_handler.TaskPacket,
		SequenceNumber: 2,
		AckNumber:      0,
		Flags: udp_handler.Flags{
			SYN: false,
			ACK: false,
			RET: false,
		},
		Data: tasks,
	}
}

func createReportPacket() *udp_handler.Packet {
	reports := []udp_handler.ReportRecord{
		{
			Name:  "cpu_usage",
			Value: "45.2",
		},
		{
			Name:  "network_latency",
			Value: "23.5",
		},
	}

	return &udp_handler.Packet{
		Type:           udp_handler.ReportPacket,
		SequenceNumber: 3,
		AckNumber:      0,
		Flags: udp_handler.Flags{
			SYN: true,
			ACK: false,
			RET: false,
		},
		Data: reports,
	}
}

func sendPacket(conn *net.UDPConn, address *net.UDPAddr, packet *udp_handler.Packet) error {
	serializedData, err := packet.Serialize()
	if err != nil {
		return fmt.Errorf("serialization error: %v", err)
	}

	_, err = conn.WriteToUDP(serializedData,address)
	if err != nil {
		return fmt.Errorf("send error: %v", err)
	}

	fmt.Printf("Sent packet of type %v with sequence number %d\n", packet.Type, packet.SequenceNumber)
	return nil
}

func main() {
	/* serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8008")
	if err != nil {
		fmt.Printf("Address resolution error: %v\n", err)
		os.Exit(1)
	} */

	// Create a local address for the client
	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:54310")
	if err != nil {
		fmt.Printf("Local address resolution error: %v\n", err)
		os.Exit(1)
	} 
	//defer conn.Close()
	// Create separate listening connection
    listenAddr := localAddr
    listenConn, err := net.ListenUDP("udp", listenAddr)
    if err != nil {
        fmt.Printf("Listen error: %v\n", err)
        os.Exit(1)
    } 

	 ch := make(chan []string)
	/*		go func(){
			for{
				s := <-ch
				s = s
			}
			
		}() */
	go udp_handler.ListenUdp("client","",listenConn,ch);
	go udp_handler.ListenClient(ch,listenConn);

	var a [] string = make([]string,8,8)/* 
	“taskId”,"name","frequencia","threshold",”client_ip”,"dest_ip",”duration”,”packet_count” */
	a = make([] string, 8,8)
	a[0] = "0"
	a[1] = "1"
	a[2] = "CPU"
	a[3] = "30"
	a[4] = "127.0.0.1:8008"
	a[5] = "127.0.0.1:8007" 
	a[6] = "10"
	a[7] = "10"
	ch <- a
	print("sent to channel")

/* 
	udp_handler.SetConnState("127.0.0.1:54310:127.0.0.1:8008",1)


	// Send report packet
	fmt.Println("\nSending Report Packet...")
	reportPacket := createReportPacket()
	err = sendPacket(listenConn, serverAddr,reportPacket)
	if err != nil {
		fmt.Printf("Failed to send report packet: %v\n", err)
	} 


	time.Sleep(3 * time.Second)


	// Send report packet
	fmt.Println("\nSending Report Packet...")
	err = sendPacket(listenConn, serverAddr,reportPacket)
	if err != nil {
		fmt.Printf("Failed to send report packet: %v\n", err)
	}  */
	// Send registration packet
	/* fmt.Println("\nSending Registration Packet...")
	regPacket := createRegistrationPacket()
	err = sendPacket(listenConn, serverAddr,regPacket)
	if err != nil {
		fmt.Printf("Failed to send registration packet: %v\n", err)
	} */

	select{}


	 time.Sleep(1 * time.Second)
/*
	// Send task packet
	fmt.Println("\nSending Task Packet...")
	taskPacket := createTaskPacket()
	err = sendPacket(conn, taskPacket)
	if err != nil {
		fmt.Printf("Failed to send task packet: %v\n", err)
	}

	time.Sleep(1 * time.Second)*/

	// Send report packet
	/* fmt.Println("\nSending Report Packet...")
	reportPacket := createReportPacket()
	err = sendPacket(conn, reportPacket)
	if err != nil {
		fmt.Printf("Failed to send report packet: %v\n", err)
	}  */
}