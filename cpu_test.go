package vigo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
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

	t.Run("1NNN jump", func(t *testing.T) {
		cpu := NewCPU()
		err := cpu.exec(0x1, 0xA, 0xB, 0xC)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0xABC), cpu.pc)
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
}
