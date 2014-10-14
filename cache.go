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
	tag      uint32 // 24 bits, for this assgn, only least sig, 4 are used.
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
	fmt.Printf("MM[0x7FF]: %X\n\n", MM[0x7FF])

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

func (c *Cache) WriteByte(addr uint32, b byte) {
	// what if slot is not valid on write? , we just write to memory

	// get slot #
	slotNum := maskAndShift(slotNumMask, addr)
	fmt.Println("slot Number on write: " + slotNum)

	// get the slot we're working with?
	slot := c[slotNum]

	if slot.validBit {
		// update the slot in cache

		// calculate location for new byte in slot.block
		blockOffset := maskAndShift(blockOffsetMask, addr)

		// put new value there
		slot.block[blockOffset] = b

		// put slot in cache

		// write through to memory

	}

	// the write through to memory regardless
	MM.WriteThrough(addr, b)

}

func (m *Memory) WriteThrough(addr uint32, b byte) {
	// simply assign byte to location in memory
	m[addr] = b
}

// wil return value, and true/false dep. on if it was a cache hit/miss
func (c *Cache) ReadByte(addr uint32) (data byte, hit bool) {
	// if in cache (i.e. valid bit true, and tags match)
	tag := getTag(addr)

	// block number is address of first byte in block.. right? so is taht blockBeginMask for this cache?
	blockBegin := maskAndShift(blockBeginMask, addr)

	// ask, is this block in cache?
	// if yes: return value from cache, and true, for cache hit
	// if no: update slot in cache with this block, update valid bit, return value and false for cache miss
	// do a "getBlock()" method with a starting address that returns a slice of bytes?

	return byte(0), false
}

// maskAndShift() returns desired bits in a 16-bit value
// depending on the mask (including a shift value)
func maskAndShift(mask Mask, addr uint32) uint32 {
	return (addr & uint32(mask.bits)) >> mask.shift
}

func getTag(addr uint32) uint32 {
	return addr >> 8
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
