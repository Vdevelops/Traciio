import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../visit_reports/presentation/widgets/searchable_dropdown.dart';
import '../application/lead_provider.dart';
import '../data/models/lead_model.dart';

class LeadFormScreen extends ConsumerStatefulWidget {
  const LeadFormScreen({super.key, this.leadId});

  final String? leadId;

  @override
  ConsumerState<LeadFormScreen> createState() => _LeadFormScreenState();
}

class _LeadFormScreenState extends ConsumerState<LeadFormScreen> {
  final _formKey = GlobalKey<FormState>();

  late TextEditingController _firstNameController;
  late TextEditingController _lastNameController;
  late TextEditingController _companyController;
  late TextEditingController _titleController;
  late TextEditingController _emailController;
  late TextEditingController _phoneController;
  late TextEditingController _addressController;
  late TextEditingController _notesController;

  String? _selectedStatus;
  String? _selectedSource;
  String? _selectedIndustry;
  String? _selectedProvince;

  bool _isEditing = false;
  bool _isLoadingData = false;

  @override
  void initState() {
    super.initState();
    _firstNameController = TextEditingController();
    _lastNameController = TextEditingController();
    _companyController = TextEditingController();
    _titleController = TextEditingController();
    _emailController = TextEditingController();
    _phoneController = TextEditingController();
    _addressController = TextEditingController();
    _notesController = TextEditingController();

    _isEditing = widget.leadId != null;

    if (_isEditing) {
      _loadLeadData();
    }
  }

  Future<void> _loadLeadData() async {
    setState(() {
      _isLoadingData = true;
    });

    try {
      final lead = await ref.read(leadDetailProvider(widget.leadId!).future);
      _firstNameController.text = lead.firstName;
      _lastNameController.text = lead.lastName ?? '';
      _companyController.text = lead.company;
      _titleController.text = lead.title ?? '';
      _emailController.text = lead.email ?? '';
      _phoneController.text = lead.phone ?? '';
      _addressController.text = lead.address ?? '';
      _notesController.text = lead.notes ?? '';
      _selectedStatus = lead.leadStatusId; // Use ID
      _selectedSource = lead.source;
      _selectedIndustry = lead.industry;
      _selectedProvince = lead.province;
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Failed to load lead data: $e')));
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoadingData = false;
        });
      }
    }
  }

  @override
  void dispose() {
    _firstNameController.dispose();
    _lastNameController.dispose();
    _companyController.dispose();
    _titleController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _addressController.dispose();
    _notesController.dispose();
    super.dispose();
  }

  void _saveLead() async {
    if (!_formKey.currentState!.validate()) return;

    final l10n = AppLocalizations.of(context)!;

    final leadData = {
      'first_name': _firstNameController.text.trim(),
      'last_name': _lastNameController.text.trim().isNotEmpty
          ? _lastNameController.text.trim()
          : null,
      'company_name': _companyController.text.trim().isNotEmpty
          ? _companyController.text.trim()
          : null,
      'job_title': _titleController.text.trim().isNotEmpty
          ? _titleController.text.trim()
          : null,
      'email': _emailController.text.trim().isNotEmpty
          ? _emailController.text.trim()
          : null,
      'phone': _phoneController.text.trim().isNotEmpty
          ? _phoneController.text.trim()
          : null,
      'address': _addressController.text.trim().isNotEmpty
          ? _addressController.text.trim()
          : null,
      'notes': _notesController.text.trim().isNotEmpty
          ? _notesController.text.trim()
          : null,
      'lead_status_id': _selectedStatus,
      'lead_source': _selectedSource,
      'industry': _selectedIndustry,
      'province': _selectedProvince,
    };

    final success = _isEditing
        ? await ref
              .read(leadFormProvider.notifier)
              .updateLead(widget.leadId!, leadData)
        : await ref.read(leadFormProvider.notifier).createLead(leadData);

    if (success != null && mounted) {
      final messenger = ScaffoldMessenger.of(context);
      messenger.showSnackBar(
        SnackBar(
          content: Text(
            _isEditing
                ? l10n.leadUpdatedSuccessfully
                : l10n.leadCreatedSuccessfully,
          ),
        ),
      );
      Navigator.pop(context);
    } else if (mounted) {
      final messenger = ScaffoldMessenger.of(context);
      messenger.showSnackBar(
        SnackBar(
          content: Text(
            _isEditing ? l10n.failedToUpdateLead : l10n.failedToCreateLead,
          ),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final formState = ref.watch(leadFormProvider);
    final formDataAsync = ref.watch(leadFormDataProvider);

    return Scaffold(
      appBar: AppBar(title: Text(_isEditing ? l10n.editLead : l10n.createLead)),
      body: _isLoadingData
          ? const Center(child: CircularProgressIndicator())
          : formDataAsync.when(
              data: (formData) => Form(
                key: _formKey,
                child: SingleChildScrollView(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      Row(
                        children: [
                          Expanded(
                            child: TextFormField(
                              controller: _firstNameController,
                              decoration: InputDecoration(
                                labelText: '${l10n.firstName} *',
                              ),
                              validator: (value) {
                                if (value == null || value.trim().isEmpty) {
                                  return '${l10n.firstName} ${l10n.required}';
                                }
                                return null;
                              },
                            ),
                          ),
                          const SizedBox(width: 16),
                          Expanded(
                            child: TextFormField(
                              controller: _lastNameController,
                              decoration: InputDecoration(
                                labelText: l10n.lastName,
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      TextFormField(
                        controller: _companyController,
                        decoration: InputDecoration(labelText: l10n.company),
                      ),
                      const SizedBox(height: 16),
                      TextFormField(
                        controller: _titleController,
                        decoration: InputDecoration(labelText: l10n.title),
                      ),
                      const SizedBox(height: 16),
                      Row(
                        children: [
                          Expanded(
                            child: TextFormField(
                              controller: _emailController,
                              decoration: InputDecoration(
                                labelText: '${l10n.email} *',
                              ),
                              keyboardType: TextInputType.emailAddress,
                              validator: (value) {
                                if (value == null || value.trim().isEmpty) {
                                  return '${l10n.email} ${l10n.required}';
                                }
                                final emailRegex = RegExp(
                                  r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$',
                                );
                                if (!emailRegex.hasMatch(value.trim())) {
                                  return 'Invalid email address';
                                }
                                return null;
                              },
                            ),
                          ),
                          const SizedBox(width: 16),
                          Expanded(
                            child: TextFormField(
                              controller: _phoneController,
                              decoration: InputDecoration(
                                labelText: l10n.phone,
                              ),
                              keyboardType: TextInputType.phone,
                              inputFormatters: [
                                FilteringTextInputFormatter.digitsOnly,
                              ],
                              validator: (value) {
                                if (value != null && value.trim().isNotEmpty) {
                                  if (value.length < 8) {
                                    return 'Minimum 8 digits';
                                  }
                                }
                                return null;
                              },
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      DropdownButtonFormField<String>(
                        initialValue: _selectedStatus,
                        decoration: InputDecoration(
                          labelText: '${l10n.status} *',
                        ),
                        items:
                            _getUniqueOptions(
                                  formData.leadStatuses,
                                  _selectedStatus,
                                )
                                .map(
                                  (e) => DropdownMenuItem(
                                    value: e.id,
                                    child: Text(e.label),
                                  ),
                                )
                                .toList(),
                        onChanged: (value) =>
                            setState(() => _selectedStatus = value),
                        validator: (value) {
                          if (_isEditing && value == null) {
                            return l10n.isRequired;
                          }
                          return null;
                        },
                      ),
                      const SizedBox(height: 16),
                      DropdownButtonFormField<String>(
                        initialValue: _selectedSource,
                        decoration: InputDecoration(
                          labelText: '${l10n.source} *',
                        ),
                        items:
                            _getUniqueOptions(
                                  formData.leadSources,
                                  _selectedSource,
                                )
                                .map(
                                  (e) => DropdownMenuItem(
                                    value: e.value,
                                    child: Text(e.label),
                                  ),
                                )
                                .toList(),
                        onChanged: (value) =>
                            setState(() => _selectedSource = value),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return '${l10n.source} ${l10n.required}';
                          }
                          return null;
                        },
                      ),
                      const SizedBox(height: 16),
                      DropdownButtonFormField<String>(
                        initialValue: _selectedIndustry,
                        decoration: InputDecoration(labelText: l10n.industry),
                        items:
                            _getUniqueStringOptions(
                                  formData.industries,
                                  _selectedIndustry,
                                )
                                .map(
                                  (e) => DropdownMenuItem(
                                    value: e,
                                    child: Text(e),
                                  ),
                                )
                                .toList(),
                        onChanged: (value) =>
                            setState(() => _selectedIndustry = value),
                      ),
                      const SizedBox(height: 16),
                      SearchableDropdown<String>(
                        labelText: l10n.province,
                        items: formData.provinces,
                        selectedValue: _selectedProvince,
                        displayText: (item) => item,
                        onChanged: (value) =>
                            setState(() => _selectedProvince = value),
                      ),
                      const SizedBox(height: 16),
                      TextFormField(
                        controller: _addressController,
                        decoration: InputDecoration(labelText: l10n.address),
                        maxLines: 2,
                      ),
                      const SizedBox(height: 16),
                      TextFormField(
                        controller: _notesController,
                        decoration: InputDecoration(labelText: l10n.notes),
                        maxLines: 3,
                      ),
                      const SizedBox(height: 24),
                      FilledButton(
                        onPressed: formState.isLoading ? null : _saveLead,
                        child: formState.isLoading
                            ? const SizedBox(
                                height: 20,
                                width: 20,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                ),
                              )
                            : Text(l10n.save),
                      ),
                    ],
                  ),
                ),
              ),
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (error, _) => Center(child: Text(error.toString())),
            ),
    );
  }

  List<FormOption> _getUniqueOptions(
    List<FormOption> options,
    String? currentValue,
  ) {
    final Map<String, FormOption> uniqueMap = {};
    for (var option in options) {
      if (option.value.isNotEmpty) {
        uniqueMap[option.value] = option;
      }
    }

    // Ensure current value is in the list to avoid assertion error
    if (currentValue != null &&
        currentValue.isNotEmpty &&
        !uniqueMap.containsKey(currentValue)) {
      uniqueMap[currentValue] = FormOption(
        label: currentValue,
        value: currentValue,
      );
    }

    return uniqueMap.values.toList();
  }

  List<String> _getUniqueStringOptions(
    List<String> options,
    String? currentValue,
  ) {
    final Set<String> uniqueSet = options.toSet();
    if (currentValue != null && currentValue.isNotEmpty) {
      uniqueSet.add(currentValue);
    }
    return uniqueSet.toList();
  }
}
