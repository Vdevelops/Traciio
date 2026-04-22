import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../application/contact_provider.dart';
import '../data/models/contact.dart';
import '../presentation/contact_form_screen.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/widgets/error_widget.dart';
import '../../../core/widgets/loading_widget.dart';
import '../../../core/permissions/permission_provider.dart';

class ContactDetailScreen extends ConsumerStatefulWidget {
  const ContactDetailScreen({
    super.key,
    required this.contactId,
  });

  final String contactId;

  @override
  ConsumerState<ContactDetailScreen> createState() => _ContactDetailScreenState();
}

class _ContactDetailScreenState extends ConsumerState<ContactDetailScreen> {
  bool _isDeleting = false;

  Future<void> _handleEdit(Contact contact) async {
    final result = await Navigator.push<Contact>(
      context,
      MaterialPageRoute(
        builder: (context) => ContactFormScreen(contact: contact),
      ),
    );

    if (result != null && mounted) {
      // Refresh detail
      ref.invalidate(contactDetailProvider(widget.contactId));
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(AppLocalizations.of(context)!.contactUpdatedSuccessfully),
          backgroundColor: Colors.green,
        ),
      );
    }
  }

  Future<void> _handleDelete(Contact contact) async {
    final l10n = AppLocalizations.of(context)!;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.deleteContact),
        content: Text(l10n.deleteContactConfirmation),
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
      final success = await ref.read(contactListProvider.notifier).deleteContact(contact.id);
      setState(() => _isDeleting = false);

      if (mounted) {
        if (success) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(l10n.contactDeleted),
              backgroundColor: Colors.green,
            ),
          );
          // Pop back and refresh contact list
          Navigator.pop(context, true); // Pass true to indicate refresh needed
        } else {
          final error = ref.read(contactListProvider).errorMessage;
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(error ?? 'Failed to delete contact'),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final contactAsync = ref.watch(contactDetailProvider(widget.contactId));
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.contactDetails),
        elevation: 0,
      ),
      body: contactAsync.when(
        loading: () => const LoadingWidget(),
        error: (error, stack) => ErrorStateWidget(
          message: error.toString().replaceFirst('Exception: ', ''),
          onRetry: () => ref.invalidate(contactDetailProvider(widget.contactId)),
        ),
        data: (contact) => SingleChildScrollView(
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
                child: Row(
                  children: [
                    // Avatar
                    Container(
                      width: 64,
                      height: 64,
                      decoration: BoxDecoration(
                        color: colorScheme.primary.withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(32),
                      ),
                      child: Icon(
                        Icons.person_outline,
                        color: colorScheme.primary,
                        size: 32,
                      ),
                    ),
                    const SizedBox(width: 16),
                    // Name & Position
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            contact.name,
                            style: theme.textTheme.headlineSmall?.copyWith(
                              fontWeight: FontWeight.bold,
                              color: colorScheme.onSurface,
                            ),
                          ),
                          if (contact.position != null) ...[
                            const SizedBox(height: 4),
                            Text(
                              contact.position!,
                              style: theme.textTheme.bodyMedium?.copyWith(
                                color: colorScheme.onSurface.withValues(alpha: 0.7),
                              ),
                            ),
                          ],
                          if (contact.role != null) ...[
                            const SizedBox(height: 8),
                            _RoleBadge(
                              roleName: contact.role!.name,
                              roleCode: contact.role!.code,
                              theme: theme,
                              colorScheme: colorScheme,
                            ),
                          ],
                        ],
                      ),
                    ),
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
                  if (contact.phone != null)
                    _InfoRow(
                      icon: Icons.phone_outlined,
                      label: l10n.phone,
                      value: contact.phone!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  if (contact.email != null)
                    _InfoRow(
                      icon: Icons.email_outlined,
                      label: l10n.email,
                      value: contact.email!,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                ],
              ),
              const SizedBox(height: 16),
              // Account Information
              if (contact.account != null) ...[
                _SectionTitle(
                  title: l10n.accounts,
                  theme: theme,
                  colorScheme: colorScheme,
                ),
                const SizedBox(height: 8),
                _InfoCard(
                  theme: theme,
                  colorScheme: colorScheme,
                  children: [
                    _InfoRow(
                      icon: Icons.business_outlined,
                      label: '${l10n.accounts} ${l10n.name.toLowerCase()}',
                      value: contact.account!.name,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                    if (contact.account!.city != null)
                      _InfoRow(
                        icon: Icons.location_city_outlined,
                        label: l10n.city,
                        value: contact.account!.city!,
                        theme: theme,
                        colorScheme: colorScheme,
                      ),
                  ],
                ),
                const SizedBox(height: 16),
              ],
              // Notes
              if (contact.notes != null && contact.notes!.isNotEmpty) ...[
                _SectionTitle(
                  title: l10n.notes,
                  theme: theme,
                  colorScheme: colorScheme,
                ),
                const SizedBox(height: 8),
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
                  child: Text(
                    contact.notes!,
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: colorScheme.onSurface,
                    ),
                  ),
                ),
              ],
              const SizedBox(height: 16),
              // Action Buttons
              if (ref.watch(canEditProvider('accounts'))) ...[
                Row(
                  children: [
                    Expanded(
                      child: FilledButton.icon(
                        onPressed: _isDeleting ? null : () => _handleEdit(contact),
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
                ),
                const SizedBox(height: 12),
              ],
              if (ref.watch(canDeleteProvider('accounts')))
                OutlinedButton.icon(
                  onPressed: _isDeleting ? null : () => _handleDelete(contact),
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
          ),
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

class _RoleBadge extends StatelessWidget {
  const _RoleBadge({
    required this.roleName,
    this.roleCode,
    required this.theme,
    required this.colorScheme,
  });

  final String roleName;
  final String? roleCode;
  final ThemeData theme;
  final ColorScheme colorScheme;

  Color _getRoleColor() {
    final isDark = theme.brightness == Brightness.dark;
    final nameLower = roleName.toLowerCase();
    final codeLower = roleCode?.toLowerCase() ?? '';
    
    // Check by code first (more reliable)
    if (codeLower.isNotEmpty) {
      if (codeLower.contains('doctor') || codeLower.contains('physician') || codeLower.contains('dr')) {
        // Teal/Cyan for doctor/physician
        return isDark ? const Color(0xFF26A69A) : const Color(0xFF00897B);
      } else if (codeLower.contains('nurse')) {
        // Pink/Rose for nurse
        return isDark ? const Color(0xFFEC407A) : const Color(0xFFC2185B);
      } else if (codeLower.contains('pharmacist') || codeLower.contains('apoteker')) {
        // Orange/Amber for pharmacist
        return isDark ? const Color(0xFFFF9800) : const Color(0xFFF57C00);
      } else if (codeLower.contains('manager') || codeLower.contains('admin')) {
        // Purple/Indigo for manager/admin
        return isDark ? const Color(0xFF9575CD) : const Color(0xFF5E35B1);
      } else if (codeLower.contains('staff') || codeLower.contains('employee')) {
        // Blue for staff/employee
        return isDark ? const Color(0xFF42A5F5) : const Color(0xFF1976D2);
      }
    }
    
    // Check by name if code not available
    if (nameLower.contains('doctor') || nameLower.contains('physician') || nameLower.contains('dokter')) {
      return isDark ? const Color(0xFF26A69A) : const Color(0xFF00897B);
    } else if (nameLower.contains('nurse') || nameLower.contains('perawat')) {
      return isDark ? const Color(0xFFEC407A) : const Color(0xFFC2185B);
    } else if (nameLower.contains('pharmacist') || nameLower.contains('apoteker')) {
      return isDark ? const Color(0xFFFF9800) : const Color(0xFFF57C00);
    } else if (nameLower.contains('manager') || nameLower.contains('admin') || nameLower.contains('manajer')) {
      return isDark ? const Color(0xFF9575CD) : const Color(0xFF5E35B1);
    } else if (nameLower.contains('staff') || nameLower.contains('employee') || nameLower.contains('karyawan')) {
      return isDark ? const Color(0xFF42A5F5) : const Color(0xFF1976D2);
    }
    
    // Default to primary color
    return colorScheme.primary;
  }

  Color _getRoleBackgroundColor() {
    final isDark = theme.brightness == Brightness.dark;
    final nameLower = roleName.toLowerCase();
    final codeLower = roleCode?.toLowerCase() ?? '';
    
    // Check by code first
    if (codeLower.isNotEmpty) {
      if (codeLower.contains('doctor') || codeLower.contains('physician') || codeLower.contains('dr')) {
        return isDark 
            ? const Color(0xFF26A69A).withValues(alpha: 0.2)
            : const Color(0xFF26A69A).withValues(alpha: 0.15);
      } else if (codeLower.contains('nurse')) {
        return isDark 
            ? const Color(0xFFEC407A).withValues(alpha: 0.2)
            : const Color(0xFFEC407A).withValues(alpha: 0.15);
      } else if (codeLower.contains('pharmacist') || codeLower.contains('apoteker')) {
        return isDark 
            ? const Color(0xFFFF9800).withValues(alpha: 0.2)
            : const Color(0xFFFF9800).withValues(alpha: 0.15);
      } else if (codeLower.contains('manager') || codeLower.contains('admin')) {
        return isDark 
            ? const Color(0xFF9575CD).withValues(alpha: 0.2)
            : const Color(0xFF9575CD).withValues(alpha: 0.15);
      } else if (codeLower.contains('staff') || codeLower.contains('employee')) {
        return isDark 
            ? const Color(0xFF42A5F5).withValues(alpha: 0.2)
            : const Color(0xFF42A5F5).withValues(alpha: 0.15);
      }
    }
    
    // Check by name
    if (nameLower.contains('doctor') || nameLower.contains('physician') || nameLower.contains('dokter')) {
      return isDark 
          ? const Color(0xFF26A69A).withValues(alpha: 0.2)
          : const Color(0xFF26A69A).withValues(alpha: 0.15);
    } else if (nameLower.contains('nurse') || nameLower.contains('perawat')) {
      return isDark 
          ? const Color(0xFFEC407A).withValues(alpha: 0.2)
          : const Color(0xFFEC407A).withValues(alpha: 0.15);
    } else if (nameLower.contains('pharmacist') || nameLower.contains('apoteker')) {
      return isDark 
          ? const Color(0xFFFF9800).withValues(alpha: 0.2)
          : const Color(0xFFFF9800).withValues(alpha: 0.15);
    } else if (nameLower.contains('manager') || nameLower.contains('admin') || nameLower.contains('manajer')) {
      return isDark 
          ? const Color(0xFF9575CD).withValues(alpha: 0.2)
          : const Color(0xFF9575CD).withValues(alpha: 0.15);
    } else if (nameLower.contains('staff') || nameLower.contains('employee') || nameLower.contains('karyawan')) {
      return isDark 
          ? const Color(0xFF42A5F5).withValues(alpha: 0.2)
          : const Color(0xFF42A5F5).withValues(alpha: 0.15);
    }
    
    // Default to primary color
    return colorScheme.primary.withValues(alpha: isDark ? 0.2 : 0.15);
  }

  @override
  Widget build(BuildContext context) {
    final textColor = _getRoleColor();
    final backgroundColor = _getRoleBackgroundColor();

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        roleName,
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.w500,
          fontSize: 12,
        ),
      ),
    );
  }
}
