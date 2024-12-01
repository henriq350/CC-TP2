package main

import (
	// "ccproj/utils"
	// "fmt"
	// "os"
	//uh "ccproj/udp_handler"
	"ccproj/server/view"
	//th "ccproj/tcp_handler"
	//"fmt"
	
)

func main() {

	// address := "localhost:8080"
	// udpAddress := "localhost:9090"
	// // Check if the user provided a configuration file
	// if len(os.Args) < 2 {
	// 	fmt.Println("Usage: go run main.go <config.json>")
	// 	return
	// }

	// configFile := os.Args[1]

	// // Parsing the configuration file
	// if !gUtils.IsJSONFile(configFile) {
	// 	fmt.Printf("Error: Configuration file must be a .json\n")
	// 	return
	// }

	// task, err := ParseTasks(configFile)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// agents := make(map[string]uh.AgentRegistration)

	// PrintTasks(task)

	// // Listener UDP, para resgistos, metricas, confirmacoes
	// go sNetTask.HandleUDP(udpAddress, agents)
	// 	// ficar a espera de registos, enviar as tasks para os agents e esperar metricas

	// // Listener TCP para alertas
	// go tcp_handler.ListenTCP(address)

	view.StartGUI()



}


