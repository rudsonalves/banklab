import '/core/extensions/datetime_extension.dart';
import '/domain/enums/user_role.dart';

class UserProfile {
  final String userId;
  final String customerId;
  final String name;
  final String email;
  final UserRole role;
  final DateTime createdAt;
  final DateTime updatedAt;

  UserProfile({
    required this.userId,
    required this.name,
    required this.email,
    required this.role,
    required this.customerId,
    required this.createdAt,
    required this.updatedAt,
  });

  factory UserProfile.fromMap(Map<String, dynamic> map) {
    return UserProfile(
      userId: map['id'] as String,
      name: map['name'] as String,
      email: map['email'] as String,
      role: UserRole.byName(map['role'] as String),
      customerId: map['customer_id'] as String,
      createdAt: DateTimeExtensions.parseOrNull(map['created_at'] as String)!,
      updatedAt: DateTimeExtensions.parseOrNull(map['updated_at'] as String)!,
    );
  }
}
