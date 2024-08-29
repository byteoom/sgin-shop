package model

import "time"

// Team 结构体定义了团队的基础信息
type Team struct {
	Id   uint   `gorm:"primary_key" json:"id"`            // ID 是团队的主键
	UUID string `gorm:"type:char(36);unique" json:"uuid"` // UUID 是团队的唯一标识符
	// 拥有者
	OwnerUuid string    `gorm:"type:char(36);index" json:"owner_uuid"` // 拥有者uuid
	Name      string    `gorm:"type:varchar(100);unique" json:"name"`  // Name 是团队的名称，它在系统中是唯一的
	Desc      string    `gorm:"type:varchar(255)" json:"desc"`         // Desc 是对团队的描述
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`      // CreatedAt 记录了团队创建的时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`      // UpdatedAt 记录了团队最后更新的时间
	IsActive  bool      `gorm:"default:true" json:"is_active"`         // IsActive 标识团队是否是活跃的
	Creater   string    `gorm:"type:char(36)" json:"creater"`          // Creater 标识团队的创建者
}

type UserTeamRes struct {
	Team
	Role          *Role `json:"role"`            // Role 是用户在团队中的角色
	IsCurrentTeam bool  `json:"is_current_team"` // IsCurrentTeam 标识用户当前选择这个团队
}

type TeamUserRes struct {
	User
	Role *Role `json:"role"` // Role 是用户在团队中的角色
}
