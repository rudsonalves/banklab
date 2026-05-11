import '/core/result/command.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/repositories/transaction/transaction_repository.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';
import '/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import '/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import '/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import 'inputs/transfer_draft.dart';

export 'inputs/transfer_draft.dart';

class TransferUsecase {
  final AccountRepository _accountRepo;
  final TransactionRepository _transactionRepo;

  TransferUsecase({
    required AccountRepository accountRepo,
    required TransactionRepository transactionRepo,
  }) : _accountRepo = accountRepo,
       _transactionRepo = transactionRepo;

  Stream<BalanceResponseDto> balance() => _accountRepo.balance();

  List<AccountSummaryResponseDto>? get accounts => _accountRepo.accounts;
  AccountSummaryResponseDto? get selectedAccount =>
      _accountRepo.selectedAccount;

  AsyncResult<TransferResponseDto> transfer(TransferDraft transfer) async {
    final account = selectedAccount;
    if (account == null) {
      return const Failure<TransferResponseDto>(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'No account selected.',
        ),
      );
    }

    if (transfer.idempotencyKey.isEmpty) {
      return const Failure<TransferResponseDto>(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Idempotency key is required for transfer.',
        ),
      );
    }

    final toAccountId = transfer.toAccountId.trim();
    if (toAccountId.isEmpty || toAccountId == account.id) {
      return const Failure<TransferResponseDto>(
        AppError(
          code: AppErrorCode.invalidData,
          message:
              'Destination account ID cannot be empty or the same '
              'as the source account.',
        ),
      );
    }

    final dto = TransferRequestDto(
      fromAccountId: account.id,
      toAccountId: toAccountId,
      amount: transfer.amount,
      description: transfer.description,
      idempotencyKey: transfer.idempotencyKey,
    );

    return _transactionRepo.transfer(dto);
  }

  AsyncResult<TransferReceiptResponseDto> getTransferReceipt(
    String reference,
  ) => _transactionRepo.getTransferReceipt(reference);

  AsyncResult<Unit> selectAccount(String accountId) async {
    _accountRepo.selectAccount(accountId);

    final account = _accountRepo.selectedAccount;
    if (account == null || account.id != accountId) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Account not found.',
        ),
      );
    }

    return _accountRepo.loadBalance();
  }

  AsyncResult<List<RecipientInfoDto>> getInternalRecipient(
    RecipientRequestDto recipient,
  ) async {
    if (recipient.toMap().isEmpty) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Recipient search query cannot be empty.',
        ),
      );
    }

    return _transactionRepo.getInternalRecipient(recipient);
  }
}
