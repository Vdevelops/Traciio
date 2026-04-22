import '../data/models/account.dart';

class AccountListState {
  const AccountListState({
    this.accounts = const [],
    this.isLoading = false,
    this.isLoadingMore = false,
    this.errorMessage,
    this.pagination,
    this.searchQuery = '',
    this.isOffline = false,
  });

  final List<Account> accounts;
  final bool isLoading;
  final bool isLoadingMore;
  final String? errorMessage;
  final Pagination? pagination;
  final String searchQuery;
  final bool isOffline;

  AccountListState copyWith({
    List<Account>? accounts,
    bool? isLoading,
    bool? isLoadingMore,
    String? errorMessage,
    Pagination? pagination,
    String? searchQuery,
    bool? isOffline,
    bool clearAccounts = false,
  }) {
    return AccountListState(
      accounts: clearAccounts ? const [] : (accounts ?? this.accounts),
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      errorMessage: errorMessage,
      pagination: pagination ?? this.pagination,
      searchQuery: searchQuery ?? this.searchQuery,
      isOffline: isOffline ?? this.isOffline,
    );
  }
}

