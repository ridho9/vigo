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
		cpu := NewCPU()
		err := cpu.exec(0, 0, 0, 0)
		assert.NoError(t, err)
		assert.True(t, cpu.halt)
	})

	t.Run("00E0 clear screen", func(t *testing.T) {
		var actualDisplay Display
		cpu := NewCPU()
		cpu.SetDisplayCallback(func(d Display) {
			actualDisplay = d
		})
		err := cpu.exec(0x0, 0x0, 0xE, 0x0)
		assert.NoError(t, err)
		assert.Equal(t, Display{}, actualDisplay)
	})

	t.Run("00EE return", func(t *testing.T) {
		cpu := NewCPU()
		curPc := cpu.pc
		cpu.exec(0x2, 0x3, 0x4, 0x0)
		err := cpu.exec(0x0, 0x0, 0xE, 0xE)
		assert.NoError(t, err)
		assert.Equal(t, curPc+2, cpu.pc)
	})

	t.Run("1NNN jump", func(t *testing.T) {
		cpu := NewCPU()
		err := cpu.exec(0x1, 0xA, 0xB, 0xC)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0xABC), cpu.pc)
	})

	t.Run("2NNN call subroutine", func(t *testing.T) {
		cpu := NewCPU()
		err := cpu.exec(0x2, 0x3, 0x4, 0x5)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0x345), cpu.pc)
	})

	t.Run("6XNN load register", func(t *testing.T) {
		cpu := NewCPU()
		err := cpu.exec(0x6, 0x0, 0xB, 0xC)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0xBC), cpu.reg[0x0])
	})

	t.Run("7XNN add val to reg", func(t *testing.T) {
		cpu := NewCPU()
		_ = cpu.exec(0x6, 0x0, 0x1, 0x0)
		err := cpu.exec(0x7, 0x0, 0x0, 0x5)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0x15), cpu.reg[0x0])
	})

	t.Run("AXXX set index register", func(t *testing.T) {
		cpu := NewCPU()
		err := cpu.exec(0xA, 0x1, 0x2, 0x3)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0x123), cpu.i)
	})

	t.Run("3XNN skip vx = nn", func(t *testing.T) {
		t.Run("skip", func(t *testing.T) {
			cpu := NewCPU()
			cpu.reg[0x1] = 0xAB
			startPc := cpu.pc
			err := cpu.exec(0x3, 0x1, 0xA, 0xB)
			assert.NoError(t, err)
			assert.Equal(t, startPc+2, cpu.pc)
		})
		t.Run("not skip", func(t *testing.T) {
			cpu := NewCPU()
			startPc := cpu.pc
			err := cpu.exec(0x3, 0x1, 0xA, 0xB)
			assert.NoError(t, err)
			assert.Equal(t, startPc, cpu.pc)
		})
	})

	t.Run("4XNN skip vx != nn", func(t *testing.T) {
		t.Run("skip", func(t *testing.T) {
			cpu := NewCPU()
			startPc := cpu.pc
			err := cpu.exec(0x4, 0x1, 0xA, 0xB)
			assert.NoError(t, err)
			assert.Equal(t, startPc+2, cpu.pc)
		})
		t.Run("not skip", func(t *testing.T) {
			cpu := NewCPU()
			startPc := cpu.pc
			cpu.reg[0x1] = 0xAB
			err := cpu.exec(0x4, 0x1, 0xA, 0xB)
			assert.NoError(t, err)
			assert.Equal(t, startPc, cpu.pc)
		})
	})

	t.Run("5XY0 skip vx = vy", func(t *testing.T) {
		t.Run("skip", func(t *testing.T) {
			cpu := NewCPU()
			cpu.reg[0x1] = 0xAB
			cpu.reg[0x2] = 0xAB
			startPc := cpu.pc
			err := cpu.exec(0x5, 0x1, 0x2, 0x0)
			assert.NoError(t, err)
			assert.Equal(t, startPc+2, cpu.pc)
		})
		t.Run("not skip", func(t *testing.T) {
			cpu := NewCPU()
			startPc := cpu.pc
			cpu.reg[0x1] = 0xAB
			err := cpu.exec(0x5, 0x1, 0x2, 0x0)
			assert.NoError(t, err)
			assert.Equal(t, startPc, cpu.pc)
		})
	})

	t.Run("8XY0 set vx=vy", func(t *testing.T) {
		cpu := NewCPU()
		cpu.reg[0x2] = 0xAB
		err := cpu.exec16(0x8120)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0xAB), cpu.reg[0x1])
	})

	t.Run("9XY0 skip vx != vy", func(t *testing.T) {
		t.Run("skip", func(t *testing.T) {
			cpu := NewCPU()
			cpu.reg[0x1] = 0xAB
			startPc := cpu.pc
			err := cpu.exec(0x9, 0x1, 0x2, 0x0)
			assert.NoError(t, err)
			assert.Equal(t, startPc+2, cpu.pc)
		})
		t.Run("not skip", func(t *testing.T) {
			cpu := NewCPU()
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
