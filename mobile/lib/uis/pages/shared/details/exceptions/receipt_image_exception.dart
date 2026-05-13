class ReceiptImageException implements Exception {
  final String message;

  const ReceiptImageException(this.message);

  @override
  String toString() => message;
}
