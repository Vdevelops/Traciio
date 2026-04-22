# Business - Route Optimization

## CRM Healthcare Mobile App - Flutter

**Module**: Business Domain  
**Sprint**: Sprint 5  
**Version**: 1.0  
**Status**: ⏳ **In Progress**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Fitur Utama](#fitur-utama)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Struktur Folder](#struktur-folder)
6. [API Endpoints](#api-endpoints)
7. [Data Models](#data-models)
8. [Configuration](#configuration)
9. [Usage Examples](#usage-examples)
10. [Cara Test Manual](#cara-test-manual)
11. [Dependencies](#dependencies)
12. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Fitur **Route Optimization** membantu sales rep untuk merencanakan kunjungan ke multiple accounts dengan urutan yang optimal berdasarkan jarak dan waktu tempuh. Fitur ini mengintegrasikan Google Maps untuk navigation dan ETA calculation.

### Goals

- **Optimal Route**: Generate rute terpendek/tercepat
- **Multi-stop Planning**: Plan visits ke multiple accounts
- **Time Estimation**: ETA per destination
- **Navigation**: Turn-by-turn navigation
- **Schedule Integration**: Sync dengan visit schedule

---

## Fitur Utama

### 1. Route Planner

**Input**:

- Starting point (current location atau custom)
- List accounts to visit
- Time constraints (start time, end time)
- Priority accounts

**Output**:

- Optimized route order
- Total distance dan time
- ETA per stop
- Alternative routes

### 2. Interactive Map

**Features**:

- Show all stops on map
- Route polyline
- Account markers
- Current location
- Traffic overlay (optional)

### 3. Turn-by-Turn Navigation

**Integration**:

- Open external maps app
- Google Maps directions
- Waze integration
- Custom navigation (if available)

### 4. Schedule Optimization

**Features**:

- Suggest optimal visit times
- Consider account availability
- Lunch/travel breaks
- Buffer time

---

## Business Rules

### 1. Route Constraints

**Optimization Criteria**:

- Shortest total distance
- Shortest total time
- Priority accounts first
- Time window constraints

### 2. Account Selection

- Max 10 accounts per route
- Filter by territory
- Filter by priority
- Exclude completed visits

### 3. Time Management

**Default Settings**:

- Visit duration: 30-60 minutes per account
- Travel buffer: 15 minutes
- Working hours: 08:00 - 17:00
- Lunch break: 12:00 - 13:00

### 4. Schedule Sync

- Integrasi dengan visit schedule
- Auto-add scheduled visits
- Conflict detection

---

## Keputusan Teknis & Trade-offs

### External Maps vs Custom Navigation

**Keputusan**: Use external maps apps (Google Maps, Waze) untuk navigation.

**Alasan**:

- Reliable navigation data
- Real-time traffic
- Voice guidance
- No development cost

**Trade-off**: Leave app untuk navigation. **Mitigasi**: Deep link dengan return functionality.

### Route Calculation

**Keputusan**: Calculate route client-side dengan Google Maps Distance Matrix API.

**Alasan**:

- Real-time calculation
- Consider current traffic
- Multiple optimization options

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── route_optimization/
│       ├── data/
│       │   ├── models/
│       │   │   ├── route_model.dart
│       │   │   └── route_stop_model.dart
│       │   └── route_repository.dart
│       ├── application/
│       │   ├── route_planner_provider.dart
│       │   └── navigation_provider.dart
│       └── presentation/
│           ├── screens/
│           │   ├── route_planner_screen.dart
│           │   ├── route_detail_screen.dart
│           │   └── route_map_screen.dart
│           └── widgets/
│               ├── stop_list_item.dart
│               ├── route_summary_card.dart
│               └── navigation_button.dart
├── core/
│   └── services/
│       ├── maps_service.dart
│       └── location_service.dart
```

---

## API Endpoints

#### POST /api/v1/route/optimize

Calculate optimal route.

**Request**:

```json
{
  "start_location": {
    "latitude": -6.2088,
    "longitude": 106.8456
  },
  "account_ids": ["account-1", "account-2", "account-3"],
  "constraints": {
    "max_travel_time": 480,
    "departure_time": "08:00",
    "consider_traffic": true
  }
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "route": {
      "total_distance": 25000,
      "total_duration": 180,
      "stops": [
        {
          "order": 1,
          "account_id": "account-1",
          "account_name": "RS Medika",
          "estimated_arrival": "08:30",
          "travel_time_from_previous": 30,
          "distance_from_previous": 8000,
          "location": {
            "latitude": -6.2,
            "longitude": 106.85
          }
        }
      ],
      "polyline": "encoded_polyline_string"
    }
  }
}
```

#### GET /api/v1/accounts/nearby

Get nearby accounts untuk route planning.

**Query Parameters**:

```
?lat=-6.2088&lng=106.8456&radius=10000&limit=20
```

---

## Cara Test Manual

1. **Plan Route**: Select accounts, generate route
2. **View Map**: Verifikasi stops dan route di map
3. **Start Navigation**: Tap navigate, verifikasi open maps app
4. **Reorder Stops**: Drag to reorder, verifikasi recalculation
5. **ETA Accuracy**: Compare estimated vs actual arrival times

---

**Document Status**: In Progress  
**Last Updated**: January 2025
