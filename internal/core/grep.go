package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func isBinary(data []byte) bool {
	return bytes.Contains(data, []byte{0})
}

func Grep(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: kitkat grep [--line-number] <pattern>")
	}

	showLineNumber := false
	i := 0

	if args[0] == "--line-number" {
		showLineNumber = true
		i++
	}

	if len(args[i:]) < 1 {
		return fmt.Errorf("usage: kitkat grep [--line-number] <pattern>")
	}

	pattern := args[i]

	// ✅ compile regex ONCE
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	// Get tracked files only
	out, err := exec.Command("git", "ls-files").Output()
	if err != nil {
		return fmt.Errorf("not a git repository")
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// skip binaries silently
		if isBinary(data) {
			continue
		}

		scanner := bufio.NewScanner(bytes.NewReader(data))
		lineNo := 1

		for scanner.Scan() {
			line := scanner.Text()

			if re.MatchString(line) {
				if showLineNumber {
					fmt.Printf("%s:%d:%s\n", file, lineNo, line)
				} else {
					fmt.Printf("%s:%s\n", file, line)
				}
			}
			lineNo++
		}
	}

	return nil
}
