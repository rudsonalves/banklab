import 'dart:io';
import 'dart:ui' as ui;

import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'package:path_provider/path_provider.dart';
import 'package:share_plus/share_plus.dart';

import '/core/extensions/datetime_extension.dart';
import '/core/routing/routes.dart';
import '/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import '/domain/common/receipt/enums/transfer_receipt_status.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/buttons/big_button.dart';
import '/uis/core/buttons/big_text_button.dart';
import '/uis/core/messages/app_snackbar.dart';
import 'viewmodel/details_viewmodel.dart';
import 'widgets/detail_line.dart';

class ReceiptImageException implements Exception {
  final String message;

  const ReceiptImageException(this.message);

  @override
  String toString() => message;
}

class DetailsPage extends StatefulWidget {
  final DetailsViewmodel viewModel;
  final String reference;

  const DetailsPage({
    super.key,
    required this.viewModel,
    required this.reference,
  });

  @override
  State<DetailsPage> createState() => _DetailsPageState();
}

class _DetailsPageState extends State<DetailsPage> {
  late final DetailsViewmodel _viewModel;
  final GlobalKey _receiptBoundaryKey = GlobalKey();

  bool _isSharing = false;

  @override
  void initState() {
    super.initState();

    _viewModel = widget.viewModel;
    _viewModel.getTransferReceipt.addListener(_onGetTransferReceiptChanged);
    _viewModel.getTransferReceipt.execute(widget.reference);
  }

  @override
  void dispose() {
    _viewModel.getTransferReceipt.removeListener(_onGetTransferReceiptChanged);

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Detalhes da transferência'),
      ),

      body: AnimatedBuilder(
        animation: _viewModel.getTransferReceipt,
        builder: (context, _) {
          if (_viewModel.getTransferReceipt.isRunning &&
              _viewModel.getTransferReceipt.value == null) {
            return const Center(
              child: CircularProgressIndicator(),
            );
          }

          final receipt = _viewModel.getTransferReceipt.value;
          if (receipt == null) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Text('Não foi possível carregar o comprovante.'),
                    const SizedBox(height: 12),
                    FilledButton(
                      onPressed: () => _viewModel.getTransferReceipt.execute(
                        widget.reference,
                      ),
                      child: const Text('Tentar novamente'),
                    ),
                  ],
                ),
              ),
            );
          }

          final theme = Theme.of(context);

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Center(
              child: ConstrainedBox(
                constraints: const BoxConstraints(maxWidth: 520),
                child: RepaintBoundary(
                  key: _receiptBoundaryKey,
                  child: Card(
                    elevation: 0,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                      side: BorderSide(
                        color: theme.colorScheme.outlineVariant,
                      ),
                    ),
                    child: Padding(
                      padding: const EdgeInsets.all(20),
                      child: Column(
                        spacing: 4,
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          Text(
                            'Comprovante de transferência',
                            style: theme.textTheme.titleLarge,
                          ),

                          Text(
                            receipt.operationDate.format(
                              context,
                              'dd/MM/yyyy HH:mm',
                            ),
                            style: theme.textTheme.bodyMedium?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
                              fontWeight: FontWeight.w700,
                            ),
                          ),

                          const SizedBox(height: 12),

                          DetailLine(
                            label: 'Valor',
                            value: receipt.amount.format(),
                            valueStyle: theme.textTheme.titleLarge,
                          ),
                          DetailLine(
                            label: 'Status',
                            value: _statusLabel(receipt.status),
                            valueColor: _statusColor(theme, receipt.status),
                          ),
                          DetailLine(
                            label: 'Tipo de operação',
                            value: receipt.operationType,
                          ),
                          DetailLine(
                            label: 'Remetente',
                            value:
                                '${receipt.sourceBranch} - ${receipt.sourceAccountNumber}',
                          ),
                          DetailLine(
                            label: 'Destinatário',
                            value: receipt.recipientName,
                          ),
                          DetailLine(
                            label: 'Conta destino',
                            value:
                                '${receipt.destinationBranch} - ${receipt.destinationAccountNumber}',
                          ),
                          if (receipt.description != null &&
                              receipt.description!.trim().isNotEmpty)
                            DetailLine(
                              label: 'Descrição',
                              value: receipt.description!,
                            ),

                          const Divider(height: 28),

                          Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      'Referência da transação',
                                      style: theme.textTheme.labelLarge
                                          ?.copyWith(
                                            color: theme
                                                .colorScheme
                                                .onSurfaceVariant,
                                          ),
                                    ),
                                    const SizedBox(height: 2),
                                    Text(
                                      receipt.transactionReference,
                                      style: theme.textTheme.bodyLarge
                                          ?.copyWith(
                                            fontWeight: FontWeight.w700,
                                          ),
                                    ),
                                  ],
                                ),
                              ),
                              IconButton(
                                onPressed: () => _copyReference(
                                  receipt.transactionReference,
                                ),
                                tooltip: 'Copiar referência',
                                visualDensity: VisualDensity.compact,
                                iconSize: 18,
                                icon: const Icon(Icons.copy_rounded),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            ),
          );
        },
      ),

      bottomNavigationBar: AnimatedBuilder(
        animation: _viewModel.getTransferReceipt,
        builder: (context, _) {
          final receipt = _viewModel.getTransferReceipt.value;
          final canShare =
              receipt != null && !_viewModel.getTransferReceipt.isRunning;

          return Row(
            spacing: 12,
            children: [
              Expanded(
                flex: 1,
                child: BigTextButton(
                  onPressed: () => _navBack(context),
                  label: 'Fechar',
                ),
              ),
              Expanded(
                flex: 2,
                child: BigButton(
                  onPressed: canShare && !_isSharing
                      ? () => _shareReceipt(receipt)
                      : null,
                  label: _isSharing ? 'Gerando imagem...' : 'Compartilhar',
                  rightIcon: Icons.share_rounded,
                ),
              ),
            ],
          );
        },
      ),
    );
  }

  void _navBack(BuildContext context) {
    context.goNamed(HomeRoutes.home.name);
  }

  void _onGetTransferReceiptChanged() {
    if (!mounted || _viewModel.getTransferReceipt.isRunning) return;

    if (_viewModel.getTransferReceipt.isFailure) {
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message:
            _viewModel.getTransferReceipt.error?.message ??
            'Não foi possível carregar o comprovante.',
      );
    }
  }

  Future<void> _copyReference(String reference) async {
    await Clipboard.setData(ClipboardData(text: reference));
    if (!mounted) return;

    AppSnackbar.show(
      context,
      type: SnackbarType.info,
      message: 'Referência copiada para a área de transferência.',
    );
  }

  Future<void> _shareReceipt(TransferReceiptResponseDto receipt) async {
    if (_isSharing) return;

    final box = context.findRenderObject() as RenderBox?;
    final shareOrigin = box == null
        ? null
        : box.localToGlobal(Offset.zero) & box.size;

    setState(() => _isSharing = true);

    try {
      final file = await _buildReceiptImage(receipt.transactionReference);

      await SharePlus.instance.share(
        ShareParams(
          files: [XFile(file.path)],
          text: 'Comprovante de transferência ${receipt.transactionReference}',
          sharePositionOrigin: shareOrigin,
        ),
      );

      if (!mounted) return;
      AppSnackbar.show(
        context,
        type: SnackbarType.success,
        message: 'Comprovante pronto para compartilhamento.',
      );
    } on ReceiptImageException catch (error) {
      if (!mounted) return;
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message: error.message,
      );
    } on PlatformException {
      if (!mounted) return;
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message: 'Não foi possível abrir o compartilhamento do comprovante.',
      );
    } catch (_) {
      if (!mounted) return;
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message: 'Não foi possível compartilhar o comprovante.',
      );
    } finally {
      if (mounted) {
        setState(() => _isSharing = false);
      }
    }
  }

  Future<File> _buildReceiptImage(String transactionReference) async {
    final deviceRatio = View.of(context).devicePixelRatio;
    await WidgetsBinding.instance.endOfFrame;
    await Future<void>.delayed(const Duration(milliseconds: 16));

    RenderRepaintBoundary render = _receiptBoundary();
    for (var attempts = 0; render.debugNeedsPaint && attempts < 5; attempts++) {
      await WidgetsBinding.instance.endOfFrame;
      render = _receiptBoundary();
    }

    if (render.debugNeedsPaint) {
      throw const ReceiptImageException(
        'Comprovante ainda não terminou de renderizar para captura.',
      );
    }

    if (!render.hasSize || render.size.isEmpty) {
      throw const ReceiptImageException(
        'Comprovante ainda não foi renderizado para captura.',
      );
    }

    final preferredRatio = deviceRatio > 2.0 ? 2.0 : deviceRatio;

    ui.Image image;
    try {
      image = await render.toImage(pixelRatio: preferredRatio);
    } catch (_) {
      image = await render.toImage(pixelRatio: 1.0);
    }

    final bytes = await image.toByteData(format: ui.ImageByteFormat.png);

    if (bytes == null) {
      throw const ReceiptImageException(
        'Falha ao gerar a imagem do comprovante.',
      );
    }

    final tempDir = await getTemporaryDirectory();
    final normalizedRef = transactionReference.replaceAll(
      RegExp(r'[^a-zA-Z0-9_-]'),
      '_',
    );
    final file = File('${tempDir.path}/comprovante_$normalizedRef.png');

    await file.writeAsBytes(bytes.buffer.asUint8List(), flush: true);

    if (!await file.exists()) {
      throw const ReceiptImageException(
        'Não foi possível salvar a imagem do comprovante.',
      );
    }

    return file;
  }

  RenderRepaintBoundary _receiptBoundary() {
    final renderObject = _receiptBoundaryKey.currentContext?.findRenderObject();
    if (renderObject is! RenderRepaintBoundary) {
      throw const ReceiptImageException(
        'Comprovante indisponível para captura.',
      );
    }

    return renderObject;
  }

  String _statusLabel(TransferReceiptStatus status) {
    return switch (status) {
      TransferReceiptStatus.completed => 'Concluida',
      TransferReceiptStatus.pending => 'Pendente',
      TransferReceiptStatus.failed => 'Falhou',
      TransferReceiptStatus.cancelled => 'Cancelada',
      TransferReceiptStatus.rejected => 'Rejeitada',
    };
  }

  Color _statusColor(ThemeData theme, TransferReceiptStatus status) {
    return switch (status) {
      TransferReceiptStatus.completed => theme.colorScheme.primary,
      TransferReceiptStatus.pending => theme.colorScheme.tertiary,
      TransferReceiptStatus.failed ||
      TransferReceiptStatus.cancelled ||
      TransferReceiptStatus.rejected => theme.colorScheme.error,
    };
  }
}
