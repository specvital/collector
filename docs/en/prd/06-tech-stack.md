---
title: Tech Stack
---

# Technology Stack

> 🇰🇷 [한국어 버전](/ko/prd/06-tech-stack.md)

## Summary

| Area     | Choice           | Reason                           |
| -------- | ---------------- | -------------------------------- |
| Parser   | Go + Tree-sitter | High performance, multi-language |
| Backend  | Go               | Performance, simple deployment   |
| Frontend | React (Next.js)  | Ecosystem, SSR                   |
| Queue    | Redis-based      | Simple, retry support            |
| DB       | PostgreSQL       | Versatility, stability           |
| Deploy   | PaaS             | DX priority                      |

## Technical Principles

1. **Type safety**: Compile-time validation
2. **Serverless first**: Minimize initial costs
3. **Avoid lock-in**: Choose standard technologies

## Risk Management

- Establish per-service migration plans
- Monitor vendor dependencies

> See go.mod, package.json in each repository for versions and details
