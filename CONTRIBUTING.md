# Contributing to kitkat

Thank you for your interest in contributing to **kitkat**! We are building an educational reimplementation of Git in Go, and we welcome contributors of all experience levels.

Whether you are a student looking for your first open-source contribution or a systems engineer wanting to solve complex algorithmic challenges, there is a place for you here.

## Ways to Contribute

We have organized contributions into three tracks. Please choose one that matches your interest:

## 🚦 Pick Your Track

We use labels to show how hard an issue is. Choose one that fits you.

### Track 1: Easy (Start Here)
**Labels:** `Easy` `documentation` `good first issue`
Great if you are new to Go or Open Source
* **What you do:** Fix typos, add simple CLI commands, update `README.md`
* **Example:** "Fix the help message for `rm`"

### Track 2: Medium (The Real Work)
**Labels:** `Medium` `bug` `core`
Great if you know some Go and want to build features.
* **What you do:** Add new logic, fix standard bugs, handle flags.
* **Example:** "Implement `kitkat log -n 5`"

### Track 3: Hard (Core Logic)
**Labels:** `Hard` `core`
For people who can think it through
* **What you do:** Graph traversal, file locking, hashing, binary formats
* **Example:** "Implement `reset --hard` with tree traversal"

> Note: This is not applied everytime
---


## Workflow (How to Contribute)

- Go to the Issues tab or click [this](https://github.com/LeeFred3042U/kitkat/issues). Look for a label (Easy, Medium, etc) 
- Comment: "I want to work on this!"
- When you are assigned the issue by the maintainer only then continue


### Prerequisites
* **Go 1.24+** installed (Check `go.mod` for the exact version)
* A text editor (VS Code recommended)

## Setup

### 1. Install Go

* Make sure you have **Go 1.24+**.
* Check with: `go version`

1. **Fork** this repo (Click the button top-right)
    - it looks like this

    - ![alt text](image-1.png)

2. **Clone** your fork:
```bash
 git clone https://github.com/username/kitkat.git
 cd kitkat
```


2. **Create a Branch:** 
    - Use a descriptive name for your branch
    - Do not work on main 
    - Make a new branch
```bash
git checkout -b feat/implement-rm-command
# or
git checkout -b docs/add-status-diagram
```

3. **Build** the project

```bash
go build -o kitkat ./cmd/main.go
```

4. **Verify** if it runs

```bash
./kitkat init
./kitkat help
```

5. **Make Changes(for code):** 
    - Write clean, idiomatic Go code 
    - If you are new to Go, feel free to ask for help in the PR!
   **documentation:**
    - Work as stated in the issue
    - keep check for typos


6. **Test:** Manual testing is required 
    - Please include (if code changes were made, else no need) a **screenshot** or **terminal output** or **screen recording** in your Pull Request description proving the command works as expected
    - run `go fmt ./...`
    - before you commit, else we have issues

5. **Push & PR:** 
    - Go to GitHub and open a Pull Request, keep the description concise, and reference the issue number (e.g., `Fixes #1`)
    - The title should be named as the issue title which is fixed by you


## Pull Request Verification Standard (MANDATORY) [only for code changes]

We require **Proof of Work** for every Pull Request 
"It works on my machine" is not enough
You must include a **Screenshot** or **Terminal Output** in your PR description showing the command running successfully.

**Acceptable Example (Terminal Output):**

> I tested the `help` command. Here is the output of terminal showing it
```bash
[terminal@terminal kitkat] $ ./kitkat help
usage: kitkat <command> [arguments]

These are the common KitKat commands:
   tag        Create a new tag for a commit
   merge      Merge a branch into the current branch.
   ls-files   Show information about files in the index
   config     Get and set repository or global options.
   commit     Record changes to the repository.
   log        Show the commit history
   clean      Remove untracked files from the working directory
   init       Initialize a new KitKat repository
   add        Add file contents to the index.
   diff       Show changes between the last commit and staging area

Use 'kitkat help <command>' for more information about a command
```
OR
> I tested `add` and `commit` command(since both go together). Here is the output of terminal shown in a screenshot

![alt text](image.png)

---

# Editing & Creating Architecture Diagrams (PlantUML)

KitKat’s architecture diagrams are stored in `.puml` format and exported as `.png`.
If you are **creating new diagrams**, follow the same workflow used for editing existing ones.

All source files live under:

```
docs/architecture/<section>/
```

Each diagram should always consist of:

```
diagram-name.puml   (source file)
diagram-name.png    (exported image checked into the repo)
```

## Required Tool (VS Code)

We use the following extension so contributors can preview and export diagrams:

* **Name:** PlantUML Viewer
* **ID:** `BenkoSoftware.plantumlviewer`
* **Publisher:** BenkoSoftware
* **Version:** 1.1.0

## How to Use It

1. Install the extension
2. Press `Ctrl + Shift + P`, search **PlantUML**
3. Add keybindings for(because it makes it easier):
* *Open Preview*
* *Export as PNG*



## Workflow for New or Updated Diagrams

1. Create or edit the `.puml` file
2. Open the preview to confirm the diagram renders correctly
3. Export to PNG
4. Commit both files inside the architecture directory following this structure:
```
docs/
└── architecture/
    └── <section-name>/
        ├── <diagram-name>.puml
        └── <diagram-name>.png
```



Pull Requests missing the PNG export will be rejected.

---

## Code of Conduct

Please note that this project is released with a [Code of Conduct](https://www.google.com/search?q=CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## License
By contributing, you agree that your contributions will be licensed under the project's MIT License.
