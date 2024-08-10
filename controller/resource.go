package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
	"time"

	"github.com/google/uuid"
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

	parentUuid := ctx.PostForm("parent_uuid")

	targetPath := "/"
	if parentUuid != "" {
		// 先获取目标文件夹信息
		parentResource, err := c.ResourceService.GetResourceByUUID(ctx, parentUuid)
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

	// Handle file uploads
	form := ctx.Request.MultipartForm

	files := form.File["files"]

	resourceList := make([]*model.Resource, 0)

	now := time.Now().Format("2006-01-02 15:04:05")

	for _, fileHeader := range files {

		// Create a unique filename and save the file
		filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
		address := filepath.Join(targetPath, filename)
		resourceList = append(resourceList, &model.Resource{
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
		})

		err := ctx.SaveUploadedFile(fileHeader, filepath.Join(ctx.Config.PkgFileDir, address))
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, "Cannot save file")
			return
		}

	}

	err := c.ResourceService.CreateResourceList(ctx, resourceList)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("ok")
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

	err = os.Remove(ctx.Config.PkgFileDir + resource.Address)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
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
		Uuid:       uuid.New().String(),
		Name:       param.Name,
		ParentUuid: param.ParentUuid,
		Type:       model.ResourceTypeFolder,
		Path:       "/" + param.Name,
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

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

		resource.Path = parentResource.Path + "/" + resource.Name

	}
	resource.Address = resource.Path
	err := c.ResourceService.CreateFolder(ctx, resource)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}
