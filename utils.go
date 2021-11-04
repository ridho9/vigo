package vigo

func splitUint16(val uint16) (uint8, uint8) {
	v1 := uint8(val >> 8)
	v2 := uint8(val)
	return v1, v2
}

func splitUint8(val uint8) (uint8, uint8) {
	v1 := val >> 4
	v2 := val & 0b0000_1111
	return v1, v2
}

func combine3u4(i1, i2, i3 uint8) uint16 {
	return uint16(i1)<<8 + uint16(i2)<<4 + uint16(i3)
}

func combine2u4(i1, i2 uint8) uint8 {
	return i1<<4 + i2
}

func u8ToBits(v uint8) [8]bool {
	result := [8]bool{}

	for i := 0; i < 8; i++ {
		if v&(1<<i) > 0 {
			result[8-i-1] = true
		}
	}

	return result
}
