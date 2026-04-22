import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../../data/models/dashboard.dart';
import '../../../../core/routing/app_router.dart';
import '../../../../core/l10n/app_localizations.dart';

class VisitCard extends StatelessWidget {
  final MobileVisit visit;

  const VisitCard({super.key, required this.visit});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final l10n = AppLocalizations.of(context)!;
    final locale = Localizations.localeOf(context);
    
    return Container(
      width: 280,
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: theme.brightness == Brightness.dark
            ? colorScheme.surfaceContainerHighest.withValues(alpha: 0.1)
            : colorScheme.surface,
        borderRadius: BorderRadius.circular(24),
        border: theme.brightness == Brightness.dark
            ? Border.all(color: colorScheme.outlineVariant.withValues(alpha: 0.2))
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
          onTap: () {
            // Navigate to visit report detail
            Navigator.pushNamed(
              context,
              '${AppRoutes.visitReports}/${visit.id}',
              arguments: visit.id,
            );
          },
          borderRadius: BorderRadius.circular(24),
            child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                // Header: Type Badge & Status Badge (konsisten dengan VisitReportCard)
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    _TypeBadge(visit: visit, l10n: l10n, colorScheme: colorScheme),
                    _StatusBadge(status: visit.status),
                  ],
                ),
                const SizedBox(height: 12),
                // Purpose (Main content)
                Text(
                  visit.purpose.isNotEmpty ? visit.purpose : 'No purpose specified',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: visit.purpose.isNotEmpty
                        ? colorScheme.onSurface
                        : colorScheme.onSurface.withValues(alpha: 0.5),
                    fontStyle: visit.purpose.isEmpty ? FontStyle.italic : FontStyle.normal,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 12),
                // Account & Contact (Below purpose)
                if (visit.accountName != null && visit.accountName!.isNotEmpty) ...[
                  _DetailRow(
                    icon: Icons.business_outlined,
                    label: visit.accountName!,
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                  if (visit.contactName != null && visit.contactName!.isNotEmpty) ...[
                    const SizedBox(height: 4),
                    _DetailRow(
                      icon: Icons.person_outline,
                      label: visit.contactName!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  ],
                ]
                else if (visit.dealTitle != null && visit.dealTitle!.isNotEmpty) ...[
                  _DetailRow(
                    icon: Icons.handshake_outlined,
                    label: visit.dealTitle!,
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                  if (visit.accountName != null && visit.accountName!.isNotEmpty) ...[
                    const SizedBox(height: 4),
                    _DetailRow(
                      icon: Icons.business_outlined,
                      label: visit.accountName!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  ],
                  if (visit.contactName != null && visit.contactName!.isNotEmpty) ...[
                    const SizedBox(height: 4),
                    _DetailRow(
                      icon: Icons.person_outline,
                      label: visit.contactName!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  ],
                ]
                else if (visit.leadName != null && visit.leadName!.isNotEmpty) ...[
                  _DetailRow(
                    icon: Icons.person_outline,
                    label: visit.leadName!,
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                ],
                // Visit Date & Time (At the bottom, konsisten dengan VisitReportCard)
                const SizedBox(height: 12),
                _DetailRow(
                  icon: Icons.calendar_today_outlined,
                  label: _formatVisitDate(visit.visitDate, visit.visitTime, locale),
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

  String _formatVisitDate(String visitDate, String? visitTime, Locale locale) {
    try {
      // Try parsing ISO 8601 format or simple date format
      DateTime dateTime;
      if (visitDate.contains('T')) {
        dateTime = DateTime.parse(visitDate).toLocal();
      } else {
        dateTime = DateTime.parse(visitDate);
      }
      
      // Use locale-aware date format
      final localeString = locale.languageCode == 'id' ? 'id_ID' : 'en_US';
      final dateFormat = DateFormat('EEEE, d MMMM yyyy', localeString);
      
      // If time is provided separately, include it
      if (visitTime != null && visitTime.isNotEmpty) {
        final timeFormat = DateFormat('HH:mm', localeString);
        try {
          final timeParts = visitTime.split(':');
          if (timeParts.length >= 2) {
            final hour = int.parse(timeParts[0]);
            final minute = int.parse(timeParts[1]);
            dateTime = DateTime(dateTime.year, dateTime.month, dateTime.day, hour, minute);
            return '${dateFormat.format(dateTime)} at ${timeFormat.format(dateTime)}';
          }
        } catch (e) {
          // If time parsing fails, just use date
        }
      }
      
      return dateFormat.format(dateTime);
    } catch (e) {
      // If parsing fails, return as-is
      return visitDate;
    }
  }
}

class _TypeBadge extends StatelessWidget {
  const _TypeBadge({
    required this.visit,
    required this.l10n,
    required this.colorScheme,
  });

  final MobileVisit visit;
  final AppLocalizations l10n;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    String typeText;
    Color backgroundColor;
    Color textColor;

    // Use type from API response
    switch (visit.type.toLowerCase()) {
      case 'lead':
        typeText = l10n.lead.toUpperCase();
        backgroundColor = Colors.purple.withValues(alpha: 0.1);
        textColor = Colors.purple;
        break;
      case 'deal':
        typeText = l10n.deal.toUpperCase();
        backgroundColor = Colors.orange.withValues(alpha: 0.1);
        textColor = Colors.orange;
        break;
      case 'account':
      default:
        typeText = l10n.accounts.toUpperCase();
        backgroundColor = colorScheme.primary.withValues(alpha: 0.1);
        textColor = colorScheme.primary;
        break;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        typeText,
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

    // Map dashboard status to visit report status format
    switch (status.toLowerCase()) {
      case 'planned':
      case 'draft':
        backgroundColor = colorScheme.onSurface.withValues(alpha: 0.1);
        textColor = colorScheme.onSurface.withValues(alpha: 0.7);
        displayText = 'DRAFT';
        break;
      case 'in_progress':
        backgroundColor = Colors.orange.withValues(alpha: 0.1);
        textColor = Colors.orange;
        displayText = 'IN PROGRESS';
        break;
      case 'completed':
      case 'approved':
        backgroundColor = Colors.green.withValues(alpha: 0.1);
        textColor = Colors.green;
        displayText = 'COMPLETED';
        break;
      case 'cancelled':
      case 'rejected':
        backgroundColor = colorScheme.error.withValues(alpha: 0.1);
        textColor = colorScheme.error;
        displayText = status.toUpperCase();
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
        Icon(icon, size: 16, color: colorScheme.onSurface.withValues(alpha: 0.7)),
        const SizedBox(width: 8),
        Expanded(
          child: Text(
            label,
            style: theme.textTheme.bodyMedium?.copyWith(
              color: colorScheme.onSurface.withValues(alpha: 0.7),
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }
}
