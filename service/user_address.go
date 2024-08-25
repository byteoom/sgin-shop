package service

import (
	"errors"
	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type UserAddressService struct{}

// NewUserAddressService creates a new instance of UserAddressService
func NewUserAddressService() *UserAddressService {
	return &UserAddressService{}
}

// CreateUserAddress creates a new user address
func (s *UserAddressService) CreateUserAddress(ctx *app.Context, address *model.UserAddress) error {
	if address.IsDefault {
		// Set other addresses to non-default if a new default address is being created
		if err := s.setAllAddressesNonDefault(ctx, address.UserID); err != nil {
			return err
		}
	}

	if err := ctx.DB.Create(address).Error; err != nil {
		ctx.Logger.Error("Failed to create user address", err)
		return errors.New("failed to create user address")
	}
	return nil
}

// GetUserAddressByID retrieves a user address by its ID
func (s *UserAddressService) GetUserAddressByID(ctx *app.Context, uuid string) (*model.UserAddress, error) {
	address := &model.UserAddress{}
	if err := ctx.DB.Where("uuid = ?", uuid).First(address).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user address not found")
		}
		ctx.Logger.Error("Failed to get user address", err)
		return nil, errors.New("failed to get user address")
	}
	return address, nil
}

// GetUserAddressesByUserID retrieves all addresses for a specific user
func (s *UserAddressService) GetUserAddressesByUserID(ctx *app.Context, userID string) ([]model.UserAddress, error) {
	var addresses []model.UserAddress
	if err := ctx.DB.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		ctx.Logger.Error("Failed to get user addresses", err)
		return nil, errors.New("failed to get user addresses")
	}
	return addresses, nil
}

// UpdateUserAddress updates an existing user address
func (s *UserAddressService) UpdateUserAddress(ctx *app.Context, address *model.UserAddress) error {
	if address.IsDefault {
		// Set other addresses to non-default if updating an address to be the default
		if err := s.setAllAddressesNonDefault(ctx, address.UserID); err != nil {
			return err
		}
	}

	if err := ctx.DB.Save(address).Error; err != nil {
		ctx.Logger.Error("Failed to update user address", err)
		return errors.New("failed to update user address")
	}
	return nil
}

// DeleteUserAddress deletes a user address by its ID
func (s *UserAddressService) DeleteUserAddress(ctx *app.Context, uuid string) error {
	if err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.UserAddress{}).Error; err != nil {
		ctx.Logger.Error("Failed to delete user address", err)
		return errors.New("failed to delete user address")
	}
	return nil
}

// Helper method to set all addresses for a user to non-default
func (s *UserAddressService) setAllAddressesNonDefault(ctx *app.Context, userID uint) error {
	if err := ctx.DB.Model(&model.UserAddress{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		ctx.Logger.Error("Failed to set addresses to non-default", err)
		return errors.New("failed to set addresses to non-default")
	}
	return nil
}
