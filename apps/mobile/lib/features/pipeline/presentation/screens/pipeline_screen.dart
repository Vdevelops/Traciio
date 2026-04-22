import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'dart:async';
import 'package:mobile/core/widgets/main_scaffold.dart';
import 'package:mobile/core/permissions/permission_provider.dart';
import 'package:mobile/core/widgets/sync_status_indicator.dart';
import 'package:mobile/features/pipeline/application/pipeline_provider.dart';
import 'package:mobile/features/pipeline/presentation/screens/deal_list_screen.dart';
import 'package:mobile/core/routing/app_router.dart';

class PipelineScreen extends ConsumerStatefulWidget {
  const PipelineScreen({super.key});

  @override
  ConsumerState<PipelineScreen> createState() => _PipelineScreenState();
}

class _PipelineScreenState extends ConsumerState<PipelineScreen>
    with SingleTickerProviderStateMixin {
  final TextEditingController _searchController = TextEditingController();
  Timer? _debounceTimer;
  TabController? _tabController;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;
      ref.read(pipelineProvider.notifier).loadStages();
    });
  }

  @override
  void dispose() {
    _searchController.dispose();
    _debounceTimer?.cancel();
    _tabController?.dispose();
    super.dispose();
  }

  void _onSearchChanged(String query) {
    _debounceTimer?.cancel();
    _debounceTimer = Timer(const Duration(milliseconds: 500), () {
      if (!mounted) return;
      ref.read(pipelineProvider.notifier).updateSearchQuery(query);
      ref
          .read(pipelineProvider.notifier)
          .loadDeals(page: 1, refresh: true, search: query);
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final state = ref.watch(pipelineProvider);
    final stages = state.stages;

    // Use CREATE permission for deals
    final hasCreatePermission = ref.watch(canCreateOpportunityProvider);

    // Initialize or update TabController when stages are loaded
    if (stages.isNotEmpty &&
        (_tabController == null || _tabController!.length != stages.length)) {
      _tabController?.dispose();
      _tabController = TabController(length: stages.length, vsync: this);
      _tabController!.addListener(() {
        if (!mounted) return;
        if (!_tabController!.indexIsChanging) {
          final stageId = stages[_tabController!.index].id;
          ref.read(pipelineProvider.notifier).selectStage(stageId);
          if (mounted) setState(() {});
        }
      });
    } else if (stages.isNotEmpty &&
        _tabController != null &&
        !_tabController!.indexIsChanging) {
      // Sync TabController with state.selectedStageId
      // This handles cases where selectedStageId is updated externally (e.g. invalid stage reset)
      // We check !indexIsChanging so we don't fight against user interactions
      final selectedStageId = state.selectedStageId;
      if (selectedStageId != null) {
        final index = stages.indexWhere((s) => s.id == selectedStageId);
        if (index != -1 && index != _tabController!.index) {
          _tabController!.animateTo(index);
        }
      }
    }

    return MainScaffold(
      currentIndex: 1, // Dashboard=0, Pipeline=1, Leads=2, Route=3, Schedules=4
      title: null,
      actions: const [
        // Sync Status Indicator
        Padding(
          padding: EdgeInsets.only(right: 16),
          child: Center(
            child: SyncStatusIndicator(featureKey: 'pipeline', iconSize: 20),
          ),
        ),
      ],
      floatingActionButton: hasCreatePermission
          ? FloatingActionButton(
              elevation: 0.5,
              onPressed: () {
                Navigator.pushNamed(context, '${AppRoutes.pipeline}/create');
              },
              child: const Icon(Icons.add),
            )
          : null,
      body: NestedScrollView(
        headerSliverBuilder: (context, innerBoxIsScrolled) {
          return [
            SliverToBoxAdapter(
              child: SafeArea(
                bottom: false,
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(20, 16, 20, 8),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Container(
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
                            hintText: 'Search deals...',
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
                    ],
                  ),
                ),
              ),
            ),
            if (stages.isNotEmpty)
              SliverToBoxAdapter(
                child: Container(
                  height: 60,
                  color: theme.scaffoldBackgroundColor,
                  padding: const EdgeInsets.symmetric(vertical: 8),
                  child: SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Row(
                      children: List.generate(stages.length, (index) {
                        final stage = stages[index];
                        final isSelected = _tabController?.index == index;
                        return Padding(
                          padding: const EdgeInsets.only(right: 8),
                          child: _buildTab(context, stage.name, isSelected, () {
                            _tabController?.animateTo(index);
                            ref
                                .read(pipelineProvider.notifier)
                                .selectStage(stage.id);
                            if (mounted) setState(() {});
                          }),
                        );
                      }),
                    ),
                  ),
                ),
              ),
          ];
        },
        body: stages.isEmpty
            ? (state.isLoading
                  ? const Center(child: CircularProgressIndicator())
                  : const SizedBox())
            : TabBarView(
                controller: _tabController,
                children: stages
                    .map((stage) => DealListScreen(stageId: stage.id))
                    .toList(),
              ),
      ),
    );
  }

  Widget _buildTab(
    BuildContext context,
    String label,
    bool isSelected,
    VoidCallback onTap,
  ) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 10),
          decoration: BoxDecoration(
            color: isSelected
                ? colorScheme.primary
                : colorScheme.surfaceContainerHighest.withValues(alpha: 0.5),
            borderRadius: BorderRadius.circular(16),
            border: isSelected
                ? null
                : Border.all(
                    color: colorScheme.outline.withValues(alpha: 0.2),
                    width: 1,
                  ),
          ),
          child: Text(
            label,
            style: theme.textTheme.labelLarge?.copyWith(
              color: isSelected
                  ? colorScheme.onPrimary
                  : colorScheme.onSurface.withValues(alpha: 0.7),
              fontWeight: isSelected ? FontWeight.bold : FontWeight.w500,
            ),
          ),
        ),
      ),
    );
  }
}
