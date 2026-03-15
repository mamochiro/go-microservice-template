---
description: Manage Docker services. Use "up", "down", or "logs <service>".
---

Manage Docker based on argument: $ARGUMENTS

- If "up": run `task docker-up`
- If "down": run `task docker-down`
- If "logs" or "logs <service>": run `docker-compose logs -f $1` to tail logs
- If "status": run `docker-compose ps` to show running containers
