import 'http_method.dart';

enum StepUpOperation {
  internalTransfer(
    method: HttpMethod.post,
    path: '/accounts/internal-transfers',
  ),
  installationRegistration(
    method: HttpMethod.post,
    path: '/security/installations',
  );

  const StepUpOperation({
    required this.method,
    required this.path,
  });

  final HttpMethod method;
  final String path;
}
