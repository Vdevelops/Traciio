# User Flow Documentation

## CRM Healthcare Mobile App - Flutter

**Module**: User Experience  
**Sprint**: All Sprints  
**Version**: 1.1  
**Status**: ✅ **Completed**  
**Last Updated**: March 2026

---

## Table of Contents

1. [Overview](#overview)
2. [User Personas](#user-personas)
3. [Onboarding Flow](#onboarding-flow)
4. [Authentication Flows](#authentication-flows)
5. [Main Navigation Flow](#main-navigation-flow)
6. [Feature User Flows](#feature-user-flows)
7. [Error Handling Flows](#error-handling-flows)
8. [Offline Mode Flows](#offline-mode-flows)
9. [Notification Flows](#notification-flows)
10. [Flow Diagrams](#flow-diagrams)

---

## Overview

Dokumen ini menjelaskan user flows lengkap untuk CRM Healthcare Mobile App, mencakup semua scenarios dari login sampai logout, termasuk error handling dan offline mode.

### Goals

- **Clear Navigation**: User dapat dengan mudah navigate antar features
- **Task Completion**: Minimal steps untuk complete tasks
- **Error Recovery**: Easy recovery dari errors
- **Offline Support**: Seamless offline experience

---

## User Personas

### 1. Sales Representative (Primary User)

**Characteristics**:

- Mobile-first user
- Needs quick access saat di lapangan
- Priority: Visit reports, accounts, tasks
- Tech-savvy level: Intermediate

**Goals**:

- Create visit reports dengan cepat
- Access account information offline
- Track daily tasks
- View performance dashboard

### 2. Supervisor

**Characteristics**:

- Manages team performance
- Needs overview dan approval workflows
- Priority: Team metrics, visit approvals
- Tech-savvy level: Intermediate

**Goals**:

- Approve/reject visit reports
- View team performance
- Monitor pipeline
- Assign tasks

---

## Onboarding Flow

### First Time Launch

```
App Launch
    │
    ▼
Splash Screen (2-3 seconds)
    │
    ├────────────────────────┐
    │                        │
    ▼                        ▼
First Time?             Returning User
    │                        │
    ▼                        ▼
Onboarding Screens      Login Screen
(3-4 slides)                │
    │                    Enter Credentials
    ▼                        │
Get Started ───────────────► Dashboard
    │
    ▼
Login Screen
```

### Onboarding Screens

**Screen 1: Welcome**

- App logo dan tagline
- "Selamat datang di CRM Healthcare"
- Brief value proposition

**Screen 2: Key Features**

- Visit Reports dengan GPS
- Account Management
- Task Tracking
- Schedule Planning

**Screen 3: Offline Support**

- "Kerja bahkan saat offline"
- Auto-sync saat online
- Data selalu tersedia

**Screen 4: Get Started**

- CTA button: "Mulai Sekarang"
- Login prompt

---

## Authentication Flows

### 1. Login Flow (Happy Path)

```
Login Screen
    │
    ├── Input Email/Username
    ├── Input Password
    │
    ▼
Validate Input
    │
    ├────────────────────────┐
    │ Valid                  │ Invalid
    ▼                        ▼
API Call: POST /auth/login  Show Error
    │                        │
    ├────────────────────────┤
    │ Success                │ Failed
    ▼                        ▼
Store Tokens              Show Error Message
    │                        │
    ▼                        │
Navigate to Dashboard ◄───────┘
    │
    ▼
Fetch Initial Data
    ├─ User Profile
    ├─ Permissions
    └─ Dashboard Overview
```

**Screen Details**:

**Login Screen**:

- Email/Username field
- Password field (obscured)
- Login button (disabled sampai valid)
- "Forgot Password?" link
- "Remember Me" checkbox

**Validation**:

- Email: Valid format
- Password: Min 6 characters
- Real-time validation dengan debounce

### 2. Login Flow (Error Scenarios)

#### Invalid Credentials

```
User Input
    │
    ▼
API Call
    │
    ▼
401 Unauthorized
    │
    ▼
Display Error:
"Email atau password salah"
(Bilingual: ID/EN)
    │
    ▼
Clear Password Field
Keep Email Field
Focus ke Password
```

#### Network Error

```
User Input
    │
    ▼
API Call
    │
    ▼
Connection Timeout
    │
    ▼
Display Error:
"Tidak dapat terhubung ke server.
Periksa koneksi internet Anda."
    │
    ▼
Show Retry Button
Keep Input Values
```

### 3. Token Refresh Flow

```
User Action (e.g., View Accounts)
    │
    ▼
API Call dengan Token
    │
    ▼
401 Token Expired
    │
    ▼
Auto Refresh Token
    ├─ POST /auth/refresh
    ├─ Store New Token
    │
    ├────────────────────────┐
    │ Success                │ Failed
    ▼                        ▼
Retry Original Call      Navigate to Login
    │                        │
    ▼                        ▼
Show Data              Clear Session
```

### 4. Logout Flow (Updated)

**Clean Logout Flow** (single navigation path):

```
User taps Logout (Profile Screen)
    │
    ▼
Confirmation Dialog
"Apakah Anda yakin ingin logout?"
    │
    ├──────────┬──────────┐
    │ Cancel   │          │ Logout
    ▼          │          ▼
 Dismiss      │    AuthNotifier.logout()
                │     ├─ Clear tokens
                │     └─ Clear user data
                │          │
                │          ▼
                │    Single Navigation:
                │    pushNamedAndRemoveUntil
                │    (clear entire nav stack)
                │          │
                └───── Login Screen
```

**Why Single Navigation Path?**

✅ **Dihapus**: Multiple navigation sources yang bersamaan  
✅ **Dihapus**: Auth state guard di profile screen `build()`  
✅ **Dihapus**: Inline `LoginScreen` render dari `AuthGate`  
✅ **Sekarang**: Satu clean navigation dari logout button handler

**Problem Sebelumnya**:

```
Logout triggered
    │
    ├─ Profile logout handler → navigate to login
    ├─ Profile build guard → navigate to login
    └─ Dashboard's AuthGate → render LoginScreen inline
    │
    ▼
Result: GLITCH ⚠️
- Loading spinner + inline login + navigation = conflict
- UI tidak responsif
- Navigation stack corrupt
```

**Solution Sekarang**:

```
Logout triggered
    │
    └─ Profile logout handler → navigate to login (SAJA)
    │
    ▼
Result: SMOOTH ✅
- Single navigation action
- Profile screen pop langsung
- AuthGate menunggu dengan loading indicator
- Clean navigation stack
```

**Interceptor-Triggered Logout** (saat token refresh gagal):

```
API Call returns 401
    │
    ▼
Auto Refresh Token
    │
    ▼
Refresh Failed
    │
    ▼
AuthGate detects unauthenticated
    │
    ├─ Set flag: _isNavigatingToLogin = true
    ├─ Show loading indicator (bukan login screen!)
    └─ Schedule navigation via addPostFrameCallback
         │
         ▼
    pushNamedAndRemoveUntil(AppRoutes.login)
         │
         ▼
    Reset flag after 500ms
         │
         ▼
    Login Screen
```

**AuthGate Behavior** (di semua protected routes):

```
AuthGate rebuild (auth state = unauthenticated)
    │
    ▼
Check _isNavigatingToLogin flag
    │
    ├─ true ────► Show loading spinner
    │              (tunggu navigasi selesai)
    │
    └─ false ───► Trigger navigation
                   ├─ Set flag = true
                   ├─ Schedule pushNamedAndRemoveUntil
                   └─ Show loading spinner
```

---

## Main Navigation Flow

### Bottom Navigation Structure

```
┌─────────────────────────────────────────┐
│                                         │
│           Screen Content                │
│                                         │
├─────────────────────────────────────────┤
│  🏠    🏢    📋    📍    👤           │
│ Home  Acc.  Tasks  Visits Profile       │
└─────────────────────────────────────────┘
```

### Navigation Logic

**Tab Selection**:

- Tap tab → Navigate ke screen
- Current tab highlighted
- Preserve scroll position
- Back button behavior:
  - Exit app dari Dashboard (first tab)
  - Navigate to previous tab dari other tabs

**Protected Routes**:

- All tabs require authentication
- Auto-redirect ke Login jika token invalid
- Restore previous route setelah login

### Deep Linking Flow

```
External Link / Push Notification
    │
    ▼
Parse URL/Deep Link
    │
    ├─ crmhealth://accounts/123
    ├─ crmhealth://tasks/456
    └─ crmhealth://visit-reports/789
    │
    ▼
Auth Check
    │
    ├──────────┬──────────┐
    │ Auth     │          │ Not Auth
    ▼          │          ▼
Navigate    │     Navigate to Login
ke Screen   │          │
    │          │          ▼
    │          └──── After Login
    │                     │
    │                     ▼
    └────────────── Navigate to
                      Original Screen
```

---

## Feature User Flows

### 1. Create Visit Report (Complete Flow)

```
User
    │
    ▼
Start from: Dashboard / Accounts / Quick Action
    │
    ▼
Visit Report Form Screen
    │
    ├─ Step 1: Select Account (jika belum dipilih)
    │     ├─ Search Account
    │     ├─ Select dari list
    │     └─ Account detail loaded
    │
    ├─ Step 2: Check-in dengan GPS
    │     ├─ Tap "Check-in"
    │     ├─ Get GPS Location
    │     ├─ Validate Accuracy (< 50m)
    │     └─ Timestamp recorded
    │
    ├─ Step 3: Add Photos
    │     ├─ Tap "Add Photo"
    │     ├─ Camera / Gallery
    │     ├─ Preview Photo
    │     ├─ Compress (auto)
    │     └─ Add to list (max 5)
    │
    ├─ Step 4: Notes
    │     └─ Type notes
    │
    ├─ Step 5: Check-out dengan GPS
    │     ├─ Tap "Check-out"
    │     ├─ Get GPS Location
    │     ├─ Calculate Duration
    │     └─ Timestamp recorded
    │
    ▼
Submit
    │
    ├──────────────────────────────┐
    │ Online                         │ Offline
    ▼                                ▼
API Call: POST /visit-reports   Save as Draft
    │                                │
    ▼                                ▼
Success                          Queue for Sync
    │                                │
    ├─ Show Success Message          ├─ Show "Saved Offline"
    ├─ Navigate to Detail            ├─ Navigate to List
    └─ Trigger Refresh               └─ Auto-sync saat online
```

**Screen Transitions**:

1. Dashboard → Account Selection (slide right)
2. Account Selection → Visit Form (fade)
3. Visit Form → Photo Preview (modal)
4. Submit Success → Detail Screen (replace)

### 2. Task Management Flow

#### Create Task

```
User
    │
    ▼
Dashboard / Tasks Tab → FAB (+)
    │
    ▼
Task Form
    │
    ├─ Title (required)
    ├─ Description (optional)
    ├─ Due Date & Time (required)
    ├─ Priority (default: Medium)
    ├─ Assigned To (default: Self)
    ├─ Related Account (optional)
    └─ Reminder (optional)
    │
    ▼
Validate
    │
    ▼
Save Task
    │
    ├─ Create local notification (jika reminder set)
    ├─ API Call (online)
    └─ Save to queue (offline)
    │
    ▼
Navigate to Task Detail
```

#### Complete Task

```
Tasks List / Task Detail
    │
    ▼
Tap Complete Button
    │
    ▼
Confirmation Dialog
"Tandai task ini sebagai selesai?"
    │
    ├──────────┬──────────┐
    │ Cancel   │          │ Yes
    ▼          │          ▼
Dismiss   │     Update Status
    │          │     ├─ status: completed
    │          │     ├─ completed_at: now
    │          │     └─ completion_notes (optional)
    │          │
    │          ▼
    │     API Call
    │          │
    │          ▼
    └──── Show Success
              │
              ▼
        Update UI
        ├─ Remove from list (filter)
        ├─ Update stats
        └─ Cancel notification
```

### 3. Account Lookup Flow

```
User needs account info
    │
    ▼
Navigate to Accounts Tab
    │
    ▼
Accounts List
    │
    ├─ Scroll through list
    ├─ Pull-to-refresh
    │
    ├─ OR Search:
    │   ├─ Tap Search Icon
    │   ├─ Type query (debounce 300ms)
    │   └─ Results filtered
    │
    ▼
Tap Account
    │
    ▼
Account Detail Screen
    │
    ├─ View Basic Info
    ├─ View Contacts
    ├─ Recent Visits
    │
    ├─ Actions:
    │   ├─ Create Visit Report
    │   ├─ Create Task
    │   ├─ Call (jika phone ada)
    │   └─ Email (jika email ada)
    │
    ▼
Tap Contact
    │
    ▼
Contact Detail
```

### 4. Dashboard Interaction Flow

```
User opens App
    │
    ▼
Dashboard (Default Screen)
    │
    ├─ View Stats Cards
    │   ├─ Total Visits
    │   ├─ Pending Tasks
    │   ├─ Active Accounts
    │   └─ Revenue
    │
    ├─ Tap Stats Card
    │   └─ Navigate to relevant screen
    │       ├─ Visits → Visit Reports List
    │       ├─ Tasks → Tasks List (filtered)
    │       └─ Accounts → Accounts List
    │
    ├─ View Pipeline Summary
    │   └─ Tap Pipeline Widget
    │       └─ Navigate to Pipeline Screen
    │
    ├─ View Upcoming Tasks
    │   └─ Tap Task
    │       └─ Navigate to Task Detail
    │
    ├─ View Recent Activities
    │   └─ Tap Activity
    │       └─ Navigate to Detail
    │
    ├─ Pull-to-Refresh
    │   └─ Refresh all data
    │
    └─ Period Selector
        ├─ Today
        ├─ This Week
        ├─ This Month
        └─ Change period → Reload data
```

### 5. Schedule Management Flow

```
User
    │
    ▼
Schedule Tab / Calendar Icon
    │
    ▼
Calendar View (Month/Week/Day)
    │
    ├─ Swipe to change month/week
    ├─ Tap date untuk Day view
    │
    ▼
Select Date
    │
    ├─ View Events on Date
    │   ├─ Visits
    │   ├─ Tasks due
    │   └─ Personal events
    │
    ├─ Tap Event → Event Detail
    │
    └─ Tap FAB (+) → Create Event
        │
        ▼
    Event Form
        ├─ Title
        ├─ Date & Time
        ├─ Duration
        ├─ Account (jika visit)
        ├─ Recurring (optional)
        └─ Reminder
        │
        ▼
    Check Conflicts
        │
        ├─ Conflict Detected
        │   └─ Show Warning
        │   └─ Suggest alternatives
        │
        └─ No Conflict
            └─ Save Event
```

---

## Error Handling Flows

### 1. Network Error Recovery

```
User Action (e.g., Load Accounts)
    │
    ▼
Network Error
    │
    ├──────────────────────────────┐
    │ Cached Data Available?       │
    ▼                              │
Show Cached Data              │
    │                              │
    ├─ Show Offline Indicator      │
    └─ Offer Refresh Button        │
    │                              │
    │                       No Cached
    │                       Data
    │                              ▼
    │                   Show Error Screen
    │                   ├─ Error Message
    │                   ├─ Retry Button
    │                   └─ Go Back Option
    │
User taps Retry
    │
    ▼
Reload Data
    │
    ├─ Success → Show Data
    └─ Fail → Show Error Again
```

### 2. Permission Denied

```
User attempts Action
(e.g., Access Location)
    │
    ▼
Permission Check
    │
    ├─ Granted → Continue
    │
    └─ Denied
        │
        ▼
Show Permission Rationale
"Aplikasi memerlukan akses lokasi
untuk check-in visit"
    │
    ├─ Deny → Disable Feature
    │
    └─ Allow → Request Permission
        │
        ├─ Granted → Continue
        │
        └─ Denied Permanently
            │
            ▼
        Show Settings Dialog
        "Izinkan akses di Settings"
            │
            ├─ Cancel → Disable Feature
            └─ Settings → Open App Settings
```

### 3. Form Validation Errors

```
User Submit Form
    │
    ▼
Validate Fields
    │
    ├─ All Valid → Submit
    │
    └─ Validation Errors
        │
        ▼
    Highlight Errors
        ├─ Red border on fields
        ├─ Error text below fields
        └─ Scroll to first error
        │
        ▼
    User Corrects
        │
        ▼
    Re-validate on Change
        │
        ├─ Still Invalid → Keep Error
        └─ Valid → Clear Error
        │
        ▼
    User Submit Again
        │
        ▼
    All Valid → Submit
```

---

## Offline Mode Flows

### 1. Initial Offline Load

```
User opens App (Offline)
    │
    ▼
Splash Screen
    │
    ▼
Check Connectivity
    │
    ▼
Offline Detected
    │
    ▼
Load from Cache
    │
    ├─ Dashboard (last cached)
    ├─ Accounts (cached list)
    └─ Tasks (cached list)
    │
    ▼
Show Offline Banner
"Mode Offline - Data mungkin tidak terbaru"
    │
    ▼
User Interacts
    ├─ View Data (read-only)
    ├─ Create Draft (queue for sync)
    └─ Navigate between screens
```

### 2. Online to Offline Transition

```
User Online
    │
    ▼
Connection Lost
    │
    ▼
Show Warning Banner
"Koneksi terputus. Mode offline aktif."
    │
    ▼
Continue with Cached Data
    │
    ├─ Ongoing operations queue
    ├─ New operations queue
    └─ Changes queue
    │
    ▼
Connection Restored
    │
    ▼
Auto-Sync Triggered
    ├─ Upload queue
    ├─ Download updates
    └─ Refresh UI
    │
    ▼
Show Success
"Data tersinkronisasi"
```

### 3. Sync Queue Management

```
User creates Visit Report (Offline)
    │
    ▼
Save as Draft
    ├─ Store locally (Hive)
    ├─ Add to sync queue
    └─ Show "Pending Sync" badge
    │
    ▼
Multiple Offline Actions
    ├─ Create Visit Report #2
    ├─ Complete Task #1
    └─ Update Account Note
    │
    ▼
All Queued Locally
    │
    ▼
Connection Restored
    │
    ▼
Process Sync Queue
    ├─ Sequential processing
    ├─ Show progress indicator
    ├─ Handle conflicts (server wins)
    └─ Retry failed items (3x)
    │
    ▼
Sync Complete
    ├─ Clear queue
    ├─ Update UI
    └─ Show notification
```

---

## Notification Flows

### 1. Push Notification Received

```
Notification Arrives
    │
    ├─ App Foreground
    │   ├─ Show in-app banner
    │   ├─ Play sound (if enabled)
    │   └─ Update badge count
    │
    ├─ App Background
    │   ├─ Show system notification
    │   ├─ Update badge count
    │   └─ Store in notification center
    │
    └─ App Terminated
        ├─ Show system notification
        ├─ Update badge count
        └─ Store for later viewing
    │
User Taps Notification
    │
    ▼
Deep Link Navigation
    ├─ Task Reminder → Task Detail
    ├─ Visit Approved → Visit Detail
    └─ General → Notification Center
```

### 2. Local Reminder Trigger

```
Reminder Time Reached
    │
    ├─ App Foreground
    │   ├─ Show alert dialog
    │   ├─ Play reminder sound
    │   └─ Offer actions:
    │       ├─ View Task
    │       ├─ Mark Complete
    │       └─ Snooze (15min/1hr/1day)
    │
    └─ App Background/Closed
        ├─ Show local notification
        ├─ Badge update
        └─ Actions in notification:
            ├─ View
            ├─ Complete
            └─ Snooze
```

---

## Flow Diagrams

### Complete App Flow Map

```
                    ┌─────────────────────────────────────┐
                    │           APP LAUNCH                │
                    └───────────────┬─────────────────────┘
                                    │
                    ┌───────────────┴───────────────┐
                    │                               │
           First Time?                    Returning User
                    │                               │
                    ▼                               ▼
            Onboarding Screens                 Login Screen
                    │                               │
                    ▼                               │
            Login Screen ◄──────────────────────────┘
                    │
                    ▼
            ┌──────────────────────────────────────┐
            │           DASHBOARD                  │
            │  (Main Navigation Hub)               │
            └──────────────┬───────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
   ┌──────────┐    ┌──────────┐    ┌──────────┐
   │ Accounts │    │  Tasks   │    │ Visits   │
   └────┬─────┘    └────┬─────┘    └────┬─────┘
        │               │               │
   ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
   │Detail   │    │Detail   │    │Create   │
   │Contact  │    │Complete │    │List     │
   │Visit    │    │         │    │Detail   │
   └─────────┘    └─────────┘    └─────────┘
        │               │               │
        └───────────────┼───────────────┘
                        │
                        ▼
               ┌──────────────────┐
               │     PROFILE      │
               │  Settings/Logout │
               └──────────────────┘
```

### Feature Decision Tree

```
User wants to:
    │
    ├─ Record Visit?
    │   └─ Visit Report Flow
    │      ├─ Check-in (GPS)
    │      ├─ Add Photos
    │      ├─ Add Notes
    │      ├─ Check-out (GPS)
    │      └─ Submit
    │
    ├─ Manage Tasks?
    │   └─ Task Flow
    │      ├─ Create
    │      ├─ View
    │      ├─ Complete
    │      └─ Set Reminder
    │
    ├─ View Accounts?
    │   └─ Account Flow
    │      ├─ Search/List
    │      ├─ View Detail
    │      ├─ View Contacts
    │      └─ Create Related Task/Visit
    │
    ├─ Check Schedule?
    │   └─ Schedule Flow
    │      ├─ Calendar View
    │      ├─ Day Detail
    │      └─ Add Event
    │
    └─ View Performance?
        └─ Dashboard Flow
           ├─ View Stats
           ├─ View Pipeline
           └─ View Activities
```

---

## User Flow Summary by Role

### Sales Representative Daily Flow

```
08:00 - Open App
        ├─ View Dashboard
        └─ Check Today's Tasks

08:15 - Start Route
        ├─ View Schedule
        └─ Open Route Optimization

09:00 - Visit Account #1
        ├─ Navigate to Account
        ├─ Check-in (GPS)
        ├─ Meeting
        ├─ Take Photos
        ├─ Check-out
        └─ Submit Visit Report

10:30 - Travel
        ├─ Next Account
        └─ Mark Task Complete (follow-up)

[... Repeat for each visit ...]

17:00 - End Day
        ├─ View Dashboard (today's stats)
        ├─ Check pending tasks
        ├─ Plan tomorrow
        └─ Logout
```

### Supervisor Daily Flow

```
09:00 - Open App
        ├─ View Team Dashboard
        └─ Check pending approvals

09:15 - Approve Visit Reports
        ├─ Review submitted reports
        ├─ View photos & notes
        ├─ Approve/Reject dengan comment
        └─ Notifications sent

10:00 - Monitor Pipeline
        ├─ View Pipeline Summary
        ├─ Check deals per stage
        └─ Follow up on stuck deals

14:00 - Assign Tasks
        ├─ Create tasks untuk team
        ├─ Set priorities
        └─ Assign to sales reps

16:00 - Review Performance
        ├─ View team metrics
        ├─ Compare performance
        └─ Identify coaching needs
```

---

**Document Status**: Active  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)
