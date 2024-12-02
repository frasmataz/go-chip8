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
	wantError     bool
}

func (test opcodeTest) doOpcodeTest() error {
	cpu := test.inputCpuState

	err := cpu.Tick()
	if err != nil {
		if !test.wantError {
			return err
		}
		return nil
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
	t.Run("random state", func(t *testing.T) {
		const n_tests = 200

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x00EE

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantError := false
			wantCpuState := new(Cpu)

			if inputCpuState.SP > 0x00 && inputCpuState.SP <= 0x10 {
				_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

				wantCpuState.PC = inputCpuState.Stack[inputCpuState.SP-1]
				wantCpuState.SP = inputCpuState.SP - 1
			} else {
				wantError = true
			}

			t.Run(fmt.Sprintf("JP %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     wantError,
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
					wantError:     false,
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
					wantError:     false,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}

func TestJP(t *testing.T) {
	t.Run("random state", func(t *testing.T) {
		const n_tests = 200

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			target := uint16(rand.Intn(0x1000))
			opcode := target | 0x1000
			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			wantError := false

			if target <= 0xFFE {
				_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

				wantCpuState.PC = opcode & 0xFFF
			} else {
				wantError = true
			}

			t.Run(fmt.Sprintf("JP %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     wantError,
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
					wantError:     false,
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
					wantError:     false,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})

	t.Run("out of range", func(t *testing.T) {
		const n_tests = 20

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x1FFF

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)

			t.Run(fmt.Sprintf("JP %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     true,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}

func TestCALL(t *testing.T) {
	t.Run("random state", func(t *testing.T) {
		const n_tests = 200

		for i := 0; i < n_tests; i++ {
			opcode := uint16(rand.Intn(0x1000)) | 0x2000

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantError := false
			wantCpuState := new(Cpu)

			if inputCpuState.SP < 0x10 {
				_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

				wantCpuState.PC = opcode & 0xFFF
				wantCpuState.SP = inputCpuState.SP + 1
				wantCpuState.Stack[inputCpuState.SP] = inputCpuState.PC + 2

			} else {
				wantError = true
			}

			t.Run(fmt.Sprintf("CALL %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     wantError,
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
					wantError:     false,
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
					wantError:     false,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})

	t.Run("stack overflow", func(t *testing.T) {
		const n_tests = 20

		for i := 0; i < n_tests; i++ {
			inputCpuState := getRandomCpuState()

			const opcode = 0x2FFE

			inputCpuState.PC = 0x200
			inputCpuState.SP = 0x10
			inputCpuState.Memory.Set16(0x200, opcode)

			wantCpuState := new(Cpu)

			t.Run(fmt.Sprintf("CALL %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     true,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}

func TestSE_v_byte(t *testing.T) {
	t.Run("random state", func(t *testing.T) {
		const n_tests = 200

		for i := 0; i < n_tests; i++ {
			r := uint16(rand.Intn(0x10))
			kk := uint8(rand.Intn(0x100))

			opcode := (0x3000 | r<<8 | uint16(kk))

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			if inputCpuState.V[r] == kk {
				wantCpuState.PC = inputCpuState.PC + 4
			} else {
				wantCpuState.PC = inputCpuState.PC + 2
			}

			t.Run(fmt.Sprintf("SE_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x3000)

			inputCpuState := getRandomCpuState()
			inputCpuState.V[0x0] = 0x00

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = inputCpuState.PC + 4

			t.Run(fmt.Sprintf("SE_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x3FFF)

			inputCpuState := getRandomCpuState()
			inputCpuState.V[0xF] = 0xFF

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = inputCpuState.PC + 4

			t.Run(fmt.Sprintf("SE_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}

func TestSNE_v_byte(t *testing.T) {
	t.Run("random state", func(t *testing.T) {
		const n_tests = 200

		for i := 0; i < n_tests; i++ {
			r := uint16(rand.Intn(0x10))
			kk := uint8(rand.Intn(0x100))

			opcode := (0x4000 | r<<8 | uint16(kk))

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			if inputCpuState.V[r] != kk {
				wantCpuState.PC = inputCpuState.PC + 4
			} else {
				wantCpuState.PC = inputCpuState.PC + 2
			}

			t.Run(fmt.Sprintf("SNE_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x4000)

			inputCpuState := getRandomCpuState()
			inputCpuState.V[0x0] = 0x01

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = inputCpuState.PC + 4

			t.Run(fmt.Sprintf("SNE_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x4FFF)

			inputCpuState := getRandomCpuState()
			inputCpuState.V[0xF] = 0xFE

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = inputCpuState.PC + 4

			t.Run(fmt.Sprintf("SNE_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}

func TestSE_v1_v2(t *testing.T) {
	t.Run("random state", func(t *testing.T) {
		const n_tests = 200

		for i := 0; i < n_tests; i++ {
			r1 := uint16(rand.Intn(0x10))
			r2 := uint16(rand.Intn(0x10))

			opcode := (0x5000 | r1<<8 | r2<<4)

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			if inputCpuState.V[r1] == inputCpuState.V[r2] {
				wantCpuState.PC = inputCpuState.PC + 4
			} else {
				wantCpuState.PC = inputCpuState.PC + 2
			}

			t.Run(fmt.Sprintf("SE_v1_v2 %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x5000)

			inputCpuState := getRandomCpuState()
			inputCpuState.V[0x0] = 0x00

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = inputCpuState.PC + 4

			t.Run(fmt.Sprintf("SE_v1_v2 %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x5FF0)

			inputCpuState := getRandomCpuState()
			inputCpuState.V[0xF] = 0xFF

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.PC = inputCpuState.PC + 4

			t.Run(fmt.Sprintf("SE_v1_v2 %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}

func TestLD_v_byte(t *testing.T) {
	t.Run("random state", func(t *testing.T) {
		const n_tests = 200

		for i := 0; i < n_tests; i++ {
			r := uint16(rand.Intn(0x10))
			kk := uint8(rand.Intn(0x100))

			opcode := (0x6000 | r<<8 | uint16(kk))

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.V[r] = kk
			wantCpuState.PC = inputCpuState.PC + 2

			t.Run(fmt.Sprintf("LD_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x6000)

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.V[0x0] = 0x00
			wantCpuState.PC = inputCpuState.PC + 2

			t.Run(fmt.Sprintf("LD_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x6FFF)

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.V[0xF] = 0xFF
			wantCpuState.PC = inputCpuState.PC + 2

			t.Run(fmt.Sprintf("LD_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}

func TestADD_v_byte(t *testing.T) {
	t.Run("random state", func(t *testing.T) {
		const n_tests = 200

		for i := 0; i < n_tests; i++ {
			r := uint16(rand.Intn(0x10))
			kk := uint8(rand.Intn(0x100))

			opcode := (0x7000 | r<<8 | uint16(kk))

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.V[r] = inputCpuState.V[r] + kk
			wantCpuState.PC = inputCpuState.PC + 2

			t.Run(fmt.Sprintf("LD_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x7010)

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)
			inputCpuState.V[0x0] = 0x00

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.V[0x0] = 0x10
			wantCpuState.PC = inputCpuState.PC + 2

			t.Run(fmt.Sprintf("LD_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
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
			opcode := uint16(0x7F10)

			inputCpuState := getRandomCpuState()

			inputCpuState.Memory.Set16(inputCpuState.PC, opcode)
			inputCpuState.V[0xF] = 0xFF

			wantCpuState := new(Cpu)
			_ = deepcopy.Copy(&wantCpuState, &inputCpuState)

			wantCpuState.V[0xF] = 0x0F
			wantCpuState.PC = inputCpuState.PC + 2

			t.Run(fmt.Sprintf("LD_v_byte %04X", opcode), func(t *testing.T) {
				err := opcodeTest{
					inputCpuState: inputCpuState,
					wantCpuState:  wantCpuState,
					wantError:     false,
				}.doOpcodeTest()

				if err != nil {
					t.Error(err.Error())
				}
			})
		}
	})
}
