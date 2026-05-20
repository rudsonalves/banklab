import '/core/result/result.dart';
import '/data/services/apis/contact_verification/dtos/contact_verification_confirm_request_dto.dart';
import '/data/services/apis/contact_verification/dtos/contact_verification_confirm_response_dto.dart';
import '/data/services/apis/contact_verification/dtos/contact_verification_request_dto.dart';
import '/data/services/apis/contact_verification/dtos/contact_verification_request_response_dto.dart';

abstract class ContactVerificationRepository {
  /// Requests a contact verification code for email or phone.
  AsyncResult<ContactVerificationRequestResponseDto> requestContactVerification(
    ContactVerificationRequestDto dto,
  );

  /// Confirms a contact verification token and returns the verified token.
  AsyncResult<ContactVerificationConfirmResponseDto> confirmContactVerification(
    ContactVerificationConfirmRequestDto dto,
  );
}
