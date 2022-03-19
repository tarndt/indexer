package main

import (
	"os"
	"strconv"
	"testing"
)

func TestGetArgs(t *testing.T) {
	/*fin, fout, linesPerPage := getArgs()
	switch {
	case fin != os.Stdin:
		t.Fatalf("fin should have been stdin")
	case fout != os.Stdout:
		t.Fatalf("fout should have been stdout")
	case linesPerPage != defLinesPerPage:
		t.Fatalf("linesPerPage should have been %d", defLinesPerPage)
	}*/

	const testInputFile = "testInput.txt"
	fin, err := os.Create(testInputFile)
	if err != nil {
		t.Fatalf("Coud not create test input file: %s", err)
	}
	defer func() {
		fin.Close()
		os.Remove(testInputFile)
	}()

	const (
		testOutputFile = "testOutput.txt"
		testNumLines   = defLinesPerPage + 1
	)
	os.Args = []string{os.Args[0], "-input=" + testInputFile, "-output=" + testOutputFile, "-lines-per-page=" + strconv.Itoa(testNumLines)}
	fin, fout, linesPerPage := getArgs()
	defer func() {
		fout.Close()
		os.Remove(testOutputFile)
	}()

	switch {
	case fin == os.Stdin && fin == nil:
		t.Fatalf("fin should have been: %q", testInputFile)
	case fout == os.Stdout && fout == nil:
		t.Fatalf("fout should have been: %q", testOutputFile)
	case linesPerPage != testNumLines:
		t.Fatalf("linesPerPage should have been %d but was %d", testNumLines, linesPerPage)
	}
}
