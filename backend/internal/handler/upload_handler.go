package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/service"
	"github.com/techblog/community/internal/util"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	uploadSvc *service.UploadService
}

// NewUploadHandler 构造上传处理器
func NewUploadHandler(uploadSvc *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadSvc: uploadSvc}
}

// Image POST /api/v1/upload/image 上传图片
func (h *UploadHandler) Image(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "上传参数校验失败: field=file")
		return
	}
	url, err := h.uploadSvc.SaveImage(c.Request.Context(), util.CurrentUserID(c), file)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, gin.H{"url": url})
}
