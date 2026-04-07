typedef AsyncResult<T extends Object> = Future<Result<T>>;

final class Unit {
  const Unit._();
}

const unit = Unit._();

sealed class Result<T extends Object> {
  const Result();

  T? get value => null;
  Exception? get error => null;

  bool get isSuccess => value != null;
  bool get isFailure => error != null;

  const factory Result.success(T value) = Success._;
  const factory Result.failure(Exception error) = Failure._;

  R fold<R>({
    required R Function(T value) onSuccess,
    required R Function(Exception error) onFailure,
  }) {
    if (isSuccess) {
      return onSuccess(value!);
    } else {
      return onFailure(error!);
    }
  }
}

final class Success<T extends Object> extends Result<T> {
  final T _value;

  const Success._(this._value);

  @override
  T get value => _value;
}

final class Failure<T extends Object> extends Result<T> {
  final Exception _error;

  const Failure._(this._error);

  @override
  Exception get error => _error;
}
