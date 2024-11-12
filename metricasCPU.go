package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func getCPUStats() (float64, float64, error) {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 1 {
		return 0, 0, fmt.Errorf("no CPU data found")
	}

	// A primeira linha contém informações sobre a CPU
	cpuInfo := strings.Fields(lines[0])[1:]

	var total, idle float64
	for i, v := range cpuInfo {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, 0, err
		}
		if i == 3 { // idle é o quarto valor
			idle = value
		}
		total += value
	}

	return total, idle, nil
}

func getCPUUsage() (float64, error) {
	total1, idle1, err := getCPUStats()
	if err != nil {
		return 0, err
	}

	// Aguarda um segundo
	time.Sleep(1 * time.Second)

	total2, idle2, err := getCPUStats()
	if err != nil {
		return 0, err
	}

	totalDiff := total2 - total1
	idleDiff := idle2 - idle1

	usage := (totalDiff - idleDiff) / totalDiff * 100
	return usage, nil
}

func main() {
	for {
		usage, err := getCPUUsage()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("CPU Usage: %.2f%%\n", usage)
		time.Sleep(1 * time.Second) // Aguarda um segundo antes da próxima leitura
	}
}
