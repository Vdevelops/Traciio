import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/routing/app_router.dart';
import '../application/lead_provider.dart';
import '../data/models/lead_model.dart';

class LeadDetailScreen extends ConsumerWidget {
  const LeadDetailScreen({super.key, required this.leadId});

  final String leadId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = AppLocalizations.of(context)!;
    final leadAsync = ref.watch(leadDetailProvider(leadId));

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.leadDetails),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit),
            onPressed: () {
              Navigator.pushNamed(
                context,
                AppRoutes.leadsForm,
                arguments: leadId,
              );
            },
          ),
        ],
      ),
      body: leadAsync.when(
        data: (lead) => _buildBody(context, ref, lead),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(child: Text(error.toString())),
      ),
    );
  }

  Widget _buildBody(BuildContext context, WidgetRef ref, Lead lead) {
    final l10n = AppLocalizations.of(context)!;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildHeader(context, lead),
          const SizedBox(height: 24),
          _buildInfoSection(context, lead),
          const SizedBox(height: 24),
          if (lead.notes != null && lead.notes!.isNotEmpty) ...[
            _buildNotesSection(context, lead),
            const SizedBox(height: 24),
          ],
          if (lead.status != 'converted')
            SizedBox(
              width: double.infinity,
              child: FilledButton.icon(
                onPressed: () => _showConvertDialog(context, lead),
                icon: const Icon(Icons.check_circle_outline),
                label: Text(l10n.convertLead),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildHeader(BuildContext context, Lead lead) {
    final theme = Theme.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Expanded(
              child: Text(
                lead.name,
                style: theme.textTheme.headlineSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            _StatusBadge(status: lead.status),
          ],
        ),
        if (lead.company.isNotEmpty) ...[
          const SizedBox(height: 8),
          Text(
            lead.company,
            style: theme.textTheme.titleMedium?.copyWith(
              color: theme.colorScheme.primary,
            ),
          ),
        ],
      ],
    );
  }

  Widget _buildInfoSection(BuildContext context, Lead lead) {
    final l10n = AppLocalizations.of(context)!;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (lead.email != null)
          _InfoRow(
            icon: Icons.email_outlined,
            label: l10n.email,
            value: lead.email!,
            onTap: () => _launchUrl('mailto:${lead.email}'),
          ),
        if (lead.phone != null)
          _InfoRow(
            icon: Icons.phone_outlined,
            label: l10n.phone,
            value: lead.phone!,
            onTap: () => _launchUrl('tel:${lead.phone}'),
          ),
        if (lead.industry != null)
          _InfoRow(
            icon: Icons.business_outlined,
            label: l10n.industry,
            value: lead.industry!,
          ),
        _InfoRow(
          icon: Icons.source_outlined,
          label: l10n.source,
          value: lead.source,
        ),
        if (lead.province != null)
          _InfoRow(
            icon: Icons.location_on_outlined,
            label: l10n.province,
            value: lead.province!,
          ),
        if (lead.address != null)
          _InfoRow(
            icon: Icons.map_outlined,
            label: l10n.address,
            value: lead.address!,
          ),
      ],
    );
  }

  Widget _buildNotesSection(BuildContext context, Lead lead) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.notes,
          style: theme.textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Container(
          width: double.infinity,
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: theme.colorScheme.surfaceContainerHighest.withValues(
              alpha: 0.5,
            ),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Text(lead.notes!, style: theme.textTheme.bodyMedium),
        ),
      ],
    );
  }

  Future<void> _launchUrl(String urlString) async {
    final url = Uri.parse(urlString);
    if (await canLaunchUrl(url)) {
      await launchUrl(url);
    }
  }

  void _showConvertDialog(BuildContext context, Lead lead) {
    Navigator.of(context).pushNamed(AppRoutes.leadsConvert, arguments: lead);
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({required this.status});

  final String status;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    Color backgroundColor;
    Color foregroundColor;

    switch (status.toLowerCase()) {
      case 'new':
        backgroundColor = colorScheme.primaryContainer;
        foregroundColor = colorScheme.onPrimaryContainer;
        break;
      case 'contacted':
        backgroundColor = colorScheme.tertiaryContainer;
        foregroundColor = colorScheme.onTertiaryContainer;
        break;
      case 'qualified':
        backgroundColor = Colors.green.shade100;
        foregroundColor = Colors.green.shade900;
        break;
      case 'lost':
        backgroundColor = colorScheme.errorContainer;
        foregroundColor = colorScheme.onErrorContainer;
        break;
      case 'converted':
        backgroundColor = colorScheme.secondaryContainer;
        foregroundColor = colorScheme.onSecondaryContainer;
        break;
      default:
        backgroundColor = colorScheme.surfaceContainerHighest;
        foregroundColor = colorScheme.onSurfaceVariant;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        status.toUpperCase(),
        style: theme.textTheme.labelSmall?.copyWith(
          color: foregroundColor,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  const _InfoRow({
    required this.icon,
    required this.label,
    required this.value,
    this.onTap,
  });

  final IconData icon;
  final String label;
  final String value;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(8),
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 4),
          child: Row(
            children: [
              Icon(icon, size: 20, color: theme.colorScheme.onSurfaceVariant),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      label,
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                    Text(
                      value,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: onTap != null ? theme.colorScheme.primary : null,
                        decoration: onTap != null
                            ? TextDecoration.underline
                            : null,
                      ),
                    ),
                  ],
                ),
              ),
              if (onTap != null)
                Icon(
                  Icons.arrow_forward_ios,
                  size: 12,
                  color: theme.colorScheme.onSurfaceVariant,
                ),
            ],
          ),
        ),
      ),
    );
  }
}
