package lookup

import "testing"

type Status int

const (
	StatusContinue                      Status = 100
	StatusSwitchingProtocols            Status = 101
	StatusProcessing                    Status = 102
	StatusEarlyHints                    Status = 103
	StatusOK                            Status = 200
	StatusCreated                       Status = 201
	StatusAccepted                      Status = 202
	StatusNonAuthoritativeInfo          Status = 203
	StatusNoContent                     Status = 204
	StatusResetContent                  Status = 205
	StatusPartialContent                Status = 206
	StatusMultiStatus                   Status = 207
	StatusAlreadyReported               Status = 208
	StatusIMUsed                        Status = 226
	StatusMultipleChoices               Status = 300
	StatusMovedPermanently              Status = 301
	StatusFound                         Status = 302
	StatusSeeOther                      Status = 303
	StatusNotModified                   Status = 304
	StatusUseProxy                      Status = 305
	StatusTemporaryRedirect             Status = 307
	StatusPermanentRedirect             Status = 308
	StatusBadRequest                    Status = 400
	StatusUnauthorized                  Status = 401
	StatusPaymentRequired               Status = 402
	StatusForbidden                     Status = 403
	StatusNotFound                      Status = 404
	StatusMethodNotAllowed              Status = 405
	StatusNotAcceptable                 Status = 406
	StatusProxyAuthRequired             Status = 407
	StatusRequestTimeout                Status = 408
	StatusConflict                      Status = 409
	StatusGone                          Status = 410
	StatusLengthRequired                Status = 411
	StatusPreconditionFailed            Status = 412
	StatusRequestEntityTooLarge         Status = 413
	StatusRequestURITooLong             Status = 414
	StatusUnsupportedMediaType          Status = 415
	StatusRequestedRangeNotSatisfiable  Status = 416
	StatusExpectationFailed             Status = 417
	StatusTeapot                        Status = 418
	StatusMisdirectedRequest            Status = 421
	StatusUnprocessableEntity           Status = 422
	StatusLocked                        Status = 423
	StatusFailedDependency              Status = 424
	StatusTooEarly                      Status = 425
	StatusUpgradeRequired               Status = 426
	StatusPreconditionRequired          Status = 428
	StatusTooManyRequests               Status = 429
	StatusRequestHeaderFieldsTooLarge   Status = 431
	StatusUnavailableForLegalReasons    Status = 451
	StatusInternalServerError           Status = 500
	StatusNotImplemented                Status = 501
	StatusBadGateway                    Status = 502
	StatusServiceUnavailable            Status = 503
	StatusGatewayTimeout                Status = 504
	StatusHTTPVersionNotSupported       Status = 505
	StatusVariantAlsoNegotiates         Status = 506
	StatusInsufficientStorage           Status = 507
	StatusLoopDetected                  Status = 508
	StatusNotExtended                   Status = 510
	StatusNetworkAuthenticationRequired Status = 511
	StatusCustomIDK                     Status = 600
	StatusCustomJustDont                Status = 601
)

var AllStatuses = [...]string{
	"StatusContinue",
	"StatusSwitchingProtocols",
	"StatusProcessing",
	"StatusEarlyHints",
	"StatusOK",
	"StatusCreated",
	"StatusAccepted",
	"StatusNonAuthoritativeInfo",
	"StatusNoContent",
	"StatusResetContent",
	"StatusPartialContent",
	"StatusMultiStatus",
	"StatusAlreadyReported",
	"StatusIMUsed",
	"StatusMultipleChoices",
	"StatusMovedPermanently",
	"StatusFound",
	"StatusSeeOther",
	"StatusNotModified",
	"StatusUseProxy",
	"StatusTemporaryRedirect",
	"StatusPermanentRedirect",
	"StatusBadRequest",
	"StatusUnauthorized",
	"StatusPaymentRequired",
	"StatusForbidden",
	"StatusNotFound",
	"StatusMethodNotAllowed",
	"StatusNotAcceptable",
	"StatusProxyAuthRequired",
	"StatusRequestTimeout",
	"StatusConflict",
	"StatusGone",
	"StatusLengthRequired",
	"StatusPreconditionFailed",
	"StatusRequestEntityTooLarge",
	"StatusRequestURITooLong",
	"StatusUnsupportedMediaType",
	"StatusRequestedRangeNotSatisfiable",
	"StatusExpectationFailed",
	"StatusTeapot",
	"StatusMisdirectedRequest",
	"StatusUnprocessableEntity",
	"StatusLocked",
	"StatusFailedDependency",
	"StatusTooEarly",
	"StatusUpgradeRequired",
	"StatusPreconditionRequired",
	"StatusTooManyRequests",
	"StatusRequestHeaderFieldsTooLarge",
	"StatusUnavailableForLegalReasons",
	"StatusInternalServerError",
	"StatusNotImplemented",
	"StatusBadGateway",
	"StatusServiceUnavailable",
	"StatusGatewayTimeout",
	"StatusHTTPVersionNotSupported",
	"StatusVariantAlsoNegotiates",
	"StatusInsufficientStorage",
	"StatusLoopDetected",
	"StatusNotExtended",
	"StatusNetworkAuthenticationRequired",
	"StatusCustomIDK",
	"StatusCustomJustDont",
}

func Switch8(s string) Status {
	switch s {
	case "StatusContinue":
		return StatusContinue
	case "StatusSwitchingProtocols":
		return StatusSwitchingProtocols
	case "StatusProcessing":
		return StatusProcessing
	case "StatusEarlyHints":
		return StatusEarlyHints
	case "StatusOK":
		return StatusOK
	case "StatusCreated":
		return StatusCreated
	case "StatusAccepted":
		return StatusAccepted
	case "StatusNonAuthoritativeInfo":
		return StatusNonAuthoritativeInfo
	default:
		return 0
	}
}

var Map8 = map[string]Status{
	"StatusContinue":             StatusContinue,
	"StatusSwitchingProtocols":   StatusSwitchingProtocols,
	"StatusProcessing":           StatusProcessing,
	"StatusEarlyHints":           StatusEarlyHints,
	"StatusOK":                   StatusOK,
	"StatusCreated":              StatusCreated,
	"StatusAccepted":             StatusAccepted,
	"StatusNonAuthoritativeInfo": StatusNonAuthoritativeInfo,
}

func Switch16(s string) Status {
	switch s {
	case "StatusContinue":
		return StatusContinue
	case "StatusSwitchingProtocols":
		return StatusSwitchingProtocols
	case "StatusProcessing":
		return StatusProcessing
	case "StatusEarlyHints":
		return StatusEarlyHints
	case "StatusOK":
		return StatusOK
	case "StatusCreated":
		return StatusCreated
	case "StatusAccepted":
		return StatusAccepted
	case "StatusNonAuthoritativeInfo":
		return StatusNonAuthoritativeInfo
	case "StatusNoContent":
		return StatusNoContent
	case "StatusResetContent":
		return StatusResetContent
	case "StatusPartialContent":
		return StatusPartialContent
	case "StatusMultiStatus":
		return StatusMultiStatus
	case "StatusAlreadyReported":
		return StatusAlreadyReported
	case "StatusIMUsed":
		return StatusIMUsed
	case "StatusMultipleChoices":
		return StatusMultipleChoices
	case "StatusMovedPermanently":
		return StatusMovedPermanently
	default:
		return 0
	}
}

var Map16 = map[string]Status{
	"StatusContinue":             StatusContinue,
	"StatusSwitchingProtocols":   StatusSwitchingProtocols,
	"StatusProcessing":           StatusProcessing,
	"StatusEarlyHints":           StatusEarlyHints,
	"StatusOK":                   StatusOK,
	"StatusCreated":              StatusCreated,
	"StatusAccepted":             StatusAccepted,
	"StatusNonAuthoritativeInfo": StatusNonAuthoritativeInfo,
	"StatusNoContent":            StatusNoContent,
	"StatusResetContent":         StatusResetContent,
	"StatusPartialContent":       StatusPartialContent,
	"StatusMultiStatus":          StatusMultiStatus,
	"StatusAlreadyReported":      StatusAlreadyReported,
	"StatusIMUsed":               StatusIMUsed,
	"StatusMultipleChoices":      StatusMultipleChoices,
	"StatusMovedPermanently":     StatusMovedPermanently,
}

func Switch32(s string) Status {
	switch s {
	case "StatusContinue":
		return StatusContinue
	case "StatusSwitchingProtocols":
		return StatusSwitchingProtocols
	case "StatusProcessing":
		return StatusProcessing
	case "StatusEarlyHints":
		return StatusEarlyHints
	case "StatusOK":
		return StatusOK
	case "StatusCreated":
		return StatusCreated
	case "StatusAccepted":
		return StatusAccepted
	case "StatusNonAuthoritativeInfo":
		return StatusNonAuthoritativeInfo
	case "StatusNoContent":
		return StatusNoContent
	case "StatusResetContent":
		return StatusResetContent
	case "StatusPartialContent":
		return StatusPartialContent
	case "StatusMultiStatus":
		return StatusMultiStatus
	case "StatusAlreadyReported":
		return StatusAlreadyReported
	case "StatusIMUsed":
		return StatusIMUsed
	case "StatusMultipleChoices":
		return StatusMultipleChoices
	case "StatusMovedPermanently":
		return StatusMovedPermanently
	case "StatusFound":
		return StatusFound
	case "StatusSeeOther":
		return StatusSeeOther
	case "StatusNotModified":
		return StatusNotModified
	case "StatusUseProxy":
		return StatusUseProxy
	case "StatusTemporaryRedirect":
		return StatusTemporaryRedirect
	case "StatusPermanentRedirect":
		return StatusPermanentRedirect
	case "StatusBadRequest":
		return StatusBadRequest
	case "StatusUnauthorized":
		return StatusUnauthorized
	case "StatusPaymentRequired":
		return StatusPaymentRequired
	case "StatusForbidden":
		return StatusForbidden
	case "StatusNotFound":
		return StatusNotFound
	case "StatusMethodNotAllowed":
		return StatusMethodNotAllowed
	case "StatusNotAcceptable":
		return StatusNotAcceptable
	case "StatusProxyAuthRequired":
		return StatusProxyAuthRequired
	case "StatusRequestTimeout":
		return StatusRequestTimeout
	case "StatusConflict":
		return StatusConflict
	default:
		return 0
	}
}

var Map32 = map[string]Status{
	"StatusContinue":             StatusContinue,
	"StatusSwitchingProtocols":   StatusSwitchingProtocols,
	"StatusProcessing":           StatusProcessing,
	"StatusEarlyHints":           StatusEarlyHints,
	"StatusOK":                   StatusOK,
	"StatusCreated":              StatusCreated,
	"StatusAccepted":             StatusAccepted,
	"StatusNonAuthoritativeInfo": StatusNonAuthoritativeInfo,
	"StatusNoContent":            StatusNoContent,
	"StatusResetContent":         StatusResetContent,
	"StatusPartialContent":       StatusPartialContent,
	"StatusMultiStatus":          StatusMultiStatus,
	"StatusAlreadyReported":      StatusAlreadyReported,
	"StatusIMUsed":               StatusIMUsed,
	"StatusMultipleChoices":      StatusMultipleChoices,
	"StatusMovedPermanently":     StatusMovedPermanently,
	"StatusFound":                StatusFound,
	"StatusSeeOther":             StatusSeeOther,
	"StatusNotModified":          StatusNotModified,
	"StatusUseProxy":             StatusUseProxy,
	"StatusTemporaryRedirect":    StatusTemporaryRedirect,
	"StatusPermanentRedirect":    StatusPermanentRedirect,
	"StatusBadRequest":           StatusBadRequest,
	"StatusUnauthorized":         StatusUnauthorized,
	"StatusPaymentRequired":      StatusPaymentRequired,
	"StatusForbidden":            StatusForbidden,
	"StatusNotFound":             StatusNotFound,
	"StatusMethodNotAllowed":     StatusMethodNotAllowed,
	"StatusNotAcceptable":        StatusNotAcceptable,
	"StatusProxyAuthRequired":    StatusProxyAuthRequired,
	"StatusRequestTimeout":       StatusRequestTimeout,
	"StatusConflict":             StatusConflict,
}

func Switch64(s string) Status {
	switch s {
	case "StatusContinue":
		return StatusContinue
	case "StatusSwitchingProtocols":
		return StatusSwitchingProtocols
	case "StatusProcessing":
		return StatusProcessing
	case "StatusEarlyHints":
		return StatusEarlyHints
	case "StatusOK":
		return StatusOK
	case "StatusCreated":
		return StatusCreated
	case "StatusAccepted":
		return StatusAccepted
	case "StatusNonAuthoritativeInfo":
		return StatusNonAuthoritativeInfo
	case "StatusNoContent":
		return StatusNoContent
	case "StatusResetContent":
		return StatusResetContent
	case "StatusPartialContent":
		return StatusPartialContent
	case "StatusMultiStatus":
		return StatusMultiStatus
	case "StatusAlreadyReported":
		return StatusAlreadyReported
	case "StatusIMUsed":
		return StatusIMUsed
	case "StatusMultipleChoices":
		return StatusMultipleChoices
	case "StatusMovedPermanently":
		return StatusMovedPermanently
	case "StatusFound":
		return StatusFound
	case "StatusSeeOther":
		return StatusSeeOther
	case "StatusNotModified":
		return StatusNotModified
	case "StatusUseProxy":
		return StatusUseProxy
	case "StatusTemporaryRedirect":
		return StatusTemporaryRedirect
	case "StatusPermanentRedirect":
		return StatusPermanentRedirect
	case "StatusBadRequest":
		return StatusBadRequest
	case "StatusUnauthorized":
		return StatusUnauthorized
	case "StatusPaymentRequired":
		return StatusPaymentRequired
	case "StatusForbidden":
		return StatusForbidden
	case "StatusNotFound":
		return StatusNotFound
	case "StatusMethodNotAllowed":
		return StatusMethodNotAllowed
	case "StatusNotAcceptable":
		return StatusNotAcceptable
	case "StatusProxyAuthRequired":
		return StatusProxyAuthRequired
	case "StatusRequestTimeout":
		return StatusRequestTimeout
	case "StatusConflict":
		return StatusConflict
	case "StatusGone":
		return StatusGone
	case "StatusLengthRequired":
		return StatusLengthRequired
	case "StatusPreconditionFailed":
		return StatusPreconditionFailed
	case "StatusRequestEntityTooLarge":
		return StatusRequestEntityTooLarge
	case "StatusRequestURITooLong":
		return StatusRequestURITooLong
	case "StatusUnsupportedMediaType":
		return StatusUnsupportedMediaType
	case "StatusRequestedRangeNotSatisfiable":
		return StatusRequestedRangeNotSatisfiable
	case "StatusExpectationFailed":
		return StatusExpectationFailed
	case "StatusTeapot":
		return StatusTeapot
	case "StatusMisdirectedRequest":
		return StatusMisdirectedRequest
	case "StatusUnprocessableEntity":
		return StatusUnprocessableEntity
	case "StatusLocked":
		return StatusLocked
	case "StatusFailedDependency":
		return StatusFailedDependency
	case "StatusTooEarly":
		return StatusTooEarly
	case "StatusUpgradeRequired":
		return StatusUpgradeRequired
	case "StatusPreconditionRequired":
		return StatusPreconditionRequired
	case "StatusTooManyRequests":
		return StatusTooManyRequests
	case "StatusRequestHeaderFieldsTooLarge":
		return StatusRequestHeaderFieldsTooLarge
	case "StatusUnavailableForLegalReasons":
		return StatusUnavailableForLegalReasons
	case "StatusInternalServerError":
		return StatusInternalServerError
	case "StatusNotImplemented":
		return StatusNotImplemented
	case "StatusBadGateway":
		return StatusBadGateway
	case "StatusServiceUnavailable":
		return StatusServiceUnavailable
	case "StatusGatewayTimeout":
		return StatusGatewayTimeout
	case "StatusHTTPVersionNotSupported":
		return StatusHTTPVersionNotSupported
	case "StatusVariantAlsoNegotiates":
		return StatusVariantAlsoNegotiates
	case "StatusInsufficientStorage":
		return StatusInsufficientStorage
	case "StatusLoopDetected":
		return StatusLoopDetected
	case "StatusNotExtended":
		return StatusNotExtended
	case "StatusNetworkAuthenticationRequired":
		return StatusNetworkAuthenticationRequired
	case "StatusCustomIDK":
		return StatusCustomIDK
	case "StatusCustomJustDont":
		return StatusCustomJustDont
	default:
		return 0
	}
}

var Map64 = map[string]Status{
	"StatusContinue":                      StatusContinue,
	"StatusSwitchingProtocols":            StatusSwitchingProtocols,
	"StatusProcessing":                    StatusProcessing,
	"StatusEarlyHints":                    StatusEarlyHints,
	"StatusOK":                            StatusOK,
	"StatusCreated":                       StatusCreated,
	"StatusAccepted":                      StatusAccepted,
	"StatusNonAuthoritativeInfo":          StatusNonAuthoritativeInfo,
	"StatusNoContent":                     StatusNoContent,
	"StatusResetContent":                  StatusResetContent,
	"StatusPartialContent":                StatusPartialContent,
	"StatusMultiStatus":                   StatusMultiStatus,
	"StatusAlreadyReported":               StatusAlreadyReported,
	"StatusIMUsed":                        StatusIMUsed,
	"StatusMultipleChoices":               StatusMultipleChoices,
	"StatusMovedPermanently":              StatusMovedPermanently,
	"StatusFound":                         StatusFound,
	"StatusSeeOther":                      StatusSeeOther,
	"StatusNotModified":                   StatusNotModified,
	"StatusUseProxy":                      StatusUseProxy,
	"StatusTemporaryRedirect":             StatusTemporaryRedirect,
	"StatusPermanentRedirect":             StatusPermanentRedirect,
	"StatusBadRequest":                    StatusBadRequest,
	"StatusUnauthorized":                  StatusUnauthorized,
	"StatusPaymentRequired":               StatusPaymentRequired,
	"StatusForbidden":                     StatusForbidden,
	"StatusNotFound":                      StatusNotFound,
	"StatusMethodNotAllowed":              StatusMethodNotAllowed,
	"StatusNotAcceptable":                 StatusNotAcceptable,
	"StatusProxyAuthRequired":             StatusProxyAuthRequired,
	"StatusRequestTimeout":                StatusRequestTimeout,
	"StatusConflict":                      StatusConflict,
	"StatusGone":                          StatusGone,
	"StatusLengthRequired":                StatusLengthRequired,
	"StatusPreconditionFailed":            StatusPreconditionFailed,
	"StatusRequestEntityTooLarge":         StatusRequestEntityTooLarge,
	"StatusRequestURITooLong":             StatusRequestURITooLong,
	"StatusUnsupportedMediaType":          StatusUnsupportedMediaType,
	"StatusRequestedRangeNotSatisfiable":  StatusRequestedRangeNotSatisfiable,
	"StatusExpectationFailed":             StatusExpectationFailed,
	"StatusTeapot":                        StatusTeapot,
	"StatusMisdirectedRequest":            StatusMisdirectedRequest,
	"StatusUnprocessableEntity":           StatusUnprocessableEntity,
	"StatusLocked":                        StatusLocked,
	"StatusFailedDependency":              StatusFailedDependency,
	"StatusTooEarly":                      StatusTooEarly,
	"StatusUpgradeRequired":               StatusUpgradeRequired,
	"StatusPreconditionRequired":          StatusPreconditionRequired,
	"StatusTooManyRequests":               StatusTooManyRequests,
	"StatusRequestHeaderFieldsTooLarge":   StatusRequestHeaderFieldsTooLarge,
	"StatusUnavailableForLegalReasons":    StatusUnavailableForLegalReasons,
	"StatusInternalServerError":           StatusInternalServerError,
	"StatusNotImplemented":                StatusNotImplemented,
	"StatusBadGateway":                    StatusBadGateway,
	"StatusServiceUnavailable":            StatusServiceUnavailable,
	"StatusGatewayTimeout":                StatusGatewayTimeout,
	"StatusHTTPVersionNotSupported":       StatusHTTPVersionNotSupported,
	"StatusVariantAlsoNegotiates":         StatusVariantAlsoNegotiates,
	"StatusInsufficientStorage":           StatusInsufficientStorage,
	"StatusLoopDetected":                  StatusLoopDetected,
	"StatusNotExtended":                   StatusNotExtended,
	"StatusNetworkAuthenticationRequired": StatusNetworkAuthenticationRequired,
	"StatusCustomIDK":                     StatusCustomIDK,
	"StatusCustomJustDont":                StatusCustomJustDont,
}

func BenchmarkLookup_Switch8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Switch8(AllStatuses[i%len(AllStatuses)])
	}
}

func BenchmarkLookup_Map8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := Map8[AllStatuses[i%len(AllStatuses)]]
		_ = ok
	}
}

func BenchmarkLookup_Switch16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Switch16(AllStatuses[i%len(AllStatuses)])
	}
}

func BenchmarkLookup_Map16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := Map16[AllStatuses[i%len(AllStatuses)]]
		_ = ok
	}
}

func BenchmarkLookup_Switch32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Switch32(AllStatuses[i%len(AllStatuses)])
	}
}

func BenchmarkLookup_Map32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := Map32[AllStatuses[i%len(AllStatuses)]]
		_ = ok
	}
}

func BenchmarkLookup_Switch64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Switch64(AllStatuses[i%len(AllStatuses)])
	}
}

func BenchmarkLookup_Map64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := Map64[AllStatuses[i%len(AllStatuses)]]
		_ = ok
	}
}
