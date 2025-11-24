# Contributing to kitkat

Thank you for your interest in contributing to **kitkat**! We are building an educational reimplementation of Git in Go, and we welcome contributors of all experience levels.

Whether you are a student looking for your first open-source contribution or a systems engineer wanting to solve complex algorithmic challenges, there is a place for you here.

## Ways to Contribute

We have organized contributions into three tracks. Please choose one that matches your interest:

### Track 1: Beginner Code (Good First Issues)
Perfect for learning Go and getting comfortable with the codebase.
* **Goal:** Implement simple commands that interact with the existing index/storage.
* **Examples:** `kitkat rm`, `kitkat mv`, `kitkat config --list`, `kitkat tag --list`.
* **Skills:** Basic Go, File I/O.

### Track 2: Visual Documentation (No Code)
Perfect for those who want to understand the architecture without writing Go.
* **Goal:** Create PlantUML diagrams explaining how KitKat works.
* **Examples:** Diagramming the "Add" flow, the "Snapshot" storage model, or the "Myers Diff" logic.
* **Skills:** Logic, PlantUML, Technical Writing.

### Track 3: Systems Engineering (Advanced)
For contributors who want deep technical challenges.
* **Goal:** Implement core version control algorithms.
* **Examples:** `kitkat rebase`, `kitkat cherry-pick`, `kitkat blame`, `kitkat stash`.
* **Skills:** Algorithms (Graphs/Trees), Data Structures, Systems Programming.

---

## Getting Started

### Prerequisites
* Go 1.22+ installed
* A text editor (VS Code recommended)

### Setup
1. **Fork** this repository.
2. **Clone** your fork:
   ```bash
   git clone https://github.com/LeeFred3042U/kitkat.git
   cd kitkat
   ```

## Build the project

```bash
go build -o kitkat ./cmd/main.go
````

## Verify it runs

```bash
./kitkat init
./kitkat help
```

## Workflow

**Find an Issue:** Hunt for labels like `good first issue` or `documentation`.

**Create a Branch:** Name it like you actually care.

```
feat/implement-rm-command
docs/add-merge-diagram
```

**Make Changes:** Write clean, idiomatic Go code instead of whatever fever dream comes first.

**Test:** Manual testing for now. Drop a screenshot or terminal output in your PR proving the command works.

**Push & PR:** Open a Pull Request, keep the description tight, and reference the issue number.

## Code Style

Run:

```bash
go fmt ./...
```

Keep functions short enough

Comment any logic that makes you squint, especially in `internal/core`.

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
3. Add keybindings for:

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
## License
By contributing, you agree that your contributions will be licensed under the project's MIT License.
