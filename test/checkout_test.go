package test

import (
	"strings"
	"testing"

	"github.com/LeeFred3042U/kitkat/internal/core"
)

// TestCheckoutFile verifies file restoration from the last commit.
// Covers uninitialized repo, successful restoration, missing files, and edge cases.
func TestCheckoutFile(t *testing.T) {
	// Test 1: Uninitialized repo should error
	runUninitTest(t, "Uninitialized", func() error {
		return core.CheckoutFile("any.txt")
	})

	// Test 2: Happy path - add file, commit, then checkout should restore exact content
	runHappyTest(t, "Happy Path", func(repo *TestRepo) {
		repo.AddFile("hello.txt", "world").
			CommitAll("initial commit").
			CheckoutFile("hello.txt")
		repo.AssertFileContent("hello.txt", "world")
	})

	// Test 3: File not in last commit should error with clear message
	runErrorTest(t, "File Not In Commit", "file not found", func(repo *TestRepo) error {
		repo.CommitAll("empty commit")
		return core.CheckoutFile("missing.txt")
	})

	// Test 4: Table-driven edge cases for comprehensive coverage
	t.Run("Table Driven Edges", func(t *testing.T) {
		tests := []struct {
			name       string
			path       string
			files      map[string]string
			wantErrMsg string
		}{
			{"Valid Restore", "a.txt", map[string]string{"a.txt": "A"}, ""},
			{"Missing File", "b.txt", map[string]string{"a.txt": "A"}, "file not found"},
			{
				"Path Traversal",
				"../secret.txt",
				map[string]string{"dummy.txt": "D"},
				"file not found",
			},
			{"Multiple Files", "x.txt", map[string]string{"x.txt": "X", "y.txt": "Y"}, ""},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := NewTestRepo(t, true)
				for path, content := range tt.files {
					repo.AddFile(path, content)
				}
				if len(tt.files) > 0 {
					repo.CommitAll("seed")
				}

				err := core.CheckoutFile(tt.path)

				if tt.wantErrMsg != "" {
					if err == nil {
						t.Errorf(
							"%s: expected error containing %q, got nil",
							tt.name,
							tt.wantErrMsg,
						)
						return
					}
					if !strings.Contains(err.Error(), tt.wantErrMsg) {
						t.Errorf(
							"%s: expected error containing %q, got: %v",
							tt.name,
							tt.wantErrMsg,
							err,
						)
					}
					return
				}

				if err != nil {
					t.Fatalf("%s: unexpected error: %v", tt.name, err)
				}
				if content, exists := tt.files[tt.path]; exists {
					repo.AssertFileContent(tt.path, content)
				}
			})
		}
	})
}

// TestCheckoutBranch verifies branch switching, HEAD updates, and working directory changes.
// Covers uninitialized repo, successful branch switches, dirty WD detection, and nonexistent branches.
func TestCheckoutBranch(t *testing.T) {
	// Test 1: Uninitialized repo should error
	runUninitTest(t, "Uninitialized", func() error {
		return core.CheckoutBranch("feature")
	})

	// Test 2: Happy path - create branches with different files, switch between them
	runHappyTest(t, "Happy Path Branch Switch", func(repo *TestRepo) {
		repo.AddFile("shared.txt", "common").
			CommitAll("init").
			CreateBranch("feature-a").
			AddFile("feat-only.txt", "feature").
			CommitAll("feat commit").
			CreateBranch("feature-b").
			CheckoutBranch("feature-b")

		repo.AssertFileContent("shared.txt", "common")
		repo.AssertFileContent("feat-only.txt", "feature")
		repo.AssertHEAD("ref: refs/heads/feature-b")
	})

	// Test 3: Nonexistent branch should error
	runErrorTest(t, "Nonexistent Branch", "not found", func(repo *TestRepo) error {
		repo.AddFile("test.txt", "test").CommitAll("init")
		return core.CheckoutBranch("phantom")
	})

	// Test 4: Dirty working directory (uncommitted changes) should abort switch
	runErrorTest(
		t,
		"Dirty Working Directory",
		"would be overwritten by checkout",
		func(repo *TestRepo) error {
			repo.AddFile("dirty.txt", "uncommitted").
				CommitAll("initial").
				CreateBranch("target")

			// Modify file without committing (dirty WD)
			if err := repo.modifyFile("dirty.txt", "modified"); err != nil {
				repo.t.Fatalf("Failed to make WD dirty: %v", err)
			}

			return core.CheckoutBranch("target")
		},
	)

	// Test 5: Table-driven multi-branch scenarios
	t.Run("Table Driven Branch Scenarios", func(t *testing.T) {
		tests := []struct {
			name       string
			setup      func(*TestRepo)
			checkout   string
			wantErr    bool
			wantErrMsg string
			asserts    func(*TestRepo)
		}{
			{
				"Switch Between Two Branches",
				func(r *TestRepo) {
					r.AddFile("b1.txt", "B1").CommitAll("b1").CreateBranch("branch1")
					r.AddFile("b2.txt", "B2").CommitAll("b2").CreateBranch("branch2")
				},
				"branch1",
				false,
				"",
				func(r *TestRepo) {
					r.AssertHEAD("ref: refs/heads/branch1")
					r.AssertFileContent("b1.txt", "B1")
				},
			},
			{
				"Switch Back To Initial",
				func(r *TestRepo) {
					r.AddFile("init.txt", "initial").
						CommitAll("init").
						CreateBranch("initial-branch").
						AddFile("feat.txt", "feature").
						CommitAll("feat").
						CreateBranch("feature-branch").
						CheckoutBranch("feature-branch")
				},
				"initial-branch",
				false,
				"",
				func(r *TestRepo) {
					r.AssertHEAD("ref: refs/heads/initial-branch")
					r.AssertFileContent("init.txt", "initial")
				},
			},
			{
				"Invalid Branch Name",
				func(r *TestRepo) {
					r.AddFile("f.txt", "F").CommitAll("init")
				},
				"nonexistent-branch",
				true,
				"not found",
				nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := NewTestRepo(t, true)
				tt.setup(repo)

				err := core.CheckoutBranch(tt.checkout)

				if tt.wantErr {
					if err == nil {
						t.Errorf("Expected error, got nil")
						return
					}
					if !strings.Contains(err.Error(), tt.wantErrMsg) {
						t.Errorf("Expected error containing %q, got: %v", tt.wantErrMsg, err)
					}
					return
				}
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if tt.asserts != nil {
					tt.asserts(repo)
				}
			})
		}
	})
}

// TestCheckoutCommit verifies checking out specific commits (detached HEAD state).
// Covers uninitialized repo, detached HEAD, invalid hashes, and transitions between states.
func TestCheckoutCommit(t *testing.T) {
	// Test 1: Uninitialized repo should error
	runUninitTest(t, "Uninitialized", func() error {
		return core.CheckoutCommit("fakehash")
	})

	// Test 2: Happy path - checkout an earlier commit in detached HEAD state
	runHappyTest(t, "Basic Detached HEAD", func(repo *TestRepo) {
		repo.AddFile("commit1.txt", "v1").
			CommitAll("c1")
		c1ID := repo.commits[0].ID

		repo.AddFile("commit2.txt", "v2").
			CommitAll("c2")

		repo.CheckoutCommit(c1ID)
		repo.AssertHEAD(c1ID)
		repo.AssertFileContent("commit1.txt", "v1")
	})

	// Test 3: Invalid commit hash should error
	runErrorTest(t, "Invalid Commit Hash", "not found", func(repo *TestRepo) error {
		repo.AddFile("test.txt", "test").CommitAll("init")
		return core.CheckoutCommit("invalidhash123")
	})

	// Test 4: Table-driven commit checkout scenarios
	t.Run("Table Driven Commit Scenarios", func(t *testing.T) {
		tests := []struct {
			name       string
			setup      func(*TestRepo) string
			wantErr    bool
			wantErrMsg string
			asserts    func(*TestRepo, string)
		}{
			{
				"Checkout First Commit",
				func(r *TestRepo) string {
					r.AddFile("f1.txt", "1").CommitAll("c1")
					id := r.commits[0].ID
					r.AddFile("f2.txt", "2").CommitAll("c2")
					return id
				},
				false,
				"",
				func(r *TestRepo, id string) {
					r.AssertHEAD(id)
					r.AssertFileContent("f1.txt", "1")
				},
			},
			{
				"Detached to Branch",
				func(r *TestRepo) string {
					r.AddFile("f1.txt", "1").CommitAll("c1")
					id := r.commits[0].ID
					r.AddFile("f2.txt", "2").CommitAll("c2").CreateBranch("stable")
					return id
				},
				false,
				"",
				func(r *TestRepo, id string) {
					r.AssertHEAD(id)
					r.CheckoutBranch("stable")
					r.AssertHEAD("ref: refs/heads/stable")
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := NewTestRepo(t, true)
				commitID := tt.setup(repo)

				err := core.CheckoutCommit(commitID)

				if tt.wantErr {
					if err == nil {
						t.Errorf("Expected error, got nil")
						return
					}
					if !strings.Contains(err.Error(), tt.wantErrMsg) {
						t.Errorf("Expected error containing %q, got: %v", tt.wantErrMsg, err)
					}
					return
				}
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if tt.asserts != nil {
					tt.asserts(repo, commitID)
				}
			})
		}
	})
}
