# Analisis: "Add Visit" Button & Form Date/Timestamp Handling

**Tanggal Analisis**: May 22, 2026  
**Scope**: Web App (`/apps/web/src/features/`)  
**Keywords Dicari**: addVisit, add-visit, createVisit, VisitForm, visit modal, visit dialog

---

## 📌 RINGKASAN EKSEKUTIF

Ada **3 lokasi utama** yang menampilkan tombol "Add Visit" dan membuka dialog/form untuk membuat visit report baru. Semua menggunakan **field `visit_date` dengan format `YYYY-MM-DD HH:mm`**, bukan timestamp Unix.

---

## 📍 LOKASI FILE YANG MENANGANI "ADD VISIT"

### 1️⃣ **LeadDetailTabs.tsx** — Activity Log Tab
**File**: `apps/web/src/features/sales-crm/lead-management/components/LeadDetailTabs.tsx`

| Aspek | Detail |
|-------|--------|
| **Line Number** | 203-204 (Button) |
| **Dialog/Modal Type** | `Dialog` (shadcn/ui) |
| **Button Label** | `Add Visit` |
| **Button Location** | "Activity Log" tab, top-right corner |
| **Dialog Title** | "Create Visit Report" (Line 316) |
| **Form Component** | `VisitReportForm` (Line 319) |
| **Date Field Usage** | `visit_date` (YYYY-MM-DD HH:mm format) |
| **Initial Lead ID** | Passed via `initialLeadId={leadId}` (Line 336) |

**Excerpt Code** (Line 203-204):
```tsx
<Button size="sm" variant="outline" className="rounded-full px-4" onClick={() => setIsVisitModalOpen(true)}>
  <Plus className="h-4 w-4 mr-1" /> Add Visit
</Button>
```

**Dialog Opening** (Line 313-327):
```tsx
<Dialog open={isVisitModalOpen} onOpenChange={setIsVisitModalOpen}>
  <DialogContent className="sm:max-w-xl max-h-[90vh] overflow-y-auto">
    <DialogHeader>
      <DialogTitle>Create Visit Report</DialogTitle>
    </DialogHeader>
    <VisitReportForm
      onSubmit={handleCreateVisit}
      onCancel={() => setIsVisitModalOpen(false)}
      isLoading={createVisitReport.isPending}
      initialLeadId={leadId}
    />
  </DialogContent>
</Dialog>
```

---

### 2️⃣ **visit-report-calendar.tsx** — Calendar View
**File**: `apps/web/src/features/sales-crm/visit-report/components/visit-report-calendar.tsx`

| Aspek | Detail |
|-------|--------|
| **Line Number** | 389-400 (Dialog Definition) |
| **Dialog/Modal Type** | `Dialog` (shadcn/ui) |
| **Trigger** | Button in top toolbar (Plus icon) - line ~180 |
| **Dialog Title** | "createVisitReport" i18n key (Line 385) |
| **Form Component** | `VisitReportForm` (Line 393) |
| **Date Field Usage** | `visit_date` (YYYY-MM-DD HH:mm format) |
| **Pre-filled Date** | Selected date from calendar slot via `selectedDate` state |

**Dialog Code** (Line 381-400):
```tsx
{/* Create Dialog */}
<Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
  <DialogContent className="sm:max-w-[600px]">
    <DialogHeader>
      <DialogTitle>{t("createVisitReport")}</DialogTitle>
    </DialogHeader>
    <VisitReportForm
      key="create-visit-report-form"
      open={isCreateDialogOpen}
      onSubmit={async (data) => {
        await handleCreate(data as CreateVisitReportFormData);
      }}
      onCancel={() => {
        setIsCreateDialogOpen(false);
      }}
      isLoading={createVisitReport.isPending}
    />
  </DialogContent>
</Dialog>
```

---

### 3️⃣ **LeadDetailShell.tsx** — Quick Action Button
**File**: `apps/web/src/features/sales-crm/lead-management/components/lead-detail-shell.tsx`

| Aspek | Detail |
|-------|--------|
| **Line Number** | 270 (Button) |
| **Dialog/Modal Type** | `Dialog` (shadcn/ui) |
| **Button Label** | Implied "Add Visit" (via state variable) |
| **Button Location** | Header section of lead detail |
| **Dialog Title** | "shortcuts.newVisitReport" i18n key (Line 576) |
| **Form Component** | `VisitReportForm` (Line 579) |
| **Date Field Usage** | `visit_date` (YYYY-MM-DD HH:mm format) |

**Button** (Line 270-271):
```tsx
<Button variant="outline" size="sm" className="cursor-pointer" onClick={() => setIsVisitModalOpen(true)}>
  {t("shortcuts.newVisitReport")}
</Button>
```

**Dialog** (Line 504-521):
```tsx
<Dialog open={isVisitModalOpen} onOpenChange={setIsVisitModalOpen}>
  <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
    <DialogHeader>
      <DialogTitle>{t("shortcuts.newVisitReport")}</DialogTitle>
    </DialogHeader>
    <VisitReportForm
      onCancel={() => setIsVisitModalOpen(false)}
      onSubmit={handleCreateVisit}
      isLoading={createVisitReport.isPending}
    />
  </DialogContent>
</Dialog>
```

---

## 📋 FORM COMPONENT YANG MENANGANI DATE PICKER

### **visit-report-form.tsx** — Main Form Component
**File**: `apps/web/src/features/sales-crm/visit-report/components/visit-report-form.tsx`

| Aspek | Detail |
|-------|--------|
| **Date Field Name** | `visit_date` |
| **Date Field Type** | String (YYYY-MM-DD HH:mm) |
| **Date Picker Component** | `DateTimePicker` (Line 613) |
| **Import** | `import { DateTimePicker } from "@/components/ui/date-time-picker"` (Line 28) |
| **Handler Function** | `handleDateTimeChange()` (Line 305-318) |

#### Date Picker Implementation (Line 613-618):
```tsx
<Field orientation="vertical">
  <FieldLabel>
    {t("fields.visitDateTimeLabel")} *
  </FieldLabel>
  <DateTimePicker
    date={selectedDate}
    time={selectedTime}
    onDateChange={handleDateTimeChange}
  />
  {errors.visit_date && <FieldError>{errors.visit_date.message}</FieldError>}
</Field>
```

#### Date Change Handler (Line 305-318):
```tsx
const handleDateTimeChange = (date: Date | null, time: string | null) => {
  // Always set a visit_date string when a date is chosen.
  // If time is missing, fall back to the currently selected time or a sensible default.
  if (date) {
    const timeToUse = time ?? selectedTime ?? "09:00";
    const dateStr = `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, "0")}-${date.getDate().toString().padStart(2, "0")}`;
    const datetimeStr = `${dateStr} ${timeToUse}`;
    setValue("visit_date", datetimeStr, { shouldValidate: true });
    setSelectedDate(date);
    setSelectedTime(timeToUse);
  } else {
    // Clear selection
    setValue("visit_date", "", { shouldValidate: true });
    setSelectedDate(null);
    setSelectedTime(null);
  }
};
```

---

## 🎯 SCHEMA VALIDATION

### **visit-report.schema.ts** — Zod Schema Definition
**File**: `apps/web/src/features/sales-crm/visit-report/schemas/visit-report.schema.ts`

```typescript
export const createVisitReportSchema = z.object({
  account_id: z.string().uuid("Invalid account ID").optional(),
  contact_id: z.string().uuid("Invalid contact ID").optional(),
  deal_id: z.string().uuid("Invalid deal ID").optional(),
  lead_id: z.string().uuid("Invalid lead ID").optional(),
  visit_date: z.string().regex(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/, "Invalid date format (YYYY-MM-DD HH:mm)"),
  purpose: z.string().min(3, "Purpose must be at least 3 characters"),
  notes: z.string().optional(),
  check_in_location: locationSchema.optional(),
  check_out_location: locationSchema.optional(),
  photos: z.array(z.string().url("Invalid photo URL")).optional(),
});
```

**Format Requirement**:
- ✅ `visit_date: "YYYY-MM-DD HH:mm"` (e.g., "2025-05-22 14:30")
- ❌ NOT Unix timestamp
- ❌ NOT ISO datetime

---

## 📊 TABEL RINGKASAN LENGKAP

| Komponen | File | Line | Type | Date Field | Format | Component |
|----------|------|------|------|-----------|--------|-----------|
| **"Add Visit" Button** | LeadDetailTabs.tsx | 203 | Dialog | visit_date | YYYY-MM-DD HH:mm | VisitReportForm |
| **"Add Visit" Button** | visit-report-calendar.tsx | 389 | Dialog | visit_date | YYYY-MM-DD HH:mm | VisitReportForm |
| **"Add Visit" Button** | LeadDetailShell.tsx | 270 | Dialog | visit_date | YYYY-MM-DD HH:mm | VisitReportForm |
| **Date Picker** | visit-report-form.tsx | 613 | DateTimePicker | visit_date | YYYY-MM-DD HH:mm | DateTimePicker |
| **Date Handler** | visit-report-form.tsx | 305 | Function | visit_date | String format | handleDateTimeChange |
| **Schema** | visit-report.schema.ts | 15 | Zod | visit_date | Regex validated | createVisitReportSchema |

---

## 🔄 FLOW DIAGRAM: "Add Visit" Button → Form → Submit

```
┌─────────────────────────────────────┐
│ LeadDetailTabs.tsx (Line 203)        │
│ OR                                  │
│ visit-report-calendar.tsx (Line 381)│
│ OR                                  │
│ LeadDetailShell.tsx (Line 270)       │
│                                     │
│ → Click "Add Visit" Button          │
└────────────┬────────────────────────┘
             │
             ↓
┌─────────────────────────────────────┐
│ Dialog Opens                        │
│ (shadcn/ui Dialog)                  │
└────────────┬────────────────────────┘
             │
             ↓
┌─────────────────────────────────────┐
│ VisitReportForm renders             │
│ (visit-report-form.tsx)             │
│                                     │
│ Fields:                             │
│ - Account/Contact/Deal/Lead Select  │
│ - Purpose (textarea)                │
│ - Notes (textarea)                  │
│ - DateTimePicker (Line 613)         │
└────────────┬────────────────────────┘
             │
             ↓
┌─────────────────────────────────────┐
│ DateTimePicker Component            │
│ (Line 613)                          │
│                                     │
│ Props:                              │
│ - date: selectedDate (Date object)  │
│ - time: selectedTime (string HH:mm) │
│ - onDateChange: handler             │
│                                     │
│ → User selects date & time          │
└────────────┬────────────────────────┘
             │
             ↓
┌─────────────────────────────────────┐
│ handleDateTimeChange triggered      │
│ (Line 305-318)                      │
│                                     │
│ Converts Date → "YYYY-MM-DD HH:mm"  │
│ Updates visit_date field via        │
│ setValue("visit_date", datetimeStr) │
└────────────┬────────────────────────┘
             │
             ↓
┌─────────────────────────────────────┐
│ Form Validation                     │
│ (Zod schema - visit-report.schema)  │
│                                     │
│ Regex check:                        │
│ /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/  │
│ ✅ "2025-05-22 14:30" → PASS       │
│ ❌ Timestamp → FAIL                 │
└────────────┬────────────────────────┘
             │
             ↓
┌─────────────────────────────────────┐
│ Form Submit                         │
│ (handleFormSubmit in form)          │
│                                     │
│ Calls onSubmit() from parent        │
│ handleCreateVisit()                 │
└────────────┬────────────────────────┘
             │
             ↓
┌─────────────────────────────────────┐
│ API Call                            │
│ (visitReportService.create())       │
│                                     │
│ Sends visit_date in format:         │
│ "YYYY-MM-DD HH:mm"                  │
└─────────────────────────────────────┘
```

---

## 🔍 DETAIL PARSING & NORMALISASI

### Default Value Parsing (visit-report-form.tsx, Line 57-64):

```typescript
const parseVisitDate = (dateString: string): { date: Date; time: string | null } => {
  if (!dateString) {
    const now = new Date();
    return { date: now, time: `${now.getHours().toString().padStart(2, "0")}:${now.getMinutes().toString().padStart(2, "0")}` };
  }
  // Try to handle several common formats:
  // - YYYY-MM-DD HH:mm
  // - YYYY-MM-DDTHH:mm[:ss] (ISO)
  // - YYYY-MM-DD
  const isoOrDatetime = dateString.match(/^(\d{4}-\d{2}-\d{2})(?:[ T](\d{2}:\d{2}))?/);
  if (isoOrDatetime) {
    const [, datePart, timePart] = isoOrDatetime;
    const [year, month, day] = datePart.split("-").map(Number);
    const timeToUse = timePart ?? "09:00";
    const [hours, minutes] = timeToUse.split(":").map(Number);
    const date = new Date(year, month - 1, day, hours, minutes);
    if (!Number.isNaN(date.getTime())) {
      return { date, time: timePart ?? "09:00" };
    }
  }
  // ... fallback to Date parser
};
```

### Normalisasi String (visit-report-form.tsx, Line 142-154):

```typescript
const normalizeVisitDateString = (raw: string, fallbackTime: string | null = "09:00"): string => {
  const fb = fallbackTime ?? "09:00";
  if (!raw) return `${new Date().toISOString().split("T")[0]} ${fb}`;
  // Already in correct format
  if (/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/.test(raw)) return raw;
  // Parse any format and reformat
  const parsed = parseVisitDate(raw);
  if (parsed.date && !Number.isNaN(parsed.date.getTime())) {
    const d = parsed.date;
    const datePart = `${d.getFullYear()}-${(d.getMonth() + 1).toString().padStart(2, "0")}-${d.getDate().toString().padStart(2, "0")}`;
    const timePart = parsed.time ?? fb;
    return `${datePart} ${timePart}`;
  }
  return `${new Date().toISOString().split("T")[0]} ${fb}`;
};
```

---

## ⚙️ COMPARISON: visit_date vs TIMESTAMP

| Aspek | Visit Date (Current) | Timestamp (Activity) |
|-------|---------------------|---------------------|
| **Format** | `YYYY-MM-DD HH:mm` | ISO String or Unix |
| **Example** | `"2025-05-22 14:30"` | `"2025-05-22T14:30:00Z"` |
| **Used In** | VisitReportForm | CreateActivityDialog |
| **Component** | DateTimePicker | Input / Picker |
| **Validation** | Regex in schema | ISO validation |
| **File Location** | visit-report-form.tsx | create-activity-dialog.tsx |
| **Line Number** | 613 | ~61 |

---

## 🔗 RELATED FILES & HOOKS

| File | Purpose | Line |
|------|---------|------|
| `visit-report.schema.ts` | Zod validation schema | 15 |
| `visitReportService.ts` | API service layer | 37 |
| `useVisitReports.ts` | Custom hook for CRUD | 50 |
| `useVisitReportList.ts` | Hook for list operations | 55 |
| `visit-report-detail-modal.tsx` | Detail view modal | 27 |
| `create-activity-dialog.tsx` | Activity creation | 43 |

---

## 📝 CATATAN PENTING

1. **Format Konsisten**: Semua "Add Visit" button menggunakan `visit_date` format `YYYY-MM-DD HH:mm`
2. **Bukan Unix Timestamp**: Form menerima string format, BUKAN numeric timestamp
3. **Parser Fleksibel**: Form dapat parse berbagai format date input dari API dan normalize menjadi format standar
4. **Dialog Type**: Semua menggunakan `Dialog` dari shadcn/ui, bukan modal custom
5. **Date Picker Component**: Custom `DateTimePicker` component yang menangani date dan time secara terpisah
6. **Initializer**: LeadDetailTabs dapat pre-fill dengan `initialLeadId` untuk Lead-scoped creation

---

## 🎯 KESIMPULAN

**SEMUA "Add Visit" button dan form menggunakan:**
- ✅ Field: `visit_date` (STRING, bukan timestamp)
- ✅ Format: `YYYY-MM-DD HH:mm`
- ✅ Component: `DateTimePicker`
- ✅ Dialog Type: `Dialog` (shadcn/ui)
- ✅ Validasi: Regex `/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/`

Jika Anda ingin mengubah ke Unix timestamp atau format lain, perlu update:
1. Schema regex di `visit-report.schema.ts`
2. Handler di `visit-report-form.tsx`
3. All 3 places yang membuka dialog
