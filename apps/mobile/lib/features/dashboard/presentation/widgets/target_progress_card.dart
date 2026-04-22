import 'package:flutter/material.dart';
import '../../data/models/dashboard.dart';
import '../../../../core/l10n/app_localizations.dart';

class TargetProgressCard extends StatelessWidget {
  final TargetSummary data;

  const TargetProgressCard({super.key, required this.data});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    
    // Use orange theme color (#F39200) for Performance Goal card
    final primaryOrange = theme.colorScheme.primary;
    
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: primaryOrange, // Orange background like in the image
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: primaryOrange.withValues(alpha: 0.15),
            blurRadius: 4,
            offset: const Offset(0, 2),
            spreadRadius: 0,
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'PERFORMANCE GOAL',
                style: theme.textTheme.labelSmall?.copyWith(
                  color: Colors.white.withValues(alpha: 0.9),
                  fontWeight: FontWeight.bold,
                  letterSpacing: 0.5,
                  fontSize: 11,
                ),
              ),
              Icon(
                Icons.show_chart,
                color: Colors.white.withValues(alpha: 0.8),
                size: 20,
              ),
            ],
          ),
          const SizedBox(height: 16),
          Text(
            data.achievedAmountFormatted,
            style: theme.textTheme.headlineLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: Colors.white,
              letterSpacing: -1,
              fontSize: 32,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          const SizedBox(height: 16),
          ClipRRect(
            borderRadius: BorderRadius.circular(12),
            child: LinearProgressIndicator(
              value: (data.progressPercent / 100).clamp(0.0, 1.0),
              backgroundColor: Colors.white.withValues(alpha: 0.3),
              valueColor: const AlwaysStoppedAnimation<Color>(Colors.white),
              minHeight: 8,
            ),
          ),
          const SizedBox(height: 12),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '${data.progressPercent.toStringAsFixed(0)}% of Target',
                style: theme.textTheme.labelMedium?.copyWith(
                  color: Colors.white.withValues(alpha: 0.9),
                  fontWeight: FontWeight.w600,
                  fontSize: 12,
                ),
              ),
              Text(
                _formatRemaining(
                  data.targetAmount - data.achievedAmount,
                  l10n,
                ),
                style: theme.textTheme.labelMedium?.copyWith(
                  color: Colors.white.withValues(alpha: 0.9),
                  fontWeight: FontWeight.w600,
                  fontSize: 12,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _formatRemaining(int remaining, AppLocalizations l10n) {
    final isIndonesian = l10n.locale.languageCode == 'id';
    
    // Miliar (Billion): >= 1,000,000,000
    if (remaining >= 1000000000) {
      final value = (remaining / 1000000000);
      final formatted = value == value.toInt() 
          ? value.toInt().toString() 
          : value.toStringAsFixed(1);
      if (isIndonesian) {
        return 'Rp $formatted M ${l10n.tersisa}';
      } else {
        return 'Rp $formatted B ${l10n.remaining}';
      }
    }
    
    // Juta (Million): >= 1,000,000
    if (remaining >= 1000000) {
      final value = (remaining / 1000000);
      final formatted = value == value.toInt() 
          ? value.toInt().toString() 
          : value.toStringAsFixed(1);
      if (isIndonesian) {
        return 'Rp $formatted Jt ${l10n.tersisa}';
      } else {
        return 'Rp $formatted M ${l10n.remaining}';
      }
    }
    
    // Ribu (Thousand): >= 1,000
    if (remaining >= 1000) {
      final value = (remaining / 1000);
      final formatted = value == value.toInt() 
          ? value.toInt().toString() 
          : value.toStringAsFixed(1);
      if (isIndonesian) {
        return 'Rp $formatted Rb ${l10n.tersisa}';
      } else {
        return 'Rp $formatted K ${l10n.remaining}';
      }
    }
    
    // Kurang dari 1,000
    if (isIndonesian) {
      return 'Rp ${remaining.toStringAsFixed(0)} ${l10n.tersisa}';
    } else {
      return 'Rp ${remaining.toStringAsFixed(0)} ${l10n.remaining}';
    }
  }
}
