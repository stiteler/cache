package main

import (
	"fmt"
)

type Memory [2048]uint16

type Slot struct {
	number   uint8 // 4 bits
	validBit bool
	//dirtyBit bool
	tag   uint32 // 24 bits
	block [16]byte
}

type Cache [16]Slot

type Mask struct {
	bits  uint32
	shift uint8
}

var slotNumMask = Mask{0xF0, 4}
var tagMask = Mask{0xFFFFFF00, 8}

var MM = &Memory{}

func main() {
	MM = initMemory()
	for _, i := range MM {
		fmt.Printf("0x%X ", MM[i])
	}
}

func initMemory() *Memory {
	memory := &Memory{}
	inc := uint16(0x00)
	for i, _ := range memory {
		memory[i] = inc % 0x100
		inc++
	}
	return memory
}

func (c *Cache) Write(addr, val uint32) {
	// create slot:
	// get block
	// set valid bit

	// put slot in cache

	// write through to memory
	MM.WriteThrough(addr, val)

}

// wil return value, and true/false dep. on if it was a cache hit/miss
func (c *Cache) Read(addr uint32) (uint32, bool) {

	return 0, false
}

func (m *Memory) WriteThrough(addr, val uint32) {

}
