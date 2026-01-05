# kitkat Testing Guidelines & Standards

To ensure the stability of the kitkat VCS engine, every new command must include a corresponding `_test.go` file in the `test/` directory
These tests serve as the primary proof of stability and must maintain high-quality systems engineering standards

## 1. Environment Isolation & Isolation Policy

Tests **must never** modify your actual filesystem or global environment.

- **Sandbox Creation:** Always use `t.TempDir()` to create a sandbox for the test
- **Process Movement:** Use `os.Chdir` to move the process into that temporary directory so commands like `init` work as expected
- **State Restoration:** Always use a `defer` to return to the original working directory and restore environment variables to prevent side effects in other tests
- **Environment Mocking:** If a command touches global state (like `$HOME` for config or `$EDITOR` for rebase), you must mock those variables in the test and restore them using `defer`

```go
tempDir := t.TempDir()
origDir, _ := os.Getwd()
origHome := os.Getenv("HOME")

os.Setenv("HOME", tempDir)
os.Chdir(tempDir)

defer func() {
    os.Chdir(origDir)
    os.Setenv("HOME", origHome)
}()
```

## 2. The "Three-State" Test Coverage

Every command test suite must cover these three core states:

### A. The Uninitialized (Zero) State

Verify that the command fails gracefully if run before `kitkat init`

- **Requirement:** It should return a clear "not a kitkat repository" error
- **Check:** Explicitly check for `err != nil` when `IsRepoInitialized()` is false

### B. The Happy Path (Success)

Verify the command works perfectly under ideal conditions

- **Requirement:**
  - Do not just check the error return value
  - You must verify the physical side effects on the filesystem

### C. Edge Cases and Errors

Verify how the command handles "illegal," "awkward," or malicious states:

- **Naming & Security:** What happens with paths or tags named `../root`? (Ensure path traversal safety)
- **Duplicates:** Ensure the system prevents overwriting existing data (e.g., tags, branches) unless a `--force` flag is explicitly used
- **Empty Inputs:** Test behavior with empty messages, empty paths, or empty hashes.
- **Missing Files:** Test how commands like `add` or `rm` handle files that were deleted from the disk before the command was executed
- **Empty State:** Ensure `List` commands return empty slices (not errors) when no items exist

## 3. Atomic & Physical Verification

When a command changes the repository state (like `commit` or `tag`), your test must verify the result manually:

1. **Action:** Perform the command
2. **Physical Inspection:** Manually use `os.Stat` or `os.ReadFile` to inspect the resulting file in the `.kitkat` directory
3. **Content Accuracy:** Ensure the stored data matches the input (e.g., the tag file contains the correct commit ID hash)
4. **Data Integrity:** Ensure the data format matches the specification (e.g., valid JSON in the `index` or `commits.log`)

## 4. Integration vs. Unit Testing

While unit tests (testing internal functions) are preferred for speed, commands that involve external processes require specific handling:

- **Editor Mocking:** Commands that launch external editors (like `rebase -i`) must use mocked environment variables (e.g., `EDITOR="touch"`) to allow tests to complete automatically without human intervention
