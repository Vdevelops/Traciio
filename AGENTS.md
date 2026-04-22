# AGENTS.md — CRM Healthcare Monorepo

> Guide for AI agents working in this repository. Last updated: 2025-03-04

## Repository Overview

Healthcare/pharmaceutical CRM platform monorepo using **Turborepo + pnpm** with three main applications:

- **Web** (`apps/web/`): Next.js 16, React 19, TypeScript, Tailwind CSS v4
- **API** (`apps/api/`): Go 1.25+, Gin, GORM, PostgreSQL
- **Mobile** (`apps/mobile/`): Flutter, Dart 3+, Riverpod, Dio, Hive

## Build Commands

### Root Level (Turborepo)

use 'npx' to run 'pnpm' command

```bash
# Development - run all apps in parallel
pnpm dev                    # Start web, api dev servers
pnpm dev:web               # Web only
pnpm dev:api               # API with Docker PostgreSQL

# Build
pnpm build                 # Build all apps
pnpm build --filter=web    # Build specific app

# Lint & Type Check
pnpm lint                  # ESLint for all
pnpm type-check           # TypeScript check for all
pnpm format               # Prettier format

# Test
pnpm test                  # Run all tests

# Clean
pnpm clean                # Clean build artifacts
```

### Web App (`apps/web/`)

```bash
cd apps/web

# Dev
pnpm dev                   # Next.js dev server

# Build
pnpm build                 # Production build
pnpm start                 # Start production server

# Quality
pnpm lint                  # ESLint (zero warnings)
pnpm check-types          # tsc --noEmit
```

### API App (`apps/api/`)

```bash
cd apps/api

# Dev
pnpm dev                   # go run ./cmd/server/main.go
pnpm dev:account          # With DROP_TABLES=true ONLY_ACCOUNT=true

# Build
pnpm build                 # go build -o bin/server
pnpm start                 # Run built binary

# Quality
pnpm lint                  # golangci-lint run
pnpm test                  # go test ./...

# Docker
make docker-up            # docker-compose up --build
make docker-down          # docker-compose down

# Direct Go commands
make run                  # go run cmd/server/main.go
make test                 # go test ./...
make build                # go build -o server
```

### Mobile App (`apps/mobile/`)

```bash
cd apps/mobile

# Dev
flutter run                                        # Default
flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080     # Android emulator
flutter run --dart-define=API_BASE_URL=http://localhost:8080    # iOS
flutter run --dart-define=API_BASE_URL=http://192.168.1.x:8080  # Physical device

# Build
flutter build apk --dart-define=API_BASE_URL=...   # Android APK
flutter build ios --dart-define=API_BASE_URL=...   # iOS

# Quality
flutter analyze            # Dart/Flutter lint
flutter test               # Run tests
flutter clean              # Clean build artifacts
```

## Testing Single Files

### Web (Jest/Vitest pattern)

```bash
# Run single test file
pnpm test -- path/to/file.test.ts
pnpm test -- --testNamePattern="test name"

# Watch mode
pnpm test -- --watch
```

### API (Go)

```bash
# Run single test
pnpm test ./internal/service/account

# Run with verbose
pnpm test ./internal/... -v

# Run specific test function
pnpm test -run TestCreateAccount ./internal/service/account
```

### Mobile (Flutter)

```bash
# Run single test file
flutter test test/unit/feature_test.dart

# Run specific test
flutter test --name "test description"

# Run tests matching pattern
flutter test --plain-name "Dashboard"
```

## Code Style Guidelines

### TypeScript/Web

**Types:**

- Never use `any` — use `unknown` with type guards
- Use `.d.ts` for ambient types in `types/` folders
- Export DTOs from `types/` and import everywhere

**Imports:**

```typescript
// External packages first
import React from "react";
import { z } from "zod";

// Internal absolute imports
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";

// Relative imports last
import { localHelper } from "./utils";
```

**Naming:**

- Directories: kebab-case (`user-profile/`)
- Components: PascalCase (`UserProfile.tsx`)
- Hooks: camelCase with `use` prefix (`useUserData.ts`)
- Services: camelCase (`userService.ts`)
- Stores: PascalCase with `use` prefix (`useUserStore.ts`)

### Go/API

**Structure:**

```
internal/
  domain/<feature>/     # Entities/structs
  repository/           # DB operations (interface + impl)
  service/<feature>/    # Business logic
  api/
    handlers/           # HTTP handlers
    routes/             # Route registration
    middleware/         # Auth, CORS, logging
```

**Naming:**

- Files: snake_case (`user_handler.go`)
- Types: PascalCase (`UserHandler`)
- Interfaces: `Repository` suffix (`UserRepository`)
- Functions: PascalCase for exported, camelCase for private

**Error Handling:**

- Always handle errors explicitly (no naked returns)
- Propagate with `fmt.Errorf("%w", err)`
- Use structured error codes from `docs/api-standart/api-error-codes.md`
- Bilingual messages (ID/EN)

**API Response Format:**

```json
{
  "success": true,
  "data": {},
  "meta": {},
  "timestamp": "2025-01-01T00:00:00Z",
  "request_id": "uuid"
}
```

### Flutter/Mobile

**Structure:**

```
lib/
  core/                 # Shared infrastructure
    config/             # Env, base URL
    routing/            # AppRouter, named routes
    theme/              # Light/dark themes
    network/            # Dio client, interceptors
    storage/            # Hive, SharedPreferences
    widgets/            # Reusable widgets
  features/<feature>/   # Per domain
    data/               # Models, repositories
    application/        # Riverpod providers
    presentation/       # Screens, widgets
```

**Naming:**

- Files: snake_case (`dashboard_screen.dart`)
- Classes: PascalCase (`DashboardScreen`)
- Providers: camelCase + suffix (`dashboardProvider`)
- Variables/methods: camelCase (`fetchData()`)

**Architecture:**

- Use **Riverpod** (`StateNotifierProvider`) — no Provider/BLoC
- Offline-first: Cache → Network → Fallback
- Never hardcode URLs — use `Env.apiBaseUrl`
- All HTTP via `ApiClient.dio` in `core/network/`

## Key Rules

1. **No business logic in UI components** — extract to hooks/providers
2. **Always handle loading/error/empty states** in UI
3. **Use optional chaining (`?.`)** with nullish coalescing (`??`)
4. **All clickable elements must have `cursor-pointer`** class
5. **Mobile endpoints**: `/api/v1/mobile/` or `/api/v1/dashboard/mobile/`
6. **Update Postman collection** when modifying APIs: `docs/postman/`
7. **Update sprint docs** when completing tasks: `docs/sprint/`

## Quick References

- API Standards: `docs/api-standart/api-response-standards.md`
- Error Codes: `docs/api-standart/api-error-codes.md`
- PRD: `docs/PRD.md`
- Sprint Planning: `docs/sprint/SPRINT_PLANNING.md`

## Agent Rules (Extended)

- Cursor Rules: `.cursor/rules/`
- Copilot Instructions: `.github/copilot-instructions.md`
- Agent Rules: `.agent/rules/`
