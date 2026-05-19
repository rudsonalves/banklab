import 'package:flutter/foundation.dart';

import '/data/repositories/auth/auth_repository.dart';

class RegisterViewmodel extends ChangeNotifier {
  final AuthRepository _authRepository;

  RegisterViewmodel({
    required AuthRepository authRepository,
  }) : _authRepository = authRepository {}
}
