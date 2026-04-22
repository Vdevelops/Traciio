import 'package:dio/dio.dart';

import 'models/account.dart';

class CategoryRepository {
  CategoryRepository(this._dio);

  final Dio _dio;

  Future<List<Category>> getCategories() async {
    try {
      final response = await _dio.get('/api/v1/categories');

      if (response.data['success'] == true) {
        final data = response.data['data'];
        if (data is List) {
          return data
              .map((item) => Category.fromJson(item as Map<String, dynamic>))
              .toList();
        } else if (data is Map<String, dynamic> && data['items'] != null) {
          return (data['items'] as List<dynamic>)
              .map((item) => Category.fromJson(item as Map<String, dynamic>))
              .toList();
        }
        return [];
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to fetch categories',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(error['error']['message'] ?? 'Failed to fetch categories');
        }
      }
      throw Exception('Failed to fetch categories: ${e.message}');
    } catch (e) {
      throw Exception('Failed to fetch categories: $e');
    }
  }
}

