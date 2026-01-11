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
Debug helper (enabled when KITCAT_DEBUG=1)
*/
func debugf(format string, args ...any) {
	if os.Getenv("KITCAT_DEBUG") == "1" {
		fmt.Printf("[DEBUG grep] "+format+"\n", args...)
	}
}

/*
Binary detection (NUL byte heuristic – same idea as git)
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
	debugf("entered Grep with args=%v", args)

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
	debugf("compiled pattern=%q", pattern)

	// Load KitCat index (NOT git)
	entries, err := LoadIndex()
	if err != nil {
		return err
	}
	debugf("index loaded, %d entries", len(entries))

	// Deterministic order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	for _, entry := range entries {
		debugf("scanning file=%s", entry.Path)

		data, err := os.ReadFile(entry.Path)
		if err != nil {
			debugf("failed to read file=%s: %v", entry.Path, err)
			continue
		}

		if isBinary(data) {
			debugf("skipping binary file=%s", entry.Path)
			continue
		}

		text, ok := decodeText(data)
		if !ok {
			debugf("skipping unsupported encoding file=%s", entry.Path)
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(text))
		lineNo := 1
		matched := false

		for scanner.Scan() {
			line := scanner.Text()
			debugf("line %d: %q", lineNo, line)

			if re.MatchString(line) {
				matched = true
				if showLineNumber {
					fmt.Printf("%s:%d:%s\n", entry.Path, lineNo, line)
				} else {
					fmt.Printf("%s:%s\n", entry.Path, line)
				}
			}
			lineNo++
		}

		if !matched {
			debugf("no matches in file=%s", entry.Path)
		}
	}

	debugf("grep finished successfully")
	return nil
}
