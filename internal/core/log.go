package core

import (
	"fmt"
	"sort"

	"github.com/LeeFred3042U/kitcat/internal/models"
	"github.com/LeeFred3042U/kitcat/internal/storage"
)

// ShowLog prints the commit log. It accepts a boolean for oneline format
// and an optional limit to restrict the number of commits shown (use -1 or 0 for no limit)
func ShowLog(oneline bool, limit int) error {
	// 1Start from HEAD (Architecture from reset-hard branch)
	// We must walk backwards from HEAD, otherwise 'reset' changes won't be reflected
	currentCommit, err := GetHeadCommit()
	if err != nil {
		// Handle the case where the repo is empty or HEAD is invalid
		return nil
	}

	commitHash := currentCommit.ID
	count := 0

	// Walk the graph (Architecture from reset-hard branch)
	for commitHash != "" {
		// Apply the Limit Check (Feature from main branch)
		if limit > 0 && count >= limit {
			break
		}

		commit, err := storage.FindCommit(commitHash)
		if err != nil {
			return err
		}

		// Print Logic
		if oneline {
			fmt.Printf("%s %s\n", commit.ID[:7], commit.Message)
		} else {
			fmt.Printf("commit %s\n", commit.ID)
			fmt.Printf("Author: %s <%s>\n", commit.AuthorName, commit.AuthorEmail)
			fmt.Printf("Date:   %s\n", commit.Timestamp.Local().Format("Mon Jan 02 15:04:05 2006 -0700"))
			fmt.Printf("\n    %s\n\n", commit.Message)
		}

		// Move to parent pointer
		commitHash = commit.Parent
		count++
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
			fmt.Printf("\t%s\n", commit.Message)
		}
		fmt.Println()
	}

	return nil
}
