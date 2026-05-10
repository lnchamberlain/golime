## GoLiME

An implementation of a reader for LiME format memory captures in Go. The Linux Memory Extractor format (LiME) creates LiME blocks that record the starting and ending physical addresses captured in each block. This is rad compared to the raw memory capture format because some regions of physical memory are not readable (reserved for device I/O and all sorts of other reasons), so there are 'holes' in the capture that either need to be padded or otherwise accounted for. Luckily, LiME skips that by recording this start and stop info at capture time, love it!


##### LiME Specs

```
typedef struct {
    unsigned int magic;        // Always 0x4C694D45 (LiME)
    unsigned int version;      // Header version number
    unsigned long long s_addr; // Starting address of range
    unsigned long long e_addr; // Ending address of range
    unsigned char reserved[8]; // Currently all zeros
} __attribute__ ((__packed__)) lime_mem_range_header;
```

### Usage 

To use this tool first import the library:
```
go get github.com/lnchamberlain/golime
```
Then add to your imports:
```
import (
    "github.com/lnchamberlain/golime"
)
```

Now you can create a LiME object that implments `.Read()`:
```
lime, err := golime.New("./test.lime")
```
Close the reader and underlying file:
```
defer lime.Close()
````
Print debug info like number of blocks and start/stop for each 
```
lime.DebugInfo(os.Stdout)
```
Example of that debug info: 
```
---- LiME Reader Info ----
Blocks: 2
BLOCK  START     END         FILE POS  SIZE
0      0x1000    0x9fbff     32        0x9ebff
1      0x100000  0x7ffeffff  650304    0x7feeffff
```
Read a physical address, accepts a physical address and a size
```
var testPhysAddress uint64 = 0x1000
data, err := lime.Read(testPhysAddress, 8)
```