package metrics

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
)

func sendCPUMetrics(cpu string, time string, device string) {}

// Função para obter as estatísticas da CPU
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

// Função para calcular a utilização da CPU
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

// Função para obter a utilização de RAM
func getRAMUsage() (float64, error) {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, err
	}

	var totalMem, freeMem float64
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		if fields[0] == "MemTotal:" {
			totalMem, err = strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return 0, err
			}
		} else if fields[0] == "MemAvailable:" {
			freeMem, err = strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return 0, err
			}
		}
	}

	usedMem := totalMem - freeMem
	usage := (usedMem / totalMem) * 100
	return usage, nil
}

// Função para listar nomes das interfaces de rede
func getNetworkInterfaces() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var names []string
	for _, iface := range interfaces {
		names = append(names, iface.Name)
	}
	return names, nil
}

func main() {
	for {
		cpuUsage, err := getCPUUsage()
		if err != nil {
			fmt.Println("Erro ao obter uso da CPU:", err)
			return
		}

		ramUsage, err := getRAMUsage()
		if err != nil {
			fmt.Println("Erro ao obter uso da RAM:", err)
			return
		}

		networkInterfaces, err := getNetworkInterfaces()
		if err != nil {
			fmt.Println("Erro ao obter interfaces de rede:", err)
			return
		}

		fmt.Printf("CPU Usage: %.2f%%\n", cpuUsage)
		fmt.Printf("RAM Usage: %.2f%%\n", ramUsage)
		fmt.Println("Network Interfaces:", networkInterfaces)

		time.Sleep(1 * time.Second) // Aguarda um segundo antes da próxima leitura
	}
}
