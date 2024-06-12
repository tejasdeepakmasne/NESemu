package cpu

// CPU represents the Central Processing Unit of the NES.

const NonMaskableInterruptVector uint16 = 0xFFFA
const ResetVector uint16 = 0xFFFC
const InterruptRequestVector uint16 = 0xFFFE
const StackBase uint16 = 0x0100
const StackReset uint8 = 0xFD

func extractBit(val uint8, pos uint8) uint8 {
	return (val & (1 << pos)) >> pos
}

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
	// Implement stack pushStack logic here
	c.memory[StackBase+uint16(c.stackPointer)] = value
	c.stackPointer--
}
func (c *CPU) popStack() uint8 {
	// Implement stack popStack logic here
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
	// Implement stack popStack logic here
	msb := uint16(c.memory[c.stackPointer])
	c.stackPointer++
	lsb := uint16(c.memory[c.stackPointer])
	c.stackPointer++
	return (msb << 8) | lsb
}
func (c *CPU) pushStack16(value uint16) {
	// Implement stack pushStack logic here
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
		if indirectVector&0x00ff == 0x00ff {
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
		accumulator:    0,
		xIndex:         0,
		yIndex:         0,
		statusRegister: 0b00100100,
		programCounter: 0,
		stackPointer:   StackReset,
		memory:         make([]uint8, 0x10000),
	}
	return cpu
}

// Helper FUnctions for setting and clearing flags
func (c *CPU) setFlag(flags ...Flags) {
	for _, f := range flags {
		c.statusRegister |= 1 << f
	}
}

func (c *CPU) clearFlag(flags ...Flags) {
	for _, f := range flags {
		c.statusRegister &= ^(1 << f)
	}
}

func (c *CPU) setFlagToValue(flag Flags, value uint8) {
	if value == 1 {
		c.setFlag(flag)
	} else {
		c.clearFlag(flag)
	}
}

func (c *CPU) getFlag(flag Flags) uint8 {
	return (c.statusRegister & (1 << flag)) >> flag
}

func (c *CPU) updateZeroAndNegativeFlag(value uint8) {
	if value == 0 {
		c.setFlag(Z)
	} else {
		c.clearFlag(Z)
	}

	if value&128 == 128 {
		c.setFlag(Z)
	} else {
		c.clearFlag(Z)
	}
}

// INSTRUCTIONS
func (c *CPU) adc(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	res := c.accumulator + value + c.getFlag(C)
	if res > 255 {
		c.setFlag(C)
	}
	if res > 127 {
		c.setFlag(V)
	}

	c.accumulator = res
	c.updateZeroAndNegativeFlag(c.accumulator)

}

func (c *CPU) and(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	c.accumulator &= value
	c.updateZeroAndNegativeFlag(c.accumulator)
}

func (c *CPU) asl(mode AddressingMode) {
	if mode == modeAccumulator {
		c.setFlagToValue(C, extractBit(c.accumulator, 7))
		c.accumulator = c.accumulator << 1
	} else {
		address := c.addressMode(mode)
		value := c.readMemory(address)
		c.setFlagToValue(C, extractBit(value, 7))
		value = value << 1
		c.writeMemory(address, value)
	}

	c.updateZeroAndNegativeFlag(c.statusRegister)
}

func (c *CPU) bcc() {
	if c.getFlag(C) == 0 {
		address := c.addressMode(modeRelative)
		value := c.readMemory(address)
		c.programCounter += uint16(value)
	}
}

func (c *CPU) bcs() {
	if c.getFlag(C) == 1 {
		address := c.addressMode(modeRelative)
		value := c.readMemory(address)
		c.programCounter += uint16(value)
	}
}

func (c *CPU) beq() {
	if c.getFlag(Z) == 1 {
		address := c.addressMode(modeRelative)
		value := c.readMemory(address)
		c.programCounter += uint16(value)
	}
}

func (c *CPU) bit(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	res := c.accumulator & value
	if res == 0 {
		c.setFlag(Z)
	} else {
		c.clearFlag(Z)
	}
	c.setFlagToValue(V, extractBit(res, 6))
	c.setFlagToValue(N, extractBit(res, 7))
}

func (c *CPU) bmi() {
	if c.getFlag(N) == 1 {
		address := c.addressMode(modeRelative)
		value := c.readMemory(address)
		c.programCounter += uint16(value)
	}
}

func (c *CPU) bne() {
	if c.getFlag(C) == 0 {
		address := c.addressMode(modeRelative)
		value := c.readMemory(address)
		c.programCounter += uint16(value)
	}
}

func (c *CPU) bpl() {
	if c.getFlag(N) == 0 {
		address := c.addressMode(modeRelative)
		value := c.readMemory(address)
		c.programCounter += uint16(value)
	}
}

func (c *CPU) brk() {
	c.pushStack16(c.programCounter)
	c.pushStack(c.statusRegister)
	c.programCounter = c.readMemory16(InterruptRequestVector)
	c.setFlag(B)
}

func (c *CPU) bvc() {
	if c.getFlag(V) == 0 {
		address := c.addressMode(modeRelative)
		value := c.readMemory(address)
		c.programCounter += uint16(value)
	}
}

func (c *CPU) bvs() {
	if c.getFlag(V) == 1 {
		address := c.addressMode(modeRelative)
		value := c.readMemory(address)
		c.programCounter += uint16(value)
	}
}

func (c *CPU) clc() {
	c.clearFlag(C)
}

func (c *CPU) cld() {
	c.clearFlag(D)
}

func (c *CPU) cli() {
	c.clearFlag(I)
}

func (c *CPU) clv() {
	c.clearFlag(V)
}

func (c *CPU) cmp(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	res := c.accumulator - value
	if res >= uint8(0) {
		c.setFlag(C)
	}
	c.updateZeroAndNegativeFlag(res)
}

func (c *CPU) cpx(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	res := c.xIndex - value
	if res >= uint8(0) {
		c.setFlag(C)
	}
	c.updateZeroAndNegativeFlag(res)

}

func (c *CPU) cpy(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	res := c.yIndex - value
	if res >= uint8(0) {
		c.setFlag(C)
	}
	c.updateZeroAndNegativeFlag(res)

}

func (c *CPU) dec(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	value--
	c.writeMemory(address, value)
	c.updateZeroAndNegativeFlag(value)
}

func (c *CPU) dex() {
	c.xIndex--
	c.updateZeroAndNegativeFlag(c.xIndex)
}

func (c *CPU) dey() {
	c.yIndex--
	c.updateZeroAndNegativeFlag(c.yIndex)
}

func (c *CPU) eor(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	c.accumulator = c.accumulator ^ value
	c.updateZeroAndNegativeFlag(c.accumulator)
}

func (c *CPU) inc(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	value++
	c.writeMemory(address, value)
	c.updateZeroAndNegativeFlag(value)
}

func (c *CPU) inx() {
	c.xIndex++
	c.updateZeroAndNegativeFlag(c.xIndex)
}

func (c *CPU) iny() {
	c.yIndex++
	c.updateZeroAndNegativeFlag(c.yIndex)
}

func (c *CPU) jmp(mode AddressingMode) {
	address := c.addressMode(mode)
	c.programCounter = address

}

func (c *CPU) jsr() {
	c.pushStack16(c.programCounter - 1)
	address := c.addressMode(modeAbsolute)
	c.programCounter = address
}

func (c *CPU) lda(mode AddressingMode) {
	address := c.addressMode(mode)
	c.accumulator = c.readMemory(address)
	c.updateZeroAndNegativeFlag(c.accumulator)
}

func (c *CPU) ldx(mode AddressingMode) {
	address := c.addressMode(mode)
	c.xIndex = c.readMemory(address)
	c.updateZeroAndNegativeFlag(c.xIndex)
}

func (c *CPU) ldy(mode AddressingMode) {
	address := c.addressMode(mode)
	c.yIndex = c.readMemory(address)
	c.updateZeroAndNegativeFlag(c.yIndex)
}

func (c *CPU) lsr(mode AddressingMode) {
	if mode == modeAccumulator {
		c.setFlagToValue(C, extractBit(c.accumulator, 0))
		c.accumulator = c.accumulator >> 1
		c.updateZeroAndNegativeFlag(c.accumulator)
	} else {
		address := c.addressMode(mode)
		value := c.readMemory(address)
		c.setFlagToValue(C, extractBit(value, 0))
		value = value >> 1
		c.writeMemory(address, value)
		c.updateZeroAndNegativeFlag(value)
	}
}

func (c *CPU) nop() {

}

func (c *CPU) ora(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	c.accumulator = c.accumulator | value
	c.updateZeroAndNegativeFlag(c.accumulator)
}

func (c *CPU) pha() {
	c.pushStack(c.accumulator)
}

func (c *CPU) php() {
	c.pushStack(c.statusRegister)
}

func (c *CPU) pla() {
	c.accumulator = c.popStack()
	c.updateZeroAndNegativeFlag(c.accumulator)
}

func (c *CPU) plp() {
	c.statusRegister = c.popStack()
	c.updateZeroAndNegativeFlag(c.statusRegister)
}

func (c *CPU) rol(mode AddressingMode) {
	if mode == modeAccumulator {
		prevCarry := extractBit(c.statusRegister, 0)
		c.setFlagToValue(C, extractBit(c.accumulator, 7))
		c.accumulator = (c.accumulator << 1) | prevCarry
		c.updateZeroAndNegativeFlag(c.accumulator)
	} else {
		address := c.addressMode(mode)
		value := c.readMemory(address)
		prevCarry := extractBit(c.statusRegister, 0)
		c.setFlagToValue(C, extractBit(value, 7))
		value = (value << 1) | prevCarry
		c.writeMemory(address, value)
		c.updateZeroAndNegativeFlag(value)
	}
}

func (c *CPU) ror(mode AddressingMode) {
	if mode == modeAccumulator {
		prevCarry := extractBit(c.statusRegister, 0)
		c.setFlagToValue(C, extractBit(c.accumulator, 0))
		c.accumulator = (c.accumulator >> 1) | (prevCarry << 7)
		c.updateZeroAndNegativeFlag(c.accumulator)
	} else {
		address := c.addressMode(mode)
		value := c.readMemory(address)
		prevCarry := extractBit(c.statusRegister, 0)
		c.setFlagToValue(C, extractBit(c.accumulator, 0))
		value = (value >> 1) | (prevCarry << 7)
		c.writeMemory(address, value)
		c.updateZeroAndNegativeFlag(value)
	}
}

func (c *CPU) rti() {
	c.statusRegister = c.popStack()
	c.programCounter = c.popStack16()
}

func (c *CPU) rts() {
	c.programCounter = c.popStack16()
}

func (c *CPU) sbc(mode AddressingMode) {
	address := c.addressMode(mode)
	value := c.readMemory(address)
	res := c.accumulator - value - (1 - c.getFlag(C))
	if res > 255 {
		c.setFlag(C)
	}
	if res > 127 {
		c.setFlag(V)
	}

	c.accumulator = res
	c.updateZeroAndNegativeFlag(c.accumulator)

}

func (c *CPU) sec() {
	c.setFlagToValue(C, 1)
}

func (c *CPU) sed() {
	c.setFlagToValue(D, 1)
}

func (c *CPU) sei() {
	c.setFlagToValue(I, 1)
}

func (c *CPU) sta(mode AddressingMode) {
	address := c.addressMode(mode)
	c.writeMemory(address, c.accumulator)
}

func (c *CPU) stx(mode AddressingMode) {
	address := c.addressMode(mode)
	c.writeMemory(address, c.xIndex)
}

func (c *CPU) sty(mode AddressingMode) {
	address := c.addressMode(mode)
	c.writeMemory(address, c.yIndex)
}

func (c *CPU) tax() {
	c.xIndex = c.accumulator
	c.updateZeroAndNegativeFlag(c.xIndex)
}

func (c *CPU) tay() {
	c.yIndex = c.accumulator
	c.updateZeroAndNegativeFlag(c.yIndex)
}

func (c *CPU) tsx() {
	c.xIndex = c.statusRegister
	c.updateZeroAndNegativeFlag(c.xIndex)
}

func (c *CPU) txa() {
	c.accumulator = c.xIndex
	c.updateZeroAndNegativeFlag(c.accumulator)
}

func (c *CPU) txs() {
	c.statusRegister = c.xIndex
	c.updateZeroAndNegativeFlag(c.xIndex)
}

func (c *CPU) tya() {
	c.accumulator = c.yIndex
	c.updateZeroAndNegativeFlag(c.accumulator)
}

// ExecuteInstruction executes one instruction on the CPU.
func (c *CPU) ExecuteInstruction() {
	// Implement instruction execution logic here
}
