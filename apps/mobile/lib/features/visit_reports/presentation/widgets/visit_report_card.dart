import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../../../../core/l10n/app_localizations.dart';
import '../../data/models/visit_report.dart';

class VisitReportCard extends StatelessWidget {
  const VisitReportCard({
    super.key,
    required this.visitReport,
    required this.onTap,
  });

  final VisitReport visitReport;
  final VoidCallback onTap;

  static String _formatVisitDate(String visitDate, Locale locale) {
    try {
      // Try parsing ISO 8601 format
      final dateTime = DateTime.parse(visitDate).toLocal();

      // Use locale-aware date format
      final localeString = locale.languageCode == 'id' ? 'id_ID' : 'en_US';
      final dateFormat = DateFormat('EEEE, d MMMM yyyy', localeString);
      final timeFormat = DateFormat('HH:mm', localeString);

      // Check if time is included
      if (visitDate.contains('T') && visitDate.length > 10) {
        final timePart = visitDate.split('T')[1];
        if (timePart.isNotEmpty && !timePart.startsWith('00:00:00')) {
          return '${dateFormat.format(dateTime)} at ${timeFormat.format(dateTime)}';
        }
      }

      return dateFormat.format(dateTime);
    } catch (e) {
      // If parsing fails, return as-is
      return visitDate;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final l10n = AppLocalizations.of(context)!;
    final locale = Localizations.localeOf(context);

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: theme.brightness == Brightness.dark
            ? colorScheme.surfaceContainerHighest.withValues(alpha: 0.1)
            : colorScheme.surface,
        borderRadius: BorderRadius.circular(24),
        border: theme.brightness == Brightness.dark
            ? Border.all(
                color: colorScheme.outlineVariant.withValues(alpha: 0.2),
              )
            : null,
        boxShadow: theme.brightness == Brightness.dark
            ? []
            : [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.05),
                  blurRadius: 3,
                  offset: const Offset(0, 1),
                ),
              ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(24),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Header: Type Badge & Status Badge
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    _TypeBadge(
                      visitReport: visitReport,
                      theme: theme,
                      colorScheme: colorScheme,
                      l10n: l10n,
                    ),
                    _StatusBadge(status: visitReport.status),
                  ],
                ),
                const SizedBox(height: 12),
                // Purpose (Main content - replaces account position)
                if (visitReport.purpose != null &&
                    visitReport.purpose!.isNotEmpty)
                  Text(
                    visitReport.purpose!,
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                      color: colorScheme.onSurface,
                    ),
                  )
                else
                  Text(
                    'No purpose specified',
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                      color: colorScheme.onSurface.withValues(alpha: 0.5),
                      fontStyle: FontStyle.italic,
                    ),
                  ),
                const SizedBox(height: 12),
                // Account & Contact (Below purpose)
                if (visitReport.account != null) ...[
                  _DetailRow(
                    icon: Icons.business_outlined,
                    label: visitReport.account!.name,
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                  if (visitReport.contact != null) ...[
                    const SizedBox(height: 4),
                    _DetailRow(
                      icon: Icons.person_outline,
                      label: visitReport.contact!.name,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  ],
                ] else if (visitReport.deal != null) ...[
                  _DetailRow(
                    icon: Icons.handshake_outlined,
                    label: visitReport.deal!.title,
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                  if (visitReport.deal!.account != null) ...[
                    const SizedBox(height: 4),
                    _DetailRow(
                      icon: Icons.business_outlined,
                      label: visitReport.deal!.account!.name,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  ],
                  if (visitReport.contact != null) ...[
                    const SizedBox(height: 4),
                    _DetailRow(
                      icon: Icons.person_outline,
                      label: visitReport.contact!.name,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  ],
                ] else if (visitReport.lead != null) ...[
                  _DetailRow(
                    icon: Icons.person_outline,
                    label:
                        '${visitReport.lead!.firstName} ${visitReport.lead!.lastName ?? ''}'
                            .trim(),
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                  if (visitReport.lead!.companyName != null) ...[
                    const SizedBox(height: 4),
                    _DetailRow(
                      icon: Icons.business_outlined,
                      label: visitReport.lead!.companyName!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  ],
                ],
                // Visit Date & Time (At the bottom)
                const SizedBox(height: 12),
                _DetailRow(
                  icon: Icons.calendar_today_outlined,
                  label: _formatVisitDate(visitReport.visitDate, locale),
                  theme: theme,
                  colorScheme: colorScheme,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _TypeBadge extends StatelessWidget {
  const _TypeBadge({
    required this.visitReport,
    required this.theme,
    required this.colorScheme,
    required this.l10n,
  });

  final VisitReport visitReport;
  final ThemeData theme;
  final ColorScheme colorScheme;
  final AppLocalizations l10n;

  @override
  Widget build(BuildContext context) {
    String typeText;
    Color backgroundColor;
    Color textColor;

    // Use type from API response instead of inferring
    // Note: Type badge colors are chosen to NOT conflict with status badge colors
    // Status colors: grey (draft), yellow (submitted), orange (in_progress),
    //                primary/blue (completed), green (approved), red (rejected)
    switch (visitReport.type.toLowerCase()) {
      case 'lead':
        typeText = l10n.lead;
        backgroundColor = Colors.purple.withValues(alpha: 0.1);
        textColor = Colors.purple;
        break;
      case 'deal':
        typeText = l10n.deal;
        // Using teal to avoid conflict with orange (in_progress status)
        backgroundColor = Colors.teal.withValues(alpha: 0.1);
        textColor = Colors.teal;
        break;
      case 'account':
      default:
        // Default to Account type - using accounts (plural) from localization
        // Using pink to avoid conflict with primary/blue (completed status)
        typeText = l10n.accounts;
        backgroundColor = Colors.pink.withValues(alpha: 0.1);
        textColor = Colors.pink;
        break;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        typeText.toUpperCase(),
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.w600,
          fontSize: 10,
        ),
      ),
    );
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({required this.status});

  final String status;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    Color backgroundColor;
    Color textColor;
    String displayText;

    switch (status.toLowerCase()) {
      case 'draft':
        backgroundColor = colorScheme.onSurface.withValues(alpha: 0.1);
        textColor = colorScheme.onSurface.withValues(alpha: 0.7);
        displayText = 'DRAFT';
        break;
      case 'submitted':
      case 'pending':
        backgroundColor = Colors.yellow.withValues(alpha: 0.2);
        textColor = Colors.orange.shade800;
        displayText = 'SUBMITTED';
        break;
      case 'in_progress':
        backgroundColor = Colors.orange.withValues(alpha: 0.1);
        textColor = Colors.orange;
        displayText = 'IN PROGRESS';
        break;
      case 'completed':
        backgroundColor = colorScheme.primary.withValues(alpha: 0.1);
        textColor = colorScheme.primary;
        displayText = 'COMPLETED';
        break;
      case 'approved':
        backgroundColor = Colors.green.withValues(alpha: 0.1);
        textColor = Colors.green;
        displayText = 'APPROVED';
        break;
      case 'rejected':
        backgroundColor = colorScheme.error.withValues(alpha: 0.1);
        textColor = colorScheme.error;
        displayText = 'REJECTED';
        break;
      default:
        backgroundColor = colorScheme.onSurface.withValues(alpha: 0.1);
        textColor = colorScheme.onSurface.withValues(alpha: 0.7);
        displayText = status.toUpperCase();
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        displayText,
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.w600,
          fontSize: 10,
        ),
      ),
    );
  }
}

class _DetailRow extends StatelessWidget {
  const _DetailRow({
    required this.icon,
    required this.label,
    required this.theme,
    required this.colorScheme,
  });

  final IconData icon;
  final String label;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(
          icon,
          size: 16,
          color: colorScheme.onSurface.withValues(alpha: 0.7),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: Text(
            label,
            style: theme.textTheme.bodyMedium?.copyWith(
              color: colorScheme.onSurface.withValues(alpha: 0.7),
            ),
          ),
        ),
      ],
    );
  }
}
