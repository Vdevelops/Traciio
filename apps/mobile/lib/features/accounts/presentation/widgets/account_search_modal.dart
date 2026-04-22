import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/l10n/app_localizations.dart';
import '../../../../core/widgets/search_modal.dart';
import '../../application/account_provider.dart';

/// Shows the same search overlay as tasks: top slide-down modal with hint,
/// debounced search, clear button, and close on submit. Uses shared [showSearchModal].
void showAccountSearchModal(BuildContext context) {
  final l10n = AppLocalizations.of(context)!;
  showSearchModal(
    context,
    hintText: l10n.searchAccounts,
    initialQuery:
        ProviderScope.containerOf(context).read(accountListProvider).searchQuery,
    onSearch: (query) {
      final container = ProviderScope.containerOf(context);
      container.read(accountListProvider.notifier).updateSearchQuery(query);
      container.read(accountListProvider.notifier).loadAccounts(
            page: 1,
            refresh: true,
            search: query,
            forceRefresh: true,
          );
    },
    onClear: () {
      final container = ProviderScope.containerOf(context);
      container.read(accountListProvider.notifier).updateSearchQuery('');
      container.read(accountListProvider.notifier).loadAccounts(
            page: 1,
            refresh: true,
            search: '',
            forceRefresh: true,
          );
    },
  );
}
