package service

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceService struct {
}

func NewResourceService() *ResourceService {
	return &ResourceService{}
}

// CreateResource 创建资源
func (s *ResourceService) CreateResource(ctx *app.Context, resource *model.Resource) error {
	resource.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	resource.UpdatedAt = resource.CreatedAt
	resource.Uuid = uuid.New().String()

	err := ctx.DB.Create(resource).Error
	if err != nil {
		ctx.Logger.Error("Failed to create resource", err)
		return errors.New("failed to create resource")
	}
	return nil
}

// 创建目录
func (s *ResourceService) CreateFolder(ctx *app.Context, param *model.Resource) error {

	// 根据path获取目录
	var oldResource model.Resource
	err := ctx.DB.Where("path = ? and address = ?", param.Path, param.Address).First(&oldResource).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Logger.Error("Failed to get resource by path", err)
		return errors.New("failed to get resource by path")
	}

	if err == nil && oldResource.Uuid != "" {
		return errors.New("folder already exists")
	}

	err = ctx.DB.Create(param).Error
	if err != nil {
		ctx.Logger.Error("Failed to create resource", err)
		return errors.New("failed to create resource")
	}
	return nil
}

func (s *ResourceService) CreateResourceList(ctx *app.Context, resource []*model.Resource) error {

	err := ctx.DB.Create(resource).Error
	if err != nil {
		ctx.Logger.Error("Failed to create resource", err)
		return errors.New("failed to create resource")
	}
	return nil
}

// GetResourceByUUID 根据UUID获取资源
func (s *ResourceService) GetResourceByUUID(ctx *app.Context, uuid string) (*model.Resource, error) {
	resource := &model.Resource{}
	err := ctx.DB.Where("uuid = ?", uuid).First(resource).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("resource not found")
		}
		ctx.Logger.Error("Failed to get resource by UUID", err)
		return nil, errors.New("failed to get resource by UUID")
	}
	return resource, nil
}

// GetResourceByPath 根据path获取文件夹资源
func (s *ResourceService) GetResourceFolderByPath(ctx *app.Context, path string) (*model.Resource, error) {
	resource := &model.Resource{}
	err := ctx.DB.Where("address = ? and type = ?", path, model.ResourceTypeFolder).First(resource).Error
	if err != nil {
		// 如果不存在则创建
		if err == gorm.ErrRecordNotFound {
			resource.Uuid = uuid.New().String()
			resource.Name = filepath.Base(path)
			resource.Address = path
			resource.Path = path
			resource.Type = model.ResourceTypeFolder
			resource.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
			resource.UpdatedAt = resource.CreatedAt
			err = ctx.DB.Create(resource).Error
			if err != nil {
				ctx.Logger.Error("Failed to create resource", err)
				return nil, errors.New("failed to create resource")
			}
			return resource, nil
		}
		return nil, err
	}
	return resource, nil
}

// UpdateResource 更新资源
func (s *ResourceService) UpdateResource(ctx *app.Context, resource *model.Resource) error {
	resource.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	err := ctx.DB.Where("uuid = ?", resource.Uuid).Updates(resource).Error
	if err != nil {
		ctx.Logger.Error("Failed to update resource", err)
		return errors.New("failed to update resource")
	}

	return nil
}

//

// DeleteResource 删除资源
func (s *ResourceService) DeleteResource(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ? ", uuid).Delete(&model.Resource{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete resource", err)
		return errors.New("failed to delete resource")
	}

	return nil
}

// 根据uuid获取子资源数量
func (s *ResourceService) GetChildResourceCountByUUID(ctx *app.Context, uuid string) (int64, error) {
	var count int64
	err := ctx.DB.Model(&model.Resource{}).Where("parent_uuid = ?", uuid).Count(&count).Error
	if err != nil {
		ctx.Logger.Error("Failed to get child resource by UUID", err)
		return 0, errors.New("failed to get child resource by UUID")
	}
	return count, nil
}

// 获取文件夹列表
func (s *ResourceService) GetFolderList(ctx *app.Context) ([]*model.ResourceRes, error) {
	var (
		folderList []*model.Resource
	)

	err := ctx.DB.Where("type = ?", model.ResourceTypeFolder).Find(&folderList).Error
	if err != nil {
		ctx.Logger.Error("Failed to get folder list", err)
		return nil, errors.New("failed to get folder list")
	}

	return s.buildResourceTree(folderList, "")
}

// MoveResource 移动资源
func (s *ResourceService) MoveResource(ctx *app.Context, uuidlist []string, parentUuid, targetPath string) error {
	err := ctx.DB.Transaction(func(tx *gorm.DB) error {

		var resourceList []*model.Resource

		err := tx.Where("uuid IN (?)", uuidlist).Find(&resourceList).Error
		if err != nil {
			ctx.Logger.Error("Failed to get resource list", err)
			tx.Rollback()
			return errors.New("failed to get resource list")
		}

		// 判断目标文件夹是否存在
		targetAddress := filepath.Join(ctx.Config.PkgFileDir, targetPath)
		if _, err := os.Stat(targetAddress); err != nil {
			if os.IsNotExist(err) {
				// 创建目标文件夹
				err := os.MkdirAll(targetAddress, os.ModePerm)
				if err != nil {
					ctx.Logger.Error("Failed to create target folder", err)
					return errors.New("failed to create target folder")
				}
			} else {
				ctx.Logger.Error("Failed to get target folder", err)
				tx.Rollback()
				return errors.New("failed to get target folder")
			}
		}

		for _, resource := range resourceList {
			address := filepath.Join(targetPath, resource.Name)
			err := os.Rename(filepath.Join(ctx.Config.PkgFileDir, resource.Address), filepath.Join(ctx.Config.PkgFileDir, address))
			if err != nil {
				// 移动失败, 记录失败信息。 不回滚
				ctx.Logger.Errorf("Failed to move resource err: %v, uuid:%s name:%s address:%s", err, resource.Uuid, resource.Name, resource.Address)

				continue
			}

			err = tx.Model(&model.Resource{}).Where("uuid = ?", resource.Uuid).Updates(map[string]interface{}{
				"parent_uuid": parentUuid,
				"address":     address,
				"path":        targetPath,
			}).Error
			if err != nil {
				ctx.Logger.Error("Failed to update resource", err)
				// 移动失败, 记录失败信息。 不回滚
				continue
			}
		}

		return nil
	})
	if err != nil {
		ctx.Logger.Error("Failed to move resource", err)
		return errors.New("failed to move resource")
	}
	return nil

}

// buildResourceTree 构建资源树
func (s *ResourceService) buildResourceTree(resources []*model.Resource, parentUuid string) ([]*model.ResourceRes, error) {
	var resourceList []*model.ResourceRes

	// 遍历所有资源，查找当前父UUID下的直接子资源
	for _, resource := range resources {
		if resource.ParentUuid == parentUuid {
			// 为当前资源创建ResourceRes结构，可能需要根据实际情况调整字段
			newResource := &model.ResourceRes{
				Resource: *resource,
				Children: nil, // 初始化子节点切片
			}

			// 递归调用自身，查找并构建所有子资源的树形结构
			children, err := s.buildResourceTree(resources, resource.Uuid)
			if err != nil {
				return nil, err // 递归过程中出错
			}

			newResource.Children = children
			resourceList = append(resourceList, newResource)
		}
	}
	return resourceList, nil
}

// GetResourceList 获取资源列表
func (s *ResourceService) GetResourceList(ctx *app.Context, params *model.ReqResourceQueryParam) (*model.PagedResponse, error) {
	var (
		resources []*model.Resource
		total     int64
	)

	db := ctx.DB.Model(&model.Resource{})

	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.Path != "" {
		db = db.Where("path = ?", params.Path)
	}

	if params.ParentUuid != "" {
		db = db.Where("parent_uuid = ?", params.ParentUuid)
	}

	if params.MimeType != "" {
		db = db.Where("mime_type LIKE ?", params.MimeType+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get resource count", err)
		return nil, errors.New("failed to get resource count")
	}

	err = db.Find(&resources).Error
	if err != nil {
		ctx.Logger.Error("Failed to get resource list", err)
		return nil, errors.New("failed to get resource list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  resources,
	}, nil
}

// 根据uuid列表获取资源列表
func (s *ResourceService) GetResourceByUUIDList(ctx *app.Context, uuidList []string) (map[string]*model.Resource, error) {
	resourceMap := make(map[string]*model.Resource, 0)
	var resources []*model.Resource
	err := ctx.DB.Where("uuid IN (?)", uuidList).Find(&resources).Error
	if err != nil {
		ctx.Logger.Error("Failed to get resource list by UUID list", err)
		return nil, errors.New("failed to get resource list by UUID list")
	}

	for _, resource := range resources {
		resourceMap[resource.Uuid] = resource
	}

	return resourceMap, nil
}
