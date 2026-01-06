package plumbing

import (
	"bytes"
	"fmt"
	"time"
)

// CommitOptions holds data for creating a commit
type CommitOptions struct {
	Tree      string   // SHA-1 of root tree
	Parents   []string // SHA-1 of parent commits
	Author    string   // "Name <email>"
	Committer string   // "Name <email>"
	Message   string   // commit message
}

// CommitTree writes a commit object and returns its SHA-1
func CommitTree(opts CommitOptions) (string, error) {
	var buf bytes.Buffer

	// ---- Commit content ----
	buf.WriteString(fmt.Sprintf("tree %s\n", opts.Tree))
	for _, p := range opts.Parents {
		buf.WriteString(fmt.Sprintf("parent %s\n", p))
	}

	now := time.Now()
	seconds := now.Unix()
	offset := now.Format("-0700") // timezone offset like "+0530"

	author := opts.Author
	if author == "" {
		author = "KitKat User <user@kitkat>"
	}
	buf.WriteString(fmt.Sprintf("author %s %d %s\n", author, seconds, offset))

	committer := opts.Committer
	if committer == "" {
		committer = author
	}
	buf.WriteString(fmt.Sprintf("committer %s %d %s\n\n", committer, seconds, offset))

	buf.WriteString(opts.Message)
	if !bytes.HasSuffix(buf.Bytes(), []byte("\n")) {
		buf.WriteByte('\n')
	}

	return HashAndWriteObject(buf.Bytes(), "commit")
}
