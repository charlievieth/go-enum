package main

import (
	"encoding/json"
	"fmt"
)

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
	Status700                           Status = 700
	Status701                           Status = 701
)

func main() {
	ck(StatusContinue, "StatusContinue", false)
	ck(StatusSwitchingProtocols, "StatusSwitchingProtocols", false)
	ck(StatusProcessing, "StatusProcessing", false)
	ck(StatusEarlyHints, "StatusEarlyHints", false)
	ck(StatusOK, "StatusOK", false)
	ck(StatusCreated, "StatusCreated", false)
	ck(StatusAccepted, "StatusAccepted", false)
	ck(StatusNonAuthoritativeInfo, "StatusNonAuthoritativeInfo", false)
	ck(StatusNoContent, "StatusNoContent", false)
	ck(StatusResetContent, "StatusResetContent", false)
	ck(StatusPartialContent, "StatusPartialContent", false)
	ck(StatusMultiStatus, "StatusMultiStatus", false)
	ck(StatusAlreadyReported, "StatusAlreadyReported", false)
	ck(StatusIMUsed, "StatusIMUsed", false)
	ck(StatusMultipleChoices, "StatusMultipleChoices", false)
	ck(StatusMovedPermanently, "StatusMovedPermanently", false)
	ck(StatusFound, "StatusFound", false)
	ck(StatusSeeOther, "StatusSeeOther", false)
	ck(StatusNotModified, "StatusNotModified", false)
	ck(StatusUseProxy, "StatusUseProxy", false)
	ck(StatusTemporaryRedirect, "StatusTemporaryRedirect", false)
	ck(StatusPermanentRedirect, "StatusPermanentRedirect", false)
	ck(StatusBadRequest, "StatusBadRequest", false)
	ck(StatusUnauthorized, "StatusUnauthorized", false)
	ck(StatusPaymentRequired, "StatusPaymentRequired", false)
	ck(StatusForbidden, "StatusForbidden", false)
	ck(StatusNotFound, "StatusNotFound", false)
	ck(StatusMethodNotAllowed, "StatusMethodNotAllowed", false)
	ck(StatusNotAcceptable, "StatusNotAcceptable", false)
	ck(StatusProxyAuthRequired, "StatusProxyAuthRequired", false)
	ck(StatusRequestTimeout, "StatusRequestTimeout", false)
	ck(StatusConflict, "StatusConflict", false)
	ck(StatusGone, "StatusGone", false)
	ck(StatusLengthRequired, "StatusLengthRequired", false)
	ck(StatusPreconditionFailed, "StatusPreconditionFailed", false)
	ck(StatusRequestEntityTooLarge, "StatusRequestEntityTooLarge", false)
	ck(StatusRequestURITooLong, "StatusRequestURITooLong", false)
	ck(StatusUnsupportedMediaType, "StatusUnsupportedMediaType", false)
	ck(StatusRequestedRangeNotSatisfiable, "StatusRequestedRangeNotSatisfiable", false)
	ck(StatusExpectationFailed, "StatusExpectationFailed", false)
	ck(StatusTeapot, "StatusTeapot", false)
	ck(StatusMisdirectedRequest, "StatusMisdirectedRequest", false)
	ck(StatusUnprocessableEntity, "StatusUnprocessableEntity", false)
	ck(StatusLocked, "StatusLocked", false)
	ck(StatusFailedDependency, "StatusFailedDependency", false)
	ck(StatusTooEarly, "StatusTooEarly", false)
	ck(StatusUpgradeRequired, "StatusUpgradeRequired", false)
	ck(StatusPreconditionRequired, "StatusPreconditionRequired", false)
	ck(StatusTooManyRequests, "StatusTooManyRequests", false)
	ck(StatusRequestHeaderFieldsTooLarge, "StatusRequestHeaderFieldsTooLarge", false)
	ck(StatusUnavailableForLegalReasons, "StatusUnavailableForLegalReasons", false)
	ck(StatusInternalServerError, "StatusInternalServerError", false)
	ck(StatusNotImplemented, "StatusNotImplemented", false)
	ck(StatusBadGateway, "StatusBadGateway", false)
	ck(StatusServiceUnavailable, "StatusServiceUnavailable", false)
	ck(StatusGatewayTimeout, "StatusGatewayTimeout", false)
	ck(StatusHTTPVersionNotSupported, "StatusHTTPVersionNotSupported", false)
	ck(StatusVariantAlsoNegotiates, "StatusVariantAlsoNegotiates", false)
	ck(StatusInsufficientStorage, "StatusInsufficientStorage", false)
	ck(StatusLoopDetected, "StatusLoopDetected", false)
	ck(StatusNotExtended, "StatusNotExtended", false)
	ck(StatusNetworkAuthenticationRequired, "StatusNetworkAuthenticationRequired", false)
	ck(StatusCustomIDK, "StatusCustomIDK", false)
	ck(StatusCustomJustDont, "StatusCustomJustDont", false)
	ck(Status700, "Status700", false)
	ck(Status701, "Status701", false)

	ck(-1, "Status(-1)", true)
	ck(0, "Status(0)", true)
	ck(1, "Status(1)", true)
	ck(209, "Status(209)", true)
}

func ck(c Status, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("http_status.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("http_status.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("http_status.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("http_status.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Status
		if err := json.Unmarshal(b, &v); err != nil {
			panic("http_status.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("http_status.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("http_status.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("http_status.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("http_status.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Status
		if err := v.UnmarshalText(b); err != nil {
			panic("http_status.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("http_status.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Status
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("http_status.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("http_status.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Status
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("http_status.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("http_status.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
