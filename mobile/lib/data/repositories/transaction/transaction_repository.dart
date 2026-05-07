import '../../../core/result/command.dart';
import '../../services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '../../services/apis/transfer/dtos/transfer_request_dto.dart';
import '../../services/apis/transfer/dtos/transfer_response_dto.dart';

abstract class TransactionRepository {
  TransferReceiptResponseDto? get lastReceipt;
  TransferResponseDto? get lastTransfer;

  AsyncResult<TransferResponseDto> transfer(TransferRequestDto dto);

  AsyncResult<TransferReceiptResponseDto> getTransferReceipt(
    String transactionReference,
  );
}
