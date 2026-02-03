# Project Bugs

## [x] Bug 1: Summarizing an empty repo praises nonexistent maintainers

**Status:** Closed (not a bug)

**Description:** When pointed at a repo with zero open issues, Claude responds with "Great job keeping your issue tracker clean!" even though the repo has 4,000 open PRs disguised as "discussions."

**Steps to reproduce:**
1. Run `gitissuesum owner/abandoned-repo`
2. Repo has 0 issues but is clearly on fire
3. Read Claude's glowing praise

**Expected:** Honest assessment or silence

**Actual:** "The maintainers are doing an exceptional job" — Claude, lying through its tokens

---
