import 'package:flutter/material.dart';

import '../../data/models/account.dart';

class AccountCard extends StatelessWidget {
  const AccountCard({super.key, required this.account, required this.onTap});

  final Account account;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final isDark = theme.brightness == Brightness.dark;

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 0, vertical: 8),
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
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Row: Category badge + Status badge (space between)
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    if (account.category != null)
                      _TypeBadge(
                        typeName: account.category!.name,
                        theme: theme,
                        colorScheme: colorScheme,
                      )
                    else
                      const SizedBox.shrink(),
                    _StatusBadge(status: account.status, theme: theme),
                  ],
                ),
                const SizedBox(height: 12),
                // Title
                Text(
                  account.name,
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: colorScheme.onSurface,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                // Full address (if exists)
                if (_hasAddress(account)) ...[
                  const SizedBox(height: 8),
                  _buildAddressRow(context, account),
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }

  // Check if account has any address data
  bool _hasAddress(Account account) {
    return account.address != null ||
        account.city != null ||
        account.province != null;
  }

  // Build full address string
  String _buildFullAddress(Account account) {
    final parts = <String>[];
    if (account.address != null && account.address!.isNotEmpty) {
      parts.add(account.address!);
    }
    if (account.city != null && account.city!.isNotEmpty) {
      parts.add(account.city!);
    }
    if (account.province != null && account.province!.isNotEmpty) {
      parts.add(account.province!);
    }
    return parts.join(', ');
  }

  // Build address row widget
  Widget _buildAddressRow(BuildContext context, Account account) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final fullAddress = _buildFullAddress(account);

    if (fullAddress.isEmpty) return const SizedBox.shrink();

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(
          Icons.location_on_outlined,
          size: 16,
          color: colorScheme.onSurface.withValues(alpha: 0.7),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: Text(
            fullAddress,
            style: theme.textTheme.bodyMedium?.copyWith(
              color: colorScheme.onSurface.withValues(alpha: 0.7),
            ),
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({required this.status, required this.theme});

  final String status;
  final ThemeData theme;

  Color _getStatusColor(ColorScheme colorScheme) {
    final statusLower = status.toLowerCase();
    if (statusLower == 'active') return Colors.green;
    return colorScheme.onSurface.withValues(alpha: 0.7);
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = theme.colorScheme;
    final color = _getStatusColor(colorScheme);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        status.toUpperCase(),
        style: theme.textTheme.bodySmall?.copyWith(
          color: color,
          fontWeight: FontWeight.w600,
          fontSize: 10,
        ),
      ),
    );
  }
}

class _TypeBadge extends StatelessWidget {
  const _TypeBadge({
    required this.typeName,
    required this.theme,
    required this.colorScheme,
  });

  final String typeName;
  final ThemeData theme;
  final ColorScheme colorScheme;

  Color _getTypeColor() {
    final typeLower = typeName.toLowerCase();
    if (typeLower == 'clinic') return Colors.blue;
    if (typeLower == 'pharmacy') return Colors.orange;
    if (typeLower == 'hospital') return Colors.purple;
    return colorScheme.primary;
  }

  @override
  Widget build(BuildContext context) {
    final color = _getTypeColor();
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        typeName,
        style: theme.textTheme.bodySmall?.copyWith(
          color: color,
          fontWeight: FontWeight.w500,
          fontSize: 10,
        ),
      ),
    );
  }
}
