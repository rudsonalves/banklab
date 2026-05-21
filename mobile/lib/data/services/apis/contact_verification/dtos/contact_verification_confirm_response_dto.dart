import '../enums/contact_verification_channel.dart';

class ContactVerificationConfirmResponseDto {
  final String verificationToken;
  final ContactVerificationChannel channel;
  final String target;
  final DateTime verifiedAt;

  ContactVerificationConfirmResponseDto({
    required this.verificationToken,
    required this.channel,
    required this.target,
    required this.verifiedAt,
  });

  factory ContactVerificationConfirmResponseDto.fromMap(
    Map<String, dynamic> map,
  ) {
    return ContactVerificationConfirmResponseDto(
      verificationToken: map['verification_token'] as String,
      channel: ContactVerificationChannel.fromString(map['channel'] as String),
      target: map['target'] as String,
      verifiedAt: DateTime.parse(map['verified_at'] as String),
    );
  }
}
