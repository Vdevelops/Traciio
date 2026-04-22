# Area Mapping & Location Capture API Documentation

## Overview
The Area Mapping & Location Capture module enables field teams to capture location data, define territories, and analyze coverage. This module leverages PostGIS for spatial operations and provides comprehensive location intelligence.

**Sprint**: Sprint 3  
**Status**: ✅ Backend Complete, Frontend Pending  
**Base URL**: `/api/v1/area-mapping`

---

## Table of Contents
1. [Features](#features)
2. [Database Schema](#database-schema)
3. [API Endpoints](#api-endpoints)
4. [Data Models](#data-models)
5. [Spatial Operations](#spatial-operations)
6. [Error Codes](#error-codes)
7. [Usage Examples](#usage-examples)

---

## Features

### Core Capabilities
- **Location Capture**: Record GPS coordinates with metadata (activity, customer, type)
- **Territory Management**: Define and manage geographic territories with polygon boundaries
- **Coverage Analysis**: Analyze field team coverage and territory penetration
- **Spatial Queries**: Find captures within territories, nearby captures, heatmaps
- **Filtering & Search**: Advanced filtering by date, user, territory, activity type
- **Pagination**: Efficient data retrieval with cursor-based pagination

### Spatial Operations (PostGIS)
- `ST_Within`: Check if point is within polygon territory
- `ST_Distance`: Calculate distance between locations
- `ST_MakePoint`: Create geographic points from coordinates
- `ST_GeomFromText`: Convert WKT to geography
- `ST_AsGeoJSON`: Export geometries as GeoJSON

---

## Database Schema

### Tables

#### 1. `area_captures`
Stores individual location captures by field team members.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `user_id` | UUID | Reference to users table |
| `location` | GEOGRAPHY(POINT) | GPS coordinates (longitude, latitude) |
| `accuracy` | NUMERIC(10,2) | GPS accuracy in meters |
| `captured_at` | TIMESTAMPTZ | Timestamp of capture |
| `activity_id` | UUID | Optional reference to activity |
| `lead_id` | UUID | Optional reference to lead |
| `contact_id` | UUID | Optional reference to contact |
| `account_id` | UUID | Optional reference to account |
| `capture_type` | VARCHAR(50) | Type: visit, check-in, route-point |
| `notes` | TEXT | Additional notes |
| `metadata` | JSONB | Extra metadata (e.g., battery level, network) |
| `created_at` | TIMESTAMPTZ | Record creation time |
| `updated_at` | TIMESTAMPTZ | Last update time |

**Indexes**:
- GIST spatial index on `location`
- B-tree index on `user_id`
- B-tree index on `captured_at`
- B-tree index on `capture_type`

#### 2. `territories`
Defines geographic territories as polygons.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `name` | VARCHAR(100) | Territory name |
| `code` | VARCHAR(50) | Unique territory code |
| `description` | TEXT | Description |
| `boundary` | GEOGRAPHY(POLYGON) | Territory boundary polygon |
| `area_km2` | NUMERIC(12,2) | Calculated area in km² |
| `assigned_user_ids` | UUID[] | Array of assigned user IDs |
| `color` | VARCHAR(7) | Hex color for map display |
| `is_active` | BOOLEAN | Active status |
| `metadata` | JSONB | Extra metadata |
| `created_at` | TIMESTAMPTZ | Record creation time |
| `updated_at` | TIMESTAMPTZ | Last update time |

**Indexes**:
- GIST spatial index on `boundary`
- B-tree index on `code` (unique)
- B-tree index on `is_active`

#### 3. `coverage_analysis`
Stores coverage analysis results.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `territory_id` | UUID | Reference to territories table |
| `user_id` | UUID | Optional reference to users table |
| `analysis_date` | DATE | Date of analysis |
| `period_start` | TIMESTAMPTZ | Analysis period start |
| `period_end` | TIMESTAMPTZ | Analysis period end |
| `total_captures` | INTEGER | Total captures in period |
| `unique_locations` | INTEGER | Unique locations visited |
| `coverage_percentage` | NUMERIC(5,2) | Coverage percentage |
| `areas_covered` | GEOGRAPHY(MULTIPOLYGON) | Optional covered areas |
| `metrics` | JSONB | Additional metrics |
| `created_at` | TIMESTAMPTZ | Record creation time |
| `updated_at` | TIMESTAMPTZ | Last update time |

**Indexes**:
- B-tree index on `territory_id`
- B-tree index on `analysis_date`
- B-tree index on `user_id`

---

## API Endpoints

### 1. Create Area Capture
**POST** `/area-mapping/captures`

Record a new location capture.

**Request Body**:
```json
{
  "user_id": "uuid",
  "latitude": -6.2088,
  "longitude": 106.8456,
  "accuracy": 10.5,
  "captured_at": "2025-01-27T10:30:00Z",
  "activity_id": "uuid",
  "lead_id": "uuid",
  "contact_id": "uuid",
  "account_id": "uuid",
  "capture_type": "visit",
  "notes": "Customer visit at headquarters",
  "metadata": {
    "battery_level": 85,
    "network_type": "4G"
  }
}
```

**Response** (201 Created):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "accuracy": 10.5,
    "captured_at": "2025-01-27T10:30:00Z",
    "activity_id": "uuid",
    "lead_id": "uuid",
    "contact_id": "uuid",
    "account_id": "uuid",
    "capture_type": "visit",
    "notes": "Customer visit at headquarters",
    "metadata": {
      "battery_level": 85,
      "network_type": "4G"
    },
    "created_at": "2025-01-27T10:30:00Z",
    "updated_at": "2025-01-27T10:30:00Z"
  }
}
```

**Validation**:
- Latitude: -90 to 90
- Longitude: -180 to 180
- Required fields: `user_id`, `latitude`, `longitude`, `captured_at`

---

### 2. Get Captures by User
**GET** `/area-mapping/captures/user/:user_id`

Retrieve all captures for a specific user with filtering and pagination.

**Query Parameters**:
- `start_date` (optional): Filter start date (ISO 8601)
- `end_date` (optional): Filter end date (ISO 8601)
- `capture_type` (optional): Filter by type (visit, check-in, route-point)
- `limit` (optional): Results per page (default: 50, max: 500)
- `offset` (optional): Pagination offset (default: 0)

**Example**:
```
GET /area-mapping/captures/user/550e8400-e29b-41d4-a716-446655440000?start_date=2025-01-01T00:00:00Z&limit=20
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "latitude": -6.2088,
      "longitude": 106.8456,
      "accuracy": 10.5,
      "captured_at": "2025-01-27T10:30:00Z",
      "capture_type": "visit",
      "notes": "Customer visit",
      "created_at": "2025-01-27T10:30:00Z"
    }
  ],
  "pagination": {
    "total": 150,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

---

### 3. Get Captures in Territory
**GET** `/area-mapping/captures/territory/:territory_id`

Retrieve all captures within a territory boundary.

**Query Parameters**:
- `start_date` (optional): Filter start date
- `end_date` (optional): Filter end date
- `user_id` (optional): Filter by specific user
- `limit` (optional): Results per page
- `offset` (optional): Pagination offset

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "latitude": -6.2088,
      "longitude": 106.8456,
      "captured_at": "2025-01-27T10:30:00Z",
      "capture_type": "visit",
      "territory_id": "uuid",
      "territory_name": "Jakarta Pusat"
    }
  ],
  "pagination": {
    "total": 75,
    "limit": 50,
    "offset": 0
  }
}
```

---

### 4. Find Nearby Captures
**GET** `/area-mapping/captures/nearby`

Find captures within a radius of a point.

**Query Parameters** (required):
- `latitude`: Center point latitude
- `longitude`: Center point longitude
- `radius_meters`: Search radius in meters

**Query Parameters** (optional):
- `limit`: Results per page (default: 50)
- `offset`: Pagination offset

**Example**:
```
GET /area-mapping/captures/nearby?latitude=-6.2088&longitude=106.8456&radius_meters=1000&limit=10
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "latitude": -6.2090,
      "longitude": 106.8460,
      "distance_meters": 45.2,
      "captured_at": "2025-01-27T09:15:00Z",
      "capture_type": "check-in"
    }
  ],
  "pagination": {
    "total": 8,
    "limit": 10,
    "offset": 0
  }
}
```

---

### 5. Get Capture Heatmap
**GET** `/area-mapping/heatmap`

Generate heatmap data for visualization (grid-based aggregation).

**Query Parameters**:
- `start_date` (optional): Filter start date
- `end_date` (optional): Filter end date
- `user_id` (optional): Filter by user
- `grid_size` (optional): Grid size in decimal degrees (default: 0.01)

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "latitude": -6.20,
      "longitude": 106.84,
      "count": 45,
      "intensity": 0.85
    },
    {
      "latitude": -6.21,
      "longitude": 106.85,
      "count": 32,
      "intensity": 0.60
    }
  ]
}
```

---

### 6. Create Territory
**POST** `/area-mapping/territories`

Create a new territory with polygon boundary.

**Request Body**:
```json
{
  "name": "Jakarta Pusat",
  "code": "JKT-CENTRAL",
  "description": "Central Jakarta territory",
  "boundary": {
    "type": "Polygon",
    "coordinates": [
      [
        [106.8000, -6.1500],
        [106.8500, -6.1500],
        [106.8500, -6.2000],
        [106.8000, -6.2000],
        [106.8000, -6.1500]
      ]
    ]
  },
  "assigned_user_ids": ["uuid1", "uuid2"],
  "color": "#FF5733",
  "is_active": true,
  "metadata": {
    "population": 928000,
    "priority": "high"
  }
}
```

**Response** (201 Created):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Jakarta Pusat",
    "code": "JKT-CENTRAL",
    "description": "Central Jakarta territory",
    "boundary": {
      "type": "Polygon",
      "coordinates": [...]
    },
    "area_km2": 48.13,
    "assigned_user_ids": ["uuid1", "uuid2"],
    "color": "#FF5733",
    "is_active": true,
    "metadata": {...},
    "created_at": "2025-01-27T10:00:00Z",
    "updated_at": "2025-01-27T10:00:00Z"
  }
}
```

**Validation**:
- Polygon must be closed (first point = last point)
- Minimum 4 points (3 vertices + 1 closing point)
- Code must be unique

---

### 7. Get All Territories
**GET** `/area-mapping/territories`

Retrieve all territories with optional filtering.

**Query Parameters**:
- `is_active` (optional): Filter by active status (true/false)
- `assigned_user_id` (optional): Filter territories assigned to user
- `limit` (optional): Results per page
- `offset` (optional): Pagination offset

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Jakarta Pusat",
      "code": "JKT-CENTRAL",
      "area_km2": 48.13,
      "assigned_user_ids": ["uuid1", "uuid2"],
      "is_active": true,
      "boundary": {...}
    }
  ],
  "pagination": {
    "total": 12,
    "limit": 50,
    "offset": 0
  }
}
```

---

### 8. Get Territory by ID
**GET** `/area-mapping/territories/:id`

Retrieve a single territory by ID.

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Jakarta Pusat",
    "code": "JKT-CENTRAL",
    "description": "Central Jakarta territory",
    "boundary": {...},
    "area_km2": 48.13,
    "assigned_user_ids": ["uuid1", "uuid2"],
    "color": "#FF5733",
    "is_active": true,
    "metadata": {...},
    "created_at": "2025-01-27T10:00:00Z",
    "updated_at": "2025-01-27T10:00:00Z"
  }
}
```

---

### 9. Update Territory
**PUT** `/area-mapping/territories/:id`

Update an existing territory.

**Request Body** (all fields optional):
```json
{
  "name": "Jakarta Pusat - Updated",
  "description": "Updated description",
  "boundary": {...},
  "assigned_user_ids": ["uuid1", "uuid2", "uuid3"],
  "color": "#00FF00",
  "is_active": false
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Jakarta Pusat - Updated",
    "updated_at": "2025-01-27T11:00:00Z",
    ...
  }
}
```

---

### 10. Delete Territory
**DELETE** `/area-mapping/territories/:id`

Soft delete a territory (sets `is_active` to false).

**Response** (204 No Content):
```json
{
  "success": true,
  "message": "Territory deleted successfully"
}
```

---

### 11. Create Coverage Analysis
**POST** `/area-mapping/coverage-analysis`

Generate a coverage analysis report.

**Request Body**:
```json
{
  "territory_id": "uuid",
  "user_id": "uuid",
  "period_start": "2025-01-01T00:00:00Z",
  "period_end": "2025-01-31T23:59:59Z"
}
```

**Response** (201 Created):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "territory_id": "uuid",
    "user_id": "uuid",
    "analysis_date": "2025-01-27",
    "period_start": "2025-01-01T00:00:00Z",
    "period_end": "2025-01-31T23:59:59Z",
    "total_captures": 156,
    "unique_locations": 78,
    "coverage_percentage": 68.42,
    "metrics": {
      "avg_captures_per_day": 5.03,
      "most_active_day": "2025-01-15",
      "visits": 120,
      "check_ins": 36
    },
    "created_at": "2025-01-27T12:00:00Z"
  }
}
```

---

## Data Models

### AreaCapture (Go Model)
```go
type AreaCapture struct {
    ID          uuid.UUID       `gorm:"type:uuid;primary_key"`
    UserID      uuid.UUID       `gorm:"type:uuid;not null;index"`
    Location    orb.Point       `gorm:"type:geography(POINT,4326);not null;index:,type:gist"`
    Accuracy    *float64        `gorm:"type:numeric(10,2)"`
    CapturedAt  time.Time       `gorm:"not null;index"`
    ActivityID  *uuid.UUID      `gorm:"type:uuid"`
    LeadID      *uuid.UUID      `gorm:"type:uuid"`
    ContactID   *uuid.UUID      `gorm:"type:uuid"`
    AccountID   *uuid.UUID      `gorm:"type:uuid"`
    CaptureType string          `gorm:"type:varchar(50);index"`
    Notes       string          `gorm:"type:text"`
    Metadata    datatypes.JSON  `gorm:"type:jsonb"`
    CreatedAt   time.Time       `gorm:"not null"`
    UpdatedAt   time.Time       `gorm:"not null"`
}
```

### Territory (Go Model)
```go
type Territory struct {
    ID              uuid.UUID       `gorm:"type:uuid;primary_key"`
    Name            string          `gorm:"type:varchar(100);not null"`
    Code            string          `gorm:"type:varchar(50);unique;not null;index"`
    Description     string          `gorm:"type:text"`
    Boundary        orb.Polygon     `gorm:"type:geography(POLYGON,4326);not null;index:,type:gist"`
    AreaKm2         *float64        `gorm:"type:numeric(12,2)"`
    AssignedUserIDs pq.StringArray  `gorm:"type:uuid[]"`
    Color           string          `gorm:"type:varchar(7)"`
    IsActive        bool            `gorm:"default:true;index"`
    Metadata        datatypes.JSON  `gorm:"type:jsonb"`
    CreatedAt       time.Time       `gorm:"not null"`
    UpdatedAt       time.Time       `gorm:"not null"`
}
```

### CoverageAnalysis (Go Model)
```go
type CoverageAnalysis struct {
    ID                 uuid.UUID         `gorm:"type:uuid;primary_key"`
    TerritoryID        uuid.UUID         `gorm:"type:uuid;not null;index"`
    UserID             *uuid.UUID        `gorm:"type:uuid;index"`
    AnalysisDate       time.Time         `gorm:"type:date;not null;index"`
    PeriodStart        time.Time         `gorm:"not null"`
    PeriodEnd          time.Time         `gorm:"not null"`
    TotalCaptures      int               `gorm:"default:0"`
    UniqueLocations    int               `gorm:"default:0"`
    CoveragePercentage *float64          `gorm:"type:numeric(5,2)"`
    AreasCovered       *orb.MultiPolygon `gorm:"type:geography(MULTIPOLYGON,4326)"`
    Metrics            datatypes.JSON    `gorm:"type:jsonb"`
    CreatedAt          time.Time         `gorm:"not null"`
    UpdatedAt          time.Time         `gorm:"not null"`
}
```

---

## Spatial Operations

### PostGIS Queries Used

#### 1. Find Captures Within Territory
```sql
SELECT * FROM area_captures
WHERE ST_Within(
    location,
    (SELECT boundary FROM territories WHERE id = $1)
);
```

#### 2. Find Nearby Captures
```sql
SELECT *,
    ST_Distance(location, ST_MakePoint($1, $2)::geography) as distance
FROM area_captures
WHERE ST_DWithin(
    location,
    ST_MakePoint($1, $2)::geography,
    $3
)
ORDER BY distance ASC;
```

#### 3. Calculate Territory Area
```sql
SELECT ST_Area(boundary::geography) / 1000000 as area_km2
FROM territories
WHERE id = $1;
```

#### 4. Heatmap Aggregation
```sql
SELECT 
    FLOOR(ST_X(location::geometry) / $1) * $1 as grid_lon,
    FLOOR(ST_Y(location::geometry) / $1) * $1 as grid_lat,
    COUNT(*) as count
FROM area_captures
GROUP BY grid_lon, grid_lat
ORDER BY count DESC;
```

---

## Error Codes

| Code | Message | HTTP Status |
|------|---------|-------------|
| `INVALID_COORDINATES` | Invalid latitude or longitude | 400 |
| `INVALID_POLYGON` | Polygon is not closed or has < 4 points | 400 |
| `TERRITORY_NOT_FOUND` | Territory not found | 404 |
| `CAPTURE_NOT_FOUND` | Capture not found | 404 |
| `DUPLICATE_TERRITORY_CODE` | Territory code already exists | 409 |
| `VALIDATION_ERROR` | Request validation failed | 422 |
| `SPATIAL_QUERY_ERROR` | PostGIS spatial query failed | 500 |

---

## Usage Examples

### Example 1: Record Check-in at Customer Location
```bash
curl -X POST http://localhost:8080/api/v1/area-mapping/captures \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "accuracy": 15.0,
    "captured_at": "2025-01-27T10:30:00Z",
    "account_id": "660e8400-e29b-41d4-a716-446655440000",
    "capture_type": "check-in",
    "notes": "Customer visit - Discussing Q1 targets"
  }'
```

### Example 2: Find All Captures in Territory This Week
```bash
curl -X GET "http://localhost:8080/api/v1/area-mapping/captures/territory/770e8400-e29b-41d4-a716-446655440000?start_date=2025-01-20T00:00:00Z&end_date=2025-01-27T23:59:59Z" \
  -H "Authorization: Bearer <token>"
```

### Example 3: Create Territory for Jakarta Selatan
```bash
curl -X POST http://localhost:8080/api/v1/area-mapping/territories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Jakarta Selatan",
    "code": "JKT-SOUTH",
    "description": "South Jakarta territory covering Kebayoran, Cilandak, and surrounding areas",
    "boundary": {
      "type": "Polygon",
      "coordinates": [[
        [106.7900, -6.2400],
        [106.8500, -6.2400],
        [106.8500, -6.3000],
        [106.7900, -6.3000],
        [106.7900, -6.2400]
      ]]
    },
    "assigned_user_ids": ["550e8400-e29b-41d4-a716-446655440000"],
    "color": "#3498db",
    "is_active": true
  }'
```

### Example 4: Generate Coverage Report
```bash
curl -X POST http://localhost:8080/api/v1/area-mapping/coverage-analysis \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "territory_id": "770e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "period_start": "2025-01-01T00:00:00Z",
    "period_end": "2025-01-31T23:59:59Z"
  }'
```

### Example 5: Find Captures Near My Current Location
```bash
curl -X GET "http://localhost:8080/api/v1/area-mapping/captures/nearby?latitude=-6.2088&longitude=106.8456&radius_meters=500&limit=20" \
  -H "Authorization: Bearer <token>"
```

---

## Integration with Other Modules

### Activity Module
- Link captures to activities via `activity_id`
- Auto-capture location on activity check-in

### Lead & Contact Module
- Link captures to leads/contacts for visit tracking
- View visit history on lead/contact detail page

### Route Optimization
- Use captures for route planning
- Optimize routes based on territory boundaries

### Analytics & Reporting
- Coverage dashboards
- Field team productivity metrics
- Territory performance analysis

---

## Performance Considerations

1. **Spatial Indexes**: GIST indexes on geometry columns for fast spatial queries
2. **Pagination**: Always use limit/offset for large result sets
3. **Caching**: Consider caching territory boundaries and static data
4. **Batch Inserts**: For offline sync, batch capture uploads in chunks of 50-100
5. **Heatmap Grid Size**: Larger grid sizes (e.g., 0.05°) for faster heatmap generation

---

## Security & Permissions

### Required Permissions
- `area_mapping.create`: Create captures and territories
- `area_mapping.read`: View captures and territories
- `area_mapping.update`: Update territories
- `area_mapping.delete`: Delete territories
- `area_mapping.admin`: Full administrative access

### Data Access Rules
- Users can only view their own captures (unless admin)
- Territory assignments control data visibility
- Coverage analysis respects user permissions

---

## Offline Support (Mobile)

### Sync Strategy
1. **Offline Queue**: Store captures locally when offline
2. **Auto-Sync**: Upload queued captures when connection restored
3. **Conflict Resolution**: Server timestamp wins on conflicts
4. **Territory Cache**: Cache territory boundaries for offline validation

### Offline Capabilities
- ✅ Record captures offline
- ✅ View cached territories
- ✅ Validate coordinates locally
- ❌ Real-time coverage analysis (requires server)
- ❌ Nearby captures query (requires server)

---

## Next Steps

### Backend (Completed ✅)
- [x] PostGIS migration
- [x] Domain models
- [x] Repository layer
- [x] Service layer
- [x] API handlers
- [x] Route integration
- [x] Demo seeders
- [x] AutoMigrate integration

### Frontend (Pending)
- [ ] TypeScript types
- [ ] Zod schemas
- [ ] API services
- [ ] Map components (Leaflet/MapLibre)
- [ ] Territory management UI
- [ ] Coverage dashboard
- [ ] i18n translations

### Documentation (Pending)
- [x] API documentation (this file)
- [ ] Postman collection update
- [ ] Menu seeder update
- [ ] Permission seeder update
- [ ] User guide
- [ ] Architecture diagrams

---

## Support & Contact
For issues or questions, contact the development team or refer to the main project documentation.

**Last Updated**: January 27, 2025  
**Version**: 1.0.0
