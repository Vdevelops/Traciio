import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'dart:async';

import '../../../core/l10n/app_localizations.dart';
// Visits hanya diakses via quick nav - tidak pakai bottom navbar
import '../../../core/permissions/permission_provider.dart';
import '../application/visit_report_provider.dart';
import 'simplified_visit_report_form_screen.dart';
import 'visit_report_list_screen.dart';

class ReportsScreen extends ConsumerStatefulWidget {
  const ReportsScreen({super.key});

  @override
  ConsumerState<ReportsScreen> createState() => _ReportsScreenState();
}

class _ReportsScreenState extends ConsumerState<ReportsScreen> {
  final TextEditingController _searchController = TextEditingController();
  Timer? _debounceTimer;

  @override
  void dispose() {
    _searchController.dispose();
    _debounceTimer?.cancel();
    super.dispose();
  }

  void _onSearchChanged(String query) {
    _debounceTimer?.cancel();
    _debounceTimer = Timer(const Duration(milliseconds: 500), () {
      ref.read(visitReportListProvider.notifier).updateSearchQuery(query);
      ref
          .read(visitReportListProvider.notifier)
          .loadVisitReports(page: 1, refresh: true, search: query);
    });
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final theme = Theme.of(context);

    // Check CREATE permission
    final hasCreatePermission = ref.watch(canCreateProvider('visit-reports'));

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: null,
        automaticallyImplyLeading: false,
      ),
      floatingActionButton: hasCreatePermission
          ? FloatingActionButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => const SimplifiedVisitReportFormScreen(),
                  ),
                ).then((result) {
                  if (context.mounted && result != null) {
                    ref.read(visitReportListProvider.notifier).refresh();
                  }
                });
              },
              child: const Icon(Icons.add),
            )
          : null,
      body: NestedScrollView(
        headerSliverBuilder: (context, innerBoxIsScrolled) {
          return [
            // COLLAPSIBLE HEADER (Search only)
            SliverToBoxAdapter(
              child: SafeArea(
                bottom: false,
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(20, 16, 20, 8),
                  child: Container(
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surfaceContainerHighest
                          .withValues(alpha: 0.3),
                      borderRadius: BorderRadius.circular(30),
                    ),
                    child: TextField(
                      controller: _searchController,
                      onChanged: _onSearchChanged,
                      style: theme.textTheme.titleMedium,
                      decoration: InputDecoration(
                        hintText: l10n.searchVisitReports,
                        hintStyle: TextStyle(
                          color: theme.hintColor,
                          fontSize: 16,
                        ),
                        prefixIcon: Icon(
                          Icons.search,
                          color: theme.hintColor,
                          size: 24,
                        ),
                        suffixIcon: _searchController.text.isNotEmpty
                            ? IconButton(
                                icon: const Icon(Icons.clear),
                                onPressed: () {
                                  _searchController.clear();
                                  _onSearchChanged('');
                                },
                              )
                            : null,
                        border: InputBorder.none,
                        focusedBorder: InputBorder.none,
                        enabledBorder: InputBorder.none,
                        contentPadding: const EdgeInsets.symmetric(
                          horizontal: 20,
                          vertical: 16,
                        ),
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ];
        },
        body: VisitReportListScreen(
          hideAppBar: true,
          searchController: _searchController,
        ),
      ),
    );
  }
}
