import 'dart:async';
import 'package:flutter/material.dart';

/// Generic search overlay widget to be reused across different features.
class SearchModal extends StatefulWidget {
  final String hintText;
  final String initialQuery;
  final Function(String) onSearch;
  final VoidCallback onClear;

  const SearchModal({
    super.key,
    required this.hintText,
    required this.initialQuery,
    required this.onSearch,
    required this.onClear,
  });

  @override
  State<SearchModal> createState() => _SearchModalState();
}

class _SearchModalState extends State<SearchModal> {
  final TextEditingController _searchController = TextEditingController();
  Timer? _debounceTimer;
  bool _isSearching = false;

  @override
  void initState() {
    super.initState();
    if (widget.initialQuery.isNotEmpty) {
      _searchController.text = widget.initialQuery;
      _isSearching = true;
    }
    _searchController.addListener(_onSearchChanged);
  }

  @override
  void dispose() {
    _debounceTimer?.cancel();
    _searchController.removeListener(_onSearchChanged);
    _searchController.dispose();
    super.dispose();
  }

  void _onSearchChanged() {
    _debounceTimer?.cancel();
    final query = _searchController.text.trim();

    setState(() {
      _isSearching = query.isNotEmpty;
    });

    _debounceTimer = Timer(const Duration(milliseconds: 500), () {
      widget.onSearch(query);
    });
  }

  void _clearSearch() {
    _debounceTimer?.cancel();
    _searchController.clear();
    widget.onClear();
    Navigator.of(context).pop();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final mediaQuery = MediaQuery.of(context);
    final safeAreaTop = mediaQuery.padding.top;

    return SafeArea(
      child: Align(
        alignment: Alignment.topCenter,
        child: Container(
          margin: EdgeInsets.only(top: safeAreaTop + 16),
          padding: const EdgeInsets.fromLTRB(20, 0, 20, 0),
          child: Material(
            color: Colors.transparent,
            child: Container(
              decoration: BoxDecoration(
                color: colorScheme.surfaceContainerHighest.withValues(
                  alpha: 0.3,
                ),
                borderRadius: BorderRadius.circular(30),
              ),
              child: TextField(
                controller: _searchController,
                autofocus: true,
                style: theme.textTheme.titleMedium,
                decoration: InputDecoration(
                  hintText: widget.hintText,
                  hintStyle: TextStyle(color: theme.hintColor, fontSize: 16),
                  prefixIcon: Icon(
                    Icons.search,
                    color: theme.hintColor,
                    size: 24,
                  ),
                  suffixIcon: _isSearching
                      ? IconButton(
                          icon: const Icon(Icons.clear),
                          onPressed: _clearSearch,
                        )
                      : null,
                  border: InputBorder.none,
                  focusedBorder: InputBorder.none,
                  enabledBorder: InputBorder.none,
                  contentPadding: const EdgeInsets.symmetric(
                    horizontal: 20,
                    vertical: 16,
                  ),
                ),
                textInputAction: TextInputAction.search,
                onSubmitted: (value) {
                  final query = value.trim();
                  _debounceTimer?.cancel();
                  widget.onSearch(query);
                  Navigator.of(context).pop();
                },
              ),
            ),
          ),
        ),
      ),
    );
  }
}

/// Helper function to show the search modal
void showSearchModal(
  BuildContext context, {
  required String hintText,
  required String initialQuery,
  required Function(String) onSearch,
  required VoidCallback onClear,
}) {
  showGeneralDialog(
    context: context,
    barrierDismissible: true,
    barrierLabel: hintText,
    barrierColor: Colors.black.withValues(alpha: 0.5),
    transitionDuration: const Duration(milliseconds: 200),
    pageBuilder: (context, animation, secondaryAnimation) {
      return SearchModal(
        hintText: hintText,
        initialQuery: initialQuery,
        onSearch: onSearch,
        onClear: onClear,
      );
    },
    transitionBuilder: (context, animation, secondaryAnimation, child) {
      return SlideTransition(
        position: Tween<Offset>(
          begin: const Offset(0, -1),
          end: Offset.zero,
        ).animate(CurvedAnimation(parent: animation, curve: Curves.easeOut)),
        child: child,
      );
    },
  );
}
