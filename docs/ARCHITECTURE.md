## 1. Repository Initialization

Logic for creating a new repository structure (`.kitcat` folders, HEAD, etc.).

![Init Flow](architecture/init/init.png)

## 2. Staging Flow (Add)

Logic for moving files from the Working Directory to the Index (Staging Area).

![Add Flow](architecture/add/add.png)

## 3. Snapshot Flow (Commit)

Logic for creating a permanent snapshot from the Index.

![Commit Flow](architecture/commit/commit.png)

## Status Command

The `kitcat status` command determines the state of files by comparing
three trees: the HEAD commit, the Index (staging area), and the Working
Directory.

The process is split into two phases:

1. Detecting staged changes by comparing the Index against HEAD.
2. Detecting unstaged and untracked changes by comparing the Working
   Directory against the Index.

![Status Command Flow](architecture/status/status.png)

## Branch Command

Logic for listing and creating branches.

![Branch Flow](architecture/branch/branch.png)
