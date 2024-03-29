package vigo

import "fmt"

func PrintInst(i1, i2, i3, i4 uint8) string {
	i34 := combine2u4(i3, i4)
	i234 := combine3u4(i2, i3, i4)

	if i1 == 0 && i2 == 0 && i3 == 0 && i4 == 0 {
		return "HALT"
	}

	if i1 == 0 && i2 == 0 && i3 == 0xE && i4 == 0 {
		return "CLEAR"
	}

	if i1 == 0 && i2 == 0 && i3 == 0xE && i4 == 0xE {
		return "RETURN"
	}

	if i1 == 0x1 {
		return fmt.Sprintf("JUMP\t0x%X", i234)
	}

	if i1 == 0x2 {
		return fmt.Sprintf("CALL\t0x%X", i234)
	}

	if i1 == 0x3 {
		return fmt.Sprintf("SEQ\tV%X %#X", i2, i34)
	}

	if i1 == 0x4 {
		return fmt.Sprintf("SNEQ\tV%X %#X", i2, i34)
	}

	if i1 == 0x5 {
		return fmt.Sprintf("SEQ\tV%X V%X", i2, i3)
	}

	if i1 == 0x6 {
		return fmt.Sprintf("MOV\tV%X %#X", i2, i34)
	}

	if i1 == 0x7 {
		return fmt.Sprintf("ADD\tV%X %#X", i2, i34)
	}

	if i1 == 0x8 && i4 == 0x0 {
		return fmt.Sprintf("MOV\tV%X V%X", i2, i3)
	}

	if i1 == 0x8 && i4 == 0x1 {
		return fmt.Sprintf("OR\tV%X V%X", i2, i3)
	}

	if i1 == 0x8 && i4 == 0x2 {
		return fmt.Sprintf("AND\tV%X V%X", i2, i3)
	}

	if i1 == 0x8 && i4 == 0x3 {
		return fmt.Sprintf("XOR\tV%X V%X", i2, i3)
	}

	if i1 == 0x8 && i4 == 0x4 {
		return fmt.Sprintf("ADD.\tV%X V%X", i2, i3)
	}

	if i1 == 0x8 && i4 == 0x5 {
		return fmt.Sprintf("SUB.\tV%X V%X", i2, i3)
	}

	if i1 == 0x8 && i4 == 0x6 {
		return fmt.Sprintf("SHR.\tV%X V%X", i2, i3)
	}

	if i1 == 0x8 && i4 == 0x7 {
		return fmt.Sprintf("SUBB.\tV%X V%X", i2, i3)
	}

	if i1 == 0x8 && i4 == 0xE {
		return fmt.Sprintf("SHL.\tV%X V%X", i2, i3)
	}

	if i1 == 0x9 {
		return fmt.Sprintf("SNEQ\tV%X V%X", i2, i3)
	}

	if i1 == 0xA {
		return fmt.Sprintf("MOV\tI %#X", i234)
	}

	if i1 == 0xD {
		x, y, n := i2, i3, i4
		return fmt.Sprintf("DRAW\tV%X V%X %d", x, y, n)
	}

	if i1 == 0xF && i3 == 0x0 && i4 == 0x7 {
		return fmt.Sprintf("MOV\tV%X DT", i2)
	}

	if i1 == 0xF && i3 == 0x1 && i4 == 0x5 {
		return fmt.Sprintf("MOV\tDT V%X", i2)
	}

	if i1 == 0xF && i3 == 0x1 && i4 == 0x8 {
		return fmt.Sprintf("MOV\tST V%X", i2)
	}

	if i1 == 0xF && i3 == 0x5 && i4 == 0x5 {
		return fmt.Sprintf("STI\t%d", i2)
	}

	if i1 == 0xF && i3 == 0x3 && i4 == 0x3 {
		return fmt.Sprintf("BCDI\tV%X", i2)
	}

	if i1 == 0xF && i3 == 0x6 && i4 == 0x5 {
		return fmt.Sprintf("LDI\t%d", i2)
	}

	if i1 == 0xF && i3 == 0x2 && i4 == 0x9 {
		return fmt.Sprintf("LDIF\tV%X", i2)
	}

	if i1 == 0xF && i3 == 0x1 && i4 == 0xE {
		return fmt.Sprintf("ADD\tI V%X", i2)
	}

	if i1 == 0xE && i3 == 0x9 && i4 == 0xE {
		return fmt.Sprintf("SKP\tV%X", i2)
	}

	if i1 == 0xE && i3 == 0xA && i4 == 0x1 {
		return fmt.Sprintf("SKNP\tV%X", i2)
	}

	return "???"
}
