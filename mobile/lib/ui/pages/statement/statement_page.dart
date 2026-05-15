import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/extensions/datetime_extension.dart';
import '/core/routing/routes.dart';
import '/data/services/apis/account/dtos/statement_query_params_dto.dart';
import '/data/services/apis/account/dtos/statement_response_dto.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/messages/app_snackbar.dart';
import 'viewmodel/statement_viewmodel.dart';
import 'widgets/day_header.dart';
import 'widgets/load_statement_error.dart';
import 'widgets/month_header.dart';
import 'widgets/no_transactions_card.dart';
import 'widgets/statement_item_card.dart';

class StatementPage extends StatefulWidget {
  final StatementViewmodel viewModel;

  const StatementPage({super.key, required this.viewModel});

  @override
  State<StatementPage> createState() => _StatementPageState();
}

class _StatementPageState extends State<StatementPage> {
  StatementViewmodel get _viewModel => widget.viewModel;

  @override
  void initState() {
    super.initState();

    _viewModel.getStatement.addListener(_onGetStatementChanged);
    _loadStatement();
  }

  @override
  void dispose() {
    _viewModel.getStatement.removeListener(_onGetStatementChanged);

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final greyColor = Theme.of(context).colorScheme.onSurface;

    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Extrato'),
      ),
      body: AnimatedBuilder(
        animation: _viewModel.getStatement,
        builder: (context, _) {
          final isLoading = _viewModel.getStatement.isRunning;
          final isFailure = _viewModel.getStatement.isFailure;
          final statement =
              _viewModel.getStatement.value ?? _viewModel.lastStatement;

          if (isLoading && statement == null) {
            return const Center(
              child: CircularProgressIndicator(),
            );
          }

          if ((isFailure || statement == null) &&
              _viewModel.lastStatement == null) {
            return LoadStatementError(
              isLoading: isLoading,
              onRetry: _loadStatement,
            );
          }

          final visibleStatement = statement;
          if (visibleStatement == null) {
            return const SizedBox.shrink();
          }

          if (visibleStatement.items.isEmpty) {
            return NoTransactionsCard(greyColor: greyColor);
          }

          return RefreshIndicator(
            onRefresh: () async {
              _loadStatement();
            },
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: _buildGroupedStatement(
                context,
                visibleStatement.items,
                greyColor,
              ),
            ),
          );
        },
      ),
    );
  }

  List<Widget> _buildGroupedStatement(
    BuildContext context,
    List<StatementItemDto> items,
    Color greyColor,
  ) {
    final groupedByMonthDay =
        <DateTime, Map<DateTime, List<StatementItemDto>>>{};
    final lastOperationByDay = <DateTime, StatementItemDto>{};

    for (final item in items) {
      final date = item.createdAt;
      final monthKey = DateTime(date.year, date.month);
      final dayKey = DateTime(date.year, date.month, date.day);

      groupedByMonthDay.putIfAbsent(
        monthKey,
        () => <DateTime, List<StatementItemDto>>{},
      );
      groupedByMonthDay[monthKey]!.putIfAbsent(
        dayKey,
        () => <StatementItemDto>[],
      );
      groupedByMonthDay[monthKey]![dayKey]!.add(item);

      final currentLast = lastOperationByDay[dayKey];
      if (currentLast == null || currentLast.createdAt.isBefore(date)) {
        lastOperationByDay[dayKey] = item;
      }
    }

    final monthKeys = groupedByMonthDay.keys.toList()
      ..sort((a, b) => b.compareTo(a));
    final children = <Widget>[];

    for (final monthKey in monthKeys) {
      children.add(MonthHeader(label: monthKey.formatMonthLabel));

      final dayMap = groupedByMonthDay[monthKey]!;
      final dayKeys = dayMap.keys.toList()..sort((a, b) => b.compareTo(a));

      for (final dayKey in dayKeys) {
        final dayBalance = lastOperationByDay[dayKey]?.balanceAfter;
        children.add(
          DayHeader(
            label: dayKey.formatDayLabel,
            greyColor: greyColor,
            balance: dayBalance,
          ),
        );

        final dayItems = dayMap[dayKey]!;
        dayItems.sort((a, b) => a.createdAt.compareTo(b.createdAt));

        for (final item in dayItems) {
          children.add(
            StatementItemCard(
              item: item,
              hourLabel: item.createdAt.formatHour,
              onTap: () => _openDetails(context, item),
            ),
          );
        }
      }
    }

    return children;
  }

  void _openDetails(BuildContext context, StatementItemDto item) {
    final reference = (item.referenceId ?? item.transactionId).trim();

    if (reference.isEmpty) {
      AppSnackbar.show(
        context,
        type: SnackbarType.info,
        title: 'Aviso',
        message: 'Esta movimentação não possui referência para detalhamento.',
      );
      return;
    }

    context.pushNamed(
      SharedRoutes.details.name,
      extra: reference,
    );
  }

  void _loadStatement() {
    _viewModel.getStatement.execute(const StatementQueryParamsDto());
  }

  void _onGetStatementChanged() {
    if (!mounted || _viewModel.getStatement.isRunning) return;

    if (_viewModel.getStatement.isFailure) {
      final message =
          _viewModel.getStatement.error?.message ??
          'Falha ao carregar o extrato.';
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message: message,
      );
    }
  }
}
