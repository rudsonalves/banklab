import 'package:flutter/material.dart';

import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';

class DropdownRecipient extends StatefulWidget {
  final List<RecipientInfoDto> receipientAccounts;
  final String? inicialValue;
  final void Function(String?)? onChanged;

  const DropdownRecipient({
    super.key,
    required this.receipientAccounts,
    this.inicialValue,
    this.onChanged,
  });

  @override
  State<DropdownRecipient> createState() => _DropdownRecipientState();
}

class _DropdownRecipientState extends State<DropdownRecipient> {
  @override
  Widget build(BuildContext context) {
    return DropdownButtonFormField<String>(
      isExpanded: true,
      initialValue: widget.inicialValue,
      items: _buildDropdownMenuItems(widget.receipientAccounts),
      selectedItemBuilder: (_) =>
          _buildSelectedItemBuilder(widget.receipientAccounts),

      onChanged: widget.onChanged,

      decoration: InputDecoration(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 12,
          vertical: 12,
        ),
      ),
    );
  }

  List<DropdownMenuItem<String>> _buildDropdownMenuItems(
    List<RecipientInfoDto> accounts,
  ) {
    return accounts
        .map(
          (recipient) => DropdownMenuItem<String>(
            value: recipient.accountId,
            child: Text(
              '${recipient.holderName} - Cc 0001/${recipient.accountNumber}',
              overflow: TextOverflow.ellipsis,
              maxLines: 1,
            ),
          ),
        )
        .toList();
  }

  List<Widget> _buildSelectedItemBuilder(List<RecipientInfoDto> accounts) {
    return accounts
        .map(
          (recipient) => Text(
            '${recipient.holderName} - Cc 0001/${recipient.accountNumber}',
            overflow: TextOverflow.ellipsis,
            maxLines: 1,
          ),
        )
        .toList();
  }
}
