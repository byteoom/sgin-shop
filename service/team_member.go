package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamMemberService struct {
}

func NewTeamMemberService() *TeamMemberService {
	return &TeamMemberService{}
}

func (s *TeamMemberService) CreateTeamMember(ctx *app.Context, teamMember *model.TeamMember) error {
	teamMember.CreatedAt = time.Now()
	teamMember.UpdatedAt = teamMember.CreatedAt
	teamMember.UUID = uuid.New().String()

	// 先查询是否存在相同的团队成员
	var isExistTeamMember model.TeamMember
	err := ctx.DB.Where("team_uuid = ? AND user_uuid = ?", teamMember.TeamUUID, teamMember.UserUUID).First(&isExistTeamMember).Error
	if err == nil && isExistTeamMember.Id > 0 {
		return errors.New("team member already exists")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Logger.Error("Failed to get team member by team UUID and user UUID", err)
		return errors.New("failed to get team member by team UUID and user UUID")
	}
	teamMember.IsCurrentTeam = true

	// 查询用户是由有选中的团队存在
	var isExistCurrentTeam model.TeamMember
	err = ctx.DB.Where("user_uuid = ? AND is_current_team = ?", teamMember.UserUUID, true).First(&isExistCurrentTeam).Error
	if err == nil && isExistCurrentTeam.Id > 0 {
		teamMember.IsCurrentTeam = false
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Logger.Error("Failed to get team member by user UUID and is current team", err)
		return errors.New("failed to get team member by user UUID and is current team")
	}

	err = ctx.DB.Create(teamMember).Error
	if err != nil {
		ctx.Logger.Error("Failed to create team member", err)
		return errors.New("failed to create team member")
	}
	return nil
}

func (s *TeamMemberService) GetTeamMemberByUUID(ctx *app.Context, uuid string) (*model.TeamMember, error) {
	teamMember := &model.TeamMember{}
	err := ctx.DB.Where("uuid = ?", uuid).First(teamMember).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("team member not found")
		}
		ctx.Logger.Error("Failed to get team member by UUID", err)
		return nil, errors.New("failed to get team member by UUID")
	}
	return teamMember, nil
}

func (s *TeamMemberService) UpdateTeamMember(ctx *app.Context, teamMember *model.TeamMember) error {
	teamMember.UpdatedAt = time.Now()
	err := ctx.DB.Save(teamMember).Error
	if err != nil {
		ctx.Logger.Error("Failed to update team member", err)
		return errors.New("failed to update team member")
	}

	return nil
}

func (s *TeamMemberService) DeleteTeamMember(ctx *app.Context, uuid string) error {

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		teamMember := &model.TeamMember{}
		err := tx.Where("uuid = ?", uuid).First(teamMember).Error
		if err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return errors.New("team member not found")
			}

			ctx.Logger.Error("Failed to get team member by UUID", err)
			return errors.New("failed to get team member by UUID")
		}

		// 如果删除的是当前团队，需要重新选择一个团队
		if teamMember.IsCurrentTeam {
			var newCurrentTeam model.TeamMember
			err = tx.Where("user_uuid = ? AND team_uuid != ?", teamMember.UserUUID, teamMember.TeamUUID).
				First(&newCurrentTeam).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// 如果没有其他团队，直接删除
					err = tx.Delete(teamMember).Error
					if err != nil {
						tx.Rollback()
						ctx.Logger.Error("Failed to delete team member", err)
						return errors.New("failed to delete team member")
					}
					return nil
				}

				tx.Rollback()
				ctx.Logger.Error("Failed to get team member by user UUID and team UUID", err)
				return errors.New("failed to get team member by user UUID and team UUID")
			}

			newCurrentTeam.IsCurrentTeam = true
			err = tx.Save(&newCurrentTeam).Error
			if err != nil {
				tx.Rollback()
				ctx.Logger.Error("Failed to update team member", err)
				return errors.New("failed to update team member")
			}
		}

		err = tx.Delete(teamMember).Error
		if err != nil {
			tx.Rollback()
			ctx.Logger.Error("Failed to delete team member", err)
			return errors.New("failed to delete team member")
		}

		return nil
	})

	if err != nil {
		ctx.Logger.Error("Failed to delete team member", err)
		return errors.New("failed to delete team member")
	}

	return nil
}

// 获取团队成员用户列表
func (s *TeamMemberService) GetTeamMemberUserList(ctx *app.Context, params *model.ReqTeamMemberQueryParam) (*model.PagedResponse, error) {
	var teamMembers []*model.TeamMember
	var users []*model.User
	var userIds []string
	var total int64

	// 查找团队成员，并获取总数
	err := ctx.DB.Model(&model.TeamMember{}).Where("team_uuid = ?", params.TeamUUID).Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to count team members by team UUID", err)
		return nil, errors.New("failed to count team members by team UUID")
	}

	// 分页查询团队成员
	err = ctx.DB.Where("team_uuid = ?", params.TeamUUID).Offset(params.GetOffset()).Limit(params.PageSize).
		Find(&teamMembers).Error
	if err != nil {
		ctx.Logger.Error("Failed to get team members by team UUID", err)
		return nil, errors.New("failed to get team members by team UUID")
	}
	roleUuids := make([]string, 0)

	mTeamMember := make(map[string]*model.TeamMember)
	// 获取成员的用户UUID列表
	for _, teamMember := range teamMembers {
		userIds = append(userIds, teamMember.UserUUID)
		roleUuids = append(roleUuids, teamMember.Role)
		mTeamMember[teamMember.UserUUID] = teamMember
	}

	// 如果没有成员，直接返回空
	if len(userIds) == 0 {
		return &model.PagedResponse{
			Total:    total,
			Data:     []*model.TeamUserRes{},
			Current:  params.Current,
			PageSize: params.PageSize,
		}, nil
	}

	// 查找用户
	err = ctx.DB.Where("uuid in ?", userIds).Find(&users).Error
	if err != nil {
		ctx.Logger.Error("Failed to get users by UUIDs", err)
		return nil, errors.New("failed to get users by UUIDs")
	}

	roleMap, err := NewRoleService().GetRoleMapByUuids(ctx, roleUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get role map by UUIDs", err)
		return nil, errors.New("failed to get role map by UUIDs")
	}

	res := make([]*model.TeamUserRes, 0)

	// 隐藏密码
	for _, user := range users {
		user.Password = ""

		itemUser := &model.TeamUserRes{
			User: *user,
		}
		if teamMember, ok := mTeamMember[user.Uuid]; ok {
			if role, ok := roleMap[teamMember.Role]; ok {
				itemUser.Role = role
			}
		}
		res = append(res, itemUser)
	}

	return &model.PagedResponse{
		Total:    total,
		Data:     res,
		Current:  params.Current,
		PageSize: params.PageSize,
	}, nil
}

// GetUserTeamList
func (s *TeamMemberService) GetUserTeamList(ctx *app.Context, params *model.ReqUserTeamQueryParam) ([]*model.UserTeamRes, error) {
	var teamMembers []*model.TeamMember
	var teams []*model.Team
	var teamIds []string

	// 查找团队成员
	err := ctx.DB.Where("user_uuid = ?", params.UserUuid).Find(&teamMembers).Error
	if err != nil {
		ctx.Logger.Error("Failed to get team members by user UUID", err)
		return nil, errors.New("failed to get team members by user UUID")
	}

	// 获取团队UUID列表
	for _, teamMember := range teamMembers {
		teamIds = append(teamIds, teamMember.TeamUUID)
	}

	// 如果没有团队，直接返回空
	if len(teamIds) == 0 {
		return []*model.UserTeamRes{}, nil
	}

	roleUuids := make([]string, 0)
	mTeamMember := make(map[string]*model.TeamMember)
	for _, teamMember := range teamMembers {
		roleUuids = append(roleUuids, teamMember.Role)
		mTeamMember[teamMember.TeamUUID] = teamMember
	}

	roleMap, err := NewRoleService().GetRoleMapByUuids(ctx, roleUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get role map by UUIDs", err)
		return nil, errors.New("failed to get role map by UUIDs")
	}

	// 查找团队
	err = ctx.DB.Where("uuid in ?", teamIds).Find(&teams).Error
	if err != nil {
		ctx.Logger.Error("Failed to get teams by UUIDs", err)
		return nil, errors.New("failed to get teams by UUIDs")
	}

	res := make([]*model.UserTeamRes, 0)
	for _, team := range teams {
		resItem := &model.UserTeamRes{
			Team: *team,
		}
		if teamMember, ok := mTeamMember[team.UUID]; ok {

			if role, ok := roleMap[teamMember.Role]; ok {
				resItem.Role = role
			}

			resItem.IsCurrentTeam = teamMember.IsCurrentTeam
		}
		res = append(res, resItem)
	}

	return res, nil
}

// SwitchTeam
func (s *TeamMemberService) SwitchTeam(ctx *app.Context, userId string, params *model.ReqSwitchTeamParam) error {

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		// 查找用户的团队成员
		var teamMember model.TeamMember
		err := tx.Where("user_uuid = ? AND team_uuid = ?", userId, params.TeamUuid).First(&teamMember).Error
		if err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return errors.New("team member not found")
			}
			ctx.Logger.Error("Failed to get team member by user UUID and team UUID", err)
			return errors.New("failed to get team member by user UUID and team UUID")
		}

		// 更新当前团队
		teamMember.IsCurrentTeam = true
		teamMember.UpdatedAt = time.Now()
		err = tx.Where("user_uuid = ? AND team_uuid = ?", userId, params.TeamUuid).Updates(teamMember).Error
		if err != nil {
			tx.Rollback()
			ctx.Logger.Error("Failed to update team member", err)
			return errors.New("failed to update team member")
		}

		// 更新其他团队
		err = tx.Model(&model.TeamMember{}).Where("user_uuid = ? AND team_uuid != ?", userId, params.TeamUuid).Updates(map[string]interface{}{
			"is_current_team": false,
			"updated_at":      time.Now().Format("2006-01-02 15:04:05"),
		}).Error
		if err != nil {
			ctx.Logger.Error("Failed to update team member", err)
			return errors.New("failed to update team member")

		}

		return nil
	})

	if err != nil {
		ctx.Logger.Error("Failed to switch team", err)
		return errors.New("failed to switch team")
	}

	return nil
}

// 获取用户当前所在团队的角色
func (s *TeamMemberService) GetUserCurrentTeam(ctx *app.Context, userId string) (*model.TeamMember, error) {
	var teamMember model.TeamMember
	err := ctx.DB.Where("user_uuid = ? AND is_current_team = ?", userId, true).First(&teamMember).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("team member not found")
		}
		ctx.Logger.Error("Failed to get team member by user UUID and is current team", err)
		return nil, errors.New("failed to get team member by user UUID and is current team")
	}

	return &teamMember, nil
}
