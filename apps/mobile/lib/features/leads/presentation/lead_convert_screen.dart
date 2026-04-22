import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';

import '../../../core/l10n/app_localizations.dart';
import '../../pipeline/application/pipeline_provider.dart';
import '../application/lead_provider.dart';
import '../data/models/lead_model.dart';

class LeadConvertScreen extends ConsumerStatefulWidget {
  final Lead lead;

  const LeadConvertScreen({super.key, required this.lead});

  @override
  ConsumerState<LeadConvertScreen> createState() => _LeadConvertScreenState();
}

class _LeadConvertScreenState extends ConsumerState<LeadConvertScreen> {
  final _formKey = GlobalKey<FormState>();

  late TextEditingController _titleController;
  late TextEditingController _valueController;
  late TextEditingController _probabilityController;
  late TextEditingController _descriptionController;

  String? _selectedStageId;
  DateTime? _expectedCloseDate;

  @override
  void initState() {
    super.initState();
    _titleController = TextEditingController(
      text:
          widget.lead.companyName ??
          '${widget.lead.firstName} ${widget.lead.lastName ?? ''}'.trim(),
    );
    _valueController = TextEditingController(text: '0');
    _probabilityController = TextEditingController(text: '0');
    _descriptionController = TextEditingController(text: widget.lead.notes);

    // Load stages if not already loaded
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(pipelineStagesProvider);
    });
  }

  @override
  void dispose() {
    _titleController.dispose();
    _valueController.dispose();
    _probabilityController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  Future<void> _selectDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate:
          _expectedCloseDate ?? DateTime.now().add(const Duration(days: 30)),
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 365 * 5)),
    );
    if (picked != null && picked != _expectedCloseDate) {
      setState(() {
        _expectedCloseDate = picked;
      });
    }
  }

  void _convertLead() async {
    if (!_formKey.currentState!.validate()) return;

    final l10n = AppLocalizations.of(context)!;

    final convertData = {
      'opportunity_title': _titleController.text.trim(),
      'value': int.tryParse(_valueController.text.trim()) ?? 0,
      'probability': int.tryParse(_probabilityController.text.trim()) ?? 0,
      'opportunity_description': _descriptionController.text.trim().isNotEmpty
          ? _descriptionController.text.trim()
          : null,
      'stage_id': _selectedStageId,
      if (_expectedCloseDate != null)
        'expected_close_date': _expectedCloseDate!.toIso8601String(),
    };

    final messenger = ScaffoldMessenger.of(context);

    final success = await ref
        .read(leadFormProvider.notifier)
        .convertLead(widget.lead.id, convertData);

    if (success && mounted) {
      messenger.showSnackBar(
        SnackBar(content: Text(l10n.leadConvertedSuccessfully)),
      );
      // Pop twice to go back to lead list (from detail and then from convert screen)
      Navigator.of(context).pop(); // Pop convert screen
      Navigator.of(context).pop(); // Pop lead detail screen
    } else if (mounted) {
      final state = ref.read(leadFormProvider);
      messenger.showSnackBar(
        SnackBar(content: Text(state.errorMessage ?? l10n.failedToConvertLead)),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final stagesAsync = ref.watch(pipelineStagesProvider);
    final formState = ref.watch(leadFormProvider);

    return Scaffold(
      appBar: AppBar(title: Text(l10n.convertLead)),
      body: formState.isLoading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              padding: const EdgeInsets.all(16.0),
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    TextFormField(
                      controller: _titleController,
                      decoration: InputDecoration(
                        labelText: l10n.opportunityTitle,
                        border: const OutlineInputBorder(),
                      ),
                      validator: (value) {
                        if (value == null || value.trim().isEmpty) {
                          return '${l10n.opportunityTitle} ${l10n.isRequired}';
                        }
                        return null;
                      },
                    ),
                    const SizedBox(height: 16),
                    stagesAsync.when(
                      data: (stages) {
                        if (_selectedStageId == null && stages.isNotEmpty) {
                          _selectedStageId = stages.first.id;
                          // Update probability based on stage
                          _probabilityController.text = stages.first.probability
                              .toString();
                        }
                        return DropdownButtonFormField<String>(
                          initialValue: _selectedStageId,
                          decoration: InputDecoration(
                            labelText: l10n.pipelineStage,
                            border: const OutlineInputBorder(),
                          ),
                          items: stages.map((stage) {
                            return DropdownMenuItem(
                              value: stage.id,
                              child: Text(stage.name),
                            );
                          }).toList(),
                          onChanged: (value) {
                            setState(() {
                              _selectedStageId = value;
                              final stage = stages.firstWhere(
                                (s) => s.id == value,
                              );
                              _probabilityController.text = stage.probability
                                  .toString();
                            });
                          },
                          validator: (value) {
                            if (value == null || value.isEmpty) {
                              return '${l10n.pipelineStage} ${l10n.isRequired}';
                            }
                            return null;
                          },
                        );
                      },
                      loading: () =>
                          const Center(child: LinearProgressIndicator()),
                      error: (err, stack) => Text('Error loading stages: $err'),
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Expanded(
                          child: TextFormField(
                            controller: _valueController,
                            decoration: InputDecoration(
                              labelText: l10n.dealValue,
                              border: const OutlineInputBorder(),
                              prefixText: 'Rp ',
                            ),
                            keyboardType: TextInputType.number,
                            inputFormatters: [
                              FilteringTextInputFormatter.digitsOnly,
                            ],
                            validator: (value) {
                              if (value == null || value.trim().isEmpty) {
                                return '${l10n.dealValue} ${l10n.isRequired}';
                              }
                              return null;
                            },
                          ),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: TextFormField(
                            controller: _probabilityController,
                            decoration: InputDecoration(
                              labelText: l10n.probability,
                              border: const OutlineInputBorder(),
                              suffixText: '%',
                            ),
                            keyboardType: TextInputType.number,
                            inputFormatters: [
                              FilteringTextInputFormatter.digitsOnly,
                              LengthLimitingTextInputFormatter(3),
                            ],
                            validator: (value) {
                              if (value == null || value.trim().isEmpty) {
                                return '${l10n.probability} ${l10n.isRequired}';
                              }
                              final prob = int.tryParse(value);
                              if (prob == null || prob < 0 || prob > 100) {
                                return '0-100';
                              }
                              return null;
                            },
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                    InkWell(
                      onTap: () => _selectDate(context),
                      child: InputDecorator(
                        decoration: InputDecoration(
                          labelText: l10n.expectedCloseDate,
                          border: const OutlineInputBorder(),
                          suffixIcon: const Icon(Icons.calendar_today),
                        ),
                        child: Text(
                          _expectedCloseDate == null
                              ? l10n.optional
                              : DateFormat(
                                  'yyyy-MM-dd',
                                ).format(_expectedCloseDate!),
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _descriptionController,
                      decoration: InputDecoration(
                        labelText: l10n.description,
                        border: const OutlineInputBorder(),
                        hintText: l10n.optional,
                      ),
                      maxLines: 3,
                    ),
                    const SizedBox(height: 32),
                    ElevatedButton(
                      onPressed: _convertLead,
                      style: ElevatedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        backgroundColor: Theme.of(context).primaryColor,
                        foregroundColor: Colors.white,
                      ),
                      child: Text(l10n.convert),
                    ),
                  ],
                ),
              ),
            ),
    );
  }
}
