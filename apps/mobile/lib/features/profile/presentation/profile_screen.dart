import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_svg/flutter_svg.dart';

import '../../../core/l10n/app_localizations.dart';
import '../../../core/l10n/locale_provider.dart';
import '../../../core/routing/app_router.dart';
import '../../../core/theme/theme_provider.dart';
import '../../../core/utils/app_info.dart';
import '../../../core/widgets/error_widget.dart';
// Profile hanya diakses via quick nav - tidak pakai bottom navbar
import '../../../core/permissions/permission_provider.dart';
import '../../auth/application/auth_provider.dart' as auth;
import '../../auth/data/models/login_response.dart';
import '../../google_calendar/application/google_calendar_provider.dart'
    as google_calendar;
import '../../google_calendar/presentation/widgets/google_calendar_connection_widget.dart';
import '../application/profile_provider.dart';
import '../data/models/profile.dart';
import '../data/profile_repository.dart';

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  String _getThemeLabel(ThemeMode themeMode, AppLocalizations l10n) {
    switch (themeMode) {
      case ThemeMode.light:
        return l10n.lightTheme;
      case ThemeMode.dark:
        return l10n.darkTheme;
      case ThemeMode.system:
        return l10n.systemTheme;
    }
  }

  String _getLanguageLabel(Locale locale, AppLocalizations l10n) {
    switch (locale.languageCode) {
      case 'id':
        return l10n.indonesian;
      case 'en':
      default:
        return l10n.english;
    }
  }

  IconData _getThemeIcon(ThemeMode themeMode) {
    switch (themeMode) {
      case ThemeMode.light:
        return Icons.light_mode_outlined;
      case ThemeMode.dark:
        return Icons.dark_mode_outlined;
      case ThemeMode.system:
        return Icons.brightness_auto_outlined;
    }
  }

  void _showThemeModal(
    BuildContext context,
    WidgetRef ref,
    ThemeMode currentThemeMode,
    ThemeModeNotifier themeNotifier,
    AppLocalizations l10n,
  ) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => _ThemeSelectorModal(
        currentThemeMode: currentThemeMode,
        onThemeSelected: (mode) async {
          await themeNotifier.setThemeMode(mode);
          if (context.mounted) {
            Navigator.pop(context);
          }
        },
        l10n: l10n,
      ),
    );
  }

  void _showLanguageModal(
    BuildContext context,
    WidgetRef ref,
    Locale currentLocale,
    LocaleNotifier localeNotifier,
    AppLocalizations l10n,
  ) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => _LanguageSelectorModal(
        currentLocale: currentLocale,
        onLanguageSelected: (locale) async {
          await localeNotifier.setLocale(locale);
          if (context.mounted) {
            Navigator.pop(context);
          }
        },
        l10n: l10n,
      ),
    );
  }

  void _showEditProfileDialog(
    BuildContext context,
    WidgetRef ref,
    String currentName,
    AppLocalizations l10n,
  ) {
    final nameController = TextEditingController(text: currentName);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Edit Profile'),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        content: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: nameController,
                decoration: InputDecoration(
                  labelText: l10n.name,
                  hintText: 'Enter your name',
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  filled: true,
                  fillColor: Theme.of(
                    context,
                  ).colorScheme.surfaceContainerHighest,
                  contentPadding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 16,
                  ),
                ),
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(l10n.cancel),
          ),
          FilledButton(
            onPressed: () async {
              final name = nameController.text.trim();
              if (name.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: const Text('Name is required'),
                    backgroundColor: Theme.of(context).colorScheme.error,
                  ),
                );
                return;
              }

              try {
                final request = UpdateProfileRequest(name: name);
                final updatedUser = await ref.read(
                  updateProfileProvider(request).future,
                );

                // Update auth state with new user data
                final authNotifier = ref.read(auth.authProvider.notifier);

                // Convert ProfileUser to UserResponse for auth state
                final userResponse = UserResponse(
                  id: updatedUser.id,
                  email: updatedUser.email,
                  name: updatedUser.name,
                  avatarUrl: updatedUser.avatarUrl,
                  role: updatedUser.role?.name ?? '',
                  status: updatedUser.status,
                  permissions:
                      ref.read(auth.authProvider).user?.permissions ?? [],
                  createdAt: updatedUser.createdAt,
                  updatedAt: updatedUser.updatedAt,
                );

                // Update auth state (this will also save to storage)
                await authNotifier.updateUser(userResponse);

                // Refresh profile
                ref.invalidate(profileProvider);

                if (context.mounted) {
                  Navigator.pop(context);
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: const Text('Profile updated successfully'),
                      backgroundColor: Colors.green,
                    ),
                  );
                }
              } catch (e) {
                if (context.mounted) {
                  final errorMessage = _extractErrorMessage(e);
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(errorMessage),
                      backgroundColor: Theme.of(context).colorScheme.error,
                      behavior: SnackBarBehavior.floating,
                      margin: const EdgeInsets.all(16),
                      duration: const Duration(seconds: 4),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                    ),
                  );
                }
              }
            },
            child: Text(l10n.save),
          ),
        ],
      ),
    );
  }

  void _showChangePasswordDialog(
    BuildContext context,
    WidgetRef ref,
    AppLocalizations l10n,
  ) {
    final currentPasswordController = TextEditingController();
    final newPasswordController = TextEditingController();
    final confirmPasswordController = TextEditingController();
    bool obscureCurrentPassword = true;
    bool obscureNewPassword = true;
    bool obscureConfirmPassword = true;

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: const Text('Change Password'),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextField(
                  controller: currentPasswordController,
                  obscureText: obscureCurrentPassword,
                  decoration: InputDecoration(
                    labelText: 'Current Password',
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    filled: true,
                    fillColor: Theme.of(
                      context,
                    ).colorScheme.surfaceContainerHighest,
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 16,
                    ),
                    suffixIcon: IconButton(
                      icon: Icon(
                        obscureCurrentPassword
                            ? Icons.visibility_outlined
                            : Icons.visibility_off_outlined,
                      ),
                      onPressed: () {
                        setState(() {
                          obscureCurrentPassword = !obscureCurrentPassword;
                        });
                      },
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: newPasswordController,
                  obscureText: obscureNewPassword,
                  decoration: InputDecoration(
                    labelText: 'New Password',
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    filled: true,
                    fillColor: Theme.of(
                      context,
                    ).colorScheme.surfaceContainerHighest,
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 16,
                    ),
                    suffixIcon: IconButton(
                      icon: Icon(
                        obscureNewPassword
                            ? Icons.visibility_outlined
                            : Icons.visibility_off_outlined,
                      ),
                      onPressed: () {
                        setState(() {
                          obscureNewPassword = !obscureNewPassword;
                        });
                      },
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: confirmPasswordController,
                  obscureText: obscureConfirmPassword,
                  decoration: InputDecoration(
                    labelText: 'Confirm Password',
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    filled: true,
                    fillColor: Theme.of(
                      context,
                    ).colorScheme.surfaceContainerHighest,
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 16,
                    ),
                    suffixIcon: IconButton(
                      icon: Icon(
                        obscureConfirmPassword
                            ? Icons.visibility_outlined
                            : Icons.visibility_off_outlined,
                      ),
                      onPressed: () {
                        setState(() {
                          obscureConfirmPassword = !obscureConfirmPassword;
                        });
                      },
                    ),
                  ),
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text(l10n.cancel),
            ),
            FilledButton(
              onPressed: () async {
                final currentPassword = currentPasswordController.text;
                final newPassword = newPasswordController.text;
                final confirmPassword = confirmPasswordController.text;

                if (currentPassword.isEmpty ||
                    newPassword.isEmpty ||
                    confirmPassword.isEmpty) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: const Text('All fields are required'),
                      backgroundColor: Theme.of(context).colorScheme.error,
                    ),
                  );
                  return;
                }

                if (newPassword != confirmPassword) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: const Text('Passwords do not match'),
                      backgroundColor: Theme.of(context).colorScheme.error,
                    ),
                  );
                  return;
                }

                try {
                  final request = ChangePasswordRequest(
                    currentPassword: currentPassword,
                    password: newPassword,
                    confirmPassword: confirmPassword,
                  );
                  await ref.read(changePasswordProvider(request).future);

                  if (context.mounted) {
                    Navigator.pop(context);
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(
                        content: const Text('Password changed successfully'),
                        backgroundColor: Colors.green,
                        behavior: SnackBarBehavior.floating,
                        margin: const EdgeInsets.all(16),
                        duration: const Duration(seconds: 3),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                    );
                  }
                } catch (e) {
                  if (context.mounted) {
                    final errorMessage = _extractErrorMessage(e);
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(
                        content: Text(errorMessage),
                        backgroundColor: Theme.of(context).colorScheme.error,
                        behavior: SnackBarBehavior.floating,
                        margin: const EdgeInsets.all(16),
                        duration: const Duration(seconds: 4),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                    );
                  }
                }
              },
              child: Text(l10n.save),
            ),
          ],
        ),
      ),
    );
  }

  String _extractErrorMessage(dynamic error) {
    // Handle ProfileException (custom exception)
    if (error is ProfileException) {
      return error.message;
    }

    // Handle regular Exception
    if (error is Exception) {
      final errorString = error.toString();
      // Remove "Exception: " prefix if present
      if (errorString.startsWith('Exception: ')) {
        return errorString.substring(11);
      }
      return errorString;
    }

    // Handle String errors
    if (error is String) {
      return error;
    }

    // Fallback to toString()
    final errorString = error.toString();
    if (errorString.startsWith('Exception: ')) {
      return errorString.substring(11);
    }

    return errorString;
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final profileAsync = ref.watch(profileProvider);
    final authState = ref.watch(auth.authProvider);
    final user = authState.user;
    final themeMode = ref.watch(themeModeProvider);
    final themeNotifier = ref.read(themeModeProvider.notifier);
    final locale = ref.watch(localeProvider);
    final localeNotifier = ref.read(localeProvider.notifier);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: null,
        automaticallyImplyLeading: false,
      ),
      body: SafeArea(
        child: profileAsync.when(
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (error, stackTrace) => ErrorStateWidget(
            message: _extractErrorMessage(error),
            onRetry: () {
              ref.invalidate(profileProvider);
            },
          ),
          data: (profile) {
            // Get profile data from API, fallback to auth state user
            final profileUser = profile.user;
            final userName = profileUser.name.isNotEmpty
                ? profileUser.name
                : (user?.name ?? 'User');
            final userEmail = profileUser.email.isNotEmpty
                ? profileUser.email
                : (user?.email ?? 'user@example.com');
            final userRole = profileUser.role?.name.isNotEmpty == true
                ? profileUser.role!.name
                : (user?.role ?? 'Sales Rep');
            final avatarUrl = profileUser.avatarUrl ?? user?.avatarUrl;

            // Get current theme mode label
            final themeLabel = _getThemeLabel(themeMode, l10n);
            final themeIcon = _getThemeIcon(themeMode);
            final languageLabel = _getLanguageLabel(locale, l10n);

            return RefreshIndicator(
              onRefresh: () async {
                // Refresh profile data
                ref.invalidate(profileProvider);
                // Refresh Google Calendar status
                await ref.read(
                  google_calendar.googleCalendarStatusProvider.future,
                );
              },
              child: SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                padding: const EdgeInsets.fromLTRB(20, 24, 20, 100),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Profile Header - Modern & Minimalistic
                    SizedBox(
                      width: double.infinity,
                      child: Container(
                        padding: const EdgeInsets.all(24),
                        decoration: BoxDecoration(
                          color: theme.colorScheme.surface,
                          borderRadius: BorderRadius.circular(20),
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
                            // Avatar
                            Container(
                              width: 80,
                              height: 80,
                              decoration: BoxDecoration(
                                shape: BoxShape.circle,
                                boxShadow: [
                                  BoxShadow(
                                    color: theme.colorScheme.primary.withValues(
                                      alpha: 0.15,
                                    ),
                                    blurRadius: 3,
                                    offset: const Offset(0, 1),
                                  ),
                                ],
                              ),
                              child: ClipOval(
                                child: avatarUrl != null && avatarUrl.isNotEmpty
                                    ? _buildAvatarImage(
                                        avatarUrl,
                                        userName,
                                        theme,
                                      )
                                    : Container(
                                        decoration: BoxDecoration(
                                          gradient: LinearGradient(
                                            begin: Alignment.topLeft,
                                            end: Alignment.bottomRight,
                                            colors: [
                                              theme.colorScheme.primary,
                                              theme.colorScheme.primary
                                                  .withValues(alpha: 0.8),
                                            ],
                                          ),
                                        ),
                                        child: Center(
                                          child: Text(
                                            userName
                                                .substring(0, 1)
                                                .toUpperCase(),
                                            style: theme.textTheme.headlineLarge
                                                ?.copyWith(
                                                  color: Colors.white,
                                                  fontWeight: FontWeight.bold,
                                                  fontSize: 32,
                                                ),
                                          ),
                                        ),
                                      ),
                              ),
                            ),
                            const SizedBox(width: 20),
                            // User Info
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  // Name
                                  Text(
                                    userName,
                                    style: theme.textTheme.titleLarge?.copyWith(
                                      fontWeight: FontWeight.bold,
                                      color: theme.colorScheme.onSurface,
                                      fontSize: 20,
                                    ),
                                  ),
                                  const SizedBox(height: 6),
                                  // Email
                                  Text(
                                    userEmail,
                                    style: theme.textTheme.bodyMedium?.copyWith(
                                      color: theme.colorScheme.onSurface
                                          .withValues(alpha: 0.7),
                                      fontSize: 14,
                                    ),
                                  ),
                                  const SizedBox(height: 8),
                                  // Role Badge
                                  Container(
                                    padding: const EdgeInsets.symmetric(
                                      horizontal: 10,
                                      vertical: 4,
                                    ),
                                    decoration: BoxDecoration(
                                      color: theme.colorScheme.primary
                                          .withValues(alpha: 0.1),
                                      borderRadius: BorderRadius.circular(16),
                                    ),
                                    child: Text(
                                      userRole,
                                      style: theme.textTheme.bodySmall
                                          ?.copyWith(
                                            color: theme.colorScheme.primary,
                                            fontWeight: FontWeight.w600,
                                            fontSize: 12,
                                          ),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(height: 32),
                    // Settings Section
                    _SectionTitle(title: l10n.settings),
                    const SizedBox(height: 16),
                    // Google Calendar Connection
                    const _GoogleCalendarSettingsCard(),
                    const SizedBox(height: 16),
                    // Edit Profile
                    _SettingsCard(
                      child: _SettingsTile(
                        icon: Icons.edit_outlined,
                        title: 'Edit Profile',
                        subtitle: 'Update your profile information',
                        onTap: () => _showEditProfileDialog(
                          context,
                          ref,
                          userName,
                          l10n,
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),
                    // Change Password (only if user has permission)
                    Builder(
                      builder: (context) {
                        final canChangePassword = ref.watch(
                          canChangePasswordProvider,
                        );
                        if (!canChangePassword) {
                          return const SizedBox.shrink();
                        }
                        return _SettingsCard(
                          child: _SettingsTile(
                            icon: Icons.lock_outline,
                            title: 'Change Password',
                            subtitle: 'Update your password',
                            onTap: () =>
                                _showChangePasswordDialog(context, ref, l10n),
                          ),
                        );
                      },
                    ),
                    const SizedBox(height: 16),
                    const SizedBox(height: 16),
                    // Notifications
                    _SettingsCard(
                      child: _SettingsTile(
                        icon: Icons.notifications_outlined,
                        title: l10n.notifications,
                        subtitle: l10n.manageNotificationSettings,
                        onTap: () {
                          // Note: Notification navigation will be implemented when notification screen is available
                        },
                      ),
                    ),
                    const SizedBox(height: 16),
                    // Language
                    _SettingsCard(
                      child: _SettingsTile(
                        icon: Icons.language_outlined,
                        title: l10n.language,
                        subtitle: languageLabel,
                        onTap: () {
                          _showLanguageModal(
                            context,
                            ref,
                            locale,
                            localeNotifier,
                            l10n,
                          );
                        },
                      ),
                    ),
                    const SizedBox(height: 16),
                    // Theme
                    _SettingsCard(
                      child: _SettingsTile(
                        icon: themeIcon,
                        title: l10n.theme,
                        subtitle: themeLabel,
                        onTap: () {
                          _showThemeModal(
                            context,
                            ref,
                            themeMode,
                            themeNotifier,
                            l10n,
                          );
                        },
                      ),
                    ),
                    const SizedBox(height: 32),
                    // About Section
                    _SectionTitle(title: l10n.about),
                    const SizedBox(height: 16),
                    // App Version
                    _SettingsCard(
                      child: _SettingsTile(
                        icon: Icons.info_outline,
                        title: l10n.appVersion,
                        subtitle: AppInfo.versionString,
                        onTap: null,
                      ),
                    ),
                    const SizedBox(height: 16),
                    // Privacy Policy
                    _SettingsCard(
                      child: _SettingsTile(
                        icon: Icons.privacy_tip_outlined,
                        title: l10n.privacyPolicy,
                        subtitle: l10n.viewPrivacyPolicy,
                        onTap: () {
                          // Note: Privacy policy navigation will be implemented when privacy policy page is available
                        },
                      ),
                    ),
                    const SizedBox(height: 16),
                    // Terms of Service
                    _SettingsCard(
                      child: _SettingsTile(
                        icon: Icons.description_outlined,
                        title: l10n.termsOfService,
                        subtitle: l10n.viewTermsOfService,
                        onTap: () {
                          // Note: Terms of service navigation will be implemented when terms of service page is available
                        },
                      ),
                    ),
                    const SizedBox(height: 32),
                    // Logout Button
                    SizedBox(
                      width: double.infinity,
                      child: OutlinedButton(
                        onPressed: () async {
                          final confirmed = await showDialog<bool>(
                            context: context,
                            builder: (context) => AlertDialog(
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(16),
                              ),
                              title: Text(l10n.logout),
                              content: Text(l10n.logoutConfirmation),
                              actions: [
                                TextButton(
                                  onPressed: () =>
                                      Navigator.pop(context, false),
                                  child: Text(l10n.cancel),
                                ),
                                FilledButton(
                                  onPressed: () => Navigator.pop(context, true),
                                  child: Text(l10n.logout),
                                ),
                              ],
                            ),
                          );

                          if (confirmed == true) {
                            await ref.read(auth.authProvider.notifier).logout();
                            if (context.mounted) {
                              Navigator.of(context).pushNamedAndRemoveUntil(
                                AppRoutes.login,
                                (route) => false,
                              );
                            }
                          }
                        },
                        style: OutlinedButton.styleFrom(
                          foregroundColor: Colors.red,
                          side: BorderSide(
                            color: Colors.red.withValues(alpha: 0.3),
                            width: 1.5,
                          ),
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(16),
                          ),
                        ),
                        child: Text(
                          l10n.logout,
                          style: const TextStyle(
                            fontWeight: FontWeight.w600,
                            fontSize: 16,
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(height: 24),
                  ],
                ),
              ),
            );
          },
        ),
      ),
    );
  }
}

Widget _buildAvatarImage(String avatarUrl, String userName, ThemeData theme) {
  // Check if URL is SVG - dicebear URLs contain .svg in the path
  final isSvg =
      avatarUrl.toLowerCase().contains('.svg') ||
      avatarUrl.toLowerCase().contains('dicebear');

  final fallbackWidget = Container(
    decoration: BoxDecoration(
      gradient: LinearGradient(
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
        colors: [
          theme.colorScheme.primary,
          theme.colorScheme.primary.withValues(alpha: 0.8),
        ],
      ),
    ),
    child: Center(
      child: Text(
        userName.substring(0, 1).toUpperCase(),
        style: theme.textTheme.headlineLarge?.copyWith(
          color: Colors.white,
          fontWeight: FontWeight.bold,
          fontSize: 32,
        ),
      ),
    ),
  );

  if (isSvg) {
    return SvgPicture.network(
      avatarUrl,
      width: 80,
      height: 80,
      fit: BoxFit.cover,
      placeholderBuilder: (context) => Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [
              theme.colorScheme.primary,
              theme.colorScheme.primary.withValues(alpha: 0.8),
            ],
          ),
        ),
        child: Center(
          child: CircularProgressIndicator(
            strokeWidth: 2,
            valueColor: const AlwaysStoppedAnimation<Color>(Colors.white),
          ),
        ),
      ),
    );
  } else {
    return Image.network(
      avatarUrl,
      width: 80,
      height: 80,
      fit: BoxFit.cover,
      errorBuilder: (context, error, stackTrace) {
        return fallbackWidget;
      },
      loadingBuilder: (context, child, loadingProgress) {
        if (loadingProgress == null) return child;
        return Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [
                theme.colorScheme.primary,
                theme.colorScheme.primary.withValues(alpha: 0.8),
              ],
            ),
          ),
          child: Center(
            child: CircularProgressIndicator(
              value: loadingProgress.expectedTotalBytes != null
                  ? loadingProgress.cumulativeBytesLoaded /
                        loadingProgress.expectedTotalBytes!
                  : null,
              strokeWidth: 2,
              valueColor: const AlwaysStoppedAnimation<Color>(Colors.white),
            ),
          ),
        );
      },
    );
  }
}

class _SectionTitle extends StatelessWidget {
  const _SectionTitle({required this.title});

  final String title;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Text(
      title,
      style: theme.textTheme.titleMedium?.copyWith(
        fontWeight: FontWeight.w600,
        color: theme.colorScheme.onSurface,
        fontSize: 16,
        letterSpacing: 0.2,
      ),
    );
  }
}

class _SettingsCard extends StatelessWidget {
  const _SettingsCard({required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 3,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: child,
    );
  }
}

class _SettingsTile extends StatelessWidget {
  const _SettingsTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    this.onTap,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(16),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: theme.colorScheme.onSurface.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(
                icon,
                color: theme.colorScheme.onSurface.withValues(alpha: 0.7),
                size: 20,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: theme.textTheme.bodyLarge?.copyWith(
                      fontWeight: FontWeight.w500,
                      color: theme.colorScheme.onSurface,
                      fontSize: 15,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    subtitle,
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: theme.colorScheme.onSurface.withValues(alpha: 0.7),
                      fontSize: 13,
                    ),
                  ),
                ],
              ),
            ),
            if (onTap != null)
              Icon(
                Icons.chevron_right,
                color: theme.colorScheme.onSurface.withValues(alpha: 0.5),
                size: 20,
              ),
          ],
        ),
      ),
    );
  }
}

class _ThemeSelectorModal extends StatelessWidget {
  const _ThemeSelectorModal({
    required this.currentThemeMode,
    required this.onThemeSelected,
    required this.l10n,
  });

  final ThemeMode currentThemeMode;
  final ValueChanged<ThemeMode> onThemeSelected;
  final AppLocalizations l10n;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      padding: const EdgeInsets.symmetric(vertical: 20),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Handle bar
          Container(
            width: 40,
            height: 4,
            margin: const EdgeInsets.only(bottom: 20),
            decoration: BoxDecoration(
              color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          // Title
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Text(
              l10n.selectTheme,
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
                color: theme.colorScheme.onSurface,
              ),
            ),
          ),
          const SizedBox(height: 24),
          // Theme Options
          _ThemeOption(
            icon: Icons.light_mode_outlined,
            title: l10n.lightTheme,
            subtitle: l10n.lightTheme,
            isSelected: currentThemeMode == ThemeMode.light,
            onTap: () => onThemeSelected(ThemeMode.light),
          ),
          const Divider(height: 1),
          _ThemeOption(
            icon: Icons.dark_mode_outlined,
            title: l10n.darkTheme,
            subtitle: l10n.darkTheme,
            isSelected: currentThemeMode == ThemeMode.dark,
            onTap: () => onThemeSelected(ThemeMode.dark),
          ),
          const Divider(height: 1),
          _ThemeOption(
            icon: Icons.brightness_auto_outlined,
            title: l10n.systemTheme,
            subtitle: l10n.systemTheme,
            isSelected: currentThemeMode == ThemeMode.system,
            onTap: () => onThemeSelected(ThemeMode.system),
          ),
          const SizedBox(height: 20),
        ],
      ),
    );
  }
}

class _ThemeOption extends StatelessWidget {
  const _ThemeOption({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.isSelected,
    required this.onTap,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final bool isSelected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
        child: Row(
          children: [
            Icon(
              icon,
              color: isSelected
                  ? theme.colorScheme.primary
                  : theme.colorScheme.onSurface.withValues(alpha: 0.7),
              size: 24,
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: theme.textTheme.bodyLarge?.copyWith(
                      fontWeight: isSelected
                          ? FontWeight.w600
                          : FontWeight.normal,
                      color: isSelected
                          ? theme.colorScheme.primary
                          : theme.colorScheme.onSurface,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    subtitle,
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: theme.colorScheme.onSurface.withValues(alpha: 0.7),
                      fontSize: 13,
                    ),
                  ),
                ],
              ),
            ),
            if (isSelected)
              Icon(
                Icons.check_circle,
                color: theme.colorScheme.primary,
                size: 24,
              ),
          ],
        ),
      ),
    );
  }
}

class _LanguageSelectorModal extends StatelessWidget {
  const _LanguageSelectorModal({
    required this.currentLocale,
    required this.onLanguageSelected,
    required this.l10n,
  });

  final Locale currentLocale;
  final ValueChanged<Locale> onLanguageSelected;
  final AppLocalizations l10n;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      padding: const EdgeInsets.symmetric(vertical: 20),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Handle bar
          Container(
            width: 40,
            height: 4,
            margin: const EdgeInsets.only(bottom: 20),
            decoration: BoxDecoration(
              color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          // Title
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Text(
              l10n.selectLanguage,
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
                color: theme.colorScheme.onSurface,
              ),
            ),
          ),
          const SizedBox(height: 24),
          // Language Options
          _LanguageOption(
            locale: const Locale('en', ''),
            title: l10n.english,
            subtitle: l10n.english,
            isSelected: currentLocale.languageCode == 'en',
            onTap: () => onLanguageSelected(const Locale('en', '')),
          ),
          const Divider(height: 1),
          _LanguageOption(
            locale: const Locale('id', ''),
            title: l10n.indonesian,
            subtitle: l10n.indonesian,
            isSelected: currentLocale.languageCode == 'id',
            onTap: () => onLanguageSelected(const Locale('id', '')),
          ),
          const SizedBox(height: 20),
        ],
      ),
    );
  }
}

class _LanguageOption extends StatelessWidget {
  const _LanguageOption({
    required this.locale,
    required this.title,
    required this.subtitle,
    required this.isSelected,
    required this.onTap,
  });

  final Locale locale;
  final String title;
  final String subtitle;
  final bool isSelected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    // Get flag emoji for locale
    final flagEmoji = locale.languageCode == 'id' ? '🇮🇩' : '🇬🇧';

    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
        child: Row(
          children: [
            // Flag emoji
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: theme.colorScheme.onSurface.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Center(
                child: Text(flagEmoji, style: const TextStyle(fontSize: 24)),
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: theme.textTheme.bodyLarge?.copyWith(
                      fontWeight: isSelected
                          ? FontWeight.w600
                          : FontWeight.normal,
                      color: isSelected
                          ? theme.colorScheme.primary
                          : theme.colorScheme.onSurface,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    subtitle,
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: theme.colorScheme.onSurface.withValues(alpha: 0.7),
                      fontSize: 13,
                    ),
                  ),
                ],
              ),
            ),
            if (isSelected)
              Icon(
                Icons.check_circle,
                color: theme.colorScheme.primary,
                size: 24,
              ),
          ],
        ),
      ),
    );
  }
}

/// Widget wrapper untuk Google Calendar Connection di Profile Screen
class _GoogleCalendarSettingsCard extends StatelessWidget {
  const _GoogleCalendarSettingsCard();

  @override
  Widget build(BuildContext context) {
    return const GoogleCalendarConnectionWidget();
  }
}
