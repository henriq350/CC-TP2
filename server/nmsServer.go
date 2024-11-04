package main

import (
	"ccproj/server/config"
	"ccproj/server/utils"
	"encoding/json"
	"fmt"
	"io"
	"os"
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

// Funcao de teste
func PrintTasks(tasks []config.Task) {
	for _, task := range tasks {
		fmt.Printf("Task ID: %d\n", task.ID)
		fmt.Printf("Task Frequency: %d\n", task.Frequency)
		for _, device := range task.Devices {
			fmt.Printf("\tDevice ID: %d\n", device.ID)
			fmt.Printf("\tDevice CPU Usage: %t\n", device.DeviceMetrics.CPUUsage)
			fmt.Printf("\tDevice RAM Usage: %t\n", device.DeviceMetrics.RAMUsage)
			fmt.Printf("\tDevice Interface Stats: %v\n", device.DeviceMetrics.InterfaceStats)
			fmt.Printf("\tLink Bandwidth Tool: %s\n", device.LinkMetrics.Bandwidth.Tool)
			fmt.Printf("\tLink Bandwidth Client: %t\n", device.LinkMetrics.Bandwidth.Client)
			fmt.Printf("\tLink Bandwidth Server Address: %s\n", device.LinkMetrics.Bandwidth.ServerAddr)
			fmt.Printf("\tLink Bandwidth Duration: %d\n", device.LinkMetrics.Bandwidth.Duration)
			fmt.Printf("\tLink Bandwidth Transport: %s\n", device.LinkMetrics.Bandwidth.Transport)
			fmt.Printf("\tLink Bandwidth Frequency: %d\n", device.LinkMetrics.Bandwidth.Frequency)
			fmt.Printf("\tLink Jitter Tool: %s\n", device.LinkMetrics.Jitter.Tool)
			fmt.Printf("\tLink Jitter Client: %t\n", device.LinkMetrics.Jitter.Client)
			fmt.Printf("\tLink Jitter Server Address: %s\n", device.LinkMetrics.Jitter.ServerAddr)
			fmt.Printf("\tLink Jitter Duration: %d\n", device.LinkMetrics.Jitter.Duration)
			fmt.Printf("\tLink Jitter Transport: %s\n", device.LinkMetrics.Jitter.Transport)
			fmt.Printf("\tLink Jitter Frequency: %d\n", device.LinkMetrics.Jitter.Frequency)
			fmt.Printf("\tLink Packet Loss Tool: %s\n", device.LinkMetrics.PacketLoss.Tool)
			fmt.Printf("\tLink Packet Loss Client: %t\n", device.LinkMetrics.PacketLoss.Client)
			fmt.Printf("\tLink Packet Loss Server Address: %s\n", device.LinkMetrics.PacketLoss.ServerAddr)
			fmt.Printf("\tLink Packet Loss Duration: %d\n", device.LinkMetrics.PacketLoss.Duration)
			fmt.Printf("\tLink Packet Loss Transport: %s\n", device.LinkMetrics.PacketLoss.Transport)
			fmt.Printf("\tLink Packet Loss Frequency: %d\n", device.LinkMetrics.PacketLoss.Frequency)
			fmt.Printf("\tLink Latency Destination: %s\n", device.LinkMetrics.Latency.Destination)
			fmt.Printf("\tLink Latency Count: %d\n", device.LinkMetrics.Latency.Count)
			fmt.Printf("\tLink Latency Frequency: %d\n", device.LinkMetrics.Latency.Frequency)
			fmt.Printf("\tLink Alert Flow Conditions CPU Usage: %f\n", device.LinkMetrics.AlertFlowConditions.CPUUsage)
			fmt.Printf("\tLink Alert Flow Conditions RAM Usage: %f\n", device.LinkMetrics.AlertFlowConditions.RAMUsage)
			fmt.Printf("\tLink Alert Flow Conditions Interface Stats: %d\n", device.LinkMetrics.AlertFlowConditions.InterfaceStats)
			fmt.Printf("\tLink Alert Flow Conditions Packet Loss: %f\n", device.LinkMetrics.AlertFlowConditions.PacketLoss)
			fmt.Printf("\tLink Alert Flow Conditions Jitter: %d\n", device.LinkMetrics.AlertFlowConditions.Jitter)
		}
	}
}

//func show_metrics() {}
