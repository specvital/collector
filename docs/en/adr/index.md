---
title: ADR
---

# Architecture Decision Records (ADR)

> 🇰🇷 [한국어 버전](/ko/adr/)

Documentation of architectural decisions made in the Specvital project.

## What is an ADR?

An Architecture Decision Record (ADR) captures an important architectural decision made along with its context and consequences. ADRs help maintain decision history across multi-repository microservices.

## When to Write an ADR

| Category         | Examples                                                |
| ---------------- | ------------------------------------------------------- |
| Technology Stack | Framework selection, library adoption, version upgrades |
| Architecture     | Service boundaries, communication patterns, data flow   |
| API Design       | Endpoint structure, versioning strategy, error handling |
| Database         | Schema design, migration strategy, indexing approach    |
| Infrastructure   | Deployment platform, scaling strategy, monitoring       |
| Cross-cutting    | Security, performance optimization, observability       |

## Templates

| Template                     | Use Case                        |
| ---------------------------- | ------------------------------- |
| [template.md](./template.md) | Standard ADR for most decisions |

## Naming Convention

```
XX-brief-decision-title.md
```

- `XX`: Two-digit sequential number (01, 02, ...)
- Lowercase with hyphens
- Brief and descriptive titles

## Technical Areas

| Area           | Affected Repositories |
| -------------- | --------------------- |
| Parser         | core                  |
| API            | web                   |
| Worker         | collector             |
| Database       | infra                 |
| Infrastructure | infra                 |
| Cross-cutting  | multiple              |

## ADR Index

| #   | Title                                                                       | Area           | Date       |
| --- | --------------------------------------------------------------------------- | -------------- | ---------- |
| 01  | [Static Analysis-Based Instant Analysis](./01-static-analysis-approach.md)  | Cross-cutting  | 2024-12-17 |
| 02  | [Competitive Differentiation Strategy](./02-competitive-differentiation.md) | Cross-cutting  | 2024-12-17 |
| 03  | [Parser Library Separation](./03-parser-library-separation.md)              | Parser         | 2024-12-17 |
| 04  | [API and Worker Service Separation](./04-api-worker-service-separation.md)  | Architecture   | 2024-12-17 |
| 05  | [Queue-Based Asynchronous Processing](./05-queue-based-async-processing.md) | Architecture   | 2024-12-17 |
| 06  | [Polyrepo Repository Strategy](./06-repository-strategy.md)                 | Architecture   | 2024-12-17 |
| 07  | [PaaS-First Infrastructure Strategy](./07-paas-first-infrastructure.md)     | Infrastructure | 2024-12-17 |
| 08  | [Shared Infrastructure Strategy](./08-shared-infrastructure.md)             | Infrastructure | 2024-12-17 |

## Process

1. **Create**: Copy [template.md](./template.md) → `XX-title.md`
2. **Write**: Fill in all sections with finalized decision
3. **Localize**: Create Korean version in `kr/adr/`
4. **Review**: Submit PR for team review
5. **Merge**: Add to index after approval

## Related Repositories

- [specvital/core](https://github.com/specvital/core) - Parser engine
- [specvital/web](https://github.com/specvital/web) - Web platform
- [specvital/collector](https://github.com/specvital/collector) - Worker service
- [specvital/infra](https://github.com/specvital/infra) - Infrastructure
