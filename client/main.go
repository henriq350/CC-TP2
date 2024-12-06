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
	if len(os.Args) < 2 {
        fmt.Println("Insert server IP\nExample: go run main.go 10.0.0.1")
        os.Exit(1)
    }

    serverIP := os.Args[1]
    validateServerIP(serverIP)

    udpServerAddr := fmt.Sprintf("%s:9090", serverIP)
    tcpServerAddr := fmt.Sprintf("%s:8080", serverIP)

	clientip, err := getLocalIP()
    if err != nil {
        fmt.Printf("Erro ao obter o IP local: %v\n", err)
        return
    }

	clientIP := fmt.Sprintf("%s:9090", clientip)

	clientID := getClientID()
	fmt.Printf("%s Running... \n", clientID)

	Tasks := make(map[string]tasks.Task)
	taskChannel := make(chan []string)
    sendChannel := make(chan []string)

    defer close(taskChannel)
    defer close(sendChannel)
	//tasksReady := make(chan struct{})
	
	go cNetTask.HandleUDP(clientIP, udpServerAddr ,taskChannel, sendChannel)

	register := []string{clientID, "","Register","","",clientIP,udpServerAddr}
	sendChannel <- register

	go func() {
		for task := range taskChannel {
            fmt.Println("Received task:", task)
            if len(task) < 2 {
                fmt.Println("Erro: Tamanho do slice task é menor que o esperado")
                continue
            }
            if task[2] == "Task" { // Verifica se o pacote contém tarefas
                tasks.AddTask(task, Tasks)
                taskID := task[1]
                go tasks.ProcessTask(taskID, Tasks[taskID], clientID, tcpServerAddr, sendChannel)
                //tasksReady <- struct{}{} // Sinaliza que as tarefas foram recebidas
            }
        }
	}()
	
	sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    
    sig := <-sigChan
    fmt.Printf("Signal recevied %s. Sending terminate packet...\n", sig)

    terminate := []string{clientID, "","Terminate","","",clientIP,""}
	sendChannel <- terminate

    close(taskChannel)
}