import 'package:flutter/material.dart';

class CriterialItemRow extends StatelessWidget {
  final bool checked;
  final String label;

  const CriterialItemRow({
    super.key,
    required this.checked,
    required this.label,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(
          checked ? Icons.check_circle : Icons.radio_button_unchecked,
          size: 18,
          color: checked ? Colors.green : Colors.grey,
        ),
        const SizedBox(width: 8),
        Text(label),
      ],
    );
  }
}
