import '/core/result/result.dart';
import '/data/services/apis/receipt/api_receipt.dart';
import '/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '/data/services/apis/transfer/api_transfer.dart';
import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import '/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import '/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import 'transfer_repository.dart';

class TransferRepositoryImpl implements TransferRepository {
  final ApiTransfer _apiTransfer;
  final ApiReceipt _apiReceipt;

  TransferRepositoryImpl({
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
  AsyncResult<TransferResponseDto> transfer({
    required String token,
    required TransferRequestDto dto,
  }) async {
    if (token.trim().isEmpty) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Step-up token is required.',
        ),
      );
    }

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

    final result = await _apiTransfer.transfer(
      token: token,
      dto: dto,
    );

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
