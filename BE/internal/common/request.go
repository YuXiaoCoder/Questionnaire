package common

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取查询参数: 字符串类型
func GetQueryParam(ctx *gin.Context, key string) (string, error) {
	value := ctx.Query(key)
	if value == "" {
		return value, errors.New(fmt.Sprintf("在查询参数中未找到[%s]", key))
	} else {
		return value, nil
	}
}

// 获取查询参数: 整数类型
func GetIntQueryParam(ctx *gin.Context, key string) (int, error) {
	value, err := GetQueryParam(ctx, key)
	if err != nil {
		return -1, err
	} else {
		return strconv.Atoi(value)
	}
}

// 获取路径参数
func GetPathParam(ctx *gin.Context, key string) (string, error) {
	value, flag := ctx.Params.Get(key)
	if !flag {
		return "", errors.New(fmt.Sprintf("在URI参数中未找到[%s]", key))
	}
	return value, nil
}

func GetIntPathParam(ctx *gin.Context, key string) (int, error) {
	value, err := GetPathParam(ctx, key)
	if err != nil {
		return -1, err
	} else {
		return strconv.Atoi(value)
	}
}
