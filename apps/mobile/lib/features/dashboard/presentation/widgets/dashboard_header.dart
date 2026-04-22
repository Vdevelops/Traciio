import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_svg/flutter_svg.dart';

import '../../../../core/widgets/sync_status_indicator.dart';
import '../../../auth/application/auth_provider.dart';
import '../../../notifications/application/notification_provider.dart';
import '../../../notifications/presentation/notification_list_screen.dart';
import '../../../profile/application/profile_provider.dart';
import '../../../profile/presentation/profile_screen.dart';
import '../../application/dashboard_provider.dart';

class DashboardHeader extends ConsumerWidget {
  const DashboardHeader({
    super.key,
    this.currentTabIndex = 0,
    this.onSearchTap,
  });

  final int currentTabIndex;
  final VoidCallback? onSearchTap;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    final dashboardState = ref.watch(dashboardProvider);
    final userState = ref.watch(authProvider);
    final notificationCountState = ref.watch(notificationCountProvider);
    final profileAsync = ref.watch(profileProvider);
    final screenWidth = MediaQuery.of(context).size.width;

    return Column(
      children: [
        // Floating Bubble Header - Glassmorphism effect
        Container(
          padding: EdgeInsets.symmetric(
            horizontal: screenWidth < 360 ? 12 : 16,
            vertical: 10,
          ),
          decoration: BoxDecoration(
            color: theme.cardColor.withValues(alpha: 0.8),
            borderRadius: BorderRadius.circular(20),
            border: Border.all(
              color: colorScheme.outline.withValues(alpha: 0.1),
              width: 1,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.05),
                blurRadius: 3,
                offset: const Offset(0, 1),
                spreadRadius: 0,
              ),
            ],
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(20),
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
              child: Container(
                color: Colors.transparent,
                child: Row(
                  children: [
                    // App Logo
                    Container(
                      width: screenWidth < 360 ? 36 : 40,
                      height: screenWidth < 360 ? 36 : 40,
                      padding: const EdgeInsets.all(4),
                      child: Image.asset(
                        'assets/logo.png',
                        fit: BoxFit.contain,
                        errorBuilder: (context, error, stackTrace) {
                          return Icon(
                            Icons.business,
                            size: screenWidth < 360 ? 24 : 28,
                            color: colorScheme.primary,
                          );
                        },
                      ),
                    ),

                    // Loading Spinner (if loading) - menggunakan theme color
                    if (dashboardState.isLoadingOverview)
                      Padding(
                        padding: const EdgeInsets.only(left: 8),
                        child: SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2.5,
                            valueColor: AlwaysStoppedAnimation<Color>(
                              colorScheme.primary,
                            ),
                          ),
                        ),
                      )
                    else
                      SizedBox(width: screenWidth < 360 ? 8 : 12),

                    const Spacer(),

                    // Center Icons: Sync Status, Search (only on Tasks tab), Notifications
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        // Sync Status Indicator
                        const Padding(
                          padding: EdgeInsets.only(right: 8),
                          child: SyncStatusIndicator(
                            featureKey: 'dashboard',
                            iconSize: 20,
                          ),
                        ),

                        // Search Icon - Only show on Tasks tab
                        if (currentTabIndex == 1)
                          IconButton(
                            onPressed: onSearchTap,
                            icon: const Icon(Icons.search_outlined),
                            color: colorScheme.onSurface.withValues(alpha: 0.7),
                            iconSize: 22,
                            padding: EdgeInsets.zero,
                            constraints: const BoxConstraints(
                              minWidth: 40,
                              minHeight: 40,
                            ),
                          ),

                        // Notifications Icon with Badge
                        Stack(
                          clipBehavior: Clip.none,
                          children: [
                            IconButton(
                              onPressed: () {
                                Navigator.push(
                                  context,
                                  MaterialPageRoute(
                                    builder: (_) =>
                                        const NotificationListScreen(),
                                  ),
                                ).then((_) {
                                  if (context.mounted) {
                                    ref
                                        .read(
                                          notificationCountProvider.notifier,
                                        )
                                        .loadUnreadCount();
                                  }
                                });
                              },
                              icon: const Icon(Icons.notifications_outlined),
                              color: colorScheme.onSurface.withValues(
                                alpha: 0.7,
                              ),
                              iconSize: 22,
                              padding: EdgeInsets.zero,
                              constraints: const BoxConstraints(
                                minWidth: 40,
                                minHeight: 40,
                              ),
                            ),
                            if (!notificationCountState.isLoading &&
                                notificationCountState.unreadCount > 0)
                              Positioned(
                                right: 6,
                                top: 6,
                                child: Container(
                                  padding: const EdgeInsets.all(3),
                                  decoration: BoxDecoration(
                                    color: colorScheme.error,
                                    shape: BoxShape.circle,
                                    border: Border.all(
                                      color: theme.cardColor,
                                      width: 2,
                                    ),
                                  ),
                                  constraints: const BoxConstraints(
                                    minWidth: 14,
                                    minHeight: 14,
                                  ),
                                  child: Text(
                                    notificationCountState.unreadCount > 9
                                        ? '9+'
                                        : notificationCountState.unreadCount
                                              .toString(),
                                    style: const TextStyle(
                                      color: Colors.white,
                                      fontSize: 9,
                                      fontWeight: FontWeight.bold,
                                      height: 1.0,
                                    ),
                                    textAlign: TextAlign.center,
                                  ),
                                ),
                              ),
                          ],
                        ),
                      ],
                    ),

                    SizedBox(width: screenWidth < 360 ? 4 : 8),

                    // User Profile Icon
                    GestureDetector(
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (_) => const ProfileScreen(),
                          ),
                        );
                      },
                      child: Container(
                        width: screenWidth < 360 ? 36 : 40,
                        height: screenWidth < 360 ? 36 : 40,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          border: Border.all(
                            color: colorScheme.outline.withValues(alpha: 0.2),
                            width: 2,
                          ),
                        ),
                        child: ClipOval(
                          child: profileAsync.when(
                            data: (profile) {
                              // Use avatar from profile API if available, otherwise fallback to auth user
                              final url =
                                  profile.user.avatarUrl?.isNotEmpty == true
                                  ? profile.user.avatarUrl!
                                  : (userState.user?.avatarUrl ?? '');

                              if (url.isNotEmpty) {
                                return _buildAvatarImage(
                                  url,
                                  profile.user.name.isNotEmpty
                                      ? profile.user.name
                                      : (userState.user?.name ?? 'U'),
                                  theme,
                                  screenWidth < 360 ? 36 : 40,
                                  screenWidth < 360 ? 16 : 18,
                                );
                              }

                              return _buildAvatarFallback(
                                profile.user.name.isNotEmpty
                                    ? profile.user.name
                                    : (userState.user?.name ?? 'U'),
                                theme,
                                screenWidth < 360 ? 16 : 18,
                              );
                            },
                            loading: () {
                              // While loading, use auth user avatar if available
                              final url = userState.user?.avatarUrl ?? '';
                              if (url.isNotEmpty) {
                                return _buildAvatarImage(
                                  url,
                                  userState.user?.name ?? 'U',
                                  theme,
                                  screenWidth < 360 ? 36 : 40,
                                  screenWidth < 360 ? 16 : 18,
                                );
                              }
                              return _buildAvatarFallback(
                                userState.user?.name ?? 'U',
                                theme,
                                screenWidth < 360 ? 16 : 18,
                              );
                            },
                            error: (_, _) {
                              // On error, use auth user avatar if available
                              final url = userState.user?.avatarUrl ?? '';
                              if (url.isNotEmpty) {
                                return _buildAvatarImage(
                                  url,
                                  userState.user?.name ?? 'U',
                                  theme,
                                  screenWidth < 360 ? 36 : 40,
                                  screenWidth < 360 ? 16 : 18,
                                );
                              }
                              return _buildAvatarFallback(
                                userState.user?.name ?? 'U',
                                theme,
                                screenWidth < 360 ? 16 : 18,
                              );
                            },
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildAvatarImage(
    String avatarUrl,
    String userName,
    ThemeData theme,
    double size,
    double fontSize,
  ) {
    // Check if URL is SVG - dicebear URLs contain .svg in the path
    final isSvg =
        avatarUrl.toLowerCase().contains('.svg') ||
        avatarUrl.toLowerCase().contains('dicebear');

    final fallbackWidget = _buildAvatarFallback(userName, theme, fontSize);

    if (isSvg) {
      return SvgPicture.network(
        avatarUrl,
        width: size,
        height: size,
        fit: BoxFit.cover,
        placeholderBuilder: (context) => fallbackWidget,
      );
    } else {
      return Image.network(
        avatarUrl,
        width: size,
        height: size,
        fit: BoxFit.cover,
        errorBuilder: (context, error, stackTrace) {
          return fallbackWidget;
        },
        loadingBuilder: (context, child, loadingProgress) {
          if (loadingProgress == null) {
            return child;
          }
          return fallbackWidget;
        },
      );
    }
  }

  Widget _buildAvatarFallback(String name, ThemeData theme, double fontSize) {
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
        child: Text(
          name.substring(0, 1).toUpperCase(),
          style: theme.textTheme.titleLarge?.copyWith(
            color: Colors.white,
            fontWeight: FontWeight.bold,
            fontSize: fontSize,
          ),
        ),
      ),
    );
  }
}
