import 'package:flutter/material.dart';

import '../../core/base/safe_scaffold.dart';
import 'viewmodel/home_viewmodel.dart';
import 'widgets/action_tite.dart';
import 'widgets/balance_tile.dart';

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
  late final HomeViewmodel _viewModel;

  @override
  void initState() {
    super.initState();

    _viewModel = widget.viewModel;

    _viewModel.initialize.execute();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Minha conta'),
        actions: [
          IconButton(
            onPressed: () => _viewModel.initialize.execute(),
            icon: const Icon(Icons.refresh_rounded),
          ),
        ],
      ),
      body: ListenableBuilder(
        listenable: _viewModel.initialize,
        builder: (context, _) => Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              ListenableBuilder(
                listenable: _viewModel.initialize,
                builder: (context, _) => BalanceTile(
                  balanceLabel: _balanceLabel(),
                  supportingLabel: _supportingLabel(),
                  isLoading: _viewModel.initialize.isRunning,
                ),
              ),

              Padding(
                padding: const EdgeInsets.only(bottom: 12, top: 20),
                child: Text(
                  'Acessos rápidos',
                  style: Theme.of(context).textTheme.titleMedium,
                ),
              ),

              Row(
                children: [
                  Expanded(
                    child: ActionTile(
                      icon: Icons.receipt_long_rounded,
                      title: 'Extrato',
                      subtitle: 'Disponibilizar em breve',
                      onTap: () => _showPendingFeature('Extrato'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: ActionTile(
                      icon: Icons.swap_horiz_rounded,
                      title: 'Transferir',
                      subtitle: 'Disponibilizar em breve',
                      onTap: () => _showPendingFeature('Transferir'),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  String _balanceLabel() {
    final balance = widget.viewModel.lastBalance;
    if (balance == null) return 'R\$ --';

    final value = (balance.balance / 100)
        .toStringAsFixed(2)
        .replaceAll('.', ',');
    return 'R\$ $value';
  }

  void _showPendingFeature(String featureName) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('$featureName sera disponibilizado em breve.')),
    );
  }

  String _supportingLabel() {
    if (_viewModel.initialize.isRunning) return 'Iniciando...';

    return 'Conta: ${_viewModel.selectedAccount!.branch} - ${_viewModel.selectedAccount!.number}';
  }
}
