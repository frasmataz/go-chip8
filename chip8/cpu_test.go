package chip8

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestGetPrettyCpuState(t *testing.T) {
	t.Log(NewCpu().GetPrettyCpuState())
}

func TestCLS(t *testing.T) {
	cpu := NewCpu()

	for y, row := range cpu.Display.framebuffer {
		for x := range row {
			cpu.Display.framebuffer[y][x] = rand.Intn(2) == 1
		}
	}

	cpu.Memory.Set16(0x200, 0x00E0)
	cpu.Tick()

	for y, row := range cpu.Display.framebuffer {
		for x := range row {
			if cpu.Display.framebuffer[y][x] {
				t.Errorf("CLS failed: screen not clear: %v", cpu.Display.PrintFrame())
			}
		}
	}
}

func TestRET(t *testing.T) {
	t.Run("stack n=1", func(t *testing.T) {
		pcWant := uint16(0x123)
		sp := uint8(0x01)
		spWant := uint8(0x00)

		cpu := NewCpu()

		cpu.Memory.Set16(0x200, 0x00EE)
		cpu.Stack[0x0] = pcWant
		cpu.SP = sp

		cpu.Tick()

		if cpu.PC != pcWant {
			t.Errorf("RET failed: expected PC %04x, got %04x", pcWant, cpu.PC)
		}

		if cpu.SP != spWant {
			t.Errorf("RET failed: expected SP %02x, got %02x", spWant, cpu.SP)
		}
	})

	t.Run("stack n=16", func(t *testing.T) {

		pcWant := uint16(0x123)
		sp := uint8(0x10)
		spWant := uint8(0x0F)

		cpu := NewCpu()

		cpu.Memory.Set16(0x200, 0x00EE)
		cpu.Stack[0xF] = pcWant
		cpu.SP = sp

		cpu.Tick()

		if cpu.PC != pcWant {
			t.Errorf("RET failed: expected PC %04x, got %04x", pcWant, cpu.PC)
		}

		if cpu.SP != spWant {
			t.Errorf("RET failed: expected SP %02x, got %02x", spWant, cpu.SP)
		}
	})

	t.Run("stack overflow", func(t *testing.T) {
		sp := uint8(0x11)

		cpu := NewCpu()

		cpu.Memory.Set16(0x200, 0x00EE)
		cpu.SP = sp

		err := cpu.Tick()
		if err == nil {
			t.Errorf("RET failed: expected overflow error: sp = %04x", cpu.SP)
		}
	})

	t.Run("stack underflow", func(t *testing.T) {
		sp := uint8(0x00)

		cpu := NewCpu()

		cpu.Memory.Set16(0x200, 0x00EE)
		cpu.SP = sp

		err := cpu.Tick()
		if err == nil {
			t.Errorf("RET failed: expected underflow error: sp = %04x", cpu.SP)
		}
	})
}

func TestJP(t *testing.T) {
	type test struct {
		opcode uint16
		pcWant uint16
	}

	doTest := func(_test test) {
		cpu := NewCpu()

		cpu.Memory.Set16(0x200, _test.opcode)
		err := cpu.Tick()
		if err != nil {
			t.Errorf("JP failed: %v", err)
		}

		if cpu.PC != _test.pcWant {
			t.Errorf("JP failed: expected PC %04x, got %04x", _test.pcWant, cpu.PC)
		}
	}

	t.Run("random valid", func(t *testing.T) {
		const n = 20
		tests := [n]test{}

		for i := 0; i < n; i++ {
			addr := uint16(rand.Intn(0x1000))
			tests[i] = test{
				opcode: addr | 0x1000,
				pcWant: addr,
			}
		}

		for _, _test := range tests {
			t.Run(fmt.Sprintf("JP %04x", _test.opcode), func(t *testing.T) {
				doTest(_test)
			})
		}
	})

	t.Run("minimum valid", func(t *testing.T) {
		doTest(test{
			opcode: 0x1000,
			pcWant: 0x0000,
		})
	})

	t.Run("maximum valid", func(t *testing.T) {
		doTest(test{
			opcode: 0x1FFE,
			pcWant: 0x0FFE,
		})
	})
}

func TestCALL(t *testing.T) {
	type test struct {
		inputCpuState *Cpu
		wantCpuState  *Cpu
	}

	doTest := func(t *testing.T, test test) {
		cpu := test.inputCpuState

		cpu.Tick()
		if !reflect.DeepEqual(cpu, test.wantCpuState) {
			t.Errorf("CALL failed: expected CPU state: %v\n\ngot state: %v", test.wantCpuState.GetPrettyCpuState(), cpu.GetPrettyCpuState())
		}
	}

	t.Run("random valid", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			addr := uint16(rand.Intn(0x1000))
			sp := uint8(rand.Intn(0x0F))
			pc := uint16(rand.Intn(0x1000-0x200) + 0x200)
			opcode := addr | 0x2000

			inputCpuState := &Cpu{
				PC:     pc,
				SP:     sp,
				Memory: NewMemory(),
			}
			inputCpuState.Memory.Set16(pc, opcode)

			wantCpuState := &Cpu{
				PC:     addr,
				SP:     sp + 1,
				Memory: NewMemory(),
			}
			wantCpuState.Memory.Set16(pc, opcode)
			wantCpuState.Stack[sp] = pc + 2

			t.Run(fmt.Sprintf("CALL %04x", opcode), func(t *testing.T) {
				doTest(t, test{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				})
			})
		}
	})

	t.Run("minimum valid", func(t *testing.T) {
		inputCpuState := &Cpu{
			PC:     0x200,
			SP:     0x00,
			Memory: NewMemory(),
		}
		inputCpuState.Memory.Set16(0x200, 0x2000)

		wantCpuState := &Cpu{
			PC:     0x000,
			SP:     0x01,
			Memory: inputCpuState.Memory,
		}
		wantCpuState.Stack[0x0] = 0x202

		doTest(t, test{
			inputCpuState: inputCpuState,
			wantCpuState:  wantCpuState,
		})
	})

	t.Run("maximum valid", func(t *testing.T) {
		inputCpuState := &Cpu{
			PC:     0x200,
			SP:     0x00,
			Memory: NewMemory(),
		}
		inputCpuState.Memory.Set16(0x200, 0x2FFE)

		wantCpuState := &Cpu{
			PC:     0xFFE,
			SP:     0x01,
			Memory: inputCpuState.Memory,
		}
		wantCpuState.Stack[0x0] = 0x202

		doTest(t, test{
			inputCpuState: inputCpuState,
			wantCpuState:  wantCpuState,
		})
	})
}
