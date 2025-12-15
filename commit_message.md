# Commit Message

## 한국어 버전

```
refactor: 인프라 패키지를 internal/infra로 이동

Clean Architecture 리팩토링 1단계로 인프라 관심사를
비즈니스 로직과 분리

- config, db, queue 패키지를 internal/infra/로 이동
- 관련 import 경로 업데이트 (main.go, repository)
- sqlc.yaml 출력 경로 수정
```

## English Version

```
refactor: move infrastructure packages to internal/infra

Separate infrastructure concerns from business logic as step 1 of Clean Architecture refactoring

- Move config, db, queue packages to internal/infra/
- Update related import paths (main.go, repository)
- Update sqlc.yaml output paths
```
