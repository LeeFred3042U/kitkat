# Pull Request

## 1. PR Type (MANDATORY)

Select **exactly one**

* [ ] **feat** – New user-facing command, flag, or engine capability
* [ ] **fix** – Bug fix correcting existing behavior
* [ ] **test** – Test-only changes (no production code)
* [ ] **chore** – Refactor, docs, tooling, or cleanup (no behavior change)

> ❗ PRs that do not clearly fit one category will be closed.

---

## 2. Scope Guard (MANDATORY)

This PR is **strictly limited** to the following areas

**Allowed files / directories:**

```
- 
```

**Explicitly NOT allowed to change:**

```
- 
```

> ❗ Changes outside the declared scope require a new PR
> ❗ “While I was here” changes will be rejected

## Pre-PR Checklist (MANDATORY)

Complete these checks before opening a PR. PRs that fail these checks will be closed.

- [ ] I have fetched `origin` and integrated `origin/main` into my branch (rebase preferred):
  - `git fetch origin`
  - `git rebase origin/main`  # preferred
  - or `git merge origin/main`  # allowed if you prefer merge
  - Verify: `git rev-list --left-right --count origin/main...HEAD` (should show your branch is ancestor/behind/ahead as expected)

- [ ] I confirm the base branch for this PR is: `__________` (fill in)

- [ ] If the base branch is `main`, I have squashed my changes into **exactly one commit** and force-pushed:
  - Interactive rebase and squash: `git rebase -i origin/main` → squash into one commit
  - Force-push safely: `git push --force-with-lease origin <branch>`
  - Verify single-commit: `git rev-list --count origin/main..HEAD` should equal `1`

Notes:
- PRs targeting `main` with more than one commit will be closed.
- This checklist is enforced by maintainers  
- Do not open a PR to `main` without satisfying these items.


---

## 3. Description (WHAT changed)

Describe **what changed**, not why it is good

* Commands / files / subsystems affected:
* User-visible behavior or CLI changes (if any):

```
<description>
```

---

## 4. Intent Declaration (CRITICAL)

Answer all that apply

**Does this PR change any user-facing command or flag?**

* [ ] Yes
* [ ] No

**Does this PR change data formats, hashing, refs, or repo state?**

* [ ] Yes
* [ ] No

**Does this PR introduce or modify filesystem interactions?**

* [ ] Yes
* [ ] No

If you answered “Yes” to any of the above, explain briefly:

```
<explanation>
```

---

## 5. Storage & Repo Safety Check (MANDATORY)

Confirm all that apply:

* [ ] This PR does NOT write to `.kitkat/objects`
* [ ] This PR does NOT change index format or index location
* [ ] This PR does NOT change hashing behavior
* [ ] This PR does NOT add new object types

If **any** box is unchecked:

* A design issue **must** be linked
* Migration or rollback notes **must** be included

---

## 6. Backward Compatibility

Does this PR change behavior for existing kitkat repositories?

* [ ] No
* [ ] Yes

If **Yes**, specify impact:

* [ ] Existing repos break immediately
* [ ] Existing repos break only for specific commands
* [ ] Migration path provided

```
<compatibility notes>
```

---

## 7. Documentation Impact

* [ ] This PR does NOT change documentation
* [ ] This PR updates documentation to reflect behavior changes
* [ ] This PR is documentation-only

If documentation was updated, list files:

```
<files>
```

---

## 8. Test Accountability (MANDATORY)

### Test Type Used

Select all that apply.

* [ ] **Unit tests** (pure logic only)
* [ ] **Integration tests** (filesystem, repo state, or disk)
* [ ] No tests (only valid for **docs / chore** PRs)

> ❗ Unit tests must NOT touch disk or process state.
> ❗ Any test touching filesystem or repo state **must** be classified as integration.
> ❗ Fix PRs **must** include a regression test.

---

### Test Expectations (REQUIRED)

This PR proves the following invariants:

1.
2.

This PR explicitly does **NOT** test:

1.

## Failure modes covered by tests:

## Failure modes NOT covered:

---

## 9. Git-Parity Risk Assessment (MANDATORY for feat / fix)

Answer **Yes / No** and explain if Yes

* Could this PR cause kitkat behavior to diverge from Git?
* Does this affect commit graphs, refs, hashes, or object semantics?
* Is this change expected to impact future `.git` compatibility?

```
<risk analysis>
```

---

## 10. Verification Steps (REQUIRED)

Exact steps a reviewer can follow to verify this PR

```
1.
2.
3.
```

---

## 11. Issue Linkage

* Related Issue(s): `Fixes #___` / `Refs #___`

If no issue exists, explain why:

```
<explanation>
```

---

## 12. Final Checklist (NO GUESSING)

Select **exactly one**:

* [ ] I have run `go fmt ./...`
* [ ] This PR contains no Go code changes

Confirm all that apply:

* [ ] PR type correctly selected
* [ ] Scope guard respected
* [ ] Intent declaration is accurate
* [ ] Test classification is correct
* [ ] No behavior change hidden as chore
* [ ] I have synced my branch with `origin/main`
* [ ] All acceptance criteria in linked issues are met
* [ ] If this PR targets `main`, it contains exactly one commit (squashed)

---

## 13. Reviewer Kill Conditions (Read Carefully)

This PR may be **closed without merge** if:

* Scope guard is violated
* Behavior changes are undeclared
* Tests do not prove stated invariants
* Storage / index / object risk is understated
* Git-parity risk is hand-waved

---

### Reminder

> If this PR changes behavior, it must say so
> If it touches storage, it must admit it
> If it relies on tests, they must prove invariants
> If it hides risk, it will be rejected
