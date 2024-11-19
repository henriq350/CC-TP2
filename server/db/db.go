package db

import (
	"os"
	"fmt"
	"path/filepath"
	"strings"
)

func ParsePacket(clientName, fileName, data string) {
	err := SaveFile(clientName, fileName, []byte(data))
    if err != nil {
        fmt.Println("Error saving file:", err)
    } else {
        fmt.Println("File saved successfully")
    }
}

// Cria pasta para cada cliente dentro da pasta ClientMetrics
func CreateFolder(client string) error {
	clientDir := filepath.Join("../ClientMetrics", client)
	err := os.MkdirAll(clientDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("Error creatring CLient's folder: %v", err)
	}
	return nil
}

func SaveFile(clientName, fileName string, data []byte) error {
    clientDir := filepath.Join("../ClientMetrics", clientName)
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