package storage

import (
	"os"
	"fmt"
	"time"
	"bufio"
	"errors"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/models"
)

const commitsPath = ".kitkat/commits.log"

// appends a new commit
func AppendCommit(commit models.Commit) error {
	f, err := os.OpenFile(commitsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := commit.Timestamp.Format(time.RFC3339)
	commitLine := fmt.Sprintf("%s %s \"%s\" %s\n", commit.ID, commit.TreeHash, commit.Message, timestamp)

	if _, err := f.WriteString(commitLine); err != nil {
		return err
	}

	return nil
}

// Reads all commits
func ReadCommits() ([]models.Commit, error) {
	var commits []models.Commit
	f, err := os.Open(commitsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return commits, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 4)
		if len(parts) < 4 {
			continue
		}

		timestamp, _ := time.Parse(time.RFC3339, parts[3])
		commit := models.Commit{
			ID:        parts[0],
			TreeHash:  parts[1],
			Message:   strings.Trim(parts[2], "\""),
			Timestamp: timestamp,
		}
		commits = append(commits, commit)
	}
	return commits, scanner.Err()
}

// Reads the last commit from commits.log file
func GetLastCommit() (models.Commit, error) {
	commits, err := ReadCommits()
	if err != nil {
		return models.Commit{}, err
	}
	if len(commits) == 0 {
		return models.Commit{}, errors.New("no commits yet")
	}
	return commits[len(commits)-1], nil
}