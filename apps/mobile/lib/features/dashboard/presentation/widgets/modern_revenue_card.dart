import 'package:flutter/material.dart';
import '../../../../core/l10n/app_localizations.dart';

class ModernRevenueCard extends StatelessWidget {
  final String revenue;
  final double revenueChange;
  final int totalDeals;
  final int wonDeals;
  final VoidCallback? onTap;

  const ModernRevenueCard({
    super.key,
    required this.revenue,
    required this.revenueChange,
    required this.totalDeals,
    required this.wonDeals,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final isPositive = revenueChange >= 0;
    final isDark = theme.brightness == Brightness.dark;

    // Colors
    final bgColor = isDark ? const Color(0xFF2D241E) : const Color(0xFFFFF7ED);
    final borderColor = isDark
        ? const Color(0xFF433025)
        : const Color(0xFFFFEDD5);
    final titleColor = isDark
        ? const Color(0xFFFFBCA0)
        : const Color(0xFF9A3412);
    final valueColor = isDark
        ? const Color(0xFFFFF1EB)
        : const Color(0xFF431407);
    final sidePanelColor = isDark
        ? Colors.black.withValues(alpha: 0.3)
        : Colors.white;

    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(24),
        child: Container(
          padding: const EdgeInsets.all(20),
          decoration: BoxDecoration(
            color: bgColor,
            borderRadius: BorderRadius.circular(24),
            border: Border.all(color: borderColor),
            boxShadow: [
              BoxShadow(
                color: Colors.orange.withValues(alpha: 0.04),
                blurRadius: 3,
                offset: const Offset(0, 1),
              ),
            ],
          ),
          child: Row(
            children: [
              Expanded(
                flex: 3,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      l10n.revenue,
                      style: theme.textTheme.titleMedium?.copyWith(
                        color: titleColor,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      revenue,
                      style: theme.textTheme.headlineSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: valueColor,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 4,
                      ),
                      decoration: BoxDecoration(
                        color: sidePanelColor,
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Icon(
                            isPositive
                                ? Icons.trending_up
                                : Icons.trending_down,
                            size: 14,
                            color: isPositive ? Colors.green : Colors.red,
                          ),
                          const SizedBox(width: 4),
                          Text(
                            '${revenueChange.abs()}% vs last period',
                            style: TextStyle(
                              fontSize: 11,
                              fontWeight: FontWeight.w600,
                              color: isPositive ? Colors.green : Colors.red,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
              Expanded(
                flex: 2,
                child: Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: sidePanelColor,
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      _buildMiniStat(
                        context,
                        'Won Deals',
                        wonDeals.toString(),
                        Colors.orange,
                        isDark,
                      ),
                      const SizedBox(height: 12),
                      _buildMiniStat(
                        context,
                        'Total Deals',
                        totalDeals.toString(),
                        isDark ? Colors.grey.shade400 : Colors.grey,
                        isDark,
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildMiniStat(
    BuildContext context,
    String label,
    String value,
    Color color,
    bool isDark,
  ) {
    return Row(
      children: [
        Container(
          width: 3,
          height: 24,
          decoration: BoxDecoration(
            color: color,
            borderRadius: BorderRadius.circular(2),
          ),
        ),
        const SizedBox(width: 8),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              value,
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
                color: isDark ? Colors.white : Colors.black87,
              ),
            ),
            Text(
              label,
              style: TextStyle(
                fontSize: 10,
                color: isDark ? Colors.white60 : Colors.grey[600],
              ),
            ),
          ],
        ),
      ],
    );
  }
}
