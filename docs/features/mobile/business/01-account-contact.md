# Business - Account & Contact Management

## CRM Healthcare Mobile App - Flutter

**Module**: Business Domain  
**Sprint**: Sprint 1  
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

Fitur **Account & Contact Management** memungkinkan sales rep untuk melihat, mencari, dan mengakses detail informasi tentang accounts (hospitals, clinics, pharmacies) dan contacts (key persons). Fitur ini merupakan fitur fundamental untuk sales activities dan terintegrasi dengan visit reports, tasks, dan deals.

### Goals

- **Account Access**: Sales rep dapat melihat semua accounts yang ditugaskan
- **Contact Management**: Akses ke contact persons untuk setiap account
- **Search & Filter**: Cari accounts dan contacts dengan cepat
- **Offline Support**: Cache data untuk akses offline
- **Quick Actions**: Quick navigation ke related features (visit reports, tasks)

---

## Fitur Utama

### 1. Account List

**Features**:

- List accounts dengan pagination (20 items per page)
- Search accounts by name, email, atau phone
- Filter by industry, type, atau status
- Pull-to-refresh untuk sync data terbaru
- Infinite scroll untuk load more
- Offline indicator saat menggunakan cached data

**UI Components**:

- AccountCard: Display account summary
- SearchBar: Real-time search dengan debounce
- FilterChips: Industry dan type filters
- EmptyState: Saat tidak ada data
- LoadingSkeleton: Loading placeholder

### 2. Account Detail

**Information Displayed**:

- Basic info: Name, type, industry, status
- Contact info: Email, phone, website, address
- Location: Map integration (view only)
- Contacts: List of key persons
- Recent Activities: Visit reports, tasks
- Quick Actions: Create visit report, add task

### 3. Contact List

**Features**:

- List contacts per account atau all contacts
- Contact cards dengan foto, name, role
- Quick call/email actions
- Search within contacts

### 4. Contact Detail

**Information Displayed**:

- Personal info: Name, title, department
- Contact info: Email, phone, mobile
- Account affiliation
- Recent interactions

---

## Business Rules

### 1. Account Visibility Rules

**Sales Rep**:

- Hanya dapat melihat accounts yang ditugaskan ke mereka (`assigned_to`)
- Tidak dapat melihat accounts milik sales rep lain (kecuali supervisor view)

**Supervisor**:

- Dapat melihat semua accounts di team-nya
- Dapat melihat assignment status

**Admin**:

- Dapat melihat semua accounts di organisasi
- Dapat edit/delete semua accounts

### 2. Contact Association

- Setiap contact harus terkait dengan minimal satu account
- Contact dapat memiliki multiple accounts (jika di different locations)
- Primary account untuk setiap contact

### 3. Data Quality Rules

**Required Fields - Account**:

- Name (unique per organization)
- Type (hospital, clinic, pharmacy, dll.)
- Status (active, inactive)

**Required Fields - Contact**:

- Name
- Account ID

### 4. Offline Rules

- Cache account list dengan TTL 15 menit
- Cache account detail dengan TTL 30 menit
- Allow view cached data saat offline
- Show offline indicator
- Auto-refresh saat connection restored

### 5. Search Rules

**Searchable Fields - Account**:

- Name (fuzzy search)
- Email
- Phone
- Address

**Searchable Fields - Contact**:

- Name
- Email
- Phone
- Title

---

## Keputusan Teknis & Trade-offs

### Menggunakan Web API vs Mobile-Specific API

**Keputusan**: Menggunakan web API endpoints (`/api/v1/accounts`) karena mobile-specific endpoints belum tersedia.

**Alasan**:

- **Availability**: Web API sudah stabil dan tested
- **Consistency**: Data format sama dengan web app
- **Development Speed**: Tidak perlu menunggu backend team

**Trade-off**: Response format mungkin tidak optimal untuk mobile (terlalu banyak data). **Mitigasi**: Implement response filtering dan field selection.

### List View vs Map View

**Keputusan**: Prioritaskan list view dengan map sebagai secondary view.

**Alasan**:

- **Use Case**: Sales rep lebih sering mencari by name daripada by location
- **Performance**: List view lebih cepat render
- **Offline**: List view bisa di-cache, map requires tiles

**Trade-off**: Kurang visual untuk melihat spatial distribution. **Mitigasi**: Add map view sebagai optional view mode.

### Contact sebagai Separate Feature atau Sub-feature

**Keputusan**: Contact sebagai separate feature dengan navigation dari Account Detail.

**Alasan**:

- **Flexibility**: User dapat browse contacts independent dari accounts
- **Searchability**: Global contact search lebih mudah
- **Code Organization**: Clean separation of concerns

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   ├── accounts/
│   │   ├── data/
│   │   │   ├── models/
│   │   │   │   ├── account_model.dart        # Account entity
│   │   │   │   └── account_filter.dart       # Filter parameters
│   │   │   └── account_repository.dart       # API & cache
│   │   ├── application/
│   │   │   ├── account_list_provider.dart    # List state management
│   │   │   ├── account_detail_provider.dart  # Detail state management
│   │   │   └── account_search_provider.dart  # Search state
│   │   └── presentation/
│   │       ├── screens/
│   │       │   ├── account_list_screen.dart
│   │       │   ├── account_detail_screen.dart
│   │       │   └── account_search_screen.dart
│   │       └── widgets/
│   │           ├── account_card.dart
│   │           ├── account_filter_chips.dart
│   │           ├── account_info_section.dart
│   │           └── account_actions.dart
│   └── contacts/
│       ├── data/
│       │   ├── models/
│       │   │   └── contact_model.dart
│       │   └── contact_repository.dart
│       ├── application/
│       │   ├── contact_list_provider.dart
│       │   └── contact_detail_provider.dart
│       └── presentation/
│           ├── screens/
│           │   ├── contact_list_screen.dart
│           │   └── contact_detail_screen.dart
│           └── widgets/
│               ├── contact_card.dart
│               └── contact_actions.dart
```

---

## API Endpoints

### Account Endpoints

#### GET /api/v1/accounts

List accounts dengan pagination dan filter.

**Query Parameters**:

```
?page=1&limit=20&search=keyword&industry=hospital&type=customer&status=active
```

**Response**:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "name": "RS Medika Hospital",
        "type": "hospital",
        "industry": "healthcare",
        "status": "active",
        "email": "info@rsmedika.com",
        "phone": "+62123456789",
        "address": "Jl. Sudirman No. 123",
        "city": "Jakarta",
        "province": "DKI Jakarta",
        "assigned_to": "user-uuid",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 100,
      "per_page": 20
    }
  },
  "timestamp": "2025-01-15T10:30:45+07:00"
}
```

#### GET /api/v1/accounts/:id

Get account detail dengan contacts.

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "RS Medika Hospital",
    "type": "hospital",
    "industry": "healthcare",
    "status": "active",
    "email": "info@rsmedika.com",
    "phone": "+62123456789",
    "website": "www.rsmedika.com",
    "address": "Jl. Sudirman No. 123",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "notes": "Major hospital in Jakarta",
    "assigned_to": {
      "id": "user-uuid",
      "name": "John Doe",
      "email": "john@example.com"
    },
    "contacts": [
      {
        "id": "contact-uuid",
        "name": "Dr. Sarah Johnson",
        "title": "Head of Procurement",
        "department": "Procurement",
        "email": "sarah@rsmedika.com",
        "phone": "+62123456780",
        "is_primary": true
      }
    ],
    "recent_activities": [
      {
        "type": "visit_report",
        "id": "vr-uuid",
        "title": "Monthly Visit",
        "date": "2025-01-10T14:00:00Z"
      }
    ],
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Contact Endpoints

#### GET /api/v1/contacts

List contacts dengan filter.

**Query Parameters**:

```
?page=1&limit=20&search=keyword&account_id=uuid
```

**Response**:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "name": "Dr. Sarah Johnson",
        "title": "Head of Procurement",
        "department": "Procurement",
        "email": "sarah@rsmedika.com",
        "phone": "+62123456780",
        "mobile": "+6281234567890",
        "account_id": "account-uuid",
        "account_name": "RS Medika Hospital",
        "is_primary": true,
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 3,
      "total_items": 50
    }
  }
}
```

#### GET /api/v1/contacts/:id

Get contact detail.

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Dr. Sarah Johnson",
    "title": "Head of Procurement",
    "department": "Procurement",
    "email": "sarah@rsmedika.com",
    "phone": "+62123456780",
    "mobile": "+6281234567890",
    "notes": "Decision maker for procurement",
    "account_id": "account-uuid",
    "account": {
      "id": "account-uuid",
      "name": "RS Medika Hospital",
      "type": "hospital"
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

## Data Models

### Account Model

```dart
@freezed
class Account with _$Account {
  const factory Account({
    required String id,
    required String name,
    required String type,
    required String industry,
    required String status,
    String? email,
    String? phone,
    String? website,
    String? address,
    String? city,
    String? province,
    String? postalCode,
    double? latitude,
    double? longitude,
    String? notes,
    required String assignedTo,
    String? assignedToName,
    @Default([]) List<Contact> contacts,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _Account;

  factory Account.fromJson(Map<String, dynamic> json) =>
      _$AccountFromJson(json);
}

// Hive adapter untuk offline storage
@HiveType(typeId: 1)
class AccountHive extends HiveObject {
  @HiveField(0)
  late String id;

  @HiveField(1)
  late String name;

  @HiveField(2)
  late String type;

  @HiveField(3)
  late String industry;

  @HiveField(4)
  late String status;

  @HiveField(5)
  String? email;

  @HiveField(6)
  String? phone;

  @HiveField(7)
  String? address;

  @HiveField(8)
  DateTime? cachedAt;
}
```

### Contact Model

```dart
@freezed
class Contact with _$Contact {
  const factory Contact({
    required String id,
    required String name,
    String? title,
    String? department,
    String? email,
    String? phone,
    String? mobile,
    String? notes,
    required String accountId,
    String? accountName,
    @Default(false) bool isPrimary,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _Contact;

  factory Contact.fromJson(Map<String, dynamic> json) =>
      _$ContactFromJson(json);
}
```

### Account Filter

```dart
@freezed
class AccountFilter with _$AccountFilter {
  const factory AccountFilter({
    String? search,
    String? industry,
    String? type,
    String? status,
    String? city,
    @Default(1) int page,
    @Default(20) int limit,
  }) = _AccountFilter;

  Map<String, dynamic> toQueryParameters() {
    return {
      if (search != null && search!.isNotEmpty) 'search': search,
      if (industry != null) 'industry': industry,
      if (type != null) 'type': type,
      if (status != null) 'status': status,
      if (city != null) 'city': city,
      'page': page,
      'limit': limit,
    };
  }
}
```

---

## Configuration

### Repository Implementation

**File**: `features/accounts/data/account_repository.dart`

```dart
class AccountRepository {
  final ApiClient _apiClient;
  final Box<AccountHive> _localBox;
  final ConnectivityService _connectivity;

  AccountRepository(
    this._apiClient,
    this._localBox,
    this._connectivity,
  );

  Future<AccountListResult> getAccounts({
    AccountFilter? filter,
    bool forceRefresh = false,
  }) async {
    final cacheKey = 'accounts_${filter?.hashCode ?? 'default'}';

    // Load dari cache dulu
    if (!forceRefresh && _localBox.isNotEmpty) {
      final cachedAccounts = _localBox.values
          .where((a) => filter == null || _matchesFilter(a, filter))
          .map((a) => _fromHive(a))
          .toList();

      if (cachedAccounts.isNotEmpty) {
        // Return cached, refresh di background
        if (_connectivity.isOnline) {
          _refreshInBackground(filter, cacheKey);
        }

        return AccountListResult(
          accounts: cachedAccounts,
          isOffline: !_connectivity.isOnline,
          hasMore: false, // Cache tidak tahu pagination
        );
      }
    }

    // Fetch dari API
    if (_connectivity.isOnline) {
      try {
        final response = await _apiClient.get(
          '/api/v1/accounts',
          queryParameters: filter?.toQueryParameters(),
        );

        final data = response.data['data'];
        final accounts = (data['items'] as List)
            .map((json) => Account.fromJson(json))
            .toList();

        // Save ke cache
        await _saveToCache(accounts, cacheKey);

        return AccountListResult(
          accounts: accounts,
          pagination: PaginationInfo.fromJson(data['pagination']),
          isOffline: false,
          hasMore: accounts.length >= (filter?.limit ?? 20),
        );
      } catch (e) {
        // Fallback ke cache jika error
        final cachedAccounts = _localBox.values.map((a) => _fromHive(a)).toList();
        if (cachedAccounts.isNotEmpty) {
          return AccountListResult(
            accounts: cachedAccounts,
            isOffline: true,
            error: e.toString(),
          );
        }
        rethrow;
      }
    }

    throw Exception('No internet connection and no cached data');
  }

  Future<Account?> getAccountById(String id) async {
    // Try cache first
    final cached = _localBox.get(id);
    if (cached != null) {
      if (_connectivity.isOnline) {
        _refreshAccountInBackground(id);
      }
      return _fromHive(cached);
    }

    // Fetch dari API
    if (_connectivity.isOnline) {
      final response = await _apiClient.get('/api/v1/accounts/$id');
      final account = Account.fromJson(response.data['data']);
      await _localBox.put(id, _toHive(account));
      return account;
    }

    return null;
  }

  Future<void> _saveToCache(List<Account> accounts, String cacheKey) async {
    for (final account in accounts) {
      await _localBox.put(account.id, _toHive(account));
    }
  }

  Account _fromHive(AccountHive hive) {
    return Account(
      id: hive.id,
      name: hive.name,
      type: hive.type,
      industry: hive.industry,
      status: hive.status,
      email: hive.email,
      phone: hive.phone,
      address: hive.address,
    );
  }

  AccountHive _toHive(Account account) {
    return AccountHive()
      ..id = account.id
      ..name = account.name
      ..type = account.type
      ..industry = account.industry
      ..status = account.status
      ..email = account.email
      ..phone = account.phone
      ..address = account.address
      ..cachedAt = DateTime.now();
  }
}
```

---

## Usage Examples

### Account List Screen

```dart
class AccountListScreen extends ConsumerWidget {
  const AccountListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(accountListProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Accounts'),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => context.push(AppRoutes.accountSearch),
          ),
        ],
      ),
      body: Column(
        children: [
          // Filter chips
          AccountFilterChips(
            onFilterChanged: (filter) {
              ref.read(accountListProvider.notifier).setFilter(filter);
            },
          ),
          // Account list
          Expanded(
            child: state.when(
              loading: () => const AccountListSkeleton(),
              error: (error) => ErrorWidget(
                error: error,
                onRetry: () => ref.read(accountListProvider.notifier).refresh(),
              ),
              data: (result) => RefreshIndicator(
                onRefresh: () => ref.read(accountListProvider.notifier).refresh(),
                child: ListView.builder(
                  itemCount: result.accounts.length + (result.hasMore ? 1 : 0),
                  itemBuilder: (context, index) {
                    if (index == result.accounts.length) {
                      // Load more indicator
                      ref.read(accountListProvider.notifier).loadMore();
                      return const LoadingIndicator();
                    }

                    final account = result.accounts[index];
                    return AccountCard(
                      account: account,
                      onTap: () => context.push(
                        AppRoutes.accountDetailPath(account.id),
                      ),
                    );
                  },
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
```

### Account Detail Screen

```dart
class AccountDetailScreen extends ConsumerWidget {
  final String accountId;

  const AccountDetailScreen({super.key, required this.accountId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(accountDetailProvider(accountId));

    return Scaffold(
      body: state.when(
        loading: () => const AccountDetailSkeleton(),
        error: (error) => ErrorWidget(error: error),
        data: (account) => CustomScrollView(
          slivers: [
            SliverAppBar(
              expandedHeight: 200,
              flexibleSpace: FlexibleSpaceBar(
                title: Text(account.name),
                background: AccountMapPreview(
                  latitude: account.latitude,
                  longitude: account.longitude,
                ),
              ),
            ),
            SliverToBoxAdapter(
              child: Column(
                children: [
                  // Account Info
                  AccountInfoSection(account: account),

                  // Quick Actions
                  AccountActions(
                    accountId: account.id,
                    onCreateVisit: () => context.push(
                      AppRoutes.visitReportCreate,
                      extra: {'accountId': account.id},
                    ),
                    onAddTask: () => context.push(
                      AppRoutes.taskCreate,
                      extra: {'accountId': account.id},
                    ),
                  ),

                  // Contacts
                  ContactListSection(
                    contacts: account.contacts,
                    onViewAll: () => context.push(
                      AppRoutes.contacts,
                      extra: {'accountId': account.id},
                    ),
                  ),

                  // Recent Activities
                  RecentActivitiesSection(
                    activities: account.recentActivities,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

---

## Cara Test Manual

### Test Account List

1. **List Loading**:
   - Buka Accounts screen
   - Verifikasi: Loading skeleton muncul
   - Verifikasi: Accounts list muncul setelah loading

2. **Search**:
   - Tap search icon
   - Type keyword
   - Verifikasi: Results terfilter dengan benar
   - Verifikasi: Debounce working (tidak search setiap keystroke)

3. **Pull-to-Refresh**:
   - Pull down pada list
   - Verifikasi: Refresh indicator muncul
   - Verifikasi: Data terrefresh

4. **Infinite Scroll**:
   - Scroll sampai akhir list
   - Verifikasi: Load more indicator muncul
   - Verifikasi: More data loaded

5. **Offline Mode**:
   - Matikan internet
   - Buka Accounts screen
   - Verifikasi: Offline indicator muncul
   - Verifikasi: Cached data ditampilkan

### Test Account Detail

1. **Detail Loading**:
   - Tap account dari list
   - Verifikasi: Navigate ke detail screen
   - Verifikasi: All account info displayed correctly

2. **Contact List**:
   - Scroll ke Contacts section
   - Verifikasi: Contacts displayed
   - Tap contact
   - Verifikasi: Navigate ke contact detail

3. **Quick Actions**:
   - Tap "Create Visit Report"
   - Verifikasi: Navigate ke visit report form dengan account pre-selected

---

## Dependencies

### Internal

- `core/network/api_client.dart` - HTTP client
- `core/storage/hive_storage.dart` - Local cache
- `features/permissions/application/permission_provider.dart` - RBAC

### External

- `flutter_riverpod: ^2.4.0` - State management
- `freezed: ^2.4.0` - Immutable models
- `hive: ^2.2.3` - Local database
- `google_maps_flutter: ^2.5.0` - Map preview (optional)

---

## Notes & Improvements

### Known Limitations

1. **No Real-time Sync**: Data tidak auto-update saat ada changes di backend.

2. **Limited Offline**: Hanya read-only offline, tidak bisa create/edit offline.

3. **No Account Creation**: Mobile app hanya view accounts, tidak bisa create/edit.

### Future Improvements

1. **Real-time Updates**: WebSocket atau polling untuk real-time updates

2. **Full Offline**: Support create/edit accounts offline dengan sync queue

3. **Account Creation**: Add functionality untuk create new accounts dari mobile

4. **Advanced Search**: Full-text search dengan autocomplete

5. **Nearby Accounts**: Show accounts nearby menggunakan GPS

6. **Account Analytics**: Show account activity charts dan statistics

---

**Document Status**: Active  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)
