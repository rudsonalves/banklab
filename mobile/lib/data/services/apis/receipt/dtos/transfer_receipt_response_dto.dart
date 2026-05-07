import 'package:money2/money2.dart';

import '/domain/common/receipt/enums/transfer_receipt_status.dart';
import '../../core/api_parse.dart';

class TransferReceiptResponseDto {
  final String operationType;
  final Money amount;
  final TransferReceiptStatus status;
  final String transactionReference;
  final DateTime operationDate;
  final String sourceBranch;
  final String sourceAccountNumber;
  final String destinationBranch;
  final String destinationAccountNumber;
  final String recipientName;

  TransferReceiptResponseDto({
    required this.operationType,
    required this.amount,
    required this.status,
    required this.transactionReference,
    required this.operationDate,
    required this.sourceBranch,
    required this.sourceAccountNumber,
    required this.destinationBranch,
    required this.destinationAccountNumber,
    required this.recipientName,
  });

  factory TransferReceiptResponseDto.fromMap(Map<String, dynamic> map) {
    return TransferReceiptResponseDto(
      operationType: map['operation_type'] as String,
      amount: ApiParse.toMoney(map['amount']),
      status: TransferReceiptStatus.fromString(map['status'] as String),
      transactionReference: map['transaction_reference'] as String,
      operationDate: DateTime.parse(map['operation_date'] as String),
      sourceBranch: map['source_branch'] as String,
      sourceAccountNumber: map['source_account_number'] as String,
      destinationBranch: map['destination_branch'] as String,
      destinationAccountNumber: map['destination_account_number'] as String,
      recipientName: map['recipient_name'] as String,
    );
  }
}
