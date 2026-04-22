import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../../data/models/deal_model.dart';

class DealCard extends StatelessWidget {
  const DealCard({super.key, required this.deal, required this.onTap});

  final Deal deal;
  final VoidCallback onTap;

  String _formatCurrency(double amount) {
    final format = NumberFormat.currency(
      locale: 'id_ID',
      symbol: 'Rp ',
      decimalDigits: 0,
    );
    return format.format(amount);
  }

  String _formatDate(String dateStr, Locale locale) {
    try {
      final dateTime = DateTime.parse(dateStr).toLocal();
      final localeString = locale.languageCode == 'id' ? 'id_ID' : 'en_US';
      return DateFormat(' d MMM yyyy', localeString).format(dateTime);
    } catch (e) {
      return dateStr;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final locale = Localizations.localeOf(context);

    final isDark = theme.brightness == Brightness.dark;

    return Container(
      margin: const EdgeInsets.symmetric(vertical: 8),
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
            padding: const EdgeInsets.all(12),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Status and Probability badges at top (space between)
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    _ProbabilityBadge(probability: deal.probability),
                    _StatusBadge(status: deal.status),
                  ],
                ),
                const SizedBox(height: 4),
                // Title
                Text(
                  deal.title,
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: colorScheme.onSurface,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                if (deal.account != null) ...[
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Icon(
                        Icons.business_outlined,
                        size: 14,
                        color: colorScheme.onSurface.withValues(alpha: 0.6),
                      ),
                      const SizedBox(width: 6),
                      Expanded(
                        child: Text(
                          deal.account!.name,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: colorScheme.onSurface.withValues(alpha: 0.6),
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                    ],
                  ),
                ],
                // Contact info below account
                if (deal.contact != null) ...[
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Icon(
                        Icons.person_outline,
                        size: 14,
                        color: colorScheme.onSurface.withValues(alpha: 0.6),
                      ),
                      const SizedBox(width: 6),
                      Expanded(
                        child: Text(
                          deal.contact!.name,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: colorScheme.onSurface.withValues(alpha: 0.6),
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                    ],
                  ),
                ],
                const SizedBox(height: 8),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Value',
                          style: theme.textTheme.labelSmall?.copyWith(
                            color: colorScheme.onSurface.withValues(alpha: 0.5),
                          ),
                        ),
                        Text(
                          _formatCurrency(deal.value.toDouble()),
                          style: theme.textTheme.titleSmall?.copyWith(
                            fontWeight: FontWeight.bold,
                            color: colorScheme.primary,
                          ),
                        ),
                      ],
                    ),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.end,
                      children: [
                        Text(
                          'Expected Close',
                          style: theme.textTheme.labelSmall?.copyWith(
                            color: colorScheme.onSurface.withValues(alpha: 0.5),
                          ),
                        ),
                        Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(
                              Icons.calendar_today_outlined,
                              size: 12,
                              color: colorScheme.onSurface.withValues(
                                alpha: 0.6,
                              ),
                            ),
                            const SizedBox(width: 4),
                            Text(
                              _formatDate(
                                deal.expectedCloseDate?.toIso8601String() ?? '',
                                locale,
                              ),
                              style: theme.textTheme.bodySmall?.copyWith(
                                color: colorScheme.onSurface.withValues(
                                  alpha: 0.6,
                                ),
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _ProbabilityBadge extends StatelessWidget {
  const _ProbabilityBadge({required this.probability});

  final int probability;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    // Color based on probability percentage
    Color backgroundColor;
    Color textColor;

    if (probability >= 70) {
      backgroundColor = Colors.green.withValues(alpha: 0.1);
      textColor = Colors.green.shade700;
    } else if (probability >= 40) {
      backgroundColor = Colors.orange.withValues(alpha: 0.1);
      textColor = Colors.orange.shade700;
    } else {
      backgroundColor = Colors.grey.withValues(alpha: 0.1);
      textColor = Colors.grey.shade700;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.trending_up, size: 12, color: textColor),
          const SizedBox(width: 4),
          Text(
            '$probability%',
            style: theme.textTheme.bodySmall?.copyWith(
              color: textColor,
              fontWeight: FontWeight.w600,
              fontSize: 10,
            ),
          ),
        ],
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
    final isDark = theme.brightness == Brightness.dark;

    Color backgroundColor;
    Color textColor;
    String displayText;

    switch (status.toLowerCase()) {
      case 'open':
        backgroundColor = isDark
            ? Colors.blue.shade800.withValues(alpha: 0.4)
            : Colors.blue.withValues(alpha: 0.1);
        textColor = isDark ? Colors.blue.shade200 : Colors.blue.shade700;
        displayText = 'OPEN';
        break;
      case 'won':
        backgroundColor = isDark
            ? Colors.green.shade800.withValues(alpha: 0.4)
            : Colors.green.withValues(alpha: 0.1);
        textColor = isDark ? Colors.green.shade200 : Colors.green.shade700;
        displayText = 'WON';
        break;
      case 'lost':
        backgroundColor = isDark
            ? Colors.red.shade800.withValues(alpha: 0.4)
            : theme.colorScheme.error.withValues(alpha: 0.1);
        textColor = isDark ? Colors.red.shade200 : theme.colorScheme.error;
        displayText = 'LOST';
        break;
      default:
        final colorScheme = theme.colorScheme;
        backgroundColor = colorScheme.onSurface.withValues(alpha: 0.1);
        textColor = colorScheme.onSurface.withValues(alpha: 0.7);
        displayText = status.toUpperCase();
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
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
