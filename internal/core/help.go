package core

import "fmt"

type CommandHelp struct {
	Summary string
	Usage   string
}

// Maps command names to their help information
var helpMessages = map[string]CommandHelp{
	"init": {
		Summary: "Initialize a new KitKat repository",
		Usage:   "Usage: kitkat init\n\nInitializes a new .kitkat directory in the current folder, preparing it for tracking files",
	},
	"add": {
		Summary: "Add file contents to the index.",
		Usage:   "Usage: kitkat add <file-path> | --all\n\nThis command adds file contents to the staging area.\nUse '--all' to stage all new, modified, and deleted files.",
	},
	"commit": {
		Summary: "Record changes to the repository.",
		Usage:   "Usage: kitkat commit <-m | -am> <message>\n\nCreates a new commit from the staging area.\nUse '-am' to automatically stage all tracked files before committing.",
	},
	"diff": {
	    Summary: "Show changes between the last commit and staging area",
	    Usage:   "Usage: kitkat diff\n\nShows content differences between the HEAD commit and the index",
	},
	"log": {
	    Summary: "Show the commit history",
	    Usage:   "Usage: kitkat log\n\nDisplays a list of all commits in the current branch's history in reverse chronological order.",
	},
	
	"tag": {
		Summary: "Create a new tag for a commit",
		Usage:   "Usage: kitkat tag <tag-name> <commit-id>\n\nCreates a new lightweight tag that points to the specified commit",
	},
	"merge": {
		Summary: "Merge a branch into the current branch.",
		Usage:   "Usage: kitkat merge <branch-name>\n\nJoins another branch's history into the current branch. Currently, only fast-forward merges are supported.",
	},
	"ls-files": {
		Summary: "Show information about files in the index",
		Usage:   "Usage: kitkat ls-files\n\nPrints a list of all files that are currently in the index (staging area)",
	},
	"clean": {
		Summary: "Remove untracked files from the working directory",
		Usage:   "Usage: kitkat clean\n\nRemoves all files from the current directory that are not tracked by KitKat",
	},
	"config": {
		Summary: "Get and set repository or global options.",
		Usage:   "Usage: kitkat config --global <key> <value>\n\nSets a global configuration value that will be used for all repositories.",
	},
}

// Displays the main help message with a list of all commands
func PrintGeneralHelp() {
	fmt.Println("usage: kitkat <command> [arguments]")
	fmt.Println("\nThese are the common KitKat commands:")
	for name, help := range helpMessages {
		fmt.Printf("   %-10s %s\n", name, help.Summary)
	}
	fmt.Println("\nUse 'kitkat help <command>' for more information about a command")
}

// Displays the detailed usage for a specific command
func PrintCommandHelp(command string) {
	if help, ok := helpMessages[command]; ok {
		fmt.Println(help.Usage)
	} else {
		fmt.Printf("Unknown help topic: '%s'. See 'kitkat --help'.\n", command)
	}
}