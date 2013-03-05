package burst

import (
	"bytes"
	"fmt"
	"unsafe"
)

var _ = fmt.Println

type action byte

const (
	ptrOffset int    = int(unsafe.Sizeof(uintptr(0)))
	_remove   action = iota + 1
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

	if len(data) == 0 {
		return 0, false
	}
	// avoid data len which is larger then current records 
	if len(c.records) <= 2 {
		return 0, false
	}

	var count, dlen, dend, dstart int

	for { // linear scan
		dlen, count = uleb128dec(c.records[dend:])
		skip := count + ptrOffset + dlen
		strRemain := c.records[dstart+count+ptrOffset : (dend + skip)]
		//fmt.Println("MADIT")
		dtest := bytes.Equal(strRemain, data)
		// search success
		if dtest {
			ptrBuff := c.records[dstart+count : dstart+count+ptrOffset]
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
	var dlen = len(byteRemaining)
	// max size is 2 bytes
	if byteRemaining == nil {
		return false
	}

	var bufferLen = make([]byte, 10)
	var bufferPtr = make([]byte, ptrOffset)
	var count = uleb128enc(bufferLen, dlen)
	insertPtr(bufferPtr, ptr)

	c.records = append(c.records, bufferLen[0:count]...)
	c.records = append(c.records, bufferPtr...)
	c.records = append(c.records, byteRemaining...)
	return true
}
