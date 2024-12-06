package sNetTask

import (
	"ccproj/server/db"
	uh "ccproj/udp_handler"
	"ccproj/server/types"
	"fmt"
	"sync"
	"time"
	"net"
	"os"
)

var agentMutex sync.Mutex

func HandleUDP(udpAddr string, agents map[string]types.Agent, lm *db.LogManager, sendChannel, receiveChannel chan []string) {

	udp_address,error := net.ResolveUDPAddr("udp",udpAddr)
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}

	connection_, _ := net.ListenUDP("udp", udp_address)
	
	go uh.ListenUdp("","",connection_ ,receiveChannel)
	go uh.ListenServer(sendChannel,connection_)

	//Receber mensagem e decidir o q fazer com ela
	for packet := range receiveChannel {
		fmt.Printf("Received packet: %v\n", packet)
		go handleUDPMessage(packet, agents, lm)
	}
	// Envia mensagem (com a task e ACKs) para os agents

}

// packet - "client_id" ,”task_id”  ,"tipo"    ,"metrica" ,"valor"  ,”client_ip”,"dest_ip"
// packet -  packet[0] , packet[1] ,packet[2] ,packet[3] packet[4] ,packet[5]   ,packet[6]

func handleUDPMessage(packet []string, agents map[string]types.Agent, lm *db.LogManager) {

	switch packet[2] {
		case "Register":
			fmt.Println("REGISTER!\nREGISTER!\nREGISTER!")
			// Cria um agente a partir do pacote
			agent := types.Agent{AgentID: packet[0], AgentIP: packet[5]}
			currentTime := time.Now().Format("2006-01-02 15:04:05")

			// Adiciona a lista de agentes
			agentMutex.Lock()
			types.AddAgent(agent, agents)
			agentMutex.Unlock()

			// Adiciona Log
			log := fmt.Sprintf("Agent %s registered", agent.AgentID)
			lm.AddLog(agent.AgentID, log, currentTime)
			
		case "Report":
			
			agentID := packet[0]

			aux := packet[1]
			metrics := packet[3:]
			metrics = append(metrics, aux)

			formatedString, currentTime := db.FormatString(metrics) 

			filename := currentTime

			db.StringToFile(agentID , filename, formatedString)

			// Adiciona nos Logs
			log := "Package received"
			lm.AddLog(agentID ,log, currentTime)

		case "Terminate":
			
			agent := types.Agent{AgentID: packet[0], AgentIP:  packet[5]}
			currentTime := time.Now().Format("2006-01-02 15:04:05")

			
			agentMutex.Lock()
			types.RemoveAgent(agent, agents)
			agentMutex.Unlock()

			// Escreve no log e remove o buffer do maps de logs
			log := fmt.Sprintf("Agent %s disconnected...", agent.AgentID)
			lm.AddLog(agent.AgentID, log, currentTime)
			lm.RemoveClientBuffer(agent.AgentID)
	}
}


