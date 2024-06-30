package ppu

// PPU represents the Picture Processing Unit of the NES.
type PPU struct {
	// Define PPU state and registers here
	// For example:
	// VRAM []uint8
	// OAM []uint8
	Registers struct {
		ppuCtrl uint8
		ppuMask uint8
		ppuStatus uint8
		oamAddr uint8
		oamData uint8
		ppuScroll uint8
		ppuAddr uint8
		ppuData uint8
	}
	//Internal registers
	internReg struct {
		v uint16
		t uint16
		x uint8
		w bool
	}
	type tile struct{
	//Uses two uint64's to store 16 bytes of a tile
		plane1 uint64
		plane2 uint64
	}
	//Left and right pattern tables(CHR-ROM/RAM)
	lptrnTable [256]tile
	rptrnTable [256]tile
	//4 Name tables
	nTable1 [32][32]uint8
	nTable2 [32][32]uint8
	nTable3 [32][32]uint8
	nTable4 [32][32]uint8
	
	OAM [64]uint32
	VRAM []uint8
	pallete [32]uint8
}

// NewPPU creates and initializes a new PPU instance.
func NewPPU() *PPU {
	ppu := &PPU{
		// Initialize PPU state and registers here
	}
	return ppu
}

// RenderFrame renders one frame of the PPU.
func (p *PPU) RenderFrame() {
	// Implement frame rendering logic here
}
