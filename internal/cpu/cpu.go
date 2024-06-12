package cpu

// CPU represents the Central Processing Unit of the NES.

const NonMaskableInterruptVector uint16 = 0xFFFA
const ResetVector uint16 = 0xFFFC
const InterruptRequestVector uint16 = 0xFFFE
const StackBase uint16 = 0x0100
const StackReset uint8 = 0xFD

type CPU struct {
	accumulator    uint8
	xIndex         uint8
	yIndex         uint8
	stackPointer   uint8
	programCounter uint16
	statusRegister uint8

	memory []uint8
}

type Flags uint8

const (
	C Flags = iota // Carry
	Z              // Zero
	I              // Interrupt Disable
	D              // Decimal Mode
	B              // B Flag
	X              // Unused
	V              // Overflow
	N              // Negative
)

type AddressingMode int

const (
	modeImmediate AddressingMode = iota
	modeZeroPage
	modeAbsolute
	modeZeroPageX
	modeZeroPageY
	modeAbsoluteX
	modeAbsoluteY
	modeIndirectX
	modeIndirectY
	modeRelative
	modeAccumulator
	modeIndirect
	modeNoneAddressing
)

func (c *CPU) readMemory(address uint16) uint8 {
	// Implement memory read logic here
	return c.memory[address]
}
func (c *CPU) writeMemory(address uint16, value uint8) {
	// Implement memory write logic here
	c.memory[address] = value
}
func (c *CPU) pushStack(value uint8) {
	// Implement stack push logic here
	c.memory[StackBase+uint16(c.stackPointer)] = value
	c.stackPointer--
}
func (c *CPU) popStack() uint8 {
	// Implement stack pop logic here
	top := c.memory[StackBase+uint16(c.stackPointer)]
	c.stackPointer++
	return top
}
func (c *CPU) readMemory16(address uint16) uint16 {
	// Implement memory read logic here
	lsb := uint16(c.memory[address])
	msb := uint16(c.memory[address+1])
	return (msb << 8) | lsb
}
func (c *CPU) writeMemory16(address uint16, value uint16) {
	// Implement memory write logic here
	lsb := uint8(value & 0xFF)
	msb := uint8(value >> 8)
	c.writeMemory(address, lsb)
	c.writeMemory(address+1, msb)
}
func (c *CPU) popStack16() uint16 {
	// Implement stack pop logic here
	msb := uint16(c.memory[c.stackPointer])
	c.stackPointer++
	lsb := uint16(c.memory[c.stackPointer])
	c.stackPointer++
	return (msb << 8) | lsb
}
func (c *CPU) pushStack16(value uint16) {
	// Implement stack push logic here
	lsb := uint8(value & 0xFF)
	msb := uint8(value >> 8)
	c.memory[StackBase+uint16(c.stackPointer)] = lsb
	c.stackPointer--
	c.memory[StackBase+uint16(c.stackPointer)] = msb
	c.stackPointer--
}
func (c *CPU) addressMode(mode AddressingMode) uint16 {
	var address uint16
	switch mode {
	case modeImmediate:
		// Immediate addressing mode: The operand is the next byte after the instruction.
		// Implement logic to fetch the operand from memory and return the address.
		address = c.programCounter

	case modeZeroPage:
		// Zero Page addressing mode: The operand is the byte at the zero page address.
		// Implement logic to fetch the operand from memory and return the address.
		address = uint16(c.readMemory(c.programCounter))

	case modeAbsolute:
		// Absolute addressing mode: The operand is the byte at the specified address.
		// Implement logic to fetch the operand from memory and return the address.
		address = c.readMemory16(c.programCounter)

	case modeZeroPageX:
		// Zero Page X addressing mode: The operand is the byte at the zero page address plus the value of the X register.
		// Implement logic to fetch the operand from memory and return the address.
		base_address := c.readMemory(c.programCounter)
		address = uint16(base_address + c.xIndex)

	case modeZeroPageY:
		// Zero Page Y addressing mode: The operand is the byte at the zero page address plus the value of the Y register.
		// Implement logic to fetch the operand from memory and return the address.
		base_address := c.readMemory(c.programCounter)
		address = uint16(base_address + c.yIndex)

	case modeAbsoluteX:
		// Absolute X addressing mode: The operand is the byte at the specified address plus the value of the X register.
		// Implement logic to fetch the operand from memory and return the address.
		base_address := c.readMemory16(c.programCounter)
		address = base_address + uint16(c.xIndex)

	case modeAbsoluteY:
		// Absolute Y addressing mode: The operand is the byte at the specified address plus the value of the Y register.
		// Implement logic to fetch the operand from memory and return the address.
		base_address := c.readMemory16(c.programCounter)
		address = base_address + uint16(c.yIndex)

	case modeIndirectX:
		// Indirect X addressing mode: The operand is the byte at the address formed by adding the X register to the zero page address.
		// Implement logic to fetch the operand from memory and return the address.
		base := c.readMemory(c.programCounter)
		var offset uint8 = base + c.xIndex
		lsb := c.readMemory(uint16(offset))
		msb := c.readMemory(uint16(offset + 1))
		address = uint16(msb)<<8 | uint16(lsb)


	case modeIndirectY:
		// Indirect Y addressing mode: The operand is the byte at the address formed by adding the Y register to the zero page address.
		// Implement logic to fetch the operand from memory and return the address.
		base := c.readMemory(c.programCounter)
		var offset uint8 = base + c.yIndex
		lsb := c.readMemory(uint16(offset))
		msb := c.readMemory(uint16(offset + 1))
		address = uint16(msb)<<8 | uint16(lsb)

	case modeRelative:
		// Relative addressing mode: The operand is a signed 8-bit offset relative to the program counter.
		// Implement logic to calculate the target address based on the offset and return it.
		address = c.programCounter

	case modeAccumulator:
		// Accumulator addressing mode: The operand is the accumulator register itself.
		// No additional logic is needed, simply return the address of the accumulator.

	case modeIndirect:
		// Indirect addressing mode: The operand is the address stored at the specified address.
		// Implement logic to fetch the operand from memory and return the address.
		lsb := uint16(c.readMemory(c.programCounter))
		msb := uint16(c.readMemory(c.programCounter + 1))
		indirectVector := (msb << 8) | lsb
		address_lsb := uint16(c.readMemory(indirectVector))
		address_msb := uint16(c.readMemory(indirectVector + 1))
		if indirectVector & 0x00ff == 0x00ff {
			address_msb = address_lsb & 0xff00
		}
		address = (address_msb << 8) | address_lsb

	case modeNoneAddressing:
		// No addressing mode: The instruction does not have an operand.
		// No additional logic is needed, simply return 0 for both addresses.
	}

	return address
}

// NewCPU creates and initializes a new CPU instance.
func NewCPU() *CPU {
	cpu := &CPU{
		// Initialize CPU state and registers here
	}
	return cpu
}

// ExecuteInstruction executes one instruction on the CPU.
func (c *CPU) ExecuteInstruction() {
	// Implement instruction execution logic here
}
