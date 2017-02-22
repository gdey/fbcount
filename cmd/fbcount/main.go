/*
Simple application that will read a file, and print out the absolute filename and the byte position for the provide line:char position.
Input should look like:
$filename#$line:$charpos
*/

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

func splitFunc(fn string) (fname string, line, char int) {
	var err error
	idx := strings.Index(fn, "#")
	// So, we don't have a line number and char, we will assume it's the first line and char.
	if idx == -1 {
		if fname, err = filepath.Abs(fn); err != nil {
			log.Println("Got the following error", err)
			return fn, 0, 0
		}
		return fname, 0, 0
	}
	// Get the function name from the start of the string to the idx value.
	if fname, err = filepath.Abs(fn[:idx]); err != nil {
		log.Println("Got the following error", err)
		fname = fn[:idx]
	}

	strs := strings.SplitN(fn[idx+1:], ":", 2)
	switch len(strs) {
	case 2:
		// Don't care about the error, want the value to be zero value if we are unable to convert the string to an number.
		char, _ = strconv.Atoi(strs[1])
		fallthrough
	case 1:
		// Don't care about the error, want the value to be zero value if we are unable to convert the string to an number.
		line, _ = strconv.Atoi(strs[0])
		line--
		if line < 0 {
			line = 0
		}
	}
	return fname, line, char
}

var ErrPosNotFound = errors.New("Line Char position not found.")

func bytePos(fname string, line, char int) (byteCount int, err error) {
	if line == 0 && char == 0 {
		return 0, nil
	}
	file, err := os.Open(fname)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	var lineCount, charCount int
	idx := 0
	bytes := make([]byte, 1024)
	next := true
	for next {
		n, err := file.Read(bytes[idx:])
		if err != nil && err != io.EOF {
			return 0, err
		}
		if err == io.EOF {
			next = false
		}
		b := bytes[:n]
		for len(b) > 0 && utf8.FullRune(b) {
			r, size := utf8.DecodeRune(b)
			b = b[size:]
			byteCount += size
			charCount++
			if r == '\n' {
				charCount = 0
				lineCount++
			}
			if lineCount >= line && charCount >= char {
				return byteCount, nil
			}
		}
		for i := range b {
			bytes[i] = b[i]
		}
		idx = len(b)
	}
	return byteCount, ErrPosNotFound
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Need file and line position.")
		os.Exit(1)
	}
	fname, line, char := splitFunc(os.Args[1])
	byteCount, _ := bytePos(fname, line, char)
	fmt.Printf("%v:#%v", fname, byteCount)
}
