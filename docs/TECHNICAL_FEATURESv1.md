# Technical Features Documentation

## 📋 Table of Contents

1. [Tech Stack Overview](#tech-stack-overview)
2. [Backend Architecture](#backend-architecture)
3. [Frontend Architecture](#frontend-architecture)
4. [Security Features](#security-features)
5. [Backend Modules](#backend-modules)
6. [Frontend Modules](#frontend-modules)
7. [Infrastructure & DevOps](#infrastructure--devops)
8. [Quality Assessment](#quality-assessment)

---

## Tech Stack Overview

### Backend (API)

**Language & Framework:**
- **Go 1.25.4** - Modern, performant, type-safe language
- **Gin Framework v1.11.0** - High-performance HTTP web framework
- **GORM v1.31.1** - Full-featured ORM for database operations

**Database:**
- **PostgreSQL** - Primary relational database
- **GORM AutoMigrate** - Database schema migration system

**Authentication & Security:**
- **JWT (golang-jwt/jwt/v5)** - Token-based authentication
- **bcrypt (golang.org/x/crypto)** - Password hashing
- **Rate Limiting (golang.org/x/time)** - Multi-level rate limiting
- **CORS (gin-contrib/cors)** - Cross-origin resource sharing
- **HSTS** - HTTP Strict Transport Security

**Storage:**
- **Local Storage** - File system storage (development)
- **Cloudflare R2** - S3-compatible object storage (production)
- **AWS SDK v2** - R2/S3 integration

**AI Integration:**
- **Cerebras AI** - LLM integration for AI features
- Custom AI service with context-aware prompts

**Real-time Communication:**
- **Gorilla WebSocket** - WebSocket support for real-time notifications

**Additional Libraries:**
- **UUID (google/uuid)** - Unique identifier generation
- **Excelize (xuri/excelize/v2)** - Excel file generation for reports
- **Validator (go-playground/validator/v10)** - Input validation

### Frontend (Web)

**Framework & Runtime:**
- **Next.js 16.0.10** - React framework with App Router
- **React 19.2.3** - UI library
- **TypeScript 5** - Type-safe JavaScript

**State Management:**
- **TanStack Query v5.90.10** - Server state management & caching
- **Zustand v5.0.8** - Client state management

**Form Handling:**
- **React Hook Form v7.66.1** - Form state management
- **Zod v4.1.12** - Schema validation
- **@hookform/resolvers** - RHF + Zod integration

**UI Components:**
- **Radix UI** - Headless UI primitives
  - Dialog, Dropdown, Select, Tabs, Popover, etc.
- **Tailwind CSS v4** - Utility-first CSS framework
- **shadcn/ui** - Component library built on Radix UI
- **Lucide React** - Icon library
- **Framer Motion** - Animation library

**Internationalization:**
- **next-intl v4.5.1** - i18n solution for Next.js

**Data Visualization:**
- **Recharts v3.5.0** - Chart library
- **React Markdown** - Markdown rendering

**HTTP Client:**
- **Axios v1.13.2** - HTTP client with interceptors

**Utilities:**
- **date-fns v4.1.0** - Date manipulation
- **class-variance-authority** - Component variants
- **cmdk** - Command palette
- **sonner** - Toast notifications

---

## Backend Architecture

### Architecture Pattern

**Layered Architecture:**
```
Handler → Service → Repository → Database
```

**Key Principles:**
- **Dependency Injection** - Constructor-based DI
- **Interface-based Design** - Repository interfaces for testability
- **Separation of Concerns** - Clear boundaries between layers
- **Domain-Driven Design** - Domain entities in `internal/domain/`

### Project Structure

```
apps/api/
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── middleware/  # HTTP middleware (auth, rate limit, CORS, HSTS)
│   │   └── routes/      # Route definitions
│   ├── config/          # Configuration management
│   ├── database/         # Database connection & migrations
│   ├── domain/           # Domain entities (20+ entities)
│   ├── repository/
│   │   ├── interfaces/  # Repository interfaces
│   │   └── postgres/    # PostgreSQL implementations
│   ├── service/          # Business logic layer
│   ├── hub/              # WebSocket hub for notifications
│   └── worker/           # Background workers
├── pkg/                  # Shared packages
│   ├── cerebras/         # AI client
│   ├── errors/           # Error handling
│   ├── jwt/              # JWT utilities
│   ├── logger/           # Logging
│   └── response/         # API response helpers
└── seeders/              # Database seeders
```

### Key Features

1. **Multi-level Rate Limiting**
   - Level 1: IP-based (5 requests/15min for login)
   - Level 2: Email-based (10 requests/15min per email)
   - Level 3: Global limit (100 requests/min globally)
   - Different limits for: login, refresh, upload, general, public endpoints

2. **JWT Authentication**
   - Access tokens (24 hours default)
   - Refresh tokens (7 days default)
   - Token rotation on refresh
   - Automatic token refresh interceptor

3. **WebSocket Hub**
   - Real-time notifications
   - User-specific message routing
   - Connection management with ping/pong
   - Cloudflare proxy compatibility

4. **Background Workers**
   - Reminder worker (runs every 1 minute)
   - Refresh token cleanup worker (runs every 24 hours)

5. **File Storage Abstraction**
   - Interface-based storage (local/R2)
   - Seamless switching between storage backends
   - File upload validation (MIME, size)

6. **Database Safety**
   - Production-safe migrations (never drops tables in production)
   - Auto-migration with constraint handling
   - Soft deletes (GORM DeletedAt)

---

## Frontend Architecture

### Architecture Pattern

**Feature-Based Architecture:**
```
features/
├── <feature-name>/
│   ├── components/    # UI components
│   ├── hooks/         # Business logic hooks
│   ├── services/      # API services
│   ├── stores/        # Zustand stores (if needed)
│   ├── schemas/       # Zod schemas
│   ├── types/         # TypeScript types
│   └── i18n/          # Translations
```

**Key Principles:**
- **Separation of Concerns** - Logic in hooks, UI in components
- **Type Safety** - Full TypeScript coverage, no `any` types
- **Server State** - TanStack Query for API data
- **Client State** - Zustand for UI state
- **Form Validation** - Zod schemas with React Hook Form

### Project Structure

```
apps/web/
├── app/                    # Next.js App Router
│   └── [locale]/          # Internationalized routes
├── src/
│   ├── components/        # Shared UI components
│   │   ├── ui/            # shadcn/ui components
│   │   ├── layouts/       # Layout components
│   │   └── navigation/    # Navigation components
│   ├── features/          # Feature modules
│   ├── lib/               # Utilities & helpers
│   │   ├── api-client.ts  # Axios instance with interceptors
│   │   └── react-query.tsx # TanStack Query provider
│   └── i18n/              # i18n configuration
└── public/                # Static assets
```

### Key Features

1. **API Client with Interceptors**
   - Automatic token injection
   - Token refresh on 401
   - Rate limit handling with countdown
   - Error formatting & toast notifications
   - Request queuing during token refresh

2. **TanStack Query Configuration**
   - 1-minute stale time
   - No refetch on window focus
   - Smart retry logic (skip network errors)
   - React Query DevTools in development

3. **Internationalization**
   - next-intl integration
   - Locale-based routing
   - Translation files per feature

4. **Error Handling**
   - Error boundaries
   - Toast notifications
   - Form-level error display
   - Network error handling

5. **UI/UX Safety**
   - Null safety with optional chaining
   - Loading states with skeletons
   - Empty states
   - Error states with retry

---

## Security Features

### Backend Security

1. **Authentication & Authorization**
   - JWT-based authentication
   - Role-based access control (RBAC)
   - Permission-based access control
   - Token expiration & rotation

2. **Rate Limiting**
   - Multi-level rate limiting (IP, email, global)
   - Configurable limits per endpoint type
   - Rate limit headers (X-RateLimit-*)
   - Stable reset time calculation

3. **CORS Protection**
   - Whitelist-based CORS
   - Credential support
   - Environment-based origin configuration

4. **HSTS (HTTP Strict Transport Security)**
   - Configurable max-age
   - Subdomain inclusion
   - Preload support

5. **Input Validation**
   - Struct tag validation
   - Parameterized queries (SQL injection prevention)
   - File upload validation (MIME, size, extension)

6. **Password Security**
   - bcrypt hashing
   - Never stored in plain text

7. **Error Handling**
   - Sanitized error messages
   - No sensitive data in logs
   - Structured error responses

### Frontend Security

1. **Token Management**
   - Secure cookie storage
   - localStorage for tokens (with refresh mechanism)
   - Automatic token refresh
   - Token cleanup on logout

2. **API Security**
   - HTTPS-only requests
   - Token injection in headers
   - Request timeout (10 seconds)

3. **Input Validation**
   - Zod schema validation
   - Client-side validation
   - Server-side validation (via API)

4. **XSS Prevention**
   - React's built-in XSS protection
   - Markdown sanitization
   - Safe HTML rendering

---

## Backend Modules

### 1. Authentication Module (`auth`)

**Features:**
- User login with email/password
- JWT token generation (access + refresh)
- Token refresh endpoint
- Password hashing with bcrypt
- Refresh token rotation
- Refresh token cleanup worker

**Endpoints:**
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

**Entities:**
- `User` - User account
- `RefreshToken` - Refresh token storage

---

### 2. User Management Module (`user`)

**Features:**
- CRUD operations for users
- User profile management
- Avatar upload
- User status management (active/inactive)
- Role assignment
- User statistics (deals, tasks, visit reports)

**Endpoints:**
- `GET /api/v1/users` - List users
- `GET /api/v1/users/:id` - Get user details
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /api/v1/users/:id/profile` - Get user profile
- `PUT /api/v1/users/:id/profile` - Update profile

**Entities:**
- `User` - User entity
- `UserProfile` - Extended user profile

---

### 3. Role & Permission Module (`role`, `permission`)

**Features:**
- Role CRUD operations
- Permission CRUD operations
- Role-permission assignment
- User-role assignment
- Menu-based permission system
- Permission checking middleware

**Endpoints:**
- Role endpoints: `/api/v1/roles/*`
- Permission endpoints: `/api/v1/permissions/*`

**Entities:**
- `Role` - User role
- `Permission` - Permission definition
- `Menu` - Menu item for permissions

---

### 4. Account Management Module (`account`)

**Features:**
- Account CRUD operations
- Account categorization
- Account status management
- Account assignment to sales reps
- Account search & filtering
- Account statistics

**Endpoints:**
- `GET /api/v1/accounts` - List accounts
- `GET /api/v1/accounts/:id` - Get account details
- `POST /api/v1/accounts` - Create account
- `PUT /api/v1/accounts/:id` - Update account
- `DELETE /api/v1/accounts/:id` - Delete account

**Entities:**
- `Account` - Account entity (Hospital, Clinic, Pharmacy)
- `Category` - Account category

---

### 5. Contact Management Module (`contact`)

**Features:**
- Contact CRUD operations
- Contact role management
- Contact-account linking
- Contact search & filtering
- Contact assignment

**Endpoints:**
- `GET /api/v1/contacts` - List contacts
- `GET /api/v1/contacts/:id` - Get contact details
- `POST /api/v1/contacts` - Create contact
- `PUT /api/v1/contacts/:id` - Update contact
- `DELETE /api/v1/contacts/:id` - Delete contact

**Entities:**
- `Contact` - Contact entity
- `ContactRole` - Contact role definition

---

### 6. Lead Management Module (`lead`)

**Features:**
- Lead CRUD operations
- Lead status management (new, contacted, qualified, converted, lost)
- Lead scoring (0-100)
- Lead source tracking
- Lead conversion to account/contact/deal
- Lead assignment
- Lead search & filtering

**Endpoints:**
- `GET /api/v1/leads` - List leads
- `GET /api/v1/leads/:id` - Get lead details
- `POST /api/v1/leads` - Create lead
- `PUT /api/v1/leads/:id` - Update lead
- `DELETE /api/v1/leads/:id` - Delete lead
- `POST /api/v1/leads/:id/convert` - Convert lead

**Entities:**
- `Lead` - Lead entity

---

### 7. Pipeline & Deal Management Module (`pipeline`, `deal`)

**Features:**
- Pipeline stage management
- Deal CRUD operations
- Deal stage movement (Kanban)
- Deal value tracking
- Deal probability tracking
- Expected close date
- Deal status (open, won, lost)
- Deal assignment
- Pipeline summary statistics

**Endpoints:**
- Pipeline: `/api/v1/pipelines/*`
- Deal: `/api/v1/deals/*`

**Entities:**
- `PipelineStage` - Pipeline stage
- `Deal` - Sales deal/opportunity

---

### 8. Product Management Module (`product`)

**Features:**
- Product CRUD operations
- Product category management
- Product search & filtering
- Product pricing
- Product status management

**Endpoints:**
- `GET /api/v1/products` - List products
- `GET /api/v1/products/:id` - Get product details
- `POST /api/v1/products` - Create product
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product

**Entities:**
- `Product` - Product entity
- `ProductCategory` - Product category

---

### 9. Task Management Module (`task`, `reminder`)

**Features:**
- Task CRUD operations
- Task status management (pending, in_progress, completed, cancelled)
- Task priority levels (low, medium, high, urgent)
- Task assignment
- Task due dates
- Task linking (account, contact, deal)
- Reminder system
- Reminder worker (background notifications)

**Endpoints:**
- `GET /api/v1/tasks` - List tasks
- `GET /api/v1/tasks/:id` - Get task details
- `POST /api/v1/tasks` - Create task
- `PUT /api/v1/tasks/:id` - Update task
- `DELETE /api/v1/tasks/:id` - Delete task
- Reminder endpoints: `/api/v1/tasks/:id/reminders/*`

**Entities:**
- `Task` - Task entity
- `Reminder` - Reminder entity

---

### 10. Visit Report Module (`visit_report`)

**Features:**
- Visit report CRUD operations
- Photo upload (multiple photos)
- Activity tracking
- Activity type management
- Visit report linking (account, contact)
- Visit report search & filtering
- Visit statistics

**Endpoints:**
- `GET /api/v1/visit-reports` - List visit reports
- `GET /api/v1/visit-reports/:id` - Get visit report details
- `POST /api/v1/visit-reports` - Create visit report
- `PUT /api/v1/visit-reports/:id` - Update visit report
- `DELETE /api/v1/visit-reports/:id` - Delete visit report
- `POST /api/v1/visit-reports/:id/photos` - Upload photos

**Entities:**
- `VisitReport` - Visit report entity
- `Activity` - Activity entity
- `ActivityType` - Activity type definition

---

### 11. Dashboard Module (`dashboard`)

**Features:**
- Dashboard overview statistics
- Deal statistics (total, open, won, lost)
- Revenue tracking
- Lead statistics
- Account statistics
- Activity trends
- Top accounts
- Top sales reps
- Upcoming tasks
- Recent activities

**Endpoints:**
- `GET /api/v1/dashboard/overview` - Get dashboard overview

**Entities:**
- `DashboardOverview` - Dashboard data DTO

---

### 12. Report Module (`report`)

**Features:**
- Sales funnel report
- Sales performance report
- Pipeline report
- Visit report analytics
- Excel export functionality

**Endpoints:**
- `GET /api/v1/reports/sales-funnel` - Sales funnel report
- `GET /api/v1/reports/sales-performance` - Sales performance report
- `GET /api/v1/reports/pipeline` - Pipeline report
- `GET /api/v1/reports/visit` - Visit report analytics

**Entities:**
- Report DTOs for various report types

---

### 13. Notification Module (`notification`)

**Features:**
- Real-time notifications via WebSocket
- Notification CRUD operations
- Notification read/unread status
- Notification broadcasting
- Notification hub (WebSocket management)

**Endpoints:**
- `GET /api/v1/notifications` - List notifications
- `GET /api/v1/notifications/:id` - Get notification
- `PUT /api/v1/notifications/:id/read` - Mark as read
- `DELETE /api/v1/notifications/:id` - Delete notification
- `WS /api/v1/ws` - WebSocket connection

**Entities:**
- `Notification` - Notification entity

---

### 14. AI Module (`ai`, `ai_settings`)

**Features:**
- AI chatbot integration (Cerebras)
- Visit report analysis
- Context-aware AI responses
- AI settings management
- Model selection
- AI usage tracking

**Endpoints:**
- `POST /api/v1/ai/chat` - Chat with AI
- `POST /api/v1/ai/analyze-visit-report/:id` - Analyze visit report
- `GET /api/v1/ai/settings` - Get AI settings
- `PUT /api/v1/ai/settings` - Update AI settings

**Entities:**
- `AISettings` - AI configuration
- `AIModelUsage` - AI usage tracking

---

### 15. Category Module (`category`)

**Features:**
- Category CRUD operations
- Category code management
- Category badge colors
- Category status management

**Endpoints:**
- `GET /api/v1/categories` - List categories
- `POST /api/v1/categories` - Create category
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

**Entities:**
- `Category` - Category entity

---

## Frontend Modules

### 1. Authentication Module (`auth`)

**Components:**
- `LoginForm` - Login form with validation
- `AuthGuard` - Route protection
- `PermissionGuard` - Permission-based access control
- `AuthLayout` - Authentication layout

**Hooks:**
- `useLogin` - Login mutation
- `useLogout` - Logout handler
- `useAuthGuard` - Auth guard logic
- `useRefreshSession` - Session refresh

**Features:**
- Email/password login
- Token management
- Automatic token refresh
- Protected routes
- Permission-based access

---

### 2. Dashboard Module (`dashboard`)

**Components:**
- `DashboardOverview` - Main dashboard view
- `ActivityTrends` - Activity trend charts
- `LeadsBySource` - Lead source chart
- `PipelineSummary` - Pipeline statistics
- `TopAccounts` - Top accounts list
- `TopSalesRep` - Top sales rep list
- `UpcomingTasks` - Upcoming tasks list
- `RecentActivities` - Recent activities list
- `LeadsTable` - Leads data table
- `VisitStatistics` - Visit statistics

**Hooks:**
- `useDashboard` - Dashboard data fetching

**Features:**
- Real-time dashboard statistics
- Interactive charts (Recharts)
- Data tables with sorting/filtering
- Responsive layout

---

### 3. User Management Module (`master-data/user-management`)

**Components:**
- `UserManagement` - Main user management view
- `UserList` - User list with CRUD
- `UserForm` - User create/edit form
- `RoleList` - Role list with CRUD
- `RoleForm` - Role create/edit form
- `PermissionList` - Permission list
- `AssignPermissionsDialog` - Permission assignment dialog

**Hooks:**
- `useUsers` - User data management
- `useRoles` - Role data management
- `usePermissions` - Permission data management
- `useHasPermission` - Permission checking
- `useUserPermissions` - User permission management

**Features:**
- User CRUD operations
- Role management
- Permission management
- Permission assignment
- User-role assignment

---

### 4. Account Management Module (`sales-crm/account-management`)

**Components:**
- `AccountManagement` - Main account management view
- `AccountList` - Account list with CRUD
- `AccountForm` - Account create/edit form
- `AccountDetailModal` - Account detail view
- `CategoryList` - Category list
- `CategoryForm` - Category form
- `ContactList` - Contact list
- `ContactForm` - Contact form
- `ContactDetailModal` - Contact detail view
- `ContactRoleList` - Contact role list
- `ContactRoleForm` - Contact role form

**Hooks:**
- `useAccounts` - Account data management
- `useAccountList` - Account list with filters
- `useCategories` - Category management
- `useCategoryList` - Category list
- `useContacts` - Contact management
- `useContactList` - Contact list
- `useContactRoles` - Contact role management
- `useContactRoleList` - Contact role list

**Features:**
- Account CRUD operations
- Category management
- Contact management
- Contact role management
- Search & filtering
- Data tables with pagination

---

### 5. Lead Management Module (`sales-crm/lead-management`)

**Components:**
- `LeadManagement` - Main lead management view
- `LeadList` - Lead list with CRUD
- `LeadForm` - Lead create/edit form
- `LeadDetailModal` - Lead detail view
- `ConvertLeadDialog` - Lead conversion dialog

**Hooks:**
- `useLeads` - Lead data management
- `useLeadList` - Lead list with filters

**Features:**
- Lead CRUD operations
- Lead status management
- Lead conversion to account/contact/deal
- Lead scoring display
- Search & filtering

---

### 6. Pipeline Management Module (`sales-crm/pipeline-management`)

**Components:**
- `KanbanBoard` - Drag-and-drop Kanban board
- `DealCard` - Deal card component
- `DealForm` - Deal create/edit form
- `DealDetailModal` - Deal detail view
- `StagesList` - Pipeline stages list
- `StageForm` - Stage create/edit form
- `PipelineSummary` - Pipeline statistics
- `Forecast` - Sales forecast view

**Hooks:**
- `useKanbanBoard` - Kanban board logic
- `useDeals` - Deal data management
- `usePipelines` - Pipeline management
- `useStages` - Stage management
- `useForecast` - Forecast calculations

**Features:**
- Kanban board with drag-and-drop
- Deal CRUD operations
- Stage management
- Deal value tracking
- Probability tracking
- Sales forecast
- Pipeline statistics

---

### 7. Product Management Module (`sales-crm/product-management`)

**Components:**
- `ProductManagement` - Main product management view
- `ProductList` - Product list with CRUD
- `ProductForm` - Product create/edit form
- `ProductDetailModal` - Product detail view
- `CategoryList` - Product category list
- `CategoryForm` - Category form

**Hooks:**
- `useProducts` - Product data management
- `useProductList` - Product list with filters
- `useCategoryList` - Category list

**Features:**
- Product CRUD operations
- Product category management
- Search & filtering
- Product pricing display

---

### 8. Task Management Module (`sales-crm/task-management`)

**Components:**
- `TaskManagement` - Main task management view
- `TaskBoard` - Kanban board for tasks
- `TaskList` - Task list view
- `TaskCard` - Task card component
- `TaskForm` - Task create/edit form
- `TaskDetailModal` - Task detail view
- `ReminderSettings` - Reminder configuration

**Hooks:**
- `useTasks` - Task data management
- `useTaskList` - Task list with filters
- `useKanbanBoard` - Kanban board logic

**Features:**
- Task CRUD operations
- Kanban board view
- Task status management
- Priority levels
- Due date tracking
- Reminder system
- Task assignment

---

### 9. Visit Report Module (`sales-crm/visit-report`)

**Components:**
- `VisitReportManagement` - Main visit report view
- `VisitReportList` - Visit report list
- `VisitReportForm` - Visit report create/edit form
- `VisitReportDetailModal` - Visit report detail view
- `PhotoUploadDialog` - Photo upload dialog
- `ActivityTimeline` - Activity timeline view
- `CreateActivityDialog` - Activity creation dialog
- `ActivityTypeList` - Activity type list
- `ActivityTypeForm` - Activity type form
- `IconPicker` - Icon picker component

**Hooks:**
- `useVisitReports` - Visit report data management
- `useVisitReportList` - Visit report list with filters
- `useActivityTypes` - Activity type management
- `useActivityTypeList` - Activity type list

**Features:**
- Visit report CRUD operations
- Photo upload (multiple photos)
- Activity tracking
- Activity type management
- Activity timeline
- Search & filtering

---

### 10. Reports Module (`reports`)

**Components:**
- `ReportGenerator` - Main report generator
- `SalesFunnelViewer` - Sales funnel visualization
- `SalesFunnelTable` - Sales funnel data table
- `SalesFunnelInsights` - Sales funnel insights
- `SalesPerformanceReportViewer` - Sales performance report
- `PipelineReportViewer` - Pipeline report
- `VisitReportViewer` - Visit report analytics

**Hooks:**
- `useReports` - Report data fetching

**Features:**
- Sales funnel reports
- Sales performance reports
- Pipeline reports
- Visit report analytics
- Data visualization
- Export functionality

---

### 11. AI Module (`ai`)

**Components:**
- `Chatbot` - AI chatbot interface
- `AISettings` - AI settings configuration
- `VisitReportInsightsButton` - Visit report analysis button

**Hooks:**
- `useChat` - Chat mutation
- `useAISettings` - AI settings management
- `useAnalyzeVisitReport` - Visit report analysis
- `useConversationStorage` - Conversation persistence

**Features:**
- AI chatbot with context awareness
- Visit report analysis
- AI settings management
- Model selection
- Conversation history
- Chat templates

---

### 12. Notifications Module (`notifications`)

**Components:**
- `NotificationDrawer` - Notification drawer
- `NotificationList` - Notification list
- `NotificationBadge` - Notification badge

**Hooks:**
- `useNotifications` - Notification data management
- `useWebSocket` - WebSocket connection

**Features:**
- Real-time notifications
- WebSocket integration
- Notification read/unread status
- Notification drawer UI

---

### 13. Profile Module (`profile`)

**Components:**
- `ProfileSettings` - Profile settings form

**Hooks:**
- `useProfile` - Profile data management

**Features:**
- Profile update
- Avatar upload
- User information management

---

## Infrastructure & DevOps

### Docker Support

**Backend:**
- `Dockerfile` - Multi-stage build
- `docker-compose.yml` - Development environment
- `docker-compose.production.yml` - Production environment

**Frontend:**
- `Dockerfile` - Next.js standalone build
- `docker-compose.production.yml` - Production environment

### Database

- **PostgreSQL** - Primary database
- **GORM AutoMigrate** - Schema migrations
- **Seeders** - Database seeding for development

### Storage

- **Local Storage** - Development (file system)
- **Cloudflare R2** - Production (S3-compatible)

### Environment Configuration

- **Environment Variables** - Configuration via `.env`
- **Config Package** - Centralized configuration management
- **Production Safety** - Environment-based behavior

---

## Quality Assessment

### Code Quality

**Backend:**
- ✅ **Layered Architecture** - Clear separation of concerns
- ✅ **Interface-based Design** - Testable and maintainable
- ✅ **Error Handling** - Structured error responses
- ✅ **Input Validation** - Comprehensive validation
- ✅ **Security Best Practices** - Rate limiting, CORS, HSTS, JWT
- ✅ **Database Safety** - Production-safe migrations
- ✅ **Type Safety** - Strong typing with Go

**Frontend:**
- ✅ **TypeScript** - Full type coverage, no `any` types
- ✅ **Feature-based Structure** - Organized and scalable
- ✅ **Separation of Concerns** - Logic in hooks, UI in components
- ✅ **Form Validation** - Zod schemas with React Hook Form
- ✅ **Error Handling** - Comprehensive error boundaries
- ✅ **Loading States** - Skeletons and spinners
- ✅ **Empty States** - User-friendly empty states
- ✅ **Null Safety** - Optional chaining throughout

### Security Quality

**Backend:**
- ✅ **Multi-level Rate Limiting** - IP, email, global
- ✅ **JWT with Rotation** - Secure token management
- ✅ **Password Hashing** - bcrypt
- ✅ **CORS Whitelist** - No wildcard CORS
- ✅ **HSTS** - HTTP Strict Transport Security
- ✅ **Input Validation** - Parameterized queries
- ✅ **File Upload Validation** - MIME, size, extension

**Frontend:**
- ✅ **Token Management** - Secure storage and refresh
- ✅ **HTTPS Only** - Secure API communication
- ✅ **Input Validation** - Client and server-side
- ✅ **XSS Prevention** - React's built-in protection

### Architecture Quality

**Backend:**
- ✅ **Clean Architecture** - Handler → Service → Repository
- ✅ **Dependency Injection** - Constructor-based
- ✅ **Interface Segregation** - Repository interfaces
- ✅ **Single Responsibility** - Each layer has clear purpose

**Frontend:**
- ✅ **Feature-based Architecture** - Scalable structure
- ✅ **Component Composition** - Reusable components
- ✅ **State Management** - TanStack Query + Zustand
- ✅ **Type Safety** - Full TypeScript coverage

### Performance Quality

**Backend:**
- ✅ **GORM Connection Pooling** - Efficient database connections
- ✅ **Rate Limiting** - Prevents abuse
- ✅ **Background Workers** - Async processing
- ✅ **WebSocket Hub** - Efficient real-time communication

**Frontend:**
- ✅ **TanStack Query Caching** - Efficient data fetching
- ✅ **Code Splitting** - Next.js automatic code splitting
- ✅ **Image Optimization** - Next.js image optimization
- ✅ **Standalone Build** - Optimized production build

### Maintainability Quality

**Backend:**
- ✅ **Clear Structure** - Well-organized codebase
- ✅ **Documentation** - Code comments and structure
- ✅ **Error Codes** - Standardized error handling
- ✅ **Configuration Management** - Centralized config

**Frontend:**
- ✅ **Feature Modules** - Self-contained features
- ✅ **Reusable Components** - Shared UI components
- ✅ **Type Definitions** - Comprehensive types
- ✅ **Internationalization** - Multi-language support

### Overall Assessment

**Strengths:**
1. **Modern Tech Stack** - Latest versions of frameworks
2. **Security First** - Comprehensive security measures
3. **Type Safety** - Strong typing throughout
4. **Scalable Architecture** - Feature-based structure
5. **Real-time Features** - WebSocket support
6. **AI Integration** - Modern AI capabilities
7. **Comprehensive Modules** - Full CRM functionality

**Areas for Improvement:**
1. **Testing** - Unit and integration tests needed
2. **Documentation** - API documentation (OpenAPI/Swagger)
3. **Monitoring** - Application monitoring and logging
4. **CI/CD** - Automated testing and deployment
5. **Performance Monitoring** - APM tools integration

**Overall Grade: A-**

The codebase demonstrates high-quality engineering practices with modern technologies, strong security measures, and a well-structured architecture. The main areas for improvement are testing coverage and operational observability.

---

## Summary

This CRM Healthcare application is a comprehensive, production-ready system with:

- **20+ Backend Modules** covering all CRM functionality
- **13+ Frontend Modules** with modern React/Next.js architecture
- **Enterprise-grade Security** with multi-level rate limiting, JWT, CORS, HSTS
- **Real-time Features** with WebSocket notifications
- **AI Integration** with Cerebras for intelligent insights
- **Modern Tech Stack** with Go, Next.js, TypeScript, TanStack Query
- **Scalable Architecture** with feature-based structure
- **Production Ready** with Docker, environment configuration, and safety measures

The system is well-architected, secure, and maintainable, making it suitable for enterprise healthcare CRM use cases.


