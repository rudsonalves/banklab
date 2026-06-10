class CustomerSession {
  final String id;
  final String name;
  final String cpf;
  final DateTime birthDate;
  final DateTime createdAt;

  CustomerSession({
    required this.id,
    required this.name,
    required this.cpf,
    required this.birthDate,
    required this.createdAt,
  });

  factory CustomerSession.fromApi(Map<String, dynamic> map) {
    return CustomerSession(
      id: map['id'] as String,
      name: map['name'] as String,
      cpf: map['cpf'] as String,
      birthDate: DateTime.parse(map['birth_date'] as String),
      createdAt: DateTime.parse(map['created_at'] as String),
    );
  }
}
