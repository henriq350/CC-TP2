package udp_handler

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// PacketType represents different types of packets in the system
type PacketType uint8

const (
	TaskPacket     PacketType = 0
	RegisterPacket PacketType = 1
	ReportPacket   PacketType = 2
)

// Flags represents the packet control flags
type Flags struct {
	SYN bool
	ACK bool
	RET bool
}

// AgentRegistration contains agent registration details
type AgentRegistration struct {
	AgentID  string
	IPv4     string
	ClientID string
}

// TaskRecord represents a task configuration
type TaskRecord struct {
	TaskID         string
	Name           string
	Value          string
	DestinationIp  string    // Added destinationIp field
	Threshold      float64
	Duration       uint32
	PacketCount    uint32
	Frequency      uint32
	ReportFreq     uint32
	CriticalValues []string
	ClientID       string
}

// ReportRecord represents a task execution report
type ReportRecord struct {
	TaskID        string
	Name          string
	Value         string
	DestinationIp string    // Added destinationIp field
	ClientID      string
}

// Packet represents the main communication structure
type Packet struct {
	Type           PacketType
	SequenceNumber uint32
	AckNumber      uint32
	Flags          Flags
	Data           interface{}
}

// Serialize converts a packet into a byte slice
func (p *Packet) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	if err := p.serializeHeader(buf); err != nil {
		return nil, fmt.Errorf("header serialization failed: %w", err)
	}

	if err := p.serializeData(buf); err != nil {
		return nil, fmt.Errorf("data serialization failed: %w", err)
	}

	return buf.Bytes(), nil
}

func (p *Packet) serializeHeader(buf *bytes.Buffer) error {
	packetType := uint8(p.Type) << 6
	if err := binary.Write(buf, binary.BigEndian, packetType); err != nil {
		return err
	}
	
	for _, v := range []interface{}{p.SequenceNumber, p.AckNumber} {
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			return err
		}
	}

	flags := p.serializeFlags()
	return binary.Write(buf, binary.BigEndian, flags)
}

func (p *Packet) serializeFlags() uint8 {
	var flags uint8
	if p.Flags.SYN {
		flags |= 1 << 2
	}
	if p.Flags.ACK {
		flags |= 1 << 1
	}
	if p.Flags.RET {
		flags |= 1
	}
	return flags
}

func (p *Packet) serializeData(buf *bytes.Buffer) error {
	switch p.Type {
	case RegisterPacket:
		return p.serializeRegistration(buf)
	case TaskPacket:
		return p.serializeTasks(buf)
	case ReportPacket:
		return p.serializeReports(buf)
	default:
		return fmt.Errorf("unknown packet type: %v", p.Type)
	}
}

func (p *Packet) serializeRegistration(buf *bytes.Buffer) error {
	reg, ok := p.Data.(AgentRegistration)
	if !ok {
		return errors.New("invalid data type for RegisterPacket")
	}
	
	for _, s := range []string{reg.AgentID, reg.IPv4, reg.ClientID} {
		if err := writeString(buf, s); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packet) serializeTasks(buf *bytes.Buffer) error {
	tasks, ok := p.Data.([]TaskRecord)
	if !ok {
		return errors.New("invalid data type for TaskPacket")
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(len(tasks))); err != nil {
		return err
	}

	for _, task := range tasks {
		if err := serializeTaskRecord(buf, task); err != nil {
			return err
		}
	}
	return nil
}

func serializeTaskRecord(buf *bytes.Buffer, task TaskRecord) error {
	strings := []string{task.TaskID, task.Name, task.Value, task.DestinationIp, task.ClientID}
	for _, s := range strings {
		if err := writeString(buf, s); err != nil {
			return err
		}
	}

	values := []interface{}{task.Threshold, task.Duration, task.PacketCount, task.Frequency, task.ReportFreq}
	for _, v := range values {
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			return err
		}
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(len(task.CriticalValues))); err != nil {
		return err
	}

	for _, cv := range task.CriticalValues {
		if err := writeString(buf, cv); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packet) serializeReports(buf *bytes.Buffer) error {
	reports, ok := p.Data.([]ReportRecord)
	if !ok {
		return errors.New("invalid data type for ReportPacket")
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(len(reports))); err != nil {
		return err
	}

	for _, report := range reports {
		strings := []string{report.TaskID, report.Name, report.Value, report.DestinationIp, report.ClientID}
		for _, s := range strings {
			if err := writeString(buf, s); err != nil {
				return err
			}
		}
	}
	return nil
}

// Deserialize converts a byte slice back into a packet
func Deserialize(data []byte) (*Packet, error) {
	buf := bytes.NewReader(data)
	packet := &Packet{}

	if err := deserializeHeader(buf, packet); err != nil {
		return nil, fmt.Errorf("header deserialization failed: %w", err)
	}

	if err := deserializeData(buf, packet); err != nil {
		return nil, fmt.Errorf("data deserialization failed: %w", err)
	}

	return packet, nil
}

func deserializeHeader(buf *bytes.Reader, packet *Packet) error {
	var packetTypeByte uint8
	if err := binary.Read(buf, binary.BigEndian, &packetTypeByte); err != nil {
		return err
	}
	packet.Type = PacketType(packetTypeByte >> 6)

	for _, v := range []interface{}{&packet.SequenceNumber, &packet.AckNumber} {
		if err := binary.Read(buf, binary.BigEndian, v); err != nil {
			return err
		}
	}

	var flags uint8
	if err := binary.Read(buf, binary.BigEndian, &flags); err != nil {
		return err
	}
	packet.Flags = deserializeFlags(flags)
	return nil
}

func deserializeFlags(flags uint8) Flags {
	return Flags{
		SYN: flags&(1<<2) != 0,
		ACK: flags&(1<<1) != 0,
		RET: flags&1 != 0,
	}
}

func deserializeData(buf *bytes.Reader, packet *Packet) error {
	switch packet.Type {
	case RegisterPacket:
		return deserializeRegistration(buf, packet)
	case TaskPacket:
		return deserializeTasks(buf, packet)
	case ReportPacket:
		return deserializeReports(buf, packet)
	default:
		return fmt.Errorf("unknown packet type: %v", packet.Type)
	}
}

func deserializeRegistration(buf *bytes.Reader, packet *Packet) error {
	reg := AgentRegistration{}
	var err error
	
	if reg.AgentID, err = readString(buf); err != nil {
		return err
	}
	if reg.IPv4, err = readString(buf); err != nil {
		return err
	}
	if reg.ClientID, err = readString(buf); err != nil {
		return err
	}
	
	packet.Data = reg
	return nil
}

func deserializeTasks(buf *bytes.Reader, packet *Packet) error {
	var numTasks uint32
	if err := binary.Read(buf, binary.BigEndian, &numTasks); err != nil {
		return err
	}

	tasks := make([]TaskRecord, numTasks)
	for i := range tasks {
		if err := deserializeTaskRecord(buf, &tasks[i]); err != nil {
			return err
		}
	}

	packet.Data = tasks
	return nil
}

func deserializeTaskRecord(buf *bytes.Reader, task *TaskRecord) error {
	var err error
	if task.TaskID, err = readString(buf); err != nil {
		return err
	}
	if task.Name, err = readString(buf); err != nil {
		return err
	}
	if task.Value, err = readString(buf); err != nil {
		return err
	}
	if task.DestinationIp, err = readString(buf); err != nil {
		return err
	}
	if task.ClientID, err = readString(buf); err != nil {
		return err
	}

	values := []interface{}{&task.Threshold, &task.Duration, &task.PacketCount, &task.Frequency, &task.ReportFreq}
	for _, v := range values {
		if err := binary.Read(buf, binary.BigEndian, v); err != nil {
			return err
		}
	}

	var numCriticalValues uint32
	if err := binary.Read(buf, binary.BigEndian, &numCriticalValues); err != nil {
		return err
	}

	task.CriticalValues = make([]string, numCriticalValues)
	for j := range task.CriticalValues {
		if task.CriticalValues[j], err = readString(buf); err != nil {
			return err
		}
	}
	return nil
}

func deserializeReports(buf *bytes.Reader, packet *Packet) error {
	var numReports uint32
	if err := binary.Read(buf, binary.BigEndian, &numReports); err != nil {
		return err
	}

	reports := make([]ReportRecord, numReports)
	for i := range reports {
		report := &reports[i]
		var err error
		if report.TaskID, err = readString(buf); err != nil {
			return err
		}
		if report.Name, err = readString(buf); err != nil {
			return err
		}
		if report.Value, err = readString(buf); err != nil {
			return err
		}
		if report.DestinationIp, err = readString(buf); err != nil {
			return err
		}
		if report.ClientID, err = readString(buf); err != nil {
			return err
		}
	}
	packet.Data = reports
	return nil
}

// Helper functions
func writeString(buf *bytes.Buffer, s string) error {
	if err := binary.Write(buf, binary.BigEndian, uint32(len(s))); err != nil {
		return err
	}
	_, err := buf.WriteString(s)
	return err
}

func readString(buf *bytes.Reader) (string, error) {
	var length uint32
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return "", err
	}
	strBytes := make([]byte, length)
	if _, err := buf.Read(strBytes); err != nil {
		return "", err
	}
	return string(strBytes), nil
}

func (pt PacketType) String() string {
	switch pt {
	case TaskPacket:
		return "TaskPacket"
	case RegisterPacket:
		return "RegisterPacket"
	case ReportPacket:
		return "ReportPacket"
	default:
		return fmt.Sprintf("Unknown(%d)", pt)
	}
}

// Print prints the packet information to console
func (p *Packet) Print() {
	fmt.Printf("=== Packet Information ===\n")
	fmt.Printf("Type: %s (value: %d)\n", p.Type, p.Type)
	fmt.Printf("Sequence Number: %d\n", p.SequenceNumber)
	fmt.Printf("Ack Number: %d\n", p.AckNumber)
	fmt.Printf("Flags: SYN=%v ACK=%v RET=%v\n", p.Flags.SYN, p.Flags.ACK, p.Flags.RET)
	
	fmt.Printf("Data: ")
	switch data := p.Data.(type) {
	case AgentRegistration:
		fmt.Printf("Registration Data\n")
		fmt.Printf("  Agent ID: %s\n", data.AgentID)
		fmt.Printf("  IPv4: %s\n", data.IPv4)
		fmt.Printf("  Client ID: %s\n", data.ClientID)
	
	case []TaskRecord:
		fmt.Printf("Task Records (%d items)\n", len(data))
		for i, task := range data {
			fmt.Printf("  Task #%d:\n", i+1)
			fmt.Printf("    Task ID: %s\n", task.TaskID)
			fmt.Printf("    Name: %s\n", task.Name)
			fmt.Printf("    Value: %s\n", task.Value)
			fmt.Printf("    Destination IP: %s\n", task.DestinationIp)
			fmt.Printf("    Threshold: %f\n", task.Threshold)
			fmt.Printf("    Duration: %d\n", task.Duration)
			fmt.Printf("    Packet Count: %d\n", task.PacketCount)
			fmt.Printf("    Frequency: %d\n", task.Frequency)
			fmt.Printf("    Report Frequency: %d\n", task.ReportFreq)
			fmt.Printf("    Critical Values: %v\n", task.CriticalValues)
			fmt.Printf("    Client ID: %s\n", task.ClientID)
		}
	
	case []ReportRecord:
		fmt.Printf("Report Records (%d items)\n", len(data))
		for i, report := range data {
			fmt.Printf("  Report #%d:\n", i+1)
			fmt.Printf("    Task ID: %s\n", report.TaskID)
			fmt.Printf("    Name: %s\n", report.Name)
			fmt.Printf("    Value: %s\n", report.Value)
			fmt.Printf("    Destination IP: %s\n", report.DestinationIp)
			fmt.Printf("    Client ID: %s\n", report.ClientID)
		}
	
	default:
		fmt.Printf("Unknown data type\n")
	}
	fmt.Println("=====================")
}