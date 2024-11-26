package chip8

import (
	"math/rand"
	"testing"
)

func TestMemory_Get(t *testing.T) {
	tests := map[string]struct {
		memoryPokes map[uint16]uint8
		getAddr     uint16
		want        uint8
		wantErr     bool
	}{
		"default": {
			getAddr: 0x000,
			want:    0xF0,
			wantErr: false,
		},
		"modified byte": {
			memoryPokes: map[uint16]uint8{
				0x800: 0x42,
			},
			getAddr: 0x800,
			want:    0x42,
			wantErr: false,
		},
		"out of range": {
			getAddr: 0x1000,
			want:    0x00,
			wantErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mem := NewMemory()

			for addr, val := range test.memoryPokes {
				mem.Memory[addr] = val
			}

			got, err := mem.Get(test.getAddr)
			if (err != nil) != test.wantErr {
				t.Errorf("Memory.Get() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("Memory.Get() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestMemory_Set(t *testing.T) {
	tests := map[string]struct {
		memoryPokes map[uint16]uint8
		setAddr     uint16
		setVal      uint8
		wantErr     bool
	}{
		"default": {
			setAddr: 0x800,
			setVal:  0x42,
			wantErr: false,
		},
		"overwrite": {
			memoryPokes: map[uint16]uint8{
				0x800: 0x42,
			},
			setAddr: 0x800,
			setVal:  0x69,
			wantErr: false,
		},
		"out of range": {
			setAddr: 0x1000,
			setVal:  0x69,
			wantErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mem := NewMemory()

			for addr, val := range test.memoryPokes {
				mem.Memory[addr] = val
			}

			err := mem.Set(test.setAddr, test.setVal)

			if err != nil {
				if !test.wantErr {
					t.Errorf("Memory.Set() error = %v, wantErr %v", err, test.wantErr)
				} else {
					return
				}
			}

			if (err == nil) && test.wantErr {
				t.Errorf("Memory.Set() did not throw error as wanted")
				return
			}

			got := mem.Memory[test.setAddr]

			if got != test.setVal {
				t.Errorf("Memory.Set() wrote a %v, want %v", got, test.setVal)
			}
		})
	}
}

func TestGetPrettyMemoryState(t *testing.T) {
	tests := map[string]struct {
		memoryPokes map[uint16]uint8
		wantErr     bool
	}{
		"default": {
			wantErr: false,
		},
	}

	for name, test := range tests {
		mem := NewMemory()

		t.Run(name, func(t *testing.T) {
			for addr, val := range test.memoryPokes {
				mem.Memory[addr] = val
			}

			output := mem.GetPrettyMemoryState()
			t.Log(output)
		})
	}
}

func BenchmarkGetPrettyMemoryState(b *testing.B) {
	mem := NewMemory()

	for i := range mem.Memory {
		mem.Memory[i] = uint8(rand.Intn(0xFF))
	}

	b.ResetTimer()
	mem.GetPrettyMemoryState()
}
