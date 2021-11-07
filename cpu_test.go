package vigo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	stdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)

	t.Run("0000 halt", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		err := cpu.exec(0, 0, 0, 0)
		assert.NoError(t, err)
		assert.True(t, cpu.halt)
	})

	t.Run("00E0 clear screen", func(t *testing.T) {
		var actualDisplay Display
		cpu := NewCPU(&actualDisplay, &[16]bool{})
		err := cpu.exec(0x0, 0x0, 0xE, 0x0)
		assert.NoError(t, err)
		assert.Equal(t, Display{}, actualDisplay)
	})

	t.Run("00EE return", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		curPc := cpu.pc
		cpu.exec(0x2, 0x3, 0x4, 0x0)
		err := cpu.exec(0x0, 0x0, 0xE, 0xE)
		assert.NoError(t, err)
		assert.Equal(t, curPc, cpu.pc)
	})

	t.Run("1NNN jump", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		err := cpu.exec(0x1, 0xA, 0xB, 0xC)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0xABC), cpu.pc)
	})

	t.Run("2NNN call subroutine", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		err := cpu.exec(0x2, 0x3, 0x4, 0x5)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0x345), cpu.pc)
	})

	t.Run("6XNN load register", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		err := cpu.exec(0x6, 0x0, 0xB, 0xC)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0xBC), cpu.reg[0x0])
	})

	t.Run("7XNN add val to reg", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		_ = cpu.exec(0x6, 0x0, 0x1, 0x0)
		err := cpu.exec(0x7, 0x0, 0x0, 0x5)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0x15), cpu.reg[0x0])
	})

	t.Run("AXXX set index register", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		err := cpu.exec(0xA, 0x1, 0x2, 0x3)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0x123), cpu.i)
	})

	t.Run("3XNN skip vx = nn", func(t *testing.T) {
		t.Run("skip", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0xAB
			startPc := cpu.pc
			err := cpu.exec(0x3, 0x1, 0xA, 0xB)
			assert.NoError(t, err)
			assert.Equal(t, startPc+2, cpu.pc)
		})
		t.Run("not skip", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			startPc := cpu.pc
			err := cpu.exec(0x3, 0x1, 0xA, 0xB)
			assert.NoError(t, err)
			assert.Equal(t, startPc, cpu.pc)
		})
	})

	t.Run("4XNN skip vx != nn", func(t *testing.T) {
		t.Run("skip", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			startPc := cpu.pc
			err := cpu.exec(0x4, 0x1, 0xA, 0xB)
			assert.NoError(t, err)
			assert.Equal(t, startPc+2, cpu.pc)
		})
		t.Run("not skip", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			startPc := cpu.pc
			cpu.reg[0x1] = 0xAB
			err := cpu.exec(0x4, 0x1, 0xA, 0xB)
			assert.NoError(t, err)
			assert.Equal(t, startPc, cpu.pc)
		})
	})

	t.Run("5XY0 skip vx = vy", func(t *testing.T) {
		t.Run("skip", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0xAB
			cpu.reg[0x2] = 0xAB
			startPc := cpu.pc
			err := cpu.exec(0x5, 0x1, 0x2, 0x0)
			assert.NoError(t, err)
			assert.Equal(t, startPc+2, cpu.pc)
		})
		t.Run("not skip", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			startPc := cpu.pc
			cpu.reg[0x1] = 0xAB
			err := cpu.exec(0x5, 0x1, 0x2, 0x0)
			assert.NoError(t, err)
			assert.Equal(t, startPc, cpu.pc)
		})
	})

	t.Run("8XY0 set vx=vy", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		cpu.reg[0x2] = 0xAB
		err := cpu.exec16(0x8120)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0xAB), cpu.reg[0x1])
	})

	t.Run("8XY1 set vx=vx|vy", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		cpu.reg[0x1] = 0x12
		cpu.reg[0x2] = 0xAB
		err := cpu.exec16(0x8121)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0xAB|0x12), cpu.reg[0x1])
	})

	t.Run("8XY2 set vx=vx&vy", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		cpu.reg[0x1] = 0x12
		cpu.reg[0x2] = 0xAB
		err := cpu.exec16(0x8122)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0xAB&0x12), cpu.reg[0x1])
	})

	t.Run("8XY3 set vx=vx^vy", func(t *testing.T) {
		cpu := NewCPU(&Display{}, &([16]bool{}))
		cpu.reg[0x1] = 0x12
		cpu.reg[0x2] = 0xAB
		err := cpu.exec16(0x8123)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0xAB^0x12), cpu.reg[0x1])
	})

	t.Run("8XY4 vx+=vy with carry in vf", func(t *testing.T) {
		t.Run("no carry", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0x12
			cpu.reg[0x2] = 0x34
			err := cpu.exec16(0x8124)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0x12+0x34), cpu.reg[0x1])
			assert.Equal(t, uint8(0), cpu.reg[0xF])
		})

		t.Run("with carry", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0xFF
			cpu.reg[0x2] = 0x01
			err := cpu.exec16(0x8124)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0x00), cpu.reg[0x1])
			assert.Equal(t, uint8(1), cpu.reg[0xF])
		})
	})

	t.Run("8XY5 vx-=vy with *flow in vf", func(t *testing.T) {
		t.Run("no carry", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0x34
			cpu.reg[0x2] = 0x12
			err := cpu.exec16(0x8125)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0x34-0x12), cpu.reg[0x1])
			assert.Equal(t, uint8(1), cpu.reg[0xF])
		})

		t.Run("with carry", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0x00
			cpu.reg[0x2] = 0x01
			err := cpu.exec16(0x8125)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0xFF), cpu.reg[0x1])
			assert.Equal(t, uint8(0), cpu.reg[0xF])
		})
	})

	t.Run("8XY7 vx=vy-vx with *flow in vf", func(t *testing.T) {
		t.Run("no carry", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0x12
			cpu.reg[0x2] = 0x34
			err := cpu.exec16(0x8127)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0x34-0x12), cpu.reg[0x1])
			assert.Equal(t, uint8(1), cpu.reg[0xF])
		})

		t.Run("with carry", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0x01
			cpu.reg[0x2] = 0x00
			err := cpu.exec16(0x8127)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0xFF), cpu.reg[0x1])
			assert.Equal(t, uint8(0), cpu.reg[0xF])
		})
	})

	t.Run("8XY6 vx shr 1 quirk true", func(t *testing.T) {
		t.Run("no shiftout", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0b1010_1010
			err := cpu.exec16(0x8126)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0b0101_0101), cpu.reg[0x1])
			assert.Equal(t, uint8(0), cpu.reg[0xF])
		})

		t.Run("shift out", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0b0101_0101
			err := cpu.exec16(0x8126)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0b0010_1010), cpu.reg[0x1])
			assert.Equal(t, uint8(1), cpu.reg[0xF])
		})
	})

	t.Run("8XYE vx shr 1 quirk true", func(t *testing.T) {
		t.Run("shiftout", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0b1010_1010
			err := cpu.exec16(0x812E)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0b0101_0100), cpu.reg[0x1])
			assert.Equal(t, uint8(1), cpu.reg[0xF])
		})

		t.Run("no shiftout", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0b0101_0101
			err := cpu.exec16(0x812E)
			assert.NoError(t, err)
			assert.Equal(t, uint8(0b1010_1010), cpu.reg[0x1])
			assert.Equal(t, uint8(0), cpu.reg[0xF])
		})
	})

	t.Run("9XY0 skip vx != vy", func(t *testing.T) {
		t.Run("skip", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			cpu.reg[0x1] = 0xAB
			startPc := cpu.pc
			err := cpu.exec(0x9, 0x1, 0x2, 0x0)
			assert.NoError(t, err)
			assert.Equal(t, startPc+2, cpu.pc)
		})
		t.Run("not skip", func(t *testing.T) {
			cpu := NewCPU(&Display{}, &([16]bool{}))
			startPc := cpu.pc
			cpu.reg[0x1] = 0xAB
			cpu.reg[0x2] = 0xAB
			err := cpu.exec(0x9, 0x1, 0x2, 0x0)
			assert.NoError(t, err)
			assert.Equal(t, startPc, cpu.pc)
		})
	})

	os.Stdout = stdout
}

func (cpu *CPU) exec16(i uint16) error {
	a1, a2 := splitUint16(i)
	b1, b2 := splitUint8(a1)
	b3, b4 := splitUint8(a2)
	return cpu.exec(b1, b2, b3, b4)
}
