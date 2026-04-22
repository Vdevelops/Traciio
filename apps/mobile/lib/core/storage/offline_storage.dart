import 'dart:convert';

import 'hive_storage.dart';

/// Offline storage helper using Hive with JSON serialization
///
/// This provides a simple way to cache data for offline access
/// without requiring Hive adapters for each model
class OfflineStorage {
  static const String _accountsBox = 'offline_accounts';
  static const String _contactsBox = 'offline_contacts';
  static const String _tasksBox = 'offline_tasks';
  static const String _visitReportsBox = 'offline_visit_reports';
  static const String _dashboardBox = 'offline_dashboard';
  static const String _dealsBox = 'offline_deals';
  static const String _stagesBox = 'offline_stages';
  static const String _productsBox = 'offline_products';
  static const String _schedulesBox = 'offline_schedules';
  static const String _leadsBox = 'offline_leads';
  static const String _userBox = 'offline_user';

  /// Initialize offline storage boxes
  static Future<void> init() async {
    await HiveStorage.openBox(_accountsBox);
    await HiveStorage.openBox(_contactsBox);
    await HiveStorage.openBox(_tasksBox);
    await HiveStorage.openBox(_visitReportsBox);
    await HiveStorage.openBox(_dashboardBox);
    await HiveStorage.openBox(_dealsBox);
    await HiveStorage.openBox(_stagesBox);
    await HiveStorage.openBox(_productsBox);
    await HiveStorage.openBox(_schedulesBox);
    await HiveStorage.openBox(_leadsBox);
    await HiveStorage.openBox(_userBox);
    await HiveStorage.openBox('offline_generic');
  }

  /// Save accounts to offline storage
  static Future<void> saveAccounts(List<Map<String, dynamic>> accounts) async {
    final box = await HiveStorage.openBox(_accountsBox);
    await box.put('list', jsonEncode(accounts));
    await box.put('cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached accounts from offline storage
  static Future<List<Map<String, dynamic>>?> getAccounts() async {
    final box = await HiveStorage.openBox(_accountsBox);
    final cachedData = box.get('list');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save account detail to offline storage
  static Future<void> saveAccountDetail(
    String id,
    Map<String, dynamic> account,
  ) async {
    final box = await HiveStorage.openBox(_accountsBox);
    await box.put('detail_$id', jsonEncode(account));
  }

  /// Get cached account detail from offline storage
  static Future<Map<String, dynamic>?> getAccountDetail(String id) async {
    final box = await HiveStorage.openBox(_accountsBox);
    final cachedData = box.get('detail_$id');
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save deals to offline storage
  static Future<void> saveDeals(List<Map<String, dynamic>> deals) async {
    final box = await HiveStorage.openBox(_dealsBox);
    await box.put('list', jsonEncode(deals));
    await box.put('cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached deals from offline storage
  static Future<List<Map<String, dynamic>>?> getDeals() async {
    final box = await HiveStorage.openBox(_dealsBox);
    final cachedData = box.get('list');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save deal detail to offline storage
  static Future<void> saveDealDetail(
    String id,
    Map<String, dynamic> deal,
  ) async {
    final box = await HiveStorage.openBox(_dealsBox);
    await box.put('detail_$id', jsonEncode(deal));
  }

  /// Get cached deal detail from offline storage
  static Future<Map<String, dynamic>?> getDealDetail(String id) async {
    final box = await HiveStorage.openBox(_dealsBox);
    final cachedData = box.get('detail_$id');
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save pipeline stages to offline storage
  static Future<void> saveStages(List<Map<String, dynamic>> stages) async {
    final box = await HiveStorage.openBox(_stagesBox);
    await box.put('list', jsonEncode(stages));
    await box.put('cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached pipeline stages from offline storage
  static Future<List<Map<String, dynamic>>?> getStages() async {
    final box = await HiveStorage.openBox(_stagesBox);
    final cachedData = box.get('list');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save products to offline storage
  static Future<void> saveProducts(List<Map<String, dynamic>> products) async {
    final box = await HiveStorage.openBox(_productsBox);
    await box.put('list', jsonEncode(products));
    await box.put('cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached products from offline storage
  static Future<List<Map<String, dynamic>>?> getProducts() async {
    final box = await HiveStorage.openBox(_productsBox);
    final cachedData = box.get('list');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save contacts to offline storage
  static Future<void> saveContacts(List<Map<String, dynamic>> contacts) async {
    final box = await HiveStorage.openBox(_contactsBox);
    await box.put('list', jsonEncode(contacts));
    await box.put('cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached contacts from offline storage
  static Future<List<Map<String, dynamic>>?> getContacts() async {
    final box = await HiveStorage.openBox(_contactsBox);
    final cachedData = box.get('list');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save contact detail to offline storage
  static Future<void> saveContactDetail(
    String id,
    Map<String, dynamic> contact,
  ) async {
    final box = await HiveStorage.openBox(_contactsBox);
    await box.put('detail_$id', jsonEncode(contact));
  }

  /// Get cached contact detail from offline storage
  static Future<Map<String, dynamic>?> getContactDetail(String id) async {
    final box = await HiveStorage.openBox(_contactsBox);
    final cachedData = box.get('detail_$id');
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save tasks to offline storage
  static Future<void> saveTasks(List<Map<String, dynamic>> tasks) async {
    final box = await HiveStorage.openBox(_tasksBox);
    await box.put('list', jsonEncode(tasks));
    await box.put('cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached tasks from offline storage
  static Future<List<Map<String, dynamic>>?> getTasks() async {
    final box = await HiveStorage.openBox(_tasksBox);
    final cachedData = box.get('list');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save task detail to offline storage
  static Future<void> saveTaskDetail(
    String id,
    Map<String, dynamic> task,
  ) async {
    final box = await HiveStorage.openBox(_tasksBox);
    await box.put('detail_$id', jsonEncode(task));
  }

  /// Get cached task detail from offline storage
  static Future<Map<String, dynamic>?> getTaskDetail(String id) async {
    final box = await HiveStorage.openBox(_tasksBox);
    final cachedData = box.get('detail_$id');
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save schedules to offline storage
  static Future<void> saveSchedules(
    List<Map<String, dynamic>> schedules,
  ) async {
    final box = await HiveStorage.openBox(_schedulesBox);
    await box.put('list', jsonEncode(schedules));
    await box.put('cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached schedules from offline storage
  static Future<List<Map<String, dynamic>>?> getSchedules() async {
    final box = await HiveStorage.openBox(_schedulesBox);
    final cachedData = box.get('list');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save schedule detail to offline storage
  static Future<void> saveScheduleDetail(
    String id,
    Map<String, dynamic> schedule,
  ) async {
    final box = await HiveStorage.openBox(_schedulesBox);
    await box.put('detail_$id', jsonEncode(schedule));
  }

  /// Get cached schedule detail from offline storage
  static Future<Map<String, dynamic>?> getScheduleDetail(String id) async {
    final box = await HiveStorage.openBox(_schedulesBox);
    final cachedData = box.get('detail_$id');
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  // Leads Storage Methods

  /// Save leads to offline storage
  static Future<void> saveLeads(List<Map<String, dynamic>> leads) async {
    final box = await HiveStorage.openBox(_leadsBox);
    await box.put('list', jsonEncode(leads));
    await box.put('cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached leads from offline storage
  static Future<List<Map<String, dynamic>>?> getLeads() async {
    final box = await HiveStorage.openBox(_leadsBox);
    final cachedData = box.get('list');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save lead detail to offline storage
  static Future<void> saveLeadDetail(
    String id,
    Map<String, dynamic> lead,
  ) async {
    final box = await HiveStorage.openBox(_leadsBox);
    await box.put('detail_$id', jsonEncode(lead));
  }

  /// Get cached lead detail from offline storage
  static Future<Map<String, dynamic>?> getLeadDetail(String id) async {
    final box = await HiveStorage.openBox(_leadsBox);
    final cachedData = box.get('detail_$id');
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save visit reports to offline storage
  static Future<void> saveVisitReports(
    List<Map<String, dynamic>> visitReports,
  ) async {
    final box = await HiveStorage.openBox(_visitReportsBox);
    await box.put('list', jsonEncode(visitReports));
    await box.put('cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached visit reports from offline storage
  static Future<List<Map<String, dynamic>>?> getVisitReports() async {
    final box = await HiveStorage.openBox(_visitReportsBox);
    final cachedData = box.get('list');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save visit report detail to offline storage
  static Future<void> saveVisitReportDetail(
    String id,
    Map<String, dynamic> visitReport,
  ) async {
    final box = await HiveStorage.openBox(_visitReportsBox);
    await box.put('detail_$id', jsonEncode(visitReport));
  }

  /// Get cached visit report detail from offline storage
  static Future<Map<String, dynamic>?> getVisitReportDetail(String id) async {
    final box = await HiveStorage.openBox(_visitReportsBox);
    final cachedData = box.get('detail_$id');
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Clear visit report detail cache
  static Future<void> clearVisitReportDetail(String id) async {
    final box = await HiveStorage.openBox(_visitReportsBox);
    await box.delete('detail_$id');
  }

  /// Save dashboard overview to offline storage
  static Future<void> saveDashboardOverview(
    Map<String, dynamic> overview, {
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'overview_$period' : 'overview';
    await box.put(key, jsonEncode(overview));
    await box.put('${key}_cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached dashboard overview from offline storage
  static Future<Map<String, dynamic>?> getDashboardOverview({
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'overview_$period' : 'overview';
    final cachedData = box.get(key);
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save dashboard visit statistics to offline storage
  static Future<void> saveDashboardVisitStatistics(
    Map<String, dynamic> visitStatistics, {
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null
        ? 'visit_statistics_$period'
        : 'visit_statistics';
    await box.put(key, jsonEncode(visitStatistics));
  }

  /// Get cached dashboard visit statistics from offline storage
  static Future<Map<String, dynamic>?> getDashboardVisitStatistics({
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null
        ? 'visit_statistics_$period'
        : 'visit_statistics';
    final cachedData = box.get(key);
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save dashboard recent activities to offline storage
  static Future<void> saveDashboardRecentActivities(
    List<Map<String, dynamic>> activities,
  ) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    await box.put('recent_activities', jsonEncode(activities));
  }

  /// Get cached dashboard recent activities from offline storage
  static Future<List<Map<String, dynamic>>?>
  getDashboardRecentActivities() async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final cachedData = box.get('recent_activities');
    if (cachedData != null) {
      return List<Map<String, dynamic>>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save dashboard activity trends to offline storage
  static Future<void> saveDashboardActivityTrends(
    Map<String, dynamic> activityTrends, {
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'activity_trends_$period' : 'activity_trends';
    await box.put(key, jsonEncode(activityTrends));
  }

  /// Get cached dashboard activity trends from offline storage
  static Future<Map<String, dynamic>?> getDashboardActivityTrends({
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'activity_trends_$period' : 'activity_trends';
    final cachedData = box.get(key);
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  // Mobile Dashboard Redesign Cache Methods

  /// Save mobile dashboard overview to offline storage
  static Future<void> saveMobileDashboardOverview(
    Map<String, dynamic> overview, {
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'mobile_overview_$period' : 'mobile_overview';
    await box.put(key, jsonEncode(overview));
    await box.put('${key}_cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached mobile dashboard overview from offline storage
  static Future<Map<String, dynamic>?> getMobileDashboardOverview({
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'mobile_overview_$period' : 'mobile_overview';
    final cachedData = box.get(key);
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save mobile dashboard visits to offline storage
  static Future<void> saveMobileVisits(
    List<dynamic> visits, {
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'mobile_visits_$period' : 'mobile_visits';
    await box.put(key, jsonEncode(visits));
  }

  /// Get cached mobile dashboard visits from offline storage
  static Future<List<dynamic>?> getMobileVisits({String? period}) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'mobile_visits_$period' : 'mobile_visits';
    final cachedData = box.get(key);
    if (cachedData != null) {
      return List<dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Save mobile dashboard tasks to offline storage
  static Future<void> saveMobileTasks(
    List<dynamic> tasks, {
    String? period,
  }) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'mobile_tasks_$period' : 'mobile_tasks';
    await box.put(key, jsonEncode(tasks));
  }

  /// Get cached mobile dashboard tasks from offline storage
  static Future<List<dynamic>?> getMobileTasks({String? period}) async {
    final box = await HiveStorage.openBox(_dashboardBox);
    final key = period != null ? 'mobile_tasks_$period' : 'mobile_tasks';
    final cachedData = box.get(key);
    if (cachedData != null) {
      return List<dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  // Profile Cache Methods

  /// Save user profile to offline storage
  static Future<void> saveUserProfile(Map<String, dynamic> profile) async {
    final box = await HiveStorage.openBox(_userBox);
    await box.put('profile', jsonEncode(profile));
    await box.put('profile_cached_at', DateTime.now().toIso8601String());
  }

  /// Get cached user profile from offline storage
  static Future<Map<String, dynamic>?> getUserProfile() async {
    final box = await HiveStorage.openBox(_userBox);
    final cachedData = box.get('profile');
    if (cachedData != null) {
      return Map<String, dynamic>.from(jsonDecode(cachedData as String));
    }
    return null;
  }

  /// Clear user profile cache
  static Future<void> clearUserProfile() async {
    final box = await HiveStorage.openBox(_userBox);
    await box.delete('profile');
    await box.delete('profile_cached_at');
  }

  /// Clear all offline cache
  static Future<void> clearAll() async {
    await HiveStorage.clearBox(_accountsBox);
    await HiveStorage.clearBox(_contactsBox);
    await HiveStorage.clearBox(_tasksBox);
    await HiveStorage.clearBox(_visitReportsBox);
    await HiveStorage.clearBox(_dashboardBox);
  }

  /// Clear specific entity cache
  static Future<void> clearEntity(String entity) async {
    switch (entity) {
      case 'accounts':
        await HiveStorage.clearBox(_accountsBox);
        break;
      case 'contacts':
        await HiveStorage.clearBox(_contactsBox);
        break;
      case 'tasks':
        await HiveStorage.clearBox(_tasksBox);
        break;
      case 'visit_reports':
        await HiveStorage.clearBox(_visitReportsBox);
        break;
      case 'dashboard':
        await HiveStorage.clearBox(_dashboardBox);
        break;
      case 'schedules':
        await HiveStorage.clearBox(_schedulesBox);
        break;
    }
  }

  /// Generic get method for any data type
  static Future<T?> get<T>(
    String key,
    T Function(Map<String, dynamic>) fromJson,
  ) async {
    // Use a generic box for permissions and other data
    const boxName = 'offline_generic';
    await HiveStorage.openBox(boxName);
    final box = await HiveStorage.openBox(boxName);
    final cachedData = box.get(key);
    if (cachedData != null) {
      try {
        final json = jsonDecode(cachedData as String);
        return fromJson(json as Map<String, dynamic>);
      } catch (e) {
        return null;
      }
    }
    return null;
  }

  /// Generic set method for any data type
  static Future<void> set<T>(String key, Map<String, dynamic> data) async {
    const boxName = 'offline_generic';
    await HiveStorage.openBox(boxName);
    final box = await HiveStorage.openBox(boxName);
    await box.put(key, jsonEncode(data));
  }
}
