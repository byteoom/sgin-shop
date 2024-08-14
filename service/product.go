package service

import (
	"errors"
	"fmt"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
}

// ProductCreate 创建产品
func (p *ProductService) ProductCreate(ctx *app.Context, params *model.ReqProductCreate) (err error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	productBase := model.ProductBase{
		Weight: params.Weight,
		Length: params.Length,
		Width:  params.Width,
		Height: params.Height,
		Unit:   params.Unit,
	}

	product := &model.Product{
		ProductBase:         productBase,
		Uuid:                uuid.New().String(),
		Name:                params.Name,
		Description:         params.Description,
		ProductCategoryUuid: params.ProductCategoryUuid,
		Images:              utils.ArrayToJsonString(params.Images),
		Videos:              utils.ArrayToJsonString(params.Videos),
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err = ctx.DB.Transaction(func(tx *gorm.DB) error {

		err := tx.Create(&product).Error
		if err != nil {
			ctx.Logger.Error("Failed to create product", err)
			tx.Rollback()
			return errors.New("failed to create product")
		}

		productItemList := make([]*model.ProductItem, 0)

		if len(params.Variants) > 0 {
			ritemList, err := p.CreateVariants(ctx, tx, params.Variants, product, params.VariantsVals)
			if err != nil {
				ctx.Logger.Error("Failed to create product variants", err)
				tx.Rollback()
				return errors.New("failed to create product variants")
			}
			productItemList = append(productItemList, ritemList...)
		} else {
			productItem := model.ProductItem{
				ProductBase:   productBase,
				Uuid:          uuid.New().String(),
				ProductUuid:   product.Uuid,
				Price:         params.Price,
				Discount:      params.Discount,
				DiscountPrice: params.DiscountPrice,
				Stock:         params.Stock,
				Description:   params.Description,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			productItemList = append(productItemList, &productItem)
		}
		err = tx.Create(&productItemList).Error
		if err != nil {
			ctx.Logger.Error("Failed to create product item", err)
			tx.Rollback()
			return errors.New("failed to create product item")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return
}

func (p *ProductService) CreateVariants(ctx *app.Context, tx *gorm.DB, variants []*model.ReqProductVariantsCreate, product *model.Product, vals []map[string]interface{}) (r []*model.ProductItem, err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	macc := make([]*model.ProductItem, 0)
	mVariants := make(map[string]string, 0)
	mVariantsOptions := make(map[string]map[string]string, 0)
	for _, variant := range variants {
		variantUuid := uuid.New().String()
		productVariant := model.ProductVariants{
			Uuid:        variantUuid,
			ProductUuid: product.Uuid,
			Name:        variant.Name,
			Description: variant.Description,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mVariants[variant.Name] = variantUuid

		err = tx.Create(&productVariant).Error
		if err != nil {
			return nil, err
		}

		mOptions := make(map[string]string, 0)

		for _, option := range variant.Options {
			optionUuid := uuid.New().String()
			productVariantOption := model.ProductVariantsOption{
				Uuid:        optionUuid,
				ProductUuid: product.Uuid,
				Name:        option,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			mOptions[option] = optionUuid
			err = tx.Create(&productVariantOption).Error
			if err != nil {
				return nil, err
			}

		}

		mVariantsOptions[variant.Name] = mOptions
	}

	for _, item := range vals {

		variantsList := make([]string, 0)
		for _, variantItem := range variants {
			if option, ok := item[variantItem.Name]; ok {
				variantsList = append(variantsList, fmt.Sprintf("%s:%s", variantItem.Name, option))
			}
		}

		variantsstr := strings.Join(variantsList, "-")

		productItem := model.ProductItem{
			ProductBase:   product.ProductBase,
			Uuid:          uuid.New().String(),
			Variants:      variantsstr,
			ProductUuid:   product.Uuid,
			Price:         p.GetFloat64ByMap(item, "price"),
			Discount:      p.GetFloat64ByMap(item, "discount"),
			DiscountPrice: p.GetFloat64ByMap(item, "discount_price"),
			Stock:         int64(p.GetFloat64ByMap(item, "stock")),
			Description:   utils.MapGetString(item, "description"),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		macc = append(macc, &productItem)
	}

	return macc, nil
}

// 根据map 字段获取float64
func (p *ProductService) GetFloat64ByMap(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if value, ok := v.(float64); ok {
			return value
		}
		if value, ok := v.(string); ok {
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return 0
			}
			return val
		}
	}
	return 0
}
