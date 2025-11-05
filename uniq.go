package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Flags struct {
	C bool
	D bool
	U bool
	F int
	S int
	I bool
}

type LineInfo struct {
	Count    int
	Index    int
	Original string
}

func processCommand(f Flags, d []string, output io.Writer) {
	counts := make(map[string]*LineInfo)
	dCopy := append([]string{}, d...)

	if f.I {
		for i, el := range dCopy {
			dCopy[i] = strings.ToLower(el)
		}
	}
	if f.F > 0 {
		for i, el := range dCopy {
			var newString string
			stringSplit := strings.Split(el, " ")
			if f.F < len(stringSplit) {
				newString = strings.Join(stringSplit[f.F:], " ")
			} else {
				newString = ""
			}

			if newString == "" {
				newString = "" + el
			}
			dCopy[i] = newString
		}
	}
	if f.S > 0 {
		for i, el := range dCopy {
			if len(el) > f.S {

				dCopy[i] = el[f.S:]
			} else {
				dCopy[i] = ""
			}
		}

	}
	meeting := make([]int, len(d))
	for i, el := range dCopy {
		if value, exists := counts[el]; exists {
			value.Count++
		} else {
			counts[el] = &LineInfo{Count: 1, Index: i, Original: d[i]}
		}
		meeting[counts[el].Index] = counts[el].Count
	}
	switch {
	case f.C:
		for i, value := range meeting {
			if value >= 1 {
				fmt.Fprintln(output, value, d[i])
			}

		}
	case f.D:
		for i, value := range meeting {
			if value > 1 {
				fmt.Fprintln(output, d[i])
			}
		}
	case f.U:
		for i, value := range meeting {
			if value == 1 {
				fmt.Fprintln(output, d[i])
			}
		}
	default:
		for i, value := range meeting {
			if value >= 1 {
				fmt.Fprintln(output, d[i])
			}
		}

	}

}

func processData(r io.Reader) []string {
	var lines []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

	}
	return lines
}

func main() {
	var commandFlags Flags
	flag.BoolVar(&commandFlags.C, "c", false, "print count string encounters with strings")
	flag.BoolVar(&commandFlags.D, "d", false, "print dublicated strings")
	flag.BoolVar(&commandFlags.U, "u", false, "print unique strings")
	flag.IntVar(&commandFlags.F, "f", 0, "skip n fields")
	flag.IntVar(&commandFlags.S, "s", 0, "skip n chars")
	flag.BoolVar(&commandFlags.I, "i", false, "ignore case")

	flag.Parse()
	activated := 0
	if commandFlags.C {
		activated++
	}
	if commandFlags.D {
		activated++
	}
	if commandFlags.U {
		activated++
	}

	if activated > 1 {
		fmt.Fprintln(os.Stderr, "Error: use only one of -c, -d, -u flags")
		flag.Usage()
		os.Exit(1)
	}

	var reader io.Reader
	if len(flag.Args()) > 0 {
		inputFile, err := os.Open(flag.Args()[0])
		if err != nil {
			log.Fatal(err)
		}
		defer inputFile.Close()
		reader = inputFile
	} else {
		input := os.Stdin
		reader = input
	}
	data := processData(reader)
	processCommand(commandFlags, data, os.Stdout)
}
