class ContactVerificationRequestResponseDto {
  final String verificationId;
  final String channel;
  final String target;
  final String token;
  final DateTime expiresAt;

  ContactVerificationRequestResponseDto({
    required this.verificationId,
    required this.channel,
    required this.target,
    required this.token,
    required this.expiresAt,
  });

  factory ContactVerificationRequestResponseDto.fromMap(
    Map<String, dynamic> map,
  ) {
    return ContactVerificationRequestResponseDto(
      verificationId: map['verification_id'] as String,
      channel: map['channel'] as String,
      target: map['target'] as String,
      token: map['token'] as String,
      expiresAt: DateTime.parse(map['expires_at'] as String),
    );
  }
}
