# Business - Sales Pipeline

## CRM Healthcare Mobile App - Flutter

**Module**: Business Domain  
**Sprint**: Sprint 4  
**Version**: 1.0  
**Status**: ✅ **Completed**  
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

Fitur **Sales Pipeline** memungkinkan sales rep untuk melihat dan mengelola deals/opportunities melalui sales stages. Pipeline membantu tracking deal progress dari initial contact sampai closing.

### Goals

- **Pipeline View**: Visual pipeline dengan stages
- **Deal Tracking**: Track deals per stage
- **Stage Management**: Move deals antara stages
- **Revenue Forecasting**: Pipeline value dan probability

---

## Fitur Utama

### 1. Pipeline Board

**Kanban-style Board**:

- Columns per stage
- Deal cards dengan key info
- Drag-and-drop (jika supported)
- Stage counts dan values

### 2. Deal List per Stage

**List View**:

- Deal name dan account
- Deal value
- Expected close date
- Assigned sales rep
- Priority/urgency

### 3. Deal Detail

**Information**:

- Deal overview
- Account info
- Products/services
- Deal value dan currency
- Expected close date
- Stage history
- Activities dan notes

### 4. Deal Actions

**Actions**:

- Move to next stage
- Move to previous stage
- Update deal value
- Change expected close date
- Mark as won/lost
- Add activities

---

## Business Rules

### 1. Pipeline Stages

**Standard Stages**:

1. **Prospecting** - Initial identification
2. **Qualification** - Needs assessment
3. **Proposal** - Quote submitted
4. **Negotiation** - Terms discussion
5. **Closed Won** - Deal secured
6. **Closed Lost** - Deal lost

### 2. Stage Movement

- Forward movement: Any stage ke next
- Backward movement: Any stage ke previous
- Won/Lost: Final stages
- Lost reason required

### 3. Deal Values

- Minimum deal value: configurable
- Currency: IDR default
- Probability per stage:
  - Prospecting: 10%
  - Qualification: 25%
  - Proposal: 50%
  - Negotiation: 75%
  - Closed Won: 100%
  - Closed Lost: 0%

### 4. Assignment

- Deals di-assign ke sales rep
- Sales rep hanya melihat assigned deals
- Supervisor melihat team pipeline

---

## Keputusan Teknis & Trade-offs

### Kanban vs List View

**Keputusan**: Support both views dengan Kanban sebagai default.

**Alasan**:

- Kanban: Visual pipeline progression
- List: Detailed information view
- User preference

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── pipeline/
│       ├── data/
│       │   ├── models/
│       │   │   ├── deal_model.dart
│       │   │   └── pipeline_stage_model.dart
│       │   └── pipeline_repository.dart
│       ├── application/
│       │   ├── pipeline_provider.dart
│       │   └── deal_detail_provider.dart
│       └── presentation/
│           ├── screens/
│           │   ├── pipeline_board_screen.dart
│           │   ├── deal_list_screen.dart
│           │   └── deal_detail_screen.dart
│           └── widgets/
│               ├── pipeline_column.dart
│               ├── deal_card.dart
│               └── stage_progress_indicator.dart
```

---

## API Endpoints

#### GET /api/v1/pipeline

Get pipeline data dengan deals per stage.

**Response**:

```json
{
  "success": true,
  "data": {
    "stages": [
      {
        "id": "stage-1",
        "name": "Prospecting",
        "order": 1,
        "probability": 10,
        "deals": [
          {
            "id": "deal-uuid",
            "name": "RS Medika - Equipment Supply",
            "account": {
              "id": "account-uuid",
              "name": "RS Medika Hospital"
            },
            "value": 50000000,
            "currency": "IDR",
            "expected_close_date": "2025-03-15",
            "assigned_to": "user-uuid",
            "created_at": "2025-01-10T10:00:00Z"
          }
        ],
        "total_value": 150000000,
        "deal_count": 3
      }
    ],
    "summary": {
      "total_deals": 12,
      "total_pipeline_value": 450000000,
      "weighted_value": 225000000
    }
  }
}
```

#### GET /api/v1/deals/:id

Get deal detail.

#### PUT /api/v1/deals/:id/stage

Move deal to different stage.

**Request**:

```json
{
  "stage_id": "stage-2",
  "notes": "Qualified after needs assessment"
}
```

#### POST /api/v1/deals/:id/close

Close deal (won/lost).

**Request**:

```json
{
  "status": "lost",
  "reason": "Budget constraints",
  "notes": "Customer decided to postpone"
}
```

---

## Cara Test Manual

1. **View Pipeline**: Verifikasi pipeline board muncul dengan stages
2. **View Deals**: Tap stage, verifikasi deals list
3. **Move Stage**: Change deal stage, verifikasi persist
4. **Close Deal**: Mark deal as won/lost, verifikasi final stage
5. **Deal Detail**: Tap deal card, verifikasi detail screen

---

**Document Status**: Active  
**Last Updated**: January 2025
