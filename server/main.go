package main

import (
	"ccproj/server/db"
	sAlertFlow "ccproj/server/serverAlertFlow"
	sNetTask "ccproj/server/serverNetTask"
	"ccproj/server/view"
	"ccproj/server/types"
	"ccproj/utils"
	"fmt"
	"os"
	"time"
)

func main() {

	// Check if the user provided a configuration file
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <config.json>")
		return
	}

	serverIP, err := getLocalIP()
    if err != nil {
        fmt.Printf("Erro ao obter o IP local: %v\n", err)
        return
    }

	udpServerAddr := fmt.Sprintf("%s:9090", serverIP)
    tcpServerAddr := fmt.Sprintf("%s:8080", serverIP)

	configFile := os.Args[1]

	// Parsing the configuration file
	if !gUtils.IsJSONFile(configFile) {
		fmt.Printf("Error: Configuration file must be a .json\n")
		return
	}

	tasks, err := ParseTasks(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	logs := db.NewLogManager()
	agents := make(map[string]types.Agent)
	
	sendChannel := make(chan []string)
	defer close(sendChannel)

	go logs.PersistLogs()

	// Listener UDP, para resgistos, metricas, confirmacoes
	go sNetTask.HandleUDP(udpServerAddr, agents, logs, sendChannel)

	// Listener TCP para alertas
	go sAlertFlow.HandleTCP(tcpServerAddr, agents, logs)

	// timer para enviar tarefas

	time.Sleep(30 * time.Second)
	fmt.Println("Sending tasks to agents...\n")
	SendTask(agents, tasks, sendChannel)
	fmt.Printf("Tasks sent to agents...\n")

	view.StartGUI()

}


