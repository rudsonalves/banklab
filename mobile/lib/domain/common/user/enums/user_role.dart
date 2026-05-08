enum UserRole {
  // customer is the default role for all users.
  customer,
  // admin users have elevated permissions to manage the system,
  // but are not customers themselves.
  admin,
  // none represents an unknown or uninitialized role. It can be used
  // as a default value when parsing or when a user has no assigned role.
  none
  ;

  factory UserRole.byName(String name) => UserRole.values.firstWhere(
    (e) => e.name == name,
    orElse: () => UserRole.none,
  );
}
