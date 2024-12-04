package metrics
import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Função para medir latência e perda de pacotes com `ping`
func PingMetrics(target string, count int) (float64, float64, error) {
	cmd := exec.Command("ping", "-c", strconv.Itoa(count), target)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return 0, 0, err
	}

	var latency float64
	var packetLoss float64
	scanner := bufio.NewScanner(&out)

	for scanner.Scan() {
		line := scanner.Text()

		// Verificar a linha que contém "rtt min/avg/max/mdev"
		if strings.Contains(line, "rtt min/avg/max/mdev") {
			// Extrair a latência média
			parts := strings.Split(line, "/")
			latency, err = strconv.ParseFloat(parts[4], 64)
			if err != nil {
				return 0, 0, err
			}
		}

		// Verificar a linha que contém a taxa de perda de pacotes
		if strings.Contains(line, "packet loss") {
			re := regexp.MustCompile(`(\d+\.?\d*)% packet loss`)
			match := re.FindStringSubmatch(line)
			if len(match) > 1 {
				packetLoss, err = strconv.ParseFloat(match[1], 64)
				if err != nil {
					return 0, 0, err
				}
			}
		}
	}

	return latency, packetLoss, nil
}

// Função para medir largura de banda e jitter com `iperf3`
func IperfMetrics(target string, duration int) (float64, float64, error) {
    cmd := exec.Command("iperf3", "-c", target, "-u", "-t", strconv.Itoa(duration))

    var out bytes.Buffer
    cmd.Stdout = &out

    err := cmd.Run()
    if err != nil {
        return 0, 0, err
    }

    var bandwidth float64
    var jitter float64
    scanner := bufio.NewScanner(&out)

    // Regexes para extrair largura de banda e jitter
    bandwidthRegex := regexp.MustCompile(`([0-9.]+)\s(Mbits|Kbits|Gbits)/sec`)
    jitterRegex := regexp.MustCompile(`([0-9.]+)\sms`)

    for scanner.Scan() {
        line := scanner.Text()

        // Largura de banda
        if bandwidthMatch := bandwidthRegex.FindStringSubmatch(line); bandwidthMatch != nil {
            bandwidth, err = strconv.ParseFloat(bandwidthMatch[1], 64)
            if err != nil {
                return 0, 0, err
            }

            // Ajustar unidade
            switch bandwidthMatch[2] {
            case "Kbits":
                bandwidth /= 1000
            case "Gbits":
                bandwidth *= 1000
            }
        }

        // Jitter
        if jitterMatch := jitterRegex.FindStringSubmatch(line); jitterMatch != nil {
            jitter, err = strconv.ParseFloat(jitterMatch[1], 64)
            if err != nil {
                return 0, 0, err
            }
        }
    }

    // Verificar se as métricas foram encontradas
    if bandwidth == 0 && jitter == 0 {
        return 0, 0, fmt.Errorf("nenhuma métrica encontrada na saída do iperf3")
    }

    return bandwidth, jitter, nil
}


// func main() {
// 	target := "127.0.0.1"  // Servidor de teste (alvo)
// 	pingCount := 4       // Número de pacotes de ping
// 	iperfDuration := 10  // Duração do teste do iperf em segundos

// 	// Medindo latência e perda de pacotes
// 	latency, packetLoss, err := pingMetrics(target, pingCount)
// 	if err != nil {
// 		fmt.Println("Erro ao medir latência e perda de pacotes:", err)
// 		return
// 	}
// 	fmt.Printf("Latência média: %.2f ms\n", latency)
// 	fmt.Printf("Perda de pacotes: %.2f%%\n", packetLoss)

// 	// Medindo largura de banda e jitter
// 	bandwidth, jitter, err := iperfMetrics(target, iperfDuration)
// 	if err != nil {
// 		fmt.Println("Erro ao medir largura de banda e jitter:", err)
// 		return
// 	}
// 	fmt.Printf("Largura de banda: %.2f Mbps\n", bandwidth)
// 	fmt.Printf("Jitter: %.2f ms\n", jitter)
// }
