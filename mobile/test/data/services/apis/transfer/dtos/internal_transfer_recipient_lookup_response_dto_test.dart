import 'package:bankflow/data/services/apis/transfer/dtos/internal_transfer_recipient_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/internal_transfer_recipient_lookup_response_dto.dart';
import 'package:flutter_test/flutter_test.dart';

Map<String, dynamic> _recipientPayload({
  String accountId = 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b',
  String holderName = 'Maria Silva',
  String document = '***.456.789-**',
  String branch = '0001',
  String accountNumber = '00067890',
  String? accountType = 'checking',
}) {
  return {
    'account_id': accountId,
    'holder_name': holderName,
    'document': document,
    'branch': branch,
    'account_number': accountNumber,
    if (accountType != null) 'account_type': accountType,
  };
}

void main() {
  group('InternalTransferRecipientDto.fromMap', () {
    test('parses a representative recipient account payload', () {
      final dto = InternalTransferRecipientDto.fromMap(_recipientPayload());

      expect(dto.accountId, 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b');
      expect(dto.holderName, 'Maria Silva');
      expect(dto.document, '***.456.789-**');
      expect(dto.branch, '0001');
      expect(dto.accountNumber, '00067890');
      expect(dto.accountType, 'checking');
    });

    test('handles optional account_type when absent', () {
      final dto = InternalTransferRecipientDto.fromMap(
        _recipientPayload(accountType: null),
      );

      expect(dto.accountType, isNull);
    });

    test('ignores prohibited fields if present in payload', () {
      final payload = _recipientPayload()
        ..addAll({
          'customer_id': 'cus_123',
          'balance': 90000,
          'full_document': '12345678901',
          'phone': '+5511999999999',
          'email': 'maria@example.com',
          'internal_transaction_id': 'tx_internal_123',
        });

      final dto = InternalTransferRecipientDto.fromMap(payload);

      expect(dto.accountId, 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b');
      expect(dto.document, '***.456.789-**');
    });

    test('does not expose prohibited fields as DTO getters', () {
      final dto = InternalTransferRecipientDto.fromMap(_recipientPayload());

      expect(() => (dto as dynamic).customerId, throwsNoSuchMethodError);
      expect(() => (dto as dynamic).balance, throwsNoSuchMethodError);
      expect(() => (dto as dynamic).fullDocument, throwsNoSuchMethodError);
      expect(() => (dto as dynamic).phone, throwsNoSuchMethodError);
      expect(() => (dto as dynamic).email, throwsNoSuchMethodError);
      expect(
        () => (dto as dynamic).internalTransactionId,
        throwsNoSuchMethodError,
      );
    });
  });

  group('InternalTransferRecipientLookupResponseDto.fromMap', () {
    test('parses zero accounts', () {
      final dto = InternalTransferRecipientLookupResponseDto.fromMap({
        'accounts': <Map<String, dynamic>>[],
      });

      expect(dto.accounts, isEmpty);
    });

    test('parses one account from account lookup payload', () {
      final dto = InternalTransferRecipientLookupResponseDto.fromMap({
        'accounts': [_recipientPayload()],
      });

      expect(dto.accounts, hasLength(1));
      expect(dto.accounts.single.holderName, 'Maria Silva');
      expect(dto.accounts.single.branch, '0001');
      expect(dto.accounts.single.accountNumber, '00067890');
    });

    test('parses multiple accounts from CPF lookup payload', () {
      final dto = InternalTransferRecipientLookupResponseDto.fromMap({
        'accounts': [
          _recipientPayload(
            accountId: 'acc_001',
            branch: '0001',
            accountNumber: '00011111',
            accountType: 'checking',
          ),
          _recipientPayload(
            accountId: 'acc_002',
            branch: '0001',
            accountNumber: '00022222',
            accountType: 'savings',
          ),
        ],
      });

      expect(dto.accounts, hasLength(2));
      expect(dto.accounts.first.accountId, 'acc_001');
      expect(dto.accounts.last.accountId, 'acc_002');
      expect(dto.accounts.last.accountType, 'savings');
    });
  });
}
