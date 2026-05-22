class CpfCheckResponseDto {
  final String cpf;
  final bool exists;
  final bool available;

  CpfCheckResponseDto({
    required this.cpf,
    required this.exists,
    required this.available,
  });

  factory CpfCheckResponseDto.fromMap(Map<String, dynamic> json) {
    return CpfCheckResponseDto(
      cpf: json['cpf'] as String,
      exists: json['exists'] as bool,
      available: json['available'] as bool,
    );
  }
}
