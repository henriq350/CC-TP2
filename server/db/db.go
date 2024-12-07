package db

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type LogManager struct {
	ClientBuffers map[string][]string
	GeneralBuffer []string
	Mutex         sync.Mutex
}

func StringToFile(clientName, fileName, data string) {
	err := SaveFile(clientName, fileName, []byte(data))
	if err != nil {
		fmt.Println("Error saving file:", err)
	} else {
		fmt.Println("File saved successfully")
	}
}

func SaveFile(clientName, fileName string, data []byte) error {
	clientDir := filepath.Join("../client_metrics", clientName)
	err := CreateFolder(clientName)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(fileName, ".txt") {
		fileName += ".txt"
	}

	filePath := filepath.Join(clientDir, fileName)
	err = os.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}
	return nil
}

// Create folder for each client on folder client_metrics
func CreateFolder(client string) error {
	clientDir := filepath.Join("../client_metrics", client)
	err := os.MkdirAll(clientDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("Error creatring CLient's folder: %v", err)
	}

	CreateLog(clientDir)
	return nil
}

func CreateLog(dir string) error {

	logFilePath := filepath.Join(dir, "log.txt")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		return fmt.Errorf("Error creating log.txt: %v", err)
	}
	defer logFile.Close()

	return nil
}

// packet -  "metrica" ,  "valor"  ,”client_ip” ,"dest_ip" ,”task_id”
// packet -  packet[0] , packet[1]  ,packet[2]  ,packet[3]  packet[4]
func FormatString(data []string) (string, string) {

	var fmtData strings.Builder

	currentTime := time.Now().Format("15:04:05")
	fmtData.WriteString(fmt.Sprintf("Received at: %s\n\n", currentTime))

    metric := data[0]
    value := data[1]
    sourceIP := data[2]
    destIP := data[3]
    taskID := data[4]
    fmtData.WriteString(fmt.Sprintf("=======[TasK %s]=======\n%s: %s\nSource IP: %s\nDestination IP: %s\n", taskID, metric, value, sourceIP, destIP))
	

	return fmtData.String(), currentTime
}

func FormatStringLog(data []string) string {

    metric := data[0]
    value := data[1]
    task := data[4]

    return fmt.Sprintf(" from %s >>> %s: %s", task, metric, value)

}

/////////////////////////////////////////// LOG LOGIC

/////////////////////////////  WRITE LOGS

// init new LogManager
func NewLogManager() *LogManager {
	return &LogManager{
		ClientBuffers: make(map[string][]string),
		GeneralBuffer: []string{},
		Mutex:         sync.Mutex{},
	}
}

// add log to client buffer and general buffer
func (lm *LogManager) AddLog(clientID, log string, time string, isRegister bool) {
	lm.Mutex.Lock()
	defer lm.Mutex.Unlock()

	if _, exists := lm.ClientBuffers[clientID]; !exists {
		lm.ClientBuffers[clientID] = []string{}
	}

	// Add log to client
    if !isRegister{
        clientLog := fmt.Sprintf("[%s] %s", time, log)
        lm.ClientBuffers[clientID] = append(lm.ClientBuffers[clientID], clientLog)
    }

	// formate log for general log
	generalLog := fmt.Sprintf("[%s] [%s] %s", time, clientID, log)
	lm.GeneralBuffer = append(lm.GeneralBuffer, generalLog)

}

// remove client buffer after client terminate
func (lm *LogManager) RemoveClientBuffer(clientID string) {
	lm.Mutex.Lock()
	defer lm.Mutex.Unlock()

	// saves information inside buffer before removing
	if logs, exists := lm.ClientBuffers[clientID]; exists && len(logs) > 0 {
		filePath := fmt.Sprintf("../client_metrics/%s/log.txt", clientID)
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening log file for client %s: %v\n", clientID, err)
			return
		}
		defer file.Close()

		for _, log := range logs {
			file.WriteString(log + "\n")
		}
	}

	delete(lm.ClientBuffers, clientID)
}

// saves information from buffer to file every X seconds
func (lm *LogManager) PersistLogs() {

	for {
		lm.Mutex.Lock()

		for clientID, logs := range lm.ClientBuffers {

			if len(logs) > 0 {

				filePath := fmt.Sprintf("../client_metrics/%s/log.txt", clientID)

				if _, err := os.Stat(filePath); os.IsNotExist(err) {
                    err := CreateLog(filepath.Join("../client_metrics", clientID))
                    if err != nil {
                        fmt.Printf("Error creating file for client %s: %v\n", clientID, err)
                        continue
                    }
                }

				file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                if err != nil {
                    fmt.Printf("Error opening file for client %s: %v\n", clientID, err)
                    continue
                }

				for _, log := range logs {
                    _, err := file.WriteString(log + "\n")
                    if err != nil {
                        fmt.Printf("Error writing log for client %s: %v\n", clientID, err)
                    }
                }
				lm.ClientBuffers[clientID] = nil
				file.Close()

            } 
	    }

		if len(lm.GeneralBuffer) > 0 {

			filePath := "../client_metrics/log.txt"
			file, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			for _, log := range lm.GeneralBuffer {
				file.WriteString(log + "\n")
			}

			lm.GeneralBuffer = nil 
            fmt.Println(lm.GeneralBuffer)
			file.Close()
		}

		lm.Mutex.Unlock()
		time.Sleep(5 * time.Second)
	}
}

/////////////////////////////  READ LOGS

//////////// Client Logs

func (lm *LogManager) GetLogsFromBuffer(clientID string) []string {
	lm.Mutex.Lock()
	defer lm.Mutex.Unlock()


	logs := make([]string, len(lm.ClientBuffers[clientID]))
	copy(logs, lm.ClientBuffers[clientID])

	return logs
}

func (lm *LogManager) GetLogsFromFile(clientID string) ([]string, error) {

	filePath := fmt.Sprintf("../client_metrics/%s/log.txt", clientID)
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logs []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}


func (lm *LogManager) GetAllLogs(clientID string) ([]string, error) {

	logsFromBuffer := lm.GetLogsFromBuffer(clientID)
	logsFromFile, err := lm.GetLogsFromFile(clientID)

	if err != nil {
		return nil, err
	}

	return append(logsFromFile, logsFromBuffer...), nil
}

//////////// General Logs

func (lm *LogManager) GetGeneralLogsFromBuffer() []string {
	lm.Mutex.Lock()
	defer lm.Mutex.Unlock()

	logs := make([]string, len(lm.GeneralBuffer))
	copy(logs, lm.GeneralBuffer)

	return logs
}

func (lm *LogManager) GetGeneralLogsFromFile() ([]string, error) {

	filePath := "../client_metrics/log.txt"
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logs []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (lm *LogManager) GetAllGeneralLogs() ([]string, error) {

	logsFromBuffer := lm.GetGeneralLogsFromBuffer()
	logsFromFile, err := lm.GetGeneralLogsFromFile()

	if err != nil {
		return nil, err
	}

	return append(logsFromFile, logsFromBuffer...), nil
}

////////////////////////////////////////


func CreateClientMetrics() {
    clientMetricsDir := "../client_metrics"

    //check if the folder exists
    if _, err := os.Stat(clientMetricsDir); os.IsNotExist(err) {
        
        //create the metrics folder
        err := os.MkdirAll(clientMetricsDir, os.ModePerm)
        if err != nil {
            fmt.Printf("Error creating folder: %v\n", err)
            return
        }
        fmt.Println("folder created successfully.")
    }

    //create log file
    logFilePath := filepath.Join(clientMetricsDir, "log.txt")
    logFile, err := os.Create(logFilePath)
    if err != nil {
        fmt.Printf("Error creating log.txt: %v\n", err)
        return
    }
    defer logFile.Close()

    fmt.Println("log.txt created successfully.")
}


// clears all folderes on client_metrics and log.txt
func Cleanup() {

	clientMetricsDir := "../client_metrics"

    // check if the folder exists
    if _, err := os.Stat(clientMetricsDir); os.IsNotExist(err) {
        fmt.Println("folder does not exist, nothing to clean.")
        return
    }

    // temove all files and folders inside client_metrics and the folder itself
    err := os.RemoveAll(clientMetricsDir)
    if err != nil {
        fmt.Printf("Error during cleanup: %v\n", err)
    } else {
        fmt.Println("All files and folders removed successfully.")
    }
}
