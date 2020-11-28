package controllers

import (
	"BE/internal/common"
	"BE/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建答案
func CreateAnswerHandler(ctx *gin.Context) {
	var answer models.Answer

	// JSON -> Struct
	err := ctx.ShouldBindJSON(&answer)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 创建答案
	err = models.CreateAnswer(&answer)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusCreated, answer)
		return
	}
}

// 获取所有答案
func GetAnswersHandler(ctx *gin.Context) {
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

	// 问题ID
	questionID, err := common.GetIntQueryParam(ctx, "question_id")
	if err == nil {
		conditions["question_id"] = questionID
	}

	// 答题卡ID
	answerSheetID, err := common.GetIntQueryParam(ctx, "answer_sheet_id")
	if err == nil {
		conditions["answer_sheet_id"] = answerSheetID
	}

	answers, total, err := models.GetAnswers(conditions)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, gin.H{
			"total":   total,
			"answers": answers,
		})
		return
	}
}

// 获取单个答案
func GetAnswerHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个答案
	answer, err := models.GetAnswer(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, answer)
		return
	}
}

// 更新答案
func UpdateAnswerHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个答案
	_, err = models.GetAnswer(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// JSON -> Struct
	var answer *models.Answer
	err = ctx.ShouldBindJSON(&answer)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 更新答案
	err = models.UpdateAnswer(id, answer)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// 返回答案
	answer, err = models.GetAnswer(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, answer)
		return
	}
}

// 删除答案
func DeleteAnswerHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个答案
	_, err = models.GetAnswer(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// 删除答案
	err = models.DeleteAnswer(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusNoContent, nil)
		return
	}
}
