package cpu

// CPU represents the Central Processing Unit of the NES.

const NonMaskableInterruptVector uint16 = 0xFFFA
const ResetVector uint16 = 0xFFFC
const InterruptRequestVector uint16 = 0xFFFE
const StackBase uint16 = 0x0100
const StackReset uint8 = 0xFD

type CPU struct {
	accumulator uint8
	xIndex uint8
	yIndex uint8
	stackPointer uint8
	programCounter uint16
	statusRegister uint8

	memory []uint8
}

type Flags uint8

const (
	C Flags = iota// Carry
	Z // Zero
	I // Interrupt Disable
	D // Decimal Mode
	B // B Flag
	X // Unused
	V // Overflow
	N // Negative
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
}
func (c *CPU) writeMemory(address uint16, value uint8) {
	// Implement memory write logic here
}
func (c *CPU) pushStack(value uint8) {
	// Implement stack push logic here
}
func (c *CPU) popStack() uint8 {
	// Implement stack pop logic here
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
