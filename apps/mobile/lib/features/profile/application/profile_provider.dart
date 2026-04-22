import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../auth/application/auth_provider.dart';
import '../../auth/application/auth_state.dart';
import '../data/models/profile.dart';
import '../data/profile_repository.dart';

final profileRepositoryProvider = Provider<ProfileRepository>((ref) {
  final connectivity = ref.watch(connectivityServiceProvider);
  return ProfileRepository(ApiClient.dio, connectivity);
});

final profileProvider = FutureProvider.autoDispose<ProfileResponse>((
  ref,
) async {
  final repository = ref.read(profileRepositoryProvider);
  final authState = ref.watch(authProvider);

  // Ensure user is authenticated
  if (authState.status != AuthStatus.authenticated) {
    throw Exception('User not authenticated');
  }

  return repository.getMyProfile();
});

final updateProfileProvider = FutureProvider.family
    .autoDispose<ProfileUser, UpdateProfileRequest>((ref, request) async {
      final repository = ref.read(profileRepositoryProvider);
      return repository.updateMyProfile(request);
    });

final changePasswordProvider = FutureProvider.family
    .autoDispose<void, ChangePasswordRequest>((ref, request) async {
      final repository = ref.read(profileRepositoryProvider);
      return repository.changeMyPassword(request);
    });
