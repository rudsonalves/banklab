import 'package:bankflow/data/services/apis/transfer/dtos/internal_transfer_recipient_lookup_query_dto.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('InternalTransferRecipientLookupQueryDto.toMap', () {
    test('serializes lookup by branch and account number', () {
      final map = InternalTransferRecipientLookupQueryDto.byAccount(
        branch: '0001',
        accountNumber: '00067890',
      ).toMap();

      expect(map['branch'], '0001');
      expect(map['account_number'], '00067890');
      expect(map.containsKey('document'), isFalse);
    });

    test('serializes CPF lookup through document', () {
      final map = InternalTransferRecipientLookupQueryDto.byCpf(
        cpf: '12345678901',
      ).toMap();

      expect(map['document'], '12345678901');
      expect(map.containsKey('branch'), isFalse);
      expect(map.containsKey('account_number'), isFalse);
    });
  });
}
