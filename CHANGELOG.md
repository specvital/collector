# Changelog

## [1.1.0](https://github.com/specvital/collector/compare/v1.0.6...v1.1.0) (2026-01-04)

### 🎯 Highlights

#### ✨ Features

- add Clone-Rename race condition detection ([ebbf443](https://github.com/specvital/collector/commit/ebbf443e3b9271c4dba82d38a567d3efdc0236a9))
- add codebase lookup queries based on external_repo_id ([d6e0b79](https://github.com/specvital/collector/commit/d6e0b797ec46aab29eb897c6c1d6a8cdbb6b47c6))
- add codebase stale handling queries and repository methods ([939e078](https://github.com/specvital/collector/commit/939e07886a8a26c385ad3206c71fe8205a6bd001))
- add GitHub API client for repository ID lookup ([46e40b8](https://github.com/specvital/collector/commit/46e40b806843f7e87ec98269baca7bce136064bd))
- determine repository visibility via reversed git ls-remote order ([0bc988e](https://github.com/specvital/collector/commit/0bc988e839f795678281ac6431b41199e4f95f95))
- integrate codebase resolution case branching into AnalyzeUseCase ([0f58440](https://github.com/specvital/collector/commit/0f58440f0012c15fd215f57a58917370ff93b2a9))
- record user analysis history on analysis completion ([e2b2095](https://github.com/specvital/collector/commit/e2b2095c47dffbd51fa57d3d24c550c19cfed851))
- store commit timestamp on analysis completion ([24bdbd7](https://github.com/specvital/collector/commit/24bdbd7050a40fbcf41965e6abfb728dc9460870))

#### 🐛 Bug Fixes

- add missing focused and xfail TestStatus types ([b24ee33](https://github.com/specvital/collector/commit/b24ee333e5c5e0f71098326f934c09853976fee6))
- add missing is_private column to test schema ([5744b95](https://github.com/specvital/collector/commit/5744b956a7b408f3dbc7583456eaedbf7fa1f4f6))
- ensure transaction atomicity for multi-step DB operations ([16834ef](https://github.com/specvital/collector/commit/16834ef0837917df4c30d31135e4f97a8a07eb3b))
- exclude stale codebases from Scheduler auto-refresh ([933c417](https://github.com/specvital/collector/commit/933c41711f375979d96cf5401ba93e6171891b49))
- fix visibility not being updated on reanalysis ([2424a5f](https://github.com/specvital/collector/commit/2424a5fadaebfd0fed1aba07045f2c86ddd5c585))
- prevent duplicate analysis job enqueue for same commit ([1a996ea](https://github.com/specvital/collector/commit/1a996ea38ad6742647317932e7acbb24939146e1))
- prevent unnecessary job retries on duplicate key error ([40eda32](https://github.com/specvital/collector/commit/40eda32b890206f3f3ef5913ce8ed4f9afdc0cdb))
- resolve stray containers from failed testcontainers cleanup ([1ef5124](https://github.com/specvital/collector/commit/1ef5124a617fcbc1ddd434b6a74baa6dd5ab390a))

#### ⚡ Performance

- improve DB save performance for large repository analysis ([200a527](https://github.com/specvital/collector/commit/200a5275cf639a2c0c65d955e79dbe65ad4f7068))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **devcontainer:** fix network creation failure in Codespaces ([2054227](https://github.com/specvital/collector/commit/2054227927b13127fb2c770323dcc17e6bba4d0a))
- isolate git ls-remote environment to fix private repo misclassification ([7d15fb8](https://github.com/specvital/collector/commit/7d15fb82534cb2c4c34ea368173265c185abf543))

#### 📚 Documentation

- add CLAUDE.md ([5194d71](https://github.com/specvital/collector/commit/5194d713b2f07fd2d4d2a66df62f861520b027bc))
- add missing version headers and improve CHANGELOG hierarchy ([d6436ab](https://github.com/specvital/collector/commit/d6436ab60b12e4bf551c23d59009fa66782e6eb4))
- rename infra repo in docs ([1bdb806](https://github.com/specvital/collector/commit/1bdb806dabc7fd082cb114e93f349aaa619d5315))

#### 💄 Styles

- format code ([8616fbd](https://github.com/specvital/collector/commit/8616fbdae4105860c87569093f302ba6a877c6c7))

#### ♻️ Refactoring

- remove unused deprecated Stop method ([c034ecc](https://github.com/specvital/collector/commit/c034ecc56660bda965a297072e7d23400e8b8e61))
- **test:** auto-sync test schema with production schema ([77668e0](https://github.com/specvital/collector/commit/77668e0e946003dc4f0d3b9e9c086c85b70f8fab))

#### 🔨 Chore

- changing the environment variable name for accessing GitHub MCP ([553c63d](https://github.com/specvital/collector/commit/553c63d358a5b1fd1c607843d41b90544d86330e))
- dump schema ([ba3fc16](https://github.com/specvital/collector/commit/ba3fc165a074f0827417ee6212002e79c9d5340e))
- dump schema ([425b609](https://github.com/specvital/collector/commit/425b6098dc1ee104189a4a33dc635f5e0b9f0352))
- dump schema ([52575e5](https://github.com/specvital/collector/commit/52575e5701088de44401abb227080800250094d8))
- dump schema ([abdaa2e](https://github.com/specvital/collector/commit/abdaa2eda93763d793b2a8a67f6fe2f3b4e14166))
- fix vscode import area not automatically collapsing ([ac92e87](https://github.com/specvital/collector/commit/ac92e87ee1be68a886e4df8b5ed006d0fc8ba0dd))
- improved the claude code status line to display the correct context window size. ([e1fa775](https://github.com/specvital/collector/commit/e1fa775b9dfd49ed57ec5d66aaf0eab4ec0e34b8))
- modified container structure to support codespaces ([0d1fec6](https://github.com/specvital/collector/commit/0d1fec6ec9af2bd3fb1df5a292242e240e13a36e))
- modify local db migration to always initialize the database ([7709a5b](https://github.com/specvital/collector/commit/7709a5b8af0a8fd7bee795ebd533dd5d3944d243))
- sync ai-config-toolkit ([0d00d4a](https://github.com/specvital/collector/commit/0d00d4a615fa3b1c162e8976b0f86b87948f0eaf))
- sync docs ([86772da](https://github.com/specvital/collector/commit/86772da7cb514400b7f7c89ea0defde95241195e))
- update core ([9092761](https://github.com/specvital/collector/commit/9092761f54e28b114b70a7dfbab14e8b82e27bdc))
- update core ([e6613c3](https://github.com/specvital/collector/commit/e6613c3a8e85189621056981ae0e3d91ff266e41))
- update core ([c163ae9](https://github.com/specvital/collector/commit/c163ae92f08d30046712de8c4b86b3162eaae758))

## [1.0.6](https://github.com/specvital/collector/compare/v1.0.5...v1.0.6) (2025-12-19)

### 🎯 Highlights

#### 🐛 Bug Fixes

- resolve 60-second timeout failure during large analysis jobs ([ed18bc3](https://github.com/specvital/collector/commit/ed18bc3f587c0a446ea41b09431126ab1d22bba5))

## [1.0.5](https://github.com/specvital/collector/compare/v1.0.4...v1.0.5) (2025-12-19)

### 🎯 Highlights

#### 🐛 Bug Fixes

- resolve 60-second timeout issue during bulk INSERT operations on NeonDB ([0b6bc9b](https://github.com/specvital/collector/commit/0b6bc9bbef0c14190ec953a14caa5b29da0422d5))

## [1.0.4](https://github.com/specvital/collector/compare/v1.0.3...v1.0.4) (2025-12-19)

### 🔧 Maintenance

#### ♻️ Refactoring

- migrate queue system from asynq(Redis) to river(PostgreSQL) ([9664002](https://github.com/specvital/collector/commit/9664002057ef1f801dd8313e9f081760c3e0af21))

#### 🔨 Chore

- missing changes ([de9c0ec](https://github.com/specvital/collector/commit/de9c0ecaa7e136fe71d70a4a231dbd194ed7a33d))

## [1.0.3](https://github.com/specvital/collector/compare/v1.0.2...v1.0.3) (2025-12-18)

### 🎯 Highlights

#### 🐛 Bug Fixes

- fix git clone failure in runtime container ([f11dfa3](https://github.com/specvital/collector/commit/f11dfa3e4090a6412af71b58a7eca6f081e49d4d))

### 🔧 Maintenance

#### ♻️ Refactoring

- remove unused dead code ([95cee17](https://github.com/specvital/collector/commit/95cee17e1307fe3e0cc23ba7d549292b05c19744))

#### 🔨 Chore

- sync docs ([9007e97](https://github.com/specvital/collector/commit/9007e97bbff365afcae69cfcba2f501732a20c8b))

## [1.0.2](https://github.com/specvital/collector/compare/v1.0.1...v1.0.2) (2025-12-17)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix asynq logs incorrectly classified as error in Railway ([d2180cc](https://github.com/specvital/collector/commit/d2180cc1182a0f1187f7dd63b982fc7816e3be47))

## [1.0.1](https://github.com/specvital/collector/compare/v1.0.0...v1.0.1) (2025-12-17)

### 🎯 Highlights

#### 🐛 Bug Fixes

- enable CGO for go-tree-sitter build ([50b1fea](https://github.com/specvital/collector/commit/50b1fea3c7834bd585c3a23615d8acb5cbae8a5f))

## [1.0.0](https://github.com/specvital/collector/releases/tag/v1.0.0) (2025-12-17)

### 🎯 Highlights

#### ✨ Features

- add adaptive decay logic for auto-refresh scheduling ([8a85854](https://github.com/specvital/collector/commit/8a858547c7b8a5190176253830b862588cda8042))
- add enqueue CLI tool ([5697cb9](https://github.com/specvital/collector/commit/5697cb9533b8dcfd8c7c90fa34d5419b181fa287))
- add focused/xfail to test_status and support modifier column ([cd60233](https://github.com/specvital/collector/commit/cd602333aa9eba34e8a3b21bbcd91040bfe59936))
- add job timeout to prevent long-running analysis jobs ([392b43e](https://github.com/specvital/collector/commit/392b43e7395aa923fd05a96872e5f0c8911a8845))
- add local development mode support to justfile ([2ca2d51](https://github.com/specvital/collector/commit/2ca2d51337cd8c425e54a5f3ffd257bde28e403c))
- add local development services to devcontainer ([a30ca9e](https://github.com/specvital/collector/commit/a30ca9eca3df2aa4a418d94b1525f8ac170ee6c1))
- add OAuth token parameter to VCS Clone interface ([de518d0](https://github.com/specvital/collector/commit/de518d05eb534c97cf5e647168ff9d8a686ae00e))
- add semaphore to limit concurrent git clones ([9ddbc06](https://github.com/specvital/collector/commit/9ddbc06291c8fb41e0e39def7f32edaa5830f1ba))
- add UserRepository for OAuth token lookup ([9a16ec1](https://github.com/specvital/collector/commit/9a16ec1f37031d15612198b27f5316d4ab066225))
- implement analysis pipeline (git clone → parse → DB save) ([66dd262](https://github.com/specvital/collector/commit/66dd2627b6e7d5b46ab7bc4358a1f6178d77cfee))
- implement asynq-based worker basic structure ([4dd16ad](https://github.com/specvital/collector/commit/4dd16ad22b355229f6c9db12e076427f7ff0c2ea))
- initialize collector service project ([1d3c8cf](https://github.com/specvital/collector/commit/1d3c8cf3a570a719301fa953c9cccb9c53a0358a))
- integrate OAuth token lookup logic into UseCase ([cb3f911](https://github.com/specvital/collector/commit/cb3f9114f89b9435f037379b6eac07c51dc53d96))
- integrate scheduler for automatic codebase refresh ([e0a1a15](https://github.com/specvital/collector/commit/e0a1a15dac6ef26b62debf1d3b52d172cf8d8ed6))
- record failure status in DB when analysis fails ([6485ac3](https://github.com/specvital/collector/commit/6485ac345c9562404e435ea32962a1b90a13fd5b))
- support external analysis ID for record creation ([5448202](https://github.com/specvital/collector/commit/54482021c84b6faa466f0ce006de06ddcd79d22d))
- support OAuth token decryption for private repo analysis ([8d0ad30](https://github.com/specvital/collector/commit/8d0ad307ce07aa7562daa6a38230e19ce2cc1644))

#### 🐛 Bug Fixes

- handle missing error logging and DB status update on analysis task failure ([64ae8d9](https://github.com/specvital/collector/commit/64ae8d9fcbc620e7470e2c7cc21ade39d7327f8d))
- parser scan failing due to unexported method type assertion ([6256673](https://github.com/specvital/collector/commit/6256673dc23d48998851f977fcd034498c591642))
- remove unnecessary wait and potential deadlock in graceful shutdown ([b78c981](https://github.com/specvital/collector/commit/b78c981c662d32fde66c0890da0e226e9b4a4d3e))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- go mod tidy ([c58f73b](https://github.com/specvital/collector/commit/c58f73b40f2de49c8c69b1d67efd45b1487c0359))

#### 💄 Styles

- format code ([5e994e2](https://github.com/specvital/collector/commit/5e994e2ab90f6ae6a8cd64d392b946c9bde0bd1d))

#### ♻️ Refactoring

- centralize dependency wiring with DI container ([c1b8215](https://github.com/specvital/collector/commit/c1b82151bdba8b62e194100e6f04271fd3f4e026))
- extract domain layer with zero infrastructure dependencies ([7ba9e51](https://github.com/specvital/collector/commit/7ba9e51a0ffa76327736e91c828fe5949cbfbcb6))
- extract repository layer from AnalyzeHandler ([464ecfa](https://github.com/specvital/collector/commit/464ecfa6087d91d4d399a7fee032ed2a9109a151))
- extract service layer from AnalyzeHandler ([d9faf20](https://github.com/specvital/collector/commit/d9faf200da77b3097417cd1671a3a1c5fbc5fe06))
- implement handler layer and clean up legacy packages ([23a093f](https://github.com/specvital/collector/commit/23a093f6d5558d10decb21272289c1d99e583101))
- implement repository adapter layer (Clean Architecture Commit 3) ([8b0e433](https://github.com/specvital/collector/commit/8b0e43372fec1b15fa7e7d76794be37d42b6988e))
- implement use case layer with dependency injection ([b2be6ff](https://github.com/specvital/collector/commit/b2be6ff3d3a1c3662f816958d614f4a89215aba6))
- implement VCS and parser adapter layer (Clean Architecture Commit 4) ([1b2e34f](https://github.com/specvital/collector/commit/1b2e34f61665d6513c2654712390b67524dfd731))
- move infrastructure packages to internal/infra ([6cc1d1c](https://github.com/specvital/collector/commit/6cc1d1caf8722317967eae5c2de6c40c71467ce2))
- separate Scheduler from Worker into independent service ([9481141](https://github.com/specvital/collector/commit/9481141e99f0adcc225e93d05c8104846f836c17))
- split entry points for Railway separate deployments ([d899192](https://github.com/specvital/collector/commit/d899192cb1772fe9ed16d426d460e016c1bbf2ee))

#### ✅ Tests

- add AnalyzeHandler unit tests ([0286e7c](https://github.com/specvital/collector/commit/0286e7cd687d65be853e503a883cd74010f8dede))
- remove unnecessary skipped tests ([f8c0eb4](https://github.com/specvital/collector/commit/f8c0eb40a4f99252650bbf9e0f9ca93a378223fb))

#### 🔧 CI/CD

- configure semantic-release automated deployment pipeline ([37f128f](https://github.com/specvital/collector/commit/37f128f2c9d113144d2af530e88f84d3209f235c))

#### 🔨 Chore

- add bootstrap command ([c8371f0](https://github.com/specvital/collector/commit/c8371f0d5f47a19353c260ecf83c5033b4e5ba53))
- add Dockerfile for collector service ([6e3b0e4](https://github.com/specvital/collector/commit/6e3b0e4225b4b0875e2ee0bae909a594b1b9f87c))
- add example env file ([64a24a4](https://github.com/specvital/collector/commit/64a24a4de88aa5e4954a59b049b915c2012da79e))
- add gitignore item ([8fc64a6](https://github.com/specvital/collector/commit/8fc64a6ab0a3abca1cf1d73f458adf17bb752ced))
- add migrate local command ([baabcfe](https://github.com/specvital/collector/commit/baabcfe97f3905122026a752ea2ba7f7ed07917b))
- add PostgreSQL connection and sqlc configuration ([eecc4a6](https://github.com/specvital/collector/commit/eecc4a69a8b8c6e5c67a8f652a97ae784ecca1c1))
- add useful action buttons ([02fa778](https://github.com/specvital/collector/commit/02fa7785ac4ba1505e06ac3add60621cf01d1be9))
- adding recommended extensions ([30d5d0b](https://github.com/specvital/collector/commit/30d5d0b0fccc3190313456433e24e1342c18d641))
- ai-config-toolkit sync ([3091cf4](https://github.com/specvital/collector/commit/3091cf46ca2e6a24f5c299fe4f8008659fe1b8c8))
- ai-config-toolkit sync ([decf96b](https://github.com/specvital/collector/commit/decf96b2c47b278ef56a3e76c6174c9688f883c3))
- delete file ([f48005c](https://github.com/specvital/collector/commit/f48005cad322fa2586e9bf315e2bce3c608dcd8b))
- dump schema ([b90bab0](https://github.com/specvital/collector/commit/b90bab0f0c33d52683b6ff6a1f132702eb54a077))
- dump schema ([370409c](https://github.com/specvital/collector/commit/370409cee67512bc3f21ac3f5835357303db9b57))
- dump schema ([d704305](https://github.com/specvital/collector/commit/d7043054a0a59ac755bc23a85a4fd39f5ce97a0a))
- Global document synchronization ([cead255](https://github.com/specvital/collector/commit/cead255f25f48397848d55cf1417f21466dae67c))
- sync ai-config-toolkit ([e559889](https://github.com/specvital/collector/commit/e55988903526ade4630d2d6516e67ad1354ff67e))
- update core ([d358131](https://github.com/specvital/collector/commit/d358131e3e6197ee8958655b3cc1cfa7d0ed9ca6))
- update core ([b47592e](https://github.com/specvital/collector/commit/b47592e2a6668c25585d0338099e83e7b72bf1d5))
- update core ([395930a](https://github.com/specvital/collector/commit/395930a21bb48b8283cac037cde1999e44ae69c6))
- update schema.sql path in justfile ([0bcbe79](https://github.com/specvital/collector/commit/0bcbe794cdbccff58e2babe75a6308aacc6ad5d0))
- update-core version ([cc65b03](https://github.com/specvital/collector/commit/cc65b0325a1e828e24270753d76fa91ff01eeb45))
