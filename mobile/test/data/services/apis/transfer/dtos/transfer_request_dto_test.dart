import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';

void main() {
  Money brl(int cents) =>
      Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

  group('TransferRequestDto.toMap', () {
    test('serializes all required fields', () {
      final map = TransferRequestDto(
        fromAccountId: 'acc-src-001',
        toAccountId: 'acc-dst-001',
        amount: brl(2500),
        idempotencyKey: 'client-key-abc',
      ).toMap();

      expect(map['from_account_id'], 'acc-src-001');
      expect(map['to_account_id'], 'acc-dst-001');
      expect(map['amount'], 2500);
      expect(map['idempotency_key'], 'client-key-abc');
    });

    test('serializes amount as int cents via ApiParse.toInt', () {
      final map = TransferRequestDto(
        fromAccountId: 'acc-src-001',
        toAccountId: 'acc-dst-001',
        amount: brl(97500),
        idempotencyKey: 'client-key-abc',
      ).toMap();

      expect(map['amount'], isA<int>());
      expect(map['amount'], 97500);
    });

    test('includes required idempotency_key', () {
      final map = TransferRequestDto(
        fromAccountId: 'acc-src-001',
        toAccountId: 'acc-dst-001',
        amount: brl(1000),
        idempotencyKey: 'client-key-abc',
      ).toMap();

      expect(map['idempotency_key'], 'client-key-abc');
    });

    test('omits description when null', () {
      final map = TransferRequestDto(
        fromAccountId: 'acc-src-001',
        toAccountId: 'acc-dst-001',
        amount: brl(1000),
        idempotencyKey: 'client-key-abc',
      ).toMap();

      expect(map.containsKey('description'), isFalse);
    });

    test('includes description when provided', () {
      final map = TransferRequestDto(
        fromAccountId: 'acc-src-001',
        toAccountId: 'acc-dst-001',
        amount: brl(1000),
        idempotencyKey: 'client-key-abc',
        description: 'Aluguel de maio',
      ).toMap();

      expect(map['description'], 'Aluguel de maio');
    });

    test('does not serialize legacy branch/account number fields', () {
      final map = TransferRequestDto(
        fromAccountId: 'acc-src-001',
        toAccountId: 'acc-dst-001',
        amount: brl(1000),
        idempotencyKey: 'client-key-abc',
      ).toMap();

      expect(map.containsKey('from_branch'), isFalse);
      expect(map.containsKey('from_account_number'), isFalse);
      expect(map.containsKey('to_branch'), isFalse);
      expect(map.containsKey('to_account_number'), isFalse);
    });
  });
}
