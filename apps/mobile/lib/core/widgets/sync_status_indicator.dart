import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../sync/auto_sync_manager.dart';

/// Widget untuk menampilkan status sync di AppBar atau header
///
/// Menampilkan:
/// - Loading indicator saat sedang sync
/// - Tidak menampilkan apa-apa saat idle (refresh button dihapus)
class SyncStatusIndicator extends ConsumerWidget {
  final String? featureKey;
  final double iconSize;

  const SyncStatusIndicator({super.key, this.featureKey, this.iconSize = 20});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    // Jika featureKey di-set, check specific feature
    // Jika tidak, check global sync status
    final bool isSyncing = featureKey != null
        ? ref.watch(isFeatureSyncingProvider(featureKey!))
        : ref.watch(isAnySyncingProvider);

    return AnimatedSwitcher(
      duration: const Duration(milliseconds: 200),
      child: isSyncing
          ? SizedBox(
              key: const ValueKey('syncing'),
              width: iconSize,
              height: iconSize,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                valueColor: AlwaysStoppedAnimation<Color>(
                  theme.colorScheme.primary,
                ),
              ),
            )
          : const SizedBox.shrink(),
    );
  }
}

/// Widget untuk menampilkan offline indicator
class OfflineIndicator extends ConsumerWidget {
  const OfflineIndicator({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Listen ke connectivity service
    // Note: This needs to be implemented based on your connectivity service
    // For now, returning a simple widget

    return Container(
      color: Colors.orange,
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: const Row(
        children: [
          Icon(Icons.wifi_off, color: Colors.white, size: 16),
          SizedBox(width: 8),
          Expanded(
            child: Text(
              'Offline Mode - Showing cached data',
              style: TextStyle(color: Colors.white, fontSize: 12),
            ),
          ),
        ],
      ),
    );
  }
}
