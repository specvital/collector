---
title: Home
---

# Specvital Documentation

> 🇰🇷 [한국어 문서](/ko/)

Welcome to the Specvital documentation. Specvital is an open-source test coverage insights tool designed to enhance code review processes.

## Documentation Structure

### [PRD (Product Requirements Document)](./prd/)

Product specifications and requirements documentation for the Specvital platform.

- [Product Overview](./prd/00-overview.md) - Product vision, target users, and GTM strategy
- [Architecture](./prd/01-architecture.md) - System architecture and service composition
- [Core Engine](./prd/02-core-engine.md) - Test parser library design
- [Web Platform](./prd/03-web-platform.md) - Web dashboard and REST API
- [Collector Service](./prd/04-collector-service.md) - Background analysis worker
- [Database Design](./prd/05-database-design.md) - Database schema and design
- [Tech Stack](./prd/06-tech-stack.md) - Technology choices and rationale

### [ADR (Architecture Decision Records)](./adr/)

Documentation of architectural decisions made during the development of Specvital.

- [ADR Overview](./adr/) - Introduction to architecture decision records
- [Static Analysis Approach](./adr/01-static-analysis-approach.md)
- [Competitive Differentiation](./adr/02-competitive-differentiation.md)
- [Parser Library Separation](./adr/03-parser-library-separation.md)
- [API Worker Service Separation](./adr/04-api-worker-service-separation.md)
- [Queue-Based Async Processing](./adr/05-queue-based-async-processing.md)
- [Repository Strategy](./adr/06-repository-strategy.md)
- [PaaS-First Infrastructure](./adr/07-paas-first-infrastructure.md)
- [Shared Infrastructure](./adr/08-shared-infrastructure.md)

### [Architecture Overview](./architecture.md)

High-level system architecture documentation.

## Related Repositories

The Specvital platform is composed of multiple repositories:

- [specvital/core](https://github.com/specvital/core) - Parser engine
- [specvital/web](https://github.com/specvital/web) - Web platform
- [specvital/collector](https://github.com/specvital/collector) - Worker service
- [specvital/infra](https://github.com/specvital/infra) - Infrastructure and schema

## Contributing

This is the main documentation repository for Specvital. For contribution guidelines, please refer to each repository's CONTRIBUTING.md file.

## License

See [LICENSE](https://github.com/specvital/.github/blob/main/LICENSE) for details.
