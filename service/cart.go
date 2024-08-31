package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartService struct {
}

func NewCartService() *CartService {
	return &CartService{}
}

// CreateCart creates a new cart item
func (s *CartService) CreateCart(ctx *app.Context, cart *model.Cart) error {
	cart.CreatedAt = time.Now().Format(time.DateTime)
	cart.UpdatedAt = cart.CreatedAt
	cart.Uuid = uuid.New().String()

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		// Check if the product exists
		// 先查询商品是否存在
		product := &model.ProductItem{}
		err := tx.Where("uuid = ?", cart.ProductItemUuid).First(product).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("product not found")
			}
			ctx.Logger.Error("Failed to get product by UUID", err)
			return errors.New("failed to get product by UUID")
		}

		// 检查购物车中是否已经存在该商品
		// Check if the product already exists in the cart
		oldCart := &model.Cart{}
		err = tx.Where("user_id = ? and product_item_uuid = ?", cart.UserID, cart.ProductItemUuid).First(oldCart).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			ctx.Logger.Error("Failed to get cart by product item UUID", err)
			return errors.New("failed to get cart by product item UUID")
		}

		if oldCart.ID > 0 {
			oldCart.Quantity += cart.Quantity
			oldCart.UpdatedAt = time.Now().Format(time.DateTime)
			err = tx.Save(oldCart).Error
			if err != nil {
				ctx.Logger.Error("Failed to update cart", err)
				return errors.New("failed to update cart")
			}
			return nil
		}

		// Create a new cart item
		// 创建购物车
		err = tx.Create(cart).Error
		if err != nil {
			ctx.Logger.Error("Failed to create cart", err)
			return errors.New("failed to create cart")
		}

		return nil
	})

	if err != nil {
		ctx.Logger.Error("Failed to create cart", err)
		return err
	}

	return nil
}

// GetCartByUUID retrieves a cart item by its UUID
func (s *CartService) GetCartByUUID(ctx *app.Context, uuid string) (*model.Cart, error) {
	cart := &model.Cart{}
	err := ctx.DB.Where("uuid = ?", uuid).First(cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("cart not found")
		}
		ctx.Logger.Error("Failed to get cart by UUID", err)
		return nil, errors.New("failed to get cart by UUID")
	}
	return cart, nil
}

// UpdateCart updates an existing cart item
func (s *CartService) UpdateCart(ctx *app.Context, cart *model.Cart) error {
	cart.UpdatedAt = time.Now().Format(time.DateTime)
	err := ctx.DB.Where("uuid = ?", cart.Uuid).Updates(cart).Error
	if err != nil {
		ctx.Logger.Error("Failed to update cart", err)
		return errors.New("failed to update cart")
	}

	return nil
}

// UpdateCartItemCount updates the quantity of a cart item
func (s *CartService) UpdateCartItemCount(ctx *app.Context, uuid string, quantity int) error {
	err := ctx.DB.Model(&model.Cart{}).Where("uuid = ?", uuid).Update("quantity", quantity).Error
	if err != nil {
		ctx.Logger.Error("Failed to update cart item count", err)
		return errors.New("failed to update cart item count")
	}

	return nil
}

// DeleteCart deletes a cart item by its UUID
func (s *CartService) DeleteCart(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.Cart{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete cart", err)
		return errors.New("failed to delete cart")
	}

	return nil
}

// 根据购物车UUID 列表获取商品列表
func (s *CartService) GetCartByUUIDList(ctx *app.Context, uuids []string) (map[string]*model.CartProductItemRes, error) {
	carts := make([]*model.Cart, 0)
	err := ctx.DB.Where("uuid in (?)", uuids).Find(&carts).Error
	if err != nil {
		ctx.Logger.Error("Failed to get cart by UUID list", err)
		return nil, errors.New("failed to get cart by UUID list")
	}

	productItemUuids := make([]string, 0)
	for _, cart := range carts {
		productItemUuids = append(productItemUuids, cart.ProductItemUuid)
	}

	productItemMap, err := NewProductService().GetProductItemByUUIDList(ctx, productItemUuids)

	if err != nil {
		ctx.Logger.Error("Failed to get product item by UUID list", err)
		return nil, errors.New("failed to get product item by UUID list")
	}

	res := make(map[string]*model.CartProductItemRes)
	for _, cart := range carts {
		item := &model.CartProductItemRes{
			Cart: *cart,
		}
		if productItem, ok := productItemMap[cart.ProductItemUuid]; ok {
			item.ProductItem = productItem
		}
		res[cart.Uuid] = item
	}

	return res, nil
}

// GetCartList retrieves a list of cart items based on query parameters
func (s *CartService) GetCartList(ctx *app.Context, params *model.ReqCartQueryParam) (r *model.PagedResponse, err error) {
	var (
		carts []*model.Cart
		total int64
	)

	db := ctx.DB.Model(&model.Cart{})

	if params.UserID != "" {
		db = db.Where("user_id = ?", params.UserID)
	}

	err = db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get cart count", err)
		return nil, errors.New("failed to get cart count")
	}

	err = db.Order("id DESC").Offset(params.GetOffset()).Limit(params.PageSize).Find(&carts).Error

	if err != nil {
		ctx.Logger.Error("Failed to get cart list", err)
		return nil, errors.New("failed to get cart list")
	}

	productItemUuids := make([]string, 0)
	for _, cart := range carts {
		productItemUuids = append(productItemUuids, cart.ProductItemUuid)
	}

	productItemMap, err := NewProductService().GetProductItemByUUIDList(ctx, productItemUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get product item by UUID list", err)
		return nil, errors.New("failed to get product item by UUID list")
	}
	res := make([]*model.CartProductItemRes, 0)
	for _, cart := range carts {
		item := &model.CartProductItemRes{
			Cart: *cart,
		}
		if productItem, ok := productItemMap[cart.ProductItemUuid]; ok {
			item.ProductItem = productItem
		}
		res = append(res, item)
	}

	return &model.PagedResponse{
		Total: total,
		Data:  res,
	}, nil
}
