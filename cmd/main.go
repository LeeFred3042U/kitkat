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
	"rm": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat rm <file>")
			return
		}
		filename := args[0]
		if err := core.RemoveFile(filename); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Removed '%s'\n", filename)
	},
	"commit": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat commit <-m | -am | --amend> <message>")
			return
		}

		var message string
		var isAmend bool

		// Check for --amend flag
		if args[0] == "--amend" {
			if len(args) < 3 || args[1] != "-m" {
				fmt.Println("Usage: kitkat commit --amend -m <message>")
				return
			}
			isAmend = true
			message = strings.Join(args[2:], " ")
		} else if len(args) < 2 {
			fmt.Println("Usage: kitkat commit <-m | -am> <message>")
			return
		} else {
			// Normal commit flow
			switch args[0] {
			case "-am":
				message = strings.Join(args[1:], " ")
				newCommit, summary, err := core.CommitAll(message)
				if err != nil {
					if err.Error() == "nothing to commit, working tree clean" {
						fmt.Println(err.Error())
					} else {
						fmt.Println("Error:", err)
					}
					return
				}
				printCommitResult(newCommit, summary)
				return
			case "-m":
				message = strings.Join(args[1:], " ")
			default:
				fmt.Println("Usage: kitkat commit <-m | -am | --amend> <message>")
				return
			}
		}

		// Handle amend or normal commit
		if isAmend {
			newCommit, err := core.AmendCommit(message)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			headState, err := core.GetHeadState()
			if err != nil {
				headData, _ := os.ReadFile(".kitkat/HEAD")
				ref := strings.TrimSpace(string(headData))
				headState = strings.TrimPrefix(ref, "ref: refs/heads/")
			}
			fmt.Printf("[%s %s] %s (amended)\n", headState, newCommit.ID[:7], newCommit.Message)
		} else {
			newCommit, summary, err := core.Commit(message)
			if err != nil {
				if err.Error() == "nothing to commit, working tree clean" {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Error:", err)
				}
				return
			}
			printCommitResult(newCommit, summary)
		}
	},
	"log": func(args []string) {
		oneline := false
		limit := -1
		i := 0
		for i < len(args) {
			switch args[i] {
			case "--oneline":
				oneline = true
				i++
			case "-n":
				if i+1 >= len(args) {
					fmt.Println("Error: -n requires a positive integer argument")
					return
				}
				var n int
				_, err := fmt.Sscanf(args[i+1], "%d", &n)
				if err != nil || n <= 0 {
					fmt.Println("Error: -n requires a positive integer argument")
					return
				}
				limit = n
				i += 2
			default:
				fmt.Printf("Error: unknown flag %s\n", args[i])
				return
			}
		}
		if err := core.ShowLog(oneline, limit); err != nil {
			fmt.Println("Error:", err)
		}
	},
	"shortlog": func(args []string) {
		if err := core.ShowShortLog(); err != nil {
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
	"reset": func(args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: kitkat reset --hard <commit-hash>")
			return
		}
		if args[0] != "--hard" {
			fmt.Println("Error: only 'reset --hard' is currently supported")
			fmt.Println("Usage: kitkat reset --hard <commit-hash>")
			return
		}
		if err := core.ResetHard(args[1]); err != nil {
			fmt.Println("Error:", err)
		}
	},
	"ls-files": func(args []string) {
		entries, err := core.LoadIndex()
		if err != nil {
			fmt.Println("Error loading index:", err)
			return
		}

		for _, entry := range entries {
			fmt.Println(entry.Path)
		}
	},
	"clean": func(args []string) {
		force := false
		includeIgnored := false

		for _, arg := range args {
			if arg == "-f" {
				force = true
			} else if arg == "-x" {
				includeIgnored = true
			}
		}

		if force {
			if err := core.Clean(false, includeIgnored); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("This will delete untracked files. Run 'kitkat clean -f' to proceed.")
			if err := core.Clean(true, includeIgnored); err != nil {
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
		if len(args) == 1 && (args[0] == "--list") {
			if err := core.PrintTags(); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}

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
			if len(args) == 1 && args[0] == "--list" {
				if err := core.PrintAllConfig(); err != nil {
					fmt.Println("Error:", err)
				}
				return
			}
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
	"show-object": func(args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: kitkat show-object <hash>")
			return
		}
		if err := core.ShowObject(args[0]); err != nil {
			fmt.Println("Error:", err)
		}
	},
	"branch": func(args []string) {
		if args[0] == "-l" {
			if err := core.ListBranches(); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
		if args[0] == "-r" {
			name := args[1]
			if err := core.RenameCurrentBranch(name); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
		name := args[0]
		if core.IsBranch(name) {
			fmt.Printf("Error: Branch '%s' already exists\n", name)
			return
		}
		if err := core.CreateBranch(name); err != nil {
			fmt.Println("Error:", err)
			return
		}
	},
}

// printCommitResult formats and prints the commit result with summary
func printCommitResult(newCommit models.Commit, summary string) {
	headState, err := core.GetHeadState()
	if err != nil {
		headData, _ := os.ReadFile(".kitkat/HEAD")
		ref := strings.TrimSpace(string(headData))
		headState = strings.TrimPrefix(ref, "ref: refs/heads/")
	}
	fmt.Printf("[%s %s] %s\n%s\n", headState, newCommit.ID[:7], newCommit.Message, summary)
}

func main() {
	if len(os.Args) >= 4 && os.Args[1] == "branch" && (os.Args[2] == "-m" || os.Args[2] == "--move") {
		newName := os.Args[3]
		err := core.RenameCurrentBranch(newName)
		if err != nil {
			fmt.Println("Error renaming branch:", err)
			os.Exit(1)
		}
		fmt.Println("Branch renamed to", newName)
		return
	}
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
