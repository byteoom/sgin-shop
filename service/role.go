package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleService struct {
}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (s *RoleService) CreateRole(ctx *app.Context, role *model.Role) error {
	role.CreatedAt = time.Now()
	role.UpdatedAt = role.CreatedAt
	role.Uuid = uuid.New().String()

	err := ctx.DB.Create(role).Error
	if err != nil {
		ctx.Logger.Error("Failed to create role", err)
		return errors.New("failed to create role")
	}
	return nil
}

func (s *RoleService) GetRoleByUUID(ctx *app.Context, uuid string) (*model.Role, error) {
	role := &model.Role{}
	err := ctx.DB.Where("uuid = ?", uuid).First(role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		ctx.Logger.Error("Failed to get role by UUID", err)
		return nil, errors.New("failed to get role by UUID")
	}
	return role, nil
}

func (s *RoleService) UpdateRole(ctx *app.Context, role *model.Role) error {
	role.UpdatedAt = time.Now()
	err := ctx.DB.Save(role).Error
	if err != nil {
		ctx.Logger.Error("Failed to update role", err)
		return errors.New("failed to update role")
	}

	return nil
}

func (s *RoleService) DeleteRole(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.Role{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete role", err)
		return errors.New("failed to delete role")
	}

	return nil
}

// 查询角色列表
func (s *RoleService) GetRoleList(ctx *app.Context, param *model.ReqRoleQueryParam) (r *model.PagedResponse, err error) {
	roles := make([]*model.Role, 0)
	query := ctx.DB.Model(&model.Role{})
	query = query.Where("team_uuid = ?", param.TeamUuid)
	if param.Name != "" {
		query = query.Where("name like ?", "%"+param.Name+"%")
	}
	if param.IsActive {
		query = query.Where("is_active = ?", param.IsActive)
	}
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get role count", err)
		return nil, errors.New("failed to get role count")
	}
	err = query.Limit(param.PageSize).Offset(param.GetOffset()).Find(&roles).Error
	if err != nil {
		ctx.Logger.Error("Failed to get role list", err)
		return nil, errors.New("failed to get role list")
	}
	return &model.PagedResponse{
		Data:     roles,
		Current:  param.Current,
		PageSize: param.PageSize,
		Total:    total,
	}, nil
}

// 根据Uuids 列表获取角色列表
func (s *RoleService) GetRoleMapByUuids(ctx *app.Context, uuids []string) (map[string]*model.Role, error) {
	roles := make([]*model.Role, 0)
	err := ctx.DB.Where("uuid in (?)", uuids).Find(&roles).Error
	if err != nil {
		ctx.Logger.Error("Failed to get role list by uuids", err)
		return nil, errors.New("failed to get role list by uuids")
	}
	roleMap := make(map[string]*model.Role)
	for _, role := range roles {
		roleMap[role.Uuid] = role
	}
	return roleMap, nil
}
