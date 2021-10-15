package i18n

import "strings"

type Response struct {
	Code    int    `json:"code"`
	Err     string `json:"error,omitempty"`
	Msg     string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}

// String unwraps a input Response object to its constituent Err or Msg field
func (r Response) String() string {
	if r.Err != "" {
		return r.Err
	} else {
		return r.Msg
	}
}

// LowerString converts the constituent message of a Response to lower-case
func (r Response) LowerString() string {
	return strings.ToLower(r.String())
}

// Error returns the Err field inside a Response object
func (r Response) Error() string {
	return r.Err
}

// Success checks if the Response object describes a successful object
func (r Response) Success() bool {
	return r.Err == "" && r.Code < 400
}

// NewResponse creates a new Response object with the given code, Message id and parameters
func NewResponse(code int, id Message, params ...interface{}) Response {
	if code < 400 {
		// Successful Response
		return Response{
			Code: code,
			Msg:  Msg(id, params...),
		}
	} else {
		// Error Response
		return Response{
			Code: code,
			Err:  Msg(id, params...),
		}
	}
}
