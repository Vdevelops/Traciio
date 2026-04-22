# Business - Leads Management

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

Fitur **Leads Management** memungkinkan sales rep untuk melihat dan mengelola leads (calon accounts) sebelum mereka di-convert menjadi accounts. Fitur ini terintegrasi dengan pipeline untuk tracking lead progression.

### Goals

- **Lead Viewing**: Melihat leads yang di-assign
- **Lead Tracking**: Update lead status dan progress
- **Conversion**: Convert leads to accounts
- **Qualification**: Assess lead quality dan potential

---

## Fitur Utama

### 1. Leads List

**Features**:

- List leads dengan status
- Filter by status (new, contacted, qualified, lost)
- Search by name/company
- Sort by date/priority

### 2. Lead Detail

**Information**:

- Company info
- Contact person
- Source (web, referral, cold call)
- Status history
- Notes dan activities

### 3. Lead Actions

**Actions**:

- Update status
- Add notes
- Schedule follow-up
- Convert to account
- Mark as lost

### 4. Lead Conversion

**Conversion Flow**:

1. Review lead info
2. Create account from lead data
3. Transfer contacts
4. Update pipeline stage
5. Archive lead

---

## Business Rules

### 1. Lead Assignment

- Leads di-assign ke sales rep berdasarkan territory/specialization
- Sales rep hanya melihat assigned leads
- Supervisor dapat reassign leads

### 2. Lead Status

**Statuses**:

- **New**: Baru di-assign, belum di-contact
- **Contacted**: Sudah initial contact
- **Qualified**: Qualified opportunity
- **Proposal**: Proposal sent
- **Negotiation**: Under negotiation
- **Converted**: Converted to account
- **Lost**: Not interested/cannot contact

### 3. Conversion Rules

- Only qualified leads dapat di-convert
- Required fields harus lengkap
- Create account record
- Transfer lead history

---

## Keputusan Teknis & Trade-offs

### Lead vs Account

**Keputusan**: Separate entity dengan conversion workflow.

**Alasan**:

- Clear distinction antara prospects dan customers
- Different workflows
- Better tracking conversion rates

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── leads/
│       ├── data/
│       │   ├── models/
│       │   │   └── lead_model.dart
│       │   └── lead_repository.dart
│       ├── application/
│       │   ├── lead_list_provider.dart
│       │   └── lead_detail_provider.dart
│       └── presentation/
│           ├── screens/
│           │   ├── lead_list_screen.dart
│           │   ├── lead_detail_screen.dart
│           │   └── lead_convert_screen.dart
│           └── widgets/
│               ├── lead_card.dart
│               ├── lead_status_badge.dart
│               └── lead_conversion_form.dart
```

---

## API Endpoints

#### GET /api/v1/leads

List leads.

**Query Parameters**:

```
?page=1&limit=20&status=new&search=keyword
```

**Response**:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "lead-uuid",
        "company_name": "RS Sejahtera",
        "contact_name": "Dr. Ahmad",
        "contact_email": "ahmad@rssejahtera.com",
        "contact_phone": "+62123456789",
        "status": "new",
        "source": "website",
        "priority": "high",
        "assigned_to": "user-uuid",
        "created_at": "2025-01-20T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 50
    }
  }
}
```

#### GET /api/v1/leads/:id

Get lead detail.

#### PUT /api/v1/leads/:id/status

Update lead status.

**Request**:

```json
{
  "status": "contacted",
  "notes": "Initial phone call completed"
}
```

#### POST /api/v1/leads/:id/convert

Convert lead to account.

**Response**:

```json
{
  "success": true,
  "data": {
    "account_id": "account-uuid",
    "message": "Lead converted successfully"
  }
}
```

---

## Cara Test Manual

1. **View Leads**: Verifikasi list leads muncul dengan benar
2. **Filter by Status**: Apply filter, verifikasi results
3. **Update Status**: Change lead status, verifikasi persist
4. **Convert Lead**: Convert qualified lead, verifikasi account created
5. **Add Notes**: Add activity notes, verifikasi tersimpan

---

**Document Status**: Active  
**Last Updated**: January 2025
