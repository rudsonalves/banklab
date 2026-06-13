import '/core/result/command.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/repositories/transaction_password/transaction_password_repository.dart';
import '/data/repositories/transfer/transfer_repository.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';
import '/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import '/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import '/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import 'inputs/protected_transfer_input.dart';

export 'inputs/transfer_draft.dart';

class TransferUsecase {
  final AccountRepository _accountRepo;
  final TransferRepository _transferRepo;
  final TransactionPasswordRepository _transactionPasswordRepo;

  TransferUsecase({
    required AccountRepository accountRepo,
    required TransferRepository transferRepo,
    required TransactionPasswordRepository transactionPasswordRepo,
  }) : _accountRepo = accountRepo,
       _transferRepo = transferRepo,
       _transactionPasswordRepo = transactionPasswordRepo;

  Stream<BalanceResponseDto> balance() => _accountRepo.balance();

  List<AccountSummaryResponseDto>? get accounts => _accountRepo.accounts;
  AccountSummaryResponseDto? get selectedAccount =>
      _accountRepo.selectedAccount;

  AsyncResult<TransferResponseDto> transfer(
    ProtectedTransferInput input,
  ) async {
    final account = selectedAccount;
    if (account == null) {
      return const Failure<TransferResponseDto>(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'No account selected.',
        ),
      );
    }

    if (input.draft.idempotencyKey.trim().isEmpty) {
      return const Failure<TransferResponseDto>(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Idempotency key is required for transfer.',
        ),
      );
    }

    final toAccountId = input.draft.toAccountId.trim();
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

    final dto = TransferRequestDto.fromTransferDraft(
      fromAccountId: account.id,
      draft: input.draft,
    );

    if (dto == null) {
      return Result.failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Failed to create transfer request.',
        ),
      );
    }

    final stepUpResult = await _transactionPasswordRepo
        .authorizeInternalTransfer(
          input.pin,
        );

    if (stepUpResult.isFailure) {
      return Failure(stepUpResult.error!);
    }

    final stepUpData = stepUpResult.value!;

    return _transferRepo.transfer(
      token: stepUpData.stepUpToken,
      dto: dto,
    );
  }

  AsyncResult<TransferReceiptResponseDto> getTransferReceipt(
    String reference,
  ) => _transferRepo.getTransferReceipt(reference);

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

    return _transferRepo.getInternalRecipient(recipient);
  }
}
