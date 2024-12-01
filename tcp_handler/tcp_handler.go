package tcp_handler

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

type AlertMetric int

const (
	CPUUsage AlertMetric = iota + 1
    RAMUsage
    InterfaceStats
    PacketLoss
    Jitter
)

type AlertMessage struct {
	AgentID string
	AlertMetric AlertMetric
	Threshold float32
	Value float32
}


func (am AlertMetric) String() string {
    switch am {
    case CPUUsage:
        return "CPU Usage"
    case RAMUsage:
        return "RAM Usage"
    case InterfaceStats:
        return "Interface Stats"
    case PacketLoss:
        return "Packet Loss"
    case Jitter:
        return "Jitter"
    default:
        return "Unknown"
    }
}

func ListenTcp(address string, alertChan chan <- AlertMessage) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor TCP:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("TCP Listener started on ", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		
		go handleTcpConnection(conn, alertChan)
	}
}

func handleTcpConnection(conn net.Conn, alertChan chan<- AlertMessage) {
	defer conn.Close()
    reader := bufio.NewReader(conn)

    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
				fmt.Println("Connection closed by client:", conn.RemoteAddr())
			} else {
				fmt.Println("Error reading from connection:", err)
			}
            return
        }

		alert, err := DeserializeAlert([]byte(message))

		if err != nil {
			return
		}
        alertChan <- alert
    }
}

func SendAlert(address string, alert AlertMessage) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Eerror connecting to TCP server:", err)
		return
	}
	defer conn.Close()

	message, err := SerializeAlert(alert)

	if err != nil {
		return
	}

	message = append(message, '\n')

	// Enviar a mensagem
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
	
}