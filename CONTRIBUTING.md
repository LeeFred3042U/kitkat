# Contributing to kitkat

Thank you for your interest in contributing to **kitkat**! We are building an educational reimplementation of Git in Go, and we welcome contributors of all experience levels.

Whether you are a student looking for your first open-source contribution or a systems engineer wanting to solve complex algorithmic challenges, there is a place for you here.

## Ways to Contribute

We have organized contributions into three tracks. Please choose one that matches your interest:

### Track 1: Beginner Code (Good First Issues)
Perfect for learning Go and getting comfortable with the codebase.
* **Goal:** Implement simple commands that interact with the existing index/storage.
* **Examples:** `kitkat rm`, `kitkat mv`, `kitkat log -n`, `kitkat version`.
* **Skills:** Basic Go, File I/O.

### Track 2: Visual Documentation (No Code)
Perfect for those who want to understand the architecture without writing Go.
* **Goal:** Create PlantUML diagrams explaining how KitKat works.
* **Examples:** Diagramming the "Branching" logic, the "Status" categorization flow, or the "Myers Diff" algorithm.
* **Skills:** Logic, PlantUML, Technical Writing.

### Track 3: Intermediate
For contributors who want deep technical challenges.
* **Goal:** Implement core version control algorithms.
* **Examples:** `kitkat reset --hard`, `kitkat stash`, `kitkat blame`.
* **Skills:** Algorithms (Graphs/Trees), Data Structures, Systems Programming.

---

## Getting Started

### Prerequisites
* **Go 1.24+** installed (Check `go.mod` for the exact version)
* A text editor (VS Code recommended)

### Setup
1. **Fork** this repository.
2. **Clone** your fork:
```bash
 git clone https://github.com/LeeFred3042U/kitkat.git
 cd kitkat
```

### Build the project

```bash
go build -o kitkat ./cmd/main.go
```

### Verify it runs

```bash
./kitkat init
./kitkat help
```

## Workflow

1. **Find an Issue:** Hunt for labels like `good first issue` or `documentation`.
2. **Create a Branch:** Use a descriptive name for your branch.
```bash
git checkout -b feat/implement-rm-command
# or
git checkout -b docs/add-status-diagram
```


3. **Make Changes:** Write clean, idiomatic Go code. If you are new to Go, feel free to ask for help in the PR!
4. **Test:** Manual testing is required. Please include a **screenshot** or **terminal output** in your Pull Request description proving the command works as expected.
5. **Push & PR:** Open a Pull Request, keep the description concise, and reference the issue number (e.g., `Fixes #1`).

## Code Style

Run the Go formatter before committing:

```bash
go fmt ./...
```

* Keep functions short and readable.
* Comment any logic that is complex, especially in `internal/core`.

## Pull Request Verification Standard (MANDATORY)

We require **Proof of Work** for every Pull Request. "It works on my machine" is not enough.
You must include a **Screenshot** or **Terminal Output** in your PR description showing the command running successfully.

**Acceptable Example (Terminal Output):**

> I tested the `help` command. Here is the output of terminal showing the file being removed:
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

**Unacceptable Examples:**
* "Tested manually."
* "It works."
* (No description at all)

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
