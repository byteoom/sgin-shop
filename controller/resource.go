package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/utils"
	"sgin/service"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceController struct {
	ResourceService *service.ResourceService
}

// 查询资源列表
// @Summary 查询资源列表
// @Tags 资源
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqResourceQueryParam false "查询参数"
// @Success 200 {object} model.ResourceQueryResponse
// @Router /api/v1/resource/list [post]
func (c *ResourceController) GetResourceList(ctx *app.Context) {
	param := &model.ReqResourceQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	resources, err := c.ResourceService.GetResourceList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(resources)
}

// 创建资源
// @Summary 创建资源
// @Tags 资源
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Resource true "资源信息"
// @Success 200 {object} model.ResourceResponse
// @Router /api/v1/resource/create [post]
func (c *ResourceController) CreateResource(ctx *app.Context) {

	targetPath := ctx.PostForm("path")

	parentUuid := ""

	if targetPath != "" {
		// 先获取目标文件夹信息
		parentResource, err := c.ResourceService.GetResourceFolderByPath(ctx, targetPath)
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, err.Error())
			return
		}
		parentUuid = parentResource.Uuid
	}

	// Handle file uploads
	form := ctx.Request.MultipartForm

	files := form.File["files"]

	resourceList := make([]*model.Resource, 0)
	addResourceList := make([]*model.Resource, 0)

	now := time.Now().Format("2006-01-02 15:04:05")

	for _, fileHeader := range files {

		// Create a unique filename and save the file
		filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
		address := filepath.Join(targetPath, filename)
		ctx.Logger.Info("address:", address)
		resourceItem := &model.Resource{
			Uuid:       uuid.New().String(),
			Name:       fileHeader.Filename,
			ParentUuid: parentUuid,
			Type:       model.ResourceTypeFile,
			MimeType:   fileHeader.Header.Get("Content-Type"),
			Size:       fileHeader.Size,
			Path:       targetPath,
			Address:    address,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		saveAddress := filepath.Join(ctx.Config.Upload.Dir, address)

		// 获取文件md5
		file, err := fileHeader.Open()
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, "Cannot open file")
			return
		}
		defer file.Close()

		md5str, err := utils.GetFileMd5(file)

		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, "Cannot get file md5")
			return
		}

		resourceItem.Md5 = md5str
		exsitRes, err := c.ResourceService.GetResourceFileByPathAndMd5(ctx, targetPath, md5str)
		if err != nil && err != gorm.ErrRecordNotFound {
			ctx.JSONError(http.StatusInternalServerError, err.Error())
			return
		}
		if err != nil && err == gorm.ErrRecordNotFound {

			ctx.Logger.Info("saveAddress:", saveAddress)
			addResourceList = append(addResourceList, resourceItem)

			err := ctx.SaveUploadedFile(fileHeader, filepath.Join(ctx.Config.Upload.Dir, address))
			if err != nil {
				ctx.JSONError(http.StatusInternalServerError, "Cannot save file")
				return
			}
		} else {
			resourceItem = exsitRes
		}
		resourceList = append(resourceList, resourceItem)

	}

	if len(addResourceList) > 0 {
		err := c.ResourceService.CreateResourceList(ctx, addResourceList)
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSONSuccess(resourceList)
}

// 更新资源
// @Summary 更新资源
// @Tags 资源
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Resource true "资源信息"
// @Success 200 {object} model.ResourceResponse
// @Router /api/v1/resource/update [post]
func (c *ResourceController) UpdateResource(ctx *app.Context) {
	param := &model.Resource{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.ResourceService.UpdateResource(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(param)
}

// 移动资源
// @Summary 移动资源
// @Tags 资源
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqResourceMoveParam true "移动参数"
// @Success 200 {object} app.Response
// @Router /api/v1/resource/move [post]
func (c *ResourceController) MoveResource(ctx *app.Context) {
	param := &model.ReqResourceMove{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	var targetPath string = "/"

	if param.ParentUuid != "" {

		// 先获取目标文件夹信息
		parentResource, err := c.ResourceService.GetResourceByUUID(ctx, param.ParentUuid)
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, err.Error())
			return
		}
		if parentResource.Type != model.ResourceTypeFolder {
			ctx.JSONError(http.StatusBadRequest, "目标资源不是文件夹")
			return
		}

		targetPath = parentResource.Path

	}

	err := c.ResourceService.MoveResource(ctx, param.UuidList, param.ParentUuid, targetPath)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("ok")

}

// 删除资源
// @Summary 删除资源
// @Tags 资源
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqResourceDeleteParam true "删除参数"
// @Success 200 {object} app.Response
// @Router /api/v1/resource/delete [post]
func (c *ResourceController) DeleteResource(ctx *app.Context) {
	param := &model.ReqResourceDeleteParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	//先获取资源信息
	resource, err := c.ResourceService.GetResourceByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	if resource.Type == model.ResourceTypeFolder {

		count, err := c.ResourceService.GetChildResourceCountByUUID(ctx, param.Uuid)
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, err.Error())
			return
		}

		if count > 0 {
			ctx.JSONError(http.StatusBadRequest, "资源下有子资源，不能删除")
			return
		}
	}

	// 先获取资源是否存在
	fileInfo, err := os.Stat(ctx.Config.Upload.Dir + resource.Address)
	if err != nil && !os.IsNotExist(err) {

		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	if err == nil && fileInfo != nil && fileInfo.Name() != "" {
		err = os.Remove(ctx.Config.Upload.Dir + resource.Address)
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = c.ResourceService.DeleteResource(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("ok")
}

// 获取文件夹列表
// @Summary 获取文件夹列表
// @Tags 资源
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Success 200 {object} []model.ResourceRes
// @Router /api/v1/resource/folders [get]
func (c *ResourceController) GetFolderList(ctx *app.Context) {
	folders, err := c.ResourceService.GetFolderList(ctx)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(folders)
}

// 创建文件夹

func (c *ResourceController) CreateFolder(ctx *app.Context) {
	param := &model.ReqResourceCreateFolder{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	resource := &model.Resource{
		Uuid:      uuid.New().String(),
		Name:      param.Name,
		Type:      model.ResourceTypeFolder,
		Path:      "/",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	if param.Path != "" {
		// 先获取目标文件夹信息
		parentResource, err := c.ResourceService.GetResourceFolderByPath(ctx, param.Path)
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, err.Error())
			return
		}
		resource.ParentUuid = parentResource.Uuid

		resource.Path = parentResource.Address

	}
	resource.Address = resource.Path + "/" + resource.Name
	err := c.ResourceService.CreateFolder(ctx, resource)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}
