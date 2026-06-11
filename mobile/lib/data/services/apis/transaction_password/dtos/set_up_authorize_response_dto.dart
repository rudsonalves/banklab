class SetUpAuthorizeResponseDto {
  final String stepUpToken;
  final int expiresIn;

  SetUpAuthorizeResponseDto({
    required this.stepUpToken,
    required this.expiresIn,
  });

  factory SetUpAuthorizeResponseDto.fromApi(Map<String, dynamic> map) {
    return SetUpAuthorizeResponseDto(
      stepUpToken: map['step_up_token'] as String,
      expiresIn: map['expires_in'] as int,
    );
  }
}
