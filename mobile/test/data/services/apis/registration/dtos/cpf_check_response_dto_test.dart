import 'package:bankflow/data/services/apis/registration/dtos/cpf_check_response_dto.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('CpfCheckResponseDto', () {
    test('parses available field from API payload', () {
      final dto = CpfCheckResponseDto.fromMap({
        'cpf': '12345678909',
        'exists': false,
        'available': true,
      });

      expect(dto.cpf, '12345678909');
      expect(dto.exists, isFalse);
      expect(dto.available, isTrue);
    });
  });
}
