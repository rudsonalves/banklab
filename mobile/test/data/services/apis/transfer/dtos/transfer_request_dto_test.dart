import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';

void main() {
  Money _brl(int cents) =>
      Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

  group('TransferRequestDto.toMap', () {
    test('serializes all required fields', () {
      final map = TransferRequestDto(
        fromBranch: '0001',
        fromAccountNumber: '00012345',
        toBranch: '0001',
        toAccountNumber: '00067890',
        amount: _brl(2500),
      ).toMap();

      expect(map['from_branch'], '0001');
      expect(map['from_account_number'], '00012345');
      expect(map['to_branch'], '0001');
      expect(map['to_account_number'], '00067890');
      expect(map['amount'], 2500);
    });

    test('serializes amount as int cents via ApiParse.toInt', () {
      final map = TransferRequestDto(
        fromBranch: '0001',
        fromAccountNumber: '00012345',
        toBranch: '0002',
        toAccountNumber: '00067890',
        amount: _brl(97500),
      ).toMap();

      expect(map['amount'], isA<int>());
      expect(map['amount'], 97500);
    });

    test('omits idempotency_key when null', () {
      final map = TransferRequestDto(
        fromBranch: '0001',
        fromAccountNumber: '00012345',
        toBranch: '0001',
        toAccountNumber: '00067890',
        amount: _brl(1000),
      ).toMap();

      expect(map.containsKey('idempotency_key'), isFalse);
    });

    test('includes idempotency_key when provided', () {
      final map = TransferRequestDto(
        fromBranch: '0001',
        fromAccountNumber: '00012345',
        toBranch: '0001',
        toAccountNumber: '00067890',
        amount: _brl(1000),
        idempotencyKey: 'client-key-abc',
      ).toMap();

      expect(map['idempotency_key'], 'client-key-abc');
    });

    test('does not serialize internal account IDs', () {
      final map = TransferRequestDto(
        fromBranch: '0001',
        fromAccountNumber: '00012345',
        toBranch: '0001',
        toAccountNumber: '00067890',
        amount: _brl(1000),
      ).toMap();

      // Identity must be expressed via branch + account_number, never internal UUIDs.
      final idKeys = map.keys.where((k) => k.endsWith('_id')).toList();
      expect(idKeys, isEmpty);
    });
  });
}
