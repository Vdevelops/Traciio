# Working Plan: ERP Integration for Leads, Pipeline & Tasks
*(Based on CRM Healthcare ERP-API and ERP-Front-End)*

This document outlines the business flow, rules, logic, and RBAC derived from the `erp-api` and `erp-front-end` repositories, and the step-by-step working plan to implement these features in the current CRM system codebase (`apps/api` and `apps/web`).

---

## 1. Business Logic, Flow, and Rules

### 1.1 Lead Management
**Logic & Rules:**
*   **Statuses:** `NEW`, `CONTACTED`, `QUALIFIED`, `CONVERTED`, `LOST`.
*   **BANT Score (0-100):** Budget, Authority, Need, Pipeline are calculated individually (20 points each + 20 points bonus if all are true). This contributes 20% to the overall Lead Score.
*   **Lead Score Calculation (0-100):** Combines Source Base Score (e.g. Referral = 30), Estimated Value bracket score, Status maturity score, Activities count score, and the BANT score contribution.
*   **Data Enrichment:** Inside a Lead, users can directly add `Visit Reports`, `Activities`, and user `Tasks` via explicit tabs.
*   **BANT Checklist:** A dedicated interface or checklist representing the BANT metrics (e.g. Budget target amount, decision maker tags, etc).
*   **Conversion:** `Lead` is converted to a `Pipeline` (Opportunity).
    *   **Relationship Continuity:** When converted, all historical data (Visit Reports, Activities, Notes) *must remain linked or implicitly carried over* to the newly created Pipeline. 
    *   A new Activity/Task of type `NOTE` is automatically created to log the conversion.

### 1.2 Pipeline / Opportunity Management
**Logic & Rules:**
*   **Pipeline Stages & Progression (Strict Forward Progression):**
    *   `NEW_LEAD` / Awareness (0%)
    *   `INITIAL_CONTACT` / Interest (25%)
    *   `PROPOSAL` / Desire (50%)
    *   `FINAL_NEGOTIATION` (75%)
    *   `DEAL_WON` / Order (100%)
    *   `CLOSED_LOST` (0% - Can happen at any point).
*   **Progression Constraints:** Opportunities cannot move backward (e.g., from 50% back to 25%).
*   **Rich Collaboration:** Just like Leads, users can view/add `Visit Reports`, `Activities`, and assigned `Tasks` directly to the Pipeline under specific tabs.
*   **Products Dependency:** Unlike Leads, Pipeline deals have an added `Product` catalog component. `EstimatedValue` is dynamically calculated based on `OpportunityProducts(quantity * price - discount)`. An Opportunity *must* have products before transitioning to `DEAL_WON`.
*   **Automated Customer Conversion:** Once a Pipeline stage hits `DEAL_WON`, it is immediately and automatically converted into a `Customer` (in master data).
*   **Customer Analytics:** Inside the resulting Customer data, there will be an analytics breakdown (Purchase History) of exactly "What this customer has bought" derived from the Pipeline products.
*   **Quotation & Sales Order Conversion:**
    *   Opportunities at `QUALIFICATION` or `PROPOSAL` can be converted into `SalesQuotation`.
    *   Opportunities at `DEAL_WON` can be formally converted into a `SalesOrder`.
*   **Stage History Logging:** Every stage movement logs a `CrmOpportunityStageHistory`.

### 1.3 Tasks (Activities)
**Logic & Rules:**
*   **Structure & Unification:** Every interactions in the Sales workflow are tracked as `Activities`. **"Tasks" and "Schedules" are unified** into a single view; scheduling a calendar event is essentially creating a Task with a specific date bounding.
*   **Relationship:** Activities are polymorphic, optionally linked to `CustomerID`, `LeadID`, or `OpportunityID`.
*   **Quick Actions:** When viewing a Task/Schedule, there is a dedicated button to **"Add Lead"**. Clicking this creates a new Lead that is instantly and automatically linked to that specific Task.
*   **Activity Types:** `CALL`, `EMAIL`, `MEETING`, `VISIT`, `TASK`, `NOTE`, `QUOTATION`, `PRESENTATION`.
*   **Status Lifecycle:** `PLANNED` -> `IN_PROGRESS` -> `COMPLETED` -> `CANCELLED`.
*   **Priorities:** `LOW`, `MEDIUM`, `HIGH`, `URGENT`. (Defaults to MEDIUM).
*   **Completion:** Completing a task allows user to write down an `Outcome` and `DurationMinutes`, locking the status to `COMPLETED`.

### 1.4 RBAC (Role Based Access Control)
**Logic:**
Permissions dictate UI interactions and Backend authorization logic:
*   **Leads:** `leads.view`, `leads.create`, `leads.edit`, `leads.delete`, `leads.convert` (Convert to customer button / API).
*   **Pipeline:** `pipeline.view`, `pipeline.create`, `pipeline.edit`, `pipeline.delete`, `pipeline.update_stage`, `pipeline.convert_quotation`, `pipeline.convert_sales_order`.
*   **Tasks:** `tasks.view`, `tasks.create`, `tasks.edit`, `tasks.delete`, `tasks.complete`.

---

## 2. Implementation Working Plan

### Phase 1: Preparation & Migration Alignment (`apps/api`) — ✅ COMPLETED
- [x] Lead Entity has BANT fields (`budget_confirmed`, `authority_confirmed`, `need_confirmed`, `timeline_confirmed`), score, probability & estimated_value. `LeadQualificationChecklist` entity already tracks detailed BANT with `CalculateScore()` method.
- [x] Pipeline Stages are fully data-driven via `pipeline_stages` table (configurable `code`, `order`, `is_won`, `is_lost`, `probability`). No hardcoded stage names.
- [x] `DealHistory` entity tracks all stage transitions with `from_stage`, `to_stage`, `probability`, `days_in_prev_stage`, `reason`, `changed_by`.
- [x] Task entity supports unified schedule (via `is_schedule_task`, `scheduled_start_time`, `scheduled_end_time`), type enum (`general`, `call`, `email`, `meeting`, `follow_up`), status enum (`pending`, `in_progress`, `completed`, `cancelled`), priority enum (`low`, `medium`, `high`, `urgent`), polymorphic linking (`lead_id`, `deal_id`, `account_id`, `contact_id`), and quick action fields.
- [x] **ENHANCED**: `ValidateStageRequirements` now enforces strict forward-only progression, blocks backward moves, requires products for Won stage, and requires positive deal value for proposal+ stages.
- [x] **ENHANCED**: `MoveStageWithValidation` now calls `convertDealToPurchaseHistory` when a deal reaches Won stage.

### Phase 2: Refactoring Backend Services (`apps/api/internal/...`) — ✅ COMPLETED
- [x] **Lead Conversion** (`Convert` function): Already implemented. Creates Account + Contact + Deal from qualified lead, batch-updates Activities and Visit Reports with a single query (no N+1), creates initial DealHistory, emits domain event.
- [x] **Pipeline Stage Transition** (`MoveStageWithValidation`): Strict forward-only progression, product requirement for Won stage, deal value check for proposal+ stages. `convertDealToPurchaseHistory` triggered on Won.
- [x] **Customer Purchase Analytics** API: `GET /accounts/:id/purchases` and `GET /accounts/:id/purchases/analytics` already available.
- [x] **NEW: `CreateLeadFromTask`** endpoint (`POST /tasks/:id/create-lead`): Quick action to create & link a lead from a task. Auto-populates assigned_to and account from task context. RBAC: `tasks.create_lead`.
- [x] **Backend RBAC enforcement**: `leads.convert` → Lead Convert handler, `pipeline.update_stage` → MoveStage handler, `tasks.create_lead` → CreateLeadFromTask handler, `tasks.create` → TaskCreate handler. All critical mutation endpoints now check permissions.
- [x] **Error handling improved**: MoveStage handler now uses `ErrStageRequirementsNotMet` sentinel error + catch-all validation error pattern instead of hardcoded error string comparisons.

### Phase 3: Frontend Refactoring (`apps/web/src/features/...`) — ✅ COMPLETED
- [x] Update **Lead Detail Page** — replaced all tab placeholders with real content: Details (Contact Info, Location, Lead Info cards with icons), BANT qualification, Visit Reports from API, Activities from API. `cursor-pointer` on all tabs.
- [x] Update **Deal Detail Page** — refactored from custom layout to `<PageDetailLayout>` with `DealDetailTabs` (6 tabs: Details, Products, Tasks, Visits, Activities, History). Added `PermissionGuard`, `PageMotion`, proper empty states. Products tab renders line items table. History tab integrates `DealHistoryTimeline`.
- [x] Refactor **Task Detail Page** — refactored from custom layout to `<PageDetailLayout>` with status/priority badges, Related Deal/Lead display, "Create & Link Lead" quick action button, Reminders section, schedule info support. Added `PermissionGuard`.
- [x] **Consistent UI**: All 3 detail pages (Lead, Deal, Task) now use `PageDetailLayout` with back button, title, subtitle (badges), and action buttons. All tabs use identical icon + label pattern with responsive hiding.
- [x] **UX safety**: All interactive elements have `cursor-pointer`, all data accesses use optional chaining (`?.`) and nullish coalescing (`??`), all lists have meaningful empty states with icons.
- [x] **Type safety**: Fixed all TypeScript errors — proper `VisitReport` and `Activity` types, renamed lucide `Activity` to `ActivityIcon` to avoid naming conflicts, fixed `formatDate` parameter type in `CustomerPurchaseHistory.tsx`, fixed `useTaskQuickActions` broken import (replaced `@/hooks/useToast` with `sonner` toast).
- [x] **Build**: `npx next build` → exit code 0, compiled successfully, TypeScript passes, all pages generated.

### Phase 4: Integrations & Testing
- [ ] Integration test between Lead conversion sequence and Pipeline Data endpoints.
- [ ] Integration test between Deal Won transition constraint (requiring Products) and automated Customer + Analytics Creation.
- [ ] Finalize the Frontend Role checking logic (Hiding Convert / Complete / Edit buttons from unauthorized viewers).
- [ ] Execute `docker compose up --build -d` & ensure local container is healthy.
