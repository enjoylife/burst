package burst

import (
	"bytes"
	"fmt"
	"unsafe"
)

var _ = fmt.Println

type action byte

//var _x interface{} = 0

const (
	maxLen    int16 = 2<<14 - 1 //int16(^uint16(0) >> 1)
	lenOffset int   = 2         // int(unsafe.Sizeof(int16(0)))
	ptrOffset int   = int(unsafe.Sizeof(uintptr(0)))
	//ptrOffset int    = int(unsafe.Sizeof(_x))
	_remove action = iota + 1
	_check
	_update
)

// need to hold cache friendly lines of an dynamic growing array(aka slices)
type container struct {
	exhaustFlag bool
	records     []byte
}

// Search through a containers records utilizing a packed byte representation.
// Quirks:
// Fails with false if byte slice is empty, or if data not matched.
// Ignores the ptrIn unless the action is  _update. 
// On _remove returns the ptr of the deleted data
// On _update returns the previously stored ptr
func (c *container) search(data []byte, ptrIn uintptr, a action) (uintptr, bool) {

	checkLen := int16(len(data))
	if checkLen == 0 || checkLen > maxLen || checkLen > int16(len(c.records)) {
		return 0, false
	}

	var dend, dstart int

	for { // linear scan
		dlen := uint16(c.records[dend]) | uint16(c.records[dend+1])<<8
		skip := lenOffset + ptrOffset + int(dlen)
		strRemain := c.records[dstart+lenOffset+ptrOffset : (dend + skip)]
		dtest := bytes.Equal(strRemain, data)
		// search success
		if dtest {
			ptrBuff := c.records[dstart+lenOffset : dstart+lenOffset+ptrOffset]
			ptr := getPtr(ptrBuff)
			switch a {
			case _check:
				return ptr, true
			case _remove:
				c.records = append(
					c.records[:dstart],
					c.records[(dend+skip):]...,
				)
				return ptr, true
			case _update:
				insertPtr(ptrBuff, ptrIn)
				return ptr, true
			}
		}

		dend += skip
		dstart += skip

		// search failed
		if int(dend) == len(c.records) {
			switch a {
			case _check:
				return 0, false
			case _remove:
				return 0, false
			case _update:
				c.extend(data, ptrIn)
				return ptrIn, true
			}
		}
	}
	panic("Massive error") // who knows what might happen :)
	return 0, false
}

// extend at the end of the slice return false only if bytes to be inserted are longer then capable
func (c *container) extend(byteRemaining []byte, ptr uintptr) bool {

	checkLen := int16(len(byteRemaining))
	if checkLen == 0 || checkLen > maxLen {
		return false
	}
	var bufferLen [2]byte
	bufferLen[0] = byte(checkLen)
	bufferLen[1] = byte(checkLen >> 8)

	var bufferPtr = make([]byte, ptrOffset)
	insertPtr(bufferPtr, ptr)

	c.records = append(c.records, bufferLen[0])
	c.records = append(c.records, bufferLen[1])
	c.records = append(c.records, bufferPtr...)
	c.records = append(c.records, byteRemaining...)
	return true
}
