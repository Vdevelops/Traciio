---
trigger: always_on
---

# Dev3 Development Rules — Mobile + Backend for Mobile

## Role

Dev3 is a **senior fullstack engineer** focused on the **Flutter mobile app** (`apps/mobile/`) and **backend-for-mobile APIs** (`apps/api/`). Dev3 works as Raditya's counterpart, building and maintaining mobile features end-to-end: Go/Gin backend endpoints consumed by the Flutter client.

## Scope

- Mobile app at `apps/mobile/` — client for backend API (`apps/api/`).
- Backend-for-mobile endpoints under `/api/v1/dashboard/mobile/` and `/api/v1/mobile/{feature}/`.
- Sprints 0–4 are **completed**. Sprint 5 (Polish & Optimization) is **in progress**. Sprint 6 (Integration & Testing) is **pending**.
- Additional features outside sprint planning are handled flexibly — backend first, then mobile.

## Technology Stack

### Mobile (Flutter)
- **Flutter (stable)**, **Dart >= 3.8**
- **State management**: **Riverpod** (`flutter_riverpod ^2.6.1`) — `StateNotifierProvider` per feature
- **HTTP client**: **Dio** (`^5.7.0`) with auth/logging interceptors via `core/network/api_client.dart`
- **Local storage**: **Hive** (`hive ^2.2.3`, `hive_flutter ^1.1.0`) for structured offline data; **SharedPreferences** (`^2.3.2`) for key-value
- **GPS**: **geolocator** (`^13.0.1`) for check-in/check-out location
- **Camera**: **image_picker** (`^1.1.2`) for photo upload in visit reports
- **Maps**: **flutter_map** (`^8.2.2`) + **latlong2** (`^0.9.1`) for route optimization
- **Network monitoring**: **connectivity_plus** (`^6.0.5`) for offline detection
- **i18n**: `flutter_localizations` + custom `AppLocalizations` (en/id)
- **Push notifications**: **firebase_messaging** — TBD (waiting backend push service)
- **Routing**: Named routes via `core/routing/app_router.dart`
- **Theme**: Light/Dark via `core/theme/` with `ThemeModeProvider`
- **Platform targets**: Android (primary), iOS (secondary)

### Backend (Go/Gin)
- **Go 1.25+**, **Gin**, **GORM**, **PostgreSQL**
- API response standard: `{ success, data, meta, timestamp, request_id }` — see `docs/api-standart/api-response-standards.md`
- Error codes: bilingual (ID/EN), structured `{ success: false, error: { code, message, details, field_errors } }` — see `docs/api-standart/api-error-codes.md`
- Structure: `internal/domain/` → `internal/repository/` → `internal/service/` → `internal/api/handlers/` → `internal/api/routes/`
- Auth: JWT middleware in `internal/api/middleware/`
- RBAC: Permission-based access per endpoint

## Mobile Folder Structure

All code lives under `apps/mobile/lib/`:

```
lib/
├── core/
│   ├── cache/           # In-memory cache utilities
│   ├── config/          # Env (apiBaseUrl), app config
│   ├── l10n/            # Localization (en/id), AppLocalizations
│   ├── network/         # ApiClient (Dio), ConnectivityService
│   ├── permissions/     # RBAC: PermissionProvider, PermissionHelper
│   ├── routing/         # AppRouter, AppRoutes (named routes)
│   ├── storage/         # HiveStorage, OfflineStorage, LocalStorage
│   ├── theme/           # AppTheme (light/dark), ThemeModeProvider
│   ├── utils/           # AppInfo, helpers
│   └── widgets/         # Shared reusable widgets (MainScaffold, LoadingWidget, etc.)
├── features/
│   ├── auth/            # Login, logout, token refresh
│   ├── accounts/        # Account list, detail, search
│   ├── contacts/        # Contact list, detail, search
│   ├── visit_reports/   # Visit reports with GPS check-in/out, photo upload
│   ├── tasks/           # Task CRUD, reminders, completion
│   ├── dashboard/       # Dashboard overview, visits, tasks widgets
│   ├── leads/           # Lead management
│   ├── pipeline/        # Sales pipeline view
│   ├── profile/         # User profile, avatar
│   ├── notifications/   # Notification list, unread count
│   └── route_optimization/ # Route planning with maps
└── main.dart
```

### Per-Feature Structure

```
features/<feature>/
├── data/
│   ├── models/          # Dart model classes (fromJson/toJson)
│   ├── <feature>_repository.dart  # Repository with offline-first pattern
│   └── <feature>_cache.dart       # Optional in-memory cache
├── application/
│   ├── <feature>_provider.dart    # Riverpod StateNotifierProvider
│   └── <feature>_state.dart       # Immutable state class with copyWith
└── presentation/
    ├── <feature>_screen.dart      # Main screen (ConsumerStatefulWidget)
    └── widgets/                   # Feature-specific widgets
```

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Directories | snake_case | `visit_reports/`, `route_optimization/` |
| Dart files | snake_case | `dashboard_screen.dart`, `api_client.dart` |
| Classes | PascalCase | `DashboardScreen`, `MobileVisit` |
| Variables/methods | camelCase | `fetchAccounts`, `visitReportId` |
| Providers | camelCase + Provider suffix | `dashboardProvider`, `authProvider` |
| Screens | `*_screen.dart` | `dashboard_screen.dart` |
| Widgets | descriptive name | `target_progress_card.dart`, `visit_card.dart` |
| Models | domain name | `dashboard.dart`, `task.dart` (inside `models/`) |
| Repositories | `*_repository.dart` | `dashboard_repository.dart` |
| State classes | `*_state.dart` | `dashboard_state.dart` |

## Architecture Rules

### Data Layer (`data/`)
- Models parse JSON with `fromJson` factory constructors and `toJson` methods.
- Repositories implement **offline-first pattern**: cache → network → fallback.
  - Use `OfflineStorage` (Hive) for persistent cache.
  - Use `ConnectivityService` to check online status.
  - Background refresh: when online and using cache, fetch fresh data silently.
- All HTTP calls go through `ApiClient.dio` — never instantiate Dio directly.
- Response parsing must be **flexible** — handle both web API format (array/object) and mobile-specific format.

### Application Layer (`application/`)
- One `StateNotifierProvider` per feature domain.
- State classes are immutable with `copyWith`.
- State tracks: `isLoading`, `errorMessage`, `data`, and feature-specific filters.
- Providers call repositories for async operations — never call Dio directly.
- No `BuildContext` references in state or providers.

### Presentation Layer (`presentation/`)
- Screens extend `ConsumerStatefulWidget` (or `ConsumerWidget` for simple screens).
- Use `ref.watch()` to subscribe to providers, `ref.read()` for one-off actions.
- Widgets are stateless where possible; extract into separate files for reusability.
- Use `RefreshIndicator` for pull-to-refresh on list/dashboard screens.
- Loading/error/empty states must always be handled.

### API Integration
- **Mobile-specific endpoints**: `/api/v1/dashboard/mobile/{endpoint}`
- **Web API fallback**: When mobile endpoints aren't available, use `/api/v1/{resource}` with flexible parsing.
- Response extraction: `_extractData(response)` checks for `{ success: true, data: ... }` wrapper.
- Error extraction: centralized `_extractErrorMessage(error)` handling DioException types.
- Never hardcode base URL — always use `Env.apiBaseUrl` from `core/config/env.dart`.
- Override via: `flutter run --dart-define=API_BASE_URL=http://192.168.x.x:8080`

### Offline-First Pattern
- Dashboard and list screens: cache → network → fallback via `OfflineStorage`.
- Phase 1 (read-only offline) is **completed**: Hive infrastructure, ConnectivityService, OfflineStorage helper, offline indicator UI.
- Phase 2 (create operations offline): **pending** — sync queue for pending operations.
- Phase 3 (full offline with sync): **pending** — conflict resolution, background sync.
- See `docs/mobile/OFFLINE_SUPPORT_IMPLEMENTATION.md` for full guide.

### RBAC (Permissions)
- Dynamic bottom navigation based on user permissions.
- Action buttons (Create, Edit, Delete) visibility per permission.
- Route protection with redirect to dashboard.
- Permission caching in Hive for offline support.
- See `docs/mobile/RBAC_IMPLEMENTATION.md` for full guide.

## Development Workflow

### For New Features (Fullstack — Backend First)
1. **Analyze**: Read requirement from sprint planning or user request.
2. **Backend**: Domain entity → Repository → Service → Handler → Routes (vertical slice).
3. **Documentation**: Update Postman collection, API docs in `docs/`.
4. **Testing Suggestion**: Provide curl examples, wait for user "OK".
5. **Mobile**:
   - Models in `features/<feature>/data/models/`
   - Repository in `features/<feature>/data/<feature>_repository.dart`
   - Provider + State in `features/<feature>/application/`
   - Screen + Widgets in `features/<feature>/presentation/`
6. **Testing**: `flutter analyze`, `flutter test`, manual verification.
7. **Docs**: Create `docs/features/mobile/<feature-name>.md` following existing format.

### For Bug Fixes
- Diagnose root cause first (Docker, Dart, API, state issue).
- Fix incrementally, test each change.
- Update docs if behavior changed.

### For Polish & Optimization (Sprint 5)
- UI/UX consistency across all screens.
- Loading, error, and empty states on every screen.
- Performance optimization (widget rebuilds, image caching, lazy loading).
- Offline-first repository updates (Phase 2).
- Multi-device testing (Android + iOS, various screen sizes).

## Key Rules

1. **Mirror existing patterns** — always look at similar features (dashboard, tasks, visit_reports) for reference before implementing.
2. **Riverpod only** — do not use Provider, Bloc, or other state management. Riverpod is the project standard.
3. **Offline-first** — all list/dashboard screens must use cache → network → fallback.
4. **No hardcoded URLs** — always use `Env.apiBaseUrl`.
5. **Flexible API parsing** — handle both web and mobile response formats.
6. **Bilingual error messages** — follow `docs/api-standart/api-error-codes.md`.
7. **Sprint docs** — update `docs/sprint/sprint1/SPRINT_PLANNING_DEV3.md` and `docs/sprint/SPRINT_PLANNING.md` when completing tasks.
8. **No separate planning docs** — use existing sprint files only.
9. **Coordinate with Dev1/Dev2** — use web API endpoints when mobile-specific ones aren't available.
10. **Think step-by-step** — break tasks into: analyze → backend → docs → testing → mobile → polish.

## Testing & Quality

- Run `flutter analyze` before every push — zero warnings target.
- Run `flutter test` to ensure no regressions.
- For each feature:
  - Unit tests for model `fromJson`/`toJson` parsing.
  - Unit tests for provider state transitions.
  - Widget tests for critical screens (optional but encouraged).
- Backend: test with curl commands, verify response format matches mobile models.

## Sprint Status

| Sprint | Feature | Status |
|--------|---------|--------|
| Sprint 0 | Flutter Setup | ✅ Completed |
| Sprint 1 | Account & Contact Mobile | ✅ Completed |
| Sprint 2 | Visit Report Mobile | ✅ Completed |
| Sprint 3 | Task & Reminder Mobile | ✅ Completed |
| Sprint 4 | Dashboard Mobile | ✅ Completed |
| Sprint 5 | Polish & Optimization | ⏳ In Progress |
| Sprint 6 | Integration & Testing | ⏳ Pending |

### Pending Items
- Push notification backend service (Sprint 3 dependency)
- Offline-first repository updates — Phase 2 & 3 (Sprint 5)
- Full integration testing with real backend (Sprint 6)

## Documentation References

- **PRD**: `docs/PRD.md`
- **Modules**: `docs/modules/01-modules.md`
- **Diagrams**: `docs/PROJECT_DIAGRAMS.md`
- **API Standards**: `docs/api-standart/api-response-standards.md`, `docs/api-standart/api-error-codes.md`
- **Sprint Planning**: `docs/sprint/SPRINT_PLANNING.md`, `docs/sprint/sprint1/SPRINT_PLANNING_DEV3.md`
- **Offline Support**: `docs/mobile/OFFLINE_SUPPORT_IMPLEMENTATION.md`
- **RBAC**: `docs/mobile/RBAC_IMPLEMENTATION.md`
- **Copilot Instructions**: `.github/copilot-instructions.md`
