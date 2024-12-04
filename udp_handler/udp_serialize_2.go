package udp_handler

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// PacketType represents the type of UDP packet
type PacketType uint8

const (
	TaskPacket     PacketType = 0
	RegisterPacket PacketType = 1
	ReportPacket   PacketType = 2
)

// Flags represents the packet flags
type Flags struct {
	SYN bool
	ACK bool
	RET bool
}

// AgentRegistration represents the registration data
type AgentRegistration struct {
	AgentID string
	IPv4    string
}

// TaskRecord represents a single task configuration
type TaskRecord struct {
	TaskID       string
	Name         string
	Value        string
	Threshold    float64
	Duration     uint32
	PacketCount  uint32
	Frequency    uint32
	ReportFreq   uint32
	CriticalValues []string
}

// ReportRecord represents a single report value
type ReportRecord struct {
	TaskID string
	Name   string
	Value  string
}

// Packet represents the complete UDP packet structure
type Packet struct {
	Type           PacketType
	SequenceNumber uint32
	AckNumber      uint32
	Flags          Flags
	Data           interface{} // Can be AgentRegistration, []TaskRecord, or []ReportRecord
}

// Serialize converts a Packet into a byte slice
func (p *Packet) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write packet type (2 bits)
	packetType := uint8(p.Type) << 6
	err := binary.Write(buf, binary.BigEndian, packetType)
	if err != nil {
		return nil, err
	}

	// Write sequence and ack numbers
	err = binary.Write(buf, binary.BigEndian, p.SequenceNumber)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, p.AckNumber)
	if err != nil {
		return nil, err
	}

	// Write flags
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
	err = binary.Write(buf, binary.BigEndian, flags)
	if err != nil {
		return nil, err
	}

	// Serialize data based on packet type
	switch p.Type {
	case RegisterPacket:
		reg, ok := p.Data.(AgentRegistration)
		if !ok {
			return nil, errors.New("invalid data type for RegisterPacket")
		}
		err = writeString(buf, reg.AgentID)
		if err != nil {
			return nil, err
		}
		err = writeString(buf, reg.IPv4)
		if err != nil {
			return nil, err
		}

	case TaskPacket:
		tasks, ok := p.Data.([]TaskRecord)
		if !ok {
			return nil, errors.New("invalid data type for TaskPacket")
		}
		// Write number of tasks
		err = binary.Write(buf, binary.BigEndian, uint32(len(tasks)))
		if err != nil {
			return nil, err
		}
		// Write each task
		for _, task := range tasks {
			err = writeString(buf, task.TaskID)
			if err != nil {
				return nil, err
			}
			err = writeString(buf, task.Name)
			if err != nil {
				return nil, err
			}
			err = writeString(buf, task.Value)
			if err != nil {
				return nil, err
			}
			err = binary.Write(buf, binary.BigEndian, task.Threshold)
			if err != nil {
				return nil, err
			}
			err = binary.Write(buf, binary.BigEndian, task.Duration)
			if err != nil {
				return nil, err
			}
			err = binary.Write(buf, binary.BigEndian, task.PacketCount)
			if err != nil {
				return nil, err
			}
			err = binary.Write(buf, binary.BigEndian, task.Frequency)
			if err != nil {
				return nil, err
			}
			err = binary.Write(buf, binary.BigEndian, task.ReportFreq)
			if err != nil {
				return nil, err
			}
			// Write critical values
			err = binary.Write(buf, binary.BigEndian, uint32(len(task.CriticalValues)))
			if err != nil {
				return nil, err
			}
			for _, cv := range task.CriticalValues {
				err = writeString(buf, cv)
				if err != nil {
					return nil, err
				}
			}
		}

	case ReportPacket:
		reports, ok := p.Data.([]ReportRecord)
		if !ok {
			return nil, errors.New("invalid data type for ReportPacket")
		}
		// Write number of reports
		err = binary.Write(buf, binary.BigEndian, uint32(len(reports)))
		if err != nil {
			return nil, err
		}
		// Write each report
		for _, report := range reports {
			err = writeString(buf, report.TaskID)
			if err != nil {
				return nil, err
			}
			err = writeString(buf, report.Name)
			if err != nil {
				return nil, err
			}
			err = writeString(buf, report.Value)
			if err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

// Deserialize converts a byte slice into a Packet
func Deserialize(data []byte) (*Packet, error) {
	buf := bytes.NewReader(data)
	packet := &Packet{}

	// Read packet type
	var packetTypeByte uint8
	err := binary.Read(buf, binary.BigEndian, &packetTypeByte)
	if err != nil {
		return nil, err
	}
	packet.Type = PacketType(packetTypeByte >> 6)

	// Read sequence and ack numbers
	err = binary.Read(buf, binary.BigEndian, &packet.SequenceNumber)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.BigEndian, &packet.AckNumber)
	if err != nil {
		return nil, err
	}

	// Read flags
	var flags uint8
	err = binary.Read(buf, binary.BigEndian, &flags)
	if err != nil {
		return nil, err
	}
	packet.Flags = Flags{
		SYN: flags&(1<<2) != 0,
		ACK: flags&(1<<1) != 0,
		RET: flags&1 != 0,
	}

	// Deserialize data based on packet type
	switch packet.Type {
	case RegisterPacket:
		reg := AgentRegistration{}
		reg.AgentID, err = readString(buf)
		if err != nil {
			return nil, err
		}
		reg.IPv4, err = readString(buf)
		if err != nil {
			return nil, err
		}
		packet.Data = reg

	case TaskPacket:
		var numTasks uint32
		err = binary.Read(buf, binary.BigEndian, &numTasks)
		if err != nil {
			return nil, err
		}
		tasks := make([]TaskRecord, numTasks)
		for i := range tasks {
			task := TaskRecord{}
			task.TaskID, err = readString(buf)
			if err != nil {
				return nil, err
			}
			task.Name, err = readString(buf)
			if err != nil {
				return nil, err
			}
			task.Value, err = readString(buf)
			if err != nil {
				return nil, err
			}
			err = binary.Read(buf, binary.BigEndian, &task.Threshold)
			if err != nil {
				return nil, err
			}
			err = binary.Read(buf, binary.BigEndian, &task.Duration)
			if err != nil {
				return nil, err
			}
			err = binary.Read(buf, binary.BigEndian, &task.PacketCount)
			if err != nil {
				return nil, err
			}
			err = binary.Read(buf, binary.BigEndian, &task.Frequency)
			if err != nil {
				return nil, err
			}
			err = binary.Read(buf, binary.BigEndian, &task.ReportFreq)
			if err != nil {
				return nil, err
			}
			var numCriticalValues uint32
			err = binary.Read(buf, binary.BigEndian, &numCriticalValues)
			if err != nil {
				return nil, err
			}
			task.CriticalValues = make([]string, numCriticalValues)
			for j := range task.CriticalValues {
				task.CriticalValues[j], err = readString(buf)
				if err != nil {
					return nil, err
				}
			}
			tasks[i] = task
		}
		packet.Data = tasks

	case ReportPacket:
		var numReports uint32
		err = binary.Read(buf, binary.BigEndian, &numReports)
		if err != nil {
			return nil, err
		}
		reports := make([]ReportRecord, numReports)
		for i := range reports {
			report := ReportRecord{}
			report.TaskID, err = readString(buf)
			if err != nil {
				return nil, err
			}
			report.Name, err = readString(buf)
			if err != nil {
				return nil, err
			}
			report.Value, err = readString(buf)
			if err != nil {
				return nil, err
			}
			reports[i] = report
		}
		packet.Data = reports
	}

	return packet, nil
}

// Helper functions remain the same
func writeString(buf *bytes.Buffer, s string) error {
	err := binary.Write(buf, binary.BigEndian, uint32(len(s)))
	if err != nil {
		return err
	}
	_, err = buf.WriteString(s)
	return err
}

func readString(buf *bytes.Reader) (string, error) {
	var length uint32
	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return "", err
	}
	strBytes := make([]byte, length)
	_, err = buf.Read(strBytes)
	if err != nil {
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
	
	case []TaskRecord:
		fmt.Printf("Task Records (%d items)\n", len(data))
		for i, task := range data {
			fmt.Printf("  Task #%d:\n", i+1)
			fmt.Printf("    Task ID: %s\n", task.TaskID)
			fmt.Printf("    Name: %s\n", task.Name)
			fmt.Printf("    Value: %s\n", task.Value)
			fmt.Printf("    Threshold: %f\n", task.Threshold)
			fmt.Printf("    Duration: %d\n", task.Duration)
			fmt.Printf("    Packet Count: %d\n", task.PacketCount)
			fmt.Printf("    Frequency: %d\n", task.Frequency)
			fmt.Printf("    Report Frequency: %d\n", task.ReportFreq)
			fmt.Printf("    Critical Values: %v\n", task.CriticalValues)
		}
	
	case []ReportRecord:
		fmt.Printf("Report Records (%d items)\n", len(data))
		for i, report := range data {
			fmt.Printf("  Report #%d:\n", i+1)
			fmt.Printf("    Task ID: %s\n", report.TaskID)
			fmt.Printf("    Name: %s\n", report.Name)
			fmt.Printf("    Value: %s\n", report.Value)
		}
	
	default:
		fmt.Printf("Unknown data type\n")
	}
	fmt.Println("=====================")
}