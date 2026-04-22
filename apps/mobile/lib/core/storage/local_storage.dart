import 'dart:convert';

import 'package:shared_preferences/shared_preferences.dart';

import '../../features/auth/data/models/login_response.dart';

class LocalStorage {
  LocalStorage(this._prefs);

  final SharedPreferences _prefs;

  static Future<LocalStorage> create() async {
    final prefs = await SharedPreferences.getInstance();
    return LocalStorage(prefs);
  }

  Future<void> saveAuthToken(String token) async {
    await _prefs.setString('auth_token', token);
  }

  String? getAuthToken() => _prefs.getString('auth_token');

  Future<void> saveRefreshToken(String refreshToken) async {
    await _prefs.setString('refresh_token', refreshToken);
  }

  String? getRefreshToken() => _prefs.getString('refresh_token');

  Future<void> setRememberMe(bool remember) async {
    await _prefs.setBool('remember_me', remember);
  }

  bool getRememberMe() => _prefs.getBool('remember_me') ?? false;

  Future<void> clearAuthToken() async {
    await _prefs.remove('auth_token');
    await _prefs.remove('refresh_token');
  }

  Future<void> saveUser(UserResponse user) async {
    final userJson = jsonEncode({
      'id': user.id,
      'email': user.email,
      'name': user.name,
      'avatar_url': user.avatarUrl,
      'role': user.role,
      'status': user.status,
      'created_at': user.createdAt.toIso8601String(),
      'updated_at': user.updatedAt.toIso8601String(),
    });
    await _prefs.setString('user_data', userJson);
  }

  UserResponse? getUser() {
    final userJson = _prefs.getString('user_data');
    if (userJson == null) {
      return null;
    }
    try {
      final userMap = jsonDecode(userJson) as Map<String, dynamic>;
      return UserResponse.fromJson(userMap);
    } catch (e) {
      return null;
    }
  }

  Future<void> clearUser() async {
    await _prefs.remove('user_data');
  }

  Future<void> setThemeMode(String themeMode) async {
    await _prefs.setString('theme_mode', themeMode);
  }

  String getThemeMode() => _prefs.getString('theme_mode') ?? 'light';

  Future<void> setLocale(String locale) async {
    await _prefs.setString('locale', locale);
  }

  String getLocale() => _prefs.getString('locale') ?? 'en';
}


