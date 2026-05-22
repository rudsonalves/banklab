import '../enums/contact_verification_channel.dart';

class ContactVerificationRequestResponseDto {
  final String verificationId;
  final ContactVerificationChannel channel;
  final String target;
  final DateTime expiresAt;

  ContactVerificationRequestResponseDto({
    required this.verificationId,
    required this.channel,
    required this.target,
    required this.expiresAt,
  });

  factory ContactVerificationRequestResponseDto.fromMap(
    Map<String, dynamic> map,
  ) {
    return ContactVerificationRequestResponseDto(
      verificationId: map['verification_id'] as String,
      channel: ContactVerificationChannel.fromString(map['channel'] as String),
      target: map['target'] as String,
      expiresAt: DateTime.parse(map['expires_at'] as String),
    );
  }
}
