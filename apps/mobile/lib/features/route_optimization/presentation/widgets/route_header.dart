import 'package:flutter/material.dart';

/// Simple header widget untuk route optimization screens
/// Konsisten dengan design halaman lain tanpa hamburger menu
class RouteHeader extends StatelessWidget {
  const RouteHeader({
    super.key,
    required this.title,
    this.actions,
    this.onBack,
  });

  final String title;
  final List<Widget>? actions;
  final VoidCallback? onBack;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final screenWidth = MediaQuery.of(context).size.width;

    return Container(
      padding: EdgeInsets.symmetric(
        horizontal: screenWidth < 360 ? 12 : 16,
        vertical: 10,
      ),
      decoration: BoxDecoration(
        color: theme.cardColor,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 3,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: Row(
        children: [
          // Back Button
          IconButton(
            onPressed: onBack ?? () => Navigator.of(context).pop(),
            icon: const Icon(Icons.arrow_back),
            color: colorScheme.onSurface.withValues(alpha: 0.7),
            iconSize: 22,
            padding: EdgeInsets.zero,
            constraints: const BoxConstraints(
              minWidth: 40,
              minHeight: 40,
            ),
          ),
          const SizedBox(width: 8),
          // Title
          Expanded(
            child: Text(
              title,
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
                fontSize: screenWidth < 360 ? 18 : 20,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ),
          // Actions
          if (actions != null) ...actions!,
        ],
      ),
    );
  }
}

