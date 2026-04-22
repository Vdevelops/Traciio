# Context Timeout Migration Guide

## Overview

This document describes the migration pattern for adding context timeout to all repository methods to prevent queries from hanging indefinitely.

## Pattern

### Before (No Timeout)
```go
type repository struct {
    db *gorm.DB
}

func (r *repository) FindByID(id string) (*Entity, error) {
    var entity Entity
    err := r.db.Where("id = ?", id).First(&entity).Error
    return &entity, err
}
```

### After (With Timeout)
```go
import (
    "github.com/gilabs/crm-healthcare/api/internal/repository"
)

type repository struct {
    *repository.BaseRepository
}

func NewRepository(db *gorm.DB) interfaces.EntityRepository {
    return &repository{
        BaseRepository: repository.NewBaseRepository(db),
    }
}

func (r *repository) FindByID(id string) (*Entity, error) {
    var entity Entity
    err := r.DBWithTimeout(repository.DefaultTimeout).
        Where("id = ?", id).
        First(&entity).Error
    return &entity, err
}
```

## Migration Steps

1. **Update repository struct** to embed `BaseRepository`:
   ```go
   type repository struct {
       *repository.BaseRepository
   }
   ```

2. **Update NewRepository** to initialize BaseRepository:
   ```go
   func NewRepository(db *gorm.DB) interfaces.EntityRepository {
       return &repository{
           BaseRepository: repository.NewBaseRepository(db),
       }
   }
   ```

3. **Update all methods** to use `DBWithTimeout`:
   ```go
   // Simple query
   err := r.DBWithTimeout(repository.DefaultTimeout).
       Where("id = ?", id).
       First(&entity).Error
   
   // Complex query with multiple operations
   err := r.WithTimeout(repository.DefaultTimeout, func(ctx context.Context, db *gorm.DB) error {
       // Use db (which has context) for all operations
       return db.Where("id = ?", id).First(&entity).Error
   })
   ```

## Timeout Values

- **DefaultTimeout**: 30 seconds (for most queries)
- **ShortTimeout**: 10 seconds (for simple lookups)
- **LongTimeout**: 60 seconds (for complex aggregations)

## Priority Order

### Phase 1 - Critical Repositories (High Traffic)
1. ✅ Lead Repository
2. ✅ Deal Repository
3. ✅ Account Repository
4. ✅ User Repository
5. ✅ Visit Report Repository

### Phase 2 - Important Repositories
6. Task Repository
7. Activity Repository
8. Pipeline Repository
9. Contact Repository

### Phase 3 - Other Repositories
10. All remaining repositories

## Notes

- All queries should use timeout to prevent indefinite hangs
- Use `DBWithTimeout` for simple queries
- Use `WithTimeout` for complex operations that need multiple queries
- Always use `repository.DefaultTimeout` unless specific timeout is needed

