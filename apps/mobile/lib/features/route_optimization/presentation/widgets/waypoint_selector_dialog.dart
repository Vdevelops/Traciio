import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/api_client.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../../../../core/l10n/app_localizations.dart';
import '../../../accounts/application/account_provider.dart';
import '../../../visit_reports/application/visit_report_provider.dart';
import '../../../visit_reports/data/models/visit_report.dart' hide AccountInfo;
import '../../data/models/waypoint.dart';

class WaypointSelectorDialog extends ConsumerStatefulWidget {
  const WaypointSelectorDialog({
    super.key,
    required this.onSelect,
  });

  final Function(List<Waypoint>) onSelect;

  @override
  ConsumerState<WaypointSelectorDialog> createState() =>
      _WaypointSelectorDialogState();
}

class _WaypointSelectorDialogState
    extends ConsumerState<WaypointSelectorDialog>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final Set<String> _selectedAccounts = {};
  final Set<String> _selectedVisitReports = {};
  final TextEditingController _searchAccountsController =
      TextEditingController();
  final TextEditingController _searchVisitsController =
      TextEditingController();
  Timer? _accountSearchDebounce;
  Timer? _visitSearchDebounce;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchAccountsController.dispose();
    _searchVisitsController.dispose();
    _accountSearchDebounce?.cancel();
    _visitSearchDebounce?.cancel();
    super.dispose();
  }

  /// Geocode address to get coordinates using backend API (more reliable)
  /// Falls back to direct Nominatim if backend fails
  /// Same implementation as web for consistency
  Future<Map<String, double>?> _geocodeAddress(String address) async {
    // Try backend geocoding first (same as web - uses backend API with fallback)
    try {
      final response = await ApiClient.dio.post(
        '/api/v1/mobile/geocoding/geocode',
        data: {'address': address},
      ).timeout(const Duration(seconds: 10)); // Increased timeout to match web

      if (response.data['success'] == true && response.data['data'] != null) {
        final data = response.data['data'] as Map<String, dynamic>;
        final lat = (data['latitude'] as num?)?.toDouble();
        final lng = (data['longitude'] as num?)?.toDouble();
        
        if (lat != null && lng != null) {
          return {
            'lat': lat,
            'lng': lng,
          };
        }
      }
    } on DioException catch (e) {
      // Handle specific error cases (same as web)
      if (e.response != null) {
        final errorData = e.response!.data;
        if (errorData is Map<String, dynamic> && 
            errorData['error'] != null &&
            errorData['error']['code'] == 'GEOCODING_NO_RESULTS') {
          debugPrint('No geocoding results found for "$address"');
          rethrow; // Re-throw to be handled by caller
        }
      }
      debugPrint('Backend geocoding failed for "$address": $e');
    } catch (e) {
      debugPrint('Backend geocoding error for "$address": $e');
    }

    // Fallback to direct Nominatim (same as web fallback, but web doesn't use this)
    // Note: Web uses backend API only, but mobile has this as extra fallback
    try {
      // Check if address already contains "Indonesia" to avoid duplication
      final queryAddress = address.toLowerCase().contains('indonesia') 
          ? address 
          : '$address, Indonesia';
      final query = Uri.encodeComponent(queryAddress);
      
      final dio = Dio(BaseOptions(
        headers: {
          'User-Agent': 'CRM-Healthcare-Mobile/1.0',
          'Accept-Language': 'id,en',
        },
        connectTimeout: const Duration(seconds: 8),
        receiveTimeout: const Duration(seconds: 8),
      ));
      final response = await dio.get(
        'https://nominatim.openstreetmap.org/search',
        queryParameters: {
          'format': 'json',
          'q': query,
          'limit': 1,
          'countrycodes': 'id',
        },
      );

      if (response.statusCode == 200) {
        final data = response.data as List<dynamic>;
        if (data.isNotEmpty) {
          final result = data[0] as Map<String, dynamic>;
          final lat = double.tryParse(result['lat'] as String? ?? '');
          final lng = double.tryParse(result['lon'] as String? ?? '');
          
          if (lat != null && lng != null) {
            return {
              'lat': lat,
              'lng': lng,
            };
          }
        }
      }
    } catch (e) {
      debugPrint('Nominatim geocoding also failed for "$address": $e');
    }
    
    return null;
  }

  Future<void> _handleConfirm() async {
    final waypoints = <Waypoint>[];

    // Ensure visit reports are loaded (same as web - fetch all visit reports with location)
    final visitReportsState = ref.read(visitReportListProvider);
    if (visitReportsState.visitReports.isEmpty && !visitReportsState.isLoading) {
      // Load visit reports if not already loaded (same as web)
      await ref.read(visitReportListProvider.notifier).loadVisitReports(
        forRouteOptimization: true,
        status: 'approved',
      );
    }

    // Get all visit reports with location (same as web - for account lookup)
    final allVisitReportsState = ref.read(visitReportListProvider);
    final allVisitReports = allVisitReportsState.visitReports
        .where((vr) => vr.checkInLocation != null || vr.checkOutLocation != null)
        .toList();

    // Add selected visit reports with location
    final selectedVisitReports = allVisitReports.where((vr) {
      return _selectedVisitReports.contains(vr.id);
    }).toList();

    for (final visitReport in selectedVisitReports) {
      final location = visitReport.checkInLocation ?? visitReport.checkOutLocation;
      if (location != null) {
        waypoints.add(
          Waypoint(
            lat: location.latitude,
            lng: location.longitude,
            address: location.address ?? visitReport.account?.name,
            accountId: visitReport.accountId,
            accountName: visitReport.account?.name,
            visitReportId: visitReport.id,
            account: visitReport.account != null
                ? AccountInfo(
                    id: visitReport.account!.id,
                    name: visitReport.account!.name,
                  )
                : null,
          ),
        );
      }
    }

    // Add selected accounts (with visit report location priority, same as web)
    if (_selectedAccounts.isNotEmpty) {
      await _handleAddAccounts(waypoints, allVisitReports);
    } else if (waypoints.isNotEmpty) {
      widget.onSelect(waypoints);
      if (mounted) {
        Navigator.pop(context);
      }
    }
  }

  Future<void> _handleAddAccounts(
    List<Waypoint> waypoints,
    List<VisitReport> allVisitReports,
  ) async {
    final accountsState = ref.read(accountListProvider);
    final accounts = accountsState.accounts
        .where((acc) => _selectedAccounts.contains(acc.id))
        .toList();

    if (accounts.isEmpty) {
      if (waypoints.isNotEmpty) {
        widget.onSelect(waypoints);
        Navigator.pop(context);
      }
      return;
    }

    // Use visit reports passed from _handleConfirm (already loaded, same as web)
    final visitReports = allVisitReports;

    // Separate accounts into those with visit report locations and those needing geocoding
    final accountsWithVisitLocation = <String, Map<String, dynamic>>{};
    final accountsNeedingGeocoding = <dynamic>[];

    for (final account in accounts) {
      // First, try to find visit report with location for this account
      final accountVisitReports = visitReports.where((vr) => vr.accountId == account.id).toList();
      final accountVisitReport = accountVisitReports.isNotEmpty ? accountVisitReports.first : null;

      if (accountVisitReport != null) {
        final location = accountVisitReport.checkInLocation ?? accountVisitReport.checkOutLocation;
        if (location != null) {
          // Use location from visit report (more accurate, no geocoding needed)
          accountsWithVisitLocation[account.id] = {
            'account': account,
            'location': location,
            'visitReport': accountVisitReport,
          };
          continue;
        }
      }

      // No visit report location found, need geocoding
      final hasAddress = account.address != null && account.address!.isNotEmpty;
      final hasCityOrProvince = account.city != null || account.province != null;
      
      if (hasAddress || hasCityOrProvince) {
        accountsNeedingGeocoding.add(account);
      } else {
        // No address info at all
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('${account.name} has no address information and no visit report location'),
              backgroundColor: Colors.orange,
              duration: const Duration(seconds: 3),
            ),
          );
        }
      }
    }

    // Add accounts with visit report locations immediately (no geocoding needed)
    for (final entry in accountsWithVisitLocation.entries) {
      final data = entry.value;
      final account = data['account'];
      final location = data['location'];
      final visitReport = data['visitReport'];

      waypoints.add(
        Waypoint(
          lat: location.latitude,
          lng: location.longitude,
          address: location.address ?? account.address ?? account.name,
          accountId: account.id,
          accountName: account.name,
          visitReportId: visitReport.id,
          account: AccountInfo(
            id: account.id,
            name: account.name,
          ),
        ),
      );
    }

    // If no accounts need geocoding, return early
    if (accountsNeedingGeocoding.isEmpty) {
      if (mounted) {
        if (waypoints.isNotEmpty) {
          widget.onSelect(waypoints);
          Navigator.pop(context);
        } else {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('No waypoints with valid locations found'),
            ),
          );
        }
      }
      return;
    }

    // Show loading dialog only if geocoding is needed
    if (!mounted) return;
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const CircularProgressIndicator(),
            const SizedBox(height: 16),
            Text('Geocoding ${accountsNeedingGeocoding.length} address${accountsNeedingGeocoding.length != 1 ? 'es' : ''}...'),
          ],
        ),
      ),
    );

    try {
      // Optimize: Process geocoding in parallel batches (max 3 concurrent)
      const batchSize = 3;
      for (var i = 0; i < accountsNeedingGeocoding.length; i += batchSize) {
        final batch = accountsNeedingGeocoding.skip(i).take(batchSize).toList();
        
        // Process batch in parallel
        final results = await Future.wait(
          batch.map((account) async {
            // Build full address with Indonesia for better geocoding accuracy
            final addressParts = <String>[];
            if (account.address != null) addressParts.add(account.address!);
            if (account.city != null) addressParts.add(account.city!);
            if (account.province != null) addressParts.add(account.province!);
            addressParts.add('Indonesia'); // Add Indonesia for better results
            final fullAddress = addressParts.join(', ');

              try {
                final coords = await _geocodeAddress(fullAddress);
                if (coords == null) {
                  throw Exception('No geocoding results found');
                }
                return {
                  'account': account,
                  'coords': coords,
                  'address': fullAddress,
                  'error': null,
                };
              } on DioException catch (e) {
                // Handle DioException specifically (same as web)
                String errorMessage = 'Unknown error';
                if (e.response != null) {
                  final errorData = e.response!.data;
                  if (errorData is Map<String, dynamic> && 
                      errorData['error'] != null &&
                      errorData['error']['message'] != null) {
                    errorMessage = errorData['error']['message'] as String;
                  } else {
                    errorMessage = 'Geocoding failed: ${e.response?.statusCode} ${e.response?.statusMessage}';
                  }
                } else if (e.message != null) {
                  errorMessage = e.message!;
                }
                return {
                  'account': account,
                  'coords': null,
                  'address': fullAddress,
                  'error': errorMessage,
                };
              } catch (e) {
                return {
                  'account': account,
                  'coords': null,
                  'address': fullAddress,
                  'error': e.toString(),
                };
              }
          }),
        );

        // Process results
        for (final result in results) {
          final account = result['account'] as dynamic;
          final coords = result['coords'] as Map<String, double>?;
          final fullAddress = result['address'] as String;
          final error = result['error'] as String?;

          if (coords != null) {
            waypoints.add(
              Waypoint(
                lat: coords['lat']!,
                lng: coords['lng']!,
                address: fullAddress,
                accountId: account.id,
                accountName: account.name,
                account: AccountInfo(
                  id: account.id,
                  name: account.name,
                ),
              ),
            );
          } else {
            // Show error for failed geocoding (same message format as web)
            if (mounted) {
              final errorMessage = error ?? 'Unknown error';
              final isNoResults = errorMessage.contains('No geocoding results found') ||
                  errorMessage.contains('GEOCODING_NO_RESULTS');
              
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(
                    isNoResults
                        ? 'Address not found for "${account.name}". Please verify the address or use Visit Reports tab to select from visit reports with GPS coordinates.'
                        : 'Failed to geocode "${account.name}": $errorMessage',
                  ),
                  backgroundColor: Colors.orange,
                  duration: const Duration(seconds: 4),
                ),
              );
            }
          }
        }
      }
    } finally {
      if (mounted) {
        Navigator.pop(context); // Close loading dialog
      }
    }

    if (mounted) {
      if (waypoints.isNotEmpty) {
        widget.onSelect(waypoints);
        Navigator.pop(context); // Close selector dialog
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('No waypoints with valid locations found'),
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;

    return Dialog(
      child: Container(
        constraints: const BoxConstraints(maxHeight: 600, maxWidth: 500),
        child: Column(
          children: [
            // Header
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                children: [
                  Text(
                    l10n.selectWaypoint,
                    style: theme.textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const Spacer(),
                  IconButton(
                    icon: const Icon(Icons.close),
                    onPressed: () => Navigator.pop(context),
                  ),
                ],
              ),
            ),
            const Divider(),

            // Tabs
            TabBar(
              controller: _tabController,
              isScrollable: false,
              tabAlignment: TabAlignment.fill,
              labelPadding: const EdgeInsets.symmetric(horizontal: 4),
              tabs: [
                Tab(
                  child: FittedBox(
                    fit: BoxFit.scaleDown,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(Icons.business, size: 16),
                        const SizedBox(width: 4),
                        Text(
                          '${l10n.selectFromAccounts} (${_selectedAccounts.length})',
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ),
                ),
                Tab(
                  child: FittedBox(
                    fit: BoxFit.scaleDown,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(Icons.calendar_today, size: 16),
                        const SizedBox(width: 4),
                        Text(
                          '${l10n.selectFromVisitReports} (${_selectedVisitReports.length})',
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),

            // Content
            Expanded(
              child: TabBarView(
                controller: _tabController,
                children: [
                  _buildAccountsTab(context, theme),
                  _buildVisitReportsTab(context, theme),
                ],
              ),
            ),

            // Footer
            const Divider(),
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  TextButton(
                    onPressed: () => Navigator.pop(context),
                    child: Text(l10n.cancel),
                  ),
                  const SizedBox(width: 8),
                  ElevatedButton(
                    onPressed: (_selectedAccounts.isNotEmpty ||
                            _selectedVisitReports.isNotEmpty)
                        ? _handleConfirm
                        : null,
                    child: Text(
                      '${l10n.add} ${_selectedAccounts.length + _selectedVisitReports.length} ${l10n.waypoints}',
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAccountsTab(BuildContext context, ThemeData theme) {
    final l10n = AppLocalizations.of(context)!;
    final accountsState = ref.watch(accountListProvider);

    // Load accounts only on initial state (no search active, never loaded before)
    if (accountsState.accounts.isEmpty &&
        !accountsState.isLoading &&
        accountsState.searchQuery.isEmpty &&
        accountsState.pagination == null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        ref.read(accountListProvider.notifier).loadAccounts();
      });
    }

    // Ensure visit reports are loaded for location indicators (same as web)
    final visitReportsState = ref.watch(visitReportListProvider);
    if (visitReportsState.visitReports.isEmpty &&
        !visitReportsState.isLoading &&
        visitReportsState.searchQuery.isEmpty &&
        visitReportsState.pagination == null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        ref.read(visitReportListProvider.notifier).loadVisitReports(
          forRouteOptimization: true,
          status: 'approved',
        );
      });
    }

    return Column(
      children: [
        // Search
        Padding(
          padding: const EdgeInsets.all(16.0),
          child: TextField(
            controller: _searchAccountsController,
            decoration: InputDecoration(
              hintText: l10n.searchAccountsOrContacts,
              prefixIcon: const Icon(Icons.search),
              suffixIcon: _searchAccountsController.text.isNotEmpty
                  ? IconButton(
                      icon: const Icon(Icons.clear),
                      onPressed: () {
                        _searchAccountsController.clear();
                        _accountSearchDebounce?.cancel();
                        ref.read(accountListProvider.notifier).loadAccounts(
                              search: '',
                              refresh: true,
                            );
                      },
                    )
                  : null,
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              filled: true,
              fillColor: theme.colorScheme.surfaceContainerHighest,
            ),
            onChanged: (value) {
              setState(() {}); // Update suffixIcon visibility
              _accountSearchDebounce?.cancel();
              _accountSearchDebounce = Timer(
                const Duration(milliseconds: 500),
                () {
                  ref.read(accountListProvider.notifier).loadAccounts(
                        search: value,
                        refresh: true,
                      );
                },
              );
            },
          ),
        ),

        // Accounts List
        Expanded(
          child: accountsState.isLoading
              ? const LoadingWidget()
              : accountsState.accounts.isEmpty
                  ? Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.business_outlined,
                            size: 48,
                            color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
                          ),
                          const SizedBox(height: 16),
                          Text(
                            l10n.noAccountsFoundForWaypoint,
                            style: theme.textTheme.bodyMedium?.copyWith(
                              color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
                            ),
                          ),
                        ],
                      ),
                    )
                  : Builder(
                      builder: (context) {
                        // Get visit reports for location check
                        final visitReportsState = ref.watch(visitReportListProvider);
                        final visitReportsWithLocation = visitReportsState.visitReports
                            .where((vr) => 
                                vr.checkInLocation != null || vr.checkOutLocation != null)
                            .toList();

                        return ListView.builder(
                          padding: const EdgeInsets.symmetric(horizontal: 16),
                          itemCount: accountsState.accounts.length,
                          itemBuilder: (context, index) {
                            final account = accountsState.accounts[index];
                            final isSelected = _selectedAccounts.contains(account.id);
                            final hasAddress = account.address != null &&
                                account.address!.isNotEmpty;
                            final hasCityOrProvince = account.city != null || account.province != null;
                            
                            // Check if account has visit report with location
                            final hasVisitReportLocation = visitReportsWithLocation
                                .any((vr) => vr.accountId == account.id);
                            
                            final canSelect = hasAddress || hasCityOrProvince || hasVisitReportLocation;

                            return Card(
                              margin: const EdgeInsets.only(bottom: 8),
                              child: CheckboxListTile(
                                value: isSelected,
                                onChanged: canSelect
                                    ? (value) {
                                        setState(() {
                                          if (value == true) {
                                            _selectedAccounts.add(account.id);
                                          } else {
                                            _selectedAccounts.remove(account.id);
                                          }
                                        });
                                      }
                                    : null,
                                title: Row(
                                  children: [
                                    Expanded(
                                      child: Text(
                                        account.name,
                                        style: theme.textTheme.bodyMedium?.copyWith(
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                    ),
                                    if (hasVisitReportLocation)
                                      Container(
                                        padding: const EdgeInsets.symmetric(
                                          horizontal: 8,
                                          vertical: 4,
                                        ),
                                        decoration: BoxDecoration(
                                          color: Colors.green.shade50,
                                          borderRadius: BorderRadius.circular(12),
                                          border: Border.all(
                                            color: Colors.green.shade300,
                                            width: 1,
                                          ),
                                        ),
                                        child: Row(
                                          mainAxisSize: MainAxisSize.min,
                                          children: [
                                            Icon(
                                              Icons.location_on,
                                              size: 12,
                                              color: Colors.green.shade700,
                                            ),
                                            const SizedBox(width: 4),
                                            Text(
                                              'GPS',
                                              style: TextStyle(
                                                fontSize: 10,
                                                fontWeight: FontWeight.bold,
                                                color: Colors.green.shade700,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                  ],
                                ),
                                subtitle: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    if (hasVisitReportLocation)
                                      Padding(
                                        padding: const EdgeInsets.only(bottom: 4),
                                        child: Text(
                                          'Location available from visit report',
                                          style: theme.textTheme.bodySmall?.copyWith(
                                            color: Colors.green.shade700,
                                            fontWeight: FontWeight.w500,
                                          ),
                                        ),
                                      ),
                                    if (account.address != null)
                                      Text(
                                        account.address!,
                                        style: theme.textTheme.bodySmall,
                                      ),
                                    if (account.city != null || account.province != null)
                                      Text(
                                        '${account.city ?? ''}${account.city != null && account.province != null ? ', ' : ''}${account.province ?? ''}',
                                        style: theme.textTheme.bodySmall?.copyWith(
                                          color: theme.colorScheme.onSurface
                                              .withValues(alpha: 0.6),
                                        ),
                                      ),
                                    if (!canSelect)
                                      Padding(
                                        padding: const EdgeInsets.only(top: 4),
                                        child: Text(
                                          'No address or visit report location available',
                                          style: theme.textTheme.bodySmall?.copyWith(
                                            color: theme.colorScheme.error,
                                            fontStyle: FontStyle.italic,
                                          ),
                                        ),
                                      ),
                                  ],
                                ),
                                secondary: const Icon(Icons.business),
                              ),
                            );
                          },
                        );
                      },
                    ),
        ),
      ],
    );
  }

  Widget _buildVisitReportsTab(BuildContext context, ThemeData theme) {
    final l10n = AppLocalizations.of(context)!;
    final visitReportsState = ref.watch(visitReportListProvider);

    // Load visit reports only on initial state (no search active, never loaded before)
    if (visitReportsState.visitReports.isEmpty &&
        !visitReportsState.isLoading &&
        visitReportsState.searchQuery.isEmpty &&
        visitReportsState.pagination == null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        ref.read(visitReportListProvider.notifier).loadVisitReports(
          forRouteOptimization: true,
          status: 'approved', // Only show approved visits for route planning
        );
      });
    }

    // Filter visit reports with location data and approved status
    final visitReportsWithLocation = visitReportsState.visitReports
        .where((vr) =>
            (vr.checkInLocation != null || vr.checkOutLocation != null) &&
            vr.status == 'approved')
        .toList();

    return Column(
      children: [
        // Search
        Padding(
          padding: const EdgeInsets.all(16.0),
          child: TextField(
            controller: _searchVisitsController,
            decoration: InputDecoration(
              hintText: l10n.searchVisitReports,
              prefixIcon: const Icon(Icons.search),
              suffixIcon: _searchVisitsController.text.isNotEmpty
                  ? IconButton(
                      icon: const Icon(Icons.clear),
                      onPressed: () {
                        _searchVisitsController.clear();
                        _visitSearchDebounce?.cancel();
                        ref.read(visitReportListProvider.notifier).loadVisitReports(
                              search: '',
                              refresh: true,
                              forRouteOptimization: true,
                              status: 'approved',
                            );
                      },
                    )
                  : null,
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              filled: true,
              fillColor: theme.colorScheme.surfaceContainerHighest,
            ),
            onChanged: (value) {
              setState(() {}); // Update suffixIcon visibility
              _visitSearchDebounce?.cancel();
              _visitSearchDebounce = Timer(
                const Duration(milliseconds: 500),
                () {
                  ref.read(visitReportListProvider.notifier).loadVisitReports(
                        search: value,
                        refresh: true,
                        forRouteOptimization: true,
                        status: 'approved',
                      );
                },
              );
            },
          ),
        ),

        // Visit Reports List
        Expanded(
          child: visitReportsState.isLoading
              ? const LoadingWidget()
              : visitReportsWithLocation.isEmpty
                  ? Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.calendar_today_outlined,
                            size: 48,
                            color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
                          ),
                          const SizedBox(height: 16),
                          Text(
                            'No visit reports with location found',
                            style: theme.textTheme.bodyMedium?.copyWith(
                              color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
                            ),
                          ),
                        ],
                      ),
                    )
                  : ListView.builder(
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      itemCount: visitReportsWithLocation.length,
                      itemBuilder: (context, index) {
                        final visitReport = visitReportsWithLocation[index];
                        final isSelected =
                            _selectedVisitReports.contains(visitReport.id);
                        final location = visitReport.checkInLocation ??
                            visitReport.checkOutLocation;

                        return Card(
                          margin: const EdgeInsets.only(bottom: 8),
                          child: CheckboxListTile(
                            value: isSelected,
                            onChanged: (value) {
                              setState(() {
                                if (value == true) {
                                  _selectedVisitReports.add(visitReport.id);
                                } else {
                                  _selectedVisitReports.remove(visitReport.id);
                                }
                              });
                            },
                            title: Text(
                              visitReport.account?.name ?? 'Visit Report',
                              style: theme.textTheme.bodyMedium?.copyWith(
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                            subtitle: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                if (location != null && location.address != null)
                                  Text(
                                    location.address!,
                                    style: theme.textTheme.bodySmall,
                                  ),
                                if (location != null)
                                  Text(
                                    '${location.latitude.toStringAsFixed(6)}, ${location.longitude.toStringAsFixed(6)}',
                                    style: theme.textTheme.bodySmall?.copyWith(
                                      color: theme.colorScheme.onSurface
                                          .withValues(alpha: 0.6),
                                    ),
                                  ),
                                Text(
                                  visitReport.visitDate,
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: theme.colorScheme.onSurface
                                        .withValues(alpha: 0.6),
                                  ),
                                ),
                              ],
                            ),
                            secondary: const Icon(Icons.calendar_today),
                          ),
                        );
                      },
                    ),
        ),
      ],
    );
  }
}

