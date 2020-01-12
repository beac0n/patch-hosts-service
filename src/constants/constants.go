package constants

import "net/http"

var HttpMethods = [...]string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

const HeaderPrefix = "X-Phs-"
const ContentLength = "Content-Length"
const WrongHttpMethod = "wrong http method"
const NoContent = "no content"

const Req = "/req"
const Res = "/res"
