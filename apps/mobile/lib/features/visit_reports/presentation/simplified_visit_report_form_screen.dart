import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';

import '../application/visit_report_provider.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../dashboard/presentation/widgets/modern_tab_switch.dart';
import 'widgets/searchable_dropdown.dart';

/// Simplified Visit Report Form untuk sales rep
/// Support tabs: Account, Deal, Lead (seperti web version)
/// Redesigned untuk native mobile dan konsisten dengan desain halaman lain
class SimplifiedVisitReportFormScreen extends ConsumerStatefulWidget {
  const SimplifiedVisitReportFormScreen({super.key});

  @override
  ConsumerState<SimplifiedVisitReportFormScreen> createState() =>
      _SimplifiedVisitReportFormScreenState();
}

class _SimplifiedVisitReportFormScreenState
    extends ConsumerState<SimplifiedVisitReportFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _purposeController = TextEditingController();
  final _notesController = TextEditingController();

  String? _selectedAccountId;
  String? _selectedContactId;
  String? _selectedDealId;
  String? _selectedLeadId;
  int _selectedTabIndex = 0; // 0: account, 1: deal, 2: lead
  DateTime _selectedDate = DateTime.now();
  TimeOfDay _selectedTime = TimeOfDay.now();

  @override
  void dispose() {
    _purposeController.dispose();
    _notesController.dispose();
    super.dispose();
  }

  Future<void> _selectDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime.now().subtract(const Duration(days: 365)),
      lastDate: DateTime.now().add(const Duration(days: 365)),
    );
    if (picked != null) {
      setState(() {
        _selectedDate = picked;
      });
    }
  }

  Future<void> _selectTime() async {
    final picked = await showTimePicker(
      context: context,
      initialTime: _selectedTime,
    );
    if (picked != null) {
      setState(() {
        _selectedTime = picked;
      });
    }
  }

  String _formatVisitDate() {
    final dateStr = DateFormat('yyyy-MM-dd').format(_selectedDate);
    final timeStr =
        '${_selectedTime.hour.toString().padLeft(2, '0')}:${_selectedTime.minute.toString().padLeft(2, '0')}';
    return '$dateStr $timeStr';
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;
    final l10n = AppLocalizations.of(context)!;

    // Business rule validation based on active tab
    if (_selectedTabIndex == 2 && _selectedLeadId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(l10n.pleaseSelectLead),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }
    if (_selectedTabIndex == 1 && _selectedDealId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(l10n.pleaseSelectDeal),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }
    if (_selectedTabIndex == 0 && _selectedAccountId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(l10n.pleaseSelectAccount),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }

    // Purpose is required
    if (_purposeController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(l10n.purposeRequired),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }

    final formNotifier = ref.read(visitReportFormProvider.notifier);
    final visitReport = await formNotifier.createVisitReport(
      accountId: _selectedAccountId,
      contactId: _selectedContactId,
      dealId: _selectedDealId,
      leadId: _selectedLeadId,
      visitDate: _formatVisitDate(),
      purpose: _purposeController.text.trim(),
      notes: _notesController.text.trim().isNotEmpty
          ? _notesController.text.trim()
          : null,
    );

    if (mounted) {
      if (visitReport != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.visitReportCreatedSuccessfully),
            backgroundColor: Colors.green,
          ),
        );
        Navigator.pop(context, visitReport);
      } else {
        final error = ref.read(visitReportFormProvider).errorMessage;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(error ?? l10n.failedToCreateVisitReport),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  void _handleTabChange(int index) {
    setState(() {
      _selectedTabIndex = index;
      // Clear all selections when switching tabs
      _selectedAccountId = null;
      _selectedContactId = null;
      _selectedDealId = null;
      _selectedLeadId = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final formState = ref.watch(visitReportFormProvider);
    final formDataAsync = ref.watch(visitReportFormDataProvider);
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(title: Text(l10n.createVisitReport), elevation: 0),
      body: SafeArea(
        child: formDataAsync.when(
          data: (formData) {
            final accounts =
                (formData['accounts'] as List<dynamic>?)
                    ?.map((e) => Map<String, dynamic>.from(e))
                    .toList() ??
                [];
            final contacts =
                (formData['contacts'] as List<dynamic>?)
                    ?.map((e) => Map<String, dynamic>.from(e))
                    .toList() ??
                [];
            final deals =
                (formData['deals'] as List<dynamic>?)
                    ?.map((e) => Map<String, dynamic>.from(e))
                    .toList() ??
                [];
            final leads =
                (formData['leads'] as List<dynamic>?)
                    ?.map((e) => Map<String, dynamic>.from(e))
                    .toList() ??
                [];

            // Filter contacts by selected account
            final List<Map<String, dynamic>> filteredContacts =
                _selectedAccountId != null
                ? contacts.where((c) {
                    final accountId = c['account_id'] as String?;
                    return accountId == _selectedAccountId;
                  }).toList()
                : <Map<String, dynamic>>[];

            // Filter deals by status (only open deals)
            final openDeals = deals
                .where((d) => (d['status'] as String?) == 'open')
                .toList();

            return Form(
              key: _formKey,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  // Visit Type Selection (Modern Tab Switch)
                  ModernTabSwitch(
                    tabs: [l10n.accounts, l10n.deal, l10n.lead],
                    selectedIndex: _selectedTabIndex,
                    onChanged: _handleTabChange,
                  ),

                  const SizedBox(height: 24),

                  // Tab Content
                  _buildTabContent(
                    context,
                    _selectedTabIndex,
                    accounts,
                    filteredContacts,
                    openDeals,
                    leads,
                    theme,
                    colorScheme,
                    l10n,
                  ),

                  const SizedBox(height: 20),

                  // Visit Date & Time
                  TextFormField(
                    readOnly: true,
                    onTap: _selectDate,
                    decoration: InputDecoration(
                      labelText: '${l10n.visitDate} *',
                      hintText: l10n.visitDate,
                      prefixIcon: const Icon(Icons.calendar_today),
                      suffixIcon: const Icon(Icons.arrow_drop_down),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      filled: true,
                      fillColor: colorScheme.surfaceContainerHighest,
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 16,
                      ),
                    ),
                    controller: TextEditingController(
                      text: DateFormat('dd/MM/yyyy').format(_selectedDate),
                    ),
                  ),

                  const SizedBox(height: 20),

                  // Visit Time
                  TextFormField(
                    readOnly: true,
                    onTap: _selectTime,
                    decoration: InputDecoration(
                      labelText: 'Time *',
                      hintText: 'Select time',
                      prefixIcon: const Icon(Icons.access_time),
                      suffixIcon: const Icon(Icons.arrow_drop_down),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      filled: true,
                      fillColor: colorScheme.surfaceContainerHighest,
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 16,
                      ),
                    ),
                    controller: TextEditingController(
                      text: _selectedTime.format(context),
                    ),
                  ),

                  const SizedBox(height: 20),

                  // Purpose (Required)
                  TextFormField(
                    controller: _purposeController,
                    decoration: InputDecoration(
                      labelText: '${l10n.purpose} *',
                      hintText: 'e.g., Product demo, Follow-up, Meeting',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      filled: true,
                      fillColor: colorScheme.surfaceContainerHighest,
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 16,
                      ),
                    ),
                    maxLines: 2,
                    textInputAction: TextInputAction.next,
                    validator: (value) {
                      if (value == null || value.trim().isEmpty) {
                        return '${l10n.purpose} ${l10n.required.toLowerCase()}';
                      }
                      if (value.trim().length < 3) {
                        return 'Purpose must be at least 3 characters';
                      }
                      return null;
                    },
                  ),

                  const SizedBox(height: 20),

                  // Notes (Optional)
                  TextFormField(
                    controller: _notesController,
                    decoration: InputDecoration(
                      labelText: l10n.notes,
                      hintText: 'Enter additional notes...',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      filled: true,
                      fillColor: colorScheme.surfaceContainerHighest,
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 16,
                      ),
                    ),
                    maxLines: 4,
                    textInputAction: TextInputAction.done,
                  ),

                  const SizedBox(height: 32),

                  // Submit Button
                  FilledButton(
                    onPressed: formState.isSubmitting ? null : _handleSubmit,
                    style: FilledButton.styleFrom(
                      minimumSize: const Size(double.infinity, 56),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      padding: const EdgeInsets.symmetric(vertical: 16),
                    ),
                    child: formState.isSubmitting
                        ? SizedBox(
                            width: 24,
                            height: 24,
                            child: CircularProgressIndicator(
                              strokeWidth: 2.5,
                              valueColor: AlwaysStoppedAnimation<Color>(
                                Colors.white,
                              ),
                            ),
                          )
                        : Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              const Icon(Icons.check_circle_outline, size: 24),
                              const SizedBox(width: 8),
                              Text(
                                l10n.createVisitReport,
                                style: theme.textTheme.titleLarge?.copyWith(
                                  fontWeight: FontWeight.bold,
                                  color: Colors.white,
                                ),
                              ),
                            ],
                          ),
                  ),
                ],
              ),
            );
          },
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (error, stack) => Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text('Error: ${error.toString()}'),
                const SizedBox(height: 16),
                ElevatedButton(
                  onPressed: () {
                    ref.invalidate(visitReportFormDataProvider);
                  },
                  child: const Text('Retry'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildTabContent(
    BuildContext context,
    int tabIndex,
    List<Map<String, dynamic>> accounts,
    List<Map<String, dynamic>> filteredContacts,
    List<Map<String, dynamic>> openDeals,
    List<Map<String, dynamic>> leads,
    ThemeData theme,
    ColorScheme colorScheme,
    AppLocalizations l10n,
  ) {
    switch (tabIndex) {
      case 0:
        return _buildAccountTab(
          context,
          accounts,
          filteredContacts,
          theme,
          colorScheme,
          l10n,
        );
      case 1:
        return _buildDealTab(
          context,
          openDeals,
          accounts,
          filteredContacts,
          theme,
          colorScheme,
          l10n,
        );
      case 2:
        return _buildLeadTab(context, leads, theme, colorScheme, l10n);
      default:
        return const SizedBox.shrink();
    }
  }

  Widget _buildAccountTab(
    BuildContext context,
    List<Map<String, dynamic>> accounts,
    List<Map<String, dynamic>> contacts,
    ThemeData theme,
    ColorScheme colorScheme,
    AppLocalizations l10n,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Account Selection (Required)
        SearchableDropdown<String>(
          labelText: '${l10n.selectAccount} *',
          hintText: 'Select account...',
          icon: Icons.business_outlined,
          searchHint: 'Search accounts...',
          items: accounts
              .where(
                (a) =>
                    (a['status'] as String?) == 'active' &&
                    a['id'] != null &&
                    a['name'] != null &&
                    (a['name'] as String).trim().isNotEmpty,
              )
              .map((account) => account['id'] as String)
              .toList(),
          displayText: (id) {
            final account = accounts.firstWhere(
              (a) => a['id'] == id,
              orElse: () => <String, dynamic>{},
            );
            return account['name'] as String? ?? '';
          },
          selectedValue: _selectedAccountId,
          onChanged: (value) {
            setState(() {
              _selectedAccountId = value;
              _selectedContactId = null; // Reset contact when account changes
            });
          },
          validator: (value) {
            if (_selectedTabIndex == 0 && (value == null || value.isEmpty)) {
              return '${l10n.selectAccount} ${l10n.required.toLowerCase()}';
            }
            return null;
          },
        ),

        // Contact Selection (Optional, only if account selected)
        if (_selectedAccountId != null) ...[
          const SizedBox(height: 20),
          SearchableDropdown<String?>(
            labelText: l10n.selectContactOptional,
            hintText: 'Select contact (optional)...',
            icon: Icons.person_outline,
            searchHint: 'Search contacts...',
            allowNone: true,
            noneText: l10n.none,
            items: contacts
                .where(
                  (c) =>
                      c['id'] != null &&
                      c['name'] != null &&
                      (c['name'] as String).trim().isNotEmpty,
                )
                .map((contact) => contact['id'] as String)
                .toList(),
            displayText: (id) {
              if (id == null) return l10n.none;
              final contact = contacts.firstWhere(
                (c) => c['id'] == id,
                orElse: () => <String, dynamic>{},
              );
              return contact['name'] as String? ?? '';
            },
            selectedValue: _selectedContactId,
            onChanged: (value) {
              setState(() {
                _selectedContactId = value;
              });
            },
          ),
        ],
      ],
    );
  }

  Widget _buildDealTab(
    BuildContext context,
    List<Map<String, dynamic>> deals,
    List<Map<String, dynamic>> accounts,
    List<Map<String, dynamic>> contacts,
    ThemeData theme,
    ColorScheme colorScheme,
    AppLocalizations l10n,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Deal Selection (Required)
        SearchableDropdown<String>(
          labelText: '${l10n.deal} *',
          hintText: 'Select deal...',
          icon: Icons.handshake_outlined,
          searchHint: 'Search deals...',
          items: deals.map((deal) => deal['id'] as String).toList(),
          displayText: (id) {
            final deal = deals.firstWhere(
              (d) => d['id'] == id,
              orElse: () => <String, dynamic>{},
            );
            final accountId = deal['account_id'] as String?;
            final account = accounts.firstWhere(
              (a) => a['id'] == accountId,
              orElse: () => <String, dynamic>{},
            );
            final accountName = account['name'] as String? ?? '';
            final dealTitle = deal['title'] as String? ?? '';
            return accountName.isNotEmpty
                ? '$dealTitle ($accountName)'
                : dealTitle;
          },
          selectedValue: _selectedDealId,
          onChanged: (value) {
            setState(() {
              _selectedDealId = value;
              // Auto-set account_id from deal (business rule: Deal requires Account)
              if (value != null) {
                final selectedDeal = deals.firstWhere(
                  (d) => d['id'] == value,
                  orElse: () => <String, dynamic>{},
                );
                final accountId = selectedDeal['account_id'] as String?;
                if (accountId != null) {
                  _selectedAccountId = accountId;
                  _selectedContactId =
                      null; // Reset contact when account changes
                }
              } else {
                _selectedAccountId = null;
              }
            });
          },
          validator: (value) {
            if (_selectedTabIndex == 1 && (value == null || value.isEmpty)) {
              return '${l10n.deal} ${l10n.required.toLowerCase()}';
            }
            return null;
          },
        ),

        // Contact Selection (Optional, only if deal selected and account is set)
        if (_selectedDealId != null && _selectedAccountId != null) ...[
          const SizedBox(height: 20),
          SearchableDropdown<String?>(
            labelText: l10n.selectContactOptional,
            hintText: 'Select contact (optional)...',
            icon: Icons.person_outline,
            searchHint: 'Search contacts...',
            allowNone: true,
            noneText: l10n.none,
            items: contacts
                .where((c) => c['account_id'] == _selectedAccountId)
                .map((contact) => contact['id'] as String)
                .toList(),
            displayText: (id) {
              if (id == null) return l10n.none;
              final contact = contacts.firstWhere(
                (c) => c['id'] == id,
                orElse: () => <String, dynamic>{},
              );
              return contact['name'] as String? ?? '';
            },
            selectedValue: _selectedContactId,
            onChanged: (value) {
              setState(() {
                _selectedContactId = value;
              });
            },
          ),
        ],
      ],
    );
  }

  Widget _buildLeadTab(
    BuildContext context,
    List<Map<String, dynamic>> leads,
    ThemeData theme,
    ColorScheme colorScheme,
    AppLocalizations l10n,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Lead Selection (Required)
        SearchableDropdown<String>(
          labelText: '${l10n.lead} *',
          hintText: leads.isEmpty ? l10n.noLeadsAvailable : 'Select lead...',
          icon: Icons.person_outline,
          searchHint: 'Search leads...',
          items: leads.map((lead) => lead['id'] as String).toList(),
          displayText: (id) {
            final lead = leads.firstWhere(
              (l) => l['id'] == id,
              orElse: () => <String, dynamic>{},
            );
            final firstName = lead['first_name'] as String? ?? '';
            final lastName = lead['last_name'] as String?;
            final companyName = lead['company_name'] as String?;
            final displayName = '$firstName ${lastName ?? ''}'.trim();
            return companyName != null
                ? '$displayName ($companyName)'
                : displayName;
          },
          selectedValue: _selectedLeadId,
          onChanged: leads.isEmpty
              ? null
              : (value) {
                  setState(() {
                    _selectedLeadId = value;
                  });
                },
          validator: (value) {
            if (_selectedTabIndex == 2 && (value == null || value.isEmpty)) {
              return '${l10n.lead} ${l10n.required.toLowerCase()}';
            }
            return null;
          },
        ),
      ],
    );
  }
}
