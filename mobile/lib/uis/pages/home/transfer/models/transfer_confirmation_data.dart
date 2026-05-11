import 'package:money2/money2.dart';

import '/core/resources/app_currencies.dart';
import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';

class TransferConfirmationData {
  final String toAccountId;
  final String toHolderName;
  final String toDocument;
  final String toNumber;

  final String fromAccountId;
  final Money amount;
  final String description;

  TransferConfirmationData({
    required this.toAccountId,
    required this.toHolderName,
    required this.toDocument,
    required this.toNumber,

    required this.fromAccountId,
    required this.amount,
    required this.description,
  });

  factory TransferConfirmationData.fromRecipientInfo({
    required RecipientInfoDto recipientInfo,
    required String fromAccountId,
    required Money amount,
    required String description,
  }) {
    return TransferConfirmationData(
      toAccountId: recipientInfo.accountId,
      toHolderName: recipientInfo.holderName,
      toDocument: recipientInfo.document,
      toNumber: recipientInfo.accountNumber,

      fromAccountId: fromAccountId,
      amount: amount,
      description: description,
    );
  }

  Map<String, dynamic> toMap() => {
    'to_account_id': toAccountId,
    'to_holder_name': toHolderName,
    'to_document': toDocument,
    'to_number': toNumber,
    'from_account_id': fromAccountId,
    'amount': amount.minorUnits.toString(),
    'description': description,
  };

  factory TransferConfirmationData.fromMap(Map<String, dynamic> map) {
    return TransferConfirmationData(
      toAccountId: map['to_account_id'],
      toHolderName: map['to_holder_name'],
      toDocument: map['to_document'],
      toNumber: map['to_number'],
      fromAccountId: map['from_account_id'],
      amount: Money.fromBigIntWithCurrency(
        BigInt.parse(map['amount']),
        appCurrency,
      ),
      description: map['description'],
    );
  }
}
