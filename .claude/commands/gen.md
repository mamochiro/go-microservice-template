---
description: Regenerate code (Wire DI, Mocks, Swagger). Optionally specify "wire", "mock", or "swagger".
---

Regenerate code based on argument: $ARGUMENTS

- If empty or "all": run `task gen` (Wire + Mocks + Swagger)
- If "wire": run `task wire`
- If "mock": run `task mock`
- If "swagger": run `task swagger`

After generation, run `task lint-api` to verify no lint issues were introduced.
