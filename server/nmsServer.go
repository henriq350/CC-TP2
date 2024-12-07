package main

import (
	"ccproj/server/config"
	sutils "ccproj/server/utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
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

	file.Seek(0, io.SeekStart)

	var tasks []config.Task
	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
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

    return "", fmt.Errorf("failed to get local IP")
}