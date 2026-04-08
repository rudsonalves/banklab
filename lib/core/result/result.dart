import 'errors/app_error.dart';

export 'errors/app_error.dart';
export 'unit.dart';

typedef AsyncResult<T extends Object> = Future<Result<T>>;

sealed class Result<T extends Object> {
  const Result();

  bool get isSuccess => this is Success<T>;
  bool get isFailure => this is Failure<T>;

  const factory Result.success(T value) = Success<T>;
  const factory Result.failure(AppError error) = Failure<T>;

  T? get value => switch (this) {
    Success(:final value) => value,
    _ => null,
  };

  AppError? get error => switch (this) {
    Failure(:final error) => error,
    _ => null,
  };

  R fold<R>({
    required R Function(T value) onSuccess,
    required R Function(AppError error) onFailure,
  }) {
    return switch (this) {
      Success(:final value) => onSuccess(value),
      Failure(:final error) => onFailure(error),
    };
  }
}

final class Success<T extends Object> extends Result<T> {
  final T _value;

  @override
  T get value => _value;

  const Success(this._value);
}

final class Failure<T extends Object> extends Result<T> {
  final AppError _error;

  @override
  AppError get error => _error;

  const Failure(this._error);
}
