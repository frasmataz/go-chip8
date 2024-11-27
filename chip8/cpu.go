package chip8

import (
	"fmt"
	"strings"
)

const debugLog = true

type Cpu struct {
	V     [0x10]uint8  // General-purpose 8-bit registers V0 - VF
	I     uint16       // 16-bit register, used to hold memory addresses
	PC    uint16       // Program counter - 16-bit
	SP    uint8        // Stack pointer - 8-bit
	DT    uint8        // Delay timer - 8-bit - dec at 60Hz when non-zero
	ST    uint8        // Sound timer - 8-bit - dec at 60Hz when non-zero
	Stack [0x10]uint16 // Stack - 16 16-bit values

	Memory  *Memory
	Display *Display
}

func NewCpu() *Cpu {
	cpu := new(Cpu)
	cpu.Memory = NewMemory()
	cpu.Display = NewDisplay()
	cpu.PC = 0x200
	return cpu
}

func (cpu *Cpu) Tick() error {
	// Fetch
	opcode, err := cpu.Memory.Get16(cpu.PC)
	if err != nil {
		return err
	}
	cpu.PC += 2

	// Decode
	if opcode&0xF000 == 0x0000 {
		if opcode == 0x00E0 {
			// CLS - 0x00E0
			cpu.CLS()
		} else if opcode == 0x00EE {
			// RET - 0x00EE
			err := cpu.RET()
			if err != nil {
				return err
			}
		} else {
			// SYS 0nnn
			return fmt.Errorf("SYS not implemented")
		}
	} else if opcode&0xF000 == 0x1000 {
		// JP - 1nnn
		err := cpu.JP(opcode)
		if err != nil {
			return err
		}
	} else if opcode&0xF000 == 0x2000 {
		// CALL - 2nnn
		err := cpu.CALL(opcode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cpu *Cpu) CLS() {
	if debugLog {
		fmt.Printf("CLS\n")
	}

	cpu.Display = NewDisplay()
}

func (cpu *Cpu) RET() error {
	if debugLog {
		fmt.Printf("RET\n")
	}

	if cpu.SP > 0x10 {
		return fmt.Errorf("stack overflow on RET - SP is > 0x10 - cpu state: %v", cpu.GetPrettyCpuState())
	} else if cpu.SP == 0x00 {
		return fmt.Errorf("stack underflow on RET - SP is 0x00 - cpu state: %v", cpu.GetPrettyCpuState())
	}

	cpu.PC = uint16(cpu.Stack[cpu.SP-1])
	cpu.SP--

	return nil
}

func (cpu *Cpu) JP(opcode uint16) error {
	target := opcode & 0x0FFF

	if debugLog {
		fmt.Printf("JP %04x\n", target)
	}

	if target > 0xFFE {
		return fmt.Errorf("target out of range for JP: %04x, max: 0ffe", target)
	}

	cpu.PC = target
	return nil
}

func (cpu *Cpu) CALL(opcode uint16) error {
	target := opcode & 0x0FFF

	if debugLog {
		fmt.Printf("CALL %04x\n", target)
	}

	if cpu.SP > 0x0F {
		return fmt.Errorf("stack overflow on CALL - SP is > 0x0F - cpu state: %v", cpu.GetPrettyCpuState())
	}

	cpu.Stack[cpu.SP] = cpu.PC
	cpu.SP++
	cpu.PC = target

	return nil
}

func (cpu *Cpu) GetPrettyCpuState() string {
	var sb strings.Builder

	sb.WriteString("\nRegisters: \n\n")

	sb.WriteString(fmt.Sprintf("PC: %04x \n\n", cpu.PC))

	sb.WriteString("       ")
	for i := range cpu.V {
		sb.WriteString(fmt.Sprintf("V%01x ", i))
	}
	sb.WriteString("\n")
	sb.WriteString("Vx:    ")
	for _, val := range cpu.V {
		sb.WriteString(fmt.Sprintf("%02x ", val))
	}
	sb.WriteString("\n\n")

	sb.WriteString("       ")
	for i := range cpu.Stack {
		sb.WriteString(fmt.Sprintf("%02x   ", i))
	}
	sb.WriteString("\n")
	sb.WriteString("Stack: ")
	for _, val := range cpu.Stack {
		sb.WriteString(fmt.Sprintf("%04x ", val))
	}
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("I:  %04x \n", cpu.I))
	sb.WriteString(fmt.Sprintf("SP: %02x \n", cpu.SP))
	sb.WriteString(fmt.Sprintf("DT: %02x \n", cpu.DT))
	sb.WriteString(fmt.Sprintf("ST: %02x \n", cpu.ST))

	sb.WriteString(cpu.Memory.GetPrettyMemoryState())

	return sb.String()
}
