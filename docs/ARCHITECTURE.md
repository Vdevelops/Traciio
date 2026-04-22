# Architecture Documentation
## CRM Healthcare Platform

**Versi**: 1.0  
**Last Updated**: 2025-01-15

---

## Overview

Dokumen ini menjelaskan arsitektur sistem CRM Healthcare Platform, termasuk struktur codebase, design patterns, dan best practices.

---

## Arsitektur Backend (Go)

### Layered Architecture

Sistem menggunakan **Layered Architecture** dengan pemisahan yang jelas:

```
┌─────────────────────────────────────┐
│   Handler Layer (HTTP)              │
│   - Request validation              │
│   - Response formatting             │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Service Layer (Business Logic)    │
│   - Business rules                  │
│   - Domain logic                    │
│   - Transaction management          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Repository Layer (Data Access)    │
│   - Database operations             │
│   - Data mapping                    │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Database (PostgreSQL)             │
└─────────────────────────────────────┘
```

### Layer Responsibilities

#### 1. Handler Layer (`internal/api/handlers/`)

**Responsibilities:**
- Parse HTTP requests
- Validate request data (using validator)
- Call service layer
- Format responses (using response helpers)
- Handle HTTP-specific concerns

**Example:**
```go
func (h *AuthHandler) Login(c *gin.Context) {
    var req auth.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        errors.HandleValidationError(c, err)
        return
    }
    
    loginResponse, err := h.authService.Login(&req)
    if err != nil {
        errors.ErrorResponse(c, "INVALID_CREDENTIALS", nil, nil)
        return
    }
    
    response.SuccessResponse(c, loginResponse, nil)
}
```

#### 2. Service Layer (`internal/service/`)

**Responsibilities:**
- Business logic
- Domain rules validation
- Orchestrate repository calls
- Transaction management
- Error handling

**Example:**
```go
func (s *Service) Login(req *auth.LoginRequest) (*auth.LoginResponse, error) {
    user, err := s.repo.FindByEmail(req.Email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    
    // Business logic: verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, ErrInvalidCredentials
    }
    
    // Generate tokens
    accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
    // ...
}
```

#### 3. Repository Layer (`internal/repository/`)

**Responsibilities:**
- Database operations
- Data mapping (DB → Domain)
- Query optimization
- No business logic

**Example:**
```go
func (r *repository) FindByEmail(email string) (*auth.User, error) {
    var user auth.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

### Dependency Injection Pattern

Semua dependencies di-inject melalui constructor:

```go
// Repository
authRepo := auth.NewRepository(database.DB)

// Service (depends on Repository)
authService := authservice.NewService(authRepo, jwtManager)

// Handler (depends on Service)
authHandler := handlers.NewAuthHandler(authService)
```

**Benefits:**
- Testable (easy to mock)
- Flexible (easy to swap implementations)
- Clear dependencies

### Interface-Based Design

Repository menggunakan interface untuk abstraction:

```go
// Interface
type AuthRepository interface {
    FindByEmail(email string) (*auth.User, error)
    FindByID(id string) (*auth.User, error)
    Create(user *auth.User) error
}

// Implementation
type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.AuthRepository {
    return &repository{db: db}
}
```

**Benefits:**
- Easy to test (mock interface)
- Easy to swap implementations
- Clear contracts

---

## Arsitektur Frontend (Next.js)

### Feature-Based Structure

```
apps/web/src/
├── features/
│   └── auth/
│       ├── components/      # Feature components
│       ├── hooks/           # Feature hooks
│       ├── services/        # API services
│       ├── stores/          # Zustand stores
│       ├── types/           # TypeScript types
│       └── schemas/         # Zod schemas
├── components/
│   ├── ui/                  # shadcn/ui components
│   └── shared/              # Shared components
└── lib/
    ├── api-client/          # API client
    └── utils/               # Utilities
```

### State Management

**Zustand** untuk client state:
- Auth state
- UI state
- Local state

**TanStack Query** untuk server state:
- API data caching
- Automatic refetching
- Background updates

### Form Handling

**React Hook Form + Zod**:
- Type-safe forms
- Client-side validation
- Server-side validation integration

---

## Design Patterns

### 1. Repository Pattern

Abstraction layer untuk data access:

```go
type AuthRepository interface {
    FindByEmail(email string) (*auth.User, error)
}

type repository struct {
    db *gorm.DB
}
```

### 2. Service Pattern

Business logic encapsulation:

```go
type Service struct {
    repo      interfaces.AuthRepository
    jwtManager *jwt.JWTManager
}
```

### 3. Dependency Injection

Constructor-based injection:

```go
func NewService(repo interfaces.AuthRepository, jwtManager *jwt.JWTManager) *Service {
    return &Service{
        repo:       repo,
        jwtManager: jwtManager,
    }
}
```

### 4. Middleware Pattern

Cross-cutting concerns:

```go
router.Use(middleware.LoggerMiddleware())
router.Use(middleware.CORSMiddleware())
router.Use(middleware.RequestIDMiddleware())
router.Use(middleware.LocaleMiddleware())
```

---

## Error Handling Strategy

### Backend

1. **Domain Errors**: Defined in service layer
2. **Error Mapping**: Map domain errors to API error codes
3. **Error Response**: Format menggunakan response helpers
4. **Bilingual Support**: Error messages dalam ID & EN

```go
// Service returns domain error
if err == authservice.ErrInvalidCredentials {
    errors.ErrorResponse(c, "INVALID_CREDENTIALS", nil, nil)
    return
}
```

### Frontend

1. **Zod Validation**: Client-side validation
2. **API Error Handling**: Parse error dari API response
3. **Error Boundary**: Catch React errors
4. **User-Friendly Messages**: Display error messages

---

## Logging Strategy

### Backend

- **Structured Logging**: Using `pkg/logger`
- **Request Logging**: Via middleware
- **Error Logging**: With context
- **Log Levels**: Info, Error, Debug

```go
logger.LogError(err, map[string]interface{}{
    "user_id": userID,
    "action":  "login",
})
```

### Frontend

- **Console Logging**: Development only
- **Error Logging**: Error boundary
- **API Error Logging**: Via interceptors

---

## Security

### Authentication

- **JWT**: Access token + Refresh token
- **Password Hashing**: bcrypt
- **Token Validation**: Middleware

### Authorization

- **Role-Based**: User roles
- **Permission-Based**: Fine-grained permissions (future)

### Data Protection

- **Input Validation**: Zod (frontend) + Validator (backend)
- **SQL Injection**: GORM parameterized queries
- **XSS**: React automatic escaping
- **CORS**: Configured middleware

---

## Scalability Considerations

### Backend

1. **Stateless API**: JWT-based, no server-side sessions
2. **Database Connection Pooling**: GORM handles this
3. **Horizontal Scaling**: Stateless design enables this
4. **Caching**: Ready for Redis integration

### Frontend

1. **Code Splitting**: Next.js automatic
2. **Image Optimization**: Next.js Image component
3. **API Caching**: TanStack Query
4. **Bundle Optimization**: Next.js automatic

---

## Testing Strategy

### Backend

- **Unit Tests**: Per layer (service, repository)
- **Integration Tests**: API endpoints
- **Mocking**: Interface-based mocking

### Frontend

- **Unit Tests**: Components, hooks
- **Integration Tests**: User flows
- **E2E Tests**: Critical paths (future)

---

## Deployment

### Backend

- **Docker**: Containerized
- **Docker Compose**: Local development
- **Environment Variables**: Config management

### Frontend

- **Next.js Build**: Static/SSR optimization
- **Environment Variables**: Runtime config
- **CDN**: Static assets

---

## Monitoring & Observability

### Backend

- **Request ID**: Track requests
- **Structured Logging**: Easy to parse
- **Error Tracking**: Error codes + context

### Frontend

- **Error Boundary**: Catch errors
- **API Error Tracking**: Request ID correlation

---

## Best Practices

### Code Organization

1. **Feature-Based**: Group by feature
2. **Separation of Concerns**: Clear layer boundaries
3. **DRY**: Don't repeat yourself
4. **SOLID Principles**: Especially Single Responsibility

### Error Handling

1. **Consistent Format**: Use response helpers
2. **Bilingual**: Always include ID & EN
3. **Context**: Include relevant details
4. **Logging**: Log errors with context

### Testing

1. **Test Each Layer**: Unit tests per layer
2. **Mock Dependencies**: Use interfaces
3. **Integration Tests**: Test full flows
4. **Coverage**: Aim for >80%

---

**Dokumen ini akan diupdate sesuai dengan perkembangan development.**

