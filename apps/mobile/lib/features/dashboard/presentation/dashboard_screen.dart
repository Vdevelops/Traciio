import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/widgets/main_scaffold.dart';
import '../../auth/application/auth_provider.dart';
import '../../notifications/application/notification_provider.dart';
import '../../profile/application/profile_provider.dart';
import '../application/dashboard_provider.dart';
import 'widgets/dashboard_header.dart';
import 'widgets/simplified_dashboard_content.dart';
import '../../tasks/presentation/widgets/task_search_modal.dart';
import 'dart:async';

class DashboardScreen extends ConsumerStatefulWidget {
  const DashboardScreen({super.key});

  @override
  ConsumerState<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends ConsumerState<DashboardScreen> {
  Timer? _debounceTimer;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) {
        ref.read(dashboardProvider.notifier).loadDashboard();
        ref.read(notificationCountProvider.notifier).loadUnreadCount();
        // Force refresh profile to get latest avatar URL
        ref.invalidate(profileProvider);
      }
    });
  }

  @override
  void dispose() {
    _debounceTimer?.cancel();
    super.dispose();
  }

  Future<void> _onRefresh() async {
    if (!mounted) return;
    await ref.read(dashboardProvider.notifier).refresh();
    if (!mounted) return;
    ref.read(notificationCountProvider.notifier).loadUnreadCount();
  }

  @override
  Widget build(BuildContext context) {
    ref.watch(authProvider);
    ref.watch(notificationCountProvider);

    return MainScaffold(
      // Phase 2: Remove AppBar
      title: null,
      currentIndex: 0,
      body: SafeArea(
        bottom: false,
        child: RefreshIndicator(
          onRefresh: _onRefresh,
          child: SingleChildScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            child: Column(
              children: [
                // 1. DASHBOARD HEADER (Top Navigation) - Tidak fixed, scroll dengan konten
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 16, 20, 12),
                  child: DashboardHeader(
                    currentTabIndex: 0,
                    onSearchTap: () {
                      // Show search overlay untuk tasks
                      showTaskSearchModal(context);
                    },
                  ),
                ),

                // 2. DASHBOARD CONTENT (includes greeting header, performance goal, quick nav, tasks)
                const SimplifiedDashboardContent(),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
