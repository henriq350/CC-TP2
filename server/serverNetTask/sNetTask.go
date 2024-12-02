package sNetTask

import (
	"ccproj/server/db"
	uh "ccproj/udp_handler"
	"fmt"
	"main"
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
		handleUDPMessage(packet, agents, lm)
	}
	// Envia mensagem (com a task e ACKs) para os agents

}

func handleUDPMessage(packet []string, agents map[string]uh.AgentRegistration, lm *db.LogManager) {

	switch packet[0] {
		case "Register":
			// Envia ACK
			
			// Cria um agente a partir do pacote
			agent := getAgent(packet)
			currentTime := time.Now().Format("2024-11-14 15:04:05")

			// Adiciona a lista de agentes
			agentMutex.Lock()
			main.AddAgent(agent, agents)
			agentMutex.Unlock()

			// Adiciona Log
			log := fmt.Sprintf("[%s] Agent %s registered", currentTime, agent.AgentID)
			lm.AddLog(agent.AgentID,log)
			
		case "Report":
			
			// Remove Tipo do array
			metrics := packet[1:]
			agent := getAgent(packet)

			formatedString, currentTime := db.FormatString(metrics)

			// TODO com a estrutura q vier do udp_handler, pegar o nome do cliente
			//clientName := 

			filename := fmt.Sprintf("%s", &currentTime)

			db.StringToFile(clientName,filename,formatedString)

			// Adiciona nos Logs
			log := fmt.Sprintf("[%s] Package received", currentTime)
			lm.AddLog(agent.AgentID ,log)

		case "Terminate":
			// Envia ACK
			agent := getAgent(packet)
			currentTime := time.Now().Format("2024-11-14 15:04:05")

			// Remove da lista de agents
			agentMutex.Lock()
			main.RemoveAgent(agent, agents)
			agentMutex.Unlock()

			// Escreve no log e remove o buffer do maps de logs
			log := fmt.Sprintf("[%s] Agent %s disconnected...", currentTime, agent.AgentID)
			lm.AddLog(agent.AgentID,log)
			lm.RemoveClientBuffer(agent.AgentID)
	}
}

// TODO
func SendTaskToAgent(agent uh.AgentRegistration, task uh.Task) {
	// Envia a task para o agente
}


// TODO refazer esta funcao, esperar pelo tipo de Agent para importar para este package
func getAgent(packet uh.Packet) uh.AgentRegistration {
	
	agent := uh.AgentRegistration{packet.Data.AgentID, packet.Data.IPv4}
	return agent
}