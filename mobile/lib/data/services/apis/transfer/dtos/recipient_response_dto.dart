import 'recipient_info_dto.dart';

class RecipientResponseDto {
  final List<RecipientInfoDto> accounts;

  RecipientResponseDto(
    this.accounts,
  );

  factory RecipientResponseDto.fromMap(Map<String, dynamic> map) {
    final accounts = map['accounts'] as List<dynamic>;

    return RecipientResponseDto(
      accounts
          .map(
            (account) => RecipientInfoDto.fromMap(
              account as Map<String, dynamic>,
            ),
          )
          .toList(),
    );
  }
}
