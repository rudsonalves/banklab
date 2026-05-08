import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import 'package:bankflow/domain/common/receipt/enums/transfer_receipt_status.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';

Map<String, dynamic> _validPayload() => {
  'operation_type': 'transfer_out',
  'amount': 2500,
  'status': 'completed',
  'transaction_reference': '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
  'operation_date': '2026-05-06T12:30:00Z',
  'source_branch': '0001',
  'source_account_number': '00012345',
  'destination_branch': '0001',
  'destination_account_number': '00067890',
  'recipient_name': 'Maria Silva',
};

Money _brl(int cents) =>
    Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

void main() {
  group('TransferReceiptResponseDto.fromMap', () {
    test(
      'parses all required fields from a representative success payload',
      () {
        final dto = TransferReceiptResponseDto.fromMap(_validPayload());

        expect(dto.operationType, 'transfer_out');
        expect(dto.amount, _brl(2500));
        expect(dto.status, TransferReceiptStatus.completed);
        expect(
          dto.transactionReference,
          '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
        );
        expect(dto.operationDate, DateTime.parse('2026-05-06T12:30:00Z'));
        expect(dto.sourceBranch, '0001');
        expect(dto.sourceAccountNumber, '00012345');
        expect(dto.destinationBranch, '0001');
        expect(dto.destinationAccountNumber, '00067890');
        expect(dto.recipientName, 'Maria Silva');
      },
    );

    test('parses amount as cents via ApiParse.toMoney', () {
      final dto = TransferReceiptResponseDto.fromMap(_validPayload());

      expect(dto.amount, _brl(2500));
    });

    test('parses optional description when present', () {
      final map = _validPayload();
      map['description'] = 'Aluguel de maio';

      final dto = TransferReceiptResponseDto.fromMap(map);

      expect(dto.description, 'Aluguel de maio');
    });

    test('keeps description null when absent', () {
      final dto = TransferReceiptResponseDto.fromMap(_validPayload());

      expect(dto.description, isNull);
    });

    test('parses status to TransferReceiptStatus enum', () {
      final dto = TransferReceiptResponseDto.fromMap(_validPayload());

      expect(dto.status, TransferReceiptStatus.completed);
      expect(dto.status.isSuccess, isTrue);
    });

    test('parses operation_date as UTC DateTime', () {
      final dto = TransferReceiptResponseDto.fromMap(_validPayload());

      expect(dto.operationDate, DateTime.utc(2026, 5, 6, 12, 30, 0));
    });

    test('throws FormatException for unknown status value', () {
      final map = _validPayload();
      map['status'] = 'unknown_status';

      expect(
        () => TransferReceiptResponseDto.fromMap(map),
        throwsA(isA<FormatException>()),
      );
    });

    test(
      'ignores internal account and customer ID fields if present in payload',
      () {
        // The backend domain model carries SourceAccountID, SourceCustomerID,
        // DestinationAccountID, DestinationCustomerID as internal UUIDs.
        // These must never appear in the API response or DTO.
        final map = _validPayload();
        map['source_account_id'] = 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b';
        map['source_customer_id'] = 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee';
        map['destination_account_id'] = 'ab3b2800-1234-5678-abcd-ef0123456789';
        map['destination_customer_id'] = '11111111-2222-3333-4444-555555555555';

        final dto = TransferReceiptResponseDto.fromMap(map);

        // Parsing must succeed and expose only branch + account_number.
        expect(dto.sourceBranch, '0001');
        expect(dto.destinationAccountNumber, '00067890');
      },
    );

    test('does not expose internal account or customer IDs as DTO fields', () {
      final dto = TransferReceiptResponseDto.fromMap(_validPayload());

      expect(
        () => (dto as dynamic).sourceAccountId,
        throwsNoSuchMethodError,
      );
      expect(
        () => (dto as dynamic).sourceCustomerId,
        throwsNoSuchMethodError,
      );
      expect(
        () => (dto as dynamic).destinationAccountId,
        throwsNoSuchMethodError,
      );
      expect(
        () => (dto as dynamic).destinationCustomerId,
        throwsNoSuchMethodError,
      );
    });

    test('throws when a required field is missing', () {
      for (final field in _validPayload().keys.toList()) {
        final map = _validPayload()..remove(field);
        expect(
          () => TransferReceiptResponseDto.fromMap(map),
          throwsA(anything),
          reason: 'expected throw when "$field" is absent',
        );
      }
    });
  });
}
