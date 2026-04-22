import 'package:flutter/material.dart';

import 'viewmodel/home_viewmodel.dart';

class HomePage extends StatefulWidget {
  final HomeViewmodel viewModel;
  final String? accountId;

  const HomePage({
    super.key,
    required this.viewModel,
    this.accountId,
  });

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  @override
  void initState() {
    super.initState();
    widget.viewModel.initialize(accountId: widget.accountId);
  }

  @override
  void dispose() {
    widget.viewModel.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Minha conta'),
        actions: [
          IconButton(
            onPressed: widget.viewModel.refreshBalance,
            icon: const Icon(Icons.refresh_rounded),
            tooltip: 'Atualizar saldo',
          ),
        ],
      ),
      body: AnimatedBuilder(
        animation: widget.viewModel,
        builder: (context, _) {
          return ListView(
            padding: const EdgeInsets.all(20),
            children: [
              Container(
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    colors: [
                      colorScheme.primary,
                      colorScheme.primaryContainer,
                    ],
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                  ),
                  borderRadius: BorderRadius.circular(24),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Saldo em conta',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        color: colorScheme.onPrimary,
                      ),
                    ),
                    const SizedBox(height: 12),
                    Text(
                      _balanceLabel(),
                      style: Theme.of(context).textTheme.headlineMedium
                          ?.copyWith(
                            color: colorScheme.onPrimary,
                            fontWeight: FontWeight.w700,
                          ),
                    ),
                    const SizedBox(height: 12),
                    Text(
                      _supportingLabel(),
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: colorScheme.onPrimary.withValues(alpha: 0.86),
                      ),
                    ),
                    if (widget.viewModel.isLoading) ...[
                      const SizedBox(height: 16),
                      LinearProgressIndicator(
                        borderRadius: BorderRadius.circular(999),
                        minHeight: 6,
                      ),
                    ],
                  ],
                ),
              ),
              const SizedBox(height: 20),
              if (widget.viewModel.balanceError != null)
                Card(
                  color: colorScheme.errorContainer,
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Text(
                      widget.viewModel.balanceError!.message,
                      style: TextStyle(color: colorScheme.onErrorContainer),
                    ),
                  ),
                ),
              const SizedBox(height: 8),
              Text(
                'Acessos rápidos',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: _ActionTile(
                      icon: Icons.receipt_long_rounded,
                      title: 'Extrato',
                      subtitle: 'Disponibilizar em breve',
                      onTap: () => _showPendingFeature('Extrato'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _ActionTile(
                      icon: Icons.swap_horiz_rounded,
                      title: 'Transferir',
                      subtitle: 'Disponibilizar em breve',
                      onTap: () => _showPendingFeature('Transferir'),
                    ),
                  ),
                ],
              ),
            ],
          );
        },
      ),
    );
  }

  String _balanceLabel() {
    final balance = widget.viewModel.balance;
    if (balance == null) return 'R\$ --';

    final value = (balance.balance / 100)
        .toStringAsFixed(2)
        .replaceAll('.', ',');
    return 'R\$ $value';
  }

  String _supportingLabel() {
    if (!widget.viewModel.hasAccountId) {
      return 'Nenhuma conta vinculada na sessao atual. Passe o accountId na rota para carregar o saldo.';
    }

    final lastUpdatedAt = widget.viewModel.lastUpdatedAt;
    final updatedText = lastUpdatedAt == null
        ? 'Aguardando primeira leitura do saldo.'
        : 'Atualizado as ${_twoDigits(lastUpdatedAt.hour)}:${_twoDigits(lastUpdatedAt.minute)}:${_twoDigits(lastUpdatedAt.second)}.';

    return 'Conta ${widget.viewModel.accountId} \u2022 $updatedText';
  }

  String _twoDigits(int value) => value.toString().padLeft(2, '0');

  void _showPendingFeature(String featureName) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('$featureName sera disponibilizado em breve.')),
    );
  }
}

class _ActionTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;

  const _ActionTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(20),
      child: Ink(
        padding: const EdgeInsets.all(18),
        decoration: BoxDecoration(
          color: Theme.of(context).colorScheme.surfaceContainerHighest,
          borderRadius: BorderRadius.circular(20),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(icon),
            const SizedBox(height: 16),
            Text(
              title,
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w700,
              ),
            ),
            const SizedBox(height: 6),
            Text(subtitle),
          ],
        ),
      ),
    );
  }
}
