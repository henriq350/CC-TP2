package sutils

import (
    "encoding/json"
    "fmt"
)

func ValidateJSON(data []byte) (bool, string) {
    var rawData interface{}
    err := json.Unmarshal(data, &rawData)
    if err != nil {
        return false, "Error: Couldn't unmarshal JSON data"
    }

    tasks, ok := rawData.([]interface{})
    if !ok {
        return false, "Error: JSON data is not a slice of tasks"
    }

    for _, task := range tasks {
        taskMap, ok := task.(map[string]interface{})
        if !ok {
            return false, "Error: Task is not a valid object"
        }

        if id, ok := taskMap["id"].(string); !ok {
            return false, fmt.Sprintf("Error: Task ID is not a valid string, got %v", id)
        }
        if frequency, ok := taskMap["frequency"].(float64); !ok || frequency < 0 || frequency > 65535 {
            return false, fmt.Sprintf("Error: Task frequency is not a valid uint16, got %v", frequency)
        }

        // Check devices
        if devices, ok := taskMap["devices"].([]interface{}); ok {
            for _, device := range devices {
                deviceMap, ok := device.(map[string]interface{})
                if !ok {
                    return false, "Error: Device is not a valid object"
                }

                if id, ok := deviceMap["id"].(string); !ok{
                    return false, fmt.Sprintf("Error: Device ID is not a valid string, got %v", id)
                }

                // Check deviceMetrics
                if deviceMetrics, ok := deviceMap["device_metrics"].(map[string]interface{}); ok {
                    if _, ok := deviceMetrics["cpu_usage"].(bool); !ok {
                        return false, "Error: Device CPU usage is not a valid bool"
                    }
                    if _, ok := deviceMetrics["ram_usage"].(bool); !ok {
                        return false, "Error: Device RAM usage is not a valid bool"
                    }
                    if interfaceStats, ok := deviceMetrics["interface_stats"].([]interface{}); ok {
                        for _, stat := range interfaceStats {
                            if _, ok := stat.(string); !ok {
                                return false, "Error: Device interface stat is not a valid string"
                            }
                        }
                    } else {
                        return false, "Error: Device interface stats is not a valid array"
                    }
                } else {
                    return false, "Error: Device metrics is not a valid object"
                }

                // Check linkMetrics
                if linkMetrics, ok := deviceMap["link_metrics"].(map[string]interface{}); ok {
                    // Check Bandwidth
                    if bandwidth, ok := linkMetrics["bandwidth"].(map[string]interface{}); ok {
                        if _, ok := bandwidth["tool"].(string); !ok {
                            return false, "Error: Bandwidth tool is not a valid string"
                        }
                        if _, ok := bandwidth["client"].(bool); !ok {
                            return false, "Error: Bandwidth client is not a valid bool"
                        }
                        if _, ok := bandwidth["server_addr"].(string); !ok {
                            return false, "Error: Bandwidth server address is not a valid string"
                        }
                        if duration, ok := bandwidth["duration"].(float64); !ok || duration < 0 || duration > 4294967295 {
                            return false, fmt.Sprintf("Error: Bandwidth duration is not a valid uint32, got %v", duration)
                        }
                        if _, ok := bandwidth["transport"].(string); !ok {
                            return false, "Error: Bandwidth transport is not a valid string"
                        }
                        if frequency, ok := bandwidth["frequency"].(float64); !ok || frequency < 0 || frequency > 65535 {
                            return false, fmt.Sprintf("Error: Bandwidth frequency is not a valid uint16, got %v", frequency)
                        }
                    } else {
                        return false, "Error: Bandwidth is not a valid object"
                    }

                    // Check Jitter
                    if jitter, ok := linkMetrics["jitter"].(map[string]interface{}); ok {
                        if _, ok := jitter["tool"].(string); !ok {
                            return false, "Error: Jitter tool is not a valid string"
                        }
                        if _, ok := jitter["client"].(bool); !ok {
                            return false, "Error: Jitter client is not a valid bool"
                        }
                        if _, ok := jitter["server_addr"].(string); !ok {
                            return false, "Error: Jitter server address is not a valid string"
                        }
                        if duration, ok := jitter["duration"].(float64); !ok || duration < 0 || duration > 4294967295 {
                            return false, fmt.Sprintf("Error: Jitter duration is not a valid uint32, got %v", duration)
                        }
                        if _, ok := jitter["transport"].(string); !ok {
                            return false, "Error: Jitter transport is not a valid string"
                        }
                        if frequency, ok := jitter["frequency"].(float64); !ok || frequency < 0 || frequency > 65535 {
                            return false, fmt.Sprintf("Error: Jitter frequency is not a valid uint16, got %v", frequency)
                        }
                    } else {
                        return false, "Error: Jitter is not a valid object"
                    }

                    // Check PacketLoss
                    if packetLoss, ok := linkMetrics["packet_loss"].(map[string]interface{}); ok {
                        if _, ok := packetLoss["destination"].(string); !ok {
                            return false, "Error: packetLoss destination is not a valid string"
                        }
                        if count, ok := packetLoss["count"].(float64); !ok || count < 0 || count > 65535 {
                            return false, fmt.Sprintf("Error: packetLoss count is not a valid uint16, got %v", count)
                        }
                        if frequency, ok := packetLoss["frequency"].(float64); !ok || frequency < 0 || frequency > 65535 {
                            return false, fmt.Sprintf("Error: packetLoss frequency is not a valid uint16, got %v", frequency)
                        }
                    } else {
                        return false, "Error: PacketLoss is not a valid object"
                    }

                    // Check Latency
                    if latency, ok := linkMetrics["latency"].(map[string]interface{}); ok {
                        if _, ok := latency["destination"].(string); !ok {
                            return false, "Error: Latency destination is not a valid string"
                        }
                        if count, ok := latency["count"].(float64); !ok || count < 0 || count > 65535 {
                            return false, fmt.Sprintf("Error: Latency count is not a valid uint16, got %v", count)
                        }
                        if frequency, ok := latency["frequency"].(float64); !ok || frequency < 0 || frequency > 65535 {
                            return false, fmt.Sprintf("Error: Latency frequency is not a valid uint16, got %v", frequency)
                        }
                    } else {
                        return false, "Error: Latency is not a valid object"
                    }

                    // Check AlertFlowConditions
                    if alertFlowConditions, ok := linkMetrics["alertflow_conditions"].(map[string]interface{}); ok {
                        if cpuUsage, ok := alertFlowConditions["cpu_usage"].(float64); !ok || cpuUsage < 0 || cpuUsage > 100 {
                            return false, fmt.Sprintf("Error: AlertFlowConditions CPU usage is not a valid float32 or out of range, got %v", cpuUsage)
                        }
                        if ramUsage, ok := alertFlowConditions["ram_usage"].(float64); !ok || ramUsage < 0 || ramUsage > 100 {
                            return false, fmt.Sprintf("Error: AlertFlowConditions RAM usage is not a valid float32 or out of range, got %v", ramUsage)
                        }
                        if _, ok := alertFlowConditions["interface_stats"].(float64); !ok {
                            return false, "Error: AlertFlowConditions interface stats is not a valid uint32"
                        }
                        if packetLoss, ok := alertFlowConditions["packet_loss"].(float64); !ok || packetLoss < 0 || packetLoss > 100 {
                            return false, fmt.Sprintf("Error: AlertFlowConditions packet loss is not a valid float32 or out of range, got %v", packetLoss)
                        }
                        if _, ok := alertFlowConditions["jitter"].(float64); !ok {
                            return false, "Error: AlertFlowConditions jitter is not a valid uint16"
                        }
                    } else {
                        return false, "Error: AlertFlowConditions is not a valid object"
                    }
                } else {
                    return false, "Error: LinkMetrics is not a valid object"
                }
            }
        } else {
            return false, "Error: Devices is not a valid array"
        }
    }
    return true, ""
}