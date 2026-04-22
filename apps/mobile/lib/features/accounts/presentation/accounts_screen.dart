import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/l10n/app_localizations.dart';
import '../../../core/permissions/permission_provider.dart';
import '../application/account_provider.dart';
import 'account_list_screen.dart';
import 'account_form_screen.dart';
import 'widgets/account_search_modal.dart';

class AccountsScreen extends ConsumerWidget {
  const AccountsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = AppLocalizations.of(context)!;
    final hasCreateAccountPermission = ref.watch(canCreateProvider('accounts'));

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Text(l10n.accounts),
        automaticallyImplyLeading: false,
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => showAccountSearchModal(context),
            tooltip: l10n.searchAccounts,
          ),
        ],
      ),
      body: const AccountListScreen(hideAppBar: true),
      floatingActionButton: hasCreateAccountPermission
          ? FloatingActionButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => const AccountFormScreen(),
                  ),
                ).then((_) {
                  ref.read(accountListProvider.notifier).refresh();
                });
              },
              child: const Icon(Icons.add),
            )
          : null,
    );
  }
}
