package golime

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

func initLiME(f *os.File) (*LiME, error) {
	reader := Reader{bufio.NewReader(f), 0}
	// read the LiME file and create the slice of blocks
	blocks := make([]limeBlock, 0)
	// scan the file using the bufio reader looking for LiME blocks
	for {
		b, err := reader.r.Peek(len(LIME_MAGIC))
		if err != nil {
			break
		}

		if bytes.Equal(b, LIME_MAGIC) {
			var candidate limeHeader

			err := binary.Read(reader.r, binary.LittleEndian, &candidate)
			if err != nil {
				return nil, err
			}
			// some sanity checks on the candidiate block
			if candidate.Version == 1 && candidate.EndAddress > candidate.StartAddress && allZero(candidate.Reserved[:]) {
				block := limeBlock{fileOffset: reader.pos + LIME_HEADER_SIZE, header: candidate}
				blocks = append(blocks, block)
				// jump to the next block, raw data between headers, can skip using the size
				nextBlock := int64(reader.pos) + int64(LIME_HEADER_SIZE) + int64(candidate.EndAddress-candidate.StartAddress)
				reader.pos = uint64(nextBlock)
				_, err := f.Seek(nextBlock, io.SeekStart)
				if err != nil {
					return nil, err
				}
				reader.r.Reset(f)
			}
		}
		// slide forward
		_, _ = reader.r.ReadByte()
		reader.pos++

	}

	l := LiME{file: f, blocks: blocks}
	return &l, nil

}

func (l *LiME) DebugInfo(w io.Writer) {
	// helper function to print some metadata about the intialized lime reader
	io.WriteString(w, "---- LiME Reader Info ----\n")
	io.WriteString(w, fmt.Sprintf("Blocks: %d\n", len(l.blocks)))
	t := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	io.WriteString(t, "BLOCK\tSTART\tEND\tFILE POS\tSIZE\n")
	for i, b := range l.blocks {
		io.WriteString(t, fmt.Sprintf("%d\t0x%x\t0x%x\t%d\t0x%x\n",
			i, b.header.StartAddress, b.header.EndAddress, b.fileOffset, b.header.EndAddress-b.header.StartAddress))
	}
	t.Flush()
}

func New(path string) (*LiME, error) {
	// main function used to create a LiME reader
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	l, err := initLiME(f)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *LiME) Close() error {
	// close the file, this is so users can use defer l.Close()
	return l.file.Close()
}

func (l *LiME) lookupOwningBlock(address uint64) (int, error) {
	if len(l.blocks) == 0 {
		return 0, errors.New("no blocks to check")
	}
	// FOR NOW - just scan through the blocks. In the future, make this much faster with binary search or a hashmap
	for i, b := range l.blocks {
		if b.header.StartAddress <= address && b.header.EndAddress > address {
			return i, nil
		}
	}
	return 0, fmt.Errorf("no owning block found for address: 0x%x", address)
}

func (l *LiME) Read(address uint64, size int) ([]byte, error) {
	// logic to quickly lookup which block has this address here
	idx, err := l.lookupOwningBlock(address)
	if err != nil {
		return nil, err
	}
	// get the owning block
	block := l.blocks[idx]
	if address+uint64(size) > block.header.EndAddress {
		return nil, fmt.Errorf("read of size %d at 0x%x would exceed block boundary at 0x%x",
			size, address, block.header.EndAddress)
	}
	// the specific offset into the file is the offset into the specific lime block
	offset := int64(address - block.header.StartAddress + block.fileOffset)
	// allocate a buffer to hold read data
	data := make([]byte, size)
	// can now read the actual data as if from a flat memory capture
	_, err = l.file.ReadAt(data, offset)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func allZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}
