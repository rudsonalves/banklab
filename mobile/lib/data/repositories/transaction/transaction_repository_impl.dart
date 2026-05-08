import '../../../core/result/result.dart';
import '../../services/apis/receipt/api_receipt.dart';
import '../../services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '../../services/apis/transfer/api_transfer.dart';
import '../../services/apis/transfer/dtos/recipient_info_dto.dart';
import '../../services/apis/transfer/dtos/recipient_request_dto.dart';
import '../../services/apis/transfer/dtos/transfer_request_dto.dart';
import '../../services/apis/transfer/dtos/transfer_response_dto.dart';
import 'transaction_repository.dart';

class TransactionRepositoryImpl implements TransactionRepository {
  final ApiTransfer _apiTransfer;
  final ApiReceipt _apiReceipt;

  TransactionRepositoryImpl({
    required ApiTransfer apiTransfer,
    required ApiReceipt apiReceipt,
  }) : _apiTransfer = apiTransfer,
       _apiReceipt = apiReceipt;

  TransferReceiptResponseDto? _lastReceipt;
  TransferResponseDto? _lastTransfer;

  @override
  TransferReceiptResponseDto? get lastReceipt => _lastReceipt;

  @override
  TransferResponseDto? get lastTransfer => _lastTransfer;

  @override
  AsyncResult<TransferResponseDto> transfer(TransferRequestDto dto) async {
    if (dto.fromAccountId.trim().isEmpty) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'No account selected.',
        ),
      );
    }

    if (dto.toAccountId.trim().isEmpty) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Destination account is required.',
        ),
      );
    }

    if (dto.amount.amount.isZero || dto.amount.isNegative) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Amount must be greater than zero.',
        ),
      );
    }

    final result = await _apiTransfer.transfer(dto);

    _lastTransfer = result.isSuccess ? result.value : null;

    return result;
  }

  @override
  AsyncResult<TransferReceiptResponseDto> getTransferReceipt(
    String transactionReference,
  ) async {
    final result = await _apiReceipt.getReceipt(transactionReference);

    _lastReceipt = result.isSuccess ? result.value : null;

    return result;
  }

  @override
  AsyncResult<List<RecipientInfoDto>> getInternalRecipient(
    RecipientRequestDto dto,
  ) async {
    if (dto.toMap().isEmpty) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Search query is required.',
        ),
      );
    }

    final result = await _apiTransfer.getInternalRecipient(dto);

    return result;
  }
}
