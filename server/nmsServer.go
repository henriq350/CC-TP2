package main

import (
	"ccproj/server/config"
	"ccproj/server/utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

type Agent struct {
	AgentID string
	AgentIP string
}


func ParseTasks(filename string) ([]config.Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if ok, errMsg := sutils.ValidateJSON(data); !ok {
		return nil, fmt.Errorf(errMsg)
	}

	// Reseta o ponteiro do ficheiro para o inicio
	file.Seek(0, io.SeekStart)

	var tasks []config.Task
	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}


// Add agent
func AddAgent(agent Agent, agents map[string]Agent) {
	agents[agent.AgentID] = agent
}

// Remove agent
func RemoveAgent(agent Agent, agents map[string]Agent) {
	delete(agents, agent.AgentID)
}

// TODO
func SendTask(agents map[string]Agent , tasks []config.Task, sendChannel chan <- []string) {

    agentList := make([]Agent, 0, len(agents))
    for _, agent := range agents {
        agentList = append(agentList, agent)
    }

    if len(agentList) == 0 {
        fmt.Println("Empty agent list")
        return
    }

    agentIndex := 0

    for _, task := range tasks {
        for i, device := range task.Devices {
			
			// Distributes the tasks to the agents equally
            agent := agentList[agentIndex]
            agentIndex = (agentIndex + 1) % len(agentList)


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
			if device.LinkMetrics.Bandwidth {
                message := []string{
                    task.ID, "Bandwidth", fmt.Sprintf("%d", task.Frequency), "",
                    agent.AgentIP, "0", task.Devices[i].LinkMetrics.Bandwidth.Duration, "",
                }
                sendChannel <- message
            }
			if device.LinkMetrics.Latency {
                message := []string{
                    task.ID, "Latency", fmt.Sprintf("%d", task.Frequency), "",
                    agent.AgentIP, task.Devices[i].LinkMetrics.Latency.Destination, "", task.Devices[i].LinkMetrics.Latency.Count,
                }
                sendChannel <- message
            }
			if device.LinkMetrics.PacketLoss {
                message := []string{
                    task.ID, "PacketLoss", fmt.Sprintf("%d", task.Frequency), fmt.Sprintf("%.2f", task.Devices[i].LinkMetrics.AlertFlowConditions.PacketLoss),
                    agent.AgentIP, task.Devices[i].LinkMetrics.PacketLoss.Destination , "", task.Devices[i].LinkMetrics.PacketLoss.Count,
                }
                sendChannel <- message
            }
			if device.LinkMetrics.Jitter {
                message := []string{
                    task.ID, "Jitter", fmt.Sprintf("%d", task.Frequency), fmt.Sprintf("%.2f", task.Devices[i].LinkMetrics.AlertFlowConditions.Jitter),
                    agent.AgentIP,  , task.Devices[i].LinkMetrics.Jitter.Duration, "",
                }
                sendChannel <- message
            }
           
        }
    }
}