package errors

// NewAuthError 
func NewAuthError(message string) *AppError {
	return NewAppError(AuthenticationError, message, "")
}

// NewPermissionError 
func NewPermissionError(message string) *AppError {
	return NewAppError(AuthorizationError, message, "")
}

// NewRateLimitError 
func NewRateLimitError(message string) *AppError {
	return NewAppError(RateLimitError, message, "")
}

// NewNotFoundError 
func NewNotFoundError(message string) *AppError {
	return NewAppError(NotFoundError, message, "")
}

// NewServerError 
func NewServerError(message string) *AppError {
	return NewAppError(InternalServerError, message, "")
}

// NewValidationError 
func NewValidationError(message string) *AppError {
	return NewAppError(ValidationError, message, "")
}