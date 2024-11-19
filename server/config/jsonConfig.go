package config

type Task struct {
	ID        string    `json:"id"`
	Frequency uint16    `json:"frequency"`
	Devices   []Device `json:"devices"`
}

type Device struct {
	ID            string          `json:"id"`
	DeviceMetrics DeviceMetrics `json:"device_metrics"`
	LinkMetrics   LinkMetrics   `json:"link_metrics"`
}

type DeviceMetrics struct {
	CPUUsage       bool     `json:"cpu_usage"`
	RAMUsage       bool     `json:"ram_usage"`
	InterfaceStats []string `json:"interface_stats"`
}

type LinkMetrics struct {
	Bandwidth           MetricsConfig       `json:"bandwidth"`
	Jitter              MetricsConfig       `json:"jitter"`
	PacketLoss          MetricsConfig       `json:"packet_loss"`
	Latency             Latency             `json:"latency"`
	AlertFlowConditions AlertFlowConditions `json:"alertflow_conditions"`
}

type MetricsConfig struct {
	Tool       string `json:"tool"`
	Client     bool   `json:"client"`
	ServerAddr string `json:"server_addr"`
	Duration   uint32  `json:"duration"`
	Transport  string `json:"transport"`
	Frequency  uint16  `json:"frequency"`
}

type Latency struct {
	Destination string `json:"destination"`
	Count       uint16  `json:"count"`
	Frequency   uint16  `json:"frequency"`
}

type AlertFlowConditions struct {
	CPUUsage       float32 `json:"cpu_usage"`
	RAMUsage       float32 `json:"ram_usage"`
	InterfaceStats uint32   `json:"interface_stats"`
	PacketLoss     float32 `json:"packet_loss"`
	Jitter         uint16   `json:"jitter"`
}
