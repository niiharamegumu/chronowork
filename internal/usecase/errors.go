package usecase

// ErrorCode represents the type of error that occurred.
type ErrorCode string

const (
	// ErrCodeNotFound indicates the requested resource was not found.
	ErrCodeNotFound ErrorCode = "NOT_FOUND"
	// ErrCodeDuplicateToday indicates a duplicate entry for today.
	ErrCodeDuplicateToday ErrorCode = "DUPLICATE_TODAY"
	// ErrCodeValidation indicates a validation error.
	ErrCodeValidation ErrorCode = "VALIDATION"
	// ErrCodePermission indicates a permission error.
	ErrCodePermission ErrorCode = "PERMISSION"
)

// UseCaseError represents an error from the use case layer.
type UseCaseError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface.
func (e *UseCaseError) Error() string {
	return e.Message
}

// Unwrap returns the wrapped error.
func (e *UseCaseError) Unwrap() error {
	return e.Err
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(message string) error {
	return &UseCaseError{
		Code:    ErrCodeNotFound,
		Message: message,
	}
}

// NewDuplicateError creates a new duplicate error.
func NewDuplicateError(message string) error {
	return &UseCaseError{
		Code:    ErrCodeDuplicateToday,
		Message: message,
	}
}

// NewValidationError creates a new validation error.
func NewValidationError(message string) error {
	return &UseCaseError{
		Code:    ErrCodeValidation,
		Message: message,
	}
}

// NewPermissionError creates a new permission error.
func NewPermissionError(message string) error {
	return &UseCaseError{
		Code:    ErrCodePermission,
		Message: message,
	}
}
