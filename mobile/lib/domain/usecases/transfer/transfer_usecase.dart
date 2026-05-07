import '/core/result/command.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/repositories/transaction/transaction_repository.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
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

  List<AccountSummaryResponseDto>? get accounts => _accountRepo.accounts;
  AccountSummaryResponseDto? get selectedAccount =>
      _accountRepo.selectedAccount;

  AsyncResult<TransferResponseDto> transfer(TransferDraft transfer) {
    final dto = TransferRequestDto(
      fromAccountNumber: selectedAccount!.number,
      fromBranch: selectedAccount!.branch,
      toAccountNumber: transfer.toAccountNumber,
      toBranch: transfer.toBranch,
      amount: transfer.amount,
      description: transfer.description,
      idempotencyKey: transfer.idempotencyKey,
    );

    return _transactionRepo.transfer(dto);
  }
}
