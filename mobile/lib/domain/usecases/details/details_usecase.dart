import '/core/result/result.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/repositories/transfer/transfer_repository.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';

class DetailsUsecase {
  final AccountRepository _accountRepo;
  final TransferRepository _transferRepo;

  DetailsUsecase({
    required AccountRepository accountRepo,
    required TransferRepository transferRepo,
  }) : _accountRepo = accountRepo,
       _transferRepo = transferRepo;

  AccountSummaryResponseDto? get selectedAccount =>
      _accountRepo.selectedAccount;

  AsyncResult<TransferReceiptResponseDto> getTransferReceipt(
    String transactionReference,
  ) => _transferRepo.getTransferReceipt(transactionReference);
}
