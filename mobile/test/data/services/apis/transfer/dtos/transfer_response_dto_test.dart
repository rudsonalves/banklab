import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';

Map<String, dynamic> _validPayload() => {
  'from_branch': '0001',
  'from_account_number': '00012345',
  'transaction_reference': '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
  'to_branch': '0001',
  'to_account_number': '00067890',
  'amount': 2500,
  'from_balance': 97500,
  'to_balance': 32500,
};

Money _brl(int cents) =>
    Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

void main() {
  group('TransferResponseDto.fromMap', () {
    test(
      'parses all required fields from a representative success payload',
      () {
        final dto = TransferResponseDto.fromMap(_validPayload());

        expect(dto.fromBranch, '0001');
        expect(dto.fromAccountNumber, '00012345');
        expect(
          dto.transactionReference,
          '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
        );
        expect(dto.toBranch, '0001');
        expect(dto.toAccountNumber, '00067890');
        expect(dto.amount, _brl(2500));
        expect(dto.fromBalance, _brl(97500));
        expect(dto.toBalance, _brl(32500));
      },
    );

    test('parses Money fields as cents via ApiParse.toMoney', () {
      final dto = TransferResponseDto.fromMap(_validPayload());

      expect(dto.amount, _brl(2500));
      expect(dto.fromBalance, _brl(97500));
      expect(dto.toBalance, _brl(32500));
    });

    test('ignores internal account ID fields if present in payload', () {
      // The backend domain model has SourceAccountID / DestinationAccountID
      // (internal UUIDs). Verify they are silently ignored when present.
      final map = _validPayload();
      map['from_account_id'] = 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b';
      map['to_account_id'] = 'ab3b2800-1234-5678-abcd-ef0123456789';

      final dto = TransferResponseDto.fromMap(map);

      // Parsing must succeed and identity is via branch + account_number only.
      expect(dto.fromBranch, '0001');
      expect(dto.toAccountNumber, '00067890');
    });

    test('does not expose internal account IDs as DTO fields', () {
      final dto = TransferResponseDto.fromMap(_validPayload());

      // Statically confirmed: TransferResponseDto has no accountId getters.
      // Dynamic check ensures no such field slips in at runtime.
      expect(() => (dto as dynamic).fromAccountId, throwsNoSuchMethodError);
      expect(() => (dto as dynamic).toAccountId, throwsNoSuchMethodError);
    });

    test('throws when transaction_reference is missing', () {
      final map = _validPayload()..remove('transaction_reference');

      expect(() => TransferResponseDto.fromMap(map), throwsA(anything));
    });

    test('throws when amount is missing', () {
      final map = _validPayload()..remove('amount');

      expect(() => TransferResponseDto.fromMap(map), throwsA(anything));
    });
  });
}
