import 'package:flutter/material.dart';

class ModernTabSwitch extends StatelessWidget {
  final List<String> tabs;
  final int selectedIndex;
  final ValueChanged<int> onChanged;

  const ModernTabSwitch({
    super.key,
    required this.tabs,
    required this.selectedIndex,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final isDark = theme.brightness == Brightness.dark;

    if (tabs.isEmpty) return const SizedBox.shrink();

    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
        color: isDark
            ? colorScheme.surfaceContainerHighest.withValues(alpha: 0.3)
            : colorScheme.surfaceContainerHighest.withValues(alpha: 0.6),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Stack(
        children: [
          // 1. Sliding Pill Background
          AnimatedAlign(
            alignment: Alignment(
              tabs.length > 1 ? (selectedIndex / (tabs.length - 1) * 2) - 1 : 0,
              0,
            ),
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOutCubic,
            child: FractionallySizedBox(
              widthFactor: 1 / tabs.length,
              child: Container(
                height: 40,
                decoration: BoxDecoration(
                  // Use white in light mode for better contrast
                  color: isDark ? colorScheme.surface : Colors.white,
                  borderRadius: BorderRadius.circular(12),
                  boxShadow: [
                    BoxShadow(
                      color: isDark
                          ? Colors.black.withValues(alpha: 0.06)
                          : Colors.black.withValues(alpha: 0.08),
                      blurRadius: 3,
                      offset: const Offset(0, 1),
                    ),
                  ],
                ),
              ),
            ),
          ),
          // 2. Text Labels
          Row(
            children: List.generate(tabs.length, (index) {
              final isSelected = selectedIndex == index;
              return Expanded(
                child: GestureDetector(
                  onTap: () => onChanged(index),
                  behavior: HitTestBehavior.opaque,
                  child: Container(
                    padding: const EdgeInsets.symmetric(vertical: 12),
                    alignment: Alignment.center,
                    child: AnimatedDefaultTextStyle(
                      duration: const Duration(milliseconds: 200),
                      style: theme.textTheme.bodyMedium!.copyWith(
                        fontWeight: isSelected
                            ? FontWeight.bold
                            : FontWeight.w500,
                        color: isSelected
                            ? colorScheme.primary
                            : colorScheme.onSurfaceVariant.withValues(alpha: 0.7),
                      ),
                      child: Text(tabs[index]),
                    ),
                  ),
                ),
              );
            }),
          ),
        ],
      ),
    );
  }
}
