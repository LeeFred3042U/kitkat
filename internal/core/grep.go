package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

/*
Binary detection (NUL byte heuristic – similar to git)
*/
func isBinary(data []byte) bool {
	return bytes.Contains(data, []byte{0})
}

/*
Decode file contents:
- UTF-8
- UTF-16 LE (Windows default)
- UTF-16 BE
*/
func decodeText(data []byte) (string, bool) {
	// UTF-8
	if utf8.Valid(data) {
		return string(data), true
	}

	// UTF-16 LE BOM (FF FE)
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		u16 := make([]uint16, 0, (len(data)-2)/2)
		for i := 2; i+1 < len(data); i += 2 {
			u16 = append(u16, uint16(data[i])|uint16(data[i+1])<<8)
		}
		return string(utf16ToRunes(u16)), true
	}

	// UTF-16 BE BOM (FE FF)
	if len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF {
		u16 := make([]uint16, 0, (len(data)-2)/2)
		for i := 2; i+1 < len(data); i += 2 {
			u16 = append(u16, uint16(data[i+1])|uint16(data[i])<<8)
		}
		return string(utf16ToRunes(u16)), true
	}

	return "", false
}

func utf16ToRunes(u16 []uint16) []rune {
	r := make([]rune, len(u16))
	for i, v := range u16 {
		r[i] = rune(v)
	}
	return r
}

/*
kitcat grep implementation
*/
func Grep(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: kitcat grep [--line-number] <pattern>")
	}

	showLineNumber := false
	i := 0

	if args[0] == "--line-number" {
		showLineNumber = true
		i++
	}

	if len(args[i:]) < 1 {
		return fmt.Errorf("usage: kitcat grep [--line-number] <pattern>")
	}

	pattern := args[i]

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	// Load KitCat index
	entries, err := LoadIndex()
	if err != nil {
		return err
	}

	// Deterministic output order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	for _, entry := range entries {
		data, err := os.ReadFile(entry.Path)
		if err != nil {
			continue
		}

		if isBinary(data) {
			continue
		}

		text, ok := decodeText(data)
		if !ok {
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(text))
		lineNo := 1

		for scanner.Scan() {
			line := scanner.Text()

			if re.MatchString(line) {
				if showLineNumber {
					fmt.Printf("%s:%d:%s\n", entry.Path, lineNo, line)
				} else {
					fmt.Printf("%s:%s\n", entry.Path, line)
				}
			}
			lineNo++
		}
	}

	return nil
}
