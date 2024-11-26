package chip8

type Cpu struct {
	V     [0xF]uint8 // General-purpose 8-bit registers V0 - VF
	I     uint16     // 16-bit register, used to hold memory addresses
	PC    uint16     // Program counter - 16-bit
	SP    uint8      // Stack pointer - 8-bit
	DT    uint8      // Delay timer - 8-bit - dec at 60Hz when non-zero
	ST    uint8      // Sound timer - 8-bit - dec at 60Hz when non-zero
	Stack [0xF]uint8 // Stack - 16 16-bit values
}
