class ContactVerificationConfirmRequestDto {
  final String verificationId;
  final String token;

  ContactVerificationConfirmRequestDto({
    required this.verificationId,
    required this.token,
  });

  Map<String, dynamic> toMap() {
    return {
      'verification_id': verificationId,
      'token': token,
    };
  }

  factory ContactVerificationConfirmRequestDto.fromMap(
    Map<String, dynamic> map,
  ) {
    return ContactVerificationConfirmRequestDto(
      verificationId: map['verification_id'] as String,
      token: map['token'] as String,
    );
  }
}
