---
description: Database migration commands. Use "up", "down", or "new <name>".
---

Run migration based on argument: $ARGUMENTS

- If "up": run `task migrate-up`
- If "down": run `task migrate-down` (rollback 1 step)
- If "new <name>": run `task migrate-new -- <name>` to create new migration files, then open both the up and down SQL files for editing
