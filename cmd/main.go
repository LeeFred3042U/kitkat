package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/core"
	"github.com/LeeFred3042U/kitkat/internal/models"
)

type CommandFunc func(args []string)

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
		if args[0] == "-A" || args[0] == "--all" {
			fmt.Println("Staging all changes...")
			if err := core.AddAll(); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
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

		var message string
		var commitFunc func(string) (models.Commit, string, error)

		switch args[0] {
		case "-am":
			message = strings.Join(args[1:], " ")
			commitFunc = core.CommitAll
		case "-m":
			message = strings.Join(args[1:], " ")
			commitFunc = core.Commit
		default:
			fmt.Println("Usage: kitkat commit <-m | -am> <message>")
			return
		}

		newCommit, summary, err := commitFunc(message)
		if err != nil {
			if err.Error() == "nothing to commit, working tree clean" {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		headState, err := core.GetHeadState()
		if err != nil {
			// Fallback in case GetHeadState fails on the very first commit before ref exists
			headData, _ := os.ReadFile(".kitkat/HEAD")
			ref := strings.TrimSpace(string(headData))
			headState = strings.TrimPrefix(ref, "ref: refs/heads/")
		}

		fmt.Printf("[%s %s] %s\n%s\n", headState, newCommit.ID[:7], newCommit.Message, summary)
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
		if len(args) == 0 {
			if err := core.ListBranches(); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
		if err := core.CreateBranch(args[0]); err != nil {
			fmt.Println("Error:", err)
		}
	},
	"checkout": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat checkout [-b] <branch-name> | <file-path>")
			return
		}
		if args[0] == "-b" {
			if len(args) != 2 {
				fmt.Println("Usage: kitkat checkout -b <branch-name>")
				return
			}
			name := args[1]
			if core.IsBranch(name) {
				fmt.Printf("Error: Branch '%s' already exists\n", name)
				return
			}
			if err := core.CreateBranch(name); err != nil {
				fmt.Println("Error:", err)
				return
			}
			if err := core.CheckoutBranch(name); err != nil {
				fmt.Println("Error:", err)
				return
			}
			return
		}
		name := args[0]
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
		if len(args) > 0 && args[0] == "-f" {
			if err := core.Clean(false); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("This will delete untracked files. Run 'kitkat clean -f' to proceed.")
			if err := core.Clean(true); err != nil {
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
		if len(args) < 2 || args[0] != "--global" {
			fmt.Println("Usage: kitkat config --global <key> [<value>]")
			return
		}
		key := args[1]
		if len(args) == 3 {
			value := args[2]
			if err := core.SetConfig(key, value); err != nil {
				fmt.Println("Error:", err)
			}
		} else if len(args) == 2 {
			value, ok, err := core.GetConfig(key)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			if ok {
				fmt.Println(value)
			}
		} else {
			fmt.Println("Usage: kitkat config --global <key> [<value>]")
		}
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
