typedef AsyncResult<T extends Object> = Future<Result<T>>;

final class Unit {
  const Unit._();
}

const unit = Unit._();

sealed class Result<T extends Object> {
  const Result();

  bool get isSuccess => this is Success<T>;
  bool get isFailure => this is Failure<T>;

  const factory Result.success(T value) = Success<T>;
  const factory Result.failure(Exception error) = Failure<T>;

  T? get value => switch (this) {
    Success(:final value) => value,
    _ => null,
  };

  Exception? get error => switch (this) {
    Failure(:final error) => error,
    _ => null,
  };

  R fold<R>({
    required R Function(T value) onSuccess,
    required R Function(Exception error) onFailure,
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
  final Exception _error;

  @override
  Exception get error => _error;

  const Failure(this._error);
}
