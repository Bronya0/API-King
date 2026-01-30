package model

import "time"

type History struct {
	ID     uint   `gorm:"primaryKey"`
	Time   int    `json:"timestamp"`
	Method string `json:"method"`
	Path   string `json:"path"`
	DSL    string `json:"dsl"`
}

// Setting 设置表模型
type Setting struct {
	Key   string `gorm:"primaryKey"`
	Value string `json:value"`
}

// ApiRequest 表示一个接口请求记录
type ApiRequest struct {
	ID        uint      `gorm:"primaryKey"`
	GroupName string    `gorm:"column:group_name;uniqueIndex:idx_uni"` // 分组名称
	ApiName   string    `gorm:"column:api_name;uniqueIndex:idx_uni"`   // 接口名称
	Method    string    `gorm:"column:method;uniqueIndex:idx_uni"`     // 请求方法
	URL       string    `gorm:"column:url;uniqueIndex:idx_uni"`        // 请求 URL
	Headers   string    `gorm:"column:headers"`                        // 请求头（JSON 字符串）
	Params    string    `gorm:"column:params"`                         //
	Type      string    `gorm:"column:type"`                           //
	Body      string    `gorm:"column:body"`                           //
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;"`     // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;"`     // 更新时间（每次使用时更新）
}

func (m *ApiRequest) TableName() string {
	return "api" // return "schema.table"
}
