import 'package:flutter/material.dart';

/// Custom searchable dropdown widget untuk selection dengan search functionality
class SearchableDropdown<T> extends StatelessWidget {
  const SearchableDropdown({
    super.key,
    required this.labelText,
    this.hintText,
    required this.items,
    required this.displayText,
    this.selectedValue,
    this.onChanged,
    this.validator,
    this.icon,
    this.allowNone = false,
    this.noneText = 'None',
    this.searchHint = 'Search...',
  });

  final String labelText;
  final String? hintText;
  final List<T> items;
  final String Function(T) displayText;
  final T? selectedValue;
  final ValueChanged<T?>? onChanged;
  final String? Function(String?)? validator;
  final IconData? icon;
  final bool allowNone;
  final String noneText;
  final String searchHint;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    // Validate that selectedValue exists in items
    // If not (e.g., after API reset), clear the selection
    final isValidSelection =
        selectedValue != null && items.any((item) => item == selectedValue);

    // If selection is invalid, notify parent to clear it
    if (selectedValue != null && !isValidSelection) {
      // Schedule the callback after build completes
      WidgetsBinding.instance.addPostFrameCallback((_) {
        onChanged?.call(null);
      });
    }

    final displayValue = isValidSelection
        ? displayText(selectedValue as T)
        : '';

    return TextFormField(
      readOnly: true,
      onTap: () => _showSearchDialog(context),
      style: theme.textTheme.bodyLarge?.copyWith(color: colorScheme.onSurface),
      decoration: InputDecoration(
        labelText: labelText,
        hintText: hintText ?? 'Tap to select...',
        prefixIcon: icon != null ? Icon(icon) : null,
        suffixIcon: const Icon(Icons.arrow_drop_down),
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
        filled: true,
        fillColor: colorScheme.surfaceContainerHighest,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 16,
          vertical: 18,
        ),
      ),
      controller: TextEditingController(text: displayValue),
      validator: validator != null
          ? (String? value) {
              // Validate based on selectedValue, not the text field value
              // Convert T? to String? for validation
              final stringValue = isValidSelection
                  ? displayText(selectedValue as T)
                  : null;
              return validator!(stringValue);
            }
          : null,
    );
  }

  void _showSearchDialog(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final searchController = TextEditingController();
    List<T> filteredItems = List.from(items);

    showDialog(
      context: context,
      builder: (dialogContext) {
        return StatefulBuilder(
          builder: (context, setState) {
            void filterItems(String query) {
              setState(() {
                if (query.isEmpty) {
                  filteredItems = List.from(items);
                } else {
                  filteredItems = items
                      .where(
                        (item) => displayText(
                          item,
                        ).toLowerCase().contains(query.toLowerCase()),
                      )
                      .toList();
                }
              });
            }

            return Dialog(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(20),
              ),
              child: Container(
                constraints: BoxConstraints(
                  maxHeight: MediaQuery.of(context).size.height * 0.7,
                ),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    // Header with search
                    Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        children: [
                          Row(
                            children: [
                              Expanded(
                                child: Text(
                                  labelText,
                                  style: theme.textTheme.titleLarge?.copyWith(
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                              ),
                              IconButton(
                                icon: const Icon(Icons.close),
                                onPressed: () => Navigator.of(context).pop(),
                              ),
                            ],
                          ),
                          const SizedBox(height: 12),
                          TextField(
                            controller: searchController,
                            decoration: InputDecoration(
                              hintText: searchHint,
                              prefixIcon: const Icon(Icons.search),
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                              filled: true,
                              fillColor: colorScheme.surfaceContainerHighest,
                              contentPadding: const EdgeInsets.symmetric(
                                horizontal: 16,
                                vertical: 12,
                              ),
                            ),
                            autofocus: true,
                            onChanged: filterItems,
                          ),
                        ],
                      ),
                    ),
                    // Divider
                    Divider(height: 1, color: colorScheme.outlineVariant),
                    // List of items
                    Flexible(
                      child: ListView.builder(
                        shrinkWrap: true,
                        itemCount: allowNone
                            ? filteredItems.length + 1
                            : filteredItems.length,
                        itemBuilder: (context, index) {
                          if (allowNone && index == 0) {
                            final isSelected = selectedValue == null;
                            return ListTile(
                              leading: Icon(
                                Icons.cancel_outlined,
                                color: isSelected
                                    ? colorScheme.primary
                                    : colorScheme.onSurfaceVariant,
                              ),
                              title: Text(
                                noneText,
                                style: theme.textTheme.bodyLarge?.copyWith(
                                  fontStyle: FontStyle.italic,
                                  color: isSelected
                                      ? colorScheme.primary
                                      : colorScheme.onSurfaceVariant,
                                  fontWeight: isSelected
                                      ? FontWeight.w600
                                      : FontWeight.normal,
                                ),
                              ),
                              trailing: isSelected
                                  ? Icon(
                                      Icons.check_circle,
                                      color: colorScheme.primary,
                                    )
                                  : null,
                              onTap: () {
                                Navigator.of(context).pop();
                                onChanged?.call(null);
                              },
                            );
                          }

                          final itemIndex = allowNone ? index - 1 : index;
                          if (itemIndex >= filteredItems.length) {
                            return const SizedBox.shrink();
                          }

                          final item = filteredItems[itemIndex];
                          final isSelected = selectedValue == item;
                          final displayName = displayText(item);

                          return ListTile(
                            leading: icon != null
                                ? Icon(
                                    icon,
                                    color: isSelected
                                        ? colorScheme.primary
                                        : colorScheme.onSurfaceVariant,
                                  )
                                : null,
                            title: Text(
                              displayName,
                              style: theme.textTheme.bodyLarge?.copyWith(
                                color: isSelected
                                    ? colorScheme.primary
                                    : colorScheme.onSurface,
                                fontWeight: isSelected
                                    ? FontWeight.w600
                                    : FontWeight.normal,
                              ),
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                            ),
                            trailing: isSelected
                                ? Icon(
                                    Icons.check_circle,
                                    color: colorScheme.primary,
                                  )
                                : null,
                            onTap: () {
                              Navigator.of(context).pop();
                              onChanged?.call(item);
                            },
                          );
                        },
                      ),
                    ),
                    // Empty state (only show if no items and not allowNone)
                    if (filteredItems.isEmpty && !allowNone && items.isNotEmpty)
                      Padding(
                        padding: const EdgeInsets.all(32),
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(
                              Icons.search_off,
                              size: 48,
                              color: colorScheme.onSurfaceVariant,
                            ),
                            const SizedBox(height: 16),
                            Text(
                              'No results found',
                              style: theme.textTheme.bodyLarge?.copyWith(
                                color: colorScheme.onSurfaceVariant,
                              ),
                            ),
                          ],
                        ),
                      ),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
  }
}
