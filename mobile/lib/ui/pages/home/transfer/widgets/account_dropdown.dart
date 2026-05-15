import 'package:flutter/material.dart';

import '../../../../../data/services/apis/account/dtos/account_summary_response_dto.dart';

class AccountDropdown extends StatelessWidget {
  final List<AccountSummaryResponseDto> accounts;
  final String selectedAccountId;
  final ValueChanged<String?> onChanged;

  const AccountDropdown({
    super.key,
    required this.accounts,
    required this.selectedAccountId,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return DropdownButtonFormField<String>(
      initialValue: selectedAccountId,
      hint: const Text('Selecione uma conta'),
      items: accounts.map((account) {
        return DropdownMenuItem(
          value: account.id,
          child: Text('${account.branch} - ${account.number}'),
        );
      }).toList(),
      onChanged: onChanged,
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
}
