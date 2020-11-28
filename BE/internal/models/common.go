package models

import (
	"BE/internal/common"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 基础模型
type Model struct {
	ID        uint           `gorm:"int;primaryKey;comment:主键;" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间;" json:"created_at"`
	UpdatedAt time.Time      `gorm:"comment:更新时间;" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"comment:删除时间;" json:"deleted_at" sql:"index"`
}

// 全局变量
var (
	db  *gorm.DB
	err error
)

// 初始化函数
func init() {
	// 获取配置
	datasource := common.GetDatasource()

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", datasource.Username, datasource.Password, datasource.Host, datasource.Port, datasource.Database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, // 缓存预编译语句
	})
	if err != nil {
		panic(err)
	}

	// 调试模式
	if datasource.Debug {
		db = db.Debug()
	}

	// 连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池中连接的最大闲置数
	sqlDB.SetMaxIdleConns(10)

	// 设置数据库的最大连接数
	sqlDB.SetMaxOpenConns(100)

	// 设置连接的最大可复用时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构, 保持迁移顺序
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(&User{}, &Questionnaire{}, &Question{}, &AnswerSheet{}, &Answer{})
	if err != nil {
		panic(err)
	}
}

// 数据类型: JSON
// https://gorm.io/zh_CN/docs/data_types.html
// https://github.com/go-gorm/gorm/issues/1935
type JSON json.RawMessage

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 || string(j) == "null" {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	result := json.RawMessage{}
	err := json.Unmarshal(data, &result)
	*j = JSON(result)
	return err
}
