package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// RebaseState tracks ongoing rebase
type RebaseState struct {
	HeadName    string   // "refs/heads/main" or empty if detached
	Onto        string   // Commit ID base
	OrigHead    string   // Commit ID where we started (for abort)
	TodoSteps   []string // List of commands
	CurrentStep int      // Index in TodoSteps (0-based)
	Message     string   // For squash/reword message accumulation
}

func EnsureRebaseDir() error {
	path := filepath.Join(RepoDir, "rebase-merge")
	return os.MkdirAll(path, 0755)
}

func SaveRebaseState(state RebaseState) error {
	if err := EnsureRebaseDir(); err != nil {
		return err
	}
	base := filepath.Join(RepoDir, "rebase-merge")

	os.WriteFile(filepath.Join(base, "head-name"), []byte(state.HeadName), 0644)
	os.WriteFile(filepath.Join(base, "onto"), []byte(state.Onto), 0644)
	os.WriteFile(filepath.Join(base, "orig-head"), []byte(state.OrigHead), 0644)
	os.WriteFile(filepath.Join(base, "git-rebase-todo"), []byte(strings.Join(state.TodoSteps, "\n")), 0644)
	os.WriteFile(filepath.Join(base, "msgnum"), []byte(fmt.Sprintf("%d", state.CurrentStep+1)), 0644)
	os.WriteFile(filepath.Join(base, "message"), []byte(state.Message), 0644) // Optional

	return nil
}

func LoadRebaseState() (*RebaseState, error) {
	base := filepath.Join(RepoDir, "rebase-merge")
	if _, err := os.Stat(base); os.IsNotExist(err) {
		return nil, fmt.Errorf("no rebase in progress")
	}

	headName, _ := os.ReadFile(filepath.Join(base, "head-name"))
	onto, _ := os.ReadFile(filepath.Join(base, "onto"))
	origHead, _ := os.ReadFile(filepath.Join(base, "orig-head"))
	todoData, _ := os.ReadFile(filepath.Join(base, "git-rebase-todo"))
	msgNumData, _ := os.ReadFile(filepath.Join(base, "msgnum"))
	message, _ := os.ReadFile(filepath.Join(base, "message"))

	step, _ := strconv.Atoi(strings.TrimSpace(string(msgNumData)))
	// step is 1-based in file, 0-based in struct
	if step > 0 {
		step--
	}

	return &RebaseState{
		HeadName:    strings.TrimSpace(string(headName)),
		Onto:        strings.TrimSpace(string(onto)),
		OrigHead:    strings.TrimSpace(string(origHead)),
		TodoSteps:   strings.Split(strings.TrimSpace(string(todoData)), "\n"),
		CurrentStep: step,
		Message:     string(message),
	}, nil
}

func IsRebaseInProgress() bool {
	_, err := os.Stat(filepath.Join(RepoDir, "rebase-merge"))
	return err == nil
}

func ClearRebaseState() error {
	return os.RemoveAll(filepath.Join(RepoDir, "rebase-merge"))
}

func ReadNextTodo() (string, *RebaseState, error) {
	state, err := LoadRebaseState()
	if err != nil {
		return "", nil, err
	}
	if state.CurrentStep >= len(state.TodoSteps) {
		return "", state, nil
	}
	return state.TodoSteps[state.CurrentStep], state, nil
}

func AdvanceRebaseStep(state *RebaseState) error {
	state.CurrentStep++
	return SaveRebaseState(*state)
}
