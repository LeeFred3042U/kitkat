package core

import (
	"crypto/sha1"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/models"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// RebaseInteractive starts an interactive rebase session
func RebaseInteractive(commitHash string) error {
	if !IsRepoInitialized() {
		return fmt.Errorf("not a kitkat repository")
	}

	// 1. Validate clean working directory
	isDirty, err := IsWorkDirDirty()
	if err != nil {
		return fmt.Errorf("failed to check working directory status: %w", err)
	}
	if isDirty {
		return fmt.Errorf("cannot rebase: you have unstaged changes")
	}

	// 2. Resolve 'onto' commit
	ontoCommit, err := storage.FindCommit(commitHash)
	if err != nil {
		return fmt.Errorf("invalid base commit '%s': %w", commitHash, err)
	}

	// 3. Get current HEAD info for potential abort
	headState, err := GetHeadState()
	if err != nil {
		return err
	}
	headHash, err := readCurrentHeadCommit()
	if err != nil {
		return err
	}
	// If headState is detached (HEAD <hash>), we store empty branch name
	rebaseHeadNameVal := ""
	if !strings.HasPrefix(headState, "HEAD") {
		rebaseHeadNameVal = "refs/heads/" + headState
	}

	// 4. Find commits to rebase (onto..HEAD)
	// We need these in chronological order (oldest to newest)
	commitsToRebase, err := getCommitsBetween(ontoCommit.ID, headHash)
	if err != nil {
		return err
	}
	if len(commitsToRebase) == 0 {
		fmt.Println("No commits to rebase.")
		return nil
	}

	// 5. Generate TODO list
	todoPath := filepath.Join(RepoDir, "rebase-todo")
	todoContent := generateTodo(commitsToRebase)
	if err := os.WriteFile(todoPath, []byte(todoContent), 0644); err != nil {
		return err
	}

	// 6. Open Editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Fallback for Windows
		editor = "notepad"
	}

	cmd := exec.Command(editor, todoPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Opening editor to modify rebase todo list...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run editor: %w", err)
	}

	// 7. Parse User Input
	newTodoContent, err := os.ReadFile(todoPath)
	if err != nil {
		return err
	}
	steps := parseTodo(string(newTodoContent))
	if len(steps) == 0 {
		fmt.Println("Nothing to do.")
		return nil
	}

	// 8. Initialize State
	state := RebaseState{
		HeadName:    rebaseHeadNameVal,
		Onto:        ontoCommit.ID,
		OrigHead:    headHash,
		TodoSteps:   steps,
		CurrentStep: 0,
	}
	if err := SaveRebaseState(state); err != nil {
		return err
	}

	// 9. Attach HEAD to temporary branch 'kitkat-rebase-tmp' pointing to 'onto'
	tmpBranch := "kitkat-rebase-tmp"
	tmpBranchPath := filepath.Join(".kitkat", "refs", "heads", tmpBranch)
	if err := os.MkdirAll(filepath.Dir(tmpBranchPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(tmpBranchPath, []byte(ontoCommit.ID), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(".kitkat/HEAD", []byte("ref: refs/heads/"+tmpBranch), 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}
	// Update workspace/index to match
	if err := UpdateWorkspaceAndIndex(ontoCommit.ID); err != nil {
		return fmt.Errorf("failed to checkout base: %w", err)
	}

	// 10. Start Loop
	return RunRebaseLoop()
}

// RebaseContinue resumes a halted rebase
func RebaseContinue() error {
	if !IsRebaseInProgress() {
		return fmt.Errorf("no rebase in progress")
	}

	state, err := LoadRebaseState()
	if err != nil {
		return err
	}

	if state.CurrentStep >= len(state.TodoSteps) {
		return fmt.Errorf("no steps remaining")
	}

	currentCmdLine := state.TodoSteps[state.CurrentStep]
	parts := strings.Fields(currentCmdLine)
	cmd := parts[0]

	// Original commit we were trying to process
	if len(parts) < 2 {
		return AdvanceRebaseStep(state)
	}
	originalHash := parts[1]
	originalCommit, _ := storage.FindCommit(originalHash)

	// Use switch instead of if/else chain
	switch cmd {
	case "pick", "reword":
		// Commit current index using original message
		msg := originalCommit.Message

		_, _, err := Commit(msg)
		if err != nil {
			if strings.Contains(err.Error(), "nothing to commit") {
				fmt.Println("Nothing to commit. Skipping step.")
			} else {
				return err
			}
		}

		if cmd == "reword" {
			head, _ := readCurrentHeadCommit()
			newMsg := promptForMessage(msg)
			if newMsg != msg {
				amendCommitMessage(head, newMsg)
			}
		}

	case "squash":
		prevHead, _ := GetHeadCommit()
		newMsg := prevHead.Message + "\n\n" + originalCommit.Message

		err := amendCommit(prevHead, newMsg)
		if err != nil {
			return err
		}
	}

	if err := AdvanceRebaseStep(state); err != nil {
		return err
	}
	return RunRebaseLoop()
}

// RebaseAbort restores state
func RebaseAbort() error {
	if !IsRebaseInProgress() {
		return fmt.Errorf("no rebase in progress")
	}
	state, err := LoadRebaseState()
	if err != nil {
		return err
	}

	fmt.Printf("Aborting rebase. restoring HEAD to %s\n", state.OrigHead[:7])

	// Also need to restore branch Ref if we were on a branch
	if state.HeadName != "" {
		refPath := filepath.Join(".kitkat", state.HeadName)
		if err := os.MkdirAll(filepath.Dir(refPath), 0755); err != nil {
			return err
		}
		// Restore original request branch pointer (ResetHard only updated the tmp branch content if we were in tmp)
		// Actually ResetHard(OrigHead) would have moved HEAD (tmp branch) to OrigHead.
		// We need to switch BACK to Main and ResetHard(OrigHead).
	}

	// Switch back to original branch
	if state.HeadName != "" {
		if err := os.WriteFile(".kitkat/HEAD", []byte("ref: "+state.HeadName), 0644); err != nil {
			return err
		}
		// Now ensure Main is at OrigHead (it should be, we never touched it directly, only tmp branch)
		// But ResetHard might have been needed if we were dirty?
		// Safety:
		if err := os.WriteFile(filepath.Join(".kitkat", state.HeadName), []byte(state.OrigHead), 0644); err != nil {
			return err
		}
		if err := UpdateWorkspaceAndIndex(state.OrigHead); err != nil {
			return err
		}
	} else {
		// Detached ORIG_HEAD
		if err := ResetHard(state.OrigHead); err != nil {
			return err
		}
	}

	// Clean up tmp branch
	os.Remove(filepath.Join(".kitkat", "refs", "heads", "kitkat-rebase-tmp"))

	return ClearRebaseState()
}

// RunRebaseLoop executes steps until done or blocked
func RunRebaseLoop() error {
	for {
		cmdLine, state, err := ReadNextTodo()
		if err != nil {
			return err
		}
		if state.CurrentStep >= len(state.TodoSteps) {
			// Done!
			fmt.Println("Rebase completed successfully.")
			return finishRebase(state)
		}

		// Execute cmd
		parts := strings.Fields(cmdLine)
		if len(parts) < 2 {
			AdvanceRebaseStep(state)
			continue
		}
		action := parts[0]
		commitHash := parts[1]

		fmt.Printf("Rebase (%d/%d): %s\n", state.CurrentStep+1, len(state.TodoSteps), cmdLine)

		var stepErr error
		switch action {
		case "pick", "p":
			stepErr = executePick(commitHash)
		case "reword", "r":
			stepErr = executeReword(commitHash)
		case "squash", "s":
			stepErr = executeSquash(commitHash)
		case "drop", "d":
			fmt.Printf("Dropping commit %s\n", commitHash)
			stepErr = nil // Just don't apply it
		default:
			fmt.Printf("Unknown command '%s'. Skipping.\n", action)
		}

		if stepErr != nil {
			fmt.Printf("Conflict or error at step %d: %v\n", state.CurrentStep+1, stepErr)
			fmt.Println("Resolve conflicts, then run 'kitkat rebase --continue'.")
			fmt.Println("To stop, run 'kitkat rebase --abort'.")
			return nil // Exit process, leave state on disk
		}

		if err := AdvanceRebaseStep(state); err != nil {
			return err
		}
	}
}

func finishRebase(state *RebaseState) error {
	headHash, err := readCurrentHeadCommit()
	if err != nil {
		return err
	}

	// Switch to original branch and Fast-Forward to headHash
	if state.HeadName != "" {
		if err := os.WriteFile(".kitkat/HEAD", []byte("ref: "+state.HeadName), 0644); err != nil {
			return err
		}

		// Update branch pointer
		refPath := filepath.Join(".kitkat", state.HeadName)
		if err := os.WriteFile(refPath, []byte(headHash), 0644); err != nil {
			return err
		}

		// Workspace should already match headHash (since we were on tmp branch at headHash)
		// But UpdateBranchPointer doesn't touch workspace.
	}

	// Delete tmp branch
	os.Remove(filepath.Join(".kitkat", "refs", "heads", "kitkat-rebase-tmp"))

	return ClearRebaseState()
}

// executePick applies the changes from the commit
func executePick(hash string) error {
	return cherryPick(hash, false)
}

func executeReword(hash string) error {
	if err := cherryPick(hash, false); err != nil {
		return err
	}

	head, _ := GetHeadCommit()
	newMsg := promptForMessage(head.Message)

	return amendCommitMessage(head.ID, newMsg)
}

func executeSquash(hash string) error {
	if err := cherryPick(hash, true); err != nil { // true = no commit, just stage
		return err
	}

	prevHead, _ := GetHeadCommit()
	targetCommit, _ := storage.FindCommit(hash)

	newMsg := prevHead.Message + "\n\n" + targetCommit.Message

	return amendCommit(prevHead, newMsg)
}

// cherryPick applies changes from a commit to the current HEAD
func cherryPick(hash string, noCommit bool) error {
	commit, err := storage.FindCommit(hash)
	if err != nil {
		return err
	}

	parentHash := commit.Parent

	changes, err := getChanges(parentHash, hash)
	if err != nil {
		return err
	}

	if err := applyChanges(changes); err != nil {
		return err
	}

	if noCommit {
		return nil
	}

	_, _, err = Commit(commit.Message)

	if err != nil && strings.Contains(err.Error(), "nothing to commit") {
		return nil
	}

	return err
}

type Change struct {
	OldHash string
	NewHash string
}

func getChanges(parentHash, childHash string) (map[string]Change, error) {
	parentTree := make(map[string]string)
	if parentHash != "" {
		pC, err := storage.FindCommit(parentHash)
		if err == nil {
			parentTree, _ = storage.ParseTree(pC.TreeHash)
		}
	}

	childCommit, err := storage.FindCommit(childHash)
	if err != nil {
		return nil, err
	}
	childTree, err := storage.ParseTree(childCommit.TreeHash)
	if err != nil {
		return nil, err
	}

	changes := make(map[string]Change)

	// Added or Modified
	for path, hash := range childTree {
		if pHash, ok := parentTree[path]; !ok || pHash != hash {
			// pHash is "" if not in parent (Added)
			changes[path] = Change{OldHash: parentTree[path], NewHash: hash}
		}
	}

	// Deleted - using empty string as marker
	for path := range parentTree {
		if _, ok := childTree[path]; !ok {
			changes[path] = Change{OldHash: parentTree[path], NewHash: ""}
		}
	}

	return changes, nil
}

func applyChanges(changes map[string]Change) error {
	headCommit, _ := GetHeadCommit()
	headTree, _ := storage.ParseTree(headCommit.TreeHash)

	for path, change := range changes {
		targetHash := change.NewHash
		if targetHash == "" {
			// Delete
			// Conflict Check: If HEAD doesn't match OldHash, we have moved.
			headFileHash, existsInHead := headTree[path]
			if existsInHead && headFileHash != change.OldHash {
				return fmt.Errorf("conflict in %s: deleted in incoming commit, but modified in HEAD", path)
			}

			if err := RemoveFile(path); err != nil {
				return err
			}
		} else {
			// Write
			content, err := storage.ReadObject(targetHash)
			if err != nil {
				return err
			}

			// Detect Conflict
			headFileHash, existsInHead := headTree[path]
			// If the file exists in HEAD, it MUST match the version we are changing FROM (OldHash).
			// If change.OldHash is empty (new file), HEAD must not have it.

			if existsInHead {
				if headFileHash != change.OldHash {
					return fmt.Errorf("conflict in %s: modified in incoming commit, but modified in HEAD", path)
				}
			} else {
				// Determine if we expected it to exist
				if change.OldHash != "" {
					// We expected it to exist (modification), but it's gone in HEAD. Conflict.
					return fmt.Errorf("conflict in %s: modified in incoming commit, but deleted in HEAD", path)
				}
			}

			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(path, content, 0644); err != nil {
				return err
			}

			if err := AddFile(path); err != nil {
				return err
			}
		}
	}
	return nil
}

// Helpers

func generateTodo(hashes []string) string {
	var sb strings.Builder
	for _, h := range hashes {
		c, _ := storage.FindCommit(h)
		sb.WriteString(fmt.Sprintf("pick %s %s\n", h, c.Message))
	}

	sb.WriteString("\n# Commands:\n")
	sb.WriteString("# p, pick <commit> = use commit\n")
	sb.WriteString("# r, reword <commit> = use commit, but edit the commit message\n")
	sb.WriteString("# s, squash <commit> = use commit, but meld into previous commit\n")
	sb.WriteString("# d, drop <commit> = remove commit\n")

	return sb.String()
}

func parseTodo(content string) []string {
	var steps []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		steps = append(steps, line)
	}
	return steps
}

func getCommitsBetween(start, end string) ([]string, error) {
	var chain []string
	curr := end
	for curr != "" {
		if curr == start {
			break
		}
		chain = append(chain, curr)

		c, err := storage.FindCommit(curr)
		if err != nil {
			return nil, err
		}
		if c.Parent == "" {
			// Reached root without finding start
			// Check if start is actually empty (rebase from root?)
			if start == "" {
				return chain, nil
			} // Should reverse first
			break
		}
		curr = c.Parent
	}

	// Reverse
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	return chain, nil
}

func promptForMessage(defaultMsg string) string {
	tmp := ".kitkat/COMMIT_EDITMSG"
	os.WriteFile(tmp, []byte(defaultMsg), 0644)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "notepad"
	}

	cmd := exec.Command(editor, tmp)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	out, _ := os.ReadFile(tmp)
	return strings.TrimSpace(string(out))
}

func amendCommitMessage(commitID, newVal string) error {
	c, err := storage.FindCommit(commitID)
	if err != nil {
		return err
	}

	parentBlock := ""
	if c.Parent != "" {
		parentBlock = "parent " + c.Parent + "\n"
	}

	content := fmt.Sprintf("tree %s\n%s\n%s", c.TreeHash, parentBlock, newVal)
	newHash, err := saveObject([]byte(content))
	if err != nil {
		return err
	}

	return UpdateBranchPointer(newHash)
}

func amendCommit(prevHead models.Commit, newMsg string) error {
	treeHash, err := storage.CreateTree()
	if err != nil {
		return err
	}

	parentBlock := ""
	if prevHead.Parent != "" {
		parentBlock = "parent " + prevHead.Parent + "\n"
	}

	content := fmt.Sprintf("tree %s\n%s\n%s", treeHash, parentBlock, newMsg)
	newHash, err := saveObject([]byte(content))
	if err != nil {
		return err
	}

	return UpdateBranchPointer(newHash)
}

// Removed unused 'objType' parameter
func saveObject(content []byte) (string, error) {
	// Manually hashing and saving, similar to storage/blob.go but for memory content
	h := sha1.New()
	h.Write(content)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	objPath := filepath.Join(".kitkat", "objects", hash)
	if err := os.MkdirAll(filepath.Dir(objPath), 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(objPath, content, 0644); err != nil {
		return "", err
	}
	return hash, nil
}
