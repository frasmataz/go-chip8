package chip8

import (
	"fmt"
	"strings"
)

type Memory struct {
	Memory [0xFFF + 1]uint8
}

func (mem *Memory) GetPrettyMemoryState() string {
	const columns = 16

	// Pre-compute output string size to minimize memory copying while building string
	// 6 chars per row adddress  (⏎xxx:␣) + 3 chars for each byte (xx␣)
	//outputSize := (len(mem.Memory) * 3) + ((len(mem.Memory) / columns) * 6)

	var sb strings.Builder
	sb.Grow(1)

	for addr, val := range mem.Memory {
		if addr%columns == 0 {
			sb.WriteString(fmt.Sprintf("\n%03x: ", addr))
		}

		sb.WriteString(fmt.Sprintf("%02x ", val))
	}

	return sb.String()
}

func (mem *Memory) Get(addr uint16) (uint8, error) {
	if addr > uint16(len(mem.Memory)-1) {
		return 0x0, fmt.Errorf("memory address out of bounds: %v, capacity %v", addr, len(mem.Memory))
	}

	return mem.Memory[addr], nil
}

func (mem *Memory) Set(addr uint16, val uint8) error {
	if addr > uint16(len(mem.Memory)-1) {
		return fmt.Errorf("memory address out of bounds: %v, capacity %v", addr, len(mem.Memory))
	}

	mem.Memory[addr] = val

	return nil
}

func NewMemory() *Memory {
	return new(Memory)
}
