package vigo

import (
	"fmt"
	"time"
)

type Display [64][32]bool
type displayCallback func(Display)

type CPU struct {
	halt bool

	memory [0x1000]uint8

	pc  uint16
	i   uint16
	reg [0x10]uint8

	delayTimer uint8
	soundTimer uint8

	display         Display
	displayCallback displayCallback

	// number if instruction per second (Hz)
	speed int64
	delay time.Duration

	callStack stack
}

func NewCPU() *CPU {
	cpu := &CPU{
		pc:    0x200,
		speed: 500,
	}
	cpu.delay = time.Duration(1000000/cpu.speed) * time.Microsecond
	return cpu
}

func (cpu *CPU) SetDisplayCallback(dc func(Display)) {
	cpu.displayCallback = dc
}

func (cpu *CPU) WriteInst(addr uint16, inst uint16) {
	inst1, inst2 := splitUint16(addr)
	cpu.memory[addr] = inst1
	cpu.memory[addr+1] = inst2
}

func (cpu *CPU) LoadRom(data []byte) {
	for i, v := range data {
		cpu.memory[0x200+i] = v
	}
}

func (cpu *CPU) Run() {
	for !cpu.halt {
		cpu.Step()
		time.Sleep(cpu.delay)
	}
}

func (cpu *CPU) Step() {
	fmt.Printf("PC=%#X\t", cpu.pc)

	i1, i2, i3, i4 := cpu.fetch()

	err := cpu.exec(i1, i2, i3, i4)
	if err != nil {
		fmt.Print("ERROR ", err, " HALT")
		cpu.halt = true
	}
	fmt.Print("\n")
}

func (cpu *CPU) fetch() (uint8, uint8, uint8, uint8) {
	x1 := cpu.memory[cpu.pc]
	x2 := cpu.memory[cpu.pc+1]
	cpu.pc += 2
	y1, y2 := splitUint8(x1)
	y3, y4 := splitUint8(x2)
	return y1, y2, y3, y4
}

func (cpu *CPU) exec(i1, i2, i3, i4 uint8) error {
	fmt.Printf("0x%X%X%X%X\t%s\t", i1, i2, i3, i4, PrintInst(i1, i2, i3, i4))

	if i1 == 0 && i2 == 0 && i3 == 0 && i4 == 0 {
		cpu.halt = true
		return nil
	}

	if i1 == 0 && i2 == 0 && i3 == 0xE && i4 == 0 {
		return cpu.opClearDisplay()
	}

	if i1 == 0 && i2 == 0 && i3 == 0xE && i4 == 0xE {
		return cpu.opReturn()
	}

	if i1 == 0x1 {
		addr := combine3u4(i2, i3, i4)
		return cpu.opJump(addr)
	}

	if i1 == 0x2 {
		addr := combine3u4(i2, i3, i4)
		return cpu.opCallSub(addr)
	}

	if i1 == 0x3 {
		reg := i2
		val := combine2u4(i3, i4)
		return cpu.opSkipRegEqLit(reg, val)
	}

	if i1 == 0x4 {
		reg := i2
		val := combine2u4(i3, i4)
		return cpu.opSkipRegNeqLit(reg, val)
	}

	if i1 == 0x5 {
		reg1 := i2
		reg2 := i3
		return cpu.opSkipRegEqReq(reg1, reg2)
	}

	if i1 == 0x6 {
		reg := i2
		val := combine2u4(i3, i4)
		return cpu.opLoadReg(reg, val)
	}

	if i1 == 0x7 {
		reg := i2
		val := combine2u4(i3, i4)
		return cpu.opAddReg(reg, val)
	}

	if i1 == 0x8 && i4 == 0x0 {
		reg1 := i2
		reg2 := i3
		return cpu.opLoadRegReg(reg1, reg2)
	}

	if i1 == 0x9 {
		reg1 := i2
		reg2 := i3
		return cpu.opSkipRegNeqReq(reg1, reg2)
	}

	if i1 == 0xA {
		addr := combine3u4(i2, i3, i4)
		return cpu.opLoadI(addr)
	}

	if i1 == 0xD {
		x, y, n := i2, i3, i4
		return cpu.opDraw(x, y, n)
	}

	return fmt.Errorf("INVALID\t")
}