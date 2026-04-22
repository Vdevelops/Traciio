---
trigger: always_on
---

# Architecture Patterns — CRM Healthcare Platform

Clean Architecture and layered design patterns applied to the CRM Healthcare monorepo, covering both Go backend and Flutter mobile.

## When to Use

- Designing new features end-to-end (backend + mobile)
- Reviewing or refactoring code for proper separation of concerns
- Deciding where to place new code within the project structure
- Ensuring dependency rules are followed

## Clean Architecture — Backend (Go/Gin)

### Layer Mapping

```
┌──────────────────────────────────────────┐
│  Handlers (apps/api/internal/api/handlers)│ ← Frameworks & Drivers
│  Routes   (apps/api/internal/api/routes)  │
├──────────────────────────────────────────┤
│  Service  (apps/api/internal/service)     │ ← Use Cases / Application
├──────────────────────────────────────────┤
│  Domain   (apps/api/internal/domain)      │ ← Entities / Business Rules
├──────────────────────────────────────────┤
│  Repository (apps/api/internal/repository)│ ← Interface Adapters (DB)
├──────────────────────────────────────────┤
│  Config   (apps/api/internal/config)      │ ← Infrastructure
│  Middleware (apps/api/internal/api/middleware)│
└──────────────────────────────────────────┘
```

### Dependency Rule (STRICT)

Dependencies flow **inward only**:
- `handlers/` → depends on `service/` (NEVER on `repository/` directly)
- `service/` → depends on `domain/` and `repository/` interfaces
- `repository/` → depends on `domain/` entities
- `domain/` → depends on NOTHING (pure Go structs)

### Backend Pattern Flow
```
HTTP Request
  → Route (authentication/RBAC middleware)
    → Handler (parse request, validate input)
      → Service (business logic, orchestration)
        → Repository (database operations via GORM)
          → PostgreSQL
```

### Implementation Rules
- **Thin handlers**: Parse input, call service, format response
- **Rich services**: All business logic lives here
- **Pure domain**: Only structs and business rules — no framework imports
- **Repository interfaces**: Define in service layer, implement in repository layer

## Clean Architecture — Mobile (Flutter)

### Layer Mapping (Per Feature)

```
┌──────────────────────────────────────────┐
│  Screens + Widgets (presentation/)        │ ← UI / Frameworks
├──────────────────────────────────────────┤
│  Providers + State (application/)         │ ← Use Cases / Application
├──────────────────────────────────────────┤
│  Models (data/models/)                    │ ← Entities
│  Repository (data/<feature>_repository)   │ ← Interface Adapters
│  Data Sources (data/datasources/)         │ ← External Communication
├──────────────────────────────────────────┤
│  Core (core/)                             │ ← Infrastructure (shared)
│    network/, storage/, config/, routing/  │
└──────────────────────────────────────────┘
```

### Mobile Dependency Rule

```
presentation/ → depends on application/ (providers)
application/  → depends on data/ (repositories, models)
data/         → depends on core/ (network, storage)
core/         → depends on NOTHING project-specific
```

### Mobile Pattern Flow
```
User Interaction
  → Screen (presentation/)
    → Provider triggers action (application/)
      → Repository (data/)
        → API Client (core/network/) or OfflineStorage (core/storage/)
```

## Feature-Based Organization

### Backend Feature Structure
```
internal/
  domain/<feature>/
    entity.go         → Request/response structs, domain types
  service/<feature>/
    service.go        → Business logic
  repository/
    <feature>_repository.go  → DB operations
  api/
    handlers/
      <feature>_handler.go   → HTTP handlers
    routes/
      <feature>_routes.go    → Route registration
```

### Mobile Feature Structure
```
lib/features/<feature>/
  data/
    models/           → Dart model classes (fromJson, toJson, copyWith)
    datasources/      → Remote/local data sources
    <feature>_repository.dart → Repository (API + offline)
  application/
    <feature>_provider.dart   → StateNotifierProvider
    <feature>_state.dart      → Immutable state class
  presentation/
    screens/          → Full page screens (*_screen.dart)
    widgets/          → Feature-specific widgets
```

## Key Patterns for This Project

### 1. Repository Pattern (Backend)
- Define repository interface in service layer
- Implement with GORM in `internal/repository/`
- Inject via constructor (dependency injection)
- GORM operations: `Create`, `Find`, `Where`, `Preload`, `Updates`, `Delete`

### 2. Repository Pattern (Mobile — Offline-First)
- Try cache first → network → update cache → fallback to stale cache
- Use `OfflineStorage` (Hive) for local persistence
- Use `ConnectivityService` to check network availability
- Always handle `DioException` gracefully

### 3. State Management Pattern (Mobile)
```dart
// Immutable state
class FeatureState {
  final bool isLoading;
  final String? error;
  final List<Model> items;
  FeatureState({ this.isLoading = false, this.error, this.items = const [] });
  FeatureState copyWith({...});
}

// StateNotifier orchestrates logic
class FeatureNotifier extends StateNotifier<FeatureState> {
  final FeatureRepository _repository;
  FeatureNotifier(this._repository) : super(FeatureState());
  // Methods trigger state transitions via state = state.copyWith(...)
}

// Provider wires it together
final featureProvider = StateNotifierProvider<FeatureNotifier, FeatureState>((ref) {
  return FeatureNotifier(ref.read(featureRepositoryProvider));
});
```

### 4. API Response Parsing (Mobile)
- Use flexible `fromJson` factories that handle both array and object responses
- Always check `response['success']` before parsing `response['data']`
- Extract meta/pagination from `response['meta']`

## Best Practices (Enforced in This Project)

1. **Single Responsibility**: Each file/class has ONE clear purpose
2. **Interface Segregation**: Small, focused interfaces over large ones
3. **Dependency Inversion**: Depend on abstractions, not implementations
4. **No circular dependencies**: Strict layer ordering
5. **Testability**: Each layer can be tested independently with mocks
6. **Consistency**: Mirror patterns from existing features (e.g., tasks, visit_reports, accounts)
