# Error and Response Standard — Bank API

## 1. Overview

This document defines the standard for:

* success responses
* error responses
* message structure
* mapping to HTTP status codes

The goal is to ensure:

* consistency across endpoints
* predictability for clients
* ease of error handling

---

## 2. Principles

### 2.1 Simplicity

* single response structure
* no unnecessary fields
* easy to consume

---

### 2.2 Consistency

* same format across all endpoints
* same codes for the same types of errors

---

### 2.3 Clarity

* objective messages
* explicit error codes

---

## 3. Response Structure

### 3.1 Success Response

```json
{
  "data": { ... },
  "error": null
}
```

---

### 3.2 Error Response

```json
{
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": {}
  }
}
```

---

## 4. Fields

### 4.1 data

* contains the operation result
* may be an object or a list

---

### 4.2 error.code

Stable identifier of the error.

Examples:

* CUSTOMER_ALREADY_EXISTS
* ACCOUNT_NOT_FOUND
* INSUFFICIENT_FUNDS
* INVALID_AMOUNT

---

### 4.3 error.message

Human-readable message.

* does not need to be localized at this stage
* must not contain sensitive details

---

### 4.4 error.details (optional)

Additional structured information.

Example:

```json
{
  "field": "amount",
  "reason": "must be greater than zero"
}
```

---

## 5. HTTP Status Codes

### 5.1 Success

| Code | Usage                |
| ---- | -------------------- |
| 200  | successful operation |
| 201  | resource created     |

---

### 5.2 Client Errors

| Code | Usage                        |
| ---- | ---------------------------- |
| 400  | invalid request              |
| 404  | resource not found           |
| 409  | conflict (e.g., duplication) |
| 422  | business rule violation      |

---

### 5.3 Server Errors

| Code | Usage          |
| ---- | -------------- |
| 500  | internal error |

---

## 6. Domain Error Mapping

| Domain Error            | HTTP | Code                    |
| ----------------------- | ---- | ----------------------- |
| customer already exists | 409  | CUSTOMER_ALREADY_EXISTS |
| account not found       | 404  | ACCOUNT_NOT_FOUND       |
| insufficient funds      | 422  | INSUFFICIENT_FUNDS      |
| invalid amount          | 400  | INVALID_AMOUNT          |
| inactive account        | 422  | ACCOUNT_INACTIVE        |
| duplicate operation     | 409  | DUPLICATE_REQUEST       |

---

## 7. Usage Rules

### 7.1 Never mix success and error

* if error → `data = null`
* if success → `error = null`

---

### 7.2 Do not expose internal details

Avoid:

* stack traces
* SQL statements
* infrastructure details

---

### 7.3 Code standardization

* codes must use UPPER_SNAKE_CASE
* must remain stable over time

---

### 7.4 Messages are not a contract

* clients must rely on `error.code`
* must not depend on message text

---

## 8. Examples

---

### 8.1 Success — Deposit

```json
{
  "data": {
    "account_id": "uuid",
    "balance": 150000
  },
  "error": null
}
```

---

### 8.2 Error — Insufficient Funds

```json
{
  "data": null,
  "error": {
    "code": "INSUFFICIENT_FUNDS",
    "message": "Insufficient balance",
    "details": {}
  }
}
```

---

### 8.3 Error — Invalid Field

```json
{
  "data": null,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "Amount must be greater than zero",
    "details": {
      "field": "amount"
    }
  }
}
```

---

## 9. Intentional Decisions

### 9.1 Single structure

Avoids multiple response formats.

---

### 9.2 Code separated from message

Enables:

* future internationalization
* client stability

---

### 9.3 Use of 422

Used for business rule violations:

* insufficient funds
* inactive account

---

## 10. Limitations

* no error versioning
* no multi-language support
* no correlation with logs (yet)

---

## 11. Future Evolution

Possible improvements:

* add `request_id`
* standardize errors by domain
* multi-language support
* centralized error catalog

---

## 12. Conclusion

The defined standard:

* is simple
* is consistent
* is sufficient for future evolution

It ensures predictability without introducing unnecessary complexity.
