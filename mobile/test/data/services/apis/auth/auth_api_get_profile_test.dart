import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/services/apis/auth/auth_api.dart';
import 'package:bankflow/domain/common/auth/models/auth_session/auth_session.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('AuthSession.fromApi', () {
    test('parses the real /auth/session JSON contract', () {
      final session = AuthSession.fromApi(_authSessionDataJson());

      expect(session.user.userId, 'd3de5f8b-4892-42e8-9680-979cf3f37844');
      expect(session.user.email, 'user@example.com');
      expect(session.user.phone, '+5527999999999');
      expect(session.customer, isNotNull);
      expect(session.customer?.id, '6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3');
      expect(session.customer?.name, 'Maria Silva');
      expect(session.customer?.cpf, '12345678901');
      expect(session.customer?.birthDate, DateTime(1990, 1, 15));
      expect(
        session.customer?.createdAt,
        DateTime.parse('2026-05-29T10:00:00Z'),
      );

      expect(session.readiness.onboardingCompleted, isTrue);
      expect(session.readiness.approved, isTrue);
      expect(session.readiness.hasOperationalAccount, isTrue);
      expect(
        session.readiness.transactionPasswordStatus,
        TransactionPasswordStatus.active,
      );
      expect(session.readiness.canAccessHome, isTrue);
    });
  });

  group('AuthApi.getProfile', () {
    test(
      'calls GET /auth/session and parses canonical envelope data',
      () async {
        final client = _FakeRestClient(
          getResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _authSessionEnvelopeJson(),
            ),
          ),
        );
        final api = AuthApi(client);

        final result = await api.getAuthSession();

        expect(result, isA<Success<AuthSession>>());
        expect(client.getCalls, 1);
        expect(client.lastGetRequest?.path, '/auth/session');

        final session = result.value!;
        expect(session.user.email, 'user@example.com');
        expect(session.customer?.name, 'Maria Silva');
        expect(
          session.readiness.transactionPasswordStatus,
          TransactionPasswordStatus.active,
        );
        expect(session.readiness.canAccessHome, isTrue);
      },
    );
  });
}

Map<String, dynamic> _authSessionEnvelopeJson() => {
  'data': _authSessionDataJson(),
  'error': null,
};

Map<String, dynamic> _authSessionDataJson() => {
  'user': {
    'id': 'd3de5f8b-4892-42e8-9680-979cf3f37844',
    'email': 'user@example.com',
    'phone': '+5527999999999',
    'role': 'customer',
  },
  'customer': {
    'id': '6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3',
    'name': 'Maria Silva',
    'cpf': '12345678901',
    'birth_date': '1990-01-15',
    'created_at': '2026-05-29T10:00:00Z',
  },
  'readiness': {
    'onboarding_completed': true,
    'approved': true,
    'has_operational_account': true,
    'transaction_password_status': 'active',
    'can_access_home': true,
  },
};

class _FakeRestClient implements RestClient {
  _FakeRestClient({required this.getResult});

  final Result<RestClientResponse> getResult;
  RestClientRequest? lastGetRequest;
  int getCalls = 0;

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async {
    getCalls++;
    lastGetRequest = request;
    return getResult;
  }

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> put(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> patch(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> delete(RestClientRequest request) async =>
      throw UnimplementedError();
}
