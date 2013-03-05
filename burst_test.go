package burst

import (
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand"
	"testing"
)

var _ = fmt.Println

func TestEncodingStuff(t *testing.T) {
	var buffer = make([]byte, 10)

	var length int = int(^uint(0) >> 1)
	length -= 2 << 20
	for u := int(^uint(0) >> 1); u > length; u-- {
		_ = uleb128enc(buffer, u)
		var u2, _ = uleb128dec(buffer)
		//fmt.Println(u)
		if u2 != u {
			t.Error("Encoding error")
		}
	}

	length = 2 << 20

	for u := 0; u < length; u++ {
		_ = uleb128enc(buffer, u)
		var u2, _ = uleb128dec(buffer)
		//fmt.Println(u)
		if u2 != u {
			t.Error("Encoding error")
		}
	}
}

func TestContainerextend(t *testing.T) {
	var c container
	var check bool
	var u uintptr
	check = c.extend([]byte{1, 1}, u)
	fmt.Println("Made it")
	if check == false {
		t.Error("Failed base case extend")
	}
	check = c.extend([]byte{2, 2, 2}, u)
	if check == false {
		t.Error("Unable trivally extend")
	}
	var b []byte
	check = c.extend(b, u)
	if check == true {
		t.Error("Need to avoid empty extend")
	}

	check = c.extend([]byte{1, 2, 3, 4}, u)
	if check == false {
		t.Error("Unable trivally extend")
	}
	//fmt.Println(c)

}

func TestContainerSearch(t *testing.T) {
	var c container
	var check bool
	var u uintptr = 0

	_, check = c.search([]byte{1, 1}, 0, _check)
	if check == true {
		t.Error("Wrong for trival base case")
	}
	c.extend([]byte{1, 1}, u)
	c.extend([]byte{2, 2, 2, 2, 2}, u)
	c.extend([]byte{4}, u)
	c.extend([]byte{5}, u)
	c.extend([]byte{6, 6, 6}, u)
	fmt.Println(c.records)
	_, check = c.search([]byte{1, 1}, 0, _check)
	if check == false {
		t.Error("Wrong for trival search case")
	}
	_, check = c.search([]byte{1}, 0, _check)
	if check == true {
		t.Error("Wrong for trival search case")
	}
	_, check = c.search([]byte{5}, 0, _check)
	if check == false {
		t.Error("Wrong for trival search case")
	}
	_, check = c.search([]byte{6, 6, 6}, 0, _check)
	if check == false {
		t.Error("Wrong for trival search case")
	}

	var b []byte
	_, check = c.search(b, 0, _check)
	if check == true {
		t.Error("Need to avoid empty search")
	}

}

func TestContainerDelete(t *testing.T) {
	var c container
	var check bool
	var u uintptr = 0

	_, check = c.search([]byte{1, 1}, 0, _check)
	if check == true {
		t.Error("Wrong for trival base case")
	}
	c.extend([]byte{1, 1}, u)
	c.extend([]byte{2, 2, 2, 2, 2}, u)
	c.extend([]byte{4}, u)
	c.extend([]byte{5}, ^uintptr(0))
	_, check = c.search([]byte{5}, 0, _remove)
	if check != true {
		t.Error("Failed to remove last element")
	}
	fmt.Println("Before: ", c.records)
	_, check = c.search([]byte{1, 1}, 0, _remove)
	if check != true {
		t.Error("Failed to remove first element")
	}
	fmt.Println("After:  ", c.records)

}

func BenchmarkContainerextend(b *testing.B) {

	var c container
	var u uintptr = 0
	for i := 0; i < b.N; i++ {
		b := make([]byte, mrand.Intn(10))
		io.ReadFull(rand.Reader, b)
		c.extend(b, u)
	}
}

func BenchmarkContainerSearch(b *testing.B) {
	b.StopTimer()
	var c container
	var u uintptr = 0
	var cache = make(map[int][]byte)
	var size = 100
	for i := 0; i < size; i++ {
		by := make([]byte, mrand.Intn(10))
		io.ReadFull(rand.Reader, by)
		cache[i] = by
		c.extend(by, u)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.search(cache[mrand.Intn(size)], 0, _check)
	}
}

func TestRandomStuff(t *testing.T) {

}

func BenchmarkEncodingStuff(b *testing.B) {

	var buffer = make([]byte, 10)
	var u2 int
	for i := 0; i < b.N; i++ {
		_ = uleb128enc(buffer, i)
		u2, _ = uleb128dec(buffer)
		_ = u2
	}
}
