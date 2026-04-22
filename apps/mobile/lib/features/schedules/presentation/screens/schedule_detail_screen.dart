import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import '../../../../core/l10n/app_localizations.dart';
import '../../../../core/widgets/error_widget.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../../../google_calendar/application/google_calendar_provider.dart'
    as google_calendar;
import '../../../google_calendar/domain/models/google_calendar_status.dart';
import '../../data/models/schedule_model.dart';
import '../../application/schedule_provider.dart';

class ScheduleDetailScreen extends ConsumerStatefulWidget {
  final String scheduleId;

  const ScheduleDetailScreen({super.key, required this.scheduleId});

  @override
  ConsumerState<ScheduleDetailScreen> createState() =>
      _ScheduleDetailScreenState();
}

class _ScheduleDetailScreenState extends ConsumerState<ScheduleDetailScreen> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() {
      ref
          .read(scheduleDetailProvider(widget.scheduleId).notifier)
          .loadSchedule();
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(scheduleDetailProvider(widget.scheduleId));
    final calendarStatusAsync = ref.watch(
      google_calendar.googleCalendarStatusProvider,
    );
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.scheduleDetails),
        actions: [
          if (state.schedule != null)
            IconButton(
              icon: const Icon(Icons.edit),
              onPressed: () {
                Navigator.pushNamed(
                  context,
                  '/schedules/form',
                  arguments: widget.scheduleId,
                ).then((_) {
                  ref
                      .read(scheduleDetailProvider(widget.scheduleId).notifier)
                      .loadSchedule();
                });
              },
            ),
        ],
      ),
      body: state.isLoading
          ? const LoadingWidget()
          : state.errorMessage != null
          ? ErrorStateWidget(
              message: state.errorMessage!,
              onRetry: () => ref
                  .read(scheduleDetailProvider(widget.scheduleId).notifier)
                  .loadSchedule(),
            )
          : state.schedule == null
          ? Center(child: Text(l10n.noData))
          : _buildContent(state.schedule!, calendarStatusAsync),
    );
  }

  Widget _buildContent(
    ScheduleModel schedule,
    AsyncValue<GoogleCalendarStatus?> calendarStatusAsync,
  ) {
    final l10n = AppLocalizations.of(context)!;
    final dateFormat = DateFormat('EEEE, dd MMM yyyy HH:mm');

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildInfoCard(schedule, dateFormat),
          const SizedBox(height: 20),
          _buildGoogleCalendarSection(schedule, calendarStatusAsync),
          const SizedBox(height: 20),
          _buildSectionTitle(l10n.description),
          Text(schedule.description ?? l10n.noDescription),
          const SizedBox(height: 20),
          _buildSectionTitle(l10n.location),
          Text(schedule.location ?? l10n.noLocation),
          const SizedBox(height: 20),
          if (schedule.visitReportId != null) ...[
            _buildSectionTitle(l10n.relatedVisitReport),
            ListTile(
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.assignment),
              title: Text(l10n.viewVisitReport),
              trailing: const Icon(Icons.chevron_right),
              onTap: () {
                Navigator.pushNamed(
                  context,
                  '/visit-reports/detail',
                  arguments: schedule.visitReportId,
                );
              },
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildGoogleCalendarSection(
    ScheduleModel schedule,
    AsyncValue<GoogleCalendarStatus?> calendarStatusAsync,
  ) {
    final theme = Theme.of(context);

    return calendarStatusAsync.when(
      data: (status) {
        if (status?.isConnected != true) {
          // Not connected - show info card
          return Card(
            margin: EdgeInsets.zero,
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  Icon(Icons.info_outline, color: theme.colorScheme.outline),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Google Calendar',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          'Connect in Profile settings to sync schedules',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.outline,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          );
        }

        // Connected - show sync/unsync button
        final isSynced = schedule.syncToCalendar;

        return Card(
          margin: EdgeInsets.zero,
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(
                      isSynced ? Icons.check_circle : Icons.calendar_today,
                      color: isSynced
                          ? Colors.green
                          : theme.colorScheme.primary,
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Google Calendar',
                            style: theme.textTheme.bodyMedium?.copyWith(
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 2),
                          Text(
                            isSynced
                                ? 'Synced to your Google Calendar'
                                : 'Not synced to Google Calendar',
                            style: theme.textTheme.bodySmall?.copyWith(
                              color: theme.colorScheme.outline,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                SizedBox(
                  width: double.infinity,
                  child: OutlinedButton.icon(
                    onPressed: () => _toggleCalendarSync(schedule),
                    icon: Icon(isSynced ? Icons.link_off : Icons.sync),
                    label: Text(
                      isSynced
                          ? 'Unsync from Google Calendar'
                          : 'Sync to Google Calendar',
                    ),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: isSynced
                          ? Colors.red
                          : theme.colorScheme.primary,
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      },
      loading: () => const SizedBox.shrink(),
      error: (_, _) => const SizedBox.shrink(),
    );
  }

  Future<void> _toggleCalendarSync(ScheduleModel schedule) async {
    final notifier = ref.read(
      scheduleDetailProvider(widget.scheduleId).notifier,
    );
    final success = schedule.syncToCalendar
        ? await notifier.unsyncFromGoogleCalendar()
        : await notifier.syncToGoogleCalendar();

    if (!mounted) return;

    if (success) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            schedule.syncToCalendar
                ? 'Schedule unsynced from Google Calendar'
                : 'Schedule synced to Google Calendar',
          ),
          backgroundColor: Colors.green,
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: const Text('Failed to sync schedule'),
          backgroundColor: Theme.of(context).colorScheme.error,
        ),
      );
    }
  }

  Widget _buildInfoCard(ScheduleModel schedule, DateFormat dateFormat) {
    final l10n = AppLocalizations.of(context)!;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            _buildInfoRow(Icons.title, l10n.title, schedule.title),
            const Divider(),
            _buildInfoRow(
              Icons.calendar_today,
              l10n.scheduledAt,
              dateFormat.format(schedule.scheduledAt),
            ),
            const Divider(),
            _buildInfoRow(
              Icons.info_outline,
              l10n.status,
              schedule.status.toUpperCase(),
            ),
            const Divider(),
            _buildInfoRow(
              Icons.category,
              l10n.activityType,
              schedule.activityType ?? '-',
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(IconData icon, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Icon(icon, size: 20, color: Colors.grey),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: const TextStyle(fontSize: 12, color: Colors.grey),
                ),
                Text(
                  value,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionTitle(String title) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Text(
        title,
        style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
      ),
    );
  }
}
