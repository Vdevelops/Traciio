import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../../../../core/l10n/app_localizations.dart';

class ScheduleStatusFilterSheet extends StatelessWidget {
  const ScheduleStatusFilterSheet({
    super.key,
    this.selectedStatus,
    required this.onSelect,
  });

  final String? selectedStatus;
  final ValueChanged<String?> onSelect;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final colorScheme = theme.colorScheme;
    final statuses = [null, 'pending', 'confirmed', 'completed', 'cancelled'];

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            l10n.filterByStatus,
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.onSurface,
            ),
          ),
          const SizedBox(height: 16),
          ...statuses.map((status) {
            final label = status == null
                ? l10n.all
                : status[0].toUpperCase() + status.substring(1);
            final isSelected = selectedStatus == status;
            return ListTile(
              title: Text(label),
              trailing: isSelected
                  ? Icon(Icons.check, color: colorScheme.primary)
                  : null,
              onTap: () => onSelect(status),
            );
          }),
        ],
      ),
    );
  }
}

class ScheduleDateFilterSheet extends StatefulWidget {
  const ScheduleDateFilterSheet({
    super.key,
    this.scheduledFrom,
    this.scheduledTo,
    required this.onSelect,
  });

  final DateTime? scheduledFrom;
  final DateTime? scheduledTo;
  final ValueChanged<Map<String, DateTime?>> onSelect;

  @override
  State<ScheduleDateFilterSheet> createState() =>
      _ScheduleDateFilterSheetState();
}

class _ScheduleDateFilterSheetState extends State<ScheduleDateFilterSheet> {
  DateTime? _selectedFrom;
  DateTime? _selectedTo;

  @override
  void initState() {
    super.initState();
    _selectedFrom = widget.scheduledFrom;
    _selectedTo = widget.scheduledTo;
  }

  Future<void> _selectDateFrom(BuildContext context) async {
    final date = await showDatePicker(
      context: context,
      initialDate: _selectedFrom ?? DateTime.now(),
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
    );
    if (date != null) {
      setState(() {
        _selectedFrom = date;
      });
      widget.onSelect({'from': _selectedFrom, 'to': _selectedTo});
    }
  }

  Future<void> _selectDateTo(BuildContext context) async {
    final date = await showDatePicker(
      context: context,
      initialDate: _selectedTo ?? (_selectedFrom ?? DateTime.now()),
      firstDate: _selectedFrom ?? DateTime(2000),
      lastDate: DateTime(2100),
    );
    if (date != null) {
      setState(() {
        _selectedTo = date;
      });
      widget.onSelect({'from': _selectedFrom, 'to': _selectedTo});
    }
  }

  void _clearDates() {
    setState(() {
      _selectedFrom = null;
      _selectedTo = null;
    });
    widget.onSelect({'from': null, 'to': null});
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final colorScheme = theme.colorScheme;
    final dateFormat = DateFormat('MMM dd, yyyy');

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                l10n.filterByScheduledDate,
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: colorScheme.onSurface,
                ),
              ),
              if (_selectedFrom != null || _selectedTo != null)
                TextButton(
                  onPressed: _clearDates,
                  child: Text(
                    l10n.clearFilters,
                    style: TextStyle(color: colorScheme.error),
                  ),
                ),
            ],
          ),
          const SizedBox(height: 16),
          InkWell(
            onTap: () => _selectDateFrom(context),
            child: InputDecorator(
              decoration: InputDecoration(
                labelText: l10n.scheduledFrom,
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
                _selectedFrom != null
                    ? dateFormat.format(_selectedFrom!)
                    : l10n.selectScheduledFrom,
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: _selectedFrom != null
                      ? colorScheme.onSurface
                      : colorScheme.onSurface.withValues(alpha: 0.5),
                ),
              ),
            ),
          ),
          const SizedBox(height: 16),
          InkWell(
            onTap: () => _selectDateTo(context),
            child: InputDecorator(
              decoration: InputDecoration(
                labelText: l10n.scheduledTo,
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
                _selectedTo != null
                    ? dateFormat.format(_selectedTo!)
                    : l10n.selectScheduledTo,
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: _selectedTo != null
                      ? colorScheme.onSurface
                      : colorScheme.onSurface.withValues(alpha: 0.5),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
