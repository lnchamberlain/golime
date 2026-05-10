package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/lnchamberlain/golime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Missing target file. Usage: test_golime /path/to/file.lime")
		os.Exit(1)
	}
	testFile := os.Args[1]
	// create a new LiME file
	lime, err := golime.New(testFile)
	if err != nil {
		fmt.Printf("Error creating LiME reader: %s\n", err.Error())
		os.Exit(1)
	}
	// close will close the underlying file
	defer lime.Close()
	// debug info prints stuff like number of blocks and start/stop for each
	lime.DebugInfo(os.Stdout)
	// example of reading a physical address from the LiME file
	var testPhysAddress uint64 = 0x1000
	data, err := lime.Read(testPhysAddress, 8)
	if err != nil {
		fmt.Printf("error reading data: %s\n", err.Error())
	}
	fmt.Printf("Data read from physical offset 0x%x:\n", testPhysAddress)
	fmt.Println(hex.Dump(data))
}
