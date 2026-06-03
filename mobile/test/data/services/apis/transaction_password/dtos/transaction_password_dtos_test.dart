import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/enums/transaction_password_status.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('TransactionPasswordStatusResponseDto.fromApi', () {
    test(
      'parses create transaction password response as active for domain compatibility',
      () {
        final dto = TransactionPasswordStatusResponseDto.fromApi({
          'user_id': 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b',
          'status': 'active',
          'created_at': '2026-05-18T12:03:00Z',
        });

        expect(dto.userId, 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b');
        expect(dto.status, TransactionPasswordStatus.active);
        expect(dto.createdAt, DateTime.parse('2026-05-18T12:03:00Z'));
      },
    );

    test('parses blocked status value', () {
      final dto = TransactionPasswordStatusResponseDto.fromApi({
        'user_id': 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b',
        'status': 'blocked',
        'created_at': '2026-05-18T12:03:00Z',
      });

      expect(dto.status, TransactionPasswordStatus.blocked);
    });

    test('throws for unknown status value', () {
      expect(
        () => TransactionPasswordStatusResponseDto.fromApi({
          'user_id': 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b',
          'status': 'invalid',
          'created_at': '2026-05-18T12:03:00Z',
        }),
        throwsA(isA<ArgumentError>()),
      );
    });
  });

  group('CreateTransactionPasswordRequestDto.toMap', () {
    test('serializes transaction password request payload', () {
      final dto = CreateTransactionPasswordRequestDto(
        password: '123456',
        confirmation: '123456',
      );

      final map = dto.toMap();

      expect(map['transaction_password'], '123456');
      expect(map['transaction_password_confirmation'], '123456');
    });
  });
}
