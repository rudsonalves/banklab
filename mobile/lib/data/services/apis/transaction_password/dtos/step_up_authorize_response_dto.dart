class StepUpAuthorizeResponseDto {
  final String stepUpToken;
  final int expiresIn;

  StepUpAuthorizeResponseDto({
    required this.stepUpToken,
    required this.expiresIn,
  });

  factory StepUpAuthorizeResponseDto.fromApi(Map<String, dynamic> map) {
    return StepUpAuthorizeResponseDto(
      stepUpToken: map['step_up_token'] as String,
      expiresIn: map['expires_in'] as int,
    );
  }
}
