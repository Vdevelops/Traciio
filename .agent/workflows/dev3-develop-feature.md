---
description: Develop Flutter Feature (Dev3 Fullstack Workflow)
---

# Workflow: Develop Flutter Feature

End-to-end workflow for Dev3 to implement a new feature — backend API first, then mobile Flutter frontend.

## Phase 1: Analyze & Plan

1. **Read requirement** from `docs/sprint/sprint1/SPRINT_PLANNING_DEV3.md` or user request.
2. **Check existing patterns** — look at similar completed features (e.g., `dashboard/`, `tasks/`, `visit_reports/`) for structure reference.
3. **Identify API needs** — determine if mobile-specific endpoints are needed (`/api/v1/dashboard/mobile/` or `/api/v1/mobile/<feature>/`) or if web endpoints (`/api/v1/<resource>`) can be reused.
4. **Plan file structure**:
   - Backend: `internal/domain/<feature>/entity.go`, `internal/service/<feature>/service.go`, `internal/api/handlers/<feature>_handler.go`, `internal/api/routes/<feature>_routes.go`
   - Mobile: `features/<feature>/data/models/`, `features/<feature>/data/<feature>_repository.dart`, `features/<feature>/application/`, `features/<feature>/presentation/`

## Phase 2: Backend Implementation (Go/Gin)

1. **Domain entities** — define request/response structs in `internal/domain/<feature>/entity.go`.
2. **Repository** — implement database operations in `internal/repository/postgres/<feature>/`.
3. **Service** — business logic in `internal/service/<feature>/service.go`.
4. **Handler** — HTTP handler methods in `internal/api/handlers/<feature>_handler.go`.
5. **Routes** — register endpoints in `internal/api/routes/<feature>_routes.go` with auth middleware.
6. **Response format** — follow `{ success, data, meta, timestamp, request_id }` standard from `docs/api-standart/api-response-standards.md`.
7. **Error codes** — use bilingual error codes from `docs/api-standart/api-error-codes.md`.

## Phase 3: Documentation & Testing

1. **Update Postman collection** in `docs/postman/` with new endpoints.
2. **Provide curl examples** for manual testing.
3. **Suggest testing** — present curl commands and expected responses to user.
4. **Wait for user "OK"** before proceeding to mobile frontend.

## Phase 4: Mobile Frontend (Flutter)

### 4a. Data Layer
1. **Models** — create Dart model classes in `features/<feature>/data/models/<feature>.dart`:
   - `fromJson` factory constructor (null-safe, with defaults).
   - `toJson` method for cache serialization.
   - Parse flexibly to handle both web and mobile API formats.
2. **Repository** — create `features/<feature>/data/<feature>_repository.dart`:
   - Implement **offline-first pattern**: cache → network → fallback.
   - Use `OfflineStorage` (Hive) for persistent cache.
   - Use `ConnectivityService` to check online status.
   - Background refresh when online and using cached data.
   - `_extractData(response)` for API response unwrapping.
3. **Cache** (optional) — create `features/<feature>/data/<feature>_cache.dart` for in-memory TTL cache.

### 4b. Application Layer
1. **State class** — create `features/<feature>/application/<feature>_state.dart`:
   - Immutable with `copyWith` method.
   - Fields: `isLoading`, `errorMessage`, data fields, filters.
2. **Provider** — create `features/<feature>/application/<feature>_provider.dart`:
   - `StateNotifierProvider<FeatureNotifier, FeatureState>`.
   - Repository provider via `Provider<FeatureRepository>`.
   - Load, refresh, filter methods.
   - Centralized error extraction from DioException.

### 4c. Presentation Layer
1. **Screen** — create `features/<feature>/presentation/<feature>_screen.dart`:
   - Extend `ConsumerStatefulWidget`.
   - Call provider's load method in `initState` via `addPostFrameCallback`.
   - Use `RefreshIndicator` for pull-to-refresh.
   - Handle loading/error/empty/data states.
2. **Widgets** — create feature-specific widgets in `features/<feature>/presentation/widgets/`:
   - Card widgets, section widgets, form widgets.
   - Extract reusable widgets into separate files.
3. **Routing** — register new routes in `core/routing/app_router.dart`.

## Phase 5: Testing & Quality

1. Run `flutter analyze` — fix all warnings.
2. Run `flutter test` — ensure no regressions.
3. Write tests:
   - **Model tests**: verify `fromJson`/`toJson` with sample API responses.
   - **Provider tests**: verify state transitions (loading → data, loading → error).
4. Manual testing on emulator/device.

## Phase 6: Documentation & Sprint Update

1. **Feature docs**: Create `docs/features/mobile/<feature-name>.md` with 11 standard sections:
   - Overview, Architecture, API Endpoints, Models, Repository, Provider, Screens, Widgets, Testing, Offline Support, Known Issues.
2. **Sprint update**: Mark tasks as `[x]` in:
   - `docs/sprint/sprint1/SPRINT_PLANNING_DEV3.md`
   - `docs/sprint/SPRINT_PLANNING.md`
3. **Copilot instructions**: Update `.github/copilot-instructions.md` if sprint status changed.

## Checklist Template

```markdown
### Feature: <Feature Name>

#### Backend
- [ ] Domain entity (`internal/domain/<feature>/entity.go`)
- [ ] Repository (`internal/repository/postgres/<feature>/`)
- [ ] Service (`internal/service/<feature>/service.go`)
- [ ] Handler (`internal/api/handlers/<feature>_handler.go`)
- [ ] Routes (`internal/api/routes/<feature>_routes.go`)
- [ ] Postman collection updated

#### Mobile
- [ ] Models (`features/<feature>/data/models/`)
- [ ] Repository with offline-first (`features/<feature>/data/<feature>_repository.dart`)
- [ ] State class (`features/<feature>/application/<feature>_state.dart`)
- [ ] Provider (`features/<feature>/application/<feature>_provider.dart`)
- [ ] Screen (`features/<feature>/presentation/<feature>_screen.dart`)
- [ ] Widgets (`features/<feature>/presentation/widgets/`)
- [ ] Route registered in `app_router.dart`

#### Quality
- [ ] `flutter analyze` passes
- [ ] `flutter test` passes
- [ ] Model unit tests
- [ ] Provider unit tests
- [ ] Manual testing on emulator

#### Docs
- [ ] Feature doc in `docs/features/mobile/`
- [ ] Sprint planning updated
```

## Notes

- **Backend first** — always implement and test backend before mobile.
- **Mirror patterns** — look at `dashboard/`, `tasks/`, or `visit_reports/` as reference implementations.
- **Riverpod only** — do not introduce other state management libraries.
- **Offline-first** — every repository should handle offline gracefully.
- **No separate planning docs** — use existing sprint files only.
