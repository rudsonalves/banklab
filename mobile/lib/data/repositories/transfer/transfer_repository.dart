import '/core/result/command.dart';
import '/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import '/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import '/data/services/apis/transfer/dtos/transfer_response_dto.dart';

abstract class TransferRepository {
  /// Returns the last transfer receipt successfully fetched in this session.
  TransferReceiptResponseDto? get lastReceipt;

  /// Returns the last transfer successfully created in this session.
  TransferResponseDto? get lastTransfer;

  /// Creates a transfer using the provided transfer data.
  ///
  /// Returns a failure when the origin account data is missing, when the
  /// destination account data is missing, or when the transfer amount is not
  /// greater than zero.
  AsyncResult<TransferResponseDto> transfer(TransferRequestDto dto);

  /// Fetches the receipt for a transfer identified by its transaction
  /// reference.
  AsyncResult<TransferReceiptResponseDto> getTransferReceipt(
    String transactionReference,
  );

  /// Fetches a list of possible internal transfer recipients based on the
  /// provided search query, which may be a partial account number or a CPF/CNPJ
  /// document number.
  AsyncResult<List<RecipientInfoDto>> getInternalRecipient(
    RecipientRequestDto dto,
  );
}
