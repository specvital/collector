# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SpecVital Collector - Worker service for analyzing test files in GitHub repositories

- Queue-based async worker (River on PostgreSQL)
- Dual-binary: Worker (scalable) + Scheduler (singleton)
- External parser: `github.com/specvital/core`

## Documentation Map

| Context                         | Reference        |
| ------------------------------- | ---------------- |
| Architecture / Data flow        | `docs/en/`       |
| Design decisions (why this way) | `docs/en/adr/`   |
| Coding rules / Test patterns    | `.claude/rules/` |

## Commands

Before running commands, read `justfile` or check available commands via `just --list`

## Project-Specific Rules

### Auto-Generated Files (NEVER modify)

- `src/internal/infra/db/{queries.sql.go,models.go,db.go}`
- Workflow: `just dump-schema` → `just gen-sqlc`

### External Dependency

- Parsing logic lives in `github.com/specvital/core`, NOT here
- For parser changes → open issue in core repo first

### Dual Binary Architecture

- **Worker** (`cmd/worker`): horizontally scalable, queue consumer
- **Scheduler** (`cmd/scheduler`): single instance only (distributed lock)
- Must remain separate for Railway deployment - NEVER merge

## Common Workflows

### DB Schema Changes

1. Modify schema in specvital-postgres repo
2. `just dump-schema` → `just gen-sqlc`
3. Update `adapter/repository/` implementation

### Adding New Worker

1. Define worker in `adapter/queue/`
2. Register in `app/container.go`
3. Write tests
