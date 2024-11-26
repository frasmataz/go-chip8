package chip8

import (
	"fmt"
	"strings"
)

type Cpu struct {
	V     [0x10]uint8 // General-purpose 8-bit registers V0 - VF
	I     uint16      // 16-bit register, used to hold memory addresses
	PC    uint16      // Program counter - 16-bit
	SP    uint8       // Stack pointer - 8-bit
	DT    uint8       // Delay timer - 8-bit - dec at 60Hz when non-zero
	ST    uint8       // Sound timer - 8-bit - dec at 60Hz when non-zero
	Stack [0x10]uint8 // Stack - 16 16-bit values
}

func (cpu *Cpu) GetPrettyCpuState() string {
	var sb strings.Builder

	sb.WriteString("\nRegisters: \n\n")

	sb.WriteString(fmt.Sprintf("PC: %02x \n\n", cpu.PC))

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
		sb.WriteString(fmt.Sprintf("%02x ", i))
	}
	sb.WriteString("\n")
	sb.WriteString("Stack: ")
	for _, val := range cpu.Stack {
		sb.WriteString(fmt.Sprintf("%02x ", val))
	}
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("I:  %02x \n", cpu.I))
	sb.WriteString(fmt.Sprintf("SP: %02x \n", cpu.SP))
	sb.WriteString(fmt.Sprintf("DT: %02x \n", cpu.DT))
	sb.WriteString(fmt.Sprintf("ST: %02x \n", cpu.ST))

	return sb.String()
}
