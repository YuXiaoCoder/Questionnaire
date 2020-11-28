package models

import (
	"encoding/json"
	"reflect"
)

// 答案
type Answer struct {
	Model
	Value         JSON `gorm:"type:json;default:null;comment:值;" json:"value" binding:"required"`
	QuestionID    uint `gorm:"comment:问题ID;" json:"question_id" binding:"required"`
	AnswerSheetID uint `gorm:"comment:答题卡ID;" json:"answer_sheet_id"`
}

type AnswerStatistic struct {
	Value JSON    `json:"value"` // 答案
	Count float64 `json:"count"` // 次数
}

// 表名
func (Answer) TableName() string {
	return "answers"
}

// 创建答案
func CreateAnswer(answer *Answer) error {
	return db.Create(answer).Error
}

// 获取所有答案
func GetAnswers(conditions map[string]interface{}) ([]*Answer, int64, error) {
	var tx = db
	var total int64
	var answer []*Answer

	for key, value := range conditions {
		switch key {
		case "limit":
			tx = tx.Limit(value.(int))
		case "offset":
			tx = tx.Offset(value.(int))
		case "question_id":
			tx = tx.Where("question_id = ?", value.(int))
		case "answer_sheet_id":
			tx = tx.Where("answer_sheet_id = ?", value.(int))
		}
	}
	err := tx.Find(&answer).Limit(-1).Offset(-1).Count(&total).Error

	return answer, total, err
}

// 获取单个答案
func GetAnswer(id string) (*Answer, error) {
	var answer Answer
	err := db.Where("id = ?", id).First(&answer).Error
	return &answer, err
}

// 更新答案
func UpdateAnswer(id string, answer *Answer) error {
	err := db.Model(Answer{}).Where("id = ?", id).Updates(answer).Error
	return err
}

// 删除答案
func DeleteAnswer(id string) error {
	return db.Where("id = ?", id).Delete(Answer{}).Error
}

type Result struct {
	Value JSON `json:"value"`
	Count int  `json:"count"`
}

// 统计答案
func StatisticalAnswer(questionID uint) (map[string]float64, float64, error) {
	// 统计结果
	var total float64
	result := make(map[string]float64)

	var statistics []AnswerStatistic
	err := db.Model(Answer{}).Select("value, COUNT(*) as count").Where("question_id = ?", questionID).Group("value").Find(&statistics).Error
	if err != nil {
		return nil, total, err
	}

	for _, item := range statistics {
		data, _ := item.Value.MarshalJSON()

		var keys interface{}
		err := json.Unmarshal(data, &keys)
		if err != nil {
			return nil, total, err
		}

		// 通过反射获取类型
		if reflect.TypeOf(keys).Kind() == reflect.String {
			key := reflect.ValueOf(keys).Interface().(string)
			result[key] = item.Count
			total += item.Count
		} else if reflect.TypeOf(keys).Kind() == reflect.Slice {
			valueOfKeys := reflect.ValueOf(keys)
			for i := 0; i < valueOfKeys.Len(); i++ {
				key := valueOfKeys.Index(i).Interface().(string)
				_, exist := result[key]
				if exist {
					result[key] += item.Count
				} else {
					result[key] = item.Count
				}
				total += item.Count
			}
		}
	}
	return result, total, nil
}
