---
title: 홈
---

# Specvital 문서

> 🇺🇸 [English Documentation](/en/)

Specvital 문서에 오신 것을 환영합니다. Specvital은 코드 리뷰 프로세스를 개선하기 위해 설계된 오픈소스 테스트 커버리지 인사이트 도구입니다.

## 문서 구조

### [PRD (Product Requirements Document)](./prd/)

Specvital 플랫폼의 제품 사양 및 요구사항 문서입니다.

- [제품 개요](./prd/00-overview.md) - 제품 비전, 타겟 사용자, GTM 전략
- [아키텍처](./prd/01-architecture.md) - 시스템 아키텍처 및 서비스 구성
- [코어 엔진](./prd/02-core-engine.md) - 테스트 파서 라이브러리 설계
- [웹 플랫폼](./prd/03-web-platform.md) - 웹 대시보드 및 REST API
- [컬렉터 서비스](./prd/04-collector-service.md) - 백그라운드 분석 워커
- [데이터베이스 설계](./prd/05-database-design.md) - 데이터베이스 스키마 및 설계
- [기술 스택](./prd/06-tech-stack.md) - 기술 선택 및 근거

### [ADR (Architecture Decision Records)](./adr/)

Specvital 개발 중 내린 아키텍처 결정에 대한 문서입니다.

- [ADR 개요](./adr/) - 아키텍처 결정 기록 소개
- [정적 분석 접근법](./adr/01-static-analysis-approach.md)
- [경쟁 차별화](./adr/02-competitive-differentiation.md)
- [파서 라이브러리 분리](./adr/03-parser-library-separation.md)
- [API 워커 서비스 분리](./adr/04-api-worker-service-separation.md)
- [큐 기반 비동기 처리](./adr/05-queue-based-async-processing.md)
- [리포지토리 전략](./adr/06-repository-strategy.md)
- [PaaS 우선 인프라](./adr/07-paas-first-infrastructure.md)
- [공유 인프라](./adr/08-shared-infrastructure.md)

### [아키텍처 개요](./architecture.md)

상위 수준의 시스템 아키텍처 문서입니다.

## 관련 리포지토리

Specvital 플랫폼은 여러 리포지토리로 구성되어 있습니다:

- [specvital/core](https://github.com/specvital/core) - 파서 엔진
- [specvital/web](https://github.com/specvital/web) - 웹 플랫폼
- [specvital/collector](https://github.com/specvital/collector) - 워커 서비스
- [specvital/infra](https://github.com/specvital/infra) - 인프라 및 스키마

## 기여하기

이것은 Specvital의 메인 문서 리포지토리입니다. 기여 가이드라인은 각 리포지토리의 CONTRIBUTING.md 파일을 참조해주세요.

## 라이선스

자세한 내용은 [LICENSE](https://github.com/specvital/.github/blob/main/LICENSE)를 참조하세요.
