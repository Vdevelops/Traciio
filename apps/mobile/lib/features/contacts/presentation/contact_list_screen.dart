import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'dart:async';

import '../application/contact_provider.dart';
import '../application/contact_state.dart';
import '../presentation/contact_form_screen.dart';
import '../../../core/permissions/permission_provider.dart';
import '../../../core/routing/app_router.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/widgets/error_widget.dart';
import '../../../core/widgets/loading_widget.dart';
import '../../../core/widgets/skeleton_widget.dart';
import '../../../core/widgets/sync_status_indicator.dart';
import 'widgets/contact_card.dart';

class ContactListScreen extends ConsumerStatefulWidget {
  const ContactListScreen({
    super.key,
    this.accountId,
    this.hideAppBar = false,
    this.searchController,
  });

  final String? accountId;
  final bool hideAppBar;
  final TextEditingController? searchController;

  @override
  ConsumerState<ContactListScreen> createState() => _ContactListScreenState();
}

class _ContactListScreenState extends ConsumerState<ContactListScreen> {
  Timer? _debounceTimer;
  // Removed custom ScrollController to allow NestedScrollView to handle scrolling
  bool _contactWasDeleted = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;
      if (widget.accountId != null) {
        ref
            .read(contactListProvider.notifier)
            .setAccountFilter(widget.accountId);
        ref
            .read(contactListProvider.notifier)
            .loadContacts(accountId: widget.accountId);
      } else {
        ref.read(contactListProvider.notifier).setAccountFilter(null);
        ref.read(contactListProvider.notifier).loadContacts(accountId: null);
      }
    });
  }

  @override
  void dispose() {
    _debounceTimer?.cancel();
    super.dispose();
  }

  bool _onScrollNotification(ScrollNotification notification) {
    if (notification is ScrollUpdateNotification) {
      if (notification.metrics.pixels >=
          notification.metrics.maxScrollExtent * 0.8) {
        ref.read(contactListProvider.notifier).loadMore();
      }
    }
    return false;
  }

  void _onSearchChanged(String query) {
    // Search is handled by parent (AccountsScreen)
    // This method is kept for compatibility but not used when searchController is provided
  }

  Future<void> _onRefresh() async {
    await ref.read(contactListProvider.notifier).refresh();
    if (!mounted) return;
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(contactListProvider);
    final theme = Theme.of(context);

    final body = widget.searchController != null
        ? NotificationListener<ScrollNotification>(
            onNotification: _onScrollNotification,
            child: RefreshIndicator(
              onRefresh: _onRefresh,
              child: _buildContent(context, state, theme),
            ),
          )
        : Column(
            children: [
              // Search Bar (only if not provided by parent)
              Container(
                margin: const EdgeInsets.fromLTRB(20, 8, 20, 16),
                decoration: BoxDecoration(
                  color: theme.colorScheme.surfaceContainerHighest.withValues(
                    alpha: 0.3,
                  ),
                  borderRadius: BorderRadius.circular(30),
                ),
                child: TextField(
                  onChanged: _onSearchChanged,
                  decoration: InputDecoration(
                    hintText: 'Search contacts...',
                    hintStyle: TextStyle(color: theme.hintColor),
                    prefixIcon: Icon(Icons.search, color: theme.hintColor),
                    border: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 20,
                      vertical: 14,
                    ),
                  ),
                ),
              ),
              // Content
              Expanded(
                child: NotificationListener<ScrollNotification>(
                  onNotification: _onScrollNotification,
                  child: RefreshIndicator(
                    onRefresh: _onRefresh,
                    child: _buildContent(context, state, theme),
                  ),
                ),
              ),
            ],
          );

    if (widget.hideAppBar) {
      return body;
    }

    // Wrap with PopScope to handle back navigation when contact is deleted
    Widget scaffold = Scaffold(
      appBar: AppBar(
        title: Text(widget.accountId != null ? 'Contacts' : 'All Contacts'),
        elevation: 0,
        actions: const [
          // Sync Status Indicator
          Padding(
            padding: EdgeInsets.only(right: 16),
            child: Center(
              child: SyncStatusIndicator(featureKey: 'contacts', iconSize: 20),
            ),
          ),
        ],
      ),
      body: body,
      floatingActionButton: widget.accountId != null
          ? (ref.watch(canCreateProvider('accounts'))
                ? FloatingActionButton(
                    onPressed: () async {
                      final result = await Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => ContactFormScreen(
                            defaultAccountId: widget.accountId,
                          ),
                        ),
                      );
                      if (result != null && mounted) {
                        await ref.read(contactListProvider.notifier).refresh();
                      }
                    },
                    child: const Icon(Icons.add),
                  )
                : null)
          : null,
    );

    // If opened from account detail and contact was deleted, wrap with PopScope
    if (widget.accountId != null) {
      return PopScope(
        canPop: true,
        onPopInvokedWithResult: (didPop, result) {
          // If contact was deleted, return true to trigger refresh in parent
          if (didPop && _contactWasDeleted) {
            // Use a post-frame callback to ensure navigation is complete
            WidgetsBinding.instance.addPostFrameCallback((_) {
              if (mounted) {
                Navigator.of(context).pop(true);
              }
            });
          }
        },
        child: scaffold,
      );
    }

    return scaffold;
  }

  Widget _buildContent(
    BuildContext context,
    ContactListState state,
    ThemeData theme,
  ) {
    if (state.isLoading && state.contacts.isEmpty) {
      return const LoadingWidget();
    }

    if (state.errorMessage != null && state.contacts.isEmpty) {
      return ErrorStateWidget(
        message: state.errorMessage!,
        onRetry: () {
          ref.read(contactListProvider.notifier).refresh();
        },
      );
    }

    if (state.contacts.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.person_outline,
              size: 64,
              color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
            ),
            const SizedBox(height: 16),
            Text(
              AppLocalizations.of(context)!.noContactsFound,
              style: theme.textTheme.titleMedium?.copyWith(
                color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
              ),
            ),
          ],
        ),
      );
    }

    // Show skeleton screens if loading first page
    if (state.isLoading && state.contacts.isEmpty) {
      return ListView.builder(
        itemCount: 5, // Show 5 skeleton items
        itemBuilder: (context, index) {
          return const SkeletonListItem(height: 80);
        },
      );
    }

    return ListView.builder(
      // controller removed for NestedScroll support
      padding: const EdgeInsets.fromLTRB(20, 16, 20, 100),
      itemCount: state.contacts.length + (state.isLoadingMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == state.contacts.length) {
          return const Padding(
            padding: EdgeInsets.all(16),
            child: LoadingWidget(size: 24),
          );
        }

        final contact = state.contacts[index];
        return ContactCard(
          contact: contact,
          onTap: () async {
            final result = await Navigator.pushNamed(
              context,
              '${AppRoutes.contacts}/${contact.id}',
              arguments: contact.id,
            );
            // Refresh list if contact was deleted or updated
            if (result == true && mounted) {
              _contactWasDeleted = true;
              await ref.read(contactListProvider.notifier).refresh();
            }
          },
        );
      },
    );
  }
}
