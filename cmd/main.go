package main

import (
	"fmt"
	"os"

	"github.com/LeeFred3042U/kitkat/internal/core"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

func main() {
	// Handle help flags
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		core.PrintGeneralHelp()
		os.Exit(0)
	}

	command := os.Args[1]
	if command == "help" {
		if len(os.Args) < 3 {
			core.PrintGeneralHelp()
		} else {
			// Handle requests like "kitkat help add"
			core.PrintCommandHelp(os.Args[2])
		}
		os.Exit(0)
	}

	// dispatch block
	switch command {
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
			fmt.Println("Usage: kitkat tag <tag-name> <commit-id>")
			return
		}
		tagName := os.Args[2]
		commitID := os.Args[3]
		if err := core.CreateTag(tagName, commitID); err != nil {
			fmt.Println("Error creating tag:", err)
		}

	case "grep":
		if len(os.Args) < 3 {
			fmt.Println("Usage: kitkat grep <term>")
			return
		}
		if err := core.GrepLogs(os.Args[2]); err != nil {
			fmt.Println("Error searching logs:", err)
		}

	case "commit":
		if len(os.Args) < 4 || os.Args[2] != "-m" {
			fmt.Println("usage: kitkat commit -m <message>")
			os.Exit(1)
		}
		message := os.Args[3]
		commitID, err := core.Commit(message)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Printf("committed with ID %s\n", commitID)

	case "clean":
		if err := core.Clean(); err != nil {
			fmt.Println("Error cleaning repository:", err)
		}
		
	default:
		fmt.Println("unknown command:", command)
	}
}