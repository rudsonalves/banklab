// import 'package:dio/dio.dart';

// import '/core/services/security/zero_thrust/device/device_identity_service.dart';

// class DeviceInterceptor extends Interceptor {
//   final DeviceIdentityService _deviceService;

//   DeviceInterceptor(this._deviceService);

//   @override
//   Future<void> onRequest(
//     RequestOptions options,
//     RequestInterceptorHandler handler,
//   ) async {
//     final result = await _deviceService.getDeviceId();

//     result.fold(
//       onSuccess: (deviceId) {
//         if (deviceId != null && deviceId.isNotEmpty) {
//           options.headers['X-Device-Id'] = deviceId;
//           options.headers['X-Device-Type'] = _deviceService.deviceType;
//         }
//       },
//       onFailure: (_) {
//         // Silent failure intentional
//         // Does not block the request
//       },
//     );

//     handler.next(options);
//   }
// }
