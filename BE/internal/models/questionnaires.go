package models

import (
	"strconv"

	"gorm.io/gorm/clause"
)

// 类型
type QuestionnaireType int

const (
	Vote QuestionnaireType = iota + 1 // 投票
	QN                                // 问卷
)

// 问卷
type Questionnaire struct {
	Model
	Title        string            `gorm:"type:varchar(30);not null;comment:标题;" json:"title" binding:"required"`
	Type         QuestionnaireType `gorm:"type:int;not null;comment:类型;" json:"type" binding:"required,min=1,max=2"`
	Status       *bool             `gorm:"type:boolean;default:false;comment:状态;" json:"status"`
	UserID       uint              `gorm:"comment:用户ID;" json:"user_id"`
	Questions    []Question        `gorm:"foreignKey:QuestionnaireID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:问题;" json:"questions,omitempty" binding:"dive"`
	AnswerSheets []AnswerSheet     `gorm:"foreignKey:QuestionnaireID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:答题卡;" json:"answer_sheets,omitempty" binding:"dive"`
}

// 表名
func (Questionnaire) TableName() string {
	return "questionnaires"
}

// 创建问卷
func CreateQuestionnaire(questionnaire *Questionnaire) error {
	// 跳过自动创建
	// https://gorm.io/zh_CN/docs/associations.html
	return db.Omit("AnswerSheets").Create(questionnaire).Error
}

// 获取所有问卷
func GetQuestionnaires(conditions map[string]interface{}) ([]*Questionnaire, int64, error) {
	var tx = db
	var total int64
	var questionnaires []*Questionnaire

	for key, value := range conditions {
		switch key {
		case "limit":
			tx = tx.Limit(value.(int))
		case "offset":
			tx = tx.Offset(value.(int))
		case "user_id":
			tx = tx.Where("user_id = ?", value.(int))
		case "type":
			tx = tx.Where("type = ?", value.(int))
		}
	}

	err := tx.Find(&questionnaires).Limit(-1).Offset(-1).Count(&total).Error
	return questionnaires, total, err
}

// 获取单个问卷
func GetQuestionnaire(id string) (*Questionnaire, error) {
	var questionnaire Questionnaire
	err := db.Preload(clause.Associations).Where("id = ?", id).First(&questionnaire).Error
	return &questionnaire, err
}

// 更新问卷
func UpdateQuestionnaire(id string, questionnaire *Questionnaire) error {
	// 更新问卷时, 支持更新问题
	tx := db.Begin()

	for i := 0; i < len(questionnaire.Questions); i++ {
		if questionnaire.Questions[i].ID == 0 {
			// 创建问题
			idU64, err := strconv.ParseUint(id, 10, 0)
			questionnaire.Questions[i].QuestionnaireID = uint(idU64)
			err = tx.Omit("Answers").Create(&questionnaire.Questions[i]).Error
			if err != nil {
				// 回滚事务
				tx.Rollback()
				return err
			}
		} else {
			// 更新问题
			err := tx.Model(Question{}).Where("id = ?", questionnaire.Questions[i].ID).Updates(questionnaire.Questions[i]).Error
			if err != nil {
				// 回滚事务
				tx.Rollback()
				return err
			}
		}
	}

	// 更新问卷
	err := tx.Model(Questionnaire{}).Where("id = ?", id).Updates(questionnaire).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// 删除问卷: 实现软删除的级联删除
func DeleteQuestionnaire(id string) error {
	// 软删除不会真实的删除数据, 只是更新了DeleteAt时间戳
	// 数据库不支持软删除的级联删除, 故需要手动实现级联删除

	// 开启事务
	tx := db.Begin()

	// 删除答案
	var answers []*Answer
	err = tx.Model(&Answer{}).Joins("LEFT JOIN answer_sheets ON answer_sheets.id AND answers.answer_sheet_id").Joins("LEFT JOIN questionnaires ON questionnaires.id = answer_sheets.id AND questionnaires.id = ?", id).Find(&answers).Error
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
	err = tx.Where("questionnaire_id = ?", id).Delete(AnswerSheet{}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 删除问题
	err = tx.Where("questionnaire_id = ?", id).Delete(Question{}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 删除问卷
	err := tx.Where("id = ?", id).Delete(Questionnaire{}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
