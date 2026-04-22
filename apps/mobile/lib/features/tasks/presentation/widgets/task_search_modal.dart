import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/l10n/app_localizations.dart';
import '../../../../core/widgets/search_modal.dart';
import '../../application/task_provider.dart';

/// Show search overlay untuk tasks (muncul di atas seperti searchbar di reports)
void showTaskSearchModal(BuildContext context) {
  final l10n = AppLocalizations.of(context)!;
  // The previous `ref` usage was incorrect for a function outside a widget build method.
  // `ProviderScope.containerOf(context).read` is the correct way to access providers
  // from a BuildContext in such scenarios.
  // The comments about `showGeneralDialog` and `TaskSearchModal` refer to a previous
  // implementation that has since been refactored to use the generic `showSearchModal`.

  showSearchModal(
    context,
    hintText: l10n.searchTasks,
    initialQuery: ProviderScope.containerOf(
      context,
    ).read(taskListProvider).searchQuery,
    onSearch: (query) {
      final container = ProviderScope.containerOf(context);
      container.read(taskListProvider.notifier).updateSearchQuery(query);
      container
          .read(taskListProvider.notifier)
          .loadTasks(page: 1, refresh: true, search: query, forceRefresh: true);
    },
    onClear: () {
      final container = ProviderScope.containerOf(context);
      container.read(taskListProvider.notifier).updateSearchQuery('');
      container
          .read(taskListProvider.notifier)
          .loadTasks(page: 1, refresh: true, search: '', forceRefresh: true);
    },
  );
}

// TaskSearchModal class removed as it is now handled by Generic SearchModal
