package golime

import (
	"bufio"
	"os"
)

var LIME_MAGIC []byte = []byte{0x45, 0x4d, 0x69, 0x4c} // LiME

type limeBlock struct {
	fileOffset uint64     // offset to startAddress in file, not to limeblock itself
	header     limeHeader // stores the actual header
}

const LIME_HEADER_SIZE = 32

type LiME struct {
	file   *os.File
	blocks []limeBlock
}

/*
	typedef struct {
		unsigned int magic;        // Always 0x4C694D45 (LiME)
		unsigned int version;      // Header version number
		unsigned long long s_addr; // Starting address of range
		unsigned long long e_addr; // Ending address of range
		unsigned char reserved[8]; // Currently all zeros
	} __attribute__ ((__packed__)) lime_mem_range_header;
*/

type limeHeader struct {
	Magic        uint32 // LIME_MAGIC
	Version      uint32 // will be 1
	StartAddress uint64
	EndAddress   uint64
	Reserved     [8]byte // currently all 0's
}

type Reader struct {
	r   *bufio.Reader
	pos uint64 // because bufio doesn't track this usually
}
