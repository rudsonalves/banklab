import '/core/result/result.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/repositories/transaction/transaction_repository.dart';
import '/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '../../../data/services/apis/account/dtos/account_summary_response_dto.dart';

class DetailsUsecase {
  final AccountRepository _accountRepo;
  final TransactionRepository _transactionRepo;

  DetailsUsecase({
    required AccountRepository accountRepo,
    required TransactionRepository transactionRepo,
  }) : _accountRepo = accountRepo,
       _transactionRepo = transactionRepo;

  AccountSummaryResponseDto? get selectedAccount =>
      _accountRepo.selectedAccount;

  AsyncResult<TransferReceiptResponseDto> getTransferReceipt(
    String transactionReference,
  ) => _transactionRepo.getTransferReceipt(transactionReference);
}
