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
			os.Exit(1)
		}
	},
	"add": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat add <file-path>")
			os.Exit(2)
		}
		if args[0] == "-A" || args[0] == "--all" {
			fmt.Println("Staging all changes...")
			if err := core.AddAll(); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
		exitCode := 0
		for _, path := range args {
			if err := core.AddFile(path); err != nil {
				fmt.Printf("Error adding %s: %v\n", path, err)
				exitCode = 1
			}
		}
		os.Exit(exitCode)
	},
	"rm": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat rm <file>")
			os.Exit(2)
		}
		filename := args[0]
		if err := core.RemoveFile(filename); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Removed '%s'\n", filename)
		os.Exit(0)
	},
	"commit": func(args []string) {
		if !core.IsRepoInitialized() {
			fmt.Println("Error: not a kitkat repository (or any of the parent directories): .kitkat")
			os.Exit(1)
		}

		if len(args) < 2 {
			fmt.Println("Usage: kitkat commit <-m | -am | --amend> <message>")
			os.Exit(2)
		}

		var isAmend bool
		var message string

		switch args[0] {
		// Checks for amending
		case "--amend":
			if len(args) < 3 || args[1] != "-m" {
				fmt.Println("Usage: kitkat commit --amend -m <message>")
				os.Exit(2)
			}
			isAmend = true
			message = strings.Join(args[2:], " ")
		// Normal commit flow
		case "-am":
			message = strings.Join(args[1:], " ")
			newCommit, summary, err := core.CommitAll(message)
			if err != nil {
				if err.Error() == "nothing to commit, working tree clean" {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				os.Exit(2)
			}
			printCommitResult(newCommit, summary)
			os.Exit(0)
		case "-m":
			message = strings.Join(args[1:], " ")
		default:
			fmt.Println("Usage: kitkat commit <-m | -am | --amend> <message>")
			os.Exit(2)
		}

		// Handle amend or normal commit
		if isAmend {
			newCommit, err := core.AmendCommit(message)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			headState, err := core.GetHeadState()
			if err != nil {
				headData, _ := os.ReadFile(".kitkat/HEAD")
				ref := strings.TrimSpace(string(headData))
				headState = strings.TrimPrefix(ref, "ref: refs/heads/")
			}
			fmt.Printf("[%s %s] %s (amended)\n", headState, newCommit.ID[:7], newCommit.Message)
			os.Exit(0)
		} else {
			newCommit, summary, err := core.Commit(message)
			if err != nil {
				if err.Error() == "nothing to commit, working tree clean" {
					fmt.Println(err.Error())
					os.Exit(1)
				} else {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
			}
			printCommitResult(newCommit, summary)
			os.Exit(0)
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
					os.Exit(2)
				}
				var n int
				_, err := fmt.Sscanf(args[i+1], "%d", &n)
				if err != nil || n <= 0 {
					fmt.Println("Error: -n requires a positive integer argument")
					os.Exit(2)
				}
				limit = n
				i += 2
			default:
				fmt.Printf("Error: unknown flag %s\n", args[i])
				os.Exit(2)
			}
		}
		if err := core.ShowLog(oneline, limit); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
	"shortlog": func(args []string) {
		if err := core.ShowShortLog(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
	"status": func(args []string) {
		if err := core.Status(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
	"diff": func(args []string) {
		staged := false
		if len(args) > 0 {
			if args[0] == "--cached" || args[0] == "--staged" {
				staged = true
			}
		}
		if err := core.Diff(staged); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
	"checkout": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat checkout [-b] <branch-name> | <file-path>")
			os.Exit(2)
		}
		if args[0] == "-b" {
			if len(args) != 2 {
				fmt.Println("Usage: kitkat checkout -b <branch-name>")
				os.Exit(2)
			}
			name := args[1]
			if core.IsBranch(name) {
				fmt.Printf("Error: Branch '%s' already exists\n", name)
				os.Exit(1)
			}
			if err := core.CreateBranch(name); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if err := core.CheckoutBranch(name); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
		name := args[0]
		if core.IsBranch(name) {
			if err := core.CheckoutBranch(name); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		} else {
			if err := core.CheckoutFile(name); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}
	},
	"merge": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat merge <branch-name>")
			os.Exit(2)
		}
		if err := core.Merge(args[0]); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		os.Exit(0)
	},
	"reset": func(args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: kitkat reset --hard <commit-hash>")
			os.Exit(2)
		}
		if args[0] != "--hard" {
			fmt.Println("Error: only 'reset --hard' is currently supported")
			fmt.Println("Usage: kitkat reset --hard <commit-hash>")
			os.Exit(2)
		}
		if err := core.ResetHard(args[1]); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		os.Exit(0)
	},
	"rebase": func(args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: kitkat rebase [-i <commit> | --continue | --abort]")
			os.Exit(2)
		}

		switch args[0] {
		case "--abort":
			if err := core.RebaseAbort(); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--continue":
			if err := core.RebaseContinue(); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "-i":
			if len(args) < 2 {
				fmt.Println("Usage: kitkat rebase -i <commit>")
				os.Exit(2)
			}
			if err := core.RebaseInteractive(args[1]); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		default:
			// If no flag, assumes simple rebase which isn't requested but we can default to error
			fmt.Println("Usage: kitkat rebase [-i <commit> | --continue | --abort]")
			os.Exit(2)
		}
	},
	"ls-files": func(args []string) {
		if !core.IsRepoInitialized() {
			fmt.Println("Error: not a kitkat repository (or any of the parent directories): .kitkat")
			os.Exit(1)
		}

		entries, err := core.LoadIndex()
		if err != nil {
			fmt.Println("Error loading index:", err)
			os.Exit(1)
		}

		for _, entry := range entries {
			fmt.Println(entry.Path)
		}
		os.Exit(0)
	},
	"clean": func(args []string) {
		force := false
		includeIgnored := false

		for _, arg := range args {
			switch arg {
			case "-f":
				force = true
			case "-x":
				includeIgnored = true
			}
		}

		if !force {
			fmt.Println("This will delete untracked files. Run 'kitkat clean -f' to proceed.")
			os.Exit(0)
		}

		if err := core.Clean(true, includeIgnored); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
	"help": func(args []string) {
		if len(args) > 0 {
			core.PrintCommandHelp(args[0])
		} else {
			core.PrintGeneralHelp()
		}
		os.Exit(0)
	},
	"tag": func(args []string) {
		if !core.IsRepoInitialized() {
			fmt.Println("Error: not a kitkat repository (or any of the parent directories): .kitkat")
			os.Exit(1)
		}

		if len(args) == 1 && (args[0] == "--list") {
			if err := core.PrintTags(); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		}

		if len(args) < 2 {
			fmt.Println("Usage: kitkat tag <tag-name> <commit-id>")
			os.Exit(2)
		}

		if err := core.CreateTag(args[0], args[1]); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
	"config": func(args []string) {
		if len(args) < 2 || args[0] != "--global" {
			if len(args) == 1 && args[0] == "--list" {
				if err := core.PrintAllConfig(); err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				os.Exit(0)
			}
			fmt.Println("Usage: kitkat config --global <key> [<value>]")
			os.Exit(2)
		}
		key := args[1]
		if len(args) == 3 {
			value := args[2]
			if err := core.SetConfig(key, value); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		} else if len(args) == 2 {
			value, ok, err := core.GetConfig(key)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if ok {
				fmt.Println(value)
				os.Exit(0)
			}
		} else {
			fmt.Println("Usage: kitkat config --global <key> [<value>]")
			os.Exit(2)
		}
	},
	"show-object": func(args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: kitkat show-object <hash>")
			os.Exit(2)
			return
		}
		if err := core.ShowObject(args[0]); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		os.Exit(0)
	},
	"branch": func(args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: kitkat branch [-l | -r <branch-name> | -d <branch-name>]")
			os.Exit(2)
		}
		switch args[0] {
		case "-l":
			if err := core.ListBranches(); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "-r":
			if len(args) < 2 {
				fmt.Fprintln(os.Stderr, "Usage: kitkat branch -r <branch-name>")
				os.Exit(2)
			}

			name := args[1]
			if err := core.RenameCurrentBranch(name); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "-d", "--delete":
			if len(args) < 2 {
				fmt.Fprintln(os.Stderr, "Usage: kitkat branch -d <branch-name>")
				os.Exit(2)
			}

			name := args[1]
			if err := core.DeleteBranch(name); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			} else {
				fmt.Println("Branch `" + name + "` deleted successfully")
				os.Exit(0)
			}
		default:
			name := args[0]
			if core.IsBranch(name) {
				fmt.Printf("Error: Branch '%s' already exists\n", name)
				os.Exit(1)
			}
			if err := core.CreateBranch(name); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	},
	"mv": func(args []string) {
		force := false
		paths := make([]string, 0, 2)

		for _, arg := range args {
			if arg == "-f" || arg == "--force" {
				force = true
				continue
			}
			paths = append(paths, arg)
		}

		if len(paths) != 2 {
			fmt.Println("Usage: kitkat mv [-f|--force] <old_path> <new_path>")
			os.Exit(2)
		}

		if err := core.MoveFile(paths[0], paths[1], force); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		os.Exit(0)
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
		os.Exit(0)
	}
	if len(os.Args) < 2 {
		fmt.Println("Usage: kitkat <command> [args]")
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	if handler, ok := commands[cmd]; ok {
		handler(args)
	} else {
		fmt.Println("Unknown command:", cmd)
		core.PrintGeneralHelp()
		os.Exit(2)
	}
}
