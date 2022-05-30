package response

// Internal server Response Codes
const (
	CodeOk = iota
	CodeUnknownException
	CodeForbidden
	CodeInvalidParams
	CodeNotFound
	CodeBadRequest
	CodeInvalidJsonConversion
	CodeUserPasswordIsEmpty
	CodeUnknownUser
	CodeCryptoError
	CodeDBError
)
