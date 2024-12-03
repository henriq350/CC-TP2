package main

import (
	"ccproj/server/config"
	sutils "ccproj/server/utils"
	"encoding/json"
    "ccproj/server/types"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)


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


// TODO
func SendTask(agents map[string]types.Agent , tasks []config.Task, sendChannel chan <- []string) {

    agentList := make([]types.Agent, 0, len(agents))
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
}


func validateServerIP(ip string) {
    if net.ParseIP(ip) == nil {
        fmt.Println("Invalid IP:", ip)
        os.Exit(1)
    }
}


func validatePort(portStr string) int {

    port, err := strconv.Atoi(portStr)

    if err != nil || port < 1024 || port > 65535 {
        fmt.Println("Invalide port:", portStr)
        os.Exit(1)
    }

    return port
}

func getLocalIP() (string, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "", err
    }

    for _, addr := range addrs {
        if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
            if ipNet.IP.To4() != nil {
                return ipNet.IP.String(), nil
            }
        }
    }

    return "", fmt.Errorf("Failed to get local IP")
}