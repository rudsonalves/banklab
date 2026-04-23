import 'package:flutter/material.dart';

import '/core/routing/route_observer.dart';
import '../../core/base/safe_scaffold.dart';
import '../../core/feedback/app_snackbar.dart';
import 'viewmodel/home_viewmodel.dart';
import 'widgets/action_tite.dart';
import 'widgets/balance_tile.dart';

class HomePage extends StatefulWidget {
  final HomeViewmodel viewModel;

  const HomePage({
    super.key,
    required this.viewModel,
  });

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> with RouteAware {
  late final HomeViewmodel _viewModel;
  bool _isRouteObserverSubscribed = false;

  @override
  void initState() {
    super.initState();

    _viewModel = widget.viewModel;
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();

    if (_isRouteObserverSubscribed) return;

    final route = ModalRoute.of(context);
    if (route == null) return;

    routeObserver.subscribe(this, route);
    _isRouteObserverSubscribed = true;
  }

  @override
  void dispose() {
    routeObserver.unsubscribe(this);
    _viewModel.dispose();

    super.dispose();
  }

  @override
  void didPush() {
    _viewModel.initialize.execute();
  }

  @override
  void didPop() {
    _viewModel.stopTimer();
  }

  @override
  void didPushNext() {
    _viewModel.stopTimer();
  }

  @override
  void didPopNext() {
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
                  balance: _viewModel.lastBalance,
                  account: _viewModel.selectedAccount,
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

  void _showPendingFeature(String featureName) {
    AppSnackbar.show(
      context,
      title: 'Em breve',
      message: '$featureName será disponibilizado em breve.',
    );
  }
}
