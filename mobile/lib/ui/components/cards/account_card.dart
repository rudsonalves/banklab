import 'package:flutter/material.dart';

import '/core/result/command.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '../text/card_text_row.dart';

class AccountCard extends StatefulWidget {
  final AccountSummaryResponseDto? selectedAccount;
  final List<AccountSummaryResponseDto>? accounts;
  final Command1<Unit, String> onSelecteAccount;

  const AccountCard({
    super.key,
    this.selectedAccount,
    this.accounts,
    required this.onSelecteAccount,
  });

  @override
  State<AccountCard> createState() => _AccountCardState();
}

class _AccountCardState extends State<AccountCard> {
  AccountSummaryResponseDto? get selectedAccount => widget.selectedAccount;

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
            CardTextRow(
              label: 'Banco',
              value: 'BankName',
            ),
            CardTextRow(
              label: 'Conta',
              value: selectedAccount != null
                  ? '${selectedAccount!.branch} - ${selectedAccount!.number}'
                  : '••••••',
            ),
          ],
        ),
      ),
    );
  }
}
