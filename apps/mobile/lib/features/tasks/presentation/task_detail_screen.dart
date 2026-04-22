import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';

import '../application/task_provider.dart';
import '../data/models/task.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/widgets/error_widget.dart';
import '../../../core/widgets/loading_widget.dart';
import '../../../core/permissions/permission_provider.dart';

class TaskDetailScreen extends ConsumerStatefulWidget {
  const TaskDetailScreen({super.key, required this.taskId});

  final String taskId;

  @override
  ConsumerState<TaskDetailScreen> createState() => _TaskDetailScreenState();
}

class _TaskDetailScreenState extends ConsumerState<TaskDetailScreen> {
  // Delete functionality removed - sales users don't need to delete tasks

  @override
  void initState() {
    super.initState();
    // Force refresh task detail when screen opens to ensure latest data
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.invalidate(taskDetailProvider(widget.taskId));
    });
  }

  @override
  Widget build(BuildContext context) {
    final taskAsync = ref.watch(taskDetailProvider(widget.taskId));
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(title: Text(l10n.taskDetails), elevation: 0),
      body: taskAsync.when(
        loading: () => const LoadingWidget(),
        error: (error, stack) => ErrorStateWidget(
          message: error.toString().replaceFirst('Exception: ', ''),
          onRetry: () => ref.invalidate(taskDetailProvider(widget.taskId)),
        ),
        data: (task) =>
            _buildTaskDetail(context, task, theme, colorScheme, l10n),
      ),
    );
  }

  Widget _buildTaskDetail(
    BuildContext context,
    Task task,
    ThemeData theme,
    ColorScheme colorScheme,
    AppLocalizations l10n,
  ) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header Card
          Container(
            decoration: BoxDecoration(
              color: colorScheme.surface,
              borderRadius: BorderRadius.circular(20),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.05),
                  blurRadius: 3,
                  offset: const Offset(0, 1),
                ),
              ],
            ),
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        task.title,
                        style: theme.textTheme.headlineSmall?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: colorScheme.onSurface,
                        ),
                      ),
                    ),
                    _StatusBadge(
                      status: task.status,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),
          // Task Information
          _SectionTitle(
            title: l10n.taskInformation,
            theme: theme,
            colorScheme: colorScheme,
          ),
          const SizedBox(height: 8),
          _InfoCard(
            theme: theme,
            colorScheme: colorScheme,
            children: [
              _InfoRow(
                icon: Icons.check_box_outlined,
                label: l10n.title,
                value: task.title,
                theme: theme,
                colorScheme: colorScheme,
              ),
              if (task.description != null && task.description!.isNotEmpty)
                _InfoRow(
                  icon: Icons.description_outlined,
                  label: l10n.description,
                  value: task.description!,
                  theme: theme,
                  colorScheme: colorScheme,
                ),
              _InfoRow(
                icon: Icons.info_outlined,
                label: l10n.status,
                value: task.status.toUpperCase().replaceAll('_', ' '),
                theme: theme,
                colorScheme: colorScheme,
                badge: _StatusBadge(
                  status: task.status,
                  theme: theme,
                  colorScheme: colorScheme,
                ),
              ),
              _InfoRow(
                icon: Icons.flag_outlined,
                label: l10n.priority,
                value: task.priority.toUpperCase(),
                theme: theme,
                colorScheme: colorScheme,
              ),
              _InfoRow(
                icon: Icons.label_outlined,
                label: l10n.type,
                value: task.type.toUpperCase().replaceAll('_', ' '),
                theme: theme,
                colorScheme: colorScheme,
              ),
              if (task.dueDate != null)
                _InfoRow(
                  icon: Icons.calendar_today_outlined,
                  label: l10n.dueDate,
                  value: _formatDueDate(task.dueDate!),
                  theme: theme,
                  colorScheme: colorScheme,
                  textColor: task.isOverdue
                      ? colorScheme.error
                      : (task.isDueToday ? Colors.orange : null),
                ),
            ],
          ),
          const SizedBox(height: 16),
          // Related Information
          if (task.account != null || task.contact != null) ...[
            _SectionTitle(
              title: l10n.relatedInformation,
              theme: theme,
              colorScheme: colorScheme,
            ),
            const SizedBox(height: 8),
            _InfoCard(
              theme: theme,
              colorScheme: colorScheme,
              children: [
                if (task.account != null)
                  _InfoRow(
                    icon: Icons.business_outlined,
                    label: l10n.accounts,
                    value: task.account!.name,
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                if (task.contact != null)
                  _InfoRow(
                    icon: Icons.person_outline,
                    label: l10n.contacts,
                    value: task.contact!.name,
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
              ],
            ),
            const SizedBox(height: 16),
          ],
          // Reminders Section
          if (task.reminders.isNotEmpty) ...[
            _SectionTitle(
              title: l10n.reminders,
              theme: theme,
              colorScheme: colorScheme,
            ),
            const SizedBox(height: 8),
            _InfoCard(
              theme: theme,
              colorScheme: colorScheme,
              children: task.reminders.map((reminder) {
                return _ReminderRow(
                  reminder: reminder,
                  theme: theme,
                  colorScheme: colorScheme,
                  l10n: l10n,
                );
              }).toList(),
            ),
            const SizedBox(height: 16),
          ],
          // Action Buttons
          // Check permissions before showing action buttons
          Builder(
            builder: (context) {
              final canComplete = ref.watch(canCompleteTaskProvider);

              return Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  if (task.status != 'completed' &&
                      task.status != 'cancelled') ...[
                    // Mark In Progress button (only for pending tasks)
                    if (task.status.toLowerCase() == 'pending') ...[
                      FilledButton.icon(
                        onPressed: () => _showMarkInProgressDialog(
                          context,
                          task,
                          l10n,
                          colorScheme,
                        ),
                        icon: const Icon(Icons.play_arrow),
                        label: Text(l10n.markInProgress),
                        style: FilledButton.styleFrom(
                          minimumSize: const Size(double.infinity, 48),
                          backgroundColor: colorScheme.primary,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                      ),
                      const SizedBox(height: 12),
                    ],
                    // Complete button (only for in_progress tasks, not for completed)
                    if (canComplete &&
                        task.status.toLowerCase() != 'completed') ...[
                      FilledButton.icon(
                        onPressed: () => _showCompleteDialog(
                          context,
                          task,
                          l10n,
                          colorScheme,
                        ),
                        icon: const Icon(Icons.check),
                        label: Text(l10n.completeTask),
                        style: FilledButton.styleFrom(
                          minimumSize: const Size(double.infinity, 48),
                          backgroundColor: Colors.green,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                      ),
                      const SizedBox(height: 12),
                    ],
                  ],
                  OutlinedButton.icon(
                    onPressed: () => _showAddReminderDialog(
                      context,
                      task,
                      l10n,
                      colorScheme,
                    ),
                    icon: const Icon(Icons.notifications_outlined),
                    label: Text(l10n.addReminder),
                    style: OutlinedButton.styleFrom(
                      minimumSize: const Size(double.infinity, 48),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  ),
                ],
              );
            },
          ),
        ],
      ),
    );
  }

  // Edit functionality removed - sales users don't need to edit tasks

  String _formatDueDate(DateTime date) {
    return DateFormat('MMM dd, yyyy • HH:mm').format(date);
  }

  void _showMarkInProgressDialog(
    BuildContext context,
    Task task,
    AppLocalizations l10n,
    ColorScheme colorScheme,
  ) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.markInProgress),
        content: Text(l10n.markInProgressConfirmation),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(l10n.cancel),
          ),
          FilledButton(
            onPressed: () async {
              final navigator = Navigator.of(context);
              final scaffoldMessenger = ScaffoldMessenger.of(context);
              navigator.pop();
              try {
                await ref
                    .read(taskFormProvider.notifier)
                    .markInProgress(task.id);
                if (!mounted) return;
                scaffoldMessenger.showSnackBar(
                  SnackBar(
                    content: Text(l10n.taskMarkedInProgress),
                    backgroundColor: Colors.green,
                  ),
                );
                // Refresh detail
                ref.invalidate(taskDetailProvider(widget.taskId));
                // Return result to trigger refresh in list screen
                if (mounted) {
                  navigator.pop(true);
                }
              } catch (e) {
                if (!mounted) return;
                final error = ref.read(taskFormProvider).errorMessage;
                scaffoldMessenger.showSnackBar(
                  SnackBar(
                    content: Text(error ?? l10n.failedToMarkInProgress),
                    backgroundColor: Colors.red,
                  ),
                );
              }
            },
            child: Text(l10n.markInProgress),
          ),
        ],
      ),
    );
  }

  void _showCompleteDialog(
    BuildContext context,
    Task task,
    AppLocalizations l10n,
    ColorScheme colorScheme,
  ) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.completeTask),
        content: Text(l10n.completeTaskConfirmation),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(l10n.cancel),
          ),
          FilledButton(
            onPressed: () async {
              final navigator = Navigator.of(context);
              final scaffoldMessenger = ScaffoldMessenger.of(context);
              navigator.pop();
              final success = await ref
                  .read(taskFormProvider.notifier)
                  .completeTask(task.id);
              if (!mounted) return;
              if (success) {
                scaffoldMessenger.showSnackBar(
                  SnackBar(
                    content: Text(l10n.taskCompletedSuccessfully),
                    backgroundColor: Colors.green,
                  ),
                );
                // Refresh detail
                ref.invalidate(taskDetailProvider(widget.taskId));
                // Return result to trigger refresh in list screen
                if (mounted) {
                  navigator.pop(true);
                }
              } else {
                final error = ref.read(taskFormProvider).errorMessage;
                scaffoldMessenger.showSnackBar(
                  SnackBar(
                    content: Text(error ?? l10n.failedToCompleteTask),
                    backgroundColor: Colors.red,
                  ),
                );
              }
            },
            child: Text(l10n.completeTask),
          ),
        ],
      ),
    );
  }

  // Delete dialog removed - sales users don't need to delete tasks

  void _showAddReminderDialog(
    BuildContext context,
    Task task,
    AppLocalizations l10n,
    ColorScheme colorScheme,
  ) async {
    DateTime? selectedDate;
    final messageController = TextEditingController();
    final messenger = ScaffoldMessenger.of(context);

    final result = await showDialog<Map<String, dynamic>>(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: Text(l10n.addReminder),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                InkWell(
                  onTap: () async {
                    final dialogContext = context;
                    final date = await showDatePicker(
                      context: dialogContext,
                      initialDate:
                          selectedDate ??
                          DateTime.now().add(const Duration(days: 1)),
                      firstDate: DateTime.now(),
                      lastDate: DateTime(2100),
                    );
                    if (date != null && dialogContext.mounted) {
                      final time = await showTimePicker(
                        context: dialogContext,
                        initialTime: selectedDate != null
                            ? TimeOfDay.fromDateTime(selectedDate!)
                            : TimeOfDay.now(),
                      );
                      if (time != null) {
                        setDialogState(() {
                          selectedDate = DateTime(
                            date.year,
                            date.month,
                            date.day,
                            time.hour,
                            time.minute,
                          );
                        });
                      } else {
                        // If time picker is cancelled, still set the date with current time
                        final now = DateTime.now();
                        setDialogState(() {
                          selectedDate = DateTime(
                            date.year,
                            date.month,
                            date.day,
                            now.hour,
                            now.minute,
                          );
                        });
                      }
                    }
                  },
                  child: InputDecorator(
                    decoration: InputDecoration(
                      labelText: l10n.selectReminderDate,
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      filled: true,
                      fillColor: colorScheme.surfaceContainerHighest,
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 16,
                      ),
                      suffixIcon: Icon(
                        Icons.calendar_today_outlined,
                        color: colorScheme.onSurface.withValues(alpha: 0.7),
                      ),
                    ),
                    child: Text(
                      selectedDate != null
                          ? DateFormat(
                              'MMM dd, yyyy • HH:mm',
                            ).format(selectedDate!)
                          : l10n.selectReminderDate,
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: selectedDate != null
                            ? colorScheme.onSurface
                            : colorScheme.onSurface.withValues(alpha: 0.5),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: messageController,
                  decoration: InputDecoration(
                    labelText: l10n.reminderMessage,
                    hintText: l10n.enterReminderMessage,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    filled: true,
                    fillColor: colorScheme.surfaceContainerHighest,
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 16,
                    ),
                  ),
                  maxLines: 3,
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pop(context, null);
              },
              child: Text(l10n.cancel),
            ),
            FilledButton(
              onPressed: selectedDate != null
                  ? () {
                      // Get message text before closing dialog
                      final message = messageController.text.trim();
                      Navigator.pop(context, {
                        'date': selectedDate,
                        'message': message.isEmpty ? null : message,
                      });
                    }
                  : null,
              child: Text(l10n.save),
            ),
          ],
        ),
      ),
    );

    // Always dispose controller after dialog is closed
    try {
      if (!mounted) return;
      if (result != null && result['date'] != null) {
        final remindAt = result['date'] as DateTime;
        final message = result['message'] as String?;

        final reminder = await ref
            .read(taskFormProvider.notifier)
            .createReminder(
              taskId: task.id,
              remindAt: remindAt,
              message: message,
            );

        if (!mounted) return;
        if (reminder != null) {
          messenger.showSnackBar(
            SnackBar(
              content: Text(l10n.reminderCreatedSuccessfully),
              backgroundColor: Colors.green,
            ),
          );
          // Refresh detail to show new reminder
          ref.invalidate(taskDetailProvider(widget.taskId));
        } else {
          final error = ref.read(taskFormProvider).errorMessage;
          messenger.showSnackBar(
            SnackBar(
              content: Text(error ?? l10n.failedToCreateReminder),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    } finally {
      // Always dispose controller in finally block to ensure it's disposed
      messageController.dispose();
    }
  }
}

class _SectionTitle extends StatelessWidget {
  const _SectionTitle({
    required this.title,
    required this.theme,
    required this.colorScheme,
  });

  final String title;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    return Text(
      title,
      style: theme.textTheme.titleMedium?.copyWith(
        fontWeight: FontWeight.w600,
        color: colorScheme.onSurface,
      ),
    );
  }
}

class _InfoCard extends StatelessWidget {
  const _InfoCard({
    required this.children,
    required this.theme,
    required this.colorScheme,
  });

  final List<Widget> children;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: colorScheme.surface,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 3,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      padding: const EdgeInsets.all(16),
      child: Column(
        children: children.asMap().entries.map((entry) {
          final index = entry.key;
          final child = entry.value;
          if (index == children.length - 1) {
            return child;
          }
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: child,
          );
        }).toList(),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  const _InfoRow({
    required this.icon,
    required this.label,
    required this.value,
    required this.theme,
    required this.colorScheme,
    this.textColor,
    this.badge,
  });

  final IconData icon;
  final String label;
  final String value;
  final ThemeData theme;
  final ColorScheme colorScheme;
  final Color? textColor;
  final Widget? badge;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(
          icon,
          size: 20,
          color: colorScheme.onSurface.withValues(alpha: 0.7),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: theme.textTheme.bodySmall?.copyWith(
                  color: colorScheme.onSurface.withValues(alpha: 0.7),
                ),
              ),
              const SizedBox(height: 4),
              if (badge != null)
                badge!
              else
                Text(
                  value,
                  style: theme.textTheme.bodyMedium?.copyWith(
                    fontWeight: FontWeight.w500,
                    color: textColor ?? colorScheme.onSurface,
                  ),
                ),
            ],
          ),
        ),
      ],
    );
  }
}

class _ReminderRow extends StatelessWidget {
  const _ReminderRow({
    required this.reminder,
    required this.theme,
    required this.colorScheme,
    required this.l10n,
  });

  final Reminder reminder;
  final ThemeData theme;
  final ColorScheme colorScheme;
  final AppLocalizations l10n;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(
          Icons.notifications_outlined,
          size: 20,
          color: reminder.isSent
              ? colorScheme.onSurface.withValues(alpha: 0.3)
              : colorScheme.primary,
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                DateFormat('MMM dd, yyyy • HH:mm').format(reminder.remindAt),
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: colorScheme.onSurface,
                  fontWeight: FontWeight.w500,
                ),
              ),
              if (reminder.message != null) ...[
                const SizedBox(height: 4),
                Text(
                  reminder.message!,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: colorScheme.onSurface.withValues(alpha: 0.7),
                  ),
                ),
              ],
            ],
          ),
        ),
        if (reminder.isSent)
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: colorScheme.onSurface.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(6),
            ),
            child: Text(
              l10n.sent,
              style: theme.textTheme.bodySmall?.copyWith(
                color: colorScheme.onSurface.withValues(alpha: 0.7),
                fontSize: 10,
              ),
            ),
          ),
      ],
    );
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({
    required this.status,
    required this.theme,
    required this.colorScheme,
  });

  final String status;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    Color backgroundColor;
    Color textColor;
    String displayText;

    switch (status.toLowerCase()) {
      case 'completed':
        backgroundColor = Colors.green.withValues(alpha: 0.1);
        textColor = Colors.green;
        displayText = 'COMPLETED';
        break;
      case 'in_progress':
        backgroundColor = Colors.blue.withValues(alpha: 0.1);
        textColor = Colors.blue;
        displayText = 'IN PROGRESS';
        break;
      case 'pending':
        backgroundColor = Colors.orange.withValues(alpha: 0.1);
        textColor = Colors.orange;
        displayText = 'PENDING';
        break;
      case 'cancelled':
        backgroundColor = colorScheme.onSurface.withValues(alpha: 0.1);
        textColor = colorScheme.onSurface.withValues(alpha: 0.7);
        displayText = 'CANCELLED';
        break;
      default:
        backgroundColor = colorScheme.onSurface.withValues(alpha: 0.1);
        textColor = colorScheme.onSurface.withValues(alpha: 0.7);
        displayText = status.toUpperCase();
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        displayText,
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.w600,
          fontSize: 12,
        ),
      ),
    );
  }
}
