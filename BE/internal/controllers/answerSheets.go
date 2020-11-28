package controllers

import (
	"BE/internal/common"
	"BE/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建答题卡
func CreateAnswerSheetHandler(ctx *gin.Context) {
	var answerSheet models.AnswerSheet

	// JSON -> Struct
	err := ctx.ShouldBindJSON(&answerSheet)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 创建答题卡
	err = models.CreateAnswerSheet(&answerSheet)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusCreated, answerSheet)
		return
	}
}

// 获取所有答题卡
func GetAnswerSheetsHandler(ctx *gin.Context) {
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

	// 用户ID
	userID, err := common.GetIntQueryParam(ctx, "user_id")
	if err == nil {
		conditions["user_id"] = userID
	}

	// 问卷ID
	questionnaireID, err := common.GetIntQueryParam(ctx, "questionnaire_id")
	if err == nil {
		conditions["questionnaire_id"] = questionnaireID
	}

	answerSheets, total, err := models.GetAnswerSheets(conditions)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, gin.H{
			"total":         total,
			"answer_sheets": answerSheets,
		})
		return
	}
}

// 获取单个答题卡
func GetAnswerSheetHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个答题卡
	answerSheet, err := models.GetAnswerSheet(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, answerSheet)
		return
	}
}

// 更新答题卡
func UpdateAnswerSheetHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个答题卡
	_, err = models.GetQuestionnaire(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// JSON -> Struct
	var answerSheet *models.AnswerSheet
	err = ctx.ShouldBindJSON(&answerSheet)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 更新答题卡
	err = models.UpdateAnswerSheet(id, answerSheet)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// 返回答题卡
	answerSheet, err = models.GetAnswerSheet(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, answerSheet)
		return
	}
}

// 删除答题卡
func DeleteAnswerSheetHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个答题卡
	_, err = models.GetAnswerSheet(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// 删除答题卡
	err = models.DeleteAnswerSheet(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusNoContent, nil)
		return
	}
}
