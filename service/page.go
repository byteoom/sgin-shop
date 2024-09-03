package service

import (
	"errors"
	"strings"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PageService struct {
}

func NewPageService() *PageService {
	return &PageService{}
}

// CreatePage 创建新的页面
func (s *PageService) CreatePage(ctx *app.Context, req *model.ReqPageCreate) error {
	page := &model.Page{
		Title:  req.Title,
		Data:   req.Data,
		Ext:    req.Ext,
		Slug:   req.Slug,
		Status: model.PageStatusDraft,
	}
	page.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	page.UpdatedAt = page.CreatedAt
	page.UUID = uuid.New().String()

	err := ctx.DB.Create(page).Error
	if err != nil {

		// Duplicate entry slug
		if strings.Contains(err.Error(), "Duplicate entry '"+req.Slug+"' for key 'slug'") {
			ctx.Logger.Error("Slug already exists", err)
			return errors.New("slug already exists")
		}

		ctx.Logger.Error("Failed to create page", err)
		return errors.New("failed to create page")
	}
	return nil
}

// GetPageByUUID 根据 UUID 获取页面
func (s *PageService) GetPageByUUID(ctx *app.Context, uuid string) (*model.Page, error) {
	page := &model.Page{}
	err := ctx.DB.Where("uuid = ?", uuid).First(page).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("page not found")
		}
		ctx.Logger.Error("Failed to get page by UUID", err)
		return nil, errors.New("failed to get page by UUID")
	}
	return page, nil
}

// UpdatePage 更新页面信息
func (s *PageService) UpdatePage(ctx *app.Context, page *model.Page) error {
	page.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	err := ctx.DB.Where("uuid = ?", page.UUID).Updates(page).Error
	if err != nil {
		ctx.Logger.Error("Failed to update page", err)
		return errors.New("failed to update page")
	}

	return nil
}

// DeletePage 根据 UUID 删除页面
func (s *PageService) DeletePage(ctx *app.Context, uuidStr string) error {
	err := ctx.DB.Where("uuid = ?", uuidStr).Delete(&model.Page{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete page", err)
		return errors.New("failed to delete page")
	}

	return nil
}

// GetPageList 获取页面列表
func (s *PageService) GetPageList(ctx *app.Context, params *model.ReqPageQueryParam) (*model.PagedResponse, error) {
	var (
		pages []*model.Page
		total int64
	)

	db := ctx.DB.Model(&model.Page{})

	if params.Title != "" {
		db = db.Where("title LIKE ?", "%"+params.Title+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get page count", err)
		return nil, errors.New("failed to get page count")
	}

	err = db.Offset(params.GetOffset()).Limit(params.PageSize).Find(&pages).Error
	if err != nil {
		ctx.Logger.Error("Failed to get page list", err)
		return nil, errors.New("failed to get page list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  pages,
	}, nil
}
