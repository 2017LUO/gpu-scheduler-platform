package errcode

const (
	CodeOK               = 0
	CodeInvalidArgument  = 1001
	CodeNotFound         = 1002
	CodeConflict         = 1003
	CodeQuotaExceeded    = 1004
	CodeUnauthenticated  = 1401
	CodePermissionDenied = 1403
	CodeRateLimited      = 1429
	CodeInternal         = 1500
)

func String(code int) string {
	switch code {
	case CodeOK:
		return "OK"
	case CodeInvalidArgument:
		return "INVALID_ARGUMENT"
	case CodeNotFound:
		return "NOT_FOUND"
	case CodeConflict:
		return "CONFLICT"
	case CodeQuotaExceeded:
		return "QUOTA_EXCEEDED"
	case CodeUnauthenticated:
		return "UNAUTHENTICATED"
	case CodePermissionDenied:
		return "PERMISSION_DENIED"
	case CodeRateLimited:
		return "RATE_LIMITED"
	default:
		return "INTERNAL_ERROR"
	}
}
