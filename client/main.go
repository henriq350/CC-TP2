package main

import (
	cNetTask "ccproj/client/clientNetTask"
	"ccproj/client/tasks"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)



func main() {

	//Get server IP
	if len(os.Args) < 4 {
        fmt.Println("Insert server IP with the following ports UDP and TCP\nExample: go run main.go 10.0.0.1 9090 8080\n")
        os.Exit(1)
    }

    serverIP := os.Args[1]
    validateServerIP(serverIP)

    udpPort := validatePort(os.Args[2])
    tcpPort := validatePort(os.Args[3])

    udpServerAddr := fmt.Sprintf("%s:%d", serverIP, udpPort)
    tcpServerAddr := fmt.Sprintf("%s:%d", serverIP, tcpPort)

	clientip, err := getLocalIP()
    if err != nil {
        fmt.Printf("Erro ao obter o IP local: %v\n", err)
        return
    }

	clientIP := fmt.Sprintf("%s:8080", clientip)

	clientID := getClientID()
	fmt.Print("%s Running... \n", clientID)

	Tasks := make(map[string]tasks.Task)
	taskChannel := make(chan []string)

	go cNetTask.HandleUDP(udpServerAddr, taskChannel)

	register := []string{clientID, "","Register","","",clientIP,""}
	taskChannel <- register

	go func() {
		for task := range taskChannel {
			tasks.AddTask(task, Tasks)
			taskID := task[1]
			go tasks.ProcessTask(taskID, Tasks[taskID], clientID, tcpServerAddr, taskChannel)
		}
	}()
	
	sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    
    sig := <-sigChan
    fmt.Printf("Signal recevied %s. Sending terminate packet...\n", sig)

    terminate := []string{clientID, "","Terminate","","",clientIP,""}
	taskChannel <- terminate

    close(taskChannel)
}