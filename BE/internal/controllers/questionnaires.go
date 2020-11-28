package controllers

import (
	"BE/internal/common"
	"BE/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QuestionStatistic struct {
	ID     uint               `json:"id"`     // 问卷ID
	Result map[string]float64 `json:"result"` // 统计结果
}

// 创建问卷
func CreateQuestionnaireHandler(ctx *gin.Context) {
	var questionnaire models.Questionnaire

	// JSON -> Struct
	err := ctx.ShouldBindJSON(&questionnaire)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 创建问卷
	err = models.CreateQuestionnaire(&questionnaire)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusCreated, questionnaire)
		return
	}
}

// 获取所有问卷
func GetQuestionnairesHandler(ctx *gin.Context) {
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

	// 用户
	user, err := common.GetIntQueryParam(ctx, "user_id")
	if err == nil {
		conditions["user_id"] = user
	}

	// 类型
	questionnaireType, err := common.GetIntQueryParam(ctx, "type")
	if err == nil {
		conditions["type"] = questionnaireType
	}

	questionnaires, total, err := models.GetQuestionnaires(conditions)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, gin.H{
			"total":          total,
			"questionnaires": questionnaires,
		})
		return
	}
}

// 获取单个问卷
func GetQuestionnaireHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个问卷
	questionnaire, err := models.GetQuestionnaire(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, questionnaire)
		return
	}
}

// 更新问卷
func UpdateQuestionnaireHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个问卷
	_, err = models.GetQuestionnaire(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// JSON -> Struct
	var questionnaire *models.Questionnaire
	err = ctx.ShouldBindJSON(&questionnaire)
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 更新问卷
	err = models.UpdateQuestionnaire(id, questionnaire)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// 返回问卷
	questionnaire, err = models.GetQuestionnaire(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusOK, questionnaire)
		return
	}
}

// 删除问卷
func DeleteQuestionnaireHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询单个问卷
	_, err = models.GetQuestionnaire(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusNotFound, err)
		return
	}

	// 删除问卷
	err = models.DeleteQuestionnaire(id)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResponse(ctx, http.StatusNoContent, nil)
		return
	}
}

// 统计分析问卷
func AnalyzeQuestionnaireHandler(ctx *gin.Context) {
	// 获取ID
	id, err := common.GetIntPathParam(ctx, "id")
	if err != nil {
		common.FailureResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// 查询条件
	conditions := make(map[string]interface{})
	conditions["questionnaire_id"] = id

	questions, _, err := models.GetQuestions(conditions)
	if err != nil {
		common.FailureResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	var statistics []QuestionStatistic
	for _, question := range questions {
		var statistic QuestionStatistic
		statistic.ID = question.ID

		result, total, err := models.StatisticalAnswer(question.ID)
		if err != nil {
			common.FailureResponse(ctx, http.StatusInternalServerError, err)
			return
		}

		statistic.Result = make(map[string]float64)
		if question.Type == models.SingleChoice || question.Type == models.MultipleChoice {
			for k, v := range result {
				statistic.Result[k] = common.Decimal(v / total)
			}
		} else {
			for k, v := range result {
				statistic.Result[k] = v
			}
		}
		statistics = append(statistics, statistic)
	}

	common.SuccessResponse(ctx, http.StatusOK, statistics)
	return
}
