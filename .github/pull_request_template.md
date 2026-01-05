# Pull Request

## 1. PR Type (MANDATORY)

Select **exactly one**.

- [ ] **feat** – New user-facing command, flag, or engine capability
- [ ] **fix** – Bug fix correcting existing behavior
- [ ] **test** – Test-only changes (no production code)
- [ ] **chore** – Refactor, docs, tooling, or cleanup (no behavior change)

> ❗ PRs that do not clearly fit one category will be closed.

---

## 2. Description (WHAT changed)

Describe **what changed**, not why it is good.

- Commands / files / subsystems affected:
- Public API, CLI flags, or behavior changes (if any):

---

## 3. Intent Declaration (CRITICAL)

Answer all that apply.

**Does this PR change any user-facing command or flag?**

- [ ] Yes
- [ ] No

**Does this PR change data formats, hashing, refs, or repo state?**

- [ ] Yes
- [ ] No

**Does this PR introduce or modify filesystem interactions?**

- [ ] Yes
- [ ] No

If you answered “Yes” to any of the above, explain briefly:

```

<explanation>
```

---

## 4. Documentation Impact

- [ ] This PR does NOT change documentation
- [ ] This PR updates documentation to reflect behavior changes
- [ ] This PR is documentation-only

If documentation was updated, specify files:

```
<list files>
```

---

## 5. Test Accountability (MANDATORY)

### Test Type Used

Select all that apply.

- [ ] **Unit tests** (pure logic, no disk, no `os.Chdir`, no `t.TempDir`)
- [ ] **Integration tests** (filesystem + repo state)
- [ ] No tests (only valid for **docs / chore** PRs)

### Test Details

- Test files added or modified:
- What behavior is proven by tests:
- What behavior is **explicitly untested** (if any):

```
<details>
```

> ❗ Unit tests that touch disk or process state will be rejected.
> ❗ Fix PRs **must** include a regression test.

---

## 6. Git-Parity Risk Assessment (MANDATORY for feat / fix)

Answer **Yes / No** and explain if Yes.

- Could this PR cause KitKat behavior to diverge from Git?
- Does this affect commit graphs, refs, hashes, or object semantics?
- Is this change expected to impact future `.git` compatibility?

```
<risk analysis>
```

---

## 7. Verification Steps (REQUIRED)

List **exact steps** a reviewer can follow to verify this PR.

Examples:

- Commands run
- Tests executed
- Files inspected

```
1.
2.
3.
```

---

## 8. Issue Linkage

- Related Issue(s): `Fixes #___` / `Refs #___`
- If no issue exists, explain why:

```
<explanation>
```

---

## 9. Final Checklist (NO GUESSING)

Select **exactly one** formatting option.

- [ ] I have run `go fmt ./...` (required for all Go code changes)
- [ ] This PR contains no Go code changes (docs / diagrams only)

Confirm all that apply:

- [ ] PR type correctly selected
- [ ] Test classification (unit vs integration) is accurate
- [ ] No behavior change hidden as chore
- [ ] All acceptance criteria in linked issues are met

---

### Reminder

> If this PR changes behavior, it must say so
> If it adds tests, it must classify them correctly
> If it hides risk, it will be rejected
