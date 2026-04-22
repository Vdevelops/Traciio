import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../application/google_calendar_provider.dart';

/// Widget untuk menampilkan status koneksi Google Calendar di Profile
class GoogleCalendarConnectionWidget extends ConsumerStatefulWidget {
  const GoogleCalendarConnectionWidget({super.key});

  @override
  ConsumerState<GoogleCalendarConnectionWidget> createState() =>
      _GoogleCalendarConnectionWidgetState();
}

class _GoogleCalendarConnectionWidgetState
    extends ConsumerState<GoogleCalendarConnectionWidget>
    with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state == AppLifecycleState.resumed) {
      // Refresh status when app comes back to foreground
      ref.invalidate(googleCalendarNotifierProvider);
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(googleCalendarNotifierProvider);

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(
                  Icons.calendar_today,
                  color: Theme.of(context).primaryColor,
                ),
                const SizedBox(width: 8),
                Text(
                  'Google Calendar',
                  style: Theme.of(context).textTheme.titleMedium,
                ),
              ],
            ),
            const SizedBox(height: 16),
            state.when(
              data: (status) {
                if (status == null) {
                  return _buildDisconnectedView();
                }

                if (status.isConnected) {
                  return _buildConnectedView(status);
                } else {
                  return _buildDisconnectedView();
                }
              },
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (error, _) => _buildErrorView(error),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildConnectedView(dynamic status) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            const Icon(Icons.check_circle, color: Colors.green),
            const SizedBox(width: 8),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Connected',
                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                      color: Colors.green,
                    ),
                  ),
                  if (status.email != null)
                    Text(
                      status.email!,
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        if (status.connectedAt != null)
          Text(
            'Connected since: ${_formatDate(status.connectedAt!)}',
            style: Theme.of(context).textTheme.bodySmall,
          ),
        const SizedBox(height: 12),
        SizedBox(
          width: double.infinity,
          child: OutlinedButton.icon(
            onPressed: _disconnect,
            icon: const Icon(Icons.link_off),
            label: const Text('Disconnect'),
            style: OutlinedButton.styleFrom(foregroundColor: Colors.red),
          ),
        ),
      ],
    );
  }

  Widget _buildDisconnectedView() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            const Icon(Icons.link_off, color: Colors.grey),
            const SizedBox(width: 8),
            const Expanded(
              child: Text(
                'Not Connected',
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  color: Colors.grey,
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Text(
          'Connect your Google Calendar to sync schedules automatically.',
          style: Theme.of(context).textTheme.bodySmall,
        ),
        const SizedBox(height: 12),
        SizedBox(
          width: double.infinity,
          child: FilledButton.icon(
            onPressed: _connect,
            icon: const Icon(Icons.link),
            label: const Text('Connect Google Calendar'),
          ),
        ),
      ],
    );
  }

  Widget _buildErrorView(Object error) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            const Icon(Icons.error_outline, color: Colors.red),
            const SizedBox(width: 8),
            const Expanded(
              child: Text(
                'Error',
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  color: Colors.red,
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Text(
          'Failed to load Google Calendar status.',
          style: Theme.of(context).textTheme.bodySmall,
        ),
        const SizedBox(height: 12),
        SizedBox(
          width: double.infinity,
          child: OutlinedButton.icon(
            onPressed: () => ref
                .read(googleCalendarNotifierProvider.notifier)
                .refreshStatus(),
            icon: const Icon(Icons.refresh),
            label: const Text('Retry'),
          ),
        ),
      ],
    );
  }

  Future<void> _connect() async {
    try {
      await ref.read(googleCalendarNotifierProvider.notifier).connect();
      // Note: Actual connection happens in browser, user needs to return to app
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Opening Google authorization...'),
            duration: Duration(seconds: 2),
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
        );
      }
    }
  }

  Future<void> _disconnect() async {
    // Show confirmation dialog
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Disconnect Google Calendar?'),
        content: const Text(
          'This will remove the connection to your Google Calendar. '
          'Your existing synced events will remain in Google Calendar.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            style: FilledButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Disconnect'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      try {
        await ref.read(googleCalendarNotifierProvider.notifier).disconnect();
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Google Calendar disconnected')),
          );
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Error disconnecting: $e'),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    }
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year}';
  }
}
