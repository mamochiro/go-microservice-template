---
description: Run full check — build, lint, and test before committing.
---

Run a full pre-commit check:

1. Run `task gen` to ensure generated code is up to date
2. Run `task build` to verify the build succeeds
3. Run `task lint` to check for lint issues (Go + Web)
4. Run `task test` to run all unit tests

Report a summary of results. If anything fails, investigate and suggest fixes.
