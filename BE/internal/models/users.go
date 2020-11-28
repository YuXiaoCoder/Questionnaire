package models

import (
	"gorm.io/gorm/clause"
)

// 用户
type User struct {
	Model
	OpenID         string          `gorm:"type:varchar(30);comment:开放ID;" json:"open_id" binding:"required"`
	Questionnaires []Questionnaire `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:问题;" json:"questionnaires"`
	AnswerSheets   []AnswerSheet   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:答题卡;" json:"answer_sheets"`
}

// 表名
func (User) TableName() string {
	return "users"
}

// 创建用户
func CreateUser(user *User) error {
	// 跳过自动创建
	// https://gorm.io/zh_CN/docs/associations.html
	return db.Omit(clause.Associations).Create(user).Error
}

// 获取所有用户
func GetUsers(conditions map[string]interface{}) ([]*User, int64, error) {
	var tx = db
	var total int64
	var users []*User

	for key, value := range conditions {
		switch key {
		case "limit":
			tx = tx.Limit(value.(int))
		case "offset":
			tx = tx.Offset(value.(int))
		}
	}

	err := tx.Find(&users).Limit(-1).Offset(-1).Count(&total).Error

	return users, total, err
}

// 获取单个用户
func GetUser(id string) (*User, error) {
	var user User
	err := db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func GetUserByOpenID(openID string) (*User, error) {
	var user User
	err := db.Where("open_id = ?", openID).First(&user).Error
	return &user, err
}

// 更新用户
func UpdateUser(id string, user *User) error {
	err := db.Model(User{}).Where("id = ?", id).Updates(user).Error
	return err
}

// 删除用户: 实现软删除的级联删除
func DeleteUser(id string) error {
	// 软删除不会真实的删除数据, 只是更新了DeleteAt时间戳
	// 数据库不支持软删除的级联删除, 故需要手动实现级联删除

	// 开启事务
	tx := db.Begin()

	// 查询答案
	var answers []*Answer
	err = tx.Model(&Answer{}).Joins("LEFT JOIN answer_sheets ON answer_sheets.id AND answers.answer_sheet_id").Joins("LEFT JOIN users ON users.id = answer_sheets.user_id AND users.id = ?", id).Find(&answers).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 删除答案
	if len(answers) != 0 {
		err = tx.Delete(&answers).Error
		if err != nil {
			// 回滚事务
			tx.Rollback()
			return err
		}
	}

	// 删除答题卡
	err = tx.Where("user_id = ?", id).Delete(AnswerSheet{}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 查询问题
	var questions []*Question
	err = tx.Model(&Question{}).Joins("LEFT JOIN questionnaires ON questionnaires.id AND questions.questionnaire_id").Joins("LEFT JOIN users ON users.id = questionnaires.user_id AND users.id = ?", id).Find(&questions).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 删除问题
	if len(questions) != 0 {
		err = tx.Delete(&questions).Error
		if err != nil {
			// 回滚事务
			tx.Rollback()
			return err
		}
	}

	// 删除问卷
	err = tx.Where("user_id = ?", id).Delete(Questionnaire{}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 删除用户
	err := tx.Where("id = ?", id).Delete(User{}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
