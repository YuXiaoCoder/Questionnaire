package common

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SuccessResponse(ctx *gin.Context, code int, data interface{}) {
	ctx.JSON(code, data)
	return
}

func FailureResponse(ctx *gin.Context, code int, message error) {
	// 类型断言
	errs, ok := message.(validator.ValidationErrors)
	if ok {
		// 数据验证错误
		ctx.JSON(code, gin.H{
			"message": RemoveTopStruct(errs.Translate(Trans)),
		})
	} else {
		// 其他错误
		ctx.JSON(code, gin.H{
			"message": message.Error(),
		})
	}
}
