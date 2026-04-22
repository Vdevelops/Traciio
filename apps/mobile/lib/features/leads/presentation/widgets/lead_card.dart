import 'package:flutter/material.dart';

import '../../data/models/lead_model.dart';

class LeadCard extends StatelessWidget {
  final Lead lead;
  final VoidCallback? onTap;

  const LeadCard({super.key, required this.lead, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;
    final colorScheme = theme.colorScheme;

    // Provide default status if null
    final status =
        lead.leadStatus ??
        LeadStatus(id: 'new', name: 'New', code: 'new', color: '#808080');
    final source = lead.source;

    return Container(
      margin: const EdgeInsets.symmetric(vertical: 6),
      decoration: BoxDecoration(
        color: isDark
            ? colorScheme.surfaceContainerHighest.withValues(alpha: 0.1)
            : colorScheme.surface,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(
          color: isDark
              ? Colors.white.withValues(alpha: 0.05)
              : Colors.transparent,
        ),
        boxShadow: isDark
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
                // Source Badge + Status Badge (space between, like pipeline)
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    if (source.isNotEmpty)
                      _SourceBadge(source: source)
                    else
                      const SizedBox.shrink(),
                    _StatusBadge(status: status),
                  ],
                ),
                const SizedBox(height: 12),
                // Name
                Text(
                  lead.name,
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: colorScheme.onSurface,
                  ),
                ),
                // Company Name
                if (lead.company.isNotEmpty) ...[
                  const SizedBox(height: 4),
                  Text(
                    lead.company,
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: colorScheme.onSurface.withValues(alpha: 0.7),
                    ),
                  ),
                ],
                if ((lead.phone != null && lead.phone!.isNotEmpty) ||
                    (lead.email != null && lead.email!.isNotEmpty) ||
                    (lead.province != null && lead.province!.isNotEmpty)) ...[
                  const SizedBox(height: 12),
                  if (lead.phone != null && lead.phone!.isNotEmpty)
                    Padding(
                      padding: const EdgeInsets.only(bottom: 8),
                      child: _buildInfoRow(
                        context,
                        Icons.phone_outlined,
                        lead.phone!,
                      ),
                    ),
                  if (lead.email != null && lead.email!.isNotEmpty)
                    Padding(
                      padding: const EdgeInsets.only(bottom: 8),
                      child: _buildInfoRow(
                        context,
                        Icons.email_outlined,
                        lead.email!,
                      ),
                    ),
                  if (lead.province != null && lead.province!.isNotEmpty)
                    _buildInfoRow(
                      context,
                      Icons.location_on_outlined,
                      lead.province!,
                    ),
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildInfoRow(BuildContext context, IconData icon, String text) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

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
            text,
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

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({required this.status});

  final LeadStatus status;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final color = _getStatusColor(status.code, colorScheme);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.2),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        status.name.toUpperCase(),
        style: theme.textTheme.bodySmall?.copyWith(
          color: color,
          fontWeight: FontWeight.bold,
          fontSize: 10,
        ),
      ),
    );
  }

  Color _getStatusColor(String code, ColorScheme colorScheme) {
    switch (code.toLowerCase()) {
      case 'new':
        return Colors.blue;
      case 'contacted':
        return Colors.orange;
      case 'qualified':
        return Colors.green;
      case 'unqualified':
        return Colors.red;
      case 'converted':
        return Colors.purple;
      case 'lost':
        return Colors.grey;
      default:
        return colorScheme.primary;
    }
  }
}

class _SourceBadge extends StatelessWidget {
  const _SourceBadge({required this.source});

  final String source;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    // Get color based on source type
    final color = _getSourceColor(source);

    final badgeColor = isDark
        ? color.withValues(alpha: 0.7)
        : color.withValues(alpha: 0.1);
    final textColor = isDark ? Colors.white : color;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: badgeColor,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        source,
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.bold,
          fontSize: 10,
        ),
      ),
    );
  }

  // Color coding for different sources
  // Avoiding status colors: blue, orange, green, red, purple, grey
  Color _getSourceColor(String source) {
    final sourceLower = source.toLowerCase();

    switch (sourceLower) {
      case 'website':
      case 'web':
        return Colors.cyan; // Not conflicting with status colors
      case 'referral':
        return Colors.pink;
      case 'social media':
      case 'social':
        return Colors.indigo;
      case 'event':
      case 'events':
        return Colors.amber;
      case 'email':
      case 'email campaign':
        return Colors.teal;
      case 'cold call':
      case 'call':
        return Colors.brown;
      case 'advertisement':
      case 'ads':
        return Colors.deepOrange;
      case 'partner':
        return Colors.lime;
      default:
        return Colors.cyan; // Default color
    }
  }
}
