---
description: Run tests. Use "unit", "integration", or "all" as argument.
---

Run tests based on the argument: $ARGUMENTS

- If "unit" or empty: run `task test`
- If "integration": run `task integration-test`
- If "all": run both unit and integration tests
- If a specific path or function name is given: run `go test -v -run $ARGUMENTS ./...`

Show the results and summarize any failures.
