// mobile/lib/core/resources/money_currency.dart

import 'package:money2/money2.dart';

final appCurrency = AppCurrencies.defaultCurrency;

/// Central registry for supported currencies.
///
/// This application currently operates in BRL.
/// Other currencies are available for future/shared-library compatibility.
abstract final class AppCurrencies {
  const AppCurrencies._();

  static final Currency brl = Currency.create(
    'BRL',
    2,
    symbol: 'R\$',
    pattern: 'S #,##0.00',
    groupSeparator: '.',
    decimalSeparator: ',',
    country: 'BR',
  );

  static final Currency usd = Currency.create(
    'USD',
    2,
    symbol: '\$',
    pattern: 'S #,##0.00',
    groupSeparator: ',',
    decimalSeparator: '.',
    country: 'US',
  );

  static final Currency eur = Currency.create(
    'EUR',
    2,
    symbol: '€',
    pattern: 'S #,##0.00',
    groupSeparator: '.',
    decimalSeparator: ',',
    country: 'EU',
  );

  static final Currency btc = Currency.create(
    'BTC',
    8,
    symbol: '₿',
    pattern: 'S #,##0.00000000',
    groupSeparator: ',',
    decimalSeparator: '.',
    country: 'BTC',
  );

  static Currency get defaultCurrency => brl;
}
