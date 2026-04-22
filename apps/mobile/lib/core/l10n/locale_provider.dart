import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../storage/local_storage.dart';

final localStorageProvider = FutureProvider<LocalStorage>((ref) async {
  return LocalStorage.create();
});

final localeProvider =
    NotifierProvider<LocaleNotifier, Locale>(LocaleNotifier.new);

class LocaleNotifier extends Notifier<Locale> {
  @override
  Locale build() {
    ref.listen(localStorageProvider, (prev, next) {
      next.whenData((localStorage) {
        state = _parseLocale(localStorage.getLocale());
      });
    });
    return const Locale('en', '');
  }

  Locale _parseLocale(String localeString) {
    switch (localeString) {
      case 'id':
        return const Locale('id', '');
      case 'en':
      default:
        return const Locale('en', '');
    }
  }

  String _localeToString(Locale locale) {
    return locale.languageCode;
  }

  Future<void> setLocale(Locale locale) async {
    await ref.read(localStorageProvider).when(
          data: (localStorage) async {
            await localStorage.setLocale(_localeToString(locale));
            state = locale;
          },
          loading: () async {},
          error: (_, _) async {},
        );
  }
}
