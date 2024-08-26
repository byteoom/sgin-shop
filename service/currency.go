package service

import (
	"errors"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CurrencyService struct {
}

func NewCurrencyService() *CurrencyService {
	return &CurrencyService{}
}

func (s *CurrencyService) CreateCurrency(ctx *app.Context, currency *model.Currency) error {

	currency.Uuid = uuid.New().String()

	err := ctx.DB.Create(currency).Error
	if err != nil {
		ctx.Logger.Error("Failed to create currency", err)
		return errors.New("failed to create currency")
	}
	return nil
}

func (s *CurrencyService) GetCurrencyByUUID(ctx *app.Context, uuid string) (*model.Currency, error) {
	currency := &model.Currency{}
	err := ctx.DB.Where("uuid = ?", uuid).First(currency).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("currency not found")
		}
		ctx.Logger.Error("Failed to get currency by UUID", err)
		return nil, errors.New("failed to get currency by UUID")
	}
	return currency, nil
}

func (s *CurrencyService) UpdateCurrency(ctx *app.Context, currency *model.Currency) error {

	err := ctx.DB.Where("uuid = ?", currency.Uuid).Updates(currency).Error
	if err != nil {
		ctx.Logger.Error("Failed to update currency", err)
		return errors.New("failed to update currency")
	}

	return nil
}

func (s *CurrencyService) DeleteCurrency(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.Currency{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete currency", err)
		return errors.New("failed to delete currency")
	}

	return nil
}

// GetCurrencyList retrieves a list of currencies based on query parameters
func (s *CurrencyService) GetCurrencyList(ctx *app.Context, params *model.ReqCurrencyQueryParam) (*model.PagedResponse, error) {
	var (
		currencies []*model.Currency
		total      int64
	)

	db := ctx.DB.Model(&model.Currency{})

	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get currency count", err)
		return nil, errors.New("failed to get currency count")
	}

	err = db.Offset(params.GetOffset()).Limit(params.PageSize).Find(&currencies).Error
	if err != nil {
		ctx.Logger.Error("Failed to get currency list", err)
		return nil, errors.New("failed to get currency list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  currencies,
	}, nil
}

// 获取全部可用的币种
func (s *CurrencyService) GetAvailableCurrencyList(ctx *app.Context) ([]*model.Currency, error) {
	var currencies []*model.Currency
	err := ctx.DB.Model(&model.Currency{}).Where("status = ?", model.CurrencyStatusEnabled).Find(&currencies).Error
	if err != nil {
		ctx.Logger.Error("Failed to get available currency list", err)
		return nil, errors.New("failed to get available currency list")
	}
	return currencies, nil
}

// 根据uuid获取币种
func (s *CurrencyService) GetCurrencyByUuids(ctx *app.Context, uuids []string) (map[string]*model.Currency, error) {
	currencies := make([]*model.Currency, 0)
	err := ctx.DB.Model(&model.Currency{}).Where("uuid IN (?)", uuids).Find(&currencies).Error
	if err != nil {
		ctx.Logger.Error("Failed to get currency by uuid", err)
		return nil, errors.New("failed to get currency by uuid")
	}
	result := make(map[string]*model.Currency)
	for _, currency := range currencies {
		result[currency.Uuid] = currency
	}
	return result, nil
}
