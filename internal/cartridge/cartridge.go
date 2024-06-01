package cartridge

// Cartridge represents the NES cartridge containing ROM data.
type Cartridge struct {
	// Define cartridge properties and ROM data here
	// For example:
	// ROM []uint8
	// Name string
	// Mapper int
}

// LoadCartridge loads an NES cartridge ROM from the specified file path.
func LoadCartridge(filePath string) *Cartridge {
	// Implement cartridge loading logic here
	return nil // Placeholder return value
}

// ReadPRGByte reads a byte from the PRG ROM of the cartridge at the specified address.
func (c *Cartridge) ReadPRGByte(address uint16) uint8 {
	// Implement PRG ROM read logic here
	return 0 // Placeholder return value
}

// WritePRGByte writes a byte to the PRG ROM of the cartridge at the specified address.
func (c *Cartridge) WritePRGByte(address uint16, data uint8) {
	// Implement PRG ROM write logic here
}
