import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '/core/routing/routes.dart';
import '/data/services/apis/account/dtos/statement_query_params_dto.dart';
import '/data/services/apis/account/dtos/statement_response_dto.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/messages/app_snackbar.dart';
import 'viewmodel/statement_viewmodel.dart';
import 'widgets/day_header.dart';
import 'widgets/month_header.dart';
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
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.error_outline,
                    size: 64,
                    color: greyColor,
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Não foi possível carregar o extrato',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  const SizedBox(height: 24),
                  ElevatedButton(
                    onPressed: isLoading ? null : _loadStatement,
                    child: isLoading
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                            ),
                          )
                        : const Text('Tentar novamente'),
                  ),
                ],
              ),
            );
          }

          final visibleStatement = statement;
          if (visibleStatement == null) {
            return const SizedBox.shrink();
          }

          if (visibleStatement.items.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.receipt_long_rounded,
                    size: 64,
                    color: greyColor,
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Nenhuma transação encontrada',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ],
              ),
            );
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
    final groupedByMonthDay = <String, Map<String, List<StatementItemDto>>>{};
    final lastOperationByDay = <String, StatementItemDto>{};

    for (final item in items) {
      final date = _tryParseDate(item.createdAt);
      final monthKey = DateFormat('yyyy-MM').format(date);
      final dayKey = DateFormat('yyyy-MM-dd').format(date);

      groupedByMonthDay.putIfAbsent(
        monthKey,
        () => <String, List<StatementItemDto>>{},
      );
      groupedByMonthDay[monthKey]!.putIfAbsent(
        dayKey,
        () => <StatementItemDto>[],
      );
      groupedByMonthDay[monthKey]![dayKey]!.add(item);

      final currentLast = lastOperationByDay[dayKey];
      if (currentLast == null ||
          _tryParseDate(currentLast.createdAt).isBefore(date)) {
        lastOperationByDay[dayKey] = item;
      }
    }

    final monthKeys = groupedByMonthDay.keys.toList()
      ..sort((a, b) => b.compareTo(a));
    final children = <Widget>[];

    for (final monthKey in monthKeys) {
      children.add(MonthHeader(label: _formatMonthLabel(monthKey)));

      final dayMap = groupedByMonthDay[monthKey]!;
      final dayKeys = dayMap.keys.toList()..sort((a, b) => b.compareTo(a));

      for (final dayKey in dayKeys) {
        final dayBalance = lastOperationByDay[dayKey]?.balanceAfter;
        children.add(
          DayHeader(
            label: _formatDayLabel(dayKey),
            greyColor: greyColor,
            balance: dayBalance,
          ),
        );

        final dayItems = dayMap[dayKey]!;
        dayItems.sort(
          (a, b) =>
              _tryParseDate(a.createdAt).compareTo(_tryParseDate(b.createdAt)),
        );

        // final lastOperation = lastOperationByDay[dayKey];
        for (final item in dayItems) {
          children.add(
            StatementItemCard(
              item: item,
              // showConsolidatedBalance:
              //     lastOperation != null &&
              //     item.transactionId == lastOperation.transactionId,
              hourLabel: _formatHour(item.createdAt),
              onTap: () => _openDetails(context, item),
            ),
          );
        }
      }
    }

    return children;
  }

  String _formatMonthLabel(String monthKey) {
    try {
      final date = DateTime.parse('$monthKey-01');
      return DateFormat('MMMM yyyy').format(date);
    } catch (_) {
      return monthKey;
    }
  }

  String _formatDayLabel(String dayKey) {
    try {
      final date = DateTime.parse(dayKey);
      return DateFormat('dd/MM/yyyy').format(date);
    } catch (_) {
      return dayKey;
    }
  }

  String _formatHour(String createdAt) {
    try {
      final date = DateTime.parse(createdAt);
      return DateFormat('HH:mm').format(date);
    } catch (_) {
      return createdAt;
    }
  }

  DateTime _tryParseDate(String input) {
    return DateTime.tryParse(input) ?? DateTime.fromMillisecondsSinceEpoch(0);
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
