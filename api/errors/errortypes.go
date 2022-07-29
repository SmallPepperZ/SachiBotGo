package errors

type ErrorMissingPermission struct {
	Permission string
}

func (e *ErrorMissingPermission) Error() string {
	return "You don't have permission to do that! Requires " + e.Permission
}

func NewErrorMissingPermission(value string) *ErrorMissingPermission {
	return &ErrorMissingPermission{Permission: value}
}
