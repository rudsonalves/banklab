import 'package:flutter/foundation.dart';

import '/core/result/command.dart';
import '/domain/common/auth/models/register_draft_snapshot.dart';
import '/domain/common/auth/models/register_draft_state.dart';
import '/domain/usecases/register/register_usecase.dart';

class RegisterViewmodel extends ChangeNotifier {
  final RegisterUsecase _usecase;

  RegisterViewmodel({
    required RegisterUsecase usecase,
  }) : _usecase = usecase {
    submitCPF = Command1(_usecase.submitCPF);
    submitName = Command1(_usecase.submitName);
    submitBirthDate = Command1(_usecase.submitBirthDate);
    submitAndRequestEmailToken = Command1(_usecase.submitAndRequestEmailToken);
    submitAndRequestPhoneToken = Command1(_usecase.submitAndRequestPhoneToken);
    confirmEmailToken = Command1(_usecase.confirmEmailToken);
    confirmPhoneToken = Command1(_usecase.confirmPhoneToken);
    submitPassword = Command1(_usecase.submitPassword);
    register = Command0(_usecase.register);
    reset = Command0(_usecase.reset);
  }

  late final Command1<Unit, String> submitCPF;
  late final Command1<Unit, String> submitName;
  late final Command1<Unit, DateTime> submitBirthDate;
  late final Command1<Unit, String> submitAndRequestEmailToken;
  late final Command1<Unit, String> submitAndRequestPhoneToken;
  late final Command1<Unit, String> confirmEmailToken;
  late final Command1<Unit, String> confirmPhoneToken;
  late final Command1<Unit, (String, String)> submitPassword;
  late final Command0<Unit> register;
  late final Command0<Unit> reset;

  RegisterDraftSnapshot? get draft => _usecase.draft;

  RegisterDraftState? get state => _usecase.state;

  void startEmptyRegisterState() => _usecase.startEmptyRegisterState();
}
