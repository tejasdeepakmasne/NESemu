package cpu

import (
	"fmt"
	"os"
)

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

func (c *CPU) reset() {
	c.accumulator = 0
	c.xIndex = 0
	c.yIndex = 0
	c.stackPointer = StackReset
	c.programCounter = c.readMemory16(0xFFFC)
	c.statusRegister = 0b00100100
}

func (c *CPU) loadProgram(program []uint8) {
	for i, value := range program {
		c.memory[0x8000+i] = value
	}
	c.writeMemory16(0xFFFC, 0x8000)
}

// ExecuteInstruction executes one instruction on the CPU.
func (c *CPU) ExecuteInstruction() {
	// Implement instruction execution logic here
	for {
		opcode := c.memory[c.programCounter]
		c.programCounter++
		switch opcode {

		case 0x69:
			c.adc(modeImmediate)
			c.programCounter++
		case 0x65:
			c.adc(modeZeroPage)
			c.programCounter++
		case 0x75:
			c.adc(modeZeroPageX)
			c.programCounter++
		case 0x6D:
			c.adc(modeAbsolute)
			c.programCounter += 2
		case 0x7d:
			c.adc(modeAbsoluteX)
			c.programCounter += 2
		case 0x79:
			c.adc(modeAbsoluteY)
			c.programCounter += 2
		case 0x61:
			c.adc(modeIndirectX)
			c.programCounter++
		case 0x71:
			c.adc(modeIndirectY)
			c.programCounter++

		case 0x29:
			c.and(modeImmediate)
			c.programCounter++
		case 0x25:
			c.and(modeZeroPage)
			c.programCounter++
		case 0x35:
			c.and(modeZeroPageX)
			c.programCounter++
		case 0x2d:
			c.and(modeAbsolute)
			c.programCounter += 2
		case 0x3d:
			c.and(modeAbsoluteX)
			c.programCounter += 2
		case 0x39:
			c.and(modeAbsoluteY)
			c.programCounter += 2
		case 0x21:
			c.and(modeIndirectX)
			c.programCounter++
		case 0x31:
			c.and(modeIndirectY)
			c.programCounter++

		case 0x0a:
			c.asl(modeAccumulator)
		case 0x06:
			c.asl(modeZeroPage)
			c.programCounter++
		case 0x16:
			c.asl(modeZeroPageX)
			c.programCounter++
		case 0x0e:
			c.asl(modeAbsolute)
			c.programCounter += 2
		case 0x1e:
			c.asl(modeAbsoluteX)
			c.programCounter += 2

		case 0x90:
			c.bcc()
			c.programCounter++

		case 0xb0:
			c.bcs()
			c.programCounter++

		case 0xf0:
			c.beq()
			c.programCounter++

		case 0x24:
			c.bit(modeZeroPage)
			c.programCounter++
		case 0x2c:
			c.bit(modeAbsolute)
			c.programCounter += 2

		case 0x30:
			c.bmi()
			c.programCounter++

		case 0xd0:
			c.bne()
			c.programCounter++

		case 0x10:
			c.bpl()
			c.programCounter++

		case 0x00:
			c.brk()

		case 0x50:
			c.bvc()
			c.programCounter++

		case 0x70:
			c.bvs()
			c.programCounter++

		case 0x18:
			c.clc()
		case 0xd8:
			c.cld()
		case 0x58:
			c.cli()
		case 0xb8:
			c.clv()

		case 0xc9:
			c.cmp(modeImmediate)
			c.programCounter++
		case 0xc5:
			c.cmp(modeZeroPage)
			c.programCounter++
		case 0xd5:
			c.cmp(modeZeroPageX)
			c.programCounter++
		case 0xcd:
			c.cmp(modeAbsolute)
			c.programCounter += 2
		case 0xdd:
			c.cmp(modeAbsoluteX)
			c.programCounter += 2
		case 0xd9:
			c.cmp(modeAbsoluteY)
			c.programCounter += 2
		case 0xc1:
			c.cmp(modeIndirectX)
			c.programCounter++
		case 0xd1:
			c.cmp(modeIndirectY)
			c.programCounter++

		case 0xe0:
			c.cpx(modeImmediate)
			c.programCounter++
		case 0xe4:
			c.cpx(modeZeroPage)
			c.programCounter++
		case 0xec:
			c.cpx(modeAbsolute)
			c.programCounter += 2

		case 0xc0:
			c.cpy(modeImmediate)
			c.programCounter++

		case 0xc4:
			c.cpy(modeZeroPage)
			c.programCounter++
		case 0xcc:
			c.cpy(modeAbsolute)
			c.programCounter += 2

		case 0xc6:
			c.dec(modeZeroPage)
			c.programCounter++
		case 0xd6:
			c.dec(modeZeroPageX)
			c.programCounter++
		case 0xce:
			c.dec(modeAbsolute)
			c.programCounter += 2
		case 0xde:
			c.dec(modeAbsoluteX)
			c.programCounter += 2

		case 0xca:
			c.dex()
		case 0x88:
			c.dey()

		case 0x49:
			c.eor(modeImmediate)
			c.programCounter++
		case 0x45:
			c.eor(modeZeroPage)
			c.programCounter++
		case 0x55:
			c.eor(modeZeroPageX)
			c.programCounter++
		case 0x4d:
			c.eor(modeAbsolute)
			c.programCounter += 2
		case 0x5d:
			c.eor(modeAbsoluteX)
			c.programCounter += 2
		case 0x59:
			c.eor(modeAbsoluteY)
			c.programCounter += 2
		case 0x41:
			c.eor(modeIndirectX)
			c.programCounter++
		case 0x51:
			c.eor(modeIndirectY)
			c.programCounter++

		case 0xe6:
			c.inc(modeZeroPage)
			c.programCounter++
		case 0xf6:
			c.inc(modeZeroPageX)
			c.programCounter++
		case 0xee:
			c.inc(modeAbsolute)
			c.programCounter += 2
		case 0xfe:
			c.inc(modeAbsoluteX)
			c.programCounter += 2

		case 0xe8:
			c.inx()
		case 0xc8:
			c.iny()

		case 0x4c:
			c.jmp(modeAbsolute)
			c.programCounter += 2
		case 0x6c:
			c.jmp(modeIndirect)
			c.programCounter += 2

		case 0x20:
			c.jsr()
			c.programCounter += 2

		case 0xa9:
			c.lda(modeImmediate)
			c.programCounter++
		case 0xa5:
			c.lda(modeZeroPage)
			c.programCounter++
		case 0xb5:
			c.lda(modeZeroPageX)
			c.programCounter++
		case 0xad:
			c.lda(modeAbsolute)
			c.programCounter += 2
		case 0xbd:
			c.lda(modeAbsoluteX)
			c.programCounter += 2
		case 0xb9:
			c.lda(modeAbsoluteY)
			c.programCounter += 2
		case 0xa1:
			c.lda(modeIndirectX)
			c.programCounter++
		case 0xb1:
			c.lda(modeIndirectY)
			c.programCounter++

		case 0xa2:
			c.ldx(modeImmediate)
			c.programCounter++
		case 0xa6:
			c.ldx(modeZeroPage)
			c.programCounter++
		case 0xae:
			c.ldx(modeAbsolute)
			c.programCounter += 2
		case 0xbe:
			c.ldx(modeAbsoluteY)
			c.programCounter += 2

		case 0xa0:
			c.ldy(modeImmediate)
			c.programCounter++
		case 0xa4:
			c.ldy(modeZeroPage)
			c.programCounter++
		case 0xb4:
			c.ldy(modeZeroPageX)
			c.programCounter++
		case 0xac:
			c.ldy(modeAbsolute)
			c.programCounter += 2
		case 0xbc:
			c.ldy(modeAbsoluteX)
			c.programCounter += 2

		case 0x4a:
			c.lsr(modeAccumulator)
		case 0x46:
			c.lsr(modeZeroPage)
			c.programCounter++
		case 0x56:
			c.lsr(modeZeroPageX)
			c.programCounter++
		case 0x4e:
			c.lsr(modeAbsolute)
			c.programCounter += 2
		case 0x5e:
			c.lsr(modeAbsoluteX)
			c.programCounter += 2

		case 0xea:
			c.nop()

		case 0x09:
			c.ora(modeImmediate)
			c.programCounter++
		case 0x05:
			c.ora(modeZeroPage)
			c.programCounter++
		case 0x015:
			c.ora(modeZeroPageX)
			c.programCounter++
		case 0x0d:
			c.ora(modeAbsolute)
			c.programCounter += 2
		case 0x1d:
			c.ora(modeAbsoluteX)
			c.programCounter += 2
		case 0x19:
			c.ora(modeAbsoluteY)
			c.programCounter += 2
		case 0x01:
			c.ora(modeIndirectX)
			c.programCounter++
		case 0x11:
			c.ora(modeIndirectY)
			c.programCounter++

		case 0x48:
			c.pha()
		case 0x08:
			c.php()
		case 0x68:
			c.pla()
		case 0x28:
			c.plp()

		case 0x2a:
			c.rol(modeAccumulator)
		case 0x26:
			c.rol(modeZeroPage)
			c.programCounter++
		case 0x36:
			c.rol(modeZeroPageX)
			c.programCounter++
		case 0x2e:
			c.rol(modeAbsolute)
			c.programCounter += 2
		case 0x3e:
			c.rol(modeAbsoluteX)
			c.programCounter += 2

		case 0x6a:
			c.ror(modeAccumulator)
		case 0x66:
			c.ror(modeZeroPage)
			c.programCounter++
		case 0x76:
			c.ror(modeZeroPageX)
			c.programCounter++
		case 0x6e:
			c.ror(modeAbsolute)
			c.programCounter += 2
		case 0x7e:
			c.ror(modeAbsoluteX)
			c.programCounter += 2

		case 0x40:
			c.rti()
		case 0x060:
			c.rts()

		case 0xe9:
			c.sbc(modeImmediate)
			c.programCounter++
		case 0xe5:
			c.sbc(modeZeroPage)
			c.programCounter++
		case 0xf5:
			c.sbc(modeZeroPageX)
			c.programCounter++
		case 0xed:
			c.sbc(modeAbsolute)
			c.programCounter += 2
		case 0xfd:
			c.sbc(modeAbsoluteX)
			c.programCounter += 2
		case 0xf9:
			c.sbc(modeAbsoluteY)
			c.programCounter += 2
		case 0xe1:
			c.sbc(modeIndirectX)
			c.programCounter++
		case 0xf1:
			c.sbc(modeIndirectY)
			c.programCounter++

		case 0x38:
			c.sec()
		case 0xf8:
			c.sed()
		case 0x78:
			c.sei()

		case 0x85:
			c.sta(modeZeroPage)
			c.programCounter++
		case 0x95:
			c.sta(modeZeroPageX)
			c.programCounter++
		case 0x8d:
			c.sta(modeAbsolute)
			c.programCounter += 2
		case 0x9d:
			c.sta(modeAbsoluteX)
			c.programCounter += 2
		case 0x99:
			c.sta(modeAbsoluteY)
			c.programCounter += 2
		case 0x81:
			c.sta(modeIndirectX)
			c.programCounter++
		case 0x91:
			c.sta(modeIndirectY)
			c.programCounter++

		case 0x86:
			c.stx(modeZeroPage)
			c.programCounter++
		case 0x96:
			c.stx(modeZeroPageY)
			c.programCounter++
		case 0x8e:
			c.stx(modeAbsolute)
			c.programCounter += 2

		case 0x84:
			c.sty(modeZeroPage)
			c.programCounter++
		case 0x94:
			c.sty(modeZeroPageX)
			c.programCounter++
		case 0x8c:
			c.sty(modeAbsolute)
			c.programCounter++

		case 0xaa:
			c.tax()
		case 0xa8:
			c.tay()
		case 0xba:
			c.tsx()
		case 0x8a:
			c.txa()
		case 0x9a:
			c.txs()
		case 0x98:
			c.tya()

		default:
			fmt.Fprintf(os.Stdout, "UNDEFINED BEHAVIOUR %v at %v", opcode, c.programCounter)

		}
	}
}

func (c *CPU) loadAndInterpret(program []uint8) {
	c.loadProgram(program)
	c.reset()
	c.ExecuteInstruction()
}
