class InstallationRegistrationResponseDto {
  final String accessToken;
  final String refreshToken;
  final String installationResourceId;
  final String installationStatus;

  const InstallationRegistrationResponseDto({
    required this.accessToken,
    required this.refreshToken,
    required this.installationResourceId,
    required this.installationStatus,
  });

  factory InstallationRegistrationResponseDto.fromMap(
    Map<String, dynamic> map,
  ) {
    return InstallationRegistrationResponseDto(
      accessToken: map['access_token'] as String,
      refreshToken: map['refresh_token'] as String,
      installationResourceId: map['installation_resource_id'] as String,
      installationStatus: map['installation_status'] as String,
    );
  }
}
