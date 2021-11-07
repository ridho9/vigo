package vigo

func (cpu *CPU) opClearDisplay() error {
	for x := 0; x < 64; x += 1 {
		for y := 0; y < 32; y += 1 {
			cpu.display[x][y] = false
		}
	}
	return nil
}

func (cpu *CPU) opJump(addr uint16) error {
	cpu.pc = addr
	return nil
}

func (cpu *CPU) opCallSub(addr uint16) error {
	err := cpu.callStack.push(cpu.pc + 2)
	if err != nil {
		return err
	}
	cpu.pc = addr
	return nil
}

func (cpu *CPU) opReturn() error {
	addr, err := cpu.callStack.pop()
	if err != nil {
		return err
	}
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

func (cpu *CPU) opLoadI(val uint16) error {
	cpu.i = val
	return nil
}

func (cpu *CPU) opDraw(x, y, n uint8) error {
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

	return nil
}

func (cpu *CPU) opSkipRegEqLit(reg uint8, lit uint8) error {
	regv := cpu.reg[reg]
	if regv == lit {
		cpu.pc += 2
	}
	return nil
}

func (cpu *CPU) opSkipRegNeqLit(reg uint8, lit uint8) error {
	regv := cpu.reg[reg]
	if regv != lit {
		cpu.pc += 2
	}
	return nil
}

func (cpu *CPU) opSkipRegEqReq(v1, v2 uint8) error {
	if cpu.reg[v1] == cpu.reg[v2] {
		cpu.pc += 2
	}
	return nil
}

func (cpu *CPU) opSkipRegNeqReq(v1, v2 uint8) error {
	if cpu.reg[v1] != cpu.reg[v2] {
		cpu.pc += 2
	}
	return nil
}

func (cpu *CPU) opLoadRegReg(v1, v2 uint8) error {
	cpu.reg[v1] = cpu.reg[v2]
	return nil
}

func (cpu *CPU) opOrRegReg(v1, v2 uint8) error {
	cpu.reg[v1] |= cpu.reg[v2]
	return nil
}

func (cpu *CPU) opAndRegReg(v1, v2 uint8) error {
	cpu.reg[v1] &= cpu.reg[v2]
	return nil
}

func (cpu *CPU) opXorRegReg(v1, v2 uint8) error {
	cpu.reg[v1] ^= cpu.reg[v2]
	return nil
}

func (cpu *CPU) opAddRegRegO(v1, v2 uint8) error {
	cpu.reg[0xF] = 0
	startV := cpu.reg[v1]
	cpu.reg[v1] += cpu.reg[v2]
	if cpu.reg[v1] < startV {
		cpu.reg[0xF] = 1
	}
	return nil
}

func (cpu *CPU) opSubRegRegO(v1, v2 uint8) error {
	cpu.reg[0xF] = 0
	startV := cpu.reg[v1]
	cpu.reg[v1] -= cpu.reg[v2]
	if cpu.reg[v1] > startV {
		cpu.reg[0xF] = 1
	}
	return nil
}

func (cpu *CPU) opSubbRegRegO(v1, v2 uint8) error {
	cpu.reg[0xF] = 0
	startV := cpu.reg[v2]
	cpu.reg[v1] = cpu.reg[v2] - cpu.reg[v1]
	if cpu.reg[v1] > startV {
		cpu.reg[0xF] = 1
	}
	return nil
}

func (cpu *CPU) opShiftRight(v1, v2 uint8) error {
	if !cpu.quirk.ShiftQuirks {
		cpu.reg[v1] = cpu.reg[v2]
	}
	cpu.reg[0xF] = cpu.reg[v1] & 1
	cpu.reg[v1] >>= 1
	return nil
}

func (cpu *CPU) opShiftLeft(v1, v2 uint8) error {
	if !cpu.quirk.ShiftQuirks {
		cpu.reg[v1] = cpu.reg[v2]
	}
	cpu.reg[0xF] = (cpu.reg[v1] & (1 << 7)) >> 7
	cpu.reg[v1] <<= 1
	return nil
}

func (cpu *CPU) opLoadDelayTimer(v1 uint8) error {
	cpu.reg[v1] = cpu.delayTimer
	return nil
}

func (cpu *CPU) opSetDelayTimer(v1 uint8) error {
	cpu.delayTimer = cpu.reg[v1]
	return nil
}

func (cpu *CPU) opSetSoundTimer(v1 uint8) error {
	cpu.soundTimer = cpu.reg[v1]
	return nil
}

func (cpu *CPU) opStoreIndex(n uint8) error {
	for i := 0; i <= int(n); i += 1 {
		cpu.memory[cpu.i+uint16(i)] = cpu.reg[i]
	}

	if !cpu.quirk.LoadStoreQuirks {
		cpu.i += uint16(n)
	}
	return nil
}

func (cpu *CPU) opLoadIndex(n uint8) error {
	for i := 0; i <= int(n); i += 1 {
		cpu.reg[i] = cpu.memory[cpu.i+uint16(i)]
	}

	if !cpu.quirk.LoadStoreQuirks {
		cpu.i += uint16(n)
	}
	return nil
}

func (cpu *CPU) opBCD(v1 uint8) error {
	v := cpu.reg[v1]
	cpu.memory[cpu.i] = v / 100
	cpu.memory[cpu.i+1] = (v / 10) % 10
	cpu.memory[cpu.i+2] = v % 10
	return nil
}
