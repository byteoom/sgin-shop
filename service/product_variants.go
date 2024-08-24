package service

import (
	"errors"
	"sgin/model"
	"sgin/pkg/app"
)

type ProductVariantService struct {
}

func NewProductVariantService() *ProductVariantService {
	return &ProductVariantService{}
}

// 根据产品uuid获取产品变体列表
func (s *ProductVariantService) GetProductVariantByProductUUID(ctx *app.Context, productUuid string) ([]*model.ProductVariants, error) {
	variants := make([]*model.ProductVariants, 0)
	err := ctx.DB.Where("product_uuid = ?", productUuid).Find(&variants).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product variants by product uuid", err)
		return nil, errors.New("failed to get product variants by product uuid")
	}

	return variants, nil
}

// 根据产品uuid获取产品变体Option列表
func (s *ProductVariantService) GetProductVariantOptionByProductUUID(ctx *app.Context, productUuid string) ([]*model.ProductVariantsOption, error) {
	variantOptions := make([]*model.ProductVariantsOption, 0)
	err := ctx.DB.Where("product_uuid = ?", productUuid).Find(&variantOptions).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product variant options by product uuid", err)
		return nil, errors.New("failed to get product variant options by product uuid")
	}

	return variantOptions, nil
}

// 根据产品uuid列表获取产品变体
func (s *ProductVariantService) GetProductVariantByProductUUIDList(ctx *app.Context, productUuids []string) (map[string][]*model.ProductVariants, error) {
	productVariants := make(map[string][]*model.ProductVariants)
	if len(productUuids) == 0 {
		return productVariants, nil
	}

	variants := make([]*model.ProductVariants, 0)
	err := ctx.DB.Where("product_uuid in (?)", productUuids).Find(&variants).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product variants by product uuid list", err)
		return nil, errors.New("failed to get product variants by product uuid list")
	}

	for _, variant := range variants {
		productVariants[variant.ProductUuid] = append(productVariants[variant.ProductUuid], variant)
	}

	return productVariants, nil
}

// 根据产品uuid列表获取产品变体Option
func (s *ProductVariantService) GetProductVariantOptionByProductUUIDList(ctx *app.Context, productUuids []string) (map[string][]*model.ProductVariantsOption, error) {
	productVariantOptions := make(map[string][]*model.ProductVariantsOption)
	if len(productUuids) == 0 {
		return productVariantOptions, nil
	}

	variantOptions := make([]*model.ProductVariantsOption, 0)
	err := ctx.DB.Where("product_uuid in (?)", productUuids).Find(&variantOptions).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product variant options by product uuid list", err)
		return nil, errors.New("failed to get product variant options by product uuid list")
	}

	for _, variantOption := range variantOptions {
		productVariantOptions[variantOption.ProductUuid] = append(productVariantOptions[variantOption.ProductUuid], variantOption)
	}

	return productVariantOptions, nil
}
