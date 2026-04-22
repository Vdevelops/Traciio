# Product Requirements Document (PRD)
## CRM Healthcare/Pharmaceutical Platform - Sales CRM

**Versi**: 2.1  
**Status**: Active  
**Last Updated**: 2025-01-20  
**Target Release**: MVP Q1 2025  
**Product Type**: Sales CRM untuk Perusahaan Farmasi

---

## 📋 Daftar Isi

1. [Executive Summary](#executive-summary)
2. [Product Overview](#product-overview)
3. [Target Users](#target-users)
4. [Business Objectives](#business-objectives)
5. [Functional Requirements](#functional-requirements)
6. [Non-Functional Requirements](#non-functional-requirements)
7. [User Stories](#user-stories)
8. [Technical Requirements](#technical-requirements)
9. [Success Metrics](#success-metrics)
10. [Timeline & Milestones](#timeline--milestones)
11. [Risks & Mitigation](#risks--mitigation)

---

## Executive Summary

CRM Healthcare/Pharmaceutical Platform adalah sistem manajemen hubungan pelanggan (CRM) yang dirancang khusus untuk **perusahaan farmasi** yang melakukan sales ke rumah sakit, klinik, dan apotek. Platform ini membantu sales team mengelola leads, accounts, contacts, visit reports, sales pipeline, dan aktivitas sales secara terintegrasi melalui **Web Application** dan **Mobile App (Flutter)**.

### Key Value Propositions

- **Lead Management**: Manajemen sales leads dari berbagai sumber dengan qualification dan conversion ke opportunities
- **Account & Contact Management**: Manajemen data rumah sakit, klinik, apotek, dan kontak (dokter, PIC) yang terpusat
- **Visit Report & Activity Tracking**: Tracking kunjungan sales dengan check-in/out, GPS, dan dokumentasi foto (dapat dikaitkan ke Lead, Deal, atau Account)
- **Sales Pipeline Management**: Manage sales pipeline dari opportunity (deal) dengan forecast yang akurat. Deals dibuat melalui Lead Conversion, bukan langsung di pipeline
- **Task & Reminder**: Task management dan reminder otomatis untuk follow-up
- **Mobile-First**: Mobile app (Flutter) untuk sales rep bekerja di lapangan
- **On-Premise Ready**: Dapat di-install on-premise di server perusahaan

---

## Product Overview

### Vision

Menjadi platform CRM Sales terdepan untuk perusahaan farmasi di Indonesia yang memungkinkan sales team bekerja lebih efisien dan efektif dalam mengelola hubungan dengan rumah sakit, klinik, dan apotek.

### Mission

Menyediakan solusi CRM Sales yang komprehensif, mudah digunakan, dan mobile-first untuk membantu perusahaan farmasi mengoptimalkan aktivitas sales dan meningkatkan revenue.

### Product Goals

1. **Improve Sales Efficiency**: Mempermudah sales rep dalam mengelola kunjungan dan aktivitas sales
2. **Better Visibility**: Memberikan visibility yang jelas kepada supervisor dan management tentang aktivitas sales
3. **Data-Driven Decisions**: Memberikan insights yang actionable melalui data dan analitik sales
4. **Mobile-First**: Mobile app yang memungkinkan sales rep bekerja di lapangan dengan mudah
5. **Scalability**: Mendukung pertumbuhan dari tim sales kecil hingga enterprise

---

## Target Users

### Primary Users

1. **Sales Rep / Sales Representative**
   - Melakukan kunjungan ke rumah sakit, klinik, apotek
   - Membuat visit report dengan check-in/out
   - Mengelola tasks dan follow-up
   - Melihat accounts dan contacts
   - Menggunakan mobile app di lapangan

2. **Supervisor / Sales Manager**
   - Review dan approve visit reports
   - Monitor aktivitas sales team
   - Manage sales pipeline
   - Melihat dashboard dan laporan
   - Assign tasks ke sales rep

3. **Administrator**
   - Mengelola users dan permissions
   - Mengelola accounts dan contacts
   - Konfigurasi sistem
   - Melihat laporan dan analitik
   - Manage products dan pricing

### Secondary Users

1. **Manager / Owner**
   - Melihat dashboard dan laporan
   - Analisis sales performance
   - Pengambilan keputusan strategis
   - Forecast dan planning

2. **Viewer / Read-Only User**
   - Melihat data dan laporan (tanpa edit)
   - Export data untuk analisis

---

## Business Objectives

### Short-term (MVP - 100 hari)

1. **Core Functionality**: Implementasi fitur-fitur inti untuk sales operations
2. **User Adoption**: Onboarding minimal 1 perusahaan farmasi dengan 10+ sales rep
3. **Data Accuracy**: Akurasi data > 95%
4. **System Uptime**: Uptime > 99%
5. **Mobile App**: Functional Flutter app untuk sales rep
6. **On-Premise**: Ready untuk install on-premise

### Medium-term (6-12 bulan)

1. **Market Expansion**: Ekspansi ke 10+ perusahaan farmasi
2. **Feature Enhancement**: Penambahan fitur advanced berdasarkan feedback
3. **Integration**: Integrasi dengan sistem eksternal (ERP, accounting)
4. **Advanced Analytics**: Advanced analytics dan reporting
5. **Mobile Enhancement**: Offline support, advanced features di mobile app

### Long-term (12+ bulan)

1. **Market Leadership**: Menjadi market leader CRM Sales untuk pharma di Indonesia
2. **AI/ML Integration**: Implementasi AI untuk sales prediction dan recommendation
3. **Advanced Features**: Territory management, route optimization, digital signature
4. **Analytics Platform**: Platform analitik advanced untuk sales insights

---

## Functional Requirements

### 1. Authentication & Authorization

#### 1.1 User Authentication
- **FR-1.1.1**: Sistem harus mendukung login dengan email dan password
- **FR-1.1.2**: Sistem harus mendukung multi-factor authentication (MFA) - optional untuk MVP
- **FR-1.1.3**: Sistem harus mendukung password reset via email
- **FR-1.1.4**: Sistem harus mendukung session management dengan token-based authentication (JWT)
- **FR-1.1.5**: Sistem harus mendukung remember me functionality
- **FR-1.1.6**: Mobile app harus mendukung login dan token refresh

#### 1.2 Role-Based Access Control (RBAC)
- **FR-1.2.1**: Sistem harus mendukung multiple roles (Admin, Supervisor, Sales, Viewer)
- **FR-1.2.2**: Sistem harus mendukung permission-based access control
- **FR-1.2.3**: Sistem harus mendukung role assignment per user
- **FR-1.2.4**: Sistem harus mendukung audit log untuk semua akses

### 2. Account & Contact Management

#### 2.1 Account Management
- **FR-2.1.1**: Sistem harus memungkinkan registrasi account baru (Rumah Sakit, Klinik, Apotek)
- **FR-2.1.2**: Sistem harus menyimpan informasi account lengkap (nama, alamat, kategori, status)
- **FR-2.1.3**: Sistem harus mendukung multiple accounts per sales rep
- **FR-2.1.4**: Sistem harus mendukung import account dari Excel (optional untuk MVP)
- **FR-2.1.5**: Sistem harus mendukung kategori account (Hospital, Clinic, Pharmacy)

#### 2.2 Contact Management
- **FR-2.2.1**: Sistem harus menyimpan informasi kontak (Dokter, PIC, Manager) yang terhubung ke account
- **FR-2.2.2**: Sistem harus menyimpan informasi kontak lengkap (nama, phone, email, jabatan)
- **FR-2.2.3**: Sistem harus mendukung multiple contacts per account
- **FR-2.2.4**: Sistem harus mendukung riwayat interaksi dengan contact
- **FR-2.2.5**: Sistem harus mendukung kategori contact (Doctor, PIC, Manager, Other)

#### 2.3 Account & Contact Search & Filter
- **FR-2.3.1**: Sistem harus memungkinkan pencarian account berdasarkan nama, kategori, lokasi
- **FR-2.3.2**: Sistem harus memungkinkan pencarian contact berdasarkan nama, phone, email
- **FR-2.3.3**: Sistem harus memungkinkan filter berdasarkan status, kategori, sales rep assignment
- **FR-2.3.4**: Sistem harus mendukung pagination untuk daftar account dan contact

### 3. Lead Management

#### 3.1 Lead Creation & Management
- **FR-3.1.1**: Sistem harus memungkinkan create lead baru dari berbagai sumber (website, referral, cold call, event, social media, email campaign, partner, other)
- **FR-3.1.2**: Sistem harus menyimpan informasi lead lengkap (nama, company, email, phone, lead source, lead status, lead score)
- **FR-3.1.3**: Sistem harus mendukung lead status flow: new → contacted → qualified → converted/lost (atau unqualified → nurturing/disqualified)
- **FR-3.1.4**: Sistem harus mendukung lead qualification process dengan lead scoring
- **FR-3.1.5**: Sistem harus memungkinkan assign lead ke sales rep
- **FR-3.1.6**: Sistem harus mendukung search dan filter leads berdasarkan status, source, assigned to

#### 3.2 Lead Conversion
- **FR-3.2.1**: Sistem harus memungkinkan convert qualified lead menjadi opportunity (deal)
- **FR-3.2.2**: Saat convert, sistem harus otomatis create account dari lead (jika belum ada) dan create contact (opsional)
- **FR-3.2.3**: Sistem harus otomatis migrate visit reports dan activities dari lead ke account/deal yang baru dibuat
- **FR-3.2.4**: Lead status harus berubah menjadi "converted" setelah conversion
- **FR-3.2.5**: Sistem harus memungkinkan pre-convert account creation (create account dari lead sebelum conversion)

#### 3.3 Lead Analytics
- **FR-3.3.1**: Sistem harus menampilkan lead analytics (total leads, leads by status, leads by source, conversion rate)
- **FR-3.3.2**: Sistem harus menampilkan conversion rate per lead source
- **FR-3.3.3**: Sistem harus menampilkan average time to conversion

### 4. Visit Report & Activity Tracking

#### 4.1 Visit Report Creation
- **FR-4.1.1**: Sistem harus memungkinkan sales rep membuat visit report dengan pilihan basis: Lead, Deal, atau Account (via tabs)
- **FR-4.1.2**: Sistem harus mendukung check-in/check-out dengan GPS location
- **FR-4.1.3**: Sistem harus menyimpan informasi kunjungan (account, contact, deal, lead, purpose, notes)
- **FR-4.1.4**: Business rule: Visit report harus dikaitkan ke minimal Lead ID atau Account ID. Jika Deal ID disediakan, Account ID wajib ada
- **FR-4.1.5**: Sistem harus mendukung upload foto dokumentasi kunjungan
- **FR-4.1.6**: Sistem harus mendukung multiple visit types (regular, follow-up, emergency)
- **FR-4.1.7**: Mobile app harus mendukung create visit report dengan GPS dan foto

#### 4.2 Visit Report Management
- **FR-4.2.1**: Sistem harus menampilkan list visit reports dengan filter (account, deal, lead, sales rep, date range)
- **FR-4.2.2**: Sistem harus mendukung status tracking (draft, submitted, approved, rejected)
- **FR-4.2.3**: Supervisor harus dapat review dan approve/reject visit reports
- **FR-4.2.4**: Sistem harus menampilkan activity timeline untuk setiap account, deal, atau lead
- **FR-4.2.5**: Sistem harus mendukung search visit reports
- **FR-4.2.6**: Sistem harus menampilkan visit reports terkait untuk deal tertentu (via GET /deals/:id/visit-reports)
- **FR-4.2.7**: Sistem harus menampilkan visit reports terkait untuk lead tertentu (via GET /leads/:id/visit-reports)

#### 4.3 Activity Tracking
- **FR-4.3.1**: Sistem harus melacak semua aktivitas sales (visit, call, email, task) yang dapat dikaitkan ke Lead, Account, atau Deal
- **FR-4.3.2**: Business rule: Activity harus dikaitkan ke minimal Lead ID, Account ID, atau Deal ID
- **FR-4.3.3**: Sistem harus menampilkan activity timeline per account, deal, atau lead
- **FR-4.3.4**: Sistem harus mendukung filter activity berdasarkan type, date, sales rep, account, deal, lead
- **FR-4.3.5**: Mobile app harus menampilkan activity timeline
- **FR-4.3.6**: Sistem harus menampilkan activities terkait untuk deal tertentu (via GET /deals/:id/activities)
- **FR-4.3.7**: Sistem harus menampilkan activities terkait untuk lead tertentu (via GET /leads/:id/activities)

### 5. Sales Pipeline Management

#### 5.1 Pipeline Stages
- **FR-5.1.1**: Sistem harus mendukung multiple pipeline stages (Qualification, Proposal, Negotiation, Closed Won, Closed Lost). Note: Stage "Lead" sudah dihapus karena deals masuk pipeline dari "Qualification" stage setelah Lead Conversion
- **FR-5.1.2**: Sistem harus memungkinkan custom pipeline stages
- **FR-5.1.3**: Sistem harus menampilkan kanban view untuk pipeline
- **FR-5.1.4**: Sistem harus mendukung drag-and-drop untuk move deal antar stages

#### 5.2 Deal/Opportunity Management
- **FR-5.2.1**: Sistem harus memungkinkan create deal/opportunity melalui Lead Conversion (POST /leads/:id/convert), bukan langsung di pipeline
- **FR-5.2.2**: Sistem harus menyimpan informasi deal (account, contact, product, value, expected close date, probability, stage)
- **FR-5.2.3**: Sistem harus mendukung notes dan attachments untuk deal
- **FR-5.2.4**: Sistem harus mendukung forecast berdasarkan deals
- **FR-5.2.5**: Sistem harus mendukung win/loss tracking
- **FR-5.2.6**: Deal harus dapat dikaitkan ke lead source (lead_id) untuk tracking
- **FR-5.2.7**: Deal detail harus menampilkan visit reports dan activities terkait

#### 5.3 Pipeline Analytics
- **FR-5.3.1**: Sistem harus menampilkan pipeline summary (total value, stage distribution)
- **FR-5.3.2**: Sistem harus menampilkan forecast berdasarkan deals
- **FR-5.3.3**: Sistem harus menampilkan conversion rate per stage
- **FR-5.3.4**: Sistem harus mendukung date range filtering

### 6. Task & Reminder Management

#### 6.1 Task Creation
- **FR-6.1.1**: Sistem harus memungkinkan create task untuk follow-up
- **FR-6.1.2**: Sistem harus menyimpan informasi task (title, description, due date, assignee, priority)
- **FR-6.1.3**: Sistem harus mendukung task assignment ke sales rep
- **FR-6.1.4**: Sistem harus mendukung task linked ke account, contact, deal, atau lead
- **FR-6.1.5**: Mobile app harus mendukung create dan view tasks

#### 6.2 Task Management
- **FR-6.2.1**: Sistem harus menampilkan task list dengan filter (status, assignee, due date, account, deal, lead)
- **FR-6.2.2**: Sistem harus mendukung task status (open, in-progress, done, cancelled)
- **FR-6.2.3**: Sistem harus mendukung task reminder (email, in-app notification)
- **FR-6.2.4**: Mobile app harus mengirim push notification untuk task reminder

### 7. Product Management

#### 7.1 Product Catalog
- **FR-7.1.1**: Sistem harus memungkinkan manage product catalog
- **FR-7.1.2**: Sistem harus menyimpan informasi product (nama, SKU, kategori, harga)
- **FR-7.1.3**: Sistem harus mendukung product categories
- **FR-7.1.4**: Sistem harus mendukung search products
- **FR-7.1.5**: Sistem harus mendukung product linked ke deals

### 8. Dashboard & Reports

#### 8.1 Dashboard
- **FR-8.1.1**: Sistem harus menampilkan dashboard dengan key metrics
- **FR-8.1.2**: Sistem harus menampilkan visit statistics (today, this week, this month)
- **FR-8.1.3**: Sistem harus menampilkan pipeline summary
- **FR-8.1.4**: Sistem harus menampilkan lead analytics (total leads, conversion rate, leads by status)
- **FR-8.1.5**: Sistem harus menampilkan top accounts, top sales rep
- **FR-8.1.6**: Sistem harus menampilkan recent activities
- **FR-8.1.7**: Mobile app harus menampilkan dashboard (basic)

#### 8.2 Reports
- **FR-8.2.1**: Sistem harus menghasilkan laporan visit reports (daily, weekly, monthly)
- **FR-8.2.2**: Sistem harus menghasilkan laporan sales pipeline
- **FR-8.2.3**: Sistem harus menghasilkan laporan sales performance per sales rep
- **FR-8.2.4**: Sistem harus menghasilkan laporan account activity
- **FR-8.2.5**: Sistem harus menghasilkan laporan lead analytics dan conversion
- **FR-8.2.6**: Sistem harus mendukung export ke PDF, Excel

### 9. Mobile App (Flutter)

#### 9.1 Core Mobile Features
- **FR-9.1.1**: Mobile app harus mendukung login dan authentication
- **FR-9.1.2**: Mobile app harus menampilkan dashboard (basic metrics)
- **FR-9.1.3**: Mobile app harus menampilkan visit reports list dan create visit report (dapat dikaitkan ke Lead, Deal, atau Account)
- **FR-9.1.4**: Mobile app harus mendukung check-in/out dengan GPS
- **FR-9.1.5**: Mobile app harus menampilkan tasks list dan create task
- **FR-9.1.6**: Mobile app harus menampilkan accounts dan contacts list
- **FR-9.1.7**: Mobile app harus mendukung upload foto untuk visit report
- **FR-9.1.8**: Mobile app harus mengirim push notification untuk task reminder
- **FR-9.1.9**: Mobile app harus menampilkan leads list (view only untuk MVP)

#### 9.2 Mobile Offline Support (Optional untuk MVP)
- **FR-9.2.1**: Mobile app harus mendukung offline mode (sync later)
- **FR-9.2.2**: Mobile app harus sync data saat online

### 10. System Administration

#### 10.1 User Management
- **FR-10.1.1**: Sistem harus memungkinkan admin mengelola users
- **FR-10.1.2**: Sistem harus memungkinkan assign roles dan permissions
- **FR-10.1.3**: Sistem harus memungkinkan enable/disable users
- **FR-10.1.4**: Sistem harus menyimpan audit log untuk semua user actions

#### 10.2 Settings
- **FR-10.2.1**: Sistem harus memungkinkan konfigurasi sistem (company info, logo, branding)
- **FR-10.2.2**: Sistem harus memungkinkan konfigurasi notification settings
- **FR-10.2.3**: Sistem harus memungkinkan konfigurasi pipeline stages (tanpa stage "Lead")
- **FR-10.2.4**: Sistem harus memungkinkan konfigurasi account categories

### 11. On-Premise Installation

#### 11.1 Installation Requirements
- **FR-11.1.1**: Sistem harus dapat di-install on-premise di server perusahaan
- **FR-11.1.2**: Sistem harus menyediakan Docker Compose untuk easy installation
- **FR-11.1.3**: Sistem harus menyediakan installation guide (PDF)
- **FR-11.1.4**: Sistem harus menyediakan database migration scripts
- **FR-11.1.5**: Sistem harus menyediakan environment configuration template

---

## Non-Functional Requirements

### 1. Performance
- **NFR-1.1**: Sistem harus dapat menangani minimal 100 concurrent users
- **NFR-1.2**: Response time untuk CRUD operations < 500ms
- **NFR-1.3**: Response time untuk reports < 3 seconds
- **NFR-1.4**: Sistem harus mendukung pagination untuk semua list views

### 2. Security
- **NFR-2.1**: Semua data harus dienkripsi dalam transit (HTTPS)
- **NFR-2.2**: Data sensitif harus dienkripsi di rest (encryption at rest)
- **NFR-2.3**: Sistem harus mematuhi HIPAA/standar privasi data kesehatan Indonesia
- **NFR-2.4**: Sistem harus mendukung audit logging untuk compliance
- **NFR-2.5**: Sistem harus mendukung data backup dan recovery

### 3. Availability
- **NFR-3.1**: Sistem harus memiliki uptime > 99.5%
- **NFR-3.2**: Sistem harus mendukung maintenance mode dengan graceful degradation
- **NFR-3.3**: Sistem harus memiliki disaster recovery plan

### 4. Scalability
- **NFR-4.1**: Sistem harus dapat scale horizontal
- **NFR-4.2**: Database harus dapat handle minimal 1M records per table
- **NFR-4.3**: Sistem harus mendukung multi-tenant architecture

### 5. Usability
- **NFR-5.1**: UI harus responsive (mobile, tablet, desktop)
- **NFR-5.2**: Sistem harus mendukung bilingual (Bahasa Indonesia & English)
- **NFR-5.3**: Sistem harus memiliki intuitive navigation
- **NFR-5.4**: Sistem harus mendukung keyboard shortcuts untuk power users

### 6. Compatibility
- **NFR-6.1**: Sistem harus support modern browsers (Chrome, Firefox, Safari, Edge)
- **NFR-6.2**: Sistem harus support mobile browsers
- **NFR-6.3**: Sistem harus support printing untuk receipts, invoices, prescriptions

---

## User Stories

### Epic 1: Account & Contact Management

**US-1.1**: Sebagai admin, saya ingin dapat mendaftarkan account baru (RS/Klinik/Apotek) dengan data lengkap agar dapat mengelola customer dengan baik.

**US-1.2**: Sebagai sales rep, saya ingin dapat melihat list accounts yang ditugaskan kepada saya agar dapat merencanakan kunjungan.

**US-1.3**: Sebagai admin, saya ingin dapat mencari account berdasarkan berbagai kriteria agar dapat menemukan data account dengan cepat.

**US-1.4**: Sebagai sales rep, saya ingin dapat melihat kontak (dokter, PIC) dari setiap account agar dapat menghubungi orang yang tepat.

### Epic 2: Lead Management

**US-2.1**: Sebagai sales rep, saya ingin dapat membuat lead baru dari berbagai sumber agar dapat melacak potential customers.

**US-2.2**: Sebagai sales rep, saya ingin dapat qualify lead dan update lead status agar dapat mengelola lead dengan efektif.

**US-2.3**: Sebagai sales rep, saya ingin dapat convert qualified lead menjadi opportunity (deal) agar dapat memindahkan lead ke sales pipeline.

**US-2.4**: Sebagai supervisor, saya ingin dapat melihat lead analytics (conversion rate, leads by source) agar dapat menganalisis efektivitas lead generation.

**US-2.5**: Sebagai sales rep, saya ingin dapat create account dari lead sebelum conversion agar dapat mempersiapkan data account lebih awal.

### Epic 3: Visit Report & Activity Tracking

**US-3.1**: Sebagai sales rep, saya ingin dapat membuat visit report dengan pilihan basis (Lead, Deal, atau Account) agar dapat mengaitkan kunjungan dengan konteks yang tepat.

**US-3.2**: Sebagai sales rep, saya ingin dapat check-in/check-out dengan GPS di mobile app agar supervisor dapat memverifikasi lokasi kunjungan.

**US-3.3**: Sebagai supervisor, saya ingin dapat review dan approve visit reports agar dapat memastikan kualitas laporan kunjungan.

**US-3.4**: Sebagai sales rep, saya ingin dapat upload foto dokumentasi kunjungan agar dapat memberikan bukti visual.

**US-3.5**: Sebagai supervisor, saya ingin dapat melihat activity timeline untuk setiap account, deal, atau lead agar dapat memahami history interaksi.

**US-3.6**: Sebagai sales rep, saya ingin dapat melihat visit reports dan activities terkait untuk deal tertentu agar dapat melacak progress deal.

### Epic 4: Sales Pipeline

**US-4.1**: Sebagai sales rep, saya ingin dapat convert qualified lead menjadi deal agar dapat memindahkan lead ke sales pipeline dengan data lengkap.

**US-4.2**: Sebagai supervisor, saya ingin dapat melihat pipeline dalam kanban view agar dapat memantau progress deals.

**US-4.3**: Sebagai manager, saya ingin dapat melihat forecast berdasarkan pipeline agar dapat melakukan planning.

**US-4.4**: Sebagai sales rep, saya ingin dapat move deal antar stages agar dapat update progress sales.

**US-4.5**: Sebagai sales rep, saya ingin dapat melihat visit reports dan activities terkait untuk deal tertentu di deal detail agar dapat memahami konteks deal.

### Epic 5: Task & Reminder

**US-4.1**: Sebagai sales rep, saya ingin dapat membuat task untuk follow-up agar tidak lupa melakukan tindak lanjut.

**US-4.2**: Sebagai supervisor, saya ingin dapat assign task ke sales rep agar dapat mengatur workload.

**US-4.3**: Sebagai sales rep, saya ingin menerima reminder untuk task yang due date-nya mendekati agar dapat menyelesaikan tepat waktu.

**US-4.4**: Sebagai sales rep, saya ingin dapat melihat task list di mobile app agar dapat mengelola tasks saat di lapangan.

### Epic 6: Dashboard & Reports

**US-5.1**: Sebagai supervisor, saya ingin dapat melihat dashboard dengan key metrics agar dapat memantau performance sales team.

**US-5.2**: Sebagai manager, saya ingin dapat melihat laporan visit reports harian/mingguan/bulanan agar dapat menganalisis aktivitas sales.

**US-5.3**: Sebagai admin, saya ingin dapat export laporan ke Excel/PDF agar dapat melakukan analisis lebih lanjut.

### Epic 7: Mobile App

**US-6.1**: Sebagai sales rep, saya ingin dapat login ke mobile app agar dapat mengakses CRM saat di lapangan.

**US-6.2**: Sebagai sales rep, saya ingin dapat membuat visit report di mobile app dengan GPS dan foto agar dapat mendokumentasikan kunjungan secara real-time.

**US-6.3**: Sebagai sales rep, saya ingin dapat melihat tasks dan accounts di mobile app agar dapat bekerja tanpa perlu membuka laptop.

---

## Technical Requirements

### Technology Stack

#### Web Frontend
- **Framework**: Next.js 16 (App Router, Server Components)
- **Styling**: Tailwind CSS v4
- **State Management**: Zustand
- **UI Components**: shadcn/ui v4
- **Internationalization**: next-intl (Bahasa Indonesia & English)

#### Mobile App
- **Framework**: Flutter
- **State Management**: 
- **HTTP Client**: 
- **Local Storage**: 

#### Backend
- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Cache**: Redis (optional untuk MVP)
- **File Storage**: Local / S3 compatible storage

#### Infrastructure
- **Hosting**: Cloud (AWS/GCP/Azure)
- **Containerization**: Docker
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus + Grafana
- **Logging**: ELK Stack

### API Standards
- Mengikuti API Response Standards yang sudah didefinisikan
- RESTful API design
- API versioning (v1, v2, dll)
- Bilingual error messages (ID & EN)

### Database Design
- Normalized database schema
- Soft delete untuk data penting
- Audit trails untuk compliance
- Indexing untuk performance

### Security
- JWT-based authentication
- Role-based access control (RBAC)
- Data encryption (at rest & in transit)
- SQL injection prevention
- XSS prevention
- CSRF protection

---

## Success Metrics

### User Adoption
- **Target**: 1 perusahaan farmasi dengan 10+ sales rep onboarded dalam 100 hari
- **Metric**: Number of active users, number of daily logins, number of visit reports created

### System Performance
- **Target**: 99.5% uptime
- **Metric**: System availability, response time, error rate

### Data Quality
- **Target**: > 95% data accuracy
- **Metric**: Data validation errors, data completeness

### User Satisfaction
- **Target**: NPS > 50
- **Metric**: User surveys, feedback, support tickets

### Business Impact
- **Target**: 30% reduction in administrative time untuk sales rep
- **Metric**: Time saved per visit report, number of visit reports per sales rep, pipeline conversion rate

---

## Timeline & Milestones

### Phase 1: MVP (100 Hari / ~14 Minggu)

#### Week 1-2: Foundation & Account Management
- Authentication & Authorization (Web + Mobile)
- User Management
- Account & Contact Management (Web)

#### Week 3-4: Visit Report & Sales Pipeline
- Visit Report & Activity Tracking (Web + Mobile)
- Sales Pipeline Management (Web)

#### Week 5-6: Task, Product, Dashboard
- Task & Reminder (Web + Mobile)
- Product Management (Web)
- Dashboard & Reports (Web)

#### Week 7-8: Mobile App & Integration
- Mobile App (Flutter) - Core Features
- Integration & Testing

#### Week 9-11: Polish & Demo
- Testing & Bug Fixes
- Demo Preparation
- On-Premise Setup

#### Week 12-14: Buffer & Delivery
- Client Feedback
- Customization
- Final Delivery

### Phase 2: Enhancement (Months 4-6)
- Advanced Reports
- Offline Support (Mobile)
- Advanced Analytics
- Performance Optimization

### Phase 3: Scale (Months 7-12)
- Multi-company Support
- Advanced Features (Territory Management, Route Optimization)
- API for Third-party Integration
- Advanced Mobile Features

---

## Risks & Mitigation

### Technical Risks

**Risk 1**: Performance issues dengan large dataset
- **Mitigation**: Implement proper indexing, pagination, caching strategy

**Risk 2**: Security vulnerabilities
- **Mitigation**: Regular security audits, penetration testing, code reviews

**Risk 3**: Integration complexity dengan BPJS
- **Mitigation**: Early research, proof of concept, phased integration

### Business Risks

**Risk 1**: Low user adoption
- **Mitigation**: User training, excellent UX, responsive support

**Risk 2**: Regulatory compliance issues
- **Mitigation**: Early consultation with legal/regulatory experts, compliance audits

**Risk 3**: Competition from established players
- **Mitigation**: Focus on unique value propositions, excellent customer service

### Operational Risks

**Risk 1**: Data loss
- **Mitigation**: Regular backups, disaster recovery plan, data redundancy

**Risk 2**: System downtime
- **Mitigation**: High availability architecture, monitoring, alerting

---

## Appendix

### Glossary

- **Account**: Rumah Sakit, Klinik, atau Apotek yang menjadi customer perusahaan farmasi
- **Contact**: Dokter, PIC, atau Manager yang merupakan kontak di account
- **Visit Report**: Laporan kunjungan sales rep ke account
- **Pipeline**: Sales pipeline dari lead hingga deal
- **Deal/Opportunity**: Potential sales yang sedang dalam proses
- **Check-in/Check-out**: Proses menandai mulai dan selesai kunjungan dengan GPS

### References

- API Response Standards: `/docs/api-standart/api-response-standards.md`
- API Error Codes: `/docs/api-standart/api-error-codes.md`
- Modules Documentation: `/docs/modules/01-modules.md`

---

**Dokumen ini akan diupdate sesuai dengan perkembangan development dan feedback dari stakeholders.**

