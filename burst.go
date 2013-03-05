package burst

import (
	"fmt"
)

var _ = fmt.Println

type Burst struct {
	//either a trieNode or a container
	root   interface{}
	height uint64
	size   uint64
}

type trieNode struct {
	exhaustFlag bool
	data        uintptr
	// will point either to another trieNode or a container
	next [256]interface{}
}

func New() *Burst {
	c := &container{exhaustFlag: false, records: make([]byte, 0)}
	return &Burst{root: c, height: 1, size: 0}
}

func (burst *Burst) Search(query []byte) bool {
	return true
}

func (b *Burst) Check(data []byte) (bool, error) {
	return true, fmt.Errorf("Error")
}
