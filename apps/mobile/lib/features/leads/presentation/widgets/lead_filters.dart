import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../application/lead_provider.dart';
import '../../../../core/l10n/app_localizations.dart';

class StatusFilterSheet extends ConsumerWidget {
  const StatusFilterSheet({
    super.key,
    this.selectedStatus,
    required this.onSelect,
  });

  final String? selectedStatus;
  final ValueChanged<String?> onSelect;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final colorScheme = theme.colorScheme;

    final formDataAsync = ref.watch(leadFormDataProvider);

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            l10n.filterByStatus,
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.onSurface,
            ),
          ),
          const SizedBox(height: 16),
          formDataAsync.when(
            data: (formData) {
              final statuses = [
                null,
                ...formData.leadStatuses.map((e) => e.value),
              ];
              return Column(
                mainAxisSize: MainAxisSize.min,
                children: statuses.map((status) {
                  final label = status == null
                      ? l10n.all
                      : formData.leadStatuses
                            .firstWhere((e) => e.value == status)
                            .label;
                  final isSelected = selectedStatus == status;
                  return ListTile(
                    title: Text(label),
                    trailing: isSelected
                        ? Icon(Icons.check, color: colorScheme.primary)
                        : null,
                    onTap: () => onSelect(status),
                  );
                }).toList(),
              );
            },
            loading: () => const Center(
              child: Padding(
                padding: EdgeInsets.all(16.0),
                child: CircularProgressIndicator(),
              ),
            ),
            error: (error, _) => Center(child: Text(error.toString())),
          ),
        ],
      ),
    );
  }
}

class SourceFilterSheet extends ConsumerWidget {
  const SourceFilterSheet({
    super.key,
    this.selectedSource,
    required this.onSelect,
  });

  final String? selectedSource;
  final ValueChanged<String?> onSelect;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final colorScheme = theme.colorScheme;

    final formDataAsync = ref.watch(leadFormDataProvider);

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            l10n.filterBySource,
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.onSurface,
            ),
          ),
          const SizedBox(height: 16),
          formDataAsync.when(
            data: (formData) {
              final sources = [
                null,
                ...formData.leadSources.map((e) => e.value),
              ];
              return Column(
                mainAxisSize: MainAxisSize.min,
                children: sources.map((source) {
                  final label = source == null
                      ? l10n.all
                      : formData.leadSources
                            .firstWhere((e) => e.value == source)
                            .label;
                  final isSelected = selectedSource == source;
                  return ListTile(
                    title: Text(label),
                    trailing: isSelected
                        ? Icon(Icons.check, color: colorScheme.primary)
                        : null,
                    onTap: () => onSelect(source),
                  );
                }).toList(),
              );
            },
            loading: () => const Center(
              child: Padding(
                padding: EdgeInsets.all(16.0),
                child: CircularProgressIndicator(),
              ),
            ),
            error: (error, _) => Center(child: Text(error.toString())),
          ),
        ],
      ),
    );
  }
}

class IndustryFilterSheet extends ConsumerWidget {
  const IndustryFilterSheet({
    super.key,
    this.selectedIndustry,
    required this.onSelect,
  });

  final String? selectedIndustry;
  final ValueChanged<String?> onSelect;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final colorScheme = theme.colorScheme;

    final formDataAsync = ref.watch(leadFormDataProvider);

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            l10n.filterByIndustry,
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.onSurface,
            ),
          ),
          const SizedBox(height: 16),
          formDataAsync.when(
            data: (formData) {
              final industries = [null, ...formData.industries];
              return Column(
                mainAxisSize: MainAxisSize.min,
                children: industries.map((industry) {
                  final label = industry ?? l10n.all;
                  final isSelected = selectedIndustry == industry;
                  return ListTile(
                    title: Text(label),
                    trailing: isSelected
                        ? Icon(Icons.check, color: colorScheme.primary)
                        : null,
                    onTap: () => onSelect(industry),
                  );
                }).toList(),
              );
            },
            loading: () => const Center(
              child: Padding(
                padding: EdgeInsets.all(16.0),
                child: CircularProgressIndicator(),
              ),
            ),
            error: (error, _) => Center(child: Text(error.toString())),
          ),
        ],
      ),
    );
  }
}

class ProvinceFilterSheet extends StatefulWidget {
  const ProvinceFilterSheet({
    super.key,
    this.selectedProvince,
    required this.onSelect,
  });

  final String? selectedProvince;
  final ValueChanged<String?> onSelect;

  @override
  State<ProvinceFilterSheet> createState() => _ProvinceFilterSheetState();
}

class _ProvinceFilterSheetState extends State<ProvinceFilterSheet> {
  String _searchQuery = '';

  @override
  Widget build(BuildContext context) {
    return Consumer(
      builder: (context, ref, child) {
        final theme = Theme.of(context);
        final l10n = AppLocalizations.of(context)!;
        final colorScheme = theme.colorScheme;
        final formDataAsync = ref.watch(leadFormDataProvider);

        return Container(
          padding: const EdgeInsets.all(16),
          constraints: BoxConstraints(
            maxHeight: MediaQuery.of(context).size.height * 0.7,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                l10n.filterByProvince,
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: colorScheme.onSurface,
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                decoration: InputDecoration(
                  hintText: l10n.search,
                  prefixIcon: const Icon(Icons.search),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  contentPadding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 8,
                  ),
                ),
                onChanged: (value) {
                  setState(() {
                    _searchQuery = value;
                  });
                },
              ),
              const SizedBox(height: 16),
              Flexible(
                child: formDataAsync.when(
                  data: (formData) {
                    final filteredProvinces = formData.provinces
                        .where(
                          (p) => p.toLowerCase().contains(
                            _searchQuery.toLowerCase(),
                          ),
                        )
                        .toList();

                    final provinces = [null, ...filteredProvinces];

                    return ListView.builder(
                      shrinkWrap: true,
                      itemCount: provinces.length,
                      itemBuilder: (context, index) {
                        final province = provinces[index];
                        final label = province ?? l10n.all;
                        final isSelected = widget.selectedProvince == province;
                        return ListTile(
                          title: Text(label),
                          trailing: isSelected
                              ? Icon(Icons.check, color: colorScheme.primary)
                              : null,
                          onTap: () => widget.onSelect(province),
                        );
                      },
                    );
                  },
                  loading: () => const Center(
                    child: Padding(
                      padding: EdgeInsets.all(16.0),
                      child: CircularProgressIndicator(),
                    ),
                  ),
                  error: (error, _) => Center(child: Text(error.toString())),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
