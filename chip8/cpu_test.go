package chip8

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/tiendc/go-deepcopy"
)

func getRandomCpuState() *Cpu {
	cpu := NewCpu()

	cpu.I = uint16(rand.Intn(0x10000))
	cpu.PC = uint16(rand.Intn(MemorySize - 1)) // Last even memory address
	cpu.SP = uint8(rand.Intn(StackSize))
	cpu.DT = uint8(rand.Intn(StackSize))
	cpu.ST = uint8(rand.Intn(StackSize))

	for i := range 0x10 {
		cpu.V[i] = uint8(rand.Intn(0x100))
	}

	for i := range StackSize {
		cpu.Stack[i] = uint16(rand.Intn(MemorySize - 1)) // Last even memory address
	}

	for i := range MemorySize {
		cpu.Memory.Set8(uint16(i), uint8(rand.Intn(0x100))) // This also randomizes 'interpreter space', containing default sprites
	}

	for y := range height {
		for x := range width {
			cpu.Display.Set(uint(x), uint(y), rand.Intn(2) == 1)
		}
	}

	return cpu
}

type opcodeTest struct {
	inputCpuState *Cpu
	wantCpuState  *Cpu
}

func (test opcodeTest) doOpcodeTest() error {
	cpu := test.inputCpuState

	err := cpu.Tick()
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(cpu, test.wantCpuState) {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(test.wantCpuState.GetPrettyCpuState(), cpu.GetPrettyCpuState(), true)
		return fmt.Errorf("opcode test failed: %v", dmp.DiffPrettyText(dmp.DiffCleanupSemantic(diffs)))
	}

	return nil
}

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
			t.Errorf("RET failed: expected PC %04X, got %04X", pcWant, cpu.PC)
		}

		if cpu.SP != spWant {
			t.Errorf("RET failed: expected SP %02X, got %02X", spWant, cpu.SP)
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
			t.Errorf("RET failed: expected PC %04X, got %04X", pcWant, cpu.PC)
		}

		if cpu.SP != spWant {
			t.Errorf("RET failed: expected SP %02X, got %02X", spWant, cpu.SP)
		}
	})

	t.Run("stack overflow", func(t *testing.T) {
		sp := uint8(0x11)

		cpu := NewCpu()

		cpu.Memory.Set16(0x200, 0x00EE)
		cpu.SP = sp

		err := cpu.Tick()
		if err == nil {
			t.Errorf("RET failed: expected overflow error: sp = %04X", cpu.SP)
		}
	})

	t.Run("stack underflow", func(t *testing.T) {
		sp := uint8(0x00)

		cpu := NewCpu()

		cpu.Memory.Set16(0x200, 0x00EE)
		cpu.SP = sp

		err := cpu.Tick()
		if err == nil {
			t.Errorf("RET failed: expected underflow error: sp = %04X", cpu.SP)
		}
	})

	t.Run("random valid", func(t *testing.T) {
		const n_tests = 200

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x00EE

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)
			inputCpuState.SP = uint8(rand.Intn(0x0F) + 1)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = inputCpuState.Stack[inputCpuState.SP-1]
			wantCpuState.SP = inputCpuState.SP - 1

			t.Run(fmt.Sprintf("JP %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				}.doOpcodeTest()

				if err != nil && inputCpuState.SP > 0x0 && inputCpuState.SP <= 0x10 {
					t.Error(err.Error())
				}
			})
		}
	})

	t.Run("minimum valid", func(t *testing.T) {
		const n_tests = 20

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x00EE

			inputCpuState.SP = uint8(0x1)
			inputCpuState.Stack[0x0] = 0x0000
			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = 0x0000
			wantCpuState.SP = inputCpuState.SP - 1

			t.Run(fmt.Sprintf("JP %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})

	t.Run("maximum valid", func(t *testing.T) {
		const n_tests = 20

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x00EE

			inputCpuState.SP = uint8(0x10)
			inputCpuState.Stack[0xF] = 0x0FFF
			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = 0xFFF
			wantCpuState.SP = inputCpuState.SP - 1

			t.Run(fmt.Sprintf("JP %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}

func TestJP(t *testing.T) {
	t.Run("random valid", func(t *testing.T) {
		const n_tests = 20

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			opcode := uint16(rand.Intn(0x1000)) | 0x1000

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = opcode & 0xFFF

			t.Run(fmt.Sprintf("JP %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})

	t.Run("minimum valid", func(t *testing.T) {
		const n_tests = 20

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x1000

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = 0x0000

			t.Run(fmt.Sprintf("JP %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})

	t.Run("maximum valid", func(t *testing.T) {
		const n_tests = 20

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x1FFE

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = 0x0FFE

			t.Run(fmt.Sprintf("JP %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}

func TestCALL(t *testing.T) {

	t.Run("random valid", func(t *testing.T) {
		const n_tests = 100

		for i := 0; i < n_tests; i++ {
			opcode := uint16(rand.Intn(0x1000)) | 0x2000

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = opcode & 0xFFF
			wantCpuState.SP = inputCpuState.SP + 1
			wantCpuState.Stack[inputCpuState.SP] = inputCpuState.PC + 2

			t.Run(fmt.Sprintf("CALL %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})

	t.Run("minimum valid", func(t *testing.T) {
		const n_tests = 20

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x2000

			inputCpuState.PC = 0x200
			inputCpuState.SP = 0x00
			inputCpuState.Memory.Set16(0x200, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = opcode & 0xFFF
			wantCpuState.SP = inputCpuState.SP + 1
			wantCpuState.Stack[inputCpuState.SP] = inputCpuState.PC + 2

			t.Run(fmt.Sprintf("CALL %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})

	t.Run("maximum valid", func(t *testing.T) {
		const n_tests = 20

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x2FFE

			inputCpuState.PC = 0x200
			inputCpuState.SP = 0x0F
			inputCpuState.Memory.Set16(0x200, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = opcode & 0xFFF
			wantCpuState.SP = inputCpuState.SP + 1
			wantCpuState.Stack[inputCpuState.SP] = inputCpuState.PC + 2

			t.Run(fmt.Sprintf("CALL %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}
