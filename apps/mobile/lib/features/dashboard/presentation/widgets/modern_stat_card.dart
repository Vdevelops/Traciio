import 'package:flutter/material.dart';

enum ModernStatCardType { primary, white, transparent }

class ModernStatCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;
  final Color? iconColor;
  final ModernStatCardType type;
  final VoidCallback? onTap;
  final double? changePercent;

  const ModernStatCard({
    super.key,
    required this.title,
    required this.value,
    required this.icon,
    this.iconColor,
    this.type = ModernStatCardType.white,
    this.onTap,
    this.changePercent,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    final isPrimary = type == ModernStatCardType.primary;
    final backgroundColor = isPrimary
        ? colorScheme.primary
        : (type == ModernStatCardType.white
              ? colorScheme.surface
              : Colors.transparent);
    final foregroundColor = isPrimary
        ? colorScheme.onPrimary
        : colorScheme.onSurface;
    final subTextColor = isPrimary
        ? colorScheme.onPrimary.withValues(alpha: 0.8)
        : colorScheme.onSurface.withValues(alpha: 0.6);
    final effectiveIconColor =
        iconColor ?? (isPrimary ? colorScheme.onPrimary : colorScheme.primary);

    // Shadow only for white card — matches web shadow-sm
    final List<BoxShadow>? shadows = type == ModernStatCardType.white
        ? [
            BoxShadow(
              color: Colors.black.withValues(alpha: 0.05),
              blurRadius: 3,
              offset: const Offset(0, 1),
            ),
          ]
        : null;

    // Enhanced border untuk better contrast
    final border = type == ModernStatCardType.white
        ? Border.all(
            color: colorScheme.outline.withValues(alpha: 0.2),
            width: 1.5,
          )
        : null;

    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(20),
        child: Container(
          constraints: const BoxConstraints(
            minHeight: 120,
          ),
          padding: const EdgeInsets.all(18),
          decoration: BoxDecoration(
            color: backgroundColor,
            borderRadius: BorderRadius.circular(20),
            boxShadow: shadows,
            border: border,
            gradient: isPrimary
                ? LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [
                      colorScheme.primary,
                      colorScheme.primary.withValues(alpha: 0.8),
                    ],
                  )
                : null,
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            mainAxisSize: MainAxisSize.min,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    padding: const EdgeInsets.all(10),
                    decoration: BoxDecoration(
                      color: isPrimary
                          ? Colors.white.withValues(alpha: 0.25)
                          : effectiveIconColor.withValues(alpha: 0.15),
                      shape: BoxShape.circle,
                      boxShadow: isPrimary
                          ? [
                              BoxShadow(
                                color: Colors.white.withValues(alpha: 0.2),
                                blurRadius: 3,
                                offset: const Offset(0, 1),
                              ),
                            ]
                          : null,
                    ),
                    child: Icon(
                      icon,
                      size: 20,
                      color: isPrimary ? Colors.white : effectiveIconColor,
                    ),
                  ),
                  if (changePercent != null) _buildTrendBadge(isPrimary),
                ],
              ),
              const SizedBox(height: 12),
              Text(
                value,
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: foregroundColor,
                  fontSize: 24,
                  letterSpacing: -0.5,
                  height: 1.1,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
              const SizedBox(height: 4),
              Text(
                title,
                style: theme.textTheme.bodySmall?.copyWith(
                  color: subTextColor,
                  fontWeight: FontWeight.w500,
                  fontSize: 12,
                  height: 1.2,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildTrendBadge(bool isPrimary) {
    if (changePercent == null) return const SizedBox.shrink();

    final isPositive = changePercent! >= 0;
    final color = isPrimary
        ? Colors.white
        : (isPositive ? Colors.green : Colors.red);
    final bgColor = isPrimary
        ? Colors.white.withValues(alpha: 0.2)
        : (isPositive
              ? Colors.green.withValues(alpha: 0.1)
              : Colors.red.withValues(alpha: 0.1));

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            isPositive ? Icons.arrow_upward : Icons.arrow_downward,
            size: 10,
            color: color,
          ),
          const SizedBox(width: 2),
          Text(
            '${changePercent!.abs()}%',
            style: TextStyle(
              fontSize: 10,
              fontWeight: FontWeight.bold,
              color: color,
            ),
          ),
        ],
      ),
    );
  }
}
