import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/core/widgets/error_widget.dart';
import 'package:mobile/core/widgets/loading_widget.dart';
import 'package:mobile/features/pipeline/application/pipeline_provider.dart';
import 'package:mobile/features/pipeline/application/pipeline_state.dart';
import 'package:mobile/features/pipeline/presentation/widgets/deal_card.dart';
import 'package:mobile/core/routing/app_router.dart';

class DealListScreen extends ConsumerStatefulWidget {
  const DealListScreen({super.key, required this.stageId});

  final String stageId;

  @override
  ConsumerState<DealListScreen> createState() => _DealListScreenState();
}

class _DealListScreenState extends ConsumerState<DealListScreen> {
  bool _onScrollNotification(ScrollNotification notification) {
    if (notification is ScrollUpdateNotification) {
      if (notification.metrics.pixels >=
          notification.metrics.maxScrollExtent * 0.8) {
        ref.read(pipelineProvider.notifier).loadMore();
      }
    }
    return false;
  }

  Future<void> _onRefresh() async {
    await ref.read(pipelineProvider.notifier).refresh();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(pipelineProvider);
    final theme = Theme.of(context);

    // Verify that this screen corresponds to the currently selected stage
    // This prevents multiple DealListScreens in a TabBarView from overwriting the shared state
    if (state.selectedStageId != widget.stageId) {
      return const LoadingWidget();
    }

    return NotificationListener<ScrollNotification>(
      onNotification: _onScrollNotification,
      child: RefreshIndicator(
        onRefresh: _onRefresh,
        child: _buildContent(context, state, theme),
      ),
    );
  }

  Widget _buildContent(
    BuildContext context,
    PipelineState state,
    ThemeData theme,
  ) {
    if (state.isLoading && state.deals.isEmpty) {
      return const LoadingWidget();
    }

    if (state.errorMessage != null && state.deals.isEmpty) {
      return ErrorStateWidget(
        message: state.errorMessage!,
        onRetry: () {
          ref.read(pipelineProvider.notifier).refresh();
        },
      );
    }

    if (state.deals.isEmpty) {
      return const EmptyStateWidget(
        message: 'No deals found in this stage',
        subtitle: 'Tap the + button to create a new deal',
        icon: Icons.handshake_outlined,
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.fromLTRB(20, 16, 20, 100),
      itemCount: state.deals.length + (state.isLoadingMore ? 1 : 0),
      addAutomaticKeepAlives: false,
      addRepaintBoundaries: true,
      cacheExtent: 500,
      itemBuilder: (context, index) {
        if (index == state.deals.length) {
          return const Padding(
            padding: EdgeInsets.all(16),
            child: LoadingWidget(size: 24),
          );
        }

        final deal = state.deals[index];
        return DealCard(
          key: ValueKey(deal.id),
          deal: deal,
          onTap: () {
            Navigator.pushNamed(
              context,
              '${AppRoutes.pipeline}/${deal.id}',
              arguments: deal.id,
            );
          },
        );
      },
    );
  }
}
