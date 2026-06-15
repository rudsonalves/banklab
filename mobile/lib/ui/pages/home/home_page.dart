import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/route_observer.dart';
import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/cards/balance_card.dart';
import '../../../core/routing/models/transaction_password_setup_origin.dart';
import '../../../data/services/apis/transaction_password/enums/transaction_password_status.dart';
import '../../components/messages/app_snackbar.dart';
import 'viewmodel/home_viewmodel.dart';
import 'widgets/action_tite.dart';

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
    _viewModel.initialize.execute();
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
              BalanceCard(
                balance: _viewModel.balance,
                isVisible: true,
                onToggleVisibility: () {},
                selectedAccount: _viewModel.selectedAccount,
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
                      onTap: _navToStatement,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: ActionTile(
                      icon: Icons.swap_horiz_rounded,
                      title: 'Transferir',
                      subtitle: 'Disponibilizar em breve',
                      onTap: _navToTransferRecipient,
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

  void _navToStatement() {
    context.pushNamed(BaseRoutes.statement.routeName);
  }

  void _navToTransferRecipient() {
    switch (_viewModel.transactionPasswordStatus) {
      case TransactionPasswordStatus.notSet:
        context.pushNamed(
          TransactionPasswordRoutes.introduction.routeName,
          extra: TransactionPasswordSetupOrigin.transfer,
        );
        break;

      case TransactionPasswordStatus.active:
        context.pushNamed(TransferRoutes.recipient.routeName);
        break;

      case TransactionPasswordStatus.locked:
        AppSnackbar.show(
          context,
          type: SnackbarType.error,
          title: 'Acesso bloqueado',
          message:
              'Sua senha transacional está bloqueada devido a múltiplas'
              ' tentativas incorretas. Por favor, entre em contato com o'
              ' suporte para desbloqueio.',
        );
        break;

      case TransactionPasswordStatus.unknown:
        AppSnackbar.show(
          context,
          type: SnackbarType.error,
          title: 'Erro',
          message:
              'Não foi possível verificar o status da sua senha transacional.'
              ' Por favor, tente novamente mais tarde.',
        );
        break;
    }
  }
}
