package main

import (
	"fmt"
	"strings"
)

type Memory [2048]uint16
type Cache [16]Slot
type Block [16]byte

type Slot struct {
	slotNum  uint8 // 4 bits
	validBit bool
	tag      uint32 // 24 bits
	block    Block
}

type Mask struct {
	bits  uint32
	shift uint8
}

var slotNumMask = Mask{0xF0, 4}
var blockOffsetMask = Mask{0x0F, 0}
var blockBeginMask = Mask{0xFFFFFFF0, 4}

//var tagMask = Mask{0xFFFFFF00, 8} just need to shift over

var MM = &Memory{}
var C = &Cache{}

func main() {
	MM.initialize()
	C.initialize()
	for _, i := range MM {
		fmt.Printf("%X ", MM[i])
	}
	fmt.Println("")
	fmt.Printf("MM[0x100]: %X\n", MM[0x100])
	fmt.Printf("MM[0x7FF]: %X\n", MM[0x7FF])

	C.Display()
}

func (m *Memory) initialize() {
	inc := uint16(0x00)
	for i, _ := range m {
		m[i] = inc % 0x100
		inc++
	}
}

func (c *Cache) initialize() {
	slotInc := uint8(0x0)
	for i, _ := range c {
		c[i] = Slot{slotNum: slotInc}
		slotInc++
	}
}

func (c *Cache) WriteByte(addr uint32, B byte) {
	// create slot:
	// get block
	// set valid bit

	// put slot in cache

	// write through to memory
	MM.WriteThrough(addr, B)

}

// wil return value, and true/false dep. on if it was a cache hit/miss
func (c *Cache) ReadByte(addr uint32) (byte, bool) {

	return byte(0), false
}

func (m *Memory) WriteThrough(addr uint32, B byte) {

}

func buildSlot() *Slot {
	return nil
}

// maskAndShift() returns desired bits in a 16-bit value
// depending on the mask (including a shift value)
func maskAndShiftShort(mask Mask, inputBits int16) int16 {
	return (inputBits & int16(mask.bits)) >> mask.shift
}

func getTag(address uint32) uint32 {
	return address >> 8
}

func (c *Cache) Display() {
	fmt.Println("Slot#|Valid| Tag | Data")
	for _, slot := range c {
		fmt.Println(slot)
	}
}

func (s Slot) String() string {
	var validPrint string = ""
	if s.validBit {
		validPrint = "1"
	} else {
		validPrint = "0"
	}
	return fmt.Sprintf("  %X  |  %s  |  %X  | %v", s.slotNum, validPrint, s.tag, s.block)
}

func (b Block) String() string {
	blockStrings := make([]string, len(b))
	for i, value := range b {
		blockStrings[i] = fmt.Sprintf("%02X", value)
	}
	return fmt.Sprintf(strings.Join(blockStrings, " "))
}
