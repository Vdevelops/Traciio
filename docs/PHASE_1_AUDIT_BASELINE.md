# Phase 1 Audit Baseline
## CRM Healthcare / Pharmaceutical Platform

**Date**: 2026-04-23  
**Scope**: Backend, Web, Mobile, Docs, Database, Security  
**Method**: Strict evidence-based audit using runtime implementation first, documentation second

---

## Scope Locked

This baseline covers 14 modules:
- Auth
- Accounts
- Contacts
- Leads
- Visit Reports
- Pipeline & Deals
- Tasks
- Dashboard
- Reports
- Master Data
- Notifications
- Activities
- Route Optimization
- AI Module

Scoring contract:
- 8 layers per module: Domain, Repo, Service, Handler, FE Web, Mobile, Tests, Docs
- TOTAL = average of 8 layer scores
- Evidence priority: runtime code > route wiring > docs
- Modules below 20% should be treated as missing/near-missing

---

## A. Module Scorecard Matrix

| Module | Domain | Repo | Service | Handler | FE Web | Mobile | Tests | Docs | TOTAL |
|--------|--------|------|---------|---------|--------|--------|-------|------|-------|
| Auth | 95% | 90% | 90% | 88% | 88% | 88% | 65% | 85% | 86% |
| Accounts | 90% | 88% | 85% | 75% | 75% | 85% | 40% | 85% | 78% |
| Contacts | 90% | 88% | 85% | 75% | 75% | 85% | 40% | 85% | 78% |
| Leads | 92% | 90% | 82% | 85% | 78% | 82% | 45% | 88% | 80% |
| Visit Reports | 92% | 90% | 86% | 88% | 82% | 85% | 55% | 88% | 83% |
| Pipeline & Deals | 92% | 90% | 80% | 80% | 80% | 78% | 45% | 84% | 79% |
| Tasks | 88% | 88% | 78% | 82% | 74% | 70% | 35% | 80% | 74% |
| Dashboard | 70% | 60% | 85% | 82% | 88% | 82% | 35% | 80% | 73% |
| Reports | 70% | 40% | 78% | 75% | 72% | 0% | 30% | 82% | 56% |
| Master Data | 92% | 92% | 88% | 86% | 84% | 0% | 70% | 78% | 74% |
| Notifications | 90% | 90% | 82% | 80% | 82% | 65% | 70% | 65% | 78% |
| Activities | 88% | 85% | 78% | 80% | 68% | 20% | 35% | 70% | 66% |
| Route Optimization | 88% | 88% | 88% | 82% | 80% | 85% | 30% | 55% | 75% |
| AI Module | 80% | 50% | 80% | 78% | 78% | 0% | 35% | 40% | 55% |

Interpretation:
- Strongest modules: Auth, Visit Reports, Leads
- At risk but usable: Accounts, Contacts, Pipeline & Deals, Tasks, Dashboard, Notifications, Route Optimization, Master Data, Activities
- Major gaps: Reports, AI Module, mobile parity for Master Data and Activities

---

## B. Critical Gap List

Production-blocking issues identified in Phase 1:

1. [Pipeline & Deals] Stage move via `PATCH /deals/:id/move` does not consistently write deal history, while the alternate `POST /deals/:id/move-stage` path does. Business impact: inconsistent audit trail and pipeline governance.
2. [Tasks] Reminder scheduling is manual. The worker only processes reminders that already exist; D-1 and day-H reminders are not auto-created. Business impact: follow-up reminders can be missed.
3. [Notifications] Approval and stage-change notification triggers are not fully wired end-to-end. Business impact: supervisors and sales users do not receive some critical alerts in real time.
4. [Security] Several route groups are auth-only without explicit scope or permission enforcement. Business impact: data scoping and authorization are weaker than expected for enterprise CRM.
5. [Mobile] Reports, Master Data, and AI are missing as dedicated mobile feature modules. Business impact: mobile users cannot access these modules natively.

---

## C. Missing Module Map

No module scored below 20% overall, but the following are effectively missing on mobile or only partially represented:

- Reports: mobile 0%
- Master Data: mobile 0%
- AI Module: mobile 0%
- Activities: mobile 20% and mostly embedded through other flows

Backend presence exists for all 14 modules, but mobile parity is incomplete for the items above.

---

## D. Database Schema Issues

Key schema observations:

1. `lead_qualification_checklist` has an entity and service layer, but no explicit SQL migration file was found in `apps/api/internal/database/migrations`.
2. Visit report approval fields exist in the entity (`approved_by`, `approved_at`, `rejection_reason`) but were not found as explicit SQL definitions in the migration set.
3. Route optimization entity exists, but no explicit CREATE TABLE migration was found for `optimized_routes` in the SQL migration folder.
4. The database strategy is hybrid: GORM `AutoMigrate` creates many core tables at startup, while the SQL migration folder mostly contains ALTER/INDEX updates. This is acceptable, but it means schema auditing must distinguish auto-migrated tables from SQL-owned changes.
5. Soft delete support is broadly consistent for transactional modules, but not every entity uses `DeletedAt`.

Notable confirmations:
- Lead BANT columns and deal history table are present in migration SQL.
- Index coverage is strong for leads, deals, activities, reminders, visit reports, notifications, and pipeline stages.

---

## E. Security Gaps

Route-level findings:

| Route Group | Auth | Scope | Permission | Notes |
|---|---|---|---|---|
| /api/v1/accounts | Yes | No | No | Auth-only; no scope middleware |
| /api/v1/contacts | Yes | No | No | Auth-only; no scope middleware |
| /api/v1/reports | Yes | No | No | Auth-only; no scope middleware |
| /api/v1/activities | Yes | No | No | Auth-only; no scope middleware |
| /api/v1/route-optimization | Yes | No | No | Auth-only; no scope middleware |
| /api/v1/deals/:id/move | Yes | Partial | No explicit HasPermission | Critical stage-change endpoint |
| /api/v1/ws/notifications | Handler-level | N/A | No | Token validated in handler, not route middleware |
| /api/v1/notifications | Yes | No | No | Query parsing is manual in handler |

Handler-level permission checks found explicitly:
- `leads.convert`
- `pipeline.update_stage`
- `tasks.create`
- `tasks.create_lead`

Input validation patterns are mostly present via `ShouldBindJSON` and `ShouldBindQuery`, but notifications still use manual parsing for some fields.

---

## F. 80% Achievement Plan

Modules currently below 80%:

| Module | Current | Effort to 80% | Main Workstream | Key Dependency |
|---|---:|---:|---|---|
| Accounts | 78% | 3-4 days | Security + tests | Scope enforcement policy |
| Contacts | 78% | 3-4 days | Security + tests | Scope enforcement policy |
| Pipeline & Deals | 79% | 4-6 days | History consistency + tests | Single stage-move path |
| Tasks | 74% | 6-8 days | Reminder automation + mobile polish | Auto reminder creation |
| Dashboard | 73% | 4-6 days | Test coverage | Consistent scoped inputs |
| Reports | 56% | 8-12 days | Mobile parity + security | Read-only report UX |
| Master Data | 74% | 7-10 days | Mobile feature + route alignment | Final backend route contract |
| Notifications | 78% | 4-5 days | Trigger wiring | Approval/stage/reminder events |
| Activities | 66% | 7-10 days | Event-driven activity creation + mobile | Event consumer wiring |
| Route Optimization | 75% | 4-6 days | Tests + security | Scope enforcement |
| AI Module | 55% | 10-14 days | Mobile module + quotas + tests | AI usage controls |

Modules that are already at or above 80% and should be protected rather than reworked:
- Auth
- Leads
- Visit Reports

---

## Evidence Highlights

### Business Rules
- Lead conversion is restricted to qualified leads in [apps/api/internal/service/lead/service.go](../apps/api/internal/service/lead/service.go)
- Visit reports enforce draft/submitted/approved/rejected state transitions in [apps/api/internal/service/visit_report/service.go](../apps/api/internal/service/visit_report/service.go)
- Deal history creation exists in [apps/api/internal/service/pipeline/service.go](../apps/api/internal/service/pipeline/service.go)
- Reminder processing exists in [apps/api/internal/worker/reminder_worker.go](../apps/api/internal/worker/reminder_worker.go)
- Notification creation exists in [apps/api/internal/service/notification/service.go](../apps/api/internal/service/notification/service.go)

### Security
- Auth middleware and scope middleware are wired broadly in [apps/api/internal/api/routes](../apps/api/internal/api/routes)
- Explicit handler-level permission checks exist in lead, pipeline, and task handlers
- WebSocket notifications authenticate inside the handler in [apps/api/internal/api/handlers/websocket_handler.go](../apps/api/internal/api/handlers/websocket_handler.go)

### Platform Parity
- Full web feature coverage is present under [apps/web/src/features](../apps/web/src/features)
- Mobile has no dedicated feature folders for Reports, Master Data, or AI under [apps/mobile/lib/features](../apps/mobile/lib/features)

---

## Baseline Conclusion

Phase 1 establishes that the backend is broadly complete, the web frontend is mostly complete, and the mobile app has strong coverage for core sales flows but still lacks three major modules. The main implementation risks are not missing CRUD surfaces, but consistency, scoping, reminder automation, notification triggers, and mobile parity for reporting and AI.

This document is intended to be the baseline contract for Phase 2 through Phase 5.
