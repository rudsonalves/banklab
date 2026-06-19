import 'package:flutter/material.dart';

const installationLimitDialogTitle = 'Limite de instalações atingido';
const installationLimitDialogPrimaryMessage =
    'Esta conta já possui 3 instalações cadastradas. A instalação atual ainda '
    'não está autorizada.';
const installationLimitDialogSecondaryMessage =
    'Acesse sua conta por uma instalação já autorizada e remova uma instalação '
    'antiga para liberar espaço. Depois, tente entrar novamente neste app.';
const installationLimitDialogButtonLabel = 'Entendi';

Future<void> showInstallationLimitDialog(BuildContext context) {
  return showDialog<void>(
    context: context,
    builder: (context) => AlertDialog(
      title: const Text(installationLimitDialogTitle),
      content: const Text(
        '$installationLimitDialogPrimaryMessage\n\n'
        '$installationLimitDialogSecondaryMessage',
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: const Text(installationLimitDialogButtonLabel),
        ),
      ],
    ),
  );
}
