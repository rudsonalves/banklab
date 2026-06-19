import '/core/result/command.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/domain/common/auth/models/auth_state.dart';

class InstallationCertificationViewModel {
  InstallationCertificationViewModel({
    required AuthRepository authRepository,
  }) : _authRepository = authRepository {
    certify = Command1(_certify);
  }

  final AuthRepository _authRepository;

  late final Command1<OperationalAuthState, String> certify;

  bool get hasRestrictedInstallationAuth =>
      _authRepository.currentUser is RestrictedInstallationAuthState;

  AsyncResult<OperationalAuthState> _certify(String transactionPassword) {
    return _authRepository.certifyInstallation(transactionPassword);
  }

  Future<void> cancel() async {
    await _authRepository.logout();
  }
}
