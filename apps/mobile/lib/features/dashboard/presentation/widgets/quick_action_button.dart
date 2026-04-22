import 'package:flutter/material.dart';

/// Simple, large button untuk quick actions di dashboard
/// Designed untuk sales rep workflow - easy to tap saat di lapangan
class QuickActionButton extends StatelessWidget {
  const QuickActionButton({
    super.key,
    required this.icon,
    required this.label,
    required this.onTap,
    this.color,
    this.iconColor,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final Color? color;
  final Color? iconColor;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    
    // Better contrast colors - use surface dengan border untuk non-primary buttons
    final isPrimary = color == colorScheme.primary;
    final effectiveColor = isPrimary
        ? colorScheme.primary
        : colorScheme.surface;
    final effectiveIconColor = isPrimary
        ? (iconColor ?? colorScheme.onPrimary)
        : (iconColor ?? colorScheme.primary);
    
    // Enhanced border untuk better contrast
    final borderColor = isPrimary
        ? colorScheme.primary.withValues(alpha: 0.3)
        : colorScheme.outline.withValues(alpha: 0.3);

    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Container(
          constraints: const BoxConstraints(
            minHeight: 80, // Prevent overflow
          ),
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 12),
          decoration: BoxDecoration(
            color: effectiveColor,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: borderColor,
              width: 1.5,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.05),
                blurRadius: 3,
                offset: const Offset(0, 1),
              ),
            ],
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Icon(
                icon,
                size: 22,
                color: effectiveIconColor,
              ),
              const SizedBox(height: 6),
              Flexible(
                child: Text(
                  label,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: effectiveIconColor,
                    fontWeight: FontWeight.w600,
                    fontSize: 11,
                  ),
                  textAlign: TextAlign.center,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

