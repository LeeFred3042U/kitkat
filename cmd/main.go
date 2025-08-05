package main

import (
	"fmt"
	"os"

	"github.com/LeeFred3042U/kitkat/internal/core"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: kitkat <command> [args]")
		os.Exit(1)
	}

	// dispatch block
	switch os.Args[1] {
	case "init":
		err := core.InitRepo()
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("Initialized empty kitkat repository in .kitkat/")
	
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("usage: kitkat add <file>")
			os.Exit(1)
		}
		err := core.AddFile(os.Args[2])
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Printf("added %s\n", os.Args[2])

	case "ls-files":
		entries, err := storage.LoadIndex()
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		for path := range entries {
			fmt.Println(path)
		}

	case "log":
		if len(os.Args) < 3 {
			fmt.Println("usage: kitkat log <message>")
			os.Exit(1)
		}
		msg := os.Args[2]
		err := core.LogMessage(msg)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("log saved.")

	case "view":
		err := core.ViewLogs()
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	
	case "tag":
		if len(os.Args) < 4 {
			fmt.Println("Usage: kitkat tag <log-id> <tag>")
			return
		}
		id := os.Args[2]
		tag := os.Args[3]
		if err := core.TagLog(id, tag); err != nil {
			fmt.Println("Error tagging log:", err)
		} else {
			fmt.Println("Log tagged successfully")
		}

	case "grep":
		if len(os.Args) < 3 {
			fmt.Println("Usage: kitkat grep <term>")
			return
		}
		if err := core.GrepLogs(os.Args[2]); err != nil {
			fmt.Println("Error searching logs:", err)
		}

	case "ls-tag":
		if len(os.Args) < 3 {
			fmt.Println("Usage: kitkat ls-tag <tag>")
			return
		}
		if err := core.ListLogsByTag(os.Args[2]); err != nil {
			fmt.Println("Error:", err)
		}

	default:
		fmt.Println("unknown command:", os.Args[1])
	}
}