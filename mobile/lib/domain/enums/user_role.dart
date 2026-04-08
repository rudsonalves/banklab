enum UserRole {
  customer,
  admin,
  none
  ;

  factory UserRole.byName(String name) => UserRole.values.firstWhere(
    (e) => e.name == name,
    orElse: () => UserRole.none,
  );
}
