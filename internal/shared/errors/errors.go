package sharederrors

type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}

func NewError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewErrorWithDetails(code, message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

var (
	ErrInvalidRequest     = NewError("INVALID_REQUEST", "Invalid request body")
	ErrInvalidData        = NewError("INVALID_DATA", "Invalid data")
	ErrInternal           = NewError("INTERNAL_ERROR", "Internal error")
	ErrUserAlreadyExists  = NewError("USER_ALREADY_EXISTS", "User already exists")
	ErrInvalidCredentials = NewError("INVALID_CREDENTIALS", "Invalid credentials")
	ErrUnauthorized       = NewError("UNAUTHORIZED", "Unauthorized")
	ErrInvalidToken       = NewError("INVALID_TOKEN", "Invalid token")
)
