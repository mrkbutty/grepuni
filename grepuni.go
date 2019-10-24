/*
grepuni
=======

Will search unicode utf-16 text files (such as those output by windows regedit).

Usage: grepuni [OPTIONS] <filename> <regexp>

	<filename> 		= File to search.
	<regex> 		= Regular expression
	-P <regex> 		= specifies end of paragraph to stop output

Author: Mark Butterworth Oct 2019
*/

package main


import (
	"flag"
    "bufio"
    "fmt"
    "log"
	"os"
	"regexp"

    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
)

var flagParagraph string
var flagVerbose bool
var flagQuiet bool
var flagTestOnly bool

type utfScanner interface {
    Read(p []byte) (n int, err error)
}

// Creates a scanner similar to os.Open() but decodes the file as UTF-16.
// Useful when reading data from MS-Windows systems that generate UTF-16BE
// files, but will do the right thing if other BOMs are found.
func NewScannerUTF16(filename string) (utfScanner, error) {

    // Read the file into a []byte:
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }

    // Make an tranformer that converts MS-Win default to UTF8:
    win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
    // Make a transformer that is like win16be, but abides by BOM:
    utf16bom := unicode.BOMOverride(win16be.NewDecoder())

    // Make a Reader that uses utf16bom:
    unicodeReader := transform.NewReader(file, utf16bom)
    return unicodeReader, nil
}

func main() {
	flagParagraph := flag.String("P", "", "Search up to end of paragraph string")
	flag.BoolVar(&flagVerbose, "v", false, "Prints detailed operations")
	flag.BoolVar(&flagQuiet, "q", false, "No output apart from errors")
	flag.BoolVar(&flagTestOnly, "t", false, "Test only do not change")
	flag.Parse()

	//items := []string{"."}  // default arguments to use if omitted

	if flag.NArg() != 2 {
		flag.PrintDefaults()
		return
	}

	args := flag.Args()
	paraflag := true
	if *flagParagraph == "" { paraflag = false }

	greparg := regexp.MustCompile(args[1])
	grepend := regexp.MustCompile(*flagParagraph)

	s, err := NewScannerUTF16(args[0])
    if err != nil {
        log.Fatal(err)
    }
	
	found := false
    scanner := bufio.NewScanner(s)
    for scanner.Scan() {
		match := greparg.MatchString(scanner.Text())
		if match { 
			found = true
			fmt.Println(scanner.Text())
			continue
		}
		if found && paraflag {
			match := grepend.MatchString(scanner.Text())
			if match { found = false }
		}
		if found {
			fmt.Println(scanner.Text())
			if paraflag == false { found = false }
		}


	}

}


