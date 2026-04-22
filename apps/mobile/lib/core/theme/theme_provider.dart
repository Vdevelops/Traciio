import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../storage/local_storage.dart';

final localStorageProvider = FutureProvider<LocalStorage>((ref) async {
  return LocalStorage.create();
});

final themeModeProvider =
    NotifierProvider<ThemeModeNotifier, ThemeMode>(ThemeModeNotifier.new);

class ThemeModeNotifier extends Notifier<ThemeMode> {
  @override
  ThemeMode build() {
    ref.listen(localStorageProvider, (prev, next) {
      next.whenData((localStorage) {
        state = _parseThemeMode(localStorage.getThemeMode());
      });
    });
    return ThemeMode.system;
  }

  ThemeMode _parseThemeMode(String themeModeString) {
    switch (themeModeString) {
      case 'light':
        return ThemeMode.light;
      case 'dark':
        return ThemeMode.dark;
      case 'system':
        return ThemeMode.system;
      default:
        return ThemeMode.system;
    }
  }

  String _themeModeToString(ThemeMode mode) {
    switch (mode) {
      case ThemeMode.light:
        return 'light';
      case ThemeMode.dark:
        return 'dark';
      case ThemeMode.system:
        return 'system';
    }
  }

  Future<void> setThemeMode(ThemeMode mode) async {
    await ref.read(localStorageProvider).when(
          data: (localStorage) async {
            await localStorage.setThemeMode(_themeModeToString(mode));
            state = mode;
          },
          loading: () async {},
          error: (_, _) async {},
        );
  }

  Future<void> toggleTheme() async {
    ThemeMode newMode;
    if (state == ThemeMode.light) {
      newMode = ThemeMode.dark;
    } else if (state == ThemeMode.dark) {
      newMode = ThemeMode.light;
    } else {
      newMode = ThemeMode.dark;
    }
    await setThemeMode(newMode);
  }
}
