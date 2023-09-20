package models

type Model struct {
	ID      int32 `gorm:"column:id;type:int;primaryKey;autoIncrement:true" json:"id"`     // ID
	Created int32 `gorm:"column:created;type:int;not null;autoCreateTime" json:"created"` // 创建时间
	Updated int32 `gorm:"column:updated;type:int;not null;autoUpdateTime" json:"updated"` // 更新时间
	Deleted int32 `gorm:"column:deleted;type:int;not null" json:"deleted"`                // 删除时间
}
