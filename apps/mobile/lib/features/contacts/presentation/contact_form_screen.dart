import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../application/contact_provider.dart';
import '../data/models/contact.dart';
import '../../accounts/application/account_provider.dart';
import '../../../core/l10n/app_localizations.dart';

class ContactFormScreen extends ConsumerStatefulWidget {
  const ContactFormScreen({super.key, this.defaultAccountId, this.contact});

  final String? defaultAccountId;
  final Contact? contact;

  @override
  ConsumerState<ContactFormScreen> createState() => _ContactFormScreenState();
}

class _ContactFormScreenState extends ConsumerState<ContactFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _phoneController = TextEditingController();
  final _emailController = TextEditingController();
  final _positionController = TextEditingController();
  final _notesController = TextEditingController();

  String? _selectedAccountId;
  String? _selectedRoleId;

  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    // Populate fields if editing
    if (widget.contact != null) {
      final contact = widget.contact!;
      _nameController.text = contact.name;
      _phoneController.text = contact.phone ?? '';
      _emailController.text = contact.email ?? '';
      _positionController.text = contact.position ?? '';
      _notesController.text = contact.notes ?? '';
      _selectedAccountId = contact.accountId;
      _selectedRoleId = contact.roleId;
    } else {
      _selectedAccountId = widget.defaultAccountId;
    }
    // Add listeners to update button state
    _nameController.addListener(_updateFormState);
    _phoneController.addListener(_updateFormState);
    _emailController.addListener(_updateFormState);
  }

  void _updateFormState() {
    setState(() {}); // Rebuild to update button state
  }

  @override
  void dispose() {
    _nameController.removeListener(_updateFormState);
    _phoneController.removeListener(_updateFormState);
    _emailController.removeListener(_updateFormState);
    _nameController.dispose();
    _phoneController.dispose();
    _emailController.dispose();
    _positionController.dispose();
    _notesController.dispose();
    super.dispose();
  }

  bool _isFormValid() {
    return _selectedAccountId != null &&
        _selectedRoleId != null &&
        _nameController.text.trim().isNotEmpty &&
        _phoneController.text.trim().isNotEmpty &&
        _emailController.text.trim().isNotEmpty &&
        _formKey.currentState?.validate() == true;
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;
    if (_selectedAccountId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(AppLocalizations.of(context)!.selectAccount),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }
    if (_selectedRoleId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(AppLocalizations.of(context)!.selectRole),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }

    setState(() => _isLoading = true);

    final l10n = AppLocalizations.of(context)!;
    Contact? result;

    if (widget.contact != null) {
      // Update existing contact
      result = await ref
          .read(contactListProvider.notifier)
          .updateContact(
            id: widget.contact!.id,
            accountId: _selectedAccountId,
            name: _nameController.text.trim(),
            roleId: _selectedRoleId,
            phone: _phoneController.text.trim(),
            email: _emailController.text.trim(),
            position: _positionController.text.trim().isEmpty
                ? null
                : _positionController.text.trim(),
            notes: _notesController.text.trim().isEmpty
                ? null
                : _notesController.text.trim(),
          );
    } else {
      // Create new contact
      result = await ref
          .read(contactListProvider.notifier)
          .createContact(
            accountId: _selectedAccountId!,
            name: _nameController.text.trim(),
            roleId: _selectedRoleId!,
            phone: _phoneController.text.trim(),
            email: _emailController.text.trim(),
            position: _positionController.text.trim().isEmpty
                ? null
                : _positionController.text.trim(),
            notes: _notesController.text.trim().isEmpty
                ? null
                : _notesController.text.trim(),
          );
    }

    if (mounted) {
      setState(() => _isLoading = false);
      if (result != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              widget.contact != null
                  ? l10n.contactUpdatedSuccessfully
                  : l10n.contactCreatedSuccessfully,
            ),
            backgroundColor: Colors.green,
          ),
        );
        Navigator.pop(context, result);
      } else {
        final error = ref.read(contactListProvider).errorMessage;
        // Show error in a dialog for better visibility
        if (mounted) {
          showDialog(
            context: context,
            builder: (context) => AlertDialog(
              title: Text('Error'),
              content: Text(
                error ??
                    (widget.contact != null
                        ? 'Failed to update contact'
                        : 'Failed to create contact'),
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
    final accountsAsync = ref.watch(accountListProvider);
    final rolesAsync = ref.watch(rolesProvider);

    // Load accounts if not loaded
    if (accountsAsync.accounts.isEmpty && !accountsAsync.isLoading) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        ref.read(accountListProvider.notifier).loadAccounts();
      });
    }

    return Scaffold(
      appBar: AppBar(
        title: Text(
          widget.contact != null ? l10n.editContact : l10n.createContact,
        ),
        elevation: 0,
      ),
      body: SafeArea(
        child: Form(
          key: _formKey,
          child: ListView(
            padding: const EdgeInsets.all(16),
            children: [
              // Account
              DropdownButtonFormField<String>(
                initialValue: _selectedAccountId,
                decoration: InputDecoration(
                  labelText: '${l10n.selectAccount} *',
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
                items: accountsAsync.accounts
                    .map(
                      (account) => DropdownMenuItem(
                        value: account.id,
                        child: Text(account.name),
                      ),
                    )
                    .toList(),
                onChanged: widget.contact != null
                    ? null
                    : (value) {
                        setState(() => _selectedAccountId = value);
                      },
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '${l10n.selectAccount} ${l10n.required.toLowerCase()}';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
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
              // Role
              rolesAsync.when(
                data: (roles) => DropdownButtonFormField<String>(
                  initialValue: _selectedRoleId,
                  decoration: InputDecoration(
                    labelText: '${l10n.role} *',
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
                  items: roles
                      .where((role) => role.name.isNotEmpty)
                      .map(
                        (role) => DropdownMenuItem(
                          value: role.id,
                          child: Text(role.name),
                        ),
                      )
                      .toList(),
                  onChanged: (value) {
                    setState(() => _selectedRoleId = value);
                  },
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return '${l10n.role} ${l10n.required.toLowerCase()}';
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
                    'Error loading roles: $error',
                    style: TextStyle(color: colorScheme.error),
                  ),
                ),
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
              // Position
              TextFormField(
                controller: _positionController,
                decoration: InputDecoration(
                  labelText: l10n.position,
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
              ),
              const SizedBox(height: 20),
              // Notes
              TextFormField(
                controller: _notesController,
                decoration: InputDecoration(
                  labelText: l10n.notes,
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
                textInputAction: TextInputAction.done,
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
