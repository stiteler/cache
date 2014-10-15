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
	runCLI()
}

func runCLI() {
	for {
		fmt.Println("(R)ead, (W)rite, or (D)isplay Cache?")
		var input string
		fmt.Scanf("%s", &input)

		switch input {
		case "R", "r":
			fmt.Println("What address would you like to read?")
			address := getHexAddressInput()
			value, hit := C.ReadByte(address)

			fmt.Printf("At that byte there is the value: %02X", value)
			if hit {
				fmt.Printf(", Cache Hit\n")
			} else {
				fmt.Printf(", Cache Miss\n")
			}
			continue

		case "W", "w":
			fmt.Println("What address would you like to write to?")
			address := getHexAddressInput()
			fmt.Println("What data would you like to write at that address?")
			data := getHexAddressInput()
			hit := C.WriteByte(address, byte(data))
			fmt.Printf("Value %X has been written to address %X", data, address)
			if hit {
				fmt.Printf(", Cache Hit\n")
			} else {
				fmt.Printf(", Cache Miss\n")
			}

		case "D", "d":
			C.Display()
			continue
		}
	}
}

func (c *Cache) WriteByte(addr uint32, b byte) bool {
	hit := false

	// get slot #
	slotNum := maskAndShift(slotNumMask, addr)
	fmt.Printf("slot Number on write: %X\n", slotNum)

	// get the slot we're working with
	slot := &c[slotNum]

	// calc blockOffset for use later
	blockOffset := maskAndShift(blockOffsetMask, addr)
	thisTag := getTag(addr)

	if slot.validBit && slot.tag == thisTag {
		// do i also need to check if the tag works out?
		// update the slot in cache

		// put new value there
		slot.block[blockOffset] = b

		hit = true
	} else {
		// bring block into slot from memory
		blockBegin := addr - blockOffset
		slot.block = getBlockFromMemory(blockBegin)
		slot.tag = thisTag
		slot.validBit = true
	}

	// the write through to memory regardless
	MM.WriteThrough(addr, b)

	// for now
	return hit

}

func getBlockFromMemory(blockBegin uint32) Block {
	var block Block
	for i, _ := range block {
		block[i] = byte(MM[blockBegin])
		blockBegin++
	}
	fmt.Println(block)
	return block
}

func (m *Memory) WriteThrough(addr uint32, b byte) {
	// simply assign byte to location in memory
	m[addr] = uint16(b)
}

// will return value, and true/false dep. on if it was a cache hit/miss
func (c *Cache) ReadByte(addr uint32) (data byte, hit bool) {
	// get details about address:
	slotNum := maskAndShift(slotNumMask, addr)
	slot := &c[slotNum]
	blockOffset := maskAndShift(blockOffsetMask, addr)
	thisTag := getTag(addr)
	slotTag := slot.tag
	isValid := slot.validBit

	// if in cache:
	if thisTag == slotTag && isValid {
		hit := true
		data := byte(slot.block[blockOffset])
		return data, hit
	} else {
		// not in cache:
		hit := false

		// bring block in from main memory
		blockBegin := addr - blockOffset
		slot.block = getBlockFromMemory(blockBegin)
		slot.tag = thisTag
		slot.validBit = true
		data := byte(slot.block[blockOffset])
		return data, hit
	}
}

func getHexAddressInput() uint32 {
	var input uint32
	fmt.Scanf("%X", &input)
	return input
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

/*func runCLIFromInputFile() {
	// todo: get hex from this channel
	inputChan := initializeInputChannel()
	switch <-inputChan {
	case "R", "r":
		fmt.Println("What address would you like to read?")
		address := getHexAddressInput()
		value, hit := C.ReadByte(address)

		fmt.Printf("At that byte there is the value: %02X", value)
		if hit {
			fmt.Printf(", Cache Hit\n")
		} else {
			fmt.Printf(", Cache Miss\n")
		}
		continue
	case "W", "w":
		fmt.Println("What address would you like to write to?")
		address := getHexAddressInput()
		hit := C.WriteByte(address)
		if hit {
			fmt.Printf(", Cache Hit\n")
		} else {
			fmt.Printf(", Cache Miss\n")
		}
	case "D", "d":
		C.Display()
		continue
	case "":
		break
	}
}

func initializeInputChannel() chan string {
	// load input from a file, fill a channel?
	inputChan := make(chan string, 46)

	return inputChan
}*/
