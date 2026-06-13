enum ContactVerificationChannel {
  email,
  phone;

  factory ContactVerificationChannel.fromString(String value) {
    switch (value.trim().toLowerCase()) {
      case 'email':
        return ContactVerificationChannel.email;
      case 'phone':
        return ContactVerificationChannel.phone;
      default:
        throw ArgumentError('Invalid contact verification channel: $value');
    }
  }
}
