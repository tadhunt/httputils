package httputils

import (
	"fmt"
	"log"
	"net/http"

	"esplitter.com/debug"

	"github.com/mholt/binding"
)

type HttpError struct {
	Skip int // callstack frames to skip for debug tracking

	Errs binding.Errors

	Msg  string
	Code int

	LogMsg string
}

// HttpError represents both the error to send to the user and the error to log.
func NewHttpError() *HttpError {
	herr := &HttpError{
		Skip:   1,
		Errs:   nil,
		Msg:    "OK",
		Code:   http.StatusOK,
		LogMsg: "OK",
	}

	return herr
}

func (herr *HttpError) mkerror(skip int, code int, msg string, logErr error) {
	herr.Code = code
	if msg == "" {
		herr.Msg = fmt.Sprintf("%d %v", code, http.StatusText(code))
	} else {
		herr.Msg = fmt.Sprintf("%d %v", code, msg)
	}

	herr.LogMsg = fmt.Sprintf("{%s} <%d> %s [%v]", debug.FuncName(skip, true), herr.Code, herr.Msg, logErr)
}

func (herr *HttpError) Error(code int, msg string, logErr error) {
	herr.mkerror(herr.Skip+1, code, msg, logErr)
}

func (herr *HttpError) Errors(errs binding.Errors) {
	herr.Errs = errs
}

func (herr *HttpError) InternalServerError(msg string, logErr error) {
	herr.mkerror(herr.Skip+1, http.StatusInternalServerError, msg, logErr)
}

func (herr *HttpError) OK() bool {
	if herr.Errs != nil {
		return false
	}

	if herr.Code != http.StatusOK {
		return false
	}

	return true
}

func (herr *HttpError) Write(w http.ResponseWriter) {
	if herr.Errs != nil {
		log.Printf("Binding Error: %s", herr.Errs.Error())
		herr.Errs.Handle(w)
		return
	}

	log.Printf("%s\n", herr.LogMsg)
	http.Error(w, herr.Msg, herr.Code)
}

// The following functions are shorthand for creating and emitting the error all at once
func InternalServerError(w http.ResponseWriter, msg string, logErr error) {
	herr := NewHttpError()
	herr.Skip++
	herr.InternalServerError(msg, logErr)
	herr.Write(w)
}

func Error(w http.ResponseWriter, code int, msg string, logErr error) {
	herr := NewHttpError()
	herr.Skip++
	herr.Error(code, msg, logErr)
	herr.Write(w)
}
