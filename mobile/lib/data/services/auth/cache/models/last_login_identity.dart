class LastLoginIdentity {
  final String name;
  final String identifier;

  LastLoginIdentity({
    required this.name,
    required this.identifier,
  });

  Map<String, dynamic> toMap() {
    return {
      'name': name,
      'identifier': identifier,
    };
  }

  factory LastLoginIdentity.fromMap(Map<String, dynamic> map) {
    return LastLoginIdentity(
      name: map['name'] as String,
      identifier: map['identifier'] as String,
    );
  }
}
