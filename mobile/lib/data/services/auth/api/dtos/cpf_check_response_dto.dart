class CpfCheckResponseDto {
  final String cpf;
  final bool exists;
  final bool avaliable;

  CpfCheckResponseDto({
    required this.cpf,
    required this.exists,
    required this.avaliable,
  });

  factory CpfCheckResponseDto.fromMap(Map<String, dynamic> json) {
    return CpfCheckResponseDto(
      cpf: json['cpf'] as String,
      exists: json['exists'] as bool,
      avaliable: json['avaliable'] as bool,
    );
  }
}
