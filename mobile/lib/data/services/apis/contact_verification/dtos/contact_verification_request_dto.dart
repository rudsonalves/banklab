class ContactVerificationRequestDto {
  final String channel;
  final String target;

  ContactVerificationRequestDto({
    required this.channel,
    required this.target,
  });

  Map<String, dynamic> toMap() {
    return {
      'channel': channel,
      'target': target,
    };
  }

  factory ContactVerificationRequestDto.fromMap(Map<String, dynamic> map) {
    return ContactVerificationRequestDto(
      channel: map['channel'] as String,
      target: map['target'] as String,
    );
  }
}
