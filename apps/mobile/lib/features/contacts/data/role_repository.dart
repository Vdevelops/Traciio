import 'package:dio/dio.dart';

import 'models/contact.dart';

class RoleRepository {
  RoleRepository(this._dio);

  final Dio _dio;

  Future<List<Role>> getRoles() async {
    try {
      final response = await _dio.get('/api/v1/contact-roles');

      if (response.data['success'] == true) {
        final data = response.data['data'];
        if (data is List) {
          return data
              .map((item) => Role.fromJson(item as Map<String, dynamic>))
              .toList();
        } else if (data is Map<String, dynamic> && data['items'] != null) {
          return (data['items'] as List<dynamic>)
              .map((item) => Role.fromJson(item as Map<String, dynamic>))
              .toList();
        }
        return [];
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to fetch roles',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(error['error']['message'] ?? 'Failed to fetch roles');
        }
      }
      throw Exception('Failed to fetch roles: ${e.message}');
    } catch (e) {
      throw Exception('Failed to fetch roles: $e');
    }
  }
}

