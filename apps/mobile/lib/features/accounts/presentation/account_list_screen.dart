import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../application/account_provider.dart';
import '../application/account_state.dart';
import '../../../core/routing/app_router.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/widgets/error_widget.dart';
import '../../../core/widgets/loading_widget.dart';
import '../../../core/widgets/skeleton_widget.dart';
import '../../../core/widgets/sync_status_indicator.dart';
import 'widgets/account_card.dart';

class AccountListScreen extends ConsumerStatefulWidget {
  const AccountListScreen({super.key, this.hideAppBar = false});

  final bool hideAppBar;

  @override
  ConsumerState<AccountListScreen> createState() => _AccountListScreenState();
}

class _AccountListScreenState extends ConsumerState<AccountListScreen> {
  @override
  void initState() {
    super.initState();
    // Initial load
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) {
        ref.read(accountListProvider.notifier).loadAccounts();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(accountListProvider);
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;

    final hasActiveSearch = state.searchQuery.isNotEmpty;

    Widget bodyContent = NotificationListener<ScrollNotification>(
      onNotification: (ScrollNotification notification) {
        if (notification is ScrollUpdateNotification &&
            notification.metrics.pixels >=
                notification.metrics.maxScrollExtent * 0.8) {
          ref.read(accountListProvider.notifier).loadMore();
        }
        return false;
      },
      child: RefreshIndicator(
        onRefresh: () => ref.read(accountListProvider.notifier).refresh(),
        child: _buildContent(context, ref, state, theme),
      ),
    );

    final body = Column(
      children: [
        if (hasActiveSearch)
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 8, 16, 0),
            child: Row(
              children: [
                Expanded(
                  child: Text(
                    '${l10n.searchResultsFor} "${state.searchQuery}"',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onSurface.withValues(alpha: 0.7),
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                const SizedBox(width: 8),
                TextButton.icon(
                  onPressed: () {
                    ref
                        .read(accountListProvider.notifier)
                        .updateSearchQuery('');
                    ref
                        .read(accountListProvider.notifier)
                        .loadAccounts(page: 1, refresh: true, search: '');
                  },
                  icon: Icon(
                    Icons.close,
                    size: 16,
                    color: theme.colorScheme.error,
                  ),
                  label: Text(
                    l10n.clearSearch,
                    style: theme.textTheme.labelLarge?.copyWith(
                      color: theme.colorScheme.error,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  style: TextButton.styleFrom(
                    visualDensity: VisualDensity.compact,
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 4,
                    ),
                    backgroundColor: theme.colorScheme.error.withValues(
                      alpha: 0.1,
                    ),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(20),
                    ),
                  ),
                ),
              ],
            ),
          ),
        Expanded(child: bodyContent),
      ],
    );

    if (widget.hideAppBar) {
      return body;
    }

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.accounts),
        elevation: 0,
        actions: const [
          // Sync Status Indicator
          Padding(
            padding: EdgeInsets.only(right: 16),
            child: Center(
              child: SyncStatusIndicator(featureKey: 'accounts', iconSize: 20),
            ),
          ),
        ],
      ),
      body: body,
    );
  }

  Widget _buildContent(
    BuildContext context,
    WidgetRef ref,
    AccountListState state,
    ThemeData theme,
  ) {
    final l10n = AppLocalizations.of(context)!;

    if (state.isLoading && state.accounts.isEmpty) {
      return ListView.builder(
        itemCount: 5,
        padding: const EdgeInsets.all(16),
        itemBuilder: (context, index) {
          return const SkeletonListItem(height: 100);
        },
      );
    }

    if (state.errorMessage != null && state.accounts.isEmpty) {
      return ErrorStateWidget(
        message: state.errorMessage!,
        onRetry: () => ref.read(accountListProvider.notifier).refresh(),
      );
    }

    if (state.accounts.isEmpty) {
      if (state.searchQuery.isNotEmpty) {
        return EmptyStateWidget(
          message: '${l10n.noSearchResults} "${state.searchQuery}"',
          subtitle: l10n.noAccountsFound,
          icon: Icons.search_off,
        );
      }
      return EmptyStateWidget(
        message: l10n.noAccountsFound,
        icon: Icons.business_outlined,
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.fromLTRB(20, 16, 20, 100),
      itemCount: state.accounts.length + (state.isLoadingMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == state.accounts.length) {
          return const Padding(
            padding: EdgeInsets.all(16),
            child: LoadingWidget(size: 24),
          );
        }
        final account = state.accounts[index];
        return AccountCard(
          account: account,
          onTap: () {
            Navigator.pushNamed(
              context,
              '${AppRoutes.accounts}/${account.id}',
              arguments: account.id,
            );
          },
        );
      },
    );
  }
}
