class ProfileResponse {
  const ProfileResponse({
    required this.user,
    this.stats,
    this.activities,
    this.transactions,
  });

  final ProfileUser user;
  final ProfileStats? stats;
  final List<ProfileActivity>? activities;
  final List<ProfileTransaction>? transactions;

  factory ProfileResponse.fromJson(Map<String, dynamic> json) {
    // Safely extract user field
    final userData = json['user'];
    if (userData == null) {
      throw Exception('Invalid profile response: user field is required');
    }
    
    Map<String, dynamic> userMap;
    if (userData is Map<String, dynamic>) {
      userMap = userData;
    } else {
      throw Exception('Invalid profile response: user field must be an object');
    }
    
    // Safely extract stats
    ProfileStats? stats;
    final statsData = json['stats'];
    if (statsData != null && statsData is Map<String, dynamic>) {
      stats = ProfileStats.fromJson(statsData);
    }
    
    // Safely extract activities
    List<ProfileActivity>? activities;
    final activitiesData = json['activities'];
    if (activitiesData != null && activitiesData is List) {
      activities = activitiesData
          .whereType<Map<String, dynamic>>()
          .map((e) => ProfileActivity.fromJson(e))
          .toList();
    }
    
    // Safely extract transactions
    List<ProfileTransaction>? transactions;
    final transactionsData = json['transactions'];
    if (transactionsData != null && transactionsData is List) {
      transactions = transactionsData
          .whereType<Map<String, dynamic>>()
          .map((e) => ProfileTransaction.fromJson(e))
          .toList();
    }
    
    return ProfileResponse(
      user: ProfileUser.fromJson(userMap),
      stats: stats,
      activities: activities,
      transactions: transactions,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'user': user.toJson(),
      'stats': stats?.toJson(),
      'activities': activities?.map((e) => e.toJson()).toList(),
      'transactions': transactions?.map((e) => e.toJson()).toList(),
    };
  }
}

class ProfileUser {
  const ProfileUser({
    required this.id,
    required this.email,
    required this.name,
    this.avatarUrl,
    this.roleId,
    this.role,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
  });

  final String id;
  final String email;
  final String name;
  final String? avatarUrl;
  final String? roleId;
  final ProfileRole? role;
  final String status;
  final DateTime createdAt;
  final DateTime updatedAt;

  factory ProfileUser.fromJson(Map<String, dynamic> json) {
    // Safely extract required fields
    final id = json['id'];
    final email = json['email'];
    final name = json['name'];
    final createdAtStr = json['created_at'];
    final updatedAtStr = json['updated_at'];
    
    if (id == null || email == null || name == null || createdAtStr == null || updatedAtStr == null) {
      throw Exception('Invalid user data: required fields missing');
    }
    
    // Safely extract optional fields
    final avatarUrl = json['avatar_url'];
    final roleId = json['role_id'];
    final status = json['status'];
    
    // Safely extract role
    ProfileRole? role;
    final roleData = json['role'];
    if (roleData != null && roleData is Map<String, dynamic>) {
      role = ProfileRole.fromJson(roleData);
    }
    
    // Parse dates safely
    DateTime createdAt;
    DateTime updatedAt;
    try {
      createdAt = DateTime.parse(createdAtStr.toString());
    } catch (e) {
      throw Exception('Invalid created_at date format: $createdAtStr');
    }
    
    try {
      updatedAt = DateTime.parse(updatedAtStr.toString());
    } catch (e) {
      throw Exception('Invalid updated_at date format: $updatedAtStr');
    }
    
    return ProfileUser(
      id: id.toString(),
      email: email.toString(),
      name: name.toString(),
      avatarUrl: avatarUrl?.toString(),
      roleId: roleId?.toString(),
      role: role,
      status: status?.toString() ?? 'active',
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'name': name,
      'avatar_url': avatarUrl,
      'role_id': roleId,
      'role': role?.toJson(),
      'status': status,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class ProfileRole {
  const ProfileRole({
    required this.id,
    required this.name,
    required this.code,
  });

  final String id;
  final String name;
  final String code;

  factory ProfileRole.fromJson(Map<String, dynamic> json) {
    final id = json['id'];
    final name = json['name'];
    final code = json['code'];
    
    if (id == null || name == null || code == null) {
      throw Exception('Invalid role data: required fields missing');
    }
    
    return ProfileRole(
      id: id.toString(),
      name: name.toString(),
      code: code.toString(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'code': code,
    };
  }
}

class ProfileStats {
  const ProfileStats({
    this.visits,
    this.deals,
    this.tasks,
  });

  final int? visits;
  final int? deals;
  final int? tasks;

  factory ProfileStats.fromJson(Map<String, dynamic> json) {
    // Safely extract int values
    int? visits;
    int? deals;
    int? tasks;
    
    final visitsData = json['visits'];
    if (visitsData != null) {
      if (visitsData is int) {
        visits = visitsData;
      } else if (visitsData is num) {
        visits = visitsData.toInt();
      }
    }
    
    final dealsData = json['deals'];
    if (dealsData != null) {
      if (dealsData is int) {
        deals = dealsData;
      } else if (dealsData is num) {
        deals = dealsData.toInt();
      }
    }
    
    final tasksData = json['tasks'];
    if (tasksData != null) {
      if (tasksData is int) {
        tasks = tasksData;
      } else if (tasksData is num) {
        tasks = tasksData.toInt();
      }
    }
    
    return ProfileStats(
      visits: visits,
      deals: deals,
      tasks: tasks,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'visits': visits,
      'deals': deals,
      'tasks': tasks,
    };
  }
}

class ProfileActivity {
  const ProfileActivity({
    required this.id,
    required this.title,
    required this.description,
    required this.type,
    required this.date,
  });

  final String id;
  final String title;
  final String description;
  final String type;
  final DateTime date;

  factory ProfileActivity.fromJson(Map<String, dynamic> json) {
    final id = json['id'];
    final title = json['title'];
    final description = json['description'];
    final type = json['type'];
    final dateData = json['date'];
    
    if (id == null || title == null || description == null || type == null || dateData == null) {
      throw Exception('Invalid activity data: required fields missing');
    }
    
    DateTime date;
    try {
      date = DateTime.parse(dateData.toString());
    } catch (e) {
      throw Exception('Invalid activity date format: $dateData');
    }
    
    return ProfileActivity(
      id: id.toString(),
      title: title.toString(),
      description: description.toString(),
      type: type.toString(),
      date: date,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'type': type,
      'date': date.toIso8601String(),
    };
  }
}

class ProfileTransaction {
  const ProfileTransaction({
    required this.id,
    required this.product,
    required this.status,
    required this.date,
    required this.amount,
  });

  final String id;
  final String product;
  final String status;
  final DateTime date;
  final int amount;

  factory ProfileTransaction.fromJson(Map<String, dynamic> json) {
    final id = json['id'];
    final product = json['product'];
    final status = json['status'];
    final dateData = json['date'];
    final amountData = json['amount'];
    
    if (id == null || product == null || status == null || dateData == null || amountData == null) {
      throw Exception('Invalid transaction data: required fields missing');
    }
    
    DateTime date;
    try {
      date = DateTime.parse(dateData.toString());
    } catch (e) {
      throw Exception('Invalid transaction date format: $dateData');
    }
    
    int amount;
    if (amountData is int) {
      amount = amountData;
    } else if (amountData is num) {
      amount = amountData.toInt();
    } else {
      throw Exception('Invalid transaction amount format: $amountData');
    }
    
    return ProfileTransaction(
      id: id.toString(),
      product: product.toString(),
      status: status.toString(),
      date: date,
      amount: amount,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'product': product,
      'status': status,
      'date': date.toIso8601String(),
      'amount': amount,
    };
  }
}

class UpdateProfileRequest {
  const UpdateProfileRequest({
    required this.name,
  });

  final String name;

  Map<String, dynamic> toJson() {
    return {
      'name': name,
    };
  }
}

class ChangePasswordRequest {
  const ChangePasswordRequest({
    required this.currentPassword,
    required this.password,
    required this.confirmPassword,
  });

  final String currentPassword;
  final String password;
  final String confirmPassword;

  Map<String, dynamic> toJson() {
    return {
      'current_password': currentPassword,
      'password': password,
      'confirm_password': confirmPassword,
    };
  }
}

