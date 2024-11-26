package chip8

import (
	"testing"
)

func TestSet(t *testing.T) {
	type displayPoke struct {
		x   uint
		y   uint
		val bool
	}
	tests := map[string]struct {
		displayPokes []displayPoke
		setx         uint
		sety         uint
		val          bool
		wantErr      bool
	}{
		"default": {
			setx:    42,
			sety:    10,
			val:     true,
			wantErr: false,
		},
		"overwrite": {
			displayPokes: []displayPoke{
				{
					x:   36,
					y:   20,
					val: true,
				},
			},
			setx:    36,
			sety:    20,
			val:     false,
			wantErr: false,
		},
		"out of range": {
			setx:    64,
			sety:    32,
			val:     true,
			wantErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			display := NewDisplay()

			for _, poke := range test.displayPokes {
				display.framebuffer[poke.y][poke.x] = poke.val
			}

			err := display.Set(test.setx, test.sety, test.val)

			if err != nil {
				if !test.wantErr {
					t.Errorf("Display.Set() error = %v, wantErr %v", err, test.wantErr)
				} else {
					return
				}
			}

			if (err == nil) && test.wantErr {
				t.Errorf("Display.Set() did not throw error as wanted")
				return
			}

			got := display.framebuffer[test.sety][test.setx]

			if got != test.val {
				t.Errorf("Display.Set() wrote a %v, want %v", got, test.val)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type displayPoke struct {
		x   uint
		y   uint
		val bool
	}
	tests := map[string]struct {
		displayPokes []displayPoke
		getx         uint
		gety         uint
		want         bool
		wantErr      bool
	}{
		"default": {
			getx:    42,
			gety:    10,
			want:    false,
			wantErr: false,
		},
		"read set pixel": {
			displayPokes: []displayPoke{
				{
					x:   36,
					y:   20,
					val: true,
				},
			},
			getx:    36,
			gety:    20,
			want:    true,
			wantErr: false,
		},
		"out of range": {
			getx:    64,
			gety:    32,
			want:    false,
			wantErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			display := NewDisplay()

			for _, poke := range test.displayPokes {
				display.framebuffer[poke.y][poke.x] = poke.val
			}

			got, err := display.Get(test.getx, test.gety)

			if err != nil {
				if !test.wantErr {
					t.Errorf("Display.Set() error = %v, wantErr %v", err, test.wantErr)
				} else {
					return
				}
			}

			if (err == nil) && test.wantErr {
				t.Errorf("Display.Set() did not throw error as wanted")
				return
			}

			if got != test.want {
				t.Errorf("Display.Set() wrote a %v, want %v", got, test.want)
			}
		})
	}
}

func TestPrintFrame(t *testing.T) {
	type displayPoke struct {
		x   uint
		y   uint
		val bool
	}
	tests := map[string]struct {
		displayPokes []displayPoke
		wantErr      bool
	}{
		"test image :)": {
			displayPokes: []displayPoke{
				{
					x:   28,
					y:   12,
					val: true,
				},
				{
					x:   36,
					y:   12,
					val: true,
				},
				{
					x:   28,
					y:   19,
					val: true,
				},
				{
					x:   28,
					y:   20,
					val: true,
				},
				{
					x:   29,
					y:   20,
					val: true,
				},
				{
					x:   30,
					y:   20,
					val: true,
				},
				{
					x:   31,
					y:   20,
					val: true,
				},
				{
					x:   32,
					y:   20,
					val: true,
				},
				{
					x:   33,
					y:   20,
					val: true,
				},
				{
					x:   34,
					y:   20,
					val: true,
				},
				{
					x:   35,
					y:   20,
					val: true,
				},
				{
					x:   36,
					y:   20,
					val: true,
				},
				{
					x:   36,
					y:   19,
					val: true,
				},
			},
			wantErr: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			display := NewDisplay()

			for _, poke := range test.displayPokes {
				display.framebuffer[poke.y][poke.x] = poke.val
			}

			output := display.PrintFrame()
			t.Log(output)
		})
	}
}
