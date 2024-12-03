package sAlertFlow

import (
	"ccproj/server/db"
	th "ccproj/tcp_handler"
	"ccproj/server"
	"fmt"
	"time"
)

func HandleTCP(tcpAddr string,  agents map[string]main.Agent, lm *db.LogManager) {
	
	receiveChannel := make(chan th.AlertMessage)

	// listener
	go th.ListenTcp(tcpAddr, receiveChannel)

	//Receber mensagem e decidir o q fazer com ela
	for packet := range receiveChannel {
		go handleAlertMessage(packet, lm)
	}
}


func handleAlertMessage(alert th.AlertMessage, lm *db.LogManager) {
	currentTime := time.Now().Format("2024-11-14 15:04:05")

	formattedLog := formatAlertMessage(alert)
	
	//Add to respetive buffer
	lm.AddLog(alert.AgentID, formattedLog, currentTime)
}

func formatAlertMessage(alert th.AlertMessage) string {
	return "[Alert][" + alert.TaskID + "] " + alert.AlertMetric.String() + ": " + fmt.Sprintf("%.1f", alert.Value) + " - Threshold: " + fmt.Sprintf("%.1f", alert.Threshold)
}