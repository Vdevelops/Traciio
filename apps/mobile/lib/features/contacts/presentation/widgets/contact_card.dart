import 'package:flutter/material.dart';

import '../../data/models/contact.dart';

class ContactCard extends StatelessWidget {
  const ContactCard({super.key, required this.contact, required this.onTap});

  final Contact contact;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: theme.brightness == Brightness.dark
            ? colorScheme.surfaceContainerHighest.withValues(alpha: 0.1)
            : colorScheme.surface,
        borderRadius: BorderRadius.circular(24),
        border: theme.brightness == Brightness.dark
            ? Border.all(color: colorScheme.outlineVariant.withValues(alpha: 0.2))
            : null,
        boxShadow: theme.brightness == Brightness.dark
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
            child: Row(
              children: [
                // Avatar
                Container(
                  width: 48,
                  height: 48,
                  decoration: BoxDecoration(
                    color: colorScheme.primary.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(24),
                  ),
                  child: Icon(
                    Icons.person_outline,
                    color: colorScheme.primary,
                    size: 24,
                  ),
                ),
                const SizedBox(width: 12),
                // Content
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        contact.name,
                        style: theme.textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.w600,
                          color: colorScheme.onSurface,
                        ),
                      ),
                      if (contact.position != null) ...[
                        const SizedBox(height: 4),
                        Text(
                          contact.position!,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: colorScheme.onSurface.withValues(alpha: 0.7),
                          ),
                        ),
                      ],
                      if (contact.account != null) ...[
                        const SizedBox(height: 4),
                        Row(
                          children: [
                            Icon(
                              Icons.business_outlined,
                              size: 14,
                              color: colorScheme.onSurface.withValues(alpha: 0.7),
                            ),
                            const SizedBox(width: 4),
                            Expanded(
                              child: Text(
                                contact.account!.name,
                                style: theme.textTheme.bodySmall?.copyWith(
                                  color: colorScheme.onSurface.withValues(alpha: 0.7),
                                ),
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                              ),
                            ),
                          ],
                        ),
                      ],
                    ],
                  ),
                ),
                // Role Badge
                if (contact.role != null)
                  _RoleBadge(
                    roleName: contact.role!.name,
                    roleCode: contact.role!.code,
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
              ],
            ),
          ),
        ),
      ),
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
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        roleName,
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.w500,
          fontSize: 10,
        ),
      ),
    );
  }
}
