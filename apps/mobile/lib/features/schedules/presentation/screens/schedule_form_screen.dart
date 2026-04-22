import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import '../../../../core/l10n/app_localizations.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../../../google_calendar/application/google_calendar_provider.dart'
    as google_calendar;
import '../../application/schedule_provider.dart';

class ScheduleFormScreen extends ConsumerStatefulWidget {
  final String? scheduleId;

  const ScheduleFormScreen({super.key, this.scheduleId});

  @override
  ConsumerState<ScheduleFormScreen> createState() => _ScheduleFormScreenState();
}

class _ScheduleFormScreenState extends ConsumerState<ScheduleFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _reminderMinutesBeforeController = TextEditingController();
  final _taskIdController = TextEditingController();
  DateTime _scheduledAt = DateTime.now();
  String _status = 'pending';
  bool _syncToCalendar = false;

  @override
  void initState() {
    super.initState();
    if (widget.scheduleId != null) {
      Future.microtask(() async {
        final schedule = await ref
            .read(scheduleFormProvider.notifier)
            .loadInitialData(widget.scheduleId!);
        if (schedule != null) {
          _titleController.text = schedule.title;
          _descriptionController.text = schedule.description ?? '';
          _reminderMinutesBeforeController.text =
              schedule.reminderMinutesBefore?.toString() ?? '';
          _taskIdController.text = schedule.taskId ?? '';
          setState(() {
            _scheduledAt = schedule.scheduledAt;
            _syncToCalendar = schedule.syncToCalendar;
            // Map common statuses or use direct value
            if ([
              'pending',
              'confirmed',
              'completed',
              'cancelled',
            ].contains(schedule.status)) {
              _status = schedule.status;
            } else {
              _status = 'pending';
            }
          });
        }
      });
    }
  }

  @override
  void dispose() {
    _titleController.dispose();
    _descriptionController.dispose();
    _reminderMinutesBeforeController.dispose();
    _taskIdController.dispose();
    super.dispose();
  }

  Future<void> _selectDateTime() async {
    final date = await showDatePicker(
      context: context,
      initialDate: _scheduledAt,
      firstDate: DateTime.now().subtract(const Duration(days: 365)),
      lastDate: DateTime.now().add(const Duration(days: 365)),
    );
    if (date != null) {
      if (!mounted) return;
      final time = await showTimePicker(
        context: context,
        initialTime: TimeOfDay.fromDateTime(_scheduledAt),
      );
      if (time != null) {
        setState(() {
          _scheduledAt = DateTime(
            date.year,
            date.month,
            date.day,
            time.hour,
            time.minute,
          );
        });
      }
    }
  }

  void _submit() {
    if (_formKey.currentState!.validate()) {
      final data = {
        'title': _titleController.text.trim(),
        'description': _descriptionController.text.trim(),
        'reminder_minutes_before':
            int.tryParse(_reminderMinutesBeforeController.text) ?? 30,
        'scheduled_at':
            '${_scheduledAt.toUtc().toIso8601String().split('.').first}Z',
        'status': _status,
        'task_id': _taskIdController.text.trim().isEmpty
            ? null
            : _taskIdController.text.trim(),
        'sync_to_calendar': _syncToCalendar,
      };

      ref.read(scheduleFormProvider.notifier).submit(data).then((success) {
        if (!mounted) return;
        if (success) {
          Navigator.pop(context, true);
        }
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(scheduleFormProvider);
    final calendarStatusAsync = ref.watch(
      google_calendar.googleCalendarStatusProvider,
    );
    final l10n = AppLocalizations.of(context)!;

    ref.listen(scheduleFormProvider, (previous, next) {
      if (next.errorMessage != null && mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(next.errorMessage!)));
      }
    });

    return Scaffold(
      appBar: AppBar(
        title: Text(
          widget.scheduleId == null ? l10n.createSchedule : l10n.editSchedule,
        ),
      ),
      body: state.isLoading
          ? const LoadingWidget()
          : Form(
              key: _formKey,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  TextFormField(
                    controller: _titleController,
                    decoration: InputDecoration(
                      labelText: '${l10n.title} *',
                      hintText: l10n.enterTitle,
                      border: const OutlineInputBorder(),
                    ),
                    validator: (v) {
                      if (v == null || v.trim().isEmpty) return l10n.required;
                      if (v.trim().length < 3) return l10n.minCharacters(3);
                      if (v.trim().length > 255) return l10n.maxCharacters(255);
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: TextEditingController(
                      text: DateFormat('yyyy-MM-dd HH:mm').format(_scheduledAt),
                    ),
                    decoration: InputDecoration(
                      labelText: '${l10n.scheduledAt} *',
                      suffixIcon: const Icon(Icons.calendar_today),
                      border: const OutlineInputBorder(),
                    ),
                    readOnly: true,
                    onTap: _selectDateTime,
                    validator: (v) =>
                        v == null || v.isEmpty ? l10n.required : null,
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _reminderMinutesBeforeController,
                    decoration: InputDecoration(
                      labelText: l10n.reminderMinutesBefore,
                      hintText: 'e.g., 30',
                      helperText: '${l10n.optional} (0 - 10080)',
                      border: const OutlineInputBorder(),
                    ),
                    keyboardType: TextInputType.number,
                    validator: (v) {
                      if (v == null || v.isEmpty) return null;
                      final min = int.tryParse(v);
                      if (min == null || min < 0 || min > 10080) {
                        return l10n.invalidReminderMinutes;
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _descriptionController,
                    decoration: InputDecoration(
                      labelText: l10n.description,
                      border: const OutlineInputBorder(),
                    ),
                    maxLines: 3,
                  ),
                  const SizedBox(height: 24),
                  // Google Calendar Sync Checkbox (only show if connected)
                  calendarStatusAsync.when(
                    data: (status) {
                      if (status.isConnected == true) {
                        return Column(
                          children: [
                            Card(
                              margin: EdgeInsets.zero,
                              child: Padding(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 12,
                                  vertical: 4,
                                ),
                                child: CheckboxListTile(
                                  value: _syncToCalendar,
                                  onChanged: (value) {
                                    setState(() {
                                      _syncToCalendar = value ?? false;
                                    });
                                  },
                                  title: Row(
                                    children: [
                                      Icon(
                                        Icons.sync,
                                        size: 20,
                                        color: Theme.of(
                                          context,
                                        ).colorScheme.primary,
                                      ),
                                      const SizedBox(width: 8),
                                      const Text('Sync to Google Calendar'),
                                    ],
                                  ),
                                  subtitle: const Text(
                                    'Add this schedule to your Google Calendar',
                                    style: TextStyle(fontSize: 12),
                                  ),
                                  contentPadding: EdgeInsets.zero,
                                  controlAffinity:
                                      ListTileControlAffinity.trailing,
                                  activeColor: Theme.of(
                                    context,
                                  ).colorScheme.primary,
                                ),
                              ),
                            ),
                            const SizedBox(height: 32),
                          ],
                        );
                      } else {
                        return Column(
                          children: [
                            // Show info when not connected
                            Card(
                              margin: EdgeInsets.zero,
                              child: Padding(
                                padding: const EdgeInsets.all(16),
                                child: Row(
                                  children: [
                                    Icon(
                                      Icons.info_outline,
                                      color: Theme.of(
                                        context,
                                      ).colorScheme.outline,
                                    ),
                                    const SizedBox(width: 12),
                                    Expanded(
                                      child: Column(
                                        crossAxisAlignment:
                                            CrossAxisAlignment.start,
                                        children: [
                                          Text(
                                            'Google Calendar',
                                            style: Theme.of(context)
                                                .textTheme
                                                .bodyMedium
                                                ?.copyWith(
                                                  fontWeight: FontWeight.w500,
                                                ),
                                          ),
                                          const SizedBox(height: 2),
                                          Text(
                                            'Connect Google Calendar in Profile settings to enable sync',
                                            style: Theme.of(context)
                                                .textTheme
                                                .bodySmall
                                                ?.copyWith(
                                                  color: Theme.of(
                                                    context,
                                                  ).colorScheme.outline,
                                                ),
                                          ),
                                        ],
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                            const SizedBox(height: 32),
                          ],
                        );
                      }
                    },
                    loading: () => const SizedBox.shrink(),
                    error: (_, _) => const SizedBox.shrink(),
                  ),
                  SizedBox(
                    width: double.infinity,
                    height: 50,
                    child: ElevatedButton(
                      onPressed: _submit,
                      child: Text(l10n.save),
                    ),
                  ),
                ],
              ),
            ),
    );
  }
}
