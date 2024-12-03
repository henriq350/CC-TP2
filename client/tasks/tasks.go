package tasks


import (
	"ccproj/client/metrics"
	th "ccproj/tcp_handler"
	"fmt"
	"strconv"
	"time"
)

// NOTE> Podera ser necessario alterar os valores do tipo de task a se enviar, uma vez que pode ser diferente no udpHandler

type Task struct {
	nTask int
	MetricType []string
	Frequency int
	Threshold float32
	IpAddress string
    Duration int
    PacketCount int
}

func AddTask(task []string, tasks map[string]Task) {
	taskID := task[1]
	tasks[taskID] = ParseTask(task[2:]) // Removes ClientID and TaskID
}

//  "name"  "frequencia" "threshold" "dest_ip"  ”duration” ”packet_count”
//  task[0]   task[1]       task[2]    task[3]    task[4]     task[5]     

func ParseTask(task []string) Task {
    nTask := 1 // Change if we get task with multiple metrics
    metricType := []string{task[0]}
    frequency, _ := strconv.Atoi(task[1])
    threshold, _ := strconv.ParseFloat(task[2], 32)
    duration, _ := strconv.Atoi(task[4])
    packetCount, _ := strconv.Atoi(task[5])

    return Task{
        nTask,
        metricType,
        frequency,
        float32(threshold),
        task[3],
        duration,
        packetCount,
    }
}


func ProcessTask(taskID string, task Task, agentID string, serverIP string, udpChannel chan <- []string) {

	fmt.Printf("Task %s started...\n", taskID)
    metricsChannel := make(chan []string)
    defer close(metricsChannel)

	// routine para enviar mensagens para o servidor
	
	for _, metric := range task.MetricType {
		switch metric {
			case "CPU":
				go monitorCPU(task.Frequency, task.Threshold, metricsChannel)
			case "RAM":
				go monitorRAM(task.Frequency, task.Threshold, metricsChannel)
			case "Bandwidth":
				go monitorBandwidth(task.Frequency, task.IpAddress, task.Duration, metricsChannel)
			case "Latency":
				go monitorLatency(task.Frequency, task.IpAddress, task.PacketCount, metricsChannel)
			case "PacketLoss":
				go monitorPacketLoss(task.Frequency, task.Threshold, task.IpAddress, task.PacketCount, metricsChannel)
			case "Jitter":
				go monitorJitter(task.Frequency, task.Threshold, task.IpAddress, task.Duration, metricsChannel)
        }
	}

    // Send Alert message
	for message := range metricsChannel {
        metricType := message[0]
        value, _ := strconv.ParseFloat(message[1], 32)
        threshold, _ := strconv.ParseFloat(message[2], 32)

        if value > threshold {
            var alertMetric th.AlertMetric
            switch metricType {
                case "CPU":
                    alertMetric = th.CPUUsage
                case "RAM":
                    alertMetric = th.RAMUsage
                case "PacketLoss":
                    alertMetric = th.PacketLoss
                case "Jitter":
                    alertMetric = th.Jitter
            }

            alert := th.AlertMessage{
                AgentID:    agentID,
                TaskID:     taskID,
                AlertMetric: alertMetric,
                Threshold:  float32(threshold),
                Value:      float32(value),
            }
            th.SendAlert(serverIP, alert)
			fmt.Println("Alert sent to server...")

            udpMessage := []string{agentID, taskID, metricType, message[1], message[2]}
            udpChannel <- udpMessage
        }
    }
}


func monitorCPU(frequency int, threshold float32, send chan <- []string) {
	ticker := time.NewTicker(time.Duration(frequency) * time.Second) 
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cpuUsage, err := metrics.GetCPUUsage()
            if err != nil {
                fmt.Println("Error getting CPU usage:", err)
                continue
            }

			message := []string{"CPU", strconv.FormatFloat(cpuUsage, 'f', 2, 64), strconv.FormatFloat(float64(threshold), 'f', 2, 64)}
            send <- message
        }
    }
}

func monitorRAM(frequency int, threshold float32, send chan <- []string) {
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            ramUsage, err := metrics.GetRAMUsage()
            if err != nil {
                fmt.Println("Error getting RAM usage:", err)
                continue
            }


            message := []string{"RAM", strconv.FormatFloat(ramUsage, 'f', 2, 64), strconv.FormatFloat(float64(threshold), 'f', 2, 64)}
            send <- message
        }
    }
}

func monitorBandwidth(frequency int, ipDest string, duration int, send chan <- []string) {
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            bandwidth, _, err := metrics.iperfMetrics(ipDest, duration)
            if err != nil {
                fmt.Println("Error getting bandwidth:", err)
                continue
            }

            message := []string{"Bandwidth", strconv.FormatFloat(bandwidth, 'f', 2, 64), ""}
            send <- message
            
        }
    }
}

func monitorLatency(frequency int, ipDest string, count int, send chan <- []string) {
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            latency, _, err := metrics.pingMetrics(ipDest, count)
            if err != nil {
                fmt.Println("Error getting latency:", err)
                continue
            }

            message := []string{"Latency", strconv.FormatFloat(latency, 'f', 2, 64), ""}
            send <- message
            
        }
    }
}

func monitorPacketLoss(frequency int, threshold float32, ipDest string, count int,send chan <- []string) {
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            _, packetLoss, err := metrics.pingMetrics(ipDest, count)
            if err != nil {
                fmt.Println("Error getting packet loss:", err)
                continue
            }

            
            message := []string{"PacketLoss", strconv.FormatFloat(packetLoss, 'f', 2, 64), strconv.FormatFloat(float64(threshold), 'f', 2, 64)}
            send <- message
            
        }
    }
}

func monitorJitter(frequency int, threshold float32, ipDest string, duration int, send chan <- []string) {
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            _, jitter, err := metrics.iperfMetrics(ipDest, duration)
            if err != nil {
                fmt.Println("Error getting jitter:", err)
                continue
            }

            
            message := []string{"Jitter", strconv.FormatFloat(jitter, 'f', 2, 64), strconv.FormatFloat(float64(threshold), 'f', 2, 64)}
            send <- message
        }
    }
}
