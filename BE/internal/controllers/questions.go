package controllers

import (
	"BE/internal/common"
	"BE/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建问题
func CreateQuestionHandler(ctx *gin.Context) {
	var question models.Question

	// JSON -> Struct
	err := ctx.ShouldBindJSON(&question)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 创建问题
	err = models.CreateQuestion(&question)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusCreated, question)
		return
	}
}

// 获取所有问题
func GetQuestionsHandler(ctx *gin.Context) {
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

	// 类型
	questionType, err := common.GetIntQueryParam(ctx, "type")
	if err == nil {
		conditions["type"] = questionType
	}

	// 问卷ID
	questionnaireID, err := common.GetIntQueryParam(ctx, "questionnaire_id")
	if err == nil {
		conditions["questionnaire_id"] = questionnaireID
	}

	questions, total, err := models.GetQuestions(conditions)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, gin.H{
			"total":     total,
			"questions": questions,
		})
		return
	}
}

// 获取单个问题
func GetQuestionHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个问题
	question, err := models.GetQuestion(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, question)
		return
	}
}

// 更新问题
func UpdateQuestionHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个问题
	_, err = models.GetQuestion(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// JSON -> Struct
	var question *models.Question
	err = ctx.ShouldBindJSON(&question)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 更新问题
	err = models.UpdateQuestion(id, question)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// 返回问题
	question, err = models.GetQuestion(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, question)
		return
	}
}

// 删除问题
func DeleteQuestionHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个问题
	_, err = models.GetQuestion(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// 删除问题
	err = models.DeleteQuestion(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusNoContent, nil)
		return
	}
}
