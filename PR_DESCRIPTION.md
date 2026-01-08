# Description

This PR adds a reproduction test `Test_CheckoutFile_OverwritesDirtyFile` in `internal/core/checkout_test.go` to document the current data-loss behavior when checking out a file that has uncommitted local modifications.

The test confirms that `CheckoutFile` blindly overwrites the working directory file with the version from the last commit, validating the reported issue.

Additionally, this PR includes a fix for `SafeWrite` in `internal/core/helpers.go` to ignore directory sync errors, which is necessary for the tests to run correctly on Windows environments.

# Changes

- **[NEW] `internal/core/checkout_test.go`**: Added `Test_CheckoutFile_OverwritesDirtyFile` to reproduce the bug.
- **[MODIFY] `internal/core/helpers.go`**: Updated `SafeWrite` to handle directory sync failures gracefully on Windows.

# Verification

Run the test:
```bash
go test ./internal/core -run Test_CheckoutFile_OverwritesDirtyFile -v
```
The test should PASS, confirming the destructive behavior.
