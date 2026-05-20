import '/core/result/result.dart';
import '/data/services/apis/contact_verification/contact_verification_api.dart';
import '/data/services/apis/contact_verification/dtos/contact_verification_confirm_request_dto.dart';
import '/data/services/apis/contact_verification/dtos/contact_verification_confirm_response_dto.dart';
import '/data/services/apis/contact_verification/dtos/contact_verification_request_dto.dart';
import '/data/services/apis/contact_verification/dtos/contact_verification_request_response_dto.dart';
import 'contact_verification_repository.dart';

class ContactVerificationRepositoryImpl
    implements ContactVerificationRepository {
  final ContactVerificationApi _api;

  ContactVerificationRepositoryImpl({
    required ContactVerificationApi api,
  }) : _api = api;

  @override
  AsyncResult<ContactVerificationRequestResponseDto> requestContactVerification(
    ContactVerificationRequestDto dto,
  ) async {
    final result = await _api.requestContactVerification(dto);
    if (result.isFailure) return Result.failure(result.error!);
    return Success(result.value!);
  }

  @override
  AsyncResult<ContactVerificationConfirmResponseDto> confirmContactVerification(
    ContactVerificationConfirmRequestDto dto,
  ) async {
    final result = await _api.confirmContactVerification(dto);
    if (result.isFailure) return Result.failure(result.error!);
    return Success(result.value!);
  }
}
