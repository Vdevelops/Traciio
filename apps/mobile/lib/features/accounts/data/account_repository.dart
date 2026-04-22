import 'package:dio/dio.dart';

import '../../../core/network/connectivity_service.dart';
import '../../../core/storage/offline_storage.dart';
import 'models/account.dart';

class AccountRepository {
  AccountRepository(this._dio, this._connectivity);

  final Dio _dio;
  final ConnectivityService _connectivity;

  Future<AccountListResponse> getAccounts({
    int page = 1,
    int perPage = 20,
    String? search,
    String? status,
    String? categoryId,
    String? assignedTo,
    bool forceRefresh = false,
    Function(AccountListResponse)? onBackgroundUpdate,
  }) async {
    // 1. Try to load from cache first (offline-first) - only for first page and no search
    if (!forceRefresh && page == 1 && (search == null || search.isEmpty)) {
      final cachedAccounts = await OfflineStorage.getAccounts();
      if (cachedAccounts != null && cachedAccounts.isNotEmpty) {
        try {
          // Convert cached data to Account objects
          final accounts = cachedAccounts
              .map((json) => Account.fromJson(json))
              .toList();

          // Return cached response immediately
          final cachedResponse = AccountListResponse(
            items: accounts,
            pagination: Pagination(
              page: 1,
              perPage: accounts.length,
              total: accounts.length,
              totalPages: 1,
            ),
          );

          // Trigger background refresh if online
          if (_connectivity.isOnline && !forceRefresh) {
            _fetchAndUpdateInBackground(
              page: page,
              perPage: perPage,
              search: search,
              status: status,
              categoryId: categoryId,
              assignedTo: assignedTo,
              onBackgroundUpdate: onBackgroundUpdate,
            );
          }

          return cachedResponse;
        } catch (e) {
          // If parsing fails, continue to API call
        }
      }
    }

    // 2. If online, fetch from API
    if (_connectivity.isOnline) {
      try {
        final queryParams = <String, dynamic>{
          'page': page,
          'per_page': perPage,
        };

        if (search != null && search.isNotEmpty) {
          queryParams['search'] = search;
        }
        if (status != null && status.isNotEmpty) {
          queryParams['status'] = status;
        }
        if (categoryId != null && categoryId.isNotEmpty) {
          queryParams['category_id'] = categoryId;
        }
        if (assignedTo != null && assignedTo.isNotEmpty) {
          queryParams['assigned_to'] = assignedTo;
        }

        final response = await _dio.get(
          '/api/v1/accounts',
          queryParameters: queryParams,
        );

        if (response.data['success'] == true) {
          try {
            final accountListResponse = AccountListResponse.fromJson(
              response.data,
            );

            // 3. Save to cache (only for first page and no search)
            if (page == 1 && (search == null || search.isEmpty)) {
              final accountsJson = accountListResponse.items
                  .map((account) => account.toJson())
                  .toList();
              await OfflineStorage.saveAccounts(accountsJson);
            }

            return accountListResponse;
          } catch (e) {
            throw Exception(
              'Failed to parse accounts response: $e. Response: ${response.data}',
            );
          }
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to fetch accounts',
          );
        }
      } on DioException catch (e) {
        // If API fails, try to return cached data if available
        if (page == 1 && (search == null || search.isEmpty)) {
          final cachedAccounts = await OfflineStorage.getAccounts();
          if (cachedAccounts != null && cachedAccounts.isNotEmpty) {
            try {
              final accounts = cachedAccounts
                  .map((json) => Account.fromJson(json))
                  .toList();
              return AccountListResponse(
                items: accounts,
                pagination: Pagination(
                  page: 1,
                  perPage: accounts.length,
                  total: accounts.length,
                  totalPages: 1,
                ),
              );
            } catch (_) {
              // Ignore parsing errors
            }
          }
        }

        if (e.response != null) {
          final error = e.response!.data;
          if (error is Map<String, dynamic> && error['error'] != null) {
            throw Exception(
              error['error']['message'] ?? 'Failed to fetch accounts',
            );
          }
        }
        throw Exception('Failed to fetch accounts: ${e.message}');
      } catch (e) {
        // If other error, try cached data
        if (page == 1 && (search == null || search.isEmpty)) {
          final cachedAccounts = await OfflineStorage.getAccounts();
          if (cachedAccounts != null && cachedAccounts.isNotEmpty) {
            try {
              final accounts = cachedAccounts
                  .map((json) => Account.fromJson(json))
                  .toList();
              return AccountListResponse(
                items: accounts,
                pagination: Pagination(
                  page: 1,
                  perPage: accounts.length,
                  total: accounts.length,
                  totalPages: 1,
                ),
              );
            } catch (_) {
              // Ignore parsing errors
            }
          }
        }
        throw Exception('Failed to fetch accounts: $e');
      }
    }

    // 4. Offline: return cached data or throw error
    if (page == 1 && (search == null || search.isEmpty)) {
      final cachedAccounts = await OfflineStorage.getAccounts();
      if (cachedAccounts != null && cachedAccounts.isNotEmpty) {
        try {
          final accounts = cachedAccounts
              .map((json) => Account.fromJson(json))
              .toList();
          return AccountListResponse(
            items: accounts,
            pagination: Pagination(
              page: 1,
              perPage: accounts.length,
              total: accounts.length,
              totalPages: 1,
            ),
          );
        } catch (e) {
          throw Exception('Failed to load cached accounts: $e');
        }
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch accounts in background and update cache + UI
  Future<void> _fetchAndUpdateInBackground({
    required int page,
    required int perPage,
    String? search,
    String? status,
    String? categoryId,
    String? assignedTo,
    Function(AccountListResponse)? onBackgroundUpdate,
  }) async {
    try {
      final queryParams = <String, dynamic>{'page': page, 'per_page': perPage};

      if (search != null && search.isNotEmpty) {
        queryParams['search'] = search;
      }
      if (status != null && status.isNotEmpty) {
        queryParams['status'] = status;
      }
      if (categoryId != null && categoryId.isNotEmpty) {
        queryParams['category_id'] = categoryId;
      }
      if (assignedTo != null && assignedTo.isNotEmpty) {
        queryParams['assigned_to'] = assignedTo;
      }

      final response = await _dio.get(
        '/api/v1/accounts',
        queryParameters: queryParams,
      );

      if (response.data['success'] == true) {
        final accountListResponse = AccountListResponse.fromJson(response.data);

        // Save to cache
        if (page == 1 && (search == null || search.isEmpty)) {
          final accountsJson = accountListResponse.items
              .map((account) => account.toJson())
              .toList();
          await OfflineStorage.saveAccounts(accountsJson);
        }

        // Notify UI to update with fresh data
        onBackgroundUpdate?.call(accountListResponse);
      }
    } catch (e) {
      // Silently fail in background
    }
  }

  Future<Account> getAccountById(String id) async {
    // 1. Try to load from cache first (offline-first)
    final cachedAccount = await OfflineStorage.getAccountDetail(id);
    if (cachedAccount != null) {
      try {
        final account = Account.fromJson(cachedAccount);
        // If online, fetch from API in background to update cache
        if (_connectivity.isOnline) {
          _fetchAndUpdateAccountDetail(id).catchError((_) {
            // Ignore errors, use cached data
          });
        }
        return account;
      } catch (e) {
        // If parsing fails, continue to API call
      }
    }

    // 2. If online, fetch from API
    if (_connectivity.isOnline) {
      try {
        final response = await _dio.get('/api/v1/accounts/$id');

        if (response.data['success'] == true) {
          final account = Account.fromJson(
            response.data['data'] as Map<String, dynamic>,
          );

          // 3. Save to cache
          await OfflineStorage.saveAccountDetail(id, account.toJson());

          return account;
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to fetch account',
          );
        }
      } on DioException catch (e) {
        // If API fails, try to return cached data if available
        if (cachedAccount != null) {
          try {
            return Account.fromJson(cachedAccount);
          } catch (_) {
            // Ignore parsing errors
          }
        }

        if (e.response != null) {
          final error = e.response!.data;
          if (error is Map<String, dynamic> && error['error'] != null) {
            throw Exception(
              error['error']['message'] ?? 'Failed to fetch account',
            );
          }
        }
        throw Exception('Failed to fetch account: ${e.message}');
      } catch (e) {
        // If other error, try cached data
        if (cachedAccount != null) {
          try {
            return Account.fromJson(cachedAccount);
          } catch (_) {
            // Ignore parsing errors
          }
        }
        throw Exception('Failed to fetch account: $e');
      }
    }

    // 4. Offline: return cached data or throw error
    if (cachedAccount != null) {
      try {
        return Account.fromJson(cachedAccount);
      } catch (e) {
        throw Exception('Failed to load cached account: $e');
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch account detail from API and update cache (background operation)
  Future<void> _fetchAndUpdateAccountDetail(String id) async {
    try {
      final response = await _dio.get('/api/v1/accounts/$id');
      if (response.data['success'] == true) {
        final account = Account.fromJson(
          response.data['data'] as Map<String, dynamic>,
        );
        await OfflineStorage.saveAccountDetail(id, account.toJson());
      }
    } catch (_) {
      // Ignore errors in background operation
    }
  }

  Future<Account> createAccount({
    required String name,
    required String categoryId,
    String? address,
    String? city,
    String? province,
    String? phone,
    String? email,
    String? status,
    String? assignedTo,
  }) async {
    try {
      final data = <String, dynamic>{'name': name, 'category_id': categoryId};

      if (address != null && address.isNotEmpty) data['address'] = address;
      if (city != null && city.isNotEmpty) data['city'] = city;
      if (province != null && province.isNotEmpty) data['province'] = province;
      if (phone != null && phone.isNotEmpty) data['phone'] = phone;
      if (email != null && email.isNotEmpty) data['email'] = email;
      if (status != null && status.isNotEmpty) data['status'] = status;
      if (assignedTo != null && assignedTo.isNotEmpty) {
        data['assigned_to'] = assignedTo;
      }

      final response = await _dio.post('/api/v1/accounts', data: data);

      if (response.data['success'] == true) {
        return Account.fromJson(response.data['data'] as Map<String, dynamic>);
      } else {
        final errorData = response.data['error'];
        if (errorData is Map<String, dynamic>) {
          final message = errorData['message'] ?? 'Failed to create account';
          throw Exception(message);
        }
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to create account',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic>) {
          if (error['error'] != null) {
            final errorObj = error['error'];
            if (errorObj is Map<String, dynamic>) {
              final message = errorObj['message'] ?? 'Failed to create account';
              throw Exception(message);
            }
            throw Exception(errorObj.toString());
          }
          // Try to get error message from response
          if (error['message'] != null) {
            throw Exception(error['message']);
          }
        }
        throw Exception('Server error: ${e.response?.statusCode}');
      }
      throw Exception('Network error: ${e.message}');
    } catch (e) {
      throw Exception('Failed to create account: $e');
    }
  }

  Future<Account> updateAccount({
    required String id,
    String? name,
    String? categoryId,
    String? address,
    String? city,
    String? province,
    String? phone,
    String? email,
    String? status,
    String? assignedTo,
  }) async {
    try {
      final data = <String, dynamic>{};

      if (name != null && name.isNotEmpty) data['name'] = name;
      if (categoryId != null && categoryId.isNotEmpty) {
        data['category_id'] = categoryId;
      }
      if (address != null && address.isNotEmpty) data['address'] = address;
      if (city != null && city.isNotEmpty) data['city'] = city;
      if (province != null && province.isNotEmpty) data['province'] = province;
      if (phone != null && phone.isNotEmpty) data['phone'] = phone;
      if (email != null && email.isNotEmpty) data['email'] = email;
      if (status != null && status.isNotEmpty) data['status'] = status;
      if (assignedTo != null && assignedTo.isNotEmpty) {
        data['assigned_to'] = assignedTo;
      }

      final response = await _dio.put('/api/v1/accounts/$id', data: data);

      if (response.data['success'] == true) {
        return Account.fromJson(response.data['data'] as Map<String, dynamic>);
      } else {
        final errorData = response.data['error'];
        if (errorData is Map<String, dynamic>) {
          final message = errorData['message'] ?? 'Failed to update account';
          throw Exception(message);
        }
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to update account',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic>) {
          if (error['error'] != null) {
            final errorObj = error['error'];
            if (errorObj is Map<String, dynamic>) {
              final message = errorObj['message'] ?? 'Failed to update account';
              throw Exception(message);
            }
            throw Exception(errorObj.toString());
          }
          if (error['message'] != null) {
            throw Exception(error['message']);
          }
        }
        throw Exception('Server error: ${e.response?.statusCode}');
      }
      throw Exception('Network error: ${e.message}');
    } catch (e) {
      throw Exception('Failed to update account: $e');
    }
  }

  Future<void> deleteAccount(String id) async {
    try {
      final response = await _dio.delete('/api/v1/accounts/$id');

      if (response.data['success'] != true) {
        final errorData = response.data['error'];
        if (errorData is Map<String, dynamic>) {
          final message = errorData['message'] ?? 'Failed to delete account';
          throw Exception(message);
        }
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to delete account',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic>) {
          if (error['error'] != null) {
            final errorObj = error['error'];
            if (errorObj is Map<String, dynamic>) {
              final message = errorObj['message'] ?? 'Failed to delete account';
              throw Exception(message);
            }
            throw Exception(errorObj.toString());
          }
          if (error['message'] != null) {
            throw Exception(error['message']);
          }
        }
        throw Exception('Server error: ${e.response?.statusCode}');
      }
      throw Exception('Network error: ${e.message}');
    } catch (e) {
      throw Exception('Failed to delete account: $e');
    }
  }
}
