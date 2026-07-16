# Copyright (C) 2026 Energy Project Team
# SPDX-License-Identifier: AGPL-3.0-only
dc = docker compose

up:
	$(dc) down
	$(dc) up -d

build:
	$(dc) build

down:
	$(dc) down

stop:
	$(dc) stop

restart:
	$(dc) restart

update:
	git pull
	$(dc) build
	$(dc) up -d
	$(dc) logs -f -t

logs:
	$(dc) logs -f -t

logs-file:
	$(dc) logs -t > compose.log
