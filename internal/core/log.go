package core

import (
	"fmt"
	"sort"

	"github.com/LeeFred3042U/kitkat/internal/models"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// ShowLog prints the commit log. It accepts a boolean to switch to a compact, one-line format.
func ShowLog(oneline bool) error {
	commits, err := storage.ReadCommits()
	if err != nil {
		return err
	}

	// Print in reverse chronological order (newest first).
	for i := len(commits) - 1; i >= 0; i-- {
		commit := commits[i]
		if oneline {
			fmt.Printf("%s %s\n", commit.ID[:7], commit.Message)
		} else {
			fmt.Printf("commit %s\n", commit.ID)
			fmt.Printf("Author: %s <%s>\n", commit.AuthorName, commit.AuthorEmail)
			fmt.Printf("Date:   %s\n", commit.Timestamp.Local().Format("Mon Jan 02 15:04:05 2006 -0700"))
			fmt.Printf("\n    %s\n\n", commit.Message)
		}
	}
	return nil
}

// ShowShortLog prints commit messages grouped by author,
// sorted by commit counts of each author.
func ShowShortLog() error {
	commits, err := storage.ReadCommits()
	if err != nil {
		return err
	}

	// Groups commits by author.
	authorCommits := make(map[string][]models.Commit)
	for _, commit := range commits {
		authorCommits[commit.AuthorName] = append(authorCommits[commit.AuthorName], commit)
	}

	// Builds a sortable slice.
	type authorLog struct {
		name    string
		commits []models.Commit
	}
	var logs []authorLog
	for author, commits := range authorCommits {
		logs = append(logs, authorLog{
			name:    author,
			commits: commits,
		})
	}

	// Sorts the slice by number of commits in descending order.
	sort.Slice(logs, func(i, j int) bool {
		return len(logs[i].commits) > len(logs[j].commits)
	})

	// Prints the shortlog.
	for _, log := range logs {
		fmt.Printf("%s (%d):\n", log.name, len(log.commits))
		for _, commit := range log.commits {
			fmt.Printf("\t%s\n", commit)
		}
		fmt.Println()
	}

	return nil
}
