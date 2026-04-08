class RegisterRequestDto {
  final String name;
  final String email;
  final String password;
  final String cpf;

  RegisterRequestDto({
    required this.name,
    required this.email,
    required this.password,
    required this.cpf,
  });

  Map<String, dynamic> toMap() {
    return {
      'name': name,
      'email': email,
      'password': password,
      'cpf': cpf,
    };
  }

  factory RegisterRequestDto.fromMap(Map<String, dynamic> map) {
    return RegisterRequestDto(
      name: map['name'] as String,
      email: map['email'] as String,
      password: map['password'] as String,
      cpf: map['cpf'] as String,
    );
  }
}
