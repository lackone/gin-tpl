package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type ResultOption func(*Result)

func WithCode(code int) ResultOption {
	return func(result *Result) {
		result.Code = code
	}
}

func WithMsg(msg string) ResultOption {
	return func(result *Result) {
		result.Msg = msg
	}
}

func WithData(data any) ResultOption {
	return func(result *Result) {
		result.Data = data
	}
}

func Success(ctx *gin.Context, opts ...ResultOption) {
	r := Result{
		Code: 200,
		Msg:  "success",
		Data: map[string]any{},
	}

	for _, opt := range opts {
		opt(&r)
	}

	ctx.JSON(http.StatusOK, r)
}

func Error(ctx *gin.Context, opts ...ResultOption) {
	r := Result{
		Code: 500,
		Msg:  "error",
		Data: map[string]any{},
	}

	for _, opt := range opts {
		opt(&r)
	}

	ctx.JSON(http.StatusInternalServerError, r)
}
