package storage

import (
	"os"
	"fmt"
	"bufio"
	"errors"
	"encoding/json"

	"github.com/LeeFred3042U/kitkat/internal/models"
)

var ErrNoCommits = errors.New("no commits yet")

const commitsPath = ".kitkat/commits.log"

// Appends commit as NDJSON
func AppendCommit(commit models.Commit) error {
	if err := os.MkdirAll(".kitkat", 0755); err != nil {
		return err
	}

	// Use the generic lock function from lock*.go for consistency
	lockFile, err := lock(commitsPath)
	if err != nil {
		return err
	}
	defer unlock(lockFile)

	f, err := os.OpenFile(commitsPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(commit); err != nil {
		return err
	}
	return f.Sync()
}

// Reads commits (NDJSON)
func ReadCommits() ([]models.Commit, error) {
	var commits []models.Commit
	if _, err := os.Stat(commitsPath); os.IsNotExist(err) {
		return commits, nil
	}

	f, err := os.Open(commitsPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var c models.Commit
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			continue
		}
		commits = append(commits, c)
	}
	return commits, scanner.Err()
}

// Returns ErrNoCommits when none exist
func GetLastCommit() (models.Commit, error) {
	commits, err := ReadCommits()
	if err != nil {
		return models.Commit{}, err
	}
	if len(commits) == 0 {
		return models.Commit{}, ErrNoCommits
	}
	return commits[len(commits)-1], nil
}

// Search the commit log for a commit with a matching hash
func FindCommit(hash string) (models.Commit, error) {
	file, err := os.Open(commitsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return models.Commit{}, ErrNoCommits
		}
		return models.Commit{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var commit models.Commit
		if err := json.Unmarshal(scanner.Bytes(), &commit); err != nil {
			continue
		}
		if commit.ID == hash {
			return commit, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return models.Commit{}, err
	}

	return models.Commit{}, fmt.Errorf("commit with hash %s not found", hash)
}

// IsAncestor returns true if ancestorHash is equal to or is an ancestor of descendantHash
func IsAncestor(ancestorHash, descendantHash string) (bool, error) {
	if ancestorHash == "" || descendantHash == "" {
		return false, nil
	}
	// A commit is its own ancestor
	if ancestorHash == descendantHash {
		return true, nil
	}

	current := descendantHash
	for current != "" {
		c, err := FindCommit(current)
		if err != nil {
			return false, err
		}
		if c.ID == ancestorHash {
			return true, nil
		}
		// walk up
		current = c.Parent
	}
	return false, nil
}