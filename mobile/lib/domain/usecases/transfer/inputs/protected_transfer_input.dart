import 'transfer_draft.dart';

class ProtectedTransferInput {
  final TransferDraft draft;
  final String pin;

  ProtectedTransferInput({
    required this.draft,
    required this.pin,
  });
}
