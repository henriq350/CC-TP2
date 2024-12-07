package sNetTask

import (
	"ccproj/server/db"
	uh "ccproj/udp_handler"
	"ccproj/server/types"
	"ccproj/server/config"
	"fmt"
	"sync"
	"time"
	"net"
	"os"
)

var (
    agentMutex    sync.Mutex
    pendingTasks  []config.Task
    originalTasks []config.Task
)

func HandleUDP(udpAddr string, agents map[string]types.Agent, lm *db.LogManager, sendChannel, receiveChannel chan []string, tasks []config.Task) {

	originalTasks = tasks
	pendingTasks = make([]config.Task, len(tasks))
	copy(pendingTasks, originalTasks)

	udp_address,error := net.ResolveUDPAddr("udp",udpAddr)
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}

	connection_, _ := net.ListenUDP("udp", udp_address)
	term := make(chan bool)
	go uh.ListenUdp("","",connection_ ,receiveChannel,term)
	
	go uh.ListenServer(sendChannel,connection_)

	//Receber mensagem e decidir o q fazer com ela
	for packet := range receiveChannel {
		fmt.Printf("Received packet: %v\n", packet)
		go handleUDPMessage(packet, agents, lm, sendChannel)
	}

}

// packet - "client_id" ,”task_id”  ,"tipo"    ,"metrica" ,"valor"  ,”client_ip”,"dest_ip"
// packet -  packet[0] , packet[1] ,packet[2] ,packet[3] packet[4] ,packet[5]   ,packet[6]

func handleUDPMessage(packet []string, agents map[string]types.Agent, lm *db.LogManager, sendChannel chan []string) {

	switch packet[2] {
		case "Register":

			agent := types.Agent{AgentID: packet[0], AgentIP: packet[5]}
			currentTime := time.Now().Format("15:04:05")

			
			agentMutex.Lock()
			types.AddAgent(agent, agents)
			agentMutex.Unlock()

			
			log := fmt.Sprintf("Agent %s registered", agent.AgentID)
			lm.AddLog(agent.AgentID, log, currentTime, true)
			
			assignTaskToAgent(agent, sendChannel)
			fmt.Println("Agent assigned task")
			
		case "Report":
			
			agentID := packet[0]

			aux := packet[1]
			metrics := packet[3:]
			metrics = append(metrics, aux)

			formatedString, currentTime := db.FormatString(metrics) 

			filename := currentTime

			db.StringToFile(agentID , filename, formatedString)

			
			log := "Packet received" + db.FormatStringLog(metrics)
			currentTime = time.Now().Format("15:04:05")
			lm.AddLog(agentID ,log, currentTime, false)

		case "Terminate":
			
			agent := types.Agent{AgentID: packet[0], AgentIP:  packet[5]}
			currentTime := time.Now().Format("15:04:05")

			
			agentMutex.Lock()
			types.RemoveAgent(agent, agents)
			agentMutex.Unlock()

			
			log := fmt.Sprintf("Agent %s disconnected...", agent.AgentID)
			lm.AddLog(agent.AgentID, log, currentTime, false)
			lm.RemoveClientBuffer(agent.AgentID)
	}
}



func assignTaskToAgent(agent types.Agent, sendChannel chan <- []string) {

    if len(pendingTasks) == 0 {
        fmt.Println("No pending tasks to assign.")
        pendingTasks = make([]config.Task, len(originalTasks))
        copy(pendingTasks, originalTasks)
    }

    task := pendingTasks[0]
    pendingTasks = pendingTasks[1:]

	fmt.Printf("Assigning task %s to agent %s\n", task.ID, agent.AgentID)

    for i, device := range task.Devices {
        if device.DeviceMetrics.CPUUsage {
            message := []string{
                task.ID, "CPU", fmt.Sprintf("%d", task.Frequency), fmt.Sprintf("%.2f", task.Devices[i].LinkMetrics.AlertFlowConditions.CPUUsage),
                agent.AgentIP, "0", "0", "0",
            }
            sendChannel <- message
        }
        if device.DeviceMetrics.RAMUsage {
            message := []string{
                task.ID, "RAM", fmt.Sprintf("%d", task.Frequency), fmt.Sprintf("%.2f", task.Devices[i].LinkMetrics.AlertFlowConditions.RAMUsage),
                agent.AgentIP, "0", "0", "0",
            }
            sendChannel <- message
        }
        for _, interfaceName := range device.DeviceMetrics.InterfaceStats{
            if device.LinkMetrics.Bandwidth.Duration > 0 {
                message := []string{
                    task.ID, "Bandwidth", fmt.Sprintf("%d", task.Frequency), "",
                    agent.AgentIP, interfaceName , fmt.Sprintf("%d",task.Devices[i].LinkMetrics.Bandwidth.Duration), "",
                }
                sendChannel <- message
            }

            if device.LinkMetrics.Jitter.Duration > 0 {
                message := []string{
                    task.ID, "Jitter", fmt.Sprintf("%d", task.Frequency), fmt.Sprintf("%.2f", task.Devices[i].LinkMetrics.AlertFlowConditions.Jitter),
                    agent.AgentIP, interfaceName, fmt.Sprintf("%d", task.Devices[i].LinkMetrics.Jitter.Duration), "",
                }
                sendChannel <- message
            }
        }
        if device.LinkMetrics.Latency.Destination != "" && device.LinkMetrics.Latency.Count > 0 {
            message := []string{
                task.ID, "Latency", fmt.Sprintf("%d", task.Frequency), "",
                agent.AgentIP, task.Devices[i].LinkMetrics.Latency.Destination, "", fmt.Sprintf("%d", task.Devices[i].LinkMetrics.Latency.Count),
            }
            sendChannel <- message
        }
        if device.LinkMetrics.PacketLoss.Destination != "" && device.LinkMetrics.PacketLoss.Count > 0 {
            message := []string{
                task.ID, "PacketLoss", fmt.Sprintf("%d", task.Frequency), fmt.Sprintf("%.2f", task.Devices[i].LinkMetrics.AlertFlowConditions.PacketLoss),
                agent.AgentIP, task.Devices[i].LinkMetrics.PacketLoss.Destination , "", fmt.Sprintf("%d", task.Devices[i].LinkMetrics.PacketLoss.Count),
            }
            sendChannel <- message
        }
    }
}