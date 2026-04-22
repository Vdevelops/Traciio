import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../network/connectivity_service.dart';

/// Widget to show offline mode indicator
/// 
/// Displays a banner at the top of the screen when device is offline
class OfflineIndicator extends ConsumerWidget {
  const OfflineIndicator({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final connectivityService = ref.watch(connectivityServiceProvider);
    
    return StreamBuilder<bool>(
      stream: connectivityService.onConnectivityChanged,
      initialData: connectivityService.isOnline,
      builder: (context, snapshot) {
        final isOnline = snapshot.data ?? true;
        
        if (!isOnline) {
          return Container(
            width: double.infinity,
            color: Colors.orange,
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(
              children: [
                const Icon(Icons.wifi_off, color: Colors.white, size: 20),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    'Offline Mode - Showing cached data',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: Colors.white,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ),
              ],
            ),
          );
        }
        
        return const SizedBox.shrink();
      },
    );
  }
}


