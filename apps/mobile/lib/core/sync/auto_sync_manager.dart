import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../network/connectivity_service.dart';

/// Callback signature untuk sync operation
typedef SyncCallback = Future<void> Function();

/// State untuk AutoSyncManager
class AutoSyncState {
  final Map<String, bool> isSyncing;
  final Map<String, DateTime?> lastSyncTime;

  const AutoSyncState({
    this.isSyncing = const {},
    this.lastSyncTime = const {},
  });

  AutoSyncState copyWith({
    Map<String, bool>? isSyncing,
    Map<String, DateTime?>? lastSyncTime,
  }) {
    return AutoSyncState(
      isSyncing: isSyncing ?? this.isSyncing,
      lastSyncTime: lastSyncTime ?? this.lastSyncTime,
    );
  }

  bool get isAnySyncing => isSyncing.values.any((syncing) => syncing);
}

/// Manager untuk handle auto-sync secara centralized
///
/// Features:
/// - Auto sync saat app resume
/// - Auto sync saat koneksi kembali
/// - Debounce untuk mencegah sync berulang kali
/// - Track sync status untuk semua features
class AutoSyncManager extends Notifier<AutoSyncState> {
  final Map<String, SyncCallback> _syncCallbacks = {};
  Timer? _debounceTimer;

  static const _debounceDuration = Duration(seconds: 2);

  @override
  AutoSyncState build() {
    _initConnectivityListener();
    return const AutoSyncState();
  }

  /// Initialize connectivity listener
  void _initConnectivityListener() {
    final connectivity = ref.read(connectivityServiceProvider);
    connectivity.onConnectivityChanged.listen((isOnline) {
      if (isOnline) {
        debugPrint('AutoSyncManager: Connection restored, triggering sync');
        syncAll();
      }
    });
  }

  /// Register sebuah feature untuk auto-sync
  ///
  /// [featureKey]: unique identifier (e.g., 'schedules', 'leads', 'accounts')
  /// [callback]: function yang akan dipanggil untuk sync
  void registerFeature(String featureKey, SyncCallback callback) {
    _syncCallbacks[featureKey] = callback;
    state = state.copyWith(
      isSyncing: {...state.isSyncing, featureKey: false},
      lastSyncTime: {...state.lastSyncTime, featureKey: null},
    );
  }

  /// Unregister feature
  void unregisterFeature(String featureKey) {
    _syncCallbacks.remove(featureKey);
    final newIsSyncing = Map<String, bool>.from(state.isSyncing);
    final newLastSyncTime = Map<String, DateTime?>.from(state.lastSyncTime);
    newIsSyncing.remove(featureKey);
    newLastSyncTime.remove(featureKey);
    state = state.copyWith(
      isSyncing: newIsSyncing,
      lastSyncTime: newLastSyncTime,
    );
  }

  /// Check apakah feature sedang syncing
  bool isSyncing(String featureKey) {
    return state.isSyncing[featureKey] ?? false;
  }

  /// Get last sync time untuk sebuah feature
  DateTime? getLastSyncTime(String featureKey) {
    return state.lastSyncTime[featureKey];
  }

  /// Trigger sync untuk semua registered features
  ///
  /// Dipanggil saat:
  /// - App resume dari background
  /// - Koneksi internet kembali
  /// - User manual refresh
  void syncAll() {
    // Cancel debounce timer yang ada
    _debounceTimer?.cancel();

    // Debounce untuk mencegah multiple sync berulang kali
    _debounceTimer = Timer(_debounceDuration, () async {
      await _performSyncAll();
    });
  }

  /// Sync specific feature only
  Future<void> syncFeature(String featureKey) async {
    final callback = _syncCallbacks[featureKey];
    if (callback == null) return;
    if (state.isSyncing[featureKey] == true) return;

    // Mark as syncing
    state = state.copyWith(isSyncing: {...state.isSyncing, featureKey: true});

    try {
      await callback();

      // Update last sync time
      state = state.copyWith(
        isSyncing: {...state.isSyncing, featureKey: false},
        lastSyncTime: {...state.lastSyncTime, featureKey: DateTime.now()},
      );
    } catch (e) {
      debugPrint('AutoSyncManager: Error syncing $featureKey: $e');
      state = state.copyWith(
        isSyncing: {...state.isSyncing, featureKey: false},
      );
    }
  }

  /// Perform sync all features
  Future<void> _performSyncAll() async {
    final connectivity = ref.read(connectivityServiceProvider);
    if (!connectivity.isOnline) {
      debugPrint('AutoSyncManager: Skip sync - offline');
      return;
    }

    debugPrint(
      'AutoSyncManager: Starting sync for ${_syncCallbacks.length} features',
    );

    // Sync semua features secara parallel
    final futures = _syncCallbacks.entries.map((entry) async {
      await syncFeature(entry.key);
    });

    await Future.wait(futures);

    debugPrint('AutoSyncManager: Sync completed');
  }
}

/// Provider untuk AutoSyncManager
final autoSyncManagerProvider =
    NotifierProvider<AutoSyncManager, AutoSyncState>(AutoSyncManager.new);

/// Provider untuk track apakah ada feature yang sedang sync
final isAnySyncingProvider = Provider<bool>((ref) {
  final state = ref.watch(autoSyncManagerProvider);
  return state.isAnySyncing;
});

/// Provider untuk check sync status specific feature
final isFeatureSyncingProvider = Provider.family<bool, String>((
  ref,
  featureKey,
) {
  final state = ref.watch(autoSyncManagerProvider);
  return state.isSyncing[featureKey] ?? false;
});

/// Provider untuk get last sync time
final lastSyncTimeProvider = Provider.family<DateTime?, String>((
  ref,
  featureKey,
) {
  final manager = ref.watch(autoSyncManagerProvider.notifier);
  return manager.getLastSyncTime(featureKey);
});
