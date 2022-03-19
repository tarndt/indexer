package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

/*
	This is a sample program for teaching that takes a text input file and returns
	a book style index of the words contained and what pages they occur on.
*/

const (
	defLinesPerPage = 106
)

func main() {
	//Get command line arguments
	fin, fout, linesPerPage := getArgs()
	defer func() {
		fin.Close()
		fout.Close()
	}()

	//Setup our scanner
	lines := bufio.NewScanner(fin)
	lines.Split(bufio.ScanLines)

	//Setup our index data structure map of words to map of pages
	idx := make(map[string]map[int]struct{})
	//Setup out line and page counters
	lineNum, pageNum := 0, 1

	//Read each line
	for lines.Scan() {
		//Increment line or page as needed
		if lineNum++; lineNum > linesPerPage {
			pageNum++
			lineNum = 1
		}

		//Split the line into words
		for _, wordBytes := range bytes.Fields(lines.Bytes()) {
			//Normalize each word and check if it was throw out
			word := cleanWord(string(wordBytes))
			if word == "" {
				continue
			}

			//Add a new page entry for the word, allocating an entry map if needed
			pages, exists := idx[word]
			if !exists {
				pages = make(map[int]struct{})
			}
			pages[pageNum] = struct{}{}
			idx[word] = pages
		}
	}
	//Did we have errors while processing?
	if err := lines.Err(); err != nil {
		log.Fatalf("ERROR: Failed to read line from input file: %s", err)
	}

	//Create a sorted slice of all words
	wordlist := make([]string, 0, len(idx))
	for word := range idx {
		wordlist = append(wordlist, word)
	}
	sort.Strings(wordlist)

	//Wrap out output file with buffering
	bufOut := bufio.NewWriter(fout)
	defer bufOut.Flush()

	//Output each word
	var pageList []int
	var lineBuf bytes.Buffer
	for _, word := range wordlist {
		pageList := pageList[:0] //Reset our page list slice for resuse
		//Create a sorted slice of all pages
		for page := range idx[word] {
			pageList = append(pageList, page)
		}
		sort.Ints(pageList)

		//Write this words output to our line buffer
		//Format: <WORD>:<tab>1,2,3,4...<newline>
		lineBuf.Reset()
		lineBuf.WriteString(word)
		lineBuf.WriteByte(':')
		lineBuf.WriteByte('\t')
		for _, page := range pageList {
			lineBuf.WriteString(strconv.Itoa(page))
			lineBuf.WriteByte(',')
		}
		lineBuf.Truncate(lineBuf.Len() - 1)
		lineBuf.WriteByte('\n')

		//Write out output line buffer to our buffered output file
		if _, err := bufOut.Write(lineBuf.Bytes()); err != nil {
			log.Fatalf("ERROR: Failed to write line to output file: %s", err)
		}
	}
}

func getArgs() (fin, fout *os.File, lpp int) {
	fin, fout = os.Stdin, os.Stdout

	//Parse command line arguments
	var linesPerPage uint
	var finName, foutName string
	flag.StringVar(&finName, "input", "", "File to read (empty implies stdin)")
	flag.StringVar(&foutName, "output", "", "File to write (empty implies stdout)")
	flag.UintVar(&linesPerPage, "lines-per-page", defLinesPerPage, "Number of lines of text per logical page")
	flag.Parse()

	//Open files as needed
	var err error
	if finName != "" {
		if fin, err = os.Open(finName); err != nil {
			log.Fatalf("Could not open input file %q: %s", finName, err)
		}
	}
	if foutName != "" {
		if fout, err = os.Create(foutName); err != nil {
			log.Fatalf("Could not create output file %q: %s", foutName, err)
		}
	}

	return fin, fout, int(linesPerPage)
}

func cleanWord(str string) string {
	//Remove any leading and trailing characters that are not unicode letters
	str = strings.TrimFunc(strings.ToLower(str), func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	//If we have any internal characters thata are not in the range a-z with the
	// exception of a apostrophe reject it!
	for _, char := range str {
		if char < 'a' || char > 'z' {
			switch char {
			case '\'', 8217: //8217 is unicode '’'
				continue
			}
			return ""
		}
	}

	return str
}
