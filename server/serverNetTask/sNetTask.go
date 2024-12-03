package sNetTask

import (
	"ccproj/server/db"
	"ccproj/server/main"
	uh "ccproj/udp_handler"
	"fmt"
	"sync"
	"time"
)

var agentMutex sync.Mutex

func HandleUDP(udpAddr string, agents map[string]uh.AgentRegistration, lm *db.LogManager) {

	receiveChannel := make(chan []string)

	// Listener de UDP
	go uh.ListenUDP(udpAddr, receiveChannel)

	//Receber mensagem e decidir o q fazer com ela
	for packet := range receiveChannel {
		go handleUDPMessage(packet, agents, lm)
	}
	// Envia mensagem (com a task e ACKs) para os agents

}

// packet - "client_id" ,”task_id”  ,"tipo"    ,"metrica" ,"valor"  ,”client_ip”,"dest_ip"
// packet -  packet[0] , packet[1] ,packet[2] ,packet[3] packet[4] ,packet[5]   ,packet[6]

func handleUDPMessage(packet []string, agents map[string]uh.AgentRegistration, lm *db.LogManager) {

	switch packet[2] {
		case "Register":
			
			// Cria um agente a partir do pacote
			agent := main.Agent{packet[0], packet[5]}
			currentTime := time.Now().Format("2024-11-14 15:04:05")

			// Adiciona a lista de agentes
			agentMutex.Lock()
			main.AddAgent(agent, agents)
			agentMutex.Unlock()

			// Adiciona Log
			log := fmt.Sprintf("[%s] Agent %s registered", currentTime, agent.AgentID)
			lm.AddLog(agent.AgentID, log, currentTime)
			
		case "Report":
			
			agentID := packet[0]

			aux := packet[1]
			metrics := packet[3:]
			metrics = append(metrics, aux)

			formatedString, currentTime := db.FormatString(metrics) 

			filename := fmt.Sprintf("%s", &currentTime)

			db.StringToFile(agentID , filename, formatedString)

			// Adiciona nos Logs
			log := fmt.Sprintf("Package received")
			lm.AddLog(agentID ,log, currentTime)

		case "Terminate":
			
			agent := main.Agent{packet[0], packet[5]}
			currentTime := time.Now().Format("2024-11-14 15:04:05")

			
			agentMutex.Lock()
			main.RemoveAgent(agent, agents)
			agentMutex.Unlock()

			// Escreve no log e remove o buffer do maps de logs
			log := fmt.Sprintf("Agent %s disconnected...", agent.AgentID)
			lm.AddLog(agent.AgentID, log, currentTime)
			lm.RemoveClientBuffer(agent.AgentID)
	}
}


