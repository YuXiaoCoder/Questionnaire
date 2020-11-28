package routers

import (
	"BE/internal/controllers"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	v1Group := engine.Group("/api/v1")
	{
		// 登录
		v1Group.GET("/login", controllers.LoginHandler)

		// 获取所有用户
		v1Group.GET("/users", controllers.GetUsersHandler)
		// 获取单个用户
		v1Group.GET("/users/:id", controllers.GetUserHandler)
		// 创建用户
		v1Group.POST("/users", controllers.CreateUserHandler)
		// 更新用户
		v1Group.PUT("/users/:id", controllers.UpdateUserHandler)
		// 删除用户
		v1Group.DELETE("/users/:id", controllers.DeleteUserHandler)

		// 获取所有问卷
		v1Group.GET("/questionnaires", controllers.GetQuestionnairesHandler)
		// 获取单个问卷
		v1Group.GET("/questionnaires/:id", controllers.GetQuestionnaireHandler)
		// 创建问卷
		v1Group.POST("/questionnaires", controllers.CreateQuestionnaireHandler)
		// 更新问卷
		v1Group.PUT("/questionnaires/:id", controllers.UpdateQuestionnaireHandler)
		// 删除问卷
		v1Group.DELETE("/questionnaires/:id", controllers.DeleteQuestionnaireHandler)
		// 统计分析
		v1Group.GET("/questionnaires/:id/analysis", controllers.AnalyzeQuestionnaireHandler)

		// 获取所有问题
		v1Group.GET("/questions", controllers.GetQuestionsHandler)
		// 获取单个问题
		v1Group.GET("/questions/:id", controllers.GetQuestionHandler)
		// 创建问题
		v1Group.POST("/questions", controllers.CreateQuestionHandler)
		// 更新问题
		v1Group.PUT("/questions/:id", controllers.UpdateQuestionHandler)
		// 删除问题
		v1Group.DELETE("/questions/:id", controllers.DeleteQuestionHandler)

		// 获取所有答题卡
		v1Group.GET("/answersheets", controllers.GetAnswerSheetsHandler)
		// 获取单个答题卡
		v1Group.GET("/answersheets/:id", controllers.GetAnswerSheetHandler)
		// 创建答题卡
		v1Group.POST("/answersheets", controllers.CreateAnswerSheetHandler)
		// 更新答题卡
		v1Group.PUT("/answersheets/:id", controllers.UpdateAnswerSheetHandler)
		// 删除答题卡
		v1Group.DELETE("/answersheets/:id", controllers.DeleteAnswerSheetHandler)

		// 获取所有答案
		v1Group.GET("/answers", controllers.GetAnswersHandler)
		// 获取单个答案
		v1Group.GET("/answers/:id", controllers.GetAnswerHandler)
		// 创建答案
		v1Group.POST("/answers", controllers.CreateAnswerHandler)
		// 更新答案
		v1Group.PUT("/answers/:id", controllers.UpdateAnswerHandler)
		// 删除答案
		v1Group.DELETE("/answers/:id", controllers.DeleteAnswerHandler)
	}
}
