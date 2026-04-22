import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../application/pipeline_provider.dart';
import '../../data/models/deal_model.dart';
import '../../../visit_reports/presentation/widgets/searchable_dropdown.dart';

class DealFormScreen extends ConsumerStatefulWidget {
  const DealFormScreen({super.key, this.dealId});

  final String? dealId;

  @override
  ConsumerState<DealFormScreen> createState() => _DealFormScreenState();
}

class _DealFormScreenState extends ConsumerState<DealFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _valueController = TextEditingController();
  final _notesController = TextEditingController();

  String? _selectedAccountId;
  String? _selectedContactId;
  String? _selectedStageId;
  List<DealProductItem> _selectedProducts = [];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(pipelineProvider.notifier).loadFormData(forceRefresh: true);
      if (widget.dealId == null) {
        // Default stage if creating
        final state = ref.read(pipelineProvider);
        if (state.stages.isNotEmpty) {
          setState(() {
            _selectedStageId = state.selectedStageId ?? state.stages.first.id;
          });
        }
      } else {
        // Load existing deal data if editing
        final state = ref.read(pipelineProvider);
        final deal = state.deals.firstWhere((d) => d.id == widget.dealId);
        _titleController.text = deal.title;
        _valueController.text = deal.value.toString();
        _notesController.text = deal.notes ?? '';
        setState(() {
          _selectedAccountId = deal.accountId;
          _selectedContactId = deal.contactId;
          _selectedStageId = deal.stageId;
          _selectedProducts = List.from(deal.productItems ?? []);
        });
      }
    });
  }

  @override
  void dispose() {
    _titleController.dispose();
    _valueController.dispose();
    _notesController.dispose();
    super.dispose();
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;

    final navigator = Navigator.of(context);
    final scaffoldMessenger = ScaffoldMessenger.of(context);

    try {
      final title = _titleController.text.trim();
      final valueText = _valueController.text.trim();
      final value =
          int.tryParse(valueText.replaceAll(RegExp(r'[^0-9]'), '')) ?? 0;
      final notes = _notesController.text.trim();

      final productItems = _selectedProducts
          .map(
            (item) => {
              'product_id': item.productId,
              'quantity': item.quantity,
              'unit_price': item.unitPrice,
              'discount_amount': item.discountAmount,
              'notes': item.notes,
            },
          )
          .toList();

      if (_selectedStageId == null) {
        scaffoldMessenger.showSnackBar(
          const SnackBar(content: Text('Please select a stage')),
        );
        return;
      }

      await ref
          .read(pipelineProvider.notifier)
          .saveDeal(
            title: title,
            value: value,
            accountId: _selectedAccountId,
            contactId: _selectedContactId,
            stageId: _selectedStageId!,
            notes: notes,
            productItems: productItems,
            dealId: widget.dealId,
          );

      if (mounted) {
        scaffoldMessenger.showSnackBar(
          SnackBar(
            content: Text(
              widget.dealId == null
                  ? 'Deal created successfully'
                  : 'Deal updated successfully',
            ),
          ),
        );
        navigator.pop(true);
      }
    } catch (e) {
      if (mounted) {
        scaffoldMessenger.showSnackBar(
          SnackBar(content: Text('Error: ${e.toString()}')),
        );
      }
    }
  }

  void _updateTotalValue() {
    int total = 0;
    for (final item in _selectedProducts) {
      total += item.subtotal;
    }
    _valueController.text = total.toString();
  }

  void _showProductSelection() {
    final products = ref.read(pipelineProvider).products;
    if (products.isEmpty) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('No products available')));
      return;
    }

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => DraggableScrollableSheet(
        initialChildSize: 0.7,
        maxChildSize: 0.9,
        minChildSize: 0.5,
        expand: false,
        builder: (context, scrollController) => Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: Text(
                'Select Product',
                style: Theme.of(context).textTheme.titleLarge,
              ),
            ),
            Expanded(
              child: ListView.builder(
                controller: scrollController,
                itemCount: products.length,
                itemBuilder: (context, index) {
                  final product = products[index];
                  return ListTile(
                    title: Text(product.name),
                    subtitle: Text('SKU: ${product.sku} - Rp ${product.price}'),
                    onTap: () {
                      Navigator.pop(context);
                      _showProductEntryDialog(product: product);
                    },
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _showProductEntryDialog({
    Product? product,
    DealProductItem? existingItem,
    int? index,
  }) async {
    if (product == null && existingItem == null) return;

    final name = product?.name ?? existingItem!.productName;
    final basePrice = product?.price ?? existingItem!.unitPrice;

    int quantity = existingItem?.quantity ?? 1;
    int unitPrice = existingItem?.unitPrice ?? basePrice;
    int discount = existingItem?.discountAmount ?? 0;

    final formKey = GlobalKey<FormState>();

    await showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(existingItem == null ? 'Add Product' : 'Edit Product'),
        content: Form(
          key: formKey,
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(name, style: Theme.of(context).textTheme.titleMedium),
                const SizedBox(height: 16),
                TextFormField(
                  initialValue: quantity.toString(),
                  decoration: const InputDecoration(
                    labelText: 'Quantity',
                    border: OutlineInputBorder(),
                  ),
                  keyboardType: TextInputType.number,
                  inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                  validator: (v) {
                    if (v == null || v.isEmpty) return 'Required';
                    final n = int.tryParse(v);
                    if (n == null || n <= 0) return 'Must be greater than 0';
                    return null;
                  },
                  onSaved: (v) => quantity = int.parse(v!),
                ),
                const SizedBox(height: 16),
                TextFormField(
                  initialValue: unitPrice.toString(),
                  decoration: const InputDecoration(
                    labelText: 'Unit Price (Rp)',
                    border: OutlineInputBorder(),
                  ),
                  keyboardType: TextInputType.number,
                  inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                  validator: (v) {
                    if (v == null || v.isEmpty) return 'Required';
                    final n = int.tryParse(v);
                    if (n == null || n < 0) return 'Must be 0 or greater';
                    return null;
                  },
                  onSaved: (v) => unitPrice = int.parse(v!),
                ),
                const SizedBox(height: 16),
                TextFormField(
                  initialValue: discount.toString(),
                  decoration: const InputDecoration(
                    labelText: 'Discount (Rp)',
                    border: OutlineInputBorder(),
                  ),
                  keyboardType: TextInputType.number,
                  inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                  validator: (v) {
                    if (v == null || v.isEmpty) return 'Required';
                    final n = int.tryParse(v);
                    if (n == null || n < 0) return 'Must be 0 or greater';
                    return null;
                  },
                  onSaved: (v) => discount = int.parse(v!),
                ),
              ],
            ),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () {
              if (formKey.currentState!.validate()) {
                formKey.currentState!.save();
                setState(() {
                  final newItem = DealProductItem(
                    id: existingItem?.id ?? '',
                    productId: product?.id ?? existingItem!.productId,
                    productName: name,
                    productSku: product?.sku ?? existingItem!.productSku,
                    unitPrice: unitPrice,
                    quantity: quantity,
                    discountAmount: discount,
                    subtotal: (unitPrice * quantity) - discount,
                  );

                  if (index != null) {
                    _selectedProducts[index] = newItem;
                  } else {
                    _selectedProducts.add(newItem);
                  }
                  _updateTotalValue();
                });
                Navigator.pop(context);
              }
            },
            child: const Text('Save'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(pipelineProvider);
    final theme = Theme.of(context);
    final formData = state.formData;

    final accounts =
        (formData['accounts'] as List?)?.cast<Map<String, dynamic>>() ?? [];
    final contacts =
        (formData['contacts'] as List?)?.cast<Map<String, dynamic>>() ?? [];
    final stages = state.stages;

    return Scaffold(
      appBar: AppBar(
        title: Text(widget.dealId == null ? 'Create Deal' : 'Edit Deal'),
      ),
      body: state.isLoading && state.stages.isEmpty
          ? const Center(child: CircularProgressIndicator())
          : Form(
              key: _formKey,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  TextFormField(
                    controller: _titleController,
                    decoration: const InputDecoration(
                      labelText: 'Title *',
                      hintText: 'Enter deal title',
                    ),
                    validator: (v) => v?.isEmpty ?? true ? 'Required' : null,
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _valueController,
                    decoration: const InputDecoration(
                      labelText: 'Value (Rp) *',
                      hintText: 'Enter deal value',
                    ),
                    keyboardType: TextInputType.number,
                    inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                    validator: (v) {
                      if (v == null || v.isEmpty) return 'Required';
                      final n = int.tryParse(v);
                      if (n == null) return 'Enter a valid number';
                      if (n <= 0) return 'Value must be greater than 0';
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  SearchableDropdown<String>(
                    labelText: 'Account *',
                    items: accounts
                        .where(
                          (e) =>
                              e['id'] != null &&
                              e['name'] != null &&
                              (e['name'] as String).trim().isNotEmpty,
                        )
                        .map((e) => e['id'] as String)
                        .toList(),
                    displayText: (id) =>
                        accounts.firstWhere(
                              (e) => e['id'] == id,
                              orElse: () => {'name': 'Unknown'},
                            )['name']
                            as String,
                    selectedValue: _selectedAccountId,
                    onChanged: (v) => setState(() {
                      _selectedAccountId = v;
                      _selectedContactId = null;
                    }),
                  ),
                  const SizedBox(height: 16),
                  SearchableDropdown<String>(
                    labelText: 'Contact',
                    items: contacts
                        .where(
                          (e) =>
                              e['account_id'] == _selectedAccountId &&
                              e['id'] != null &&
                              e['name'] != null &&
                              (e['name'] as String).trim().isNotEmpty,
                        )
                        .map((e) => e['id'] as String)
                        .toList(),
                    displayText: (id) =>
                        contacts.firstWhere(
                              (e) => e['id'] == id,
                              orElse: () => {'name': 'Unknown'},
                            )['name']
                            as String,
                    selectedValue: _selectedContactId,
                    onChanged: (v) => setState(() => _selectedContactId = v),
                    allowNone: true,
                    noneText: 'None (Clear Selection)',
                  ),
                  const SizedBox(height: 16),
                  DropdownButtonFormField<String>(
                    key: ValueKey('stage_$_selectedStageId'),
                    initialValue: stages.any((s) => s.id == _selectedStageId)
                        ? _selectedStageId
                        : null,
                    decoration: const InputDecoration(labelText: 'Stage *'),
                    items: stages
                        .map(
                          (s) => DropdownMenuItem(
                            value: s.id,
                            child: Text(s.name),
                          ),
                        )
                        .toList(),
                    onChanged: (v) => setState(() => _selectedStageId = v),
                    validator: (v) =>
                        v == null ? 'Please select a stage' : null,
                  ),
                  const SizedBox(height: 24),
                  TextFormField(
                    controller: _notesController,
                    decoration: const InputDecoration(
                      labelText: 'Notes',
                      hintText: 'Enter notes',
                    ),
                    maxLines: 3,
                  ),
                  const SizedBox(height: 24),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'Product Items',
                        style: theme.textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      TextButton.icon(
                        onPressed: _showProductSelection,
                        icon: const Icon(Icons.add),
                        label: const Text('Add Product'),
                      ),
                    ],
                  ),
                  const Divider(),
                  if (_selectedProducts.isEmpty)
                    Padding(
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      child: Center(
                        child: Text(
                          'No products added yet',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                        ),
                      ),
                    )
                  else
                    ..._selectedProducts.asMap().entries.map((entry) {
                      final index = entry.key;
                      final item = entry.value;
                      return ListTile(
                        contentPadding: EdgeInsets.zero,
                        title: Text(item.productName),
                        subtitle: Text(
                          '${item.quantity} x Rp ${item.unitPrice} - Disc Rp ${item.discountAmount} = Rp ${item.subtotal}',
                        ),
                        onTap: () => _showProductEntryDialog(
                          existingItem: item,
                          index: index,
                        ),
                        trailing: IconButton(
                          icon: const Icon(Icons.remove_circle_outline),
                          onPressed: () => setState(() {
                            _selectedProducts.removeAt(index);
                            _updateTotalValue();
                          }),
                        ),
                      );
                    }),
                  const SizedBox(height: 32),
                  FilledButton(
                    onPressed: _handleSubmit,
                    child: const Text('Save Deal'),
                  ),
                ],
              ),
            ),
    );
  }
}
