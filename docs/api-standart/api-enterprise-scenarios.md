# API Enterprise Scenarios

## Enterprise-Level Use Cases and Requirements

**Version**: 1.0  
**Status**: Active

---

## Overview

This document outlines critical enterprise-level scenarios and use cases that must be handled correctly in API development. These scenarios are common in production environments and require specific handling patterns.

---

## Multi-Tenancy

### Scenario: Tenant Isolation

**Requirement**: All data must be isolated by tenant. Users from one tenant cannot access data from another tenant.

**Implementation Pattern**:
```go
// Always filter by tenant_id
func (r *Repository) List(req *ListRequest) ([]Entity, error) {
    query := r.db.Where("tenant_id = ?", req.TenantID)
    // ... additional filters
    return query.Find(&entities).Error
}
```

**Checklist**:
- [ ] All queries include tenant_id filter
- [ ] Tenant ID validated in middleware
- [ ] Cross-tenant access prevented
- [ ] Tenant context passed through all layers

---

## Concurrency & Race Conditions

### Scenario: Concurrent Updates

**Requirement**: Prevent race conditions when multiple users update the same resource simultaneously.

**Implementation Pattern**:
```go
// Use database transactions with row-level locking
err := db.Transaction(func(tx *gorm.DB) error {
    var entity Entity
    // Lock row for update
    if err := tx.Set("gorm:query_option", "FOR UPDATE").
        Where("id = ?", id).First(&entity).Error; err != nil {
        return err
    }
    
    // Update entity
    entity.Field = newValue
    return tx.Save(&entity).Error
})
```

**Checklist**:
- [ ] Critical operations use transactions
- [ ] Row-level locking for concurrent updates
- [ ] Optimistic locking with version fields
- [ ] Handle transaction conflicts gracefully

---

## Bulk Operations

### Scenario: Bulk Import/Export

**Requirement**: Handle large-scale data import/export operations efficiently.

**Implementation Pattern**:
```go
// Process in batches
const batchSize = 1000
for i := 0; i < len(items); i += batchSize {
    end := i + batchSize
    if end > len(items) {
        end = len(items)
    }
    batch := items[i:end]
    
    // Process batch
    if err := processBatch(batch); err != nil {
        return err
    }
}
```

**Checklist**:
- [ ] Process in batches (1000-5000 records)
- [ ] Use background jobs for large operations
- [ ] Provide progress tracking
- [ ] Handle partial failures
- [ ] Validate data before processing

---

## Audit Trail

### Scenario: Track All Changes

**Requirement**: Maintain complete audit trail of all data changes for compliance and debugging.

**Implementation Pattern**:
```go
// Audit log entry
type AuditLog struct {
    ID          string
    ResourceType string
    ResourceID   string
    Action       string  // CREATE, UPDATE, DELETE
    UserID       string
    Changes      json.RawMessage  // Before/after values
    Timestamp    time.Time
}
```

**Checklist**:
- [ ] Log all CREATE, UPDATE, DELETE operations
- [ ] Include user ID and timestamp
- [ ] Store before/after values for updates
- [ ] Include request ID for correlation
- [ ] Retain audit logs per compliance requirements

---

## Soft Delete

### Scenario: Recoverable Deletions

**Requirement**: Support soft delete to allow data recovery and maintain referential integrity.

**Implementation Pattern**:
```go
// Use GORM soft delete
type Entity struct {
    ID        string
    DeletedAt gorm.DeletedAt `gorm:"index"`
    // ... other fields
}

// Soft delete
db.Delete(&entity)

// Hard delete (if needed)
db.Unscoped().Delete(&entity)
```

**Checklist**:
- [ ] Use soft delete by default
- [ ] Filter deleted records in queries
- [ ] Provide restore functionality
- [ ] Implement hard delete for compliance (GDPR, etc.)
- [ ] Clean up old soft-deleted records periodically

---

## Data Validation & Sanitization

### Scenario: Input Validation

**Requirement**: Validate and sanitize all user inputs to prevent security vulnerabilities.

**Implementation Pattern**:
```go
// Use struct validation
type CreateRequest struct {
    Name  string `binding:"required,min=3,max=100"`
    Email string `binding:"required,email"`
    Age   int    `binding:"min=18,max=100"`
}

// Sanitize input
func sanitizeInput(input string) string {
    // Remove HTML tags
    // Escape special characters
    // Trim whitespace
    return strings.TrimSpace(html.EscapeString(input))
}
```

**Checklist**:
- [ ] Validate all inputs with struct tags
- [ ] Sanitize string inputs (XSS prevention)
- [ ] Validate file uploads (type, size, content)
- [ ] Use parameterized queries (SQL injection prevention)
- [ ] Validate business rules

---

## File Upload & Storage

### Scenario: Secure File Handling

**Requirement**: Handle file uploads securely with validation and proper storage.

**Implementation Pattern**:
```go
// Validate file
func validateFile(fileHeader *multipart.FileHeader) error {
    // Check size
    if fileHeader.Size > 5*1024*1024 { // 5MB
        return errors.New("FILE_TOO_LARGE")
    }
    
    // Check MIME type
    file, _ := fileHeader.Open()
    buffer := make([]byte, 512)
    file.Read(buffer)
    mimeType := http.DetectContentType(buffer)
    
    allowedTypes := []string{"image/jpeg", "image/png"}
    if !contains(allowedTypes, mimeType) {
        return errors.New("INVALID_FILE_TYPE")
    }
    
    // Generate safe filename
    ext := filepath.Ext(fileHeader.Filename)
    filename := uuid.New().String() + ext
    
    return nil
}
```

**Checklist**:
- [ ] Validate file size (max 5-10MB for images)
- [ ] Validate MIME type (not just extension)
- [ ] Generate unique filenames (UUID)
- [ ] Store files outside web root
- [ ] Scan for malware (if applicable)
- [ ] Implement file access controls

---

## Background Jobs & Async Processing

### Scenario: Long-Running Operations

**Requirement**: Handle time-consuming operations asynchronously to avoid blocking API responses.

**Implementation Pattern**:
```go
// Queue job for async processing
func (h *Handler) GenerateReport(c *gin.Context) {
    jobID := uuid.New().String()
    
    // Queue job
    jobQueue.Enqueue(ReportJob{
        ID:     jobID,
        UserID: userID,
        Params: params,
    })
    
    // Return job ID immediately
    response.SuccessResponse(c, gin.H{
        "job_id": jobID,
        "status": "queued",
    }, nil)
}

// Poll for job status
func (h *Handler) GetJobStatus(c *gin.Context) {
    jobID := c.Param("id")
    status := jobQueue.GetStatus(jobID)
    response.SuccessResponse(c, status, nil)
}
```

**Checklist**:
- [ ] Use job queue for long operations (> 5 seconds)
- [ ] Return job ID immediately
- [ ] Provide status polling endpoint
- [ ] Handle job failures gracefully
- [ ] Implement job retry mechanism
- [ ] Set job timeout limits

---

## Rate Limiting & Throttling

### Scenario: Prevent Abuse

**Requirement**: Protect API from abuse and ensure fair resource usage.

**Implementation Pattern**:
```go
// Per-user rate limiting
func UserRateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := getUserID(c)
        limiter := getUserLimiter(userID)
        
        if !limiter.Allow() {
            errors.RateLimitResponse(c, ...)
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**Checklist**:
- [ ] Implement per-user rate limits
- [ ] Implement per-IP rate limits
- [ ] Different limits for different endpoints
- [ ] Return rate limit headers (X-RateLimit-*)
- [ ] Handle rate limit errors gracefully

---

## Error Recovery & Retry

### Scenario: Transient Failures

**Requirement**: Handle transient failures (network, database) with automatic retry.

**Implementation Pattern**:
```go
// Retry with exponential backoff
func retryWithBackoff(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        // Check if error is retryable
        if !isRetryableError(err) {
            return err
        }
        
        // Exponential backoff
        backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
        time.Sleep(backoff)
    }
    return errors.New("MAX_RETRIES_EXCEEDED")
}
```

**Checklist**:
- [ ] Identify retryable errors (network, timeout)
- [ ] Implement exponential backoff
- [ ] Set maximum retry attempts
- [ ] Log retry attempts
- [ ] Return appropriate error after max retries

---

## Data Migration & Versioning

### Scenario: Schema Changes

**Requirement**: Handle database schema changes without downtime.

**Implementation Pattern**:
```go
// Use migrations
func Migrate() error {
    return db.AutoMigrate(
        &EntityV1{},
        &EntityV2{},  // New version
    )
}

// Support multiple versions
func (h *Handler) GetEntity(c *gin.Context) {
    version := c.GetHeader("API-Version")
    switch version {
    case "v2":
        return h.getEntityV2(c)
    default:
        return h.getEntityV1(c)
    }
}
```

**Checklist**:
- [ ] Use database migrations
- [ ] Support backward compatibility
- [ ] Version API responses
- [ ] Test migrations on staging
- [ ] Rollback plan for failed migrations

---

## Security Scenarios

### Scenario: SQL Injection Prevention

**Requirement**: Prevent SQL injection attacks.

**Implementation Pattern**:
```go
// ✅ GOOD: Parameterized query
db.Where("name = ?", userInput).First(&entity)

// ❌ BAD: String concatenation
query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", userInput)
db.Raw(query).Scan(&entity)
```

### Scenario: XSS Prevention

**Requirement**: Prevent cross-site scripting attacks.

**Implementation Pattern**:
```go
// Sanitize user input
import "html"

sanitized := html.EscapeString(userInput)
```

### Scenario: CSRF Protection

**Requirement**: Protect against cross-site request forgery.

**Implementation Pattern**:
```go
// Use CSRF tokens
func CSRFMiddleware() gin.HandlerFunc {
    return csrf.Middleware(csrf.Options{
        Secret: csrfSecret,
        ErrorFunc: func(c *gin.Context) {
            errors.ForbiddenResponse(c, "CSRF_TOKEN_INVALID", ...)
        },
    })
}
```

**Security Checklist**:
- [ ] All inputs validated and sanitized
- [ ] Parameterized queries used (no string concatenation)
- [ ] XSS prevention (escape HTML)
- [ ] CSRF protection for state-changing operations
- [ ] Authentication tokens properly validated
- [ ] Sensitive data not logged
- [ ] HTTPS enforced
- [ ] Security headers set (HSTS, CSP, etc.)

---

## Monitoring & Alerting Scenarios

### Scenario: Production Incident Detection

**Requirement**: Detect and alert on production issues quickly.

**Implementation Pattern**:
```go
// Monitor critical metrics
func monitorEndpoint(endpoint string) {
    metrics := getMetrics(endpoint)
    
    // Alert on high error rate
    if metrics.ErrorRate > 0.05 { // 5%
        alert.Send("High error rate", endpoint)
    }
    
    // Alert on slow response time
    if metrics.P95ResponseTime > 1*time.Second {
        alert.Send("Slow response time", endpoint)
    }
}
```

**Checklist**:
- [ ] Monitor error rates
- [ ] Monitor response times
- [ ] Monitor resource usage (CPU, memory)
- [ ] Set up alerts for critical metrics
- [ ] Log all errors with context
- [ ] Track request IDs for tracing

---

## Distributed Caching & Locking

### Scenario: High-Traffic Read Operations

**Requirement**: Reduce database load for frequently accessed data using Redis.

**Implementation Pattern (Cache-Aside)**:
```go
// Get with Cache-Aside
func (u *Usecase) GetByID(ctx context.Context, id string) (*Entity, error) {
    key := fmt.Sprintf("entity:%s", id)
    
    // 1. Try Cache
    val, err := redis.Get(ctx, key).Result()
    if err == nil {
        var entity Entity
        if json.Unmarshal([]byte(val), &entity) == nil {
            return &entity, nil
        }
    }
    
    // 2. Fallback to DB
    entity, err := repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. Set Cache (Async or Blocking)
    go func() {
        data, _ := json.Marshal(entity)
        redis.Set(ctx, key, data, 15*time.Minute)
    }()
    
    return entity, nil
}
```

### Scenario: Distributed Resource Locking

**Requirement**: Prevent concurrent modification of shared resources across multiple instances (e.g., Inventory allocation).

**Implementation Pattern (Redis Lock)**:
```go
// Distributed Lock using Redis
func (u *Usecase) AllocateInventory(ctx context.Context, itemID string) error {
    lockKey := fmt.Sprintf("lock:inventory:%s", itemID)
    
    // Acquire Lock (10s TTL)
    acquired, err := redis.SetNX(ctx, lockKey, "locked", 10*time.Second).Result()
    if err != nil || !acquired {
        return errors.New("RESOURCE_LOCKED_PLEASE_RETRY")
    }
    defer redis.Del(ctx, lockKey)
    
    // ... Critical Section (DB Update) ...
    return repo.UpdateStock(ctx, itemID)
}
```

**Checklist**:
- [ ] Use Redis for shared cache
- [ ] Define TTL for all cached keys
- [ ] Handle cache failures gracefully (fallback to DB)
- [ ] Implement Distributed Locks for critical sections
- [ ] Use `SETNX` for locking
- [ ] Ensure locks are released (defer/TTL)

---

## Best Practices Summary

### Critical Requirements

1. **Multi-tenancy**: Always filter by tenant_id
2. **Concurrency**: Use transactions and row-level locking
3. **Bulk operations**: Process in batches
4. **Audit trail**: Log all changes
5. **Soft delete**: Use soft delete by default
6. **Validation**: Validate and sanitize all inputs
7. **File upload**: Validate type, size, and content
8. **Async processing**: Use jobs for long operations
9. **Rate limiting**: Protect against abuse
10. **Error recovery**: Retry transient failures
11. **Security**: Prevent common vulnerabilities
12. **Monitoring**: Track metrics and alert on issues

---

**This document will be updated according to enterprise requirements and feedback from the development team.**
