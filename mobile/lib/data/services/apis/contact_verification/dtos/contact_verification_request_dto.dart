import '../enums/contact_verification_channel.dart';

class ContactVerificationRequestDto {
  final ContactVerificationChannel channel;
  final String target;

  ContactVerificationRequestDto({
    required this.channel,
    required this.target,
  });

  Map<String, dynamic> toMap() {
    return {
      'channel': channel.name,
      'target': target,
    };
  }

  factory ContactVerificationRequestDto.fromMap(Map<String, dynamic> map) {
    return ContactVerificationRequestDto(
      channel: ContactVerificationChannel.fromString(map['channel']),
      target: map['target'] as String,
    );
  }
}
