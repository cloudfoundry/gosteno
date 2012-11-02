package main

import (
	"bufio"
	"flag"
	"fmt"
	steno "github.com/cloudfoundry/gosteno"
	"io"
	"os"
)

var excludedFields = steno.EXCLUDE_NONE
var ignoreParseError = false
var helptext = `Usage: steno-prettify [OPTS] [FILE(s)]

Parses json formatted log lines from FILE(s), or stdin,
and displays a more human friendly version of each line to stdout.

Examples :
    steno-prettify FILE1 FILE2
        Prettify contents of FILE1, FILE2
    
    steno-prettify
        Prettify contents of stdin.

Options:

    -l
        Exclude File, Method and Line fields
    -d
        Exclude Data filed
    -i
        Ignore errors in parsing logs
    -help
        Shows this help
`

func prettifyFromIO(src io.Reader) {
	buf := bufio.NewReader(src)

	lineno := 0
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		lineno++
		p := steno.NewJsonPrettifier(excludedFields)
		r, e := p.DecodeJsonLogEntry(line)
		if e != nil {
			if !ignoreParseError {
				fmt.Fprintf(os.Stderr, "steno-prettify: Malformed json at line %d : %s", lineno, line)
			}
			continue
		}
		s, _ := p.EncodeRecord(r)
		fmt.Println(string(s))
	}
}

func prettifyFromFile(f *os.File) {
	defer f.Close()

	prettifyFromIO(f)
}

func main() {
	l := flag.Bool("l", false, "Exclude File, Method and Line fields")
	d := flag.Bool("d", false, "Exclude Data filed")
	i := flag.Bool("i", false, "Ignore errors in parsing logs")
	h := flag.Bool("help", false, "Show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, helptext)
	}

	flag.Parse()

	if *l {
		excludedFields |= steno.EXCLUDE_FILE | steno.EXCLUDE_LINE | steno.EXCLUDE_METHOD
	}
	if *d {
		excludedFields |= steno.EXCLUDE_DATA
	}
	if *i {
		ignoreParseError = true
	}
	if *h {
		flag.Usage()
		return
	}

	args := flag.Args()
	if len(args) > 0 {
		for _, a := range args {
			f, e := os.Open(a)
			if e != nil {
				fmt.Fprintln(os.Stderr, e)
				continue
			}
			prettifyFromFile(f)
		}
	} else {
		prettifyFromIO(os.Stdin)
	}
}
