package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
}

func NewProductService() *ProductService {
	return &ProductService{}
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
		AliasName:           params.AliasName,
		ProductType:         params.ProductType,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err = ctx.DB.Transaction(func(tx *gorm.DB) error {

		// 查询产品别名是否存在
		var oldProduct model.Product
		err := tx.Where("alias_name = ?", params.AliasName).First(&oldProduct).Error
		if err == nil && oldProduct.Uuid != "" {
			return errors.New("产品别名已存在，请重新输入")
		}

		if err != nil && err != gorm.ErrRecordNotFound {
			ctx.Logger.Error("Failed to get product by alias name", err)
			tx.Rollback()
			return errors.New("failed to get product by alias name")
		}

		err = tx.Create(&product).Error
		if err != nil {
			ctx.Logger.Error("Failed to create product", err)
			tx.Rollback()
			return errors.New("failed to create product")
		}

		productItemList := make([]*model.ProductItem, 0)

		if params.ProductType == model.ProductTypeVariant && len(params.Variants) > 0 { // 变体产品
			ritemList, err := p.CreateVariants(ctx, tx, params.Variants, product, params.VariantsVals)
			if err != nil {
				ctx.Logger.Error("Failed to create product variants", err)
				tx.Rollback()
				return errors.New("failed to create product variants")
			}
			productItemList = append(productItemList, ritemList...)
		}

		if params.ProductType == model.ProductTypeSingle { // 单个产品
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
				Uuid:                optionUuid,
				ProductUuid:         product.Uuid,
				ProductVariantsUuid: variantUuid,
				Name:                option,
				CreatedAt:           now,
				UpdatedAt:           now,
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

// 获取产品列表
func (p *ProductService) ProductList(ctx *app.Context, params *model.ReqProductQueryParam) (r *model.PagedResponse, err error) {
	productList := make([]*model.Product, 0)
	query := ctx.DB.Model(&model.Product{})

	var total int64

	if params.Name != "" {
		query = query.Where("name like ?", "%"+params.Name+"%")
	}

	err = query.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product count", err)
		return nil, errors.New("failed to get product count")
	}

	err = query.Order("id DESC").Limit(params.PageSize).Offset(params.GetOffset()).Limit(params.PageSize).Find(&productList).Error

	if err != nil {
		ctx.Logger.Error("Failed to get product list", err)
		return nil, errors.New("failed to get product list")
	}

	imageUuids := make([]string, 0)
	mImage := make(map[string]bool, 0)
	mProductImages := make(map[string][]string, 0)
	for _, product := range productList {
		if product.Images == "" {
			continue
		}
		var images []string
		err = json.Unmarshal([]byte(product.Images), &images)
		if err != nil {
			ctx.Logger.Error("Failed to get product images", err)
			return nil, errors.New("failed to get product images")
		}

		for _, image := range images {
			if _, ok := mImage[image]; !ok {
				mImage[image] = true
				imageUuids = append(imageUuids, image)
			}
		}
		mProductImages[product.Uuid] = images
	}

	resourceMap, err := NewResourceService().GetResourceByUUIDList(ctx, imageUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	res := make([]*model.ProductRes, 0)

	for _, product := range productList {
		productRes := &model.ProductRes{
			Product:         *product,
			ProductCategory: model.ProductCategory{},
		}
		if images, ok := mProductImages[product.Uuid]; ok {
			for _, image := range images {
				if resource, ok := resourceMap[image]; ok {
					productRes.ImageList = append(productRes.ImageList, resource.Address)
				}
			}
		}
		res = append(res, productRes)
	}

	return &model.PagedResponse{
		Total: total,
		Data:  res,
	}, nil
}

// 删除产品 ， 根据uuid列表
func (p *ProductService) DeleteProductByUUIDList(ctx *app.Context, uuidList []string) (err error) {

	err = ctx.DB.Transaction(func(tx *gorm.DB) error {

		err := tx.Where("uuid IN (?)", uuidList).Delete(&model.Product{}).Error
		if err != nil {
			ctx.Logger.Error("Failed to delete product by UUID list", err)
			tx.Rollback()
			return errors.New("failed to delete product by UUID list")
		}

		err = tx.Where("product_uuid IN (?)", uuidList).Delete(&model.ProductItem{}).Error
		if err != nil {
			ctx.Logger.Error("Failed to delete product item by UUID list", err)
			tx.Rollback()
			return errors.New("failed to delete product item by UUID list")
		}

		// 删除产品变体

		err = tx.Where("product_uuid IN (?)", uuidList).Find(&model.ProductVariants{}).Error
		if err != nil {
			ctx.Logger.Error("Failed to get product variants by product uuid list", err)
			tx.Rollback()

			return errors.New("failed to get product variants by product uuid list")
		}

		err = tx.Where("product_uuid IN (?)", uuidList).Delete(&model.ProductVariantsOption{}).Error
		if err != nil {
			ctx.Logger.Error("Failed to delete product variants option by product uuid list", err)
			tx.Rollback()
			return errors.New("failed to delete product variants option by product uuid list")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// GetProductByUUIDList 根据uuid列表获取产品列表
func (p *ProductService) GetProductByUUIDList(ctx *app.Context, uuidList []string) (map[string]*model.ProductRes, error) {
	var products []*model.Product
	err := ctx.DB.Where("uuid IN (?)", uuidList).Find(&products).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product list by UUID list", err)
		return nil, errors.New("failed to get product list by UUID list")
	}

	imageUuids := make([]string, 0)
	mImage := make(map[string]bool, 0)
	mProductImages := make(map[string][]string, 0)
	for _, product := range products {
		if product.Images == "" {
			continue
		}
		var images []string
		err = json.Unmarshal([]byte(product.Images), &images)
		if err != nil {
			ctx.Logger.Error("Failed to get product images", err)
			return nil, errors.New("failed to get product images")
		}

		for _, image := range images {
			if _, ok := mImage[image]; !ok {
				mImage[image] = true
				imageUuids = append(imageUuids, image)
			}
		}
		mProductImages[product.Uuid] = images
	}

	resourceMap, err := NewResourceService().GetResourceByUUIDList(ctx, imageUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	res := make(map[string]*model.ProductRes, 0)

	for _, product := range products {
		productRes := &model.ProductRes{
			Product:         *product,
			ProductCategory: model.ProductCategory{},
		}
		if images, ok := mProductImages[product.Uuid]; ok {
			for _, image := range images {
				if resource, ok := resourceMap[image]; ok {
					productRes.ImageList = append(productRes.ImageList, resource.Address)
				}
			}
		}
		res[product.Uuid] = productRes
	}

	return res, nil
}

// 根据产品uuid列表获取产品item列表
func (p *ProductService) GetProductItemByProductUUIDList(ctx *app.Context, prodouctUuids []string) (r map[string][]*model.ProductItem, err error) {

	productItemList := make([]*model.ProductItem, 0)
	err = ctx.DB.Where("product_uuid IN (?)", prodouctUuids).Find(&productItemList).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product item by product uuid list", err)
		return nil, errors.New("failed to get product item by product uuid list")
	}

	mProductItem := make(map[string][]*model.ProductItem, 0)
	for _, productItem := range productItemList {
		if _, ok := mProductItem[productItem.ProductUuid]; !ok {
			mProductItem[productItem.ProductUuid] = make([]*model.ProductItem, 0)
		}
		mProductItem[productItem.ProductUuid] = append(mProductItem[productItem.ProductUuid], productItem)
	}

	return mProductItem, nil
}

// 获取产品sku列表
func (p *ProductService) GetProductSkuList(ctx *app.Context, params *model.ReqProductQueryParam) (r *model.PagedResponse, err error) {

	productList := make([]*model.ProductItem, 0)
	query := ctx.DB.Model(&model.ProductItem{})

	var total int64

	err = query.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product count", err)
		return nil, errors.New("failed to get product count")
	}

	err = query.Order("id DESC").Limit(params.PageSize).Offset(params.GetOffset()).Limit(params.PageSize).Find(&productList).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product list", err)
		return nil, errors.New("failed to get product list")
	}

	productUuids := make([]string, 0)
	mProduct := make(map[string]bool, 0)
	imageUuids := make([]string, 0)
	mImage := make(map[string]bool, 0)
	mProductImages := make(map[string][]string, 0)
	for _, product := range productList {
		if _, ok := mProduct[product.ProductUuid]; !ok {
			mProduct[product.ProductUuid] = true
			productUuids = append(productUuids, product.ProductUuid)
		}

		if product.Images != "" {
			var images []string
			err = json.Unmarshal([]byte(product.Images), &images)
			if err != nil {
				ctx.Logger.Error("Failed to get product images", err)
				return nil, errors.New("failed to get product images")
			}
			for _, image := range images {
				if _, ok := mImage[image]; !ok {
					mImage[image] = true
					imageUuids = append(imageUuids, image)
				}
			}
			mProductImages[product.Uuid] = images
		}
	}

	productMap, err := p.GetProductByUUIDList(ctx, productUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get product list by UUID list", err)
		return nil, errors.New("failed to get product list by UUID list")
	}

	resourceMap, err := NewResourceService().GetResourceByUUIDList(ctx, imageUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	res := make([]*model.ProductItemRes, 0)

	for _, product := range productList {
		productRes := &model.ProductItemRes{
			ProductItem: *product,
			ImageList:   make([]string, 0),
		}
		if product, ok := productMap[product.ProductUuid]; ok {
			productRes.ProductInfo = product
		}
		if images, ok := mProductImages[product.Uuid]; ok {
			for _, image := range images {
				if resource, ok := resourceMap[image]; ok {
					productRes.ImageList = append(productRes.ImageList, resource.Address)
				}
			}
		}
		res = append(res, productRes)
	}

	return &model.PagedResponse{
		Total: total,
		Data:  res,
	}, nil
}

// DeleteProductSkuByUUIDList
func (p *ProductService) DeleteProductSkuByUUIDList(ctx *app.Context, uuidList []string) (err error) {
	err = ctx.DB.Where("uuid IN (?)", uuidList).Delete(&model.ProductItem{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete product item by UUID list", err)
		return errors.New("failed to delete product item by UUID list")
	}
	return nil
}

// GetProductInfo
func (p *ProductService) GetProductInfo(ctx *app.Context, uuid string) (r *model.ProductRes, err error) {
	product := &model.Product{}
	err = ctx.DB.Where("uuid = ?", uuid).First(&product).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product by UUID", err)
		return nil, errors.New("failed to get product by UUID")
	}

	imageUuids := make([]string, 0)
	mImage := make(map[string]bool, 0)

	var images []string

	if product.Images != "" {
		err = json.Unmarshal([]byte(product.Images), &images)
		if err != nil {
			ctx.Logger.Error("Failed to get product images", err)
			return nil, errors.New("failed to get product images")
		}
	}

	for _, image := range images {
		if _, ok := mImage[image]; !ok {
			mImage[image] = true
			imageUuids = append(imageUuids, image)
		}
	}

	resourceMap, err := NewResourceService().GetResourceByUUIDList(ctx, imageUuids)
	if err != nil {

		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	productRes := &model.ProductRes{
		Product:         *product,
		ProductCategory: model.ProductCategory{},
		ImageList:       make([]string, 0),
	}

	for _, image := range images {
		if resource, ok := resourceMap[image]; ok {
			productRes.ImageList = append(productRes.ImageList, resource.Address)
		}
	}

	return productRes, nil
}

// GetShowProductInfo
func (p *ProductService) GetShowProductInfo(ctx *app.Context, uuid string) (r *model.ProductShowItem, err error) {
	product := &model.Product{}
	err = ctx.DB.Where("uuid = ?", uuid).First(&product).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product by UUID", err)
		return nil, errors.New("failed to get product by UUID")
	}

	imageUuids := make([]string, 0)
	mImage := make(map[string]bool, 0)

	var images []string

	if product.Images != "" {
		err = json.Unmarshal([]byte(product.Images), &images)
		if err != nil {
			ctx.Logger.Error("Failed to get product images", err)
			return nil, errors.New("failed to get product images")
		}
	}

	for _, image := range images {
		if _, ok := mImage[image]; !ok {
			mImage[image] = true
			imageUuids = append(imageUuids, image)
		}
	}

	resourceMap, err := NewResourceService().GetResourceByUUIDList(ctx, imageUuids)
	if err != nil {

		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	variants, err := NewProductVariantService().GetProductVariantByProductUUID(ctx, uuid)
	if err != nil {
		ctx.Logger.Error("Failed to get product variants by product uuid", err)
		return nil, errors.New("failed to get product variants by product uuid")
	}

	variantsOptions, err := NewProductVariantService().GetProductVariantOptionByProductUUID(ctx, uuid)
	if err != nil {
		ctx.Logger.Error("Failed to get product variant options by product uuid", err)
		return nil, errors.New("failed to get product variant options by product uuid")
	}

	productItems, err := p.GetProductItemByProductUUID(ctx, uuid)
	if err != nil {
		ctx.Logger.Error("Failed to get product item by product uuid", err)
		return nil, errors.New("failed to get product item by product uuid")
	}

	productShow := &model.ProductShow{
		ProductUuid:         product.Uuid,
		ProductType:         product.ProductType,
		Name:                product.Name,
		Description:         product.Description,
		Images:              images,
		Videos:              []string{},
		ProductCategoryUuid: product.ProductCategoryUuid,
	}

	if len(productItems) > 0 {
		sort.Sort(model.ProductItemResByPrice(productItems))
		productShow.Price = productItems[0].Price
		productShow.ProductItemUuid = productItems[0].Uuid
	}

	productShowItem := &model.ProductShowItem{
		ProductShow:           *productShow,
		ProductVariants:       variants,
		ProductVariantsOption: variantsOptions,
		ProductItems:          productItems,
	}

	productShowItem.Images = make([]string, 0)
	for _, image := range images {
		if resource, ok := resourceMap[image]; ok {
			productShowItem.Images = append(productShowItem.Images, resource.Address)
		}
	}

	return productShowItem, nil
}

// GetProductSkuInfo
func (p *ProductService) GetProductSkuInfo(ctx *app.Context, uuid string) (r *model.ProductItemRes, err error) {
	product := &model.ProductItem{}
	err = ctx.DB.Where("uuid = ?", uuid).First(&product).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product item by UUID", err)
		return nil, errors.New("failed to get product item by UUID")
	}

	productMap, err := p.GetProductByUUIDList(ctx, []string{product.ProductUuid})
	if err != nil {
		ctx.Logger.Error("Failed to get product list by UUID list", err)
		return nil, errors.New("failed to get product list by UUID list")
	}

	imageUuids := make([]string, 0)
	mImage := make(map[string]bool, 0)
	mProductImages := make(map[string][]string, 0)

	if product.Images != "" {
		var images []string
		err = json.Unmarshal([]byte(product.Images), &images)
		if err != nil {
			ctx.Logger.Error("Failed to get product images", err)
			return nil, errors.New("failed to get product images")
		}
		for _, image := range images {
			if _, ok := mImage[image]; !ok {
				mImage[image] = true
				imageUuids = append(imageUuids, image)
			}
		}
		mProductImages[product.Uuid] = images
	}

	productRes := &model.ProductItemRes{
		ProductItem: *product,
		ImageList:   make([]string, 0),
	}
	if product, ok := productMap[product.ProductUuid]; ok {
		productRes.ProductInfo = product
	}

	resourceMap, err := NewResourceService().GetResourceByUUIDList(ctx, imageUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	if images, ok := mProductImages[product.Uuid]; ok {
		for _, image := range images {
			if resource, ok := resourceMap[image]; ok {
				productRes.ImageList = append(productRes.ImageList, resource.Address)
			}
		}
	}

	return productRes, nil
}

// 根据产品uuid获取产品item列表
func (p *ProductService) GetProductItemByProductUUID(ctx *app.Context, productUuid string) (r []*model.ProductItemRes, err error) {
	productItemList := make([]*model.ProductItem, 0)
	err = ctx.DB.Where("product_uuid = ?", productUuid).Find(&productItemList).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product item by product uuid", err)
		return nil, errors.New("failed to get product item by product uuid")
	}
	imageUuids := make([]string, 0)
	mImage := make(map[string]bool, 0)

	for _, product := range productItemList {
		if product.Images != "" {
			var images []string
			err = json.Unmarshal([]byte(product.Images), &images)
			if err != nil {
				ctx.Logger.Error("Failed to get product images", err)
				return nil, errors.New("failed to get product images")
			}
			for _, image := range images {
				if _, ok := mImage[image]; !ok {
					mImage[image] = true
					imageUuids = append(imageUuids, image)
				}
			}
		}
	}

	resourceMap, err := NewResourceService().GetResourceByUUIDList(ctx, imageUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	res := make([]*model.ProductItemRes, 0)

	for _, product := range productItemList {
		itemProduct := &model.ProductItemRes{
			ProductItem: *product,
			ImageList:   make([]string, 0),
		}

		if product.Images != "" {

			var images []string
			err = json.Unmarshal([]byte(product.Images), &images)
			if err != nil {
				ctx.Logger.Error("Failed to get product images", err)
				return nil, errors.New("failed to get product images")
			}

			for _, image := range images {
				if resource, ok := resourceMap[image]; ok {
					itemProduct.ImageList = append(itemProduct.ImageList, resource.Address)
				}

			}
		}

		res = append(res, itemProduct)

	}

	return res, nil
}

// 根据产品item uuid列表获取产品item
func (p *ProductService) GetProductItemByUUIDList(ctx *app.Context, uuidList []string) (map[string]*model.ProductItemRes, error) {
	productItemList := make([]*model.ProductItem, 0)
	err := ctx.DB.Where("uuid IN (?)", uuidList).Find(&productItemList).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product item list by UUID list", err)
		return nil, errors.New("failed to get product item list by UUID list")
	}

	imageUuids := make([]string, 0)
	mImage := make(map[string]bool, 0)
	mProductImages := make(map[string][]string, 0)

	for _, product := range productItemList {

		if product.Images != "" {
			var images []string
			err = json.Unmarshal([]byte(product.Images), &images)
			if err != nil {
				ctx.Logger.Error("Failed to get product images", err)
				return nil, errors.New("failed to get product images")
			}
			for _, image := range images {
				if _, ok := mImage[image]; !ok {
					mImage[image] = true
					imageUuids = append(imageUuids, image)
				}
			}
			mProductImages[product.Uuid] = images
		}
	}

	resourceMap, err := NewResourceService().GetResourceByUUIDList(ctx, imageUuids)

	if err != nil {
		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	res := make(map[string]*model.ProductItemRes, 0)

	for _, product := range productItemList {
		productRes := &model.ProductItemRes{
			ProductItem: *product,
			ImageList:   make([]string, 0),
		}

		if images, ok := mProductImages[product.Uuid]; ok {
			for _, image := range images {
				if resource, ok := resourceMap[image]; ok {
					productRes.ImageList = append(productRes.ImageList, resource.Address)
				}
			}
		}
		res[product.Uuid] = productRes
	}

	return res, nil
}

// 获取前端展示得产品列表
func (p *ProductService) GetShowProductList(ctx *app.Context, params *model.ReqProductQueryParam) (r *model.PagedResponse, err error) {

	productList := make([]*model.Product, 0)
	query := ctx.DB.Model(&model.Product{})

	var total int64

	if params.Name != "" {
		query = query.Where("name like ?", "%"+params.Name+"%")
	}

	err = query.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product count", err)
		return nil, errors.New("failed to get product count")
	}

	err = query.Order("id DESC").Limit(params.PageSize).Offset(params.GetOffset()).Limit(params.PageSize).Find(&productList).Error

	if err != nil {
		ctx.Logger.Error("Failed to get product list", err)
		return nil, errors.New("failed to get product list")
	}

	imageUuids := make([]string, 0)
	mImage := make(map[string]bool, 0)
	mProductImages := make(map[string][]string, 0)
	for _, product := range productList {
		if product.Images == "" {
			continue
		}
		var images []string
		err = json.Unmarshal([]byte(product.Images), &images)
		if err != nil {
			ctx.Logger.Error("Failed to get product images", err)
			return nil, errors.New("failed to get product images")
		}

		for _, image := range images {
			if _, ok := mImage[image]; !ok {
				mImage[image] = true
				imageUuids = append(imageUuids, image)
			}
		}
		mProductImages[product.Uuid] = images
	}

	resourceMap, err := NewResourceService().GetResourceByUUIDList(ctx, imageUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	res := make([]*model.ProductShow, 0)

	productUuids := make([]string, 0)

	for _, product := range productList {
		productUuids = append(productUuids, product.Uuid)
	}

	productItemsMap, err := p.GetProductItemByProductUUIDList(ctx, productUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get product item by product uuid list", err)
		return nil, errors.New("failed to get product item by product uuid list")
	}

	for _, product := range productList {
		productRes := &model.ProductShow{

			ProductUuid:         product.Uuid,
			ProductItemUuid:     "",
			Name:                product.Name,
			Description:         product.Description,
			Price:               0,
			Discount:            0,
			DiscountPrice:       0,
			Stock:               0,
			ProductCategoryUuid: product.ProductCategoryUuid,
			Type:                "",
		}

		if productItems, ok := productItemsMap[product.Uuid]; ok {

			if len(productItems) > 0 {
				sort.Sort(model.ProductItemByPrice(productItems))
				productRes.ProductItemUuid = productItems[0].Uuid
				productRes.Price = productItems[0].Price
				productRes.Discount = productItems[0].Discount
				productRes.DiscountPrice = productItems[0].DiscountPrice
				productRes.Stock = productItems[0].Stock
				productRes.Type = productItems[0].Variants
			}
		}

		if images, ok := mProductImages[product.Uuid]; ok {
			for _, image := range images {
				if resource, ok := resourceMap[image]; ok {
					productRes.Images = append(productRes.Images, resource.Address)
				}
			}
		}

		res = append(res, productRes)
	}

	return &model.PagedResponse{
		Total: total,
		Data:  res,
	}, nil
}

// UpdateProduct
func (p *ProductService) UpdateProduct(ctx *app.Context, params *model.ReqProductUpdate) (err error) {

	product := &model.Product{}
	err = ctx.DB.Where("uuid = ?", params.Uuid).First(&product).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product by UUID", err)
		return errors.New("failed to get product by UUID")
	}

	if params.Name != "" {
		product.Name = params.Name
	}

	if params.Description != "" {
		product.Description = params.Description
	}

	if params.ProductCategoryUuid != "" {
		product.ProductCategoryUuid = params.ProductCategoryUuid
	}

	if len(params.Images) > 0 {
		product.Images = utils.ArrayToJsonString(params.Images)
	}

	if len(params.Videos) > 0 {
		product.Videos = utils.ArrayToJsonString(params.Videos)
	}

	if params.AliasName != "" {
		product.AliasName = params.AliasName
	}

	if params.ProductStatus != "" {
		product.Status = params.ProductStatus
	}

	if params.StockWarning > 0 {
		product.StockWarning = params.StockWarning
	}

	product.StockWarningSell = params.StockWarningSell

	if params.Weight > 0 {
		product.Weight = params.Weight
	}

	if params.Length > 0 {
		product.Length = params.Length
	}

	if params.Height > 0 {
		product.Height = params.Height
	}

	if params.Weight > 0 {
		product.Weight = params.Weight
	}

	if params.Unit != "" {
		product.Unit = params.Unit
	}

	product.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	err = ctx.DB.Where("uuid = ?", params.Uuid).Updates(&product).Error
	if err != nil {
		ctx.Logger.Error("Failed to update product", err)
		return errors.New("failed to update product")
	}

	return nil
}

// UpdateProductSku
func (p *ProductService) UpdateProductSku(ctx *app.Context, params *model.ReqProductItemUpdate) (err error) {

	productItem := &model.ProductItem{}
	err = ctx.DB.Where("uuid = ?", params.Uuid).First(&productItem).Error
	if err != nil {
		ctx.Logger.Error("Failed to get product item by UUID", err)
		return errors.New("failed to get product item by UUID")
	}

	if params.Price > 0 {
		productItem.Price = params.Price
	}

	if params.Discount > 0 {
		productItem.Discount = params.Discount
	}

	if params.DiscountPrice > 0 {
		productItem.DiscountPrice = params.DiscountPrice
	}

	if params.Stock > 0 {
		productItem.Stock = params.Stock
	}

	if params.Description != "" {
		productItem.Description = params.Description
	}

	productItem.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	err = ctx.DB.Where("uuid = ?", params.Uuid).Updates(&productItem).Error
	if err != nil {
		ctx.Logger.Error("Failed to update product item", err)
		return errors.New("failed to update product item")
	}

	return nil
}
