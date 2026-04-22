import 'package:dio/dio.dart';

import '../../../core/network/connectivity_service.dart';
import '../../../core/storage/offline_storage.dart';
import 'models/contact.dart';

class ContactRepository {
  ContactRepository(this._dio, this._connectivity);

  final Dio _dio;
  final ConnectivityService _connectivity;

  Future<ContactListResponse> getContacts({
    int page = 1,
    int perPage = 20,
    String? search,
    String? accountId,
    String? roleId,
    bool forceRefresh = false,
    Function(ContactListResponse)? onBackgroundUpdate,
  }) async {
    // 1. Try to load from cache first (offline-first) - only for first page and no filters
    if (!forceRefresh &&
        page == 1 &&
        (search == null || search.isEmpty) &&
        accountId == null) {
      final cachedContacts = await OfflineStorage.getContacts();
      if (cachedContacts != null && cachedContacts.isNotEmpty) {
        try {
          final contacts = cachedContacts
              .map((json) => Contact.fromJson(json))
              .toList();
          final cachedResponse = ContactListResponse(
            items: contacts,
            pagination: Pagination(
              page: 1,
              perPage: contacts.length,
              total: contacts.length,
              totalPages: 1,
            ),
          );

          // Trigger background refresh if online
          if (_connectivity.isOnline && !forceRefresh) {
            _fetchAndUpdateInBackground(
              page: page,
              perPage: perPage,
              search: search,
              accountId: accountId,
              roleId: roleId,
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
        if (accountId != null && accountId.isNotEmpty) {
          queryParams['account_id'] = accountId;
        }
        if (roleId != null && roleId.isNotEmpty) {
          queryParams['role_id'] = roleId;
        }

        final response = await _dio.get(
          '/api/v1/contacts',
          queryParameters: queryParams,
        );

        if (response.data['success'] == true) {
          try {
            final contactListResponse = ContactListResponse.fromJson(
              response.data,
            );

            // 3. Save to cache (only for first page and no filters)
            if (page == 1 &&
                (search == null || search.isEmpty) &&
                accountId == null) {
              final contactsJson = contactListResponse.items
                  .map((contact) => contact.toJson())
                  .toList();
              await OfflineStorage.saveContacts(contactsJson);
            }

            return contactListResponse;
          } catch (e) {
            throw Exception(
              'Failed to parse contacts response: $e. Response: ${response.data}',
            );
          }
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to fetch contacts',
          );
        }
      } on DioException catch (e) {
        // If API fails, try to return cached data if available
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            accountId == null) {
          final cachedContacts = await OfflineStorage.getContacts();
          if (cachedContacts != null && cachedContacts.isNotEmpty) {
            try {
              final contacts = cachedContacts
                  .map((json) => Contact.fromJson(json))
                  .toList();
              return ContactListResponse(
                items: contacts,
                pagination: Pagination(
                  page: 1,
                  perPage: contacts.length,
                  total: contacts.length,
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
              error['error']['message'] ?? 'Failed to fetch contacts',
            );
          }
        }
        throw Exception('Failed to fetch contacts: ${e.message}');
      } catch (e) {
        // If other error, try cached data
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            accountId == null) {
          final cachedContacts = await OfflineStorage.getContacts();
          if (cachedContacts != null && cachedContacts.isNotEmpty) {
            try {
              final contacts = cachedContacts
                  .map((json) => Contact.fromJson(json))
                  .toList();
              return ContactListResponse(
                items: contacts,
                pagination: Pagination(
                  page: 1,
                  perPage: contacts.length,
                  total: contacts.length,
                  totalPages: 1,
                ),
              );
            } catch (_) {
              // Ignore parsing errors
            }
          }
        }
        throw Exception('Failed to fetch contacts: $e');
      }
    }

    // 4. Offline: return cached data or throw error
    if (page == 1 && (search == null || search.isEmpty) && accountId == null) {
      final cachedContacts = await OfflineStorage.getContacts();
      if (cachedContacts != null && cachedContacts.isNotEmpty) {
        try {
          final contacts = cachedContacts
              .map((json) => Contact.fromJson(json))
              .toList();
          return ContactListResponse(
            items: contacts,
            pagination: Pagination(
              page: 1,
              perPage: contacts.length,
              total: contacts.length,
              totalPages: 1,
            ),
          );
        } catch (e) {
          throw Exception('Failed to load cached contacts: $e');
        }
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch contacts in background and update cache + UI
  Future<void> _fetchAndUpdateInBackground({
    required int page,
    required int perPage,
    String? search,
    String? accountId,
    String? roleId,
    Function(ContactListResponse)? onBackgroundUpdate,
  }) async {
    try {
      final queryParams = <String, dynamic>{'page': page, 'per_page': perPage};

      if (search != null && search.isNotEmpty) {
        queryParams['search'] = search;
      }
      if (accountId != null && accountId.isNotEmpty) {
        queryParams['account_id'] = accountId;
      }
      if (roleId != null && roleId.isNotEmpty) {
        queryParams['role_id'] = roleId;
      }

      final response = await _dio.get(
        '/api/v1/contacts',
        queryParameters: queryParams,
      );

      if (response.data['success'] == true) {
        final contactListResponse = ContactListResponse.fromJson(response.data);

        // Save to cache
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            accountId == null) {
          final contactsJson = contactListResponse.items
              .map((contact) => contact.toJson())
              .toList();
          await OfflineStorage.saveContacts(contactsJson);
        }

        // Notify UI to update with fresh data
        onBackgroundUpdate?.call(contactListResponse);
      }
    } catch (e) {
      // Silently fail in background
    }
  }

  Future<Contact> getContactById(String id) async {
    // 1. Try to load from cache first (offline-first)
    final cachedContact = await OfflineStorage.getContactDetail(id);
    if (cachedContact != null) {
      try {
        final contact = Contact.fromJson(cachedContact);
        // If online, fetch from API in background to update cache
        if (_connectivity.isOnline) {
          _fetchAndUpdateContactDetail(id).catchError((_) {
            // Ignore errors, use cached data
          });
        }
        return contact;
      } catch (e) {
        // If parsing fails, continue to API call
      }
    }

    // 2. If online, fetch from API
    if (_connectivity.isOnline) {
      try {
        final response = await _dio.get('/api/v1/contacts/$id');

        if (response.data['success'] == true) {
          final contact = Contact.fromJson(
            response.data['data'] as Map<String, dynamic>,
          );

          // 3. Save to cache
          await OfflineStorage.saveContactDetail(id, contact.toJson());

          return contact;
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to fetch contact',
          );
        }
      } on DioException catch (e) {
        // If API fails, try to return cached data if available
        if (cachedContact != null) {
          try {
            return Contact.fromJson(cachedContact);
          } catch (_) {
            // Ignore parsing errors
          }
        }

        if (e.response != null) {
          final error = e.response!.data;
          if (error is Map<String, dynamic> && error['error'] != null) {
            throw Exception(
              error['error']['message'] ?? 'Failed to fetch contact',
            );
          }
        }
        throw Exception('Failed to fetch contact: ${e.message}');
      } catch (e) {
        // If other error, try cached data
        if (cachedContact != null) {
          try {
            return Contact.fromJson(cachedContact);
          } catch (_) {
            // Ignore parsing errors
          }
        }
        throw Exception('Failed to fetch contact: $e');
      }
    }

    // 4. Offline: return cached data or throw error
    if (cachedContact != null) {
      try {
        return Contact.fromJson(cachedContact);
      } catch (e) {
        throw Exception('Failed to load cached contact: $e');
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch contact detail from API and update cache (background operation)
  Future<void> _fetchAndUpdateContactDetail(String id) async {
    try {
      final response = await _dio.get('/api/v1/contacts/$id');
      if (response.data['success'] == true) {
        final contact = Contact.fromJson(
          response.data['data'] as Map<String, dynamic>,
        );
        await OfflineStorage.saveContactDetail(id, contact.toJson());
      }
    } catch (_) {
      // Ignore errors in background operation
    }
  }

  Future<Contact> createContact({
    required String accountId,
    required String name,
    required String roleId,
    String? phone,
    String? email,
    String? position,
    String? notes,
  }) async {
    try {
      final data = <String, dynamic>{
        'account_id': accountId,
        'name': name,
        'role_id': roleId,
      };

      if (phone != null && phone.isNotEmpty) data['phone'] = phone;
      if (email != null && email.isNotEmpty) data['email'] = email;
      if (position != null && position.isNotEmpty) data['position'] = position;
      if (notes != null && notes.isNotEmpty) data['notes'] = notes;

      final response = await _dio.post('/api/v1/contacts', data: data);

      if (response.data['success'] == true) {
        return Contact.fromJson(response.data['data'] as Map<String, dynamic>);
      } else {
        final errorData = response.data['error'];
        if (errorData is Map<String, dynamic>) {
          final message = errorData['message'] ?? 'Failed to create contact';
          throw Exception(message);
        }
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to create contact',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic>) {
          if (error['error'] != null) {
            final errorObj = error['error'];
            if (errorObj is Map<String, dynamic>) {
              final message = errorObj['message'] ?? 'Failed to create contact';
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
      throw Exception('Failed to create contact: $e');
    }
  }

  Future<Contact> updateContact({
    required String id,
    String? accountId,
    String? name,
    String? roleId,
    String? phone,
    String? email,
    String? position,
    String? notes,
  }) async {
    try {
      final data = <String, dynamic>{};

      if (accountId != null && accountId.isNotEmpty) {
        data['account_id'] = accountId;
      }
      if (name != null && name.isNotEmpty) data['name'] = name;
      if (roleId != null && roleId.isNotEmpty) data['role_id'] = roleId;
      if (phone != null && phone.isNotEmpty) data['phone'] = phone;
      if (email != null && email.isNotEmpty) data['email'] = email;
      if (position != null && position.isNotEmpty) data['position'] = position;
      if (notes != null && notes.isNotEmpty) data['notes'] = notes;

      final response = await _dio.put('/api/v1/contacts/$id', data: data);

      if (response.data['success'] == true) {
        return Contact.fromJson(response.data['data'] as Map<String, dynamic>);
      } else {
        final errorData = response.data['error'];
        if (errorData is Map<String, dynamic>) {
          final message = errorData['message'] ?? 'Failed to update contact';
          throw Exception(message);
        }
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to update contact',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic>) {
          if (error['error'] != null) {
            final errorObj = error['error'];
            if (errorObj is Map<String, dynamic>) {
              final message = errorObj['message'] ?? 'Failed to update contact';
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
      throw Exception('Failed to update contact: $e');
    }
  }

  Future<void> deleteContact(String id) async {
    try {
      final response = await _dio.delete('/api/v1/contacts/$id');

      if (response.data['success'] != true) {
        final errorData = response.data['error'];
        if (errorData is Map<String, dynamic>) {
          final message = errorData['message'] ?? 'Failed to delete contact';
          throw Exception(message);
        }
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to delete contact',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic>) {
          if (error['error'] != null) {
            final errorObj = error['error'];
            if (errorObj is Map<String, dynamic>) {
              final message = errorObj['message'] ?? 'Failed to delete contact';
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
      throw Exception('Failed to delete contact: $e');
    }
  }
}
