package controller

import (
	"net/http"
	"path/filepath"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

// UserController handles the operations related to User.
type UserController struct {
	Service           *service.UserService
	TeamMemberService *service.TeamMemberService
	MenuService       *service.MenuService
	OrderService      *service.OrderService
}

// CreateUser creates a new User.
// @Summary 创建用户
// @Description Create a new user with the input payload
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param user body model.User true "Create user"
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/create [post]
func (uc *UserController) CreateUser(c *app.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := uc.Service.CreateUser(c, &user)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(user)
}

// GetUserByUUID gets a User by UUID.
// @Summary 获取用户信息
// @Description Get a user by its UUID
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param params body model.ReqUserQueryParam false "查询参数"
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/info [post]
func (uc *UserController) GetUserByUUID(c *app.Context) {
	param := &model.ReqUserQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	user, err := uc.Service.GetUserByUUID(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(user)
}

// UpdateUser updates an existing User.
// @Summary 更新用户
// @Description Update a user with the input payload
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param user body model.User true "Update user"
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/update [post]
func (uc *UserController) UpdateUser(c *app.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := uc.Service.UpdateUser(c, &user)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(user)
}

// DeleteUser deletes a User by UUID.
// @Summary 删除用户
// @Description Delete a user by its UUID
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param params body model.ReqUserDeleteParam true "Delete user"
// @Success 200 {object} app.Response
// @Router /api/v1/user/delete [post]
func (uc *UserController) DeleteUser(c *app.Context) {
	params := &model.ReqUserDeleteParam{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	if params.Uuid == c.GetString("user_id") {
		c.JSONError(http.StatusBadRequest, "You can't delete yourself")
		return
	}

	err := uc.Service.DeleteUser(c, params.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess("User deleted successfully")
}

// 获取用户列表
// @Summary 获取用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param params body model.ReqUserQueryParam true "获取用户列表参数"
// @Success 200 {object} model.UserQueryResponse
// @Router /api/v1/user/list [post]
func (uc *UserController) GetUserList(c *app.Context) {
	param := &model.ReqUserQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	users, err := uc.Service.GetUserList(c, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(users)
}

// 获取自己的信息
// @Summary 获取自己的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/myinfo [get]
func (uc *UserController) GetMyInfo(c *app.Context) {
	userId := c.GetString("user_id")
	user, err := uc.Service.GetUserByUUID(c, userId)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(user)
}

// 修改头像
// @Summary 修改头像
// @Tags 用户
// @Accept json
// @Produce json
// @Param file formData file true "头像文件"
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/avatar [post]
func (uc *UserController) UpdateAvatar(c *app.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 保存头像

	userid := c.GetString("user_id")
	if userid == "" {
		c.JSONError(http.StatusBadRequest, "user_id is required")
		return
	}

	// 获取文件后缀

	extfile := filepath.Ext(file.Filename)

	filename := "/avatar/" + userid + extfile

	err = c.SaveUploadedFile(file, c.Config.Upload.Dir+filename)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	user := model.User{
		Avatar: filename,
		Uuid:   userid,
	}

	err = uc.Service.UpdateUser(c, &user)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(user)
}

// 获取所有用户

func (uc *UserController) GetAllUsers(c *app.Context) {
	users, err := uc.Service.GetAllUsers(c)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(users)
}

// GetUserTeamList
// @Summary 获取用户团队列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param params body model.ReqUserTeamQueryParam true "获取用户团队列表参数"
// @Success 200 {object} model.TeamListResponse
// @Router /api/v1/user/team/list [post]
func (uc *UserController) GetUserTeamList(c *app.Context) {
	param := &model.ReqUserTeamQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	teams, err := uc.TeamMemberService.GetUserTeamList(c, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(teams)
}

// GetMyTeamList
// @Summary 获取我的团队列表
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} model.TeamListResponse
// @Router /api/v1/user/teams [post]
func (uc *UserController) GetMyTeamList(c *app.Context) {
	userId := c.GetString("user_id")

	if userId == "" {
		c.JSONError(http.StatusBadRequest, "user_id is required")
		return

	}

	param := &model.ReqUserTeamQueryParam{
		UserUuid: userId,
	}
	teams, err := uc.TeamMemberService.GetUserTeamList(c, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(teams)
}

// SwitchTeam
// @Summary 切换团队
// @Tags 用户
// @Accept json
// @Produce json
// @Param params body model.ReqSwitchTeamParam true "切换团队参数"
// @Success 200 {object} model.StringDataResponse "ok"
// @Router /api/v1/user/team/switch [post]
func (uc *UserController) SwitchTeam(c *app.Context) {
	param := &model.ReqSwitchTeamParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	userId := c.GetString("user_id")
	if userId == "" {
		c.JSONError(http.StatusBadRequest, "user_id is required")
		return
	}

	err := uc.TeamMemberService.SwitchTeam(c, userId, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess("ok")
}

// GetUserMenu
// @Summary 获取用户菜单
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} model.MenuListResponse
// @Router /api/v1/user/menus [post]
func (uc *UserController) GetUserMenu(c *app.Context) {
	userId := c.GetString("user_id")
	if userId == "" {
		c.JSONError(http.StatusBadRequest, "user_id is required")
		return
	}

	menus, err := uc.MenuService.GetUserMenu(c, userId)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(menus)
}

// GetUserOrderList
// @Summary 获取用户订单列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param params body model.ReqOrderQueryParam true "获取用户订单列表参数"
// @Success 200 {object} model.OrderListPageResponse
// @Router /api/v1/user/orders [post]
func (uc *UserController) GetUserOrderList(c *app.Context) {
	param := &model.ReqOrderQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	userId := c.GetString("user_id")
	if userId == "" {
		c.JSONError(http.StatusBadRequest, "user_id is required")
		return
	}

	orders, err := uc.OrderService.GetOrderList(c, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(orders)
}
