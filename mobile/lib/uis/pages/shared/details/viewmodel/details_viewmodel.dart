import '/core/result/command.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '/domain/usecases/details/details_usecase.dart';

class DetailsViewmodel {
  final DetailsUsecase _usecase;

  DetailsViewmodel({
    required DetailsUsecase usecase,
  }) : _usecase = usecase {
    getTransferReceipt = Command1(_usecase.getTransferReceipt);
  }

  late final Command1<TransferReceiptResponseDto, String> getTransferReceipt;

  AccountSummaryResponseDto? get selectedAccount => _usecase.selectedAccount;
}
