# Business - Schedule Management

## CRM Healthcare Mobile App - Flutter

**Module**: Business Domain  
**Sprint**: Sprint 5  
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

Fitur **Schedule Management** memungkinkan sales rep untuk melihat dan mengelola jadwal kunjungan mereka dalam format kalender. Schedule terintegrasi dengan visit reports dan tasks untuk comprehensive time management.

### Goals

- **Calendar View**: Visual schedule dalam format kalender
- **Visit Scheduling**: Plan dan manage visit appointments
- **Conflict Detection**: Hindari schedule conflicts
- **Reminders**: Notifikasi untuk upcoming schedules
- **Integration**: Sync dengan visit reports dan tasks

---

## Fitur Utama

### 1. Calendar Views

**View Modes**:

- **Month View**: Overview bulanan
- **Week View**: Detail mingguan
- **Day View**: Detail harian dengan timeline
- **Agenda View**: List semua events

### 2. Schedule Events

**Event Types**:

- 📍 Scheduled Visit (dengan account)
- 📋 Task Due
- 🎯 Follow-up reminder
- 📝 Meeting
- 🚗 Travel time

### 3. Event Management

**Features**:

- Create new schedule
- Edit existing schedule
- Delete/cancel schedule
- Recurring schedules
- Color coding by type

### 4. Schedule Integration

**Integrations**:

- Visit reports (scheduled visits)
- Tasks (due dates)
- Route optimization
- Reminders
- Google Calendar (two-way sync)

#### Google Calendar Sync

**Overview**: Schedule dapat di-sync dengan Google Calendar untuk integrasi yang lebih baik dengan ekosistem Google user.

**Features**:

- **Bidirectional Sync**: Changes di CRM otomatis sync ke Google Calendar, dan sebaliknya
- **Selective Sync**: User dapat memilih schedule mana yang di-sync
- **Real-time**: Auto-sync ketika schedule created/updated/deleted
- **Conflict Resolution**: Handle conflicts dengan priority-based rules

**Requirements**:

- User harus connect Google Calendar terlebih dahulu (via Profile Settings)
- OAuth authorization dengan scope `calendar.events`
- Deep link handling untuk OAuth callback

**Implementation Details**:

- See [Google Calendar Documentation](../google-calendar/README.md)
- OAuth Flow: Option 2 (Direct Deep Link)
- Deep Link: `crmhealth://google-calendar/callback`

**UI Integration**:

- Schedule Form: Checkbox "Sync to Google Calendar" (only visible if connected)
- Schedule Detail: Toggle sync status dengan button Sync/Unsync
- Profile Screen: Google Calendar connection widget

**API Endpoints**:

- `POST /api/v1/schedules/:id/sync-google-calendar` - Sync schedule
- `POST /api/v1/schedules/:id/unsync-google-calendar` - Unsync schedule

---

## Business Rules

### 1. Schedule Types

**Visit Schedule**:

- Terkait dengan account
- Duration: 30-120 minutes
- Include travel time buffer
- Check-in required

**Task Schedule**:

- Deadline atau reminder
- Can be all-day atau specific time
- Link ke task detail

**Personal Schedule**:

- Meeting, training, leave
- Private events
- No business logic

### 2. Time Management

**Working Hours**:

- Default: 08:00 - 17:00
- Configurable per user
- Exclude weekends (optional)

**Buffer Time**:

- 15 minutes between visits
- Travel time calculation
- Lunch break: 12:00 - 13:00

### 3. Conflict Rules

**Overlap Detection**:

- Warning untuk overlapping schedules
- Block jika same time slot
- Consider travel time

**Resolution**:

- Suggest alternative times
- Reschedule conflicting events
- Priority-based conflict handling

### 4. Recurring Schedules

**Patterns**:

- Daily
- Weekly (specific days)
- Monthly (specific date)
- Custom interval

**End Conditions**:

- After N occurrences
- End date
- Never (ongoing)

---

## Keputusan Teknis & Trade-offs

### Calendar Library Selection

**Keputusan**: Menggunakan `table_calendar` package.

**Alasan**:

- Highly customizable
- Active maintenance
- Good performance
- Support multiple views

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── schedules/
│       ├── data/
│       │   ├── models/
│       │   │   ├── schedule_event_model.dart
│       │   │   └── schedule_filter_model.dart
│       │   └── schedule_repository.dart
│       ├── application/
│       │   ├── calendar_provider.dart
│       │   ├── schedule_list_provider.dart
│       │   └── schedule_form_provider.dart
│       └── presentation/
│           ├── screens/
│           │   ├── calendar_screen.dart
│           │   ├── schedule_detail_screen.dart
│           │   └── schedule_form_screen.dart
│           └── widgets/
│               ├── calendar_widget.dart
│               ├── event_card.dart
│               ├── day_timeline.dart
│               └── conflict_warning.dart
```

---

## API Endpoints

#### GET /api/v1/schedules

Get schedules untuk date range.

**Query Parameters**:

```
?start_date=2025-01-01&end_date=2025-01-31&type=visit
```

**Response**:

```json
{
  "success": true,
  "data": {
    "events": [
      {
        "id": "schedule-uuid",
        "title": "Visit RS Medika",
        "type": "visit",
        "start_time": "2025-01-20T09:00:00Z",
        "end_time": "2025-01-20T10:30:00Z",
        "account_id": "account-uuid",
        "account_name": "RS Medika Hospital",
        "location": {
          "address": "Jl. Sudirman No. 123",
          "latitude": -6.2088,
          "longitude": 106.8456
        },
        "description": "Monthly routine visit",
        "is_recurring": false,
        "color": "#2196F3",
        "created_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

#### POST /api/v1/schedules

Create new schedule.

**Request**:

```json
{
  "title": "Follow-up Visit",
  "type": "visit",
  "start_time": "2025-01-21T14:00:00Z",
  "end_time": "2025-01-21T15:00:00Z",
  "account_id": "account-uuid",
  "description": "Discuss contract terms",
  "is_recurring": false
}
```

#### PUT /api/v1/schedules/:id

Update schedule.

#### DELETE /api/v1/schedules/:id

Delete schedule.

#### POST /api/v1/schedules/:id/sync-google-calendar

Sync schedule to Google Calendar.

**Requirements**: User must be connected to Google Calendar.

**Response**:

```json
{
  "success": true,
  "data": {
    "synced": true,
    "google_event_id": "google-event-uuid",
    "schedule": {
      "id": "schedule-uuid",
      "title": "Visit RS Medika",
      "sync_to_calendar": true
    }
  }
}
```

#### POST /api/v1/schedules/:id/unsync-google-calendar

Unsync schedule from Google Calendar.

**Response**:

```json
{
  "success": true,
  "data": {
    "synced": false,
    "schedule": {
      "id": "schedule-uuid",
      "title": "Visit RS Medika",
      "sync_to_calendar": false
    }
  }
}
```

#### GET /api/v1/schedules/conflicts

Check schedule conflicts.

**Query Parameters**:

```
?start_time=2025-01-20T09:00:00Z&end_time=2025-01-20T10:30:00Z
```

**Response**:

```json
{
  "success": true,
  "data": {
    "has_conflicts": true,
    "conflicting_events": [
      {
        "id": "existing-schedule",
        "title": "Team Meeting",
        "start_time": "2025-01-20T09:30:00Z",
        "end_time": "2025-01-20T10:00:00Z"
      }
    ]
  }
}
```

---

## Cara Test Manual

1. **View Calendar**: Switch antara month/week/day views
2. **Create Schedule**: Add new visit schedule
3. **Check Conflicts**: Create overlapping schedule, verifikasi warning
4. **Edit Schedule**: Change time, verifikasi update
5. **Recurring Schedule**: Create weekly recurring schedule
6. **Sync Integration**: Verifikasi tasks muncul di calendar

---

## Dependencies

### Internal Dependencies

- `accounts` - Account reference untuk visit schedules
- `visit_reports` - Integration dengan visit workflow
- `tasks` - Task due dates muncul di calendar
- `google_calendar` - Google Calendar sync (see [Google Calendar Docs](../google-calendar/README.md))

### External Dependencies

- `table_calendar` - Calendar widget
- `intl` - Date formatting

### Related Documentation

- [Google Calendar Integration](../google-calendar/README.md)
- [Profile Management](./05-profile.md) - Google Calendar connection

---

**Document Status**: Active  
**Last Updated**: March 2026
