import 'package:bankflow/core/services/client_http/client/rest_client_request.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('RestClientRequest', () {
    test('copyWith should keep original values when not provided', () {
      const request = RestClientRequest(
        path: '/users',
        headers: {'Authorization': 'Bearer token'},
        queryParameters: {'page': 1},
        body: {'name': 'Ada'},
      );

      final copied = request.copyWith();

      expect(copied.path, '/users');
      expect(copied.headers, {'Authorization': 'Bearer token'});
      expect(copied.queryParameters, {'page': 1});
      expect(copied.body, {'name': 'Ada'});
    });

    test('copyWith should replace provided fields', () {
      const request = RestClientRequest(
        path: '/users',
        headers: {'Authorization': 'Bearer token'},
        queryParameters: {'page': 1},
        body: {'name': 'Ada'},
      );

      final copied = request.copyWith(
        path: '/accounts',
        headers: {'X-Trace-Id': 'abc'},
        queryParameters: {'limit': 10},
        body: {'active': true},
      );

      expect(copied.path, '/accounts');
      expect(copied.headers, {'X-Trace-Id': 'abc'});
      expect(copied.queryParameters, {'limit': 10});
      expect(copied.body, {'active': true});
    });
  });
}
