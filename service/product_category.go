package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategoryService struct {
}

func NewProductCategoryService() *ProductCategoryService {
	return &ProductCategoryService{}
}

func (s *ProductCategoryService) CreateCategory(ctx *app.Context, category *model.ProductCategory) error {
	category.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	category.UpdatedAt = category.CreatedAt
	category.Uuid = uuid.New().String()

	err := ctx.DB.Create(category).Error
	if err != nil {
		ctx.Logger.Error("Failed to create product category", err)
		return errors.New("failed to create product category")
	}
	return nil
}

func (s *ProductCategoryService) UpdateCategory(ctx *app.Context, category *model.ProductCategory) error {
	category.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	err := ctx.DB.Where("uuid = ?", category.Uuid).Updates(category).Error
	if err != nil {
		ctx.Logger.Error("Failed to update product category", err)
		return errors.New("failed to update product category")
	}
	return nil
}

func (s *ProductCategoryService) DeleteCategory(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.ProductCategory{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete product category", err)
		return errors.New("failed to delete product category")
	}
	return nil
}

func (s *ProductCategoryService) GetCategoryByUUID(ctx *app.Context, uuid string) (*model.ProductCategory, error) {
	category := &model.ProductCategory{}
	err := ctx.DB.Where("uuid = ?", uuid).First(category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("product category not found")
		}
		ctx.Logger.Error("Failed to get product category by UUID", err)
		return nil, errors.New("failed to get product category by UUID")
	}
	return category, nil
}

func (s *ProductCategoryService) GetCategoryList(ctx *app.Context) ([]*model.ProductCategory, error) {
	var categories []*model.ProductCategory
	err := ctx.DB.Find(&categories).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product categories list", err)
		return nil, errors.New("failed to get product categories list")
	}
	return categories, nil
}

// 获取所有分类
func (s *ProductCategoryService) GetAllCategory(ctx *app.Context) ([]*model.ProductCategory, error) {
	var categories []*model.ProductCategory
	err := ctx.DB.Find(&categories).Error
	if err != nil {
		ctx.Logger.Error("Failed to get all product categories", err)
		return nil, errors.New("failed to get all product categories")
	}
	return categories, nil
}
