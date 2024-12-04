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
		fmt.Println("Usage: go run main.go nmsServer.go <config.json>")
		return
	}

	// serverIP, err := getLocalIP()
    // if err != nil {
    //     fmt.Printf("Erro ao obter o IP local: %v\n", err)
    //     return
    // }

	serverIP := "127.0.0.1"

	udpServerAddr := fmt.Sprintf("%s:9090", serverIP)
    tcpServerAddr := fmt.Sprintf("%s:8080", serverIP)

	configFile := os.Args[1]

	// Parsing the configuration file
	if !gUtils.IsJSONFile(configFile) {
		fmt.Printf("Error: Configuration file must be a .json\n")
		return
	}

	fmt.Println("Parsing configuration file...")
	tasks, err := ParseTasks(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	logs := db.NewLogManager()
	fmt.Println("Agent maps created...")
	agents := make(map[string]types.Agent)
	
	sendChannel := make(chan []string)
	defer close(sendChannel)

	fmt.Println("Logs created...")
	go logs.PersistLogs()

	// Listener UDP, para resgistos, metricas, confirmacoes
	fmt.Println("Starting UDP listener...")
	go sNetTask.HandleUDP(udpServerAddr, agents, logs, sendChannel)

	// Listener TCP para alertas
	fmt.Println("Starting TCP listener...")
	go sAlertFlow.HandleTCP(tcpServerAddr, agents, logs)

	// timer para enviar tarefas

	time.Sleep(7* time.Second)

	fmt.Println("Sending tasks to agents...\n")
	SendTask(agents, tasks, sendChannel)
	fmt.Printf("Tasks sent to agents...\n")

	fmt.Println("Starting GUI...")
	time.Sleep(2 * time.Second)
	view.StartGUI(agents)

}


