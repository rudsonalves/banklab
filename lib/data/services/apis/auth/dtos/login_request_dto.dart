class LoginRequestDto {
  final String email;
  final String password;

  LoginRequestDto({
    required this.email,
    required this.password,
  });

  Map<String, dynamic> toMap() => {
    'email': email,
    'password': password,
  };

  factory LoginRequestDto.fromMap(Map<String, dynamic> map) {
    return LoginRequestDto(
      email: map['email'] as String,
      password: map['password'] as String,
    );
  }
}
