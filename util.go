package burst

// assumes only 64 or  32 bit pointers
func getPtr(b []byte) uintptr {

	switch ptrOffset {
	case 8:
		return uintptr(b[7]) | uintptr(b[6])<<8 | uintptr(b[5])<<16 | uintptr(b[4])<<24 | uintptr(b[3])<<32 | uintptr(b[2])<<40 | uintptr(b[1])<<48 | uintptr(b[0])<<56
	case 4:
		return uintptr(b[3]) | uintptr(b[2])<<8 | uintptr(b[1])<<16 | uintptr(b[0])<<24

	}
	panic("Invalid size detected")
}

// assumes only 64 or  32 bit pointers
func insertPtr(b []byte, ptr uintptr) {
	switch ptrOffset {
	case 8:
		b[0] = byte(ptr >> 56)
		b[1] = byte(ptr >> 48)
		b[2] = byte(ptr >> 40)
		b[3] = byte(ptr >> 32)
		b[4] = byte(ptr >> 24)
		b[5] = byte(ptr >> 16)
		b[6] = byte(ptr >> 8)
		b[7] = byte(ptr)
		return
	case 4:
		b[0] = byte(ptr >> 24)
		b[1] = byte(ptr >> 16)
		b[2] = byte(ptr >> 8)
		b[3] = byte(ptr)
		return
	}
	panic("Invalid size detected")
}

// To be completely safe Must give array larer then 10
func uleb128enc(bout []byte, value int) int {
	/* To encode:
	Grab the lowest 7 bits of your value and store them in a byte,
	this is what you're going to output.
	Shift the value 7 bits to the right, getting rid of those 7 bits you just grabbed.
	If the value is non-zero (ie. after you shifted away 7 bits from it),
	set the high bit of the byte you're going to output before you output it.
	Output the byte
	If the value is non-zero (ie. same check that resulted in setting the high bit),
	go back and repeat the steps from the start*/
	if len(bout) < 10 {
		panic("Need to give a buffer of at least 10")
	}
	count := 0
	for first, i := true, 0; first || value > 0; i++ {
		first = false
		lower7bits := byte(value & 0x7f)
		value >>= 7
		if value > 0 {
			lower7bits |= 128
		}
		bout[i] = lower7bits
		count++

	}
	return count
}

// value is the logical representation,
//count is the phsycial number of bits  used for representation 
func uleb128dec(bout []byte) (value int, count int) {
	/* To decode:
	Start at bit-position 0
	Read one byte from the file
	Store whether the high bit is set, and mask it away
	OR in the rest of the byte into your final value, at the bit-position you're at
	If the high bit was set, increase the bit-position by 7, and repeat the steps,
	skipping the first one (don't reset the bit-position). */

	if len(bout) < 10 {
		panic("Need to give a buffer of at least 10")
	}
	var lower7bits, shift byte

	for more, i := true, 0; more; i++ {
		lower7bits = bout[i]
		more = (lower7bits & 128) != 0
		value |= int(lower7bits&0x7f) << shift
		shift += 7
		count++

	}
	return value, count

}
