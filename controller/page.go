package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type PageController struct {
	PageService *service.PageService
}

// @Tags Page
// @Summary 获取页面列表
// @Description 获取页面列表
// @Accept  json
// @Produce  json
// @Param params body model.ReqPageQueryParam false "查询参数"
// @Success 200 {object} model.PageQueryResponse
// @Router /api/v1/page/list [post]
func (pc *PageController) GetPageList(c *app.Context) {
	param := &model.ReqPageQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	pages, err := pc.PageService.GetPageList(c, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(pages)
}

// GetPageInfo 获取页面详情
// @Tags Page
// @Summary 获取页面详情
// @Description 获取页面详情
// @Accept  json
// @Produce  json
// @Param params body model.ReqUuidParam true "Get page info"
// @Success 200 {object} model.PageInfoResponse
// @Router /api/v1/page/info [post]
func (pc *PageController) GetPageInfo(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	page, err := pc.PageService.GetPageByUUID(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(page)
}

// @Tags Page
// @Summary 创建页面
// @Description 创建页面
// @Accept  json
// @Produce  json
// @Param params body model.Page true "Create page"
// @Success 200 {object} model.PageInfoResponse
// @Router /api/v1/page/create [post]
func (pc *PageController) CreatePage(c *app.Context) {
	var page model.ReqPageCreate
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.PageService.CreatePage(c, &page)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(page)
}

// @Tags Page
// @Summary 更新页面
// @Description 更新页面
// @Accept  json
// @Produce  json
// @Param params body model.Page true "Update page"
// @Success 200 {object} model.PageInfoResponse
// @Router /api/v1/page/update [post]
func (pc *PageController) UpdatePage(c *app.Context) {
	var page model.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.PageService.UpdatePage(c, &page)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(page)
}

// @Tags Page
// @Summary 删除页面
// @Description 删除页面
// @Accept  json
// @Produce  json
// @Param params body model.ReqUuidParam true "Delete page"
// @Success 200 {object} model.StringDataResponse "Successfully deleted page data"
// @Router /api/v1/page/delete [post]
func (pc *PageController) DeletePage(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.PageService.DeletePage(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess("删除成功")
}
