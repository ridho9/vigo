package vigo

import (
	"fmt"
	"time"
)

type Display [64][32]bool
type Quirk struct {
	ShiftQuirks     bool
	LoadStoreQuirks bool
}

func DefaultQuirk() Quirk {
	return Quirk{
		ShiftQuirks:     true,
		LoadStoreQuirks: true,
	}
}

const TIMER_DURATION = time.Duration(1000/60) * time.Millisecond

type CPU struct {
	halt bool

	memory [0x1000]uint8

	pc  uint16
	i   uint16
	reg [0x10]uint8

	delayTimer uint8
	soundTimer uint8

	display *Display

	// number of instruction per second (Hz)
	speed int64
	delay time.Duration

	callStack stack
	quirk     Quirk

	keyDown *[0x10]bool
}

func NewCPU(d *Display, kD *[0x10]bool) *CPU {
	cpu := &CPU{
		pc:      0x200,
		speed:   1000,
		display: d,
		quirk:   DefaultQuirk(),
		keyDown: kD,
	}
	cpu.delay = time.Duration(1000/cpu.speed) * time.Millisecond
	cpu.setupFont()
	return cpu
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
	lastTimerUpdate := time.Now()
	for !cpu.halt {
		cpu.Step()

		if time.Since(lastTimerUpdate) >= TIMER_DURATION {
			cpu.TimerStep()
			lastTimerUpdate = time.Now()
		}

		time.Sleep(cpu.delay)
	}
}

func (cpu *CPU) TimerStep() {
	if cpu.delayTimer > 0 {
		cpu.delayTimer -= 1
	}

	if cpu.soundTimer > 0 {
		cpu.soundTimer -= 1
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

	if i1 == 0 && i2 == 0 && i3 == 0xF && i4 == 0xF {
		return nil
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

	if i1 == 0x8 && i4 == 0x1 {
		reg1 := i2
		reg2 := i3
		return cpu.opOrRegReg(reg1, reg2)
	}

	if i1 == 0x8 && i4 == 0x2 {
		reg1 := i2
		reg2 := i3
		return cpu.opAndRegReg(reg1, reg2)
	}

	if i1 == 0x8 && i4 == 0x3 {
		reg1 := i2
		reg2 := i3
		return cpu.opXorRegReg(reg1, reg2)
	}

	if i1 == 0x8 && i4 == 0x4 {
		reg1 := i2
		reg2 := i3
		return cpu.opAddRegRegO(reg1, reg2)
	}

	if i1 == 0x8 && i4 == 0x5 {
		reg1 := i2
		reg2 := i3
		return cpu.opSubRegRegO(reg1, reg2)
	}

	if i1 == 0x8 && i4 == 0x6 {
		reg1 := i2
		reg2 := i3
		return cpu.opShiftRight(reg1, reg2)
	}

	if i1 == 0x8 && i4 == 0x7 {
		reg1 := i2
		reg2 := i3
		return cpu.opSubbRegRegO(reg1, reg2)
	}

	if i1 == 0x8 && i4 == 0xE {
		reg1 := i2
		reg2 := i3
		return cpu.opShiftLeft(reg1, reg2)
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

	if i1 == 0xB {
		reg := i2
		xnn := combine3u4(i2, i3, i4)
		return cpu.opJumpOffset(reg, xnn)
	}

	if i1 == 0xC {
		reg := i2
		val := combine2u4(i3, i4)
		return cpu.opRandom(reg, val)
	}

	if i1 == 0xD {
		x, y, n := i2, i3, i4
		return cpu.opDraw(x, y, n)
	}

	if i1 == 0xF && i3 == 0x0 && i4 == 0x7 {
		reg1 := i2
		return cpu.opLoadDelayTimer(reg1)
	}

	if i1 == 0xF && i3 == 0x1 && i4 == 0x5 {
		reg1 := i2
		return cpu.opSetDelayTimer(reg1)
	}

	if i1 == 0xF && i3 == 0x1 && i4 == 0x8 {
		reg1 := i2
		return cpu.opSetSoundTimer(reg1)
	}

	if i1 == 0xF && i3 == 0x5 && i4 == 0x5 {
		reg1 := i2
		return cpu.opStoreIndex(reg1)
	}

	if i1 == 0xF && i3 == 0x3 && i4 == 0x3 {
		reg1 := i2
		return cpu.opBCD(reg1)
	}

	if i1 == 0xF && i3 == 0x6 && i4 == 0x5 {
		reg1 := i2
		return cpu.opLoadIndex(reg1)
	}

	if i1 == 0xF && i3 == 0x2 && i4 == 0x9 {
		reg1 := i2
		return cpu.opLoadFont(reg1)
	}

	if i1 == 0xF && i3 == 0x0 && i4 == 0xA {
		reg1 := i2
		return cpu.opWaitKeypress(reg1)
	}

	if i1 == 0xF && i3 == 0x1 && i4 == 0xE {
		reg1 := i2
		return cpu.opAddIndex(reg1)
	}

	if i1 == 0xE && i3 == 0x9 && i4 == 0xE {
		reg1 := i2
		return cpu.opSkipIfPressed(reg1)
	}

	if i1 == 0xE && i3 == 0xA && i4 == 0x1 {
		reg1 := i2
		return cpu.opSkipIfNotPressed(reg1)
	}

	return fmt.Errorf("INVALID\t")
}

var FONT_TABLE = [0x10][5]uint8{
	{0xF0, 0x90, 0x90, 0x90, 0xF0}, //0x0
	{0x20, 0x60, 0x20, 0x20, 0x70}, //0x1
	{0xF0, 0x10, 0xF0, 0x80, 0xF0}, //0x2
	{0xF0, 0x10, 0xF0, 0x10, 0xF0}, //0x3
	{0x90, 0x90, 0xF0, 0x10, 0x10}, //0x4
	{0xF0, 0x80, 0xF0, 0x10, 0xF0}, //0x5
	{0xF0, 0x80, 0xF0, 0x90, 0xF0}, //0x6
	{0xF0, 0x10, 0x20, 0x40, 0x40}, //0x7
	{0xF0, 0x90, 0xF0, 0x90, 0xF0}, //0x8
	{0xF0, 0x90, 0xF0, 0x10, 0xF0}, //0x9
	{0xF0, 0x90, 0xF0, 0x90, 0x90}, //0xa
	{0xE0, 0x90, 0xE0, 0x90, 0xE0}, //0xb
	{0xF0, 0x80, 0x80, 0x80, 0xF0}, //0xc
	{0xE0, 0x90, 0x90, 0x90, 0xE0}, //0xd
	{0xF0, 0x80, 0xF0, 0x80, 0xF0}, //0xe
	{0xF0, 0x80, 0xF0, 0x80, 0x80}, //0xf
}

func (cpu *CPU) setupFont() {
	for i := 0; i <= 0xF; i++ {
		for j := 0; j < 5; j++ {
			cpu.memory[(i*5)+j] = FONT_TABLE[i][j]
		}
	}
}
