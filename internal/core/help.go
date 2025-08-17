package core

import "fmt"

type CommandHelp struct {
	Summary string
	Usage   string
}

// Maps command names to their help information
var helpMessages = map[string]CommandHelp{
	"init": {
		Summary: "Initialize a new KitKat repository.",
		Usage:   "Usage: kitkat init\n\nInitializes a new .kitkat directory in the current folder, preparing it for tracking files.",
	},
	"add": {
		Summary: "Add file contents to the index.",
		Usage:   "Usage: kitkat add <file>\n\nThis command hashes the file's content and adds it to the index (the staging area), preparing it to be included in the next commit.",
	},
	"commit": {
		Summary: "Record changes to the repository.",
		Usage:   "Usage: kitkat commit -m <message>\n\nCreates a new commit with the current contents of the index. The commit records a snapshot of the tracked files at this point in time.",
	},
	"log": {
		Summary: "Log a message for personal tracking.",
		Usage:   "Usage: kitkat log <message>\n\nAdds a new entry to the logs.txt file with a unique ID and timestamp.",
	},
	"view": {
		Summary: "Show all log entries.",
		Usage:   "Usage: kitkat view\n\nDisplays all entries from the logs.txt file in a formatted way.",
	},
	"tag": {
		Summary: "Create a new tag for a commit.",
		Usage:   "Usage: kitkat tag <tag-name> <commit-id>\n\nCreates a new lightweight tag that points to the specified commit.",
	},
	"grep": {
		Summary: "Search for a term in log messages and tags.",
		Usage:   "Usage: kitkat grep <term>\n\nSearches through all logs and prints entries where the message or tag contains the search term.",
	},
	"ls-files": {
		Summary: "Show information about files in the index.",
		Usage:   "Usage: kitkat ls-files\n\nPrints a list of all files that are currently in the index (staging area).",
	},
	"clean": {
		Summary: "Remove untracked files from the working directory.",
		Usage:   "Usage: kitkat clean\n\nRemoves all files from the current directory that are not tracked by KitKat.",
	},
}

// Displays the main help message with a list of all commands
func PrintGeneralHelp() {
	fmt.Println("usage: kitkat <command> [arguments]")
	fmt.Println("\nThese are the common KitKat commands:")
	for name, help := range helpMessages {
		fmt.Printf("   %-10s %s\n", name, help.Summary)
	}
	fmt.Println("\nUse 'kitkat help <command>' for more information about a command.")
}

// Displays the detailed usage for a specific command
func PrintCommandHelp(command string) {
	if help, ok := helpMessages[command]; ok {
		fmt.Println(help.Usage)
	} else {
		fmt.Printf("Unknown help topic: '%s'. See 'kitkat --help'.\n", command)
	}
}