# kitkat

A toy Git clone written in Go for learning the fundamentals of version control

## Getting Started

Follow these steps to get `kitkat` up and running on your local machine
> Note: Doesn't support any remote commands

### Prerequisites

Before you begin, make sure you have the **Go programming language** installed
You can download it from the [official Go website](https://go.dev/dl/)

### Installation

1.  **Clone the Repository**
    First, get a local copy of the project using `git clone`

    ```sh
    git clone https://github.com/LeeFred3042U/kitkat.git
    ```

2.  **Navigate into the Directory**
    Move into the project folder you just cloned.

    ```sh
    cd kitkat
    ```

3.  **Build the Executable**
    Compile the Go source code. This command creates a runnable program named `kitkat` in your current directory

    ```sh
    go build -o kitkat ./cmd/main.go
    ```

    **Note:** If you do not add the executable to your system's PATH (step 4), you must run the command with `./kitkat` from within this directory (e.g., `./kitkat init`)

4.  **(Optional) Add to Your System's PATH**
    To use the `kitkat` command from anywhere, move the executable to a location in your system's `PATH`

    ```sh
    # For Linux or macOS
    sudo mv kitkat /usr/local/bin/
    ```

### First-Time Configuration

Before you start using `kitkat`, you should set your name and email. This information will be used in your commits

```sh
# Set your name
kitkat config --global user.name "Your Name"

# Set your email
kitkat config --global user.email "you@example.com"
```

-----

## Flags and Help

You can get help directly from the command line

**General Help**
To see a list of all available commands and their summaries, use the `help` command

```sh
kitkat help
```

**Specific Command Help**
To get detailed usage information for a specific command, add the command's name after `help`

```sh
kitkat help add
kitkat help commit
```

-----

## Features

  * **Repository Initialization**: Create a new `.kitkat` repository
  * **Staging Area**: Add new, modified, or deleted files to the index
  * **Commit History**: View a log of all commits
  * **Branching & Merging**: Create and switch between branches, and perform fast-forward merges
  * **Status & Diff**: Check the status of your working directory and view colorized diffs
  * **Global Config**: Set user information like name and email

-----

## Core Concept

KitKat is an educational toy project designed to mimic the core functionality of Git. It operates on the same fundamental principles of version control, taking **snapshots** of your project. Each snapshot is built from three key objects:

  * **Blobs:** The content of your files
  * **Trees:** The directory structure that organizes blobs
  * **Commits:** A pointer to a tree, representing the state of your project at a specific point in time, linked to a parent commit to form a history

-----

## Command Reference

### `init`

Initializes a new, empty KitKat repository in the current directory.

```sh
kitkat init
```

### `add`

Adds file contents to the staging area (index). It can stage specific files or all changes.

```sh
# Stage a specific file
kitkat add <file-path>

# Stage all new, modified, and deleted files
kitkat add --all
```

### `log`

Shows the commit history for the current branch.

```sh
# Show the detailed, multi-line history
kitkat log

# Show a compact, single-line view
kitkat log --oneline
```

### `status`

Displays the state of the working directory and the staging area. It shows which files are staged, unstaged, and untracked.

```sh
kitkat status
```

### `diff`

Shows the colorized differences between the last commit and the staging area.

```sh
kitkat diff
```

### `branch`

Manages branches. Running it without arguments lists all branches. Providing a name creates a new branch.

```sh
# List all local branches
kitkat branch

# Create a new branch named 'new-feature'
kitkat branch new-feature
```

### `checkout`

Switches branches, checks out a specific commit (detached HEAD), or restores a file to its last committed state.

```sh
# Switch to the 'new-feature' branch
kitkat checkout new-feature

# Checkout a specific commit hash (detached HEAD)
kitkat checkout <commit-hash>

# Revert a file to its state in the last commit
kitkat checkout <file-path>
```

### `merge`

Joins another branch's history into the current branch. Currently only supports fast-forward merges.

```sh
kitkat merge <branch-name>
```

### `ls-files`

Shows a simple list of all files currently in the staging area.

```sh
kitkat ls-files
```

### `clean`

Removes untracked files from the working directory. Requires a `-f` flag for safety.

```sh
# Show which files would be removed (dry run)
kitkat clean

# Forcefully remove untracked files
kitkat clean -f
```

### `config`

Gets or sets the global user configuration, such as name and email.

```sh
# Set your name
kitkat config --global user.name "Your Name"

# Get your name
kitkat config --global user.name

# Set your email
kitkat config --global user.email "you@example.com"

# Get your name
kitkat config --global user.email
```
