import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../application/account_provider.dart';
import '../data/models/account.dart';
import '../presentation/account_form_screen.dart';
import '../../../core/permissions/permission_provider.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/routing/app_router.dart';
import '../../../core/widgets/error_widget.dart';
import '../../../core/widgets/loading_widget.dart';

class AccountDetailScreen extends ConsumerStatefulWidget {
  const AccountDetailScreen({
    super.key,
    required this.accountId,
  });

  final String accountId;

  @override
  ConsumerState<AccountDetailScreen> createState() => _AccountDetailScreenState();
}

class _AccountDetailScreenState extends ConsumerState<AccountDetailScreen> {
  bool _isDeleting = false;

  Future<void> _handleEdit(Account account) async {
    final result = await Navigator.push<Account>(
      context,
      MaterialPageRoute(
        builder: (context) => AccountFormScreen(account: account),
      ),
    );

    if (result != null && mounted) {
      // Refresh detail
      ref.invalidate(accountDetailProvider(widget.accountId));
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(AppLocalizations.of(context)!.accountUpdated),
          backgroundColor: Colors.green,
        ),
      );
    }
  }

  Future<void> _handleDelete(Account account) async {
    final l10n = AppLocalizations.of(context)!;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.deleteAccount),
        content: Text(l10n.deleteConfirmation),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: Text(l10n.cancel),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            style: FilledButton.styleFrom(
              backgroundColor: Theme.of(context).colorScheme.error,
            ),
            child: Text(l10n.delete),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      setState(() => _isDeleting = true);
      final success = await ref.read(accountListProvider.notifier).deleteAccount(account.id);
      setState(() => _isDeleting = false);

      if (mounted) {
        if (success) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(l10n.accountDeleted),
              backgroundColor: Colors.green,
            ),
          );
          Navigator.pop(context);
        } else {
          final error = ref.read(accountListProvider).errorMessage;
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(error ?? 'Failed to delete account'),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final accountAsync = ref.watch(accountDetailProvider(widget.accountId));
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.accountDetails),
        elevation: 0,
      ),
      body: accountAsync.when(
        data: (account) => SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Header Card
              Container(
                decoration: BoxDecoration(
                  color: colorScheme.surface,
                  borderRadius: BorderRadius.circular(20),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withValues(alpha: 0.05),
                      blurRadius: 3,
                      offset: const Offset(0, 1),
                    ),
                  ],
                ),
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            account.name,
                            style: theme.textTheme.headlineSmall?.copyWith(
                              fontWeight: FontWeight.bold,
                              color: colorScheme.onSurface,
                            ),
                          ),
                        ),
                        _StatusBadge(
                          status: account.status,
                          theme: theme,
                          colorScheme: colorScheme,
                        ),
                      ],
                    ),
                    if (account.category != null) ...[
                      const SizedBox(height: 12),
                      _TypeBadge(
                        typeName: account.category!.name,
                        theme: theme,
                        colorScheme: colorScheme,
                      ),
                    ],
                  ],
                ),
              ),
              const SizedBox(height: 16),
              // Contact Information
              _SectionTitle(
                title: 'Contact Information',
                theme: theme,
                colorScheme: colorScheme,
              ),
              const SizedBox(height: 8),
              _InfoCard(
                theme: theme,
                colorScheme: colorScheme,
                children: [
                  if (account.phone != null)
                    _InfoRow(
                      icon: Icons.phone_outlined,
                      label: l10n.phone,
                      value: account.phone!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  if (account.email != null)
                    _InfoRow(
                      icon: Icons.email_outlined,
                      label: l10n.email,
                      value: account.email!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                ],
              ),
              const SizedBox(height: 16),
              // Address Information
              _SectionTitle(
                title: 'Address Information',
                theme: theme,
                colorScheme: colorScheme,
              ),
              const SizedBox(height: 8),
              _InfoCard(
                theme: theme,
                colorScheme: colorScheme,
                children: [
                  if (account.address != null)
                    _InfoRow(
                      icon: Icons.location_on_outlined,
                      label: l10n.address,
                      value: account.address!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  if (account.city != null)
                    _InfoRow(
                      icon: Icons.location_city_outlined,
                      label: l10n.city,
                      value: account.city!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  if (account.province != null)
                    _InfoRow(
                      icon: Icons.map_outlined,
                      label: l10n.province,
                      value: account.province!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                ],
              ),
              const SizedBox(height: 16),
              // Action Buttons
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton.icon(
                      onPressed: () async {
                        final result = await Navigator.pushNamed(
                          context,
                          AppRoutes.contacts,
                          arguments: {'accountId': account.id},
                        );
                        // Refresh account detail if contact was deleted
                        if (result == true && mounted) {
                          ref.invalidate(accountDetailProvider(widget.accountId));
                        }
                      },
                      icon: const Icon(Icons.people_outline),
                      label: Text(l10n.viewContacts),
                      style: OutlinedButton.styleFrom(
                        minimumSize: const Size(double.infinity, 48),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                    ),
                  ),
                  if (ref.watch(canEditProvider('accounts'))) ...[
                    const SizedBox(width: 12),
                    Expanded(
                      child: FilledButton.icon(
                        onPressed: _isDeleting ? null : () => _handleEdit(account),
                        icon: const Icon(Icons.edit_outlined),
                        label: Text(l10n.edit),
                        style: FilledButton.styleFrom(
                          minimumSize: const Size(double.infinity, 48),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                      ),
                    ),
                  ],
                ],
              ),
              if (ref.watch(canDeleteProvider('accounts'))) ...[
                const SizedBox(height: 12),
                OutlinedButton.icon(
                  onPressed: _isDeleting ? null : () => _handleDelete(account),
                  icon: _isDeleting
                      ? const SizedBox(
                          width: 16,
                          height: 16,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Icon(Icons.delete_outline),
                  label: Text(l10n.delete),
                  style: OutlinedButton.styleFrom(
                    minimumSize: const Size(double.infinity, 48),
                    foregroundColor: colorScheme.error,
                    side: BorderSide(color: colorScheme.error),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                ),
              ],
            ],
          ),
        ),
        loading: () => const LoadingWidget(),
        error: (error, stack) => ErrorStateWidget(
          message: error.toString().replaceFirst('Exception: ', ''),
          onRetry: () => ref.invalidate(accountDetailProvider(widget.accountId)),
        ),
      ),
    );
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({
    required this.status,
    required this.theme,
    required this.colorScheme,
  });

  final String status;
  final ThemeData theme;
  final ColorScheme colorScheme;

  Color _getStatusColor() {
    final isDark = theme.brightness == Brightness.dark;
    final statusLower = status.toLowerCase();
    
    if (statusLower == 'active') {
      // Green for active - high contrast
      return isDark ? const Color(0xFF4CAF50) : const Color(0xFF2E7D32);
    } else {
      // Gray/Red for inactive - high contrast
      return isDark ? const Color(0xFF9E9E9E) : const Color(0xFF616161);
    }
  }

  Color _getStatusBackgroundColor() {
    final isDark = theme.brightness == Brightness.dark;
    final statusLower = status.toLowerCase();
    
    if (statusLower == 'active') {
      // Light green background
      return isDark 
          ? const Color(0xFF4CAF50).withValues(alpha: 0.2)
          : const Color(0xFF4CAF50).withValues(alpha: 0.15);
    } else {
      // Light gray background
      return isDark 
          ? const Color(0xFF9E9E9E).withValues(alpha: 0.2)
          : const Color(0xFF9E9E9E).withValues(alpha: 0.15);
    }
  }

  @override
  Widget build(BuildContext context) {
    final textColor = _getStatusColor();
    final backgroundColor = _getStatusBackgroundColor();

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        status.toUpperCase(),
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.w600,
          fontSize: 12,
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
    final isDark = theme.brightness == Brightness.dark;
    final typeLower = typeName.toLowerCase();
    
    if (typeLower == 'clinic') {
      // Blue for clinic
      return isDark ? const Color(0xFF42A5F5) : const Color(0xFF1976D2);
    } else if (typeLower == 'pharmacy') {
      // Orange/Amber for pharmacy
      return isDark ? const Color(0xFFFF9800) : const Color(0xFFF57C00);
    } else if (typeLower == 'hospital') {
      // Purple/Indigo for hospital
      return isDark ? const Color(0xFF9575CD) : const Color(0xFF5E35B1);
    } else {
      // Default to primary color
      return colorScheme.primary;
    }
  }

  Color _getTypeBackgroundColor() {
    final isDark = theme.brightness == Brightness.dark;
    final typeLower = typeName.toLowerCase();
    
    if (typeLower == 'clinic') {
      return isDark 
          ? const Color(0xFF42A5F5).withValues(alpha: 0.2)
          : const Color(0xFF42A5F5).withValues(alpha: 0.15);
    } else if (typeLower == 'pharmacy') {
      return isDark 
          ? const Color(0xFFFF9800).withValues(alpha: 0.2)
          : const Color(0xFFFF9800).withValues(alpha: 0.15);
    } else if (typeLower == 'hospital') {
      return isDark 
          ? const Color(0xFF9575CD).withValues(alpha: 0.2)
          : const Color(0xFF9575CD).withValues(alpha: 0.15);
    } else {
      return colorScheme.primary.withValues(alpha: isDark ? 0.2 : 0.15);
    }
  }

  @override
  Widget build(BuildContext context) {
    final textColor = _getTypeColor();
    final backgroundColor = _getTypeBackgroundColor();

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        typeName,
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.w500,
          fontSize: 12,
        ),
      ),
    );
  }
}

class _SectionTitle extends StatelessWidget {
  const _SectionTitle({
    required this.title,
    required this.theme,
    required this.colorScheme,
  });

  final String title;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    return Text(
      title,
      style: theme.textTheme.titleMedium?.copyWith(
        fontWeight: FontWeight.w600,
        color: colorScheme.onSurface,
      ),
    );
  }
}

class _InfoCard extends StatelessWidget {
  const _InfoCard({
    required this.children,
    required this.theme,
    required this.colorScheme,
  });

  final List<Widget> children;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: colorScheme.surface,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 3,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      padding: const EdgeInsets.all(16),
      child: Column(
        children: children.asMap().entries.map((entry) {
          final index = entry.key;
          final child = entry.value;
          if (index == children.length - 1) {
            return child;
          }
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: child,
          );
        }).toList(),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  const _InfoRow({
    required this.icon,
    required this.label,
    required this.value,
    required this.theme,
    required this.colorScheme,
  });

  final IconData icon;
  final String label;
  final String value;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(
          icon,
          size: 20,
          color: colorScheme.onSurface.withValues(alpha: 0.7),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: theme.textTheme.bodySmall?.copyWith(
                  color: colorScheme.onSurface.withValues(alpha: 0.7),
                ),
              ),
              const SizedBox(height: 4),
              Text(
                value,
                style: theme.textTheme.bodyMedium?.copyWith(
                  fontWeight: FontWeight.w500,
                  color: colorScheme.onSurface,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

