import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../application/account_provider.dart';
import '../data/models/account.dart';
import '../../../core/l10n/app_localizations.dart';

class AccountFormScreen extends ConsumerStatefulWidget {
  const AccountFormScreen({super.key, this.account});

  final Account? account;

  @override
  ConsumerState<AccountFormScreen> createState() => _AccountFormScreenState();
}

class _AccountFormScreenState extends ConsumerState<AccountFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _addressController = TextEditingController();
  final _cityController = TextEditingController();
  final _provinceController = TextEditingController();
  final _phoneController = TextEditingController();
  final _emailController = TextEditingController();

  String? _selectedCategoryId;
  String _selectedStatus = 'active';

  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    // Populate fields if editing
    if (widget.account != null) {
      final account = widget.account!;
      _nameController.text = account.name;
      _addressController.text = account.address ?? '';
      _cityController.text = account.city ?? '';
      _provinceController.text = account.province ?? '';
      _phoneController.text = account.phone ?? '';
      _emailController.text = account.email ?? '';
      _selectedCategoryId = account.categoryId;
      _selectedStatus = account.status;
    }
    // Add listeners to update button state
    _nameController.addListener(_updateFormState);
    _addressController.addListener(_updateFormState);
    _cityController.addListener(_updateFormState);
    _provinceController.addListener(_updateFormState);
    _phoneController.addListener(_updateFormState);
    _emailController.addListener(_updateFormState);
  }

  void _updateFormState() {
    setState(() {}); // Rebuild to update button state
  }

  @override
  void dispose() {
    _nameController.removeListener(_updateFormState);
    _addressController.removeListener(_updateFormState);
    _cityController.removeListener(_updateFormState);
    _provinceController.removeListener(_updateFormState);
    _phoneController.removeListener(_updateFormState);
    _emailController.removeListener(_updateFormState);
    _nameController.dispose();
    _addressController.dispose();
    _cityController.dispose();
    _provinceController.dispose();
    _phoneController.dispose();
    _emailController.dispose();
    super.dispose();
  }

  bool _isFormValid() {
    return _nameController.text.trim().isNotEmpty &&
        _selectedCategoryId != null &&
        _addressController.text.trim().isNotEmpty &&
        _cityController.text.trim().isNotEmpty &&
        _provinceController.text.trim().isNotEmpty &&
        _phoneController.text.trim().isNotEmpty &&
        _emailController.text.trim().isNotEmpty &&
        _formKey.currentState?.validate() == true;
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;
    if (_selectedCategoryId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(AppLocalizations.of(context)!.selectCategory),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }

    setState(() => _isLoading = true);

    final l10n = AppLocalizations.of(context)!;
    Account? result;

    if (widget.account != null) {
      // Update existing account
      result = await ref
          .read(accountListProvider.notifier)
          .updateAccount(
            id: widget.account!.id,
            name: _nameController.text.trim(),
            categoryId: _selectedCategoryId,
            address: _addressController.text.trim(),
            city: _cityController.text.trim(),
            province: _provinceController.text.trim(),
            phone: _phoneController.text.trim(),
            email: _emailController.text.trim(),
            status: _selectedStatus,
          );
    } else {
      // Create new account
      result = await ref
          .read(accountListProvider.notifier)
          .createAccount(
            name: _nameController.text.trim(),
            categoryId: _selectedCategoryId!,
            address: _addressController.text.trim(),
            city: _cityController.text.trim(),
            province: _provinceController.text.trim(),
            phone: _phoneController.text.trim(),
            email: _emailController.text.trim(),
            status: _selectedStatus,
          );
    }

    if (mounted) {
      setState(() => _isLoading = false);
      if (result != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              widget.account != null
                  ? l10n.accountUpdated
                  : l10n.accountCreatedSuccessfully,
            ),
            backgroundColor: Colors.green,
          ),
        );
        Navigator.pop(context, result);
      } else {
        final error = ref.read(accountListProvider).errorMessage;
        // Show error in a dialog for better visibility
        if (mounted) {
          showDialog(
            context: context,
            builder: (context) => AlertDialog(
              title: Text('Error'),
              content: Text(
                error ??
                    (widget.account != null
                        ? 'Failed to update account'
                        : 'Failed to create account'),
              ),
              actions: [
                TextButton(
                  onPressed: () => Navigator.pop(context),
                  child: Text('OK'),
                ),
              ],
            ),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final categoriesAsync = ref.watch(categoriesProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          widget.account != null ? l10n.editAccount : l10n.createAccount,
        ),
        elevation: 0,
      ),
      body: SafeArea(
        child: Form(
          key: _formKey,
          child: ListView(
            padding: const EdgeInsets.all(16),
            children: [
              // Name
              TextFormField(
                controller: _nameController,
                decoration: InputDecoration(
                  labelText: '${l10n.name} *',
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
                textInputAction: TextInputAction.next,
                validator: (value) {
                  if (value == null || value.trim().isEmpty) {
                    return '${l10n.name} ${l10n.required.toLowerCase()}';
                  }
                  if (value.trim().length < 3) {
                    return 'Name must be at least 3 characters';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              // Category
              categoriesAsync.when(
                data: (categories) => DropdownButtonFormField<String>(
                  initialValue: _selectedCategoryId,
                  decoration: InputDecoration(
                    labelText: '${l10n.category} *',
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
                  items: categories
                      .where((cat) => cat.name.isNotEmpty)
                      .map(
                        (category) => DropdownMenuItem(
                          value: category.id,
                          child: Text(category.name),
                        ),
                      )
                      .toList(),
                  onChanged: (value) {
                    setState(() => _selectedCategoryId = value);
                  },
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return '${l10n.category} ${l10n.required.toLowerCase()}';
                    }
                    return null;
                  },
                ),
                loading: () => const Padding(
                  padding: EdgeInsets.all(16),
                  child: Center(child: CircularProgressIndicator()),
                ),
                error: (error, stack) => Padding(
                  padding: const EdgeInsets.all(16),
                  child: Text(
                    'Error loading categories: $error',
                    style: TextStyle(color: colorScheme.error),
                  ),
                ),
              ),
              const SizedBox(height: 20),
              // Address
              TextFormField(
                controller: _addressController,
                decoration: InputDecoration(
                  labelText: '${l10n.address} *',
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
                maxLines: 3,
                textInputAction: TextInputAction.next,
                validator: (value) {
                  if (value == null || value.trim().isEmpty) {
                    return '${l10n.address} ${l10n.required.toLowerCase()}';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              // City
              TextFormField(
                controller: _cityController,
                decoration: InputDecoration(
                  labelText: '${l10n.city} *',
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
                textInputAction: TextInputAction.next,
                validator: (value) {
                  if (value == null || value.trim().isEmpty) {
                    return '${l10n.city} ${l10n.required.toLowerCase()}';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              // Province
              TextFormField(
                controller: _provinceController,
                decoration: InputDecoration(
                  labelText: '${l10n.province} *',
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
                textInputAction: TextInputAction.next,
                validator: (value) {
                  if (value == null || value.trim().isEmpty) {
                    return '${l10n.province} ${l10n.required.toLowerCase()}';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              // Phone
              TextFormField(
                controller: _phoneController,
                decoration: InputDecoration(
                  labelText: '${l10n.phone} *',
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
                keyboardType: TextInputType.number,
                textInputAction: TextInputAction.next,
                inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                validator: (value) {
                  if (value == null || value.trim().isEmpty) {
                    return '${l10n.phone} ${l10n.required.toLowerCase()}';
                  }
                  // Phone validation: at least 10 digits, max 15 digits
                  final digitsOnly = value.replaceAll(RegExp(r'[^\d]'), '');
                  if (digitsOnly.length < 10 || digitsOnly.length > 15) {
                    return 'Phone number must be between 10-15 digits';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              // Email
              TextFormField(
                controller: _emailController,
                decoration: InputDecoration(
                  labelText: '${l10n.email} *',
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
                keyboardType: TextInputType.emailAddress,
                textInputAction: TextInputAction.next,
                validator: (value) {
                  if (value == null || value.trim().isEmpty) {
                    return '${l10n.email} ${l10n.required.toLowerCase()}';
                  }
                  // Email validation
                  final emailRegex = RegExp(
                    r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
                  );
                  if (!emailRegex.hasMatch(value.trim())) {
                    return 'Invalid email format';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              // Status
              DropdownButtonFormField<String>(
                initialValue: _selectedStatus,
                decoration: InputDecoration(
                  labelText: l10n.status,
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
                items: [
                  DropdownMenuItem(value: 'active', child: Text(l10n.active)),
                  DropdownMenuItem(
                    value: 'inactive',
                    child: Text(l10n.inactive),
                  ),
                ],
                onChanged: (value) {
                  if (value != null) {
                    setState(() => _selectedStatus = value);
                  }
                },
              ),
              const SizedBox(height: 32),
              // Submit Button
              FilledButton(
                onPressed: (_isLoading || !_isFormValid())
                    ? null
                    : _handleSubmit,
                style: FilledButton.styleFrom(
                  minimumSize: const Size(double.infinity, 52),
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                child: _isLoading
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : Text(
                        l10n.save,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
              ),
              const SizedBox(height: 16),
            ],
          ),
        ),
      ),
    );
  }
}
