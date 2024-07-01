package ppu
// PPU represents the Picture Processing Unit of the NES.
type tile struct{
	//Uses two uint64's to store 16 bytes of a tile
		plane1 uint64
		plane2 uint64
	}
type VRAMtype struct {
	pTable0 [256]tile
	pTable1 [256]tile
	nTable0 [32][32]uint8
	nTable1 [32][32]uint8
	nTable2 [32][32]uint8
	nTable3 [32][32]uint8
	palette [32]uint8 
}

type PPU struct {
	// Define PPU state and registers here
	// For example:
	// VRAM []uint8
	// OAM []uint8
	//Registers:
	ppuCtrl uint8
	ppuMask uint8
	ppuStatus uint8
	oamAddr uint8
	oamData uint8
	ppuScroll uint8
	ppuAddr uint8
	ppuData uint8

	//Internal registers
	v uint16
	t uint16
	x uint8
	w bool
	oddFrame bool

	VRAM VRAMtype
	OAM [64]uint32
}

// NewPPU creates and initializes a new PPU instance.
func NewPPU() *PPU {
	ppu := &PPU{
		// Initialize PPU state and registers here
		//All are implicitly set to zero at power up
		VRAM: VRAMtype{               
                        pTable0: [256]tile{
                        },                  
                        pTable1: [256]tile{
                        },                                                      
                }, 
	}
	return ppu
}
// RenderFrame renders one frame of the PPU.
func (p *PPU) RenderFrame() {
	// Implement frame rendering logic here
}
