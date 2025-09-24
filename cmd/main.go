package main

import (
	"os"
	"fmt"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/core"
)

type CommandFunc func(args []string)

// Command registry
var commands = map[string]CommandFunc{
	"init": func(args []string) {
		if err := core.InitRepo(); err != nil {
			fmt.Println("Error:", err)
		}
	},

	"add": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat add <file-path>")
			return
		}

		// Check for the "-A" or "-all" flag
		if args[0] == "-A" || args[0] == "--all" {
			fmt.Println("Staging all changes...")
			if err := core.AddAll(); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}

		// Allow adding multiple files
		for _, path := range args {
			if err := core.AddFile(path); err != nil {
				fmt.Printf("Error adding %s: %v\n", path, err)
			}
		}
	},

	"commit": func(args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: kitkat commit <-m | -am> <message>")
			return
		}

		// Check for the combined "-am" flag.
		if args[0] == "-am" {
			message := strings.Join(args[1:], " ")
			id, err := core.CommitAll(message) // Call the new CommitAll function
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Commit created:", id)
			return
		}

		// Check for the standard "-m" flag.
		if args[0] == "-m" {
			message := strings.Join(args[1:], " ")
			id, err := core.Commit(message)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Commit created:", id)
			return
		}
		
		// If neither flag is recognized
		fmt.Println("Usage: kitkat commit <-m | -am> <message>")
	},

	"log": func(args []string) {
		oneline := false
		if len(args) > 0 && args[0] == "--oneline" {
			oneline = true
		}
		if err := core.ShowLog(oneline); err != nil {
			fmt.Println("Error:", err)
		}
	},

	"status": func(args []string) {
		if err := core.Status(); err != nil {
			fmt.Println("Error:", err)
		}
	},

	"diff": func(args []string) {
		if err := core.Diff(); err != nil {
			fmt.Println("Error:", err)
		}
	},

	"branch": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat branch <branch-name>")
			return
		}
		if err := core.CreateBranch(args[0]); err != nil {
			fmt.Println("Error:", err)
		}
	},

	"checkout": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat checkout <branch-name | file-path>")
			return
		}
		name := args[0]
		// Check if the argument is a branch or a file
		if core.IsBranch(name) {
			if err := core.CheckoutBranch(name); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			if err := core.CheckoutFile(name); err != nil {
				fmt.Println("Error:", err)
			}
		}
	},

	"merge": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat merge <branch-name>")
			return
		}
		if err := core.Merge(args[0]); err != nil {
			fmt.Println("Error:", err)
		}
	},

	"ls-files": func(args []string) {
		if err := core.ListFiles(); err != nil {
			fmt.Println("Error:", err)
		}
	},

	"clean": func(args []string) {
		// For safety, let's make it require a -f flag
		if len(args) > 0 && args[0] == "-f" {
			if err := core.Clean(false); err != nil { // false means not a dry run
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("This will delete untracked files. Run 'kitkat clean -f' to proceed.")
			if err := core.Clean(true); err != nil { // true means dry run
				fmt.Println("Error:", err)
			}
		}
	},

	"help": func(args []string) {
		if len(args) > 0 {
			core.PrintCommandHelp(args[0])
		} else {
			core.PrintGeneralHelp()
		}
	},

	"tag": func(args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: kitkat tag <tag-name> <commit-id>")
			return
		}
		if err := core.CreateTag(args[0], args[1]); err != nil {
			fmt.Println("Error:", err)
		}
	},
	"config": func(args []string) {
		if len(args) != 3 || args[0] != "--global" {
			fmt.Println("Usage: kitkat config --global <key> <value>")
			return
		}
		key := args[1]
		value := args[2]

		if err := core.SetConfig(key, value); err != nil{
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Set config: %s = %s\n", key, value)
	},
	
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kitkat <command> [args]")
		return
	}

	cmd, args := os.Args[1], os.Args[2:]
	if handler, ok := commands[cmd]; ok {
		handler(args)
	} else {
		fmt.Println("Unknown command:", cmd)
		core.PrintGeneralHelp()
	}
}