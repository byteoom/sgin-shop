package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type ConfigurationService struct {
}

func NewConfigurationService() *ConfigurationService {
	return &ConfigurationService{}
}

func (s *ConfigurationService) CreateOrUpdateConfiguration(ctx *app.Context, config *model.Configuration) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	config.CreatedAt = now
	config.UpdatedAt = now

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {

		// 查询配置是否存在
		oldConfig := &model.Configuration{}
		err := tx.Where("category = ? AND name = ?", config.Category, config.Name).First(oldConfig).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				ctx.Logger.Error("Failed to get configuration by category and name", err)
				return errors.New("failed to get configuration by category and name")
			}
		}
		if err == nil && oldConfig.Id > 0 {
			// 更新
			config.Id = oldConfig.Id
			err = tx.Where("id = ?", config.Id).Updates(config).Error
			if err != nil {
				ctx.Logger.Error("Failed to update configuration", err)
				return errors.New("failed to update configuration")
			}
			return nil
		}

		err = tx.Create(config).Error
		if err != nil {
			ctx.Logger.Error("Failed to create configuration", err)
			return errors.New("failed to create configuration")
		}
		return nil
	})

	if err != nil {
		ctx.Logger.Error("Failed to create configuration", err)
		return errors.New("failed to create configuration")
	}
	return nil
}

// CreateOrUpdateConfigurationMap creates or updates a map of configurations
func (s *ConfigurationService) CreateOrUpdateConfigurationMap(ctx *app.Context, configMap map[string]string, category string) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		for name, value := range configMap {
			config := &model.Configuration{
				Category:  category,
				Name:      name,
				Value:     value,
				CreatedAt: now,
				UpdatedAt: now,
			}

			// 查询配置是否存在
			oldConfig := &model.Configuration{}
			err := tx.Where("category = ? AND name = ?", config.Category, config.Name).First(oldConfig).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					ctx.Logger.Error("Failed to get configuration by category and name", err)
					tx.Rollback()
					return errors.New("failed to get configuration by category and name")
				}
			}
			if err == nil && oldConfig.Id > 0 {
				// 更新
				config.Id = oldConfig.Id
				err = tx.Where("id = ?", config.Id).Updates(config).Error
				if err != nil {
					ctx.Logger.Error("Failed to update configuration", err)
					tx.Rollback()
					return errors.New("failed to update configuration")
				}
				continue
			}

			err = tx.Create(config).Error
			if err != nil {
				ctx.Logger.Error("Failed to create configuration", err)
				tx.Rollback()
				return errors.New("failed to create configuration")
			}
		}
		return nil
	})

	if err != nil {
		ctx.Logger.Error("Failed to create configuration map", err)
		return errors.New("failed to create configuration map")
	}
	return nil
}

func (s *ConfigurationService) GetConfigurationByID(ctx *app.Context, id int) (*model.Configuration, error) {
	config := &model.Configuration{}
	err := ctx.DB.Where("id = ?", id).First(config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("configuration not found")
		}
		ctx.Logger.Error("Failed to get configuration by ID", err)
		return nil, errors.New("failed to get configuration by ID")
	}
	return config, nil
}

func (s *ConfigurationService) UpdateConfiguration(ctx *app.Context, config *model.Configuration) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	config.UpdatedAt = now
	err := ctx.DB.Where("id = ?", config.Id).Updates(config).Error
	if err != nil {
		ctx.Logger.Error("Failed to update configuration", err)
		return errors.New("failed to update configuration")
	}

	return nil
}

func (s *ConfigurationService) DeleteConfiguration(ctx *app.Context, id int) error {
	err := ctx.DB.Model(&model.Configuration{}).Where("id = ?", id).Update("is_deleted", 1).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete configuration", err)
		return errors.New("failed to delete configuration")
	}

	return nil
}

// 根据category获取配置列表
func (s *ConfigurationService) GetConfigurationListByCategory(ctx *app.Context, category string) ([]*model.Configuration, error) {
	configs := make([]*model.Configuration, 0)
	err := ctx.DB.Where("category = ?", category).Find(&configs).Error
	if err != nil {
		ctx.Logger.Error("Failed to get configuration list by category", err)
		return nil, errors.New("failed to get configuration list by category")
	}
	return configs, nil
}

func (s *ConfigurationService) GetConfigurationMapByCategory(ctx *app.Context, category string) (map[string]string, error) {
	list, err := s.GetConfigurationListByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	configMap := make(map[string]string)
	for _, config := range list {
		configMap[config.Name] = config.Value
	}
	return configMap, nil
}

// 根据配置category和name获取配置
func (s *ConfigurationService) GetConfigurationByCategoryAndName(ctx *app.Context, category, name string) (*model.Configuration, error) {
	config := &model.Configuration{}
	err := ctx.DB.Where("category = ? AND name = ?", category, name).First(config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return config, nil
		}
		ctx.Logger.Error("Failed to get configuration by category and name", err)
		return nil, errors.New("failed to get configuration by category and name")
	}
	return config, nil
}

// GetConfigurationList retrieves a list of configurations based on query parameters
func (s *ConfigurationService) GetConfigurationList(ctx *app.Context, params *model.ReqConfigQueryParam) (*model.PagedResponse, error) {
	var (
		configurations []*model.Configuration
		total          int64
	)

	db := ctx.DB.Model(&model.Configuration{})

	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.Category != "" {
		db = db.Where("category = ?", params.Category)
	}

	db = db.Where("is_deleted = ?", 0)

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get configuration count", err)
		return nil, errors.New("failed to get configuration count")
	}

	err = db.Offset(params.GetOffset()).Limit(params.PageSize).Find(&configurations).Error
	if err != nil {
		ctx.Logger.Error("Failed to get configuration list", err)
		return nil, errors.New("failed to get configuration list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  configurations,
	}, nil
}
