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
    Mutex   sync.Mutex
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

// Cria pasta para cada cliente dentro da pasta client_metrics
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

// TODO - Implementar uma formataca em uma string para ser salva em um arquivo
func FormatString(data []string) (string, string) {

    var fmtData strings.Builder

    currentTime := time.Now().Format("2024-11-14 15:04:05")
    fmtData.WriteString(fmt.Sprintf("Received at: %s\n\n", currentTime))

    for i := 0; i < len(data); i += 5 {
        metric := data[i]
        value := data[i+1]
        device := data[i+2]
        sourceIP := data[i+3]
        destIP := data[i+4]
        fmtData.WriteString(fmt.Sprintf("%s: %s\nDevice: %s\nSource IP: %s\nDestination IP: %s\n", metric, value, device, sourceIP, destIP))
    }

    return fmtData.String(), currentTime
}



/////////////////////////////////////////// LOG LOGIC


/////////////////////////////  WRITE LOGS

// Construtor do LogManager
func NewLogManager() *LogManager {
    return &LogManager{
        ClientBuffers: make(map[string][]string),
        GeneralBuffer: []string{},
    }
}


// log to client
func (lm *LogManager) AddLog(clientID, log string) {
    lm.Mutex.Lock()
    defer lm.Mutex.Unlock()

    if _, exists := lm.ClientBuffers[clientID]; !exists {
        lm.ClientBuffers[clientID] = []string{}
    }

    // Adiciona o log ao client
    lm.ClientBuffers[clientID] = append(lm.ClientBuffers[clientID], log)

    // Formata o log para ser usado no log geral
    entry := fmt.Sprintf("[%s] %s", clientID, log)
    lm.GeneralBuffer = append(lm.GeneralBuffer, entry)

}   

// remove buffer quando client da terminate
func (lm *LogManager) RemoveClientBuffer(clientID string) {
    lm.Mutex.Lock()
    defer lm.Mutex.Unlock()

    // Guarda o que tem dentro do buffer no disco antes de o eliminar
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

// guarda periodicamente os logs do buffer no disco
func (lm *LogManager) PersistLogs() {

    for {
        lm.Mutex.Lock()

        for clientID, logs := range lm.ClientBuffers {
        
            if len(logs) > 0 {

                filePath := fmt.Sprintf("../client_metrics/%s/log.txt", clientID)
                file, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

                for _, log := range logs {
                    file.WriteString(log + "\n")
                }

                lm.ClientBuffers[clientID] = nil // limpa o buffer dps de guardar dados
                file.Close()
            }
        }

        if len(lm.GeneralBuffer) > 0 {

            filePath := "../client_metrics/log.txt"
            file, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

            for _, log := range lm.GeneralBuffer {
                file.WriteString(log + "\n")
            }

            lm.GeneralBuffer = nil // Limpar buffer geral após salvar
            file.Close()
        }

        lm.Mutex.Unlock()
        time.Sleep(60 * time.Second) // intervalo de escrita
    }
}


/////////////////////////////  READ LOGS

//////////// Client Logs

// Lê os logs do buffer para um cliente específico
func (lm *LogManager) GetLogsFromBuffer(clientID string) []string {
    lm.Mutex.Lock()
    defer lm.Mutex.Unlock()

    // Copia para prevenir concorrencia
    logs := make([]string, len(lm.ClientBuffers[clientID]))
    copy(logs, lm.ClientBuffers[clientID] )

    return lm.ClientBuffers[clientID]
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

// combinacao dos logs em buffer e dos logs em memoria
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
