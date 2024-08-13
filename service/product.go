package service

import (
	"errors"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
}

// ProductCreate 创建产品
func (p *ProductService) ProductCreate(ctx *app.Context, params *model.ReqProductCreate) (err error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	product := model.Product{
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
			productVariantList := make([]*model.ProductVariants, 0)
			for _, variant := range params.Variants {
				productVariant := model.ProductVariants{
					Uuid:        uuid.New().String(),
					ProductID:   product.ID,
					Name:        variant.Name,
					Description: variant.Description,
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				productVariantList = append(productVariantList, &productVariant)

				if len(variant.Options) > 0 {

					optionList := make([]*model.ProductVariantsOption, 0)

					for _, option := range variant.Options {
						productVariantOption := model.ProductVariantsOption{
							Uuid:                uuid.New().String(),
							ProductVariantsUuid: productVariant.Uuid,
							Name:                option.Name,
							Unit:                option.Unit,
							Description:         option.Description,
							CreatedAt:           now,
							UpdatedAt:           now,
						}

						optionList = append(optionList, &productVariantOption)

						productItem := model.ProductItem{
							Uuid:                      uuid.New().String(),
							ProductUuid:               product.Uuid,
							ProductVariantsUuid:       productVariant.Uuid,
							ProductVariantsOptionUuid: productVariantOption.Uuid,
							Price:                     option.Price,
							Discount:                  option.Discount,
							DiscountPrice:             option.DiscountPrice,
							Stock:                     option.Stock,
							Description:               option.Description,
							CreatedAt:                 now,
							UpdatedAt:                 now,
						}
						productItemList = append(productItemList, &productItem)

					}
					err = tx.Create(&optionList).Error
					if err != nil {
						ctx.Logger.Error("Failed to create product variants option", err)
						tx.Rollback()
						return errors.New("failed to create product variants option")
					}
					continue
				}

				productItem := model.ProductItem{
					Uuid:                uuid.New().String(),
					ProductUuid:         product.Uuid,
					ProductVariantsUuid: productVariant.Uuid,
					Price:               variant.Price,
					Discount:            variant.Discount,
					DiscountPrice:       variant.DiscountPrice,
					Stock:               variant.Stock,
					Description:         variant.Description,
					CreatedAt:           now,
					UpdatedAt:           now,
				}
				productItemList = append(productItemList, &productItem)
			}
			err = tx.Create(&productVariantList).Error
			if err != nil {
				ctx.Logger.Error("Failed to create product variants", err)
				tx.Rollback()
				return errors.New("failed to create product variants")
			}
		} else {
			productItem := model.ProductItem{
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
