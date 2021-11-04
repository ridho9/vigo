package vigo

func (cpu *CPU) opClearDisplay() error {
	cpu.display = Display{}
	if cpu.displayCallback != nil {
		cpu.displayCallback(cpu.display)
	}
	return nil
}

func (cpu *CPU) opJump(addr uint16) error {
	cpu.pc = addr
	return nil
}

func (cpu *CPU) opLoadReg(reg uint8, val uint8) error {
	cpu.reg[reg] = val
	return nil
}

func (cpu *CPU) opAddReg(reg uint8, val uint8) error {
	cpu.reg[reg] += val
	return nil
}

func (cpu *CPU) LoadI(val uint16) error {
	cpu.i = val
	return nil
}

func (cpu *CPU) Draw(x, y, n uint8) error {
	// draw stuff
	startX := cpu.reg[x] % 64
	startY := cpu.reg[y] % 32
	cpu.reg[0xF] = 0

	// for each row
	for a := uint8(0); a < n && (a+startY) < 32; a += 1 {
		curY := a + startY
		bits := u8ToBits(cpu.memory[cpu.i+uint16(a)])

		for b := uint8(0); b < 8 && (b+startX) < 64; b += 1 {
			curX := startX + b
			bit := bits[b]
			screen := cpu.display[curX][curY]
			cpu.display[curX][curY] = screen != bit
			if screen && bit {
				cpu.reg[0xF] = 1
			}
		}
	}

	// callback
	if cpu.displayCallback != nil {
		cpu.displayCallback(cpu.display)
	}
	return nil
}
