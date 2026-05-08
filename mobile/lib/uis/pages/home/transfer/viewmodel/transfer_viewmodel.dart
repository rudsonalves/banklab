import 'package:uuid/uuid.dart';

import '/core/result/command.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import '/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import '/domain/usecases/transfer/transfer_usecase.dart';

class TransferViewmodel {
  final TransferUsecase _usecase;

  TransferViewmodel(this._usecase) {
    transfer = Command1(_transfer);
    receipt = Command1(_usecase.getTransferReceipt);
    selectAccount = Command1(_usecase.selectAccount);
    getInternalRecipient = Command1(_usecase.getInternalRecipient);
  }

  late final Command1<TransferResponseDto, TransferDraft> transfer;
  late final Command1<TransferReceiptResponseDto, String> receipt;
  late final Command1<Unit, String> selectAccount;
  late final Command1<List<RecipientInfoDto>, RecipientRequestDto>
  getInternalRecipient;

  List<AccountSummaryResponseDto>? get accounts => _usecase.accounts;
  AccountSummaryResponseDto? get selectedAccount => _usecase.selectedAccount;

  late final String idempotencyKey;

  AsyncResult<TransferResponseDto> _transfer(TransferDraft draft) {
    final transferWithIdempotency = draft.copyWith(
      idempotencyKey: const Uuid().v7(),
    );
    return _usecase.transfer(transferWithIdempotency);
  }
}
