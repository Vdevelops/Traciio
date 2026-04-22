import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../auth/application/auth_provider.dart';
import '../../../profile/application/profile_provider.dart';

/// Dashboard greeting header dengan nama user dan subtitle yang berubah setiap hari
class DashboardGreetingHeader extends ConsumerWidget {
  const DashboardGreetingHeader({super.key});

  /// Get subtitle yang berubah setiap hari berdasarkan hari dalam seminggu
  String _getDailySubtitle(Locale locale) {
    final now = DateTime.now();
    final dayOfWeek = now.weekday; // 1 = Monday, 7 = Sunday
    final isIndonesian = locale.languageCode == 'id';

    // List subtitle yang berbeda setiap hari
    final subtitles = isIndonesian
        ? [
            'Mari mulai minggu dengan semangat baru!', // Monday
            'Terus tingkatkan performa Anda!', // Tuesday
            'Setengah perjalanan minggu, tetap semangat!', // Wednesday
            'Hampir akhir minggu, jangan menyerah!', // Thursday
            'Akhir minggu sudah dekat, finish strong!', // Friday
            'Akhir pekan, waktu untuk recharge!', // Saturday
            'Persiapkan diri untuk minggu baru!', // Sunday
          ]
        : [
            'Start the week with fresh energy!', // Monday
            'Keep pushing your performance!', // Tuesday
            'Mid-week momentum, stay strong!', // Wednesday
            'Almost there, don\'t give up!', // Thursday
            'Weekend is near, finish strong!', // Friday
            'Weekend time, recharge yourself!', // Saturday
            'Prepare for a new week ahead!', // Sunday
          ];

    return subtitles[dayOfWeek - 1];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final locale = Localizations.localeOf(context);

    final authState = ref.watch(authProvider);
    final profileAsync = ref.watch(profileProvider);

    // Get user name from profile or auth state
    final userName = profileAsync.when(
      data: (profile) => profile.user.name.isNotEmpty
          ? profile.user.name
          : (authState.user?.name ?? 'Sales Pro'),
      loading: () => authState.user?.name ?? 'Sales Pro',
      error: (_, _) => authState.user?.name ?? 'Sales Pro',
    );

    // Get first name only (before space)
    final firstName = userName.split(' ').first;

    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 0, 20, 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Hello, $firstName',
            style: theme.textTheme.headlineMedium?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.onSurface,
              fontSize: 28,
              letterSpacing: -0.5,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            _getDailySubtitle(locale),
            style: theme.textTheme.bodyMedium?.copyWith(
              color: colorScheme.onSurface.withValues(alpha: 0.7),
              fontSize: 14,
            ),
          ),
        ],
      ),
    );
  }
}
