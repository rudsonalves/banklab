import 'internal_transfer_recipient_dto.dart';

class InternalTransferRecipientLookupResponseDto {
  final List<InternalTransferRecipientDto> accounts;

  const InternalTransferRecipientLookupResponseDto({required this.accounts});

  factory InternalTransferRecipientLookupResponseDto.fromMap(
    Map<String, dynamic> map,
  ) {
    final accounts = map['accounts'] as List<dynamic>;

    return InternalTransferRecipientLookupResponseDto(
      accounts: accounts
          .map(
            (account) => InternalTransferRecipientDto.fromMap(
              account as Map<String, dynamic>,
            ),
          )
          .toList(),
    );
  }
}
