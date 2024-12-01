package sAlertFlow

import (
	"ccproj/server/db"
	th "ccproj/tcp_handler"
	uh "ccproj/udp_handler"
	"fmt"
	"time"
)

func HandleTCP(tcpAddr string,  agents map[string]uh.AgentRegistration, lm *db.LogManager) {
	
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

	formattedLog := formatAlertMessage(alert, currentTime)
	
	//Add to respetive buffer
	lm.AddLog(alert.AgentID, formattedLog)
}

func formatAlertMessage(alert th.AlertMessage, time string) string {
	return "[" + time + "]" + "[Alert] " + alert.AlertMetric.String() + ": " + fmt.Sprintf("%.1f", alert.Value) + " - Threshold: " + fmt.Sprintf("%.1f", alert.Threshold)
}