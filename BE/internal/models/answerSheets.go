package models

import (
	"strconv"

	"gorm.io/gorm/clause"
)

// 答题卡
type AnswerSheet struct {
	Model
	UserID          uint     `gorm:"comment:用户ID;" json:"user_id" binding:"required"`
	QuestionnaireID uint     `gorm:"comment:问卷ID;" json:"questionnaire_id" binding:"required"`
	Answers         []Answer `gorm:"foreignKey:AnswerSheetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"answers" binding:"dive"`
}

// 表名
func (AnswerSheet) TableName() string {
	return "answer_sheets"
}

// 创建答题卡
func CreateAnswerSheet(answerSheet *AnswerSheet) error {
	return db.Create(answerSheet).Error
}

// 获取所有答题卡
func GetAnswerSheets(conditions map[string]interface{}) ([]*AnswerSheet, int64, error) {
	var tx = db
	var total int64
	var answerSheets []*AnswerSheet

	for key, value := range conditions {
		switch key {
		case "limit":
			tx = tx.Limit(value.(int))
		case "offset":
			tx = tx.Offset(value.(int))
		case "user_id":
			tx = tx.Where("user_id = ?", value.(int))
		case "questionnaire_id":
			tx = tx.Where("questionnaire_id = ?", value.(int))
		}
	}

	err := tx.Preload(clause.Associations).Find(&answerSheets).Limit(-1).Offset(-1).Count(&total).Error
	return answerSheets, total, err
}

// 获取单个答题卡
func GetAnswerSheet(id string) (*AnswerSheet, error) {
	var answerSheet AnswerSheet
	err := db.Preload(clause.Associations).Where("id = ?", id).First(&answerSheet).Error
	return &answerSheet, err
}

// 更新答题卡
func UpdateAnswerSheet(id string, answerSheet *AnswerSheet) error {
	// 更新答题卡时, 支持更新答案
	tx := db.Begin()

	for i := 0; i < len(answerSheet.Answers); i++ {
		if answerSheet.Answers[i].ID == 0 {
			// 创建答案
			idU64, err := strconv.ParseUint(id, 10, 0)
			answerSheet.Answers[i].AnswerSheetID = uint(idU64)
			err = tx.Create(&answerSheet.Answers[i]).Error
			if err != nil {
				// 回滚事务
				tx.Rollback()
				return err
			}
		} else {
			// 更新答案
			err := tx.Model(Answer{}).Where("id = ?", answerSheet.Answers[i].ID).Updates(answerSheet.Answers[i]).Error
			if err != nil {
				// 回滚事务
				tx.Rollback()
				return err
			}
		}
	}

	// 更新答题卡
	err := tx.Model(AnswerSheet{}).Where("id = ?", id).Updates(answerSheet).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// 删除答题卡: 实现软删除的级联删除
func DeleteAnswerSheet(id string) error {
	// 软删除不会真实的删除数据, 只是更新了DeleteAt时间戳
	// 数据库不支持软删除的级联删除, 故需要手动实现级联删除

	// 开启事务
	tx := db.Begin()

	// 删除答案
	var answers []*Answer
	err = tx.Model(&Answer{}).Joins("LEFT JOIN answer_sheets ON answer_sheets.id AND answers.answer_sheet_id AND answer_sheets.id = ?", id).Find(&answers).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	if len(answers) != 0 {
		err = tx.Delete(&answers).Error
		if err != nil {
			// 回滚事务
			tx.Rollback()
			return err
		}
	}

	// 删除答题卡
	err = tx.Where("id = ?", id).Delete(AnswerSheet{}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
