package memory

// Memory represents the memory of the NES.
type Memory struct {
	// Define memory regions and variables here
	// For example:
	// RAM []uint8
	// ROM []uint8
}

// NewMemory creates and initializes a new Memory instance.
func NewMemory() *Memory {
	mem := &Memory{
		// Initialize memory regions here
	}
	return mem
}

// Read reads data from memory at the specified address.
func (m *Memory) Read(address uint16) uint8 {
	// Implement memory read logic here
	return 0 // Placeholder return value
}

// Write writes data to memory at the specified address.
func (m *Memory) Write(address uint16, data uint8) {
	// Implement memory write logic here
}
