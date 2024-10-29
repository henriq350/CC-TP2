package config

type Task struct {
	ID        int16    `json:"id"`
	Frequency int16    `json:"frequency"`
	Devices   []Device `json:"devices"`
}

type Device struct {
	ID            int8          `json:"id"`
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
	Duration   int32  `json:"duration"`
	Transport  string `json:"transport"`
	Frequency  int16  `json:"frequency"`
}

type Latency struct {
	Destination string `json:"destination"`
	Count       int16  `json:"count"`
	Frequency   int16  `json:"frequency"`
}

type AlertFlowConditions struct {
	CPUUsage       float32 `json:"cpu_usage"`
	RAMUsage       float32 `json:"ram_usage"`
	InterfaceStats int32   `json:"interface_stats"`
	PacketLoss     float32 `json:"packet_loss"`
	Jitter         int16   `json:"jitter"`
}
