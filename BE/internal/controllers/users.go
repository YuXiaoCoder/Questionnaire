package controllers

import (
	"BE/internal/common"
	"BE/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v2"
	_ "github.com/medivhzhan/weapp/v2"
)

func LoginHandler(ctx *gin.Context) {
	// 获取 Code
	code, err := common.GetQueryParam(ctx, "code")
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	//获取配置
	wx := common.GetWeiXin()

	// https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/login.html
	// https://github.com/medivhzhan/miniapp#getAccessToken
	result, err := weapp.Login(wx.APPID, wx.Secret, code)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// 查询用户
	user, _ := models.GetUserByOpenID(result.OpenID)
	if user.ID == 0 {
		// 创建用户
		user.OpenID = result.OpenID
		err := models.CreateUser(user)
		if err != nil {
			common.FailureResponse(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// 获取会话
	session := sessions.Default(ctx)
	sessionValue := fmt.Sprintf("%d.%s", user.ID, result.SessionKey)
	session.Set("UserID", sessionValue)
	err = session.Save()
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResponse(ctx, http.StatusOK, gin.H{
		"user_id": user.ID,
	})
	return
}

// 创建用户
func CreateUserHandler(ctx *gin.Context) {
	var user models.User

	// JSON -> Struct
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 创建问卷
	err = models.CreateUser(&user)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusCreated, user)
		return
	}
}

// 获取所有用户
func GetUsersHandler(ctx *gin.Context) {
	// 查询条件
	conditions := make(map[string]interface{})

	// 限制
	limit, err := common.GetIntQueryParam(ctx, "limit")
	if err == nil {
		conditions["limit"] = limit
	}

	// 偏移量
	offset, err := common.GetIntQueryParam(ctx, "offset")
	if err == nil {
		conditions["offset"] = offset
	}

	users, total, err := models.GetUsers(conditions)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, gin.H{
			"total": total,
			"users": users,
		})
		return
	}
}

// 获取单个用户
func GetUserHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个用户
	user, err := models.GetUser(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, user)
		return
	}
}

// 更新用户
func UpdateUserHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个用户
	_, err = models.GetUser(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// JSON -> Struct
	var user *models.User
	err = ctx.ShouldBindJSON(&user)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 更新用户
	err = models.UpdateUser(id, user)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// 返回用户
	user, err = models.GetUser(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, user)
		return
	}
}

// 删除用户
func DeleteUserHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个用户
	_, err = models.GetUser(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// 删除用户
	err = models.DeleteUser(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusNoContent, nil)
		return
	}
}
