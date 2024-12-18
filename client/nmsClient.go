package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

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

func getClientID() string {
    clientID, err := os.Hostname()
    if err != nil {
        fmt.Println("Error getting clientID:", err)
        os.Exit(1)
    }
    return clientID
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