import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/routing/app_router.dart';
import '../../../core/widgets/main_scaffold.dart';
import '../../../core/widgets/sync_status_indicator.dart';
import '../application/lead_provider.dart';
import '../application/lead_state.dart';
import 'widgets/lead_card.dart';
import 'widgets/lead_filters.dart';
import 'widgets/lead_search_modal.dart';

class LeadListScreen extends ConsumerStatefulWidget {
  const LeadListScreen({super.key});

  @override
  ConsumerState<LeadListScreen> createState() => _LeadListScreenState();
}

class _LeadListScreenState extends ConsumerState<LeadListScreen> {
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) {
        ref.read(leadListProvider.notifier).loadLeads();
      }
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!mounted) return;
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent - 200) {
      ref.read(leadListProvider.notifier).loadMore();
    }
  }

  void _showStatusFilter() {
    if (!mounted) return;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => StatusFilterSheet(
        selectedStatus: ref.watch(leadListProvider).selectedStatus,
        onSelect: (status) {
          if (!mounted) return;
          ref.read(leadListProvider.notifier).updateStatusFilter(status);
          Navigator.pop(context);
        },
      ),
    );
  }

  void _showSourceFilter() {
    if (!mounted) return;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => SourceFilterSheet(
        selectedSource: ref.watch(leadListProvider).selectedSource,
        onSelect: (source) {
          if (!mounted) return;
          ref.read(leadListProvider.notifier).updateSourceFilter(source);
          Navigator.pop(context);
        },
      ),
    );
  }

  void _showIndustryFilter() {
    if (!mounted) return;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => IndustryFilterSheet(
        selectedIndustry: ref.watch(leadListProvider).selectedIndustry,
        onSelect: (industry) {
          if (!mounted) return;
          ref.read(leadListProvider.notifier).updateIndustryFilter(industry);
          Navigator.pop(context);
        },
      ),
    );
  }

  void _showProvinceFilter() {
    if (!mounted) return;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => ProvinceFilterSheet(
        selectedProvince: ref.watch(leadListProvider).selectedProvince,
        onSelect: (province) {
          if (!mounted) return;
          ref.read(leadListProvider.notifier).updateProvinceFilter(province);
          Navigator.pop(context);
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final state = ref.watch(leadListProvider);
    final hasActiveFilters =
        state.selectedStatus.isNotEmpty ||
        state.selectedSource.isNotEmpty ||
        state.selectedIndustry.isNotEmpty ||
        state.selectedProvince.isNotEmpty;

    return MainScaffold(
      currentIndex: 2, // Dashboard=0, Pipeline=1, Leads=2, Route=3, Schedules=4
      title: l10n.leads,
      actions: [
        // Sync Status Indicator
        const Padding(
          padding: EdgeInsets.only(right: 8),
          child: Center(
            child: SyncStatusIndicator(featureKey: 'leads', iconSize: 20),
          ),
        ),
        IconButton(
          icon: const Icon(Icons.search),
          onPressed: () => showLeadSearchModal(context),
          tooltip: l10n.search,
        ),
        PopupMenuButton<String>(
          icon: Stack(
            children: [
              const Icon(Icons.filter_list),
              if (hasActiveFilters)
                Positioned(
                  right: 0,
                  top: 0,
                  child: Container(
                    width: 8,
                    height: 8,
                    decoration: BoxDecoration(
                      color: theme.colorScheme.error,
                      shape: BoxShape.circle,
                    ),
                  ),
                ),
            ],
          ),
          tooltip: l10n.filter,
          onSelected: (value) {
            switch (value) {
              case 'status':
                _showStatusFilter();
                break;
              case 'source':
                _showSourceFilter();
                break;
              case 'industry':
                _showIndustryFilter();
                break;
              case 'province':
                _showProvinceFilter();
                break;
              case 'clear':
                if (mounted) {
                  ref.read(leadListProvider.notifier).clearFilters();
                }
                break;
            }
          },
          itemBuilder: (context) => [
            PopupMenuItem(
              value: 'status',
              child: Row(
                children: [
                  const Icon(Icons.flag_outlined, size: 20),
                  const SizedBox(width: 12),
                  Text(l10n.status),
                ],
              ),
            ),
            PopupMenuItem(
              value: 'source',
              child: Row(
                children: [
                  const Icon(Icons.source_outlined, size: 20),
                  const SizedBox(width: 12),
                  Text(l10n.source),
                ],
              ),
            ),
            PopupMenuItem(
              value: 'industry',
              child: Row(
                children: [
                  const Icon(Icons.business_outlined, size: 20),
                  const SizedBox(width: 12),
                  Text(l10n.industry),
                ],
              ),
            ),
            PopupMenuItem(
              value: 'province',
              child: Row(
                children: [
                  const Icon(Icons.map_outlined, size: 20),
                  const SizedBox(width: 12),
                  Text(l10n.province),
                ],
              ),
            ),
            if (hasActiveFilters) ...[
              const PopupMenuDivider(),
              PopupMenuItem(
                value: 'clear',
                child: Row(
                  children: [
                    Icon(
                      Icons.clear_all,
                      size: 20,
                      color: theme.colorScheme.error,
                    ),
                    const SizedBox(width: 12),
                    Text(
                      l10n.clear,
                      style: TextStyle(color: theme.colorScheme.error),
                    ),
                  ],
                ),
              ),
            ],
          ],
        ),
      ],
      floatingActionButton: FloatingActionButton(
        elevation: 0.5,
        onPressed: () {
          // Navigate to create lead
          Navigator.pushNamed(context, AppRoutes.leadsForm);
        },
        child: const Icon(Icons.add),
      ),
      body: Column(
        children: [
          if (state.searchQuery.isNotEmpty)
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 8, 16, 0),
              child: Row(
                children: [
                  Expanded(
                    child: Text(
                      '${l10n.search}: "${state.searchQuery}"',
                      style: theme.textTheme.bodyMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  const SizedBox(width: 8),
                  TextButton.icon(
                    onPressed: () {
                      if (!mounted) return;
                      ref.read(leadListProvider.notifier).updateSearchQuery('');
                      ref
                          .read(leadListProvider.notifier)
                          .loadLeads(page: 1, refresh: true);
                    },
                    icon: Icon(
                      Icons.close,
                      size: 16,
                      color: theme.colorScheme.error,
                    ),
                    label: Text(
                      l10n.clear,
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
          Expanded(
            child: RefreshIndicator(
              onRefresh: () async {
                if (!mounted) return;
                await ref
                    .read(leadListProvider.notifier)
                    .loadLeads(page: 1, refresh: true);
              },
              child: _buildBody(context, state),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBody(BuildContext context, LeadListState state) {
    final l10n = AppLocalizations.of(context)!;

    if (state.isLoading && state.leads.isEmpty) {
      return const SingleChildScrollView(
        physics: AlwaysScrollableScrollPhysics(),
        child: Center(
          child: Padding(
            padding: EdgeInsets.all(100.0),
            child: CircularProgressIndicator(),
          ),
        ),
      );
    }

    if (state.error != null && state.leads.isEmpty) {
      return SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        child: Center(
          child: Padding(
            padding: const EdgeInsets.all(100.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(state.error!),
                const SizedBox(height: 16),
                FilledButton(
                  onPressed: () {
                    if (!mounted) return;
                    ref
                        .read(leadListProvider.notifier)
                        .loadLeads(page: 1, refresh: true);
                  },
                  child: Text(l10n.retry),
                ),
              ],
            ),
          ),
        ),
      );
    }

    if (state.leads.isEmpty) {
      return SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        child: Center(
          child: Padding(
            padding: const EdgeInsets.all(100.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(Icons.people_outline, size: 64, color: Colors.grey),
                const SizedBox(height: 16),
                Text(
                  l10n.noLeadsFound,
                  style: Theme.of(
                    context,
                  ).textTheme.titleMedium?.copyWith(color: Colors.grey),
                ),
              ],
            ),
          ),
        ),
      );
    }

    return ListView.builder(
      controller: _scrollController,
      physics: const AlwaysScrollableScrollPhysics(),
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 100),
      itemCount: state.leads.length + (state.isLoadingMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == state.leads.length) {
          return const Center(
            child: Padding(
              padding: EdgeInsets.all(16),
              child: CircularProgressIndicator(),
            ),
          );
        }

        final lead = state.leads[index];
        return LeadCard(
          lead: lead,
          onTap: () {
            Navigator.pushNamed(
              context,
              AppRoutes.leadsDetail,
              arguments: lead.id,
            );
          },
        );
      },
    );
  }
}
