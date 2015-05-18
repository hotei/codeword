// codeword.go (c) 2013-2015 David Rook

/*

	typical usage: codeword | tee codeword-results.txt

*/

package main

import (
	// standard lib pkgs
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	//
	"github.com/hotei/mdr"
)

type VerboseType bool

var (
	Verbose VerboseType
)

func (v VerboseType) Printf(s string, a ...interface{}) {
	if v {
		fmt.Printf(s, a...)
	}
}

type targetRec struct {
	text  []byte
	count int
}

const (
	G_version = "codeword.go   (c) 2013-2015 David Rook version 0.0.3"

	// these could be flags
	g_foldASCII = false
	g_testMode  = true
	g_maxBlock  = int64(4000000)  // can set this to limiting value for test only
)

var (
	doVerbose    bool
	doVersion    bool
	g_targets    []targetRec
	g_hitMap     map[string]int
	g_alphaTable []byte
	g_blksRead   int64
)

func init() {
	// flag setup
	flag.BoolVar(&doVerbose, "verbose", false, "use Verbose messages")
	flag.BoolVar(&doVersion, "version", false, "print version and quit")

	g_alphaTable = make([]byte, 256, 256)

	// in some environments it makes sense to put lowest classification at top so
	// search will run microscopicly faster - but only if you bail after first match
	// we don't bail so makes no difference. Alphabetic order is as good as any.
	g_targets = []targetRec{
		{text: []byte("N O F O R N")},
		{text: []byte("T O P  S E C R E T")},
		{text: []byte("(TS)")},
	}
	for ndx, target := range g_targets {
		fmt.Printf("Target [%d] = %q\n", ndx, target)
	}

	/* some obvious targets and potential problems with a few of them
	//	"T O P  S E C R E T", not needed - will get same hit as S E C R E T
		"TOP SECRET",
		"(TS)",
	// "SECRET",  avoid - matches SECRETARY
		"S E C R E T",
	//	"(S)",  avoid - occurs in object code frequently
		"CONFIDENTIAL",
		"C O N F I D ",
	//	"(C)",  avoid - also the copyright symbol
		"C L A S S I F",
		"CLASSIFIED",
		"N O F O R N"
		"NOFORN" }
	*/

	for i := byte('a'); i <= byte('z'); i++ {
		g_alphaTable[i] = i
	}
	for i := byte('A'); i <= byte('Z'); i++ {
		g_alphaTable[i] = i
	}
	g_alphaTable[byte(' ')] = ' '
	// add the characters from targets - mostly catches punctuation
	for _, targ := range g_targets {
		str := targ.text
		for _, b := range str {
			if b == ' ' {
				continue
			}
			g_alphaTable[b] = b
			Verbose.Printf("%c\n", b)
		}
	}

	Verbose.Printf("%d %v\n", len(g_alphaTable), g_alphaTable)
	g_hitMap = make(map[string]int)
}

// remove punctuation and decimal digits
// optionally fold ascii > 127 (early word processors did this)
// WordStar mucked with high bits for some reason I can't recall at present

func scrub(b []byte) []byte {
	rv := make([]byte, len(b), len(b))
	copy(rv, b)
	for i := 0; i < len(b); i++ {
		val := b[i]
		if g_foldASCII {
			val &= 0x7f
		}
		if g_alphaTable[val] == 0 {
			rv[i] = ' '
		}
	}
	return rv
}

func scanBlock(blk []byte) {
	blen := len(blk)
	if blen <= 0 {
		return
	}
	// need to keep original intact if we have to print it later
	newBlock := scrub(blk)
	Verbose.Printf("len(%d) -> %s\n", len(newBlock), string(newBlock))
	for ndx, target := range g_targets {
		Verbose.Printf("search for %v\n", string(target.text))
		if bytes.Contains(newBlock, target.text) {
			g_targets[ndx].count++
			sha := mdr.BufSHA256(blk)
			_, found := g_hitMap[sha]
			if found {
				g_hitMap[sha]++
			} else {
				g_hitMap[sha] = 1
				fmt.Printf("\nAt block %d found %s \n", g_blksRead, target.text)
				fmt.Printf("%s\n", string(blk)) // only print first hit
			}
		} else {
			Verbose.Printf("No hits\n")
		}
	}
}

func dumpHitMap() {
	for ndx, hits := range g_hitMap {
		fmt.Printf("%s had %d hits\n", ndx, hits)
	}
	for ndx, target := range g_targets {
		fmt.Printf("Target[%d] %s has %d hits\n", ndx, target.text, target.count)
	}
}

func main() {
	flag.Parse()
	if doVerbose {
		Verbose = VerboseType(true)
	}
	if doVersion {
		fmt.Printf("%s\n", G_version)
		os.Exit(0)
	}
	startTime := time.Now()
	diskName := "/dev/sda"
	input, err := os.Open(diskName)
	if err != nil {
		log.Fatalf("Cant open %s for read\n", diskName)
	}
	// how big is the disk?
	dlen, err := input.Seek(0, 2)
	if err != nil {
		log.Fatalf("cant get size of disk %s\n", diskName)
	}
	fmt.Printf("size = %d\n", dlen)
	_, err = input.Seek(0, 0)
	if err != nil {
		log.Fatalf("cant get size of disk %s\n", diskName)
	}
	Verbose.Printf("open on disk(%v) returned err(%v)\n", input, err)
	// BUG(mdr): TODO - tune BufSize to match OS block - how?
	BufSize := 512
	inBuf := make([]byte, BufSize)
	dBlks := dlen >> 9
	fmt.Printf("Disk %s has %s blocks (of 512)\n", diskName, mdr.CommaFmtInt64(dBlks))
	g_blksRead = 0
	var barA *mdr.ProgStateT
	if g_testMode {
		barA = mdr.OneProgressBar(g_maxBlock)
	} else {
		barA = mdr.OneProgressBar(dBlks)
	}
	for {
		nRead, err := input.Read(inBuf)
		if nRead != BufSize {
			if (nRead == 0) && (err == io.EOF) { // normal end of file
				break
			} else {
				log.Panicf("bad read nRead(%d) err(%v)\n", nRead, err)
			}
		}
		g_blksRead++
		barA.Update(g_blksRead)
		barA.Tag(fmt.Sprintf("%d of %d have been done", g_blksRead, dBlks))
		// BUG(mdr): early exit for test only
		if g_testMode {
			if g_blksRead >= g_maxBlock {
				fmt.Printf("\nQuitting after %d blocks in test mode\n", g_maxBlock)
				break
			}
		}
		scanBlock(inBuf)
	}
	barA.Stop()
	fmt.Printf("\n\n")
	dumpHitMap()
	elapsedTime := time.Now().Sub(startTime)
	elapsedSeconds := elapsedTime.Seconds()
	bytesRead := g_blksRead * int64(BufSize)
	readRate := float64(bytesRead) / elapsedSeconds
	fmt.Printf("Read %s bytes from disk at %.2f MB/sec\n", mdr.CommaFmtInt64(bytesRead), readRate/1e6)
	fmt.Printf("<fini>\n")
}
