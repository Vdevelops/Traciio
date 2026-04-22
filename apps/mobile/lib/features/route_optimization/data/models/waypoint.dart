class Waypoint {
  final int? order;
  final double lat;
  final double lng;
  final String? address;
  final String? accountId;
  final String? accountName;
  final String? contactId;
  final String? contactName;
  final String? visitReportId;
  final AccountInfo? account;
  // Time window constraints (optional)
  final DateTime? earliestArrival; // Earliest time customer is available
  final DateTime? latestArrival; // Latest time customer is available
  final int? serviceDuration; // Service duration in minutes
  final int? priority; // Priority: 1 (highest) to 5 (lowest), default: 3

  Waypoint({
    this.order,
    required this.lat,
    required this.lng,
    this.address,
    this.accountId,
    this.accountName,
    this.contactId,
    this.contactName,
    this.visitReportId,
    this.account,
    this.earliestArrival,
    this.latestArrival,
    this.serviceDuration,
    this.priority,
  });

  factory Waypoint.fromJson(Map<String, dynamic> json) {
    return Waypoint(
      order: json['order'] as int?,
      lat: (json['lat'] as num).toDouble(),
      lng: (json['lng'] as num).toDouble(),
      address: json['address'] as String?,
      accountId: json['account_id'] as String?,
      accountName: json['account_name'] as String?,
      contactId: json['contact_id'] as String?,
      contactName: json['contact_name'] as String?,
      visitReportId: json['visit_report_id'] as String?,
      account: json['account'] != null
          ? AccountInfo.fromJson(json['account'] as Map<String, dynamic>)
          : null,
      earliestArrival: json['earliest_arrival'] != null
          ? DateTime.parse(json['earliest_arrival'] as String)
          : null,
      latestArrival: json['latest_arrival'] != null
          ? DateTime.parse(json['latest_arrival'] as String)
          : null,
      serviceDuration: json['service_duration'] as int?,
      priority: json['priority'] as int?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      if (order != null) 'order': order,
      'lat': lat,
      'lng': lng,
      if (address != null) 'address': address,
      if (accountId != null) 'account_id': accountId,
      if (accountName != null) 'account_name': accountName,
      if (contactId != null) 'contact_id': contactId,
      if (contactName != null) 'contact_name': contactName,
      if (visitReportId != null) 'visit_report_id': visitReportId,
      if (account != null) 'account': account!.toJson(),
      if (earliestArrival != null)
        'earliest_arrival': earliestArrival!.toIso8601String(),
      if (latestArrival != null)
        'latest_arrival': latestArrival!.toIso8601String(),
      if (serviceDuration != null) 'service_duration': serviceDuration,
      if (priority != null) 'priority': priority,
    };
  }
}

class AccountInfo {
  final String id;
  final String name;

  AccountInfo({required this.id, required this.name});

  factory AccountInfo.fromJson(Map<String, dynamic> json) {
    return AccountInfo(id: json['id'] as String, name: json['name'] as String);
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name};
  }
}
