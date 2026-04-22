import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../../../../core/l10n/app_localizations.dart';

class StatusFilterSheet extends StatelessWidget {
  const StatusFilterSheet({
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
    final statuses = [null, 'pending', 'in_progress', 'completed', 'cancelled'];

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
                : status.toUpperCase().replaceAll('_', ' ');
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

class PriorityFilterSheet extends StatelessWidget {
  const PriorityFilterSheet({
    super.key,
    this.selectedPriority,
    required this.onSelect,
  });

  final String? selectedPriority;
  final ValueChanged<String?> onSelect;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final colorScheme = theme.colorScheme;
    final priorities = [null, 'low', 'medium', 'high', 'urgent'];

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            l10n.filterByPriority,
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.onSurface,
            ),
          ),
          const SizedBox(height: 16),
          ...priorities.map((priority) {
            final label = priority == null ? l10n.all : priority.toUpperCase();
            final isSelected = selectedPriority == priority;
            return ListTile(
              title: Text(label),
              trailing: isSelected
                  ? Icon(Icons.check, color: colorScheme.primary)
                  : null,
              onTap: () => onSelect(priority),
            );
          }),
        ],
      ),
    );
  }
}

class TypeFilterSheet extends StatelessWidget {
  const TypeFilterSheet({
    super.key,
    this.selectedType,
    required this.onSelect,
  });

  final String? selectedType;
  final ValueChanged<String?> onSelect;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final colorScheme = theme.colorScheme;
    final types = [null, 'general', 'call', 'email', 'meeting', 'follow_up'];

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            l10n.filterByType,
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.onSurface,
            ),
          ),
          const SizedBox(height: 16),
          ...types.map((type) {
            final label = type == null
                ? l10n.all
                : type.toUpperCase().replaceAll('_', ' ');
            final isSelected = selectedType == type;
            return ListTile(
              title: Text(label),
              trailing: isSelected
                  ? Icon(Icons.check, color: colorScheme.primary)
                  : null,
              onTap: () => onSelect(type),
            );
          }),
        ],
      ),
    );
  }
}

class DueDateFilterSheet extends StatefulWidget {
  const DueDateFilterSheet({
    super.key,
    this.dueDateFrom,
    this.dueDateTo,
    required this.onSelect,
  });

  final DateTime? dueDateFrom;
  final DateTime? dueDateTo;
  final ValueChanged<Map<String, DateTime?>> onSelect;

  @override
  State<DueDateFilterSheet> createState() => _DueDateFilterSheetState();
}

class _DueDateFilterSheetState extends State<DueDateFilterSheet> {
  DateTime? _selectedFrom;
  DateTime? _selectedTo;

  @override
  void initState() {
    super.initState();
    _selectedFrom = widget.dueDateFrom;
    _selectedTo = widget.dueDateTo;
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
      widget.onSelect({
        'from': _selectedFrom,
        'to': _selectedTo,
      });
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
      widget.onSelect({
        'from': _selectedFrom,
        'to': _selectedTo,
      });
    }
  }

  void _clearDates() {
    setState(() {
      _selectedFrom = null;
      _selectedTo = null;
    });
    widget.onSelect({
      'from': null,
      'to': null,
    });
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
                l10n.filterByDueDate,
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
          // Due Date From
          InkWell(
            onTap: () => _selectDateFrom(context),
            child: InputDecorator(
              decoration: InputDecoration(
                labelText: l10n.dueDateFrom,
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
                    : l10n.selectDueDateFrom,
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: _selectedFrom != null
                      ? colorScheme.onSurface
                      : colorScheme.onSurface.withValues(alpha: 0.5),
                ),
              ),
            ),
          ),
          const SizedBox(height: 16),
          // Due Date To
          InkWell(
            onTap: () => _selectDateTo(context),
            child: InputDecorator(
              decoration: InputDecoration(
                labelText: l10n.dueDateTo,
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
                    : l10n.selectDueDateTo,
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
