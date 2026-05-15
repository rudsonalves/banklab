import 'package:flutter/material.dart';

import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';

class RecipientCard extends StatefulWidget {
  final RecipientInfoDto selectedRecipient;

  const RecipientCard({super.key, required this.selectedRecipient});

  @override
  State<RecipientCard> createState() => _RecipientCardState();
}

class _RecipientCardState extends State<RecipientCard> {
  RecipientInfoDto get selectedRecipient => widget.selectedRecipient;

  @override
  Widget build(BuildContext context) {
    return Card(
      color: Theme.of(context).colorScheme.surfaceContainerHigh,
      margin: EdgeInsets.zero,
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            _buildCardTextRow('Destinatário', selectedRecipient.holderName),
            _buildCardTextRow(
              'Conta',
              '0001 - ${selectedRecipient.accountNumber}',
            ),
            _buildCardTextRow('Documento', selectedRecipient.document),
          ],
        ),
      ),
    );
  }

  Widget _buildCardTextRow(String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          '$label: ',
          textAlign: TextAlign.left,
          style: Theme.of(context).textTheme.bodyLarge!.copyWith(
            fontWeight: FontWeight.w500,
          ),
        ),
        Text(
          value,
          textAlign: TextAlign.right,
          overflow: TextOverflow.ellipsis,
          style: Theme.of(context).textTheme.bodyLarge!.copyWith(
            fontWeight: FontWeight.w700,
          ),
        ),
      ],
    );
  }
}
