package models

// 类型
type QuestionType int

const (
	SingleChoice      QuestionType = iota + 1 // 单选
	MultipleChoice                            // 多选
	FillBlank                                 // 填空
	QuestionAndAnswer                         // 问答
)

// 问题
type Question struct {
	Model
	Title           string       `gorm:"type:varchar(30);not null;comment:标题;" json:"title" binding:"required"`
	Options         JSON         `gorm:"type:json;default:null;comment:选项;" json:"options,omitempty"`
	Type            QuestionType `gorm:"type:int;not null;comment:类型;" json:"type" binding:"required,min=1,max=4"`
	Priority        int          `gorm:"type:int;default:100;comment:优先级;" json:"priority"`
	QuestionnaireID uint         `gorm:"comment:问卷ID;" json:"questionnaire_id"`
	Answers         []Answer     `gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"answers,omitempty"`
}

// 选项
type Option struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// 表名
func (Question) TableName() string {
	return "questions"
}

// 创建问题
func CreateQuestion(question *Question) error {
	// 跳过自动创建
	// https://gorm.io/zh_CN/docs/associations.html
	return db.Omit("Answers").Create(question).Error
}

// 获取所有问题
func GetQuestions(conditions map[string]interface{}) ([]*Question, int64, error) {
	var tx = db
	var total int64
	var question []*Question

	for key, value := range conditions {
		switch key {
		case "limit":
			tx = tx.Limit(value.(int))
		case "offset":
			tx = tx.Offset(value.(int))
		case "type":
			tx = tx.Where("type = ?", value.(int))
		case "questionnaire_id":
			tx = tx.Where("questionnaire_id = ?", value.(int))
		}
	}
	err := tx.Find(&question).Limit(-1).Offset(-1).Count(&total).Error

	return question, total, err
}

// 获取单个问题
func GetQuestion(id string) (*Question, error) {
	var question Question
	err := db.Where("id = ?", id).First(&question).Error
	return &question, err
}

// 更新问题
func UpdateQuestion(id string, question *Question) error {
	err := db.Model(Question{}).Where("id = ?", id).Updates(question).Error
	return err
}

// 删除问题: 实现软删除的级联删除
func DeleteQuestion(id string) error {
	// 软删除不会真实的删除数据, 只是更新了DeleteAt时间戳
	// 数据库不支持软删除的级联删除, 故需要手动实现级联删除

	// 开启事务
	tx := db.Begin()

	// 删除答案
	err = tx.Where("answers.question_id = ?", id).Delete(&Answer{}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 删除问题
	err = tx.Where("id = ?", id).Delete(Question{}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
