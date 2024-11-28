package chip8

import (
	"testing"
)

func TestMemory_Get8(t *testing.T) {
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

			got, err := mem.Get8(test.getAddr)
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

func TestMemory_Get16(t *testing.T) {
	tests := map[string]struct {
		memoryPokes map[uint16]uint8
		getAddr     uint16
		want        uint16
		wantErr     bool
	}{
		"16-bit get": {
			memoryPokes: map[uint16]uint8{
				0x800: 0x42,
				0x801: 0x69,
			},
			getAddr: 0x800,
			want:    0x4269,
			wantErr: false,
		},
		"out of range": {
			getAddr: 0xFFF,
			want:    0x0000,
			wantErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mem := NewMemory()

			for addr, val := range test.memoryPokes {
				mem.Memory[addr] = val
			}

			got, err := mem.Get16(test.getAddr)
			if (err != nil) != test.wantErr {
				t.Errorf("Memory.Get() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("Memory.Get() = 0x%02X, want 0x%02X", got, test.want)
			}
		})
	}
}

func TestMemory_Set8(t *testing.T) {
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

			err := mem.Set8(test.setAddr, test.setVal)

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

func TestMemory_Set16(t *testing.T) {
	tests := map[string]struct {
		memoryPokes map[uint16]uint8
		setAddr     uint16
		setVal      uint16
		wantErr     bool
	}{
		"16-bit set": {
			setAddr: 0x800,
			setVal:  0x4269,
			wantErr: false,
		},
		"out of range": {
			setAddr: 0xFFF,
			setVal:  0x6942,
			wantErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mem := NewMemory()

			for addr, val := range test.memoryPokes {
				mem.Memory[addr] = val
			}

			err := mem.Set16(test.setAddr, test.setVal)

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

			got := uint16(mem.Memory[test.setAddr])<<8 | uint16(mem.Memory[test.setAddr+1])

			if got != test.setVal {
				t.Errorf("Memory.Set() wrote a %v, want %v", got, test.setVal)
			}
		})
	}
}

func TestGetPrettyMemoryState(t *testing.T) {
	NewMemory().GetPrettyMemoryState()
}
