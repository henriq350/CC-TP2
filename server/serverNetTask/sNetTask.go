package sNetTask

import (
	uh "ccproj/udp_handler"
	"main"
)



func handleUDP(udpAddr string, agents map[string]uh.AgentRegistration) {

	packetChannel := make(chan []uh.Packet)
	// Listener de UDP
	go udp_handler.ListenUDP(udpAddr, packetChannel)

	//Receber mensagem e decidir o q fazer com ela
	// for packet := range packetChannel {
	// 	switch packet.Type {
	// 		case RegisterPacket:
	// 			// Envia ACK
				
	// 			// Cria um agente a partir do pacote
	// 			agent := getAgent(packet)

	// 			// Adiciona a lista de agentes
	// 			main.AddAgent(agent, agents)
				
	// 		case ReportPacket:
	// 			// Recebe metricas

			

	// 			// Parse das metricas

	// 		case TerminatePacket:
	// 			// Envia ACK

	// 			agent := getAgent(packet)
	// 			// Remove da lista de agents
	// 			main.RemoveAgent(agent, agents)
	// 	}

	// }
	// Envia mensagem (com a task e ACKs) para os agents

}

func getAgent(packet uh.Packet) uh.AgentRegistration {
	
	agent := uh.AgentRegistration{packet.Data.AgentID, packet.Data.IPv4}
	return agent
}