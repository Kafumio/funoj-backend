package response

import (
	e "funoj-backend/consts/error"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Result @Description: 统一Result，并返回数据给前端
type Result struct {
	ctx *gin.Context
}

type ResultContent struct {
	Code    int         `json:"code"`    //response code
	Message string      `json:"message"` //response message
	Data    interface{} `json:"data"`    //response data
}

func NewResult(ctx *gin.Context) *Result {
	return &Result{ctx: ctx}
}

func (r *Result) SuccessData(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := &ResultContent{
		Code:    200,
		Message: "request success",
		Data:    data,
	}
	r.ctx.JSON(http.StatusOK, res)
}

func (r *Result) SuccessMessage(message string) {
	res := &ResultContent{
		Code:    200,
		Message: message,
		Data:    nil,
	}
	r.ctx.JSON(http.StatusOK, res)
}

func (r *Result) Success(message string, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := &ResultContent{
		Code:    200,
		Message: message,
		Data:    data,
	}
	r.ctx.JSON(http.StatusOK, res)
}

// 返回异常信息
func (r *Result) Error(e *e.Error) {
	res := &ResultContent{
		Code:    e.Code,
		Message: e.Message,
	}
	r.ctx.JSON(e.HttpCode, res)
}

func (r *Result) SimpleError(code int, message string, data interface{}) {
	res := &ResultContent{
		Code:    code,
		Message: message,
		Data:    data,
	}
	r.ctx.JSON(http.StatusOK, res)
}

func (r *Result) SimpleErrorMessage(message string) {
	res := &ResultContent{
		Code:    400,
		Message: message,
	}
	r.ctx.JSON(http.StatusOK, res)
}
