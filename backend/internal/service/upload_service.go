package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/util"
)

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true,
}

// UploadService 图片上传服务
type UploadService struct {
	uploadDir string
}

// NewUploadService 构造上传服务
func NewUploadService(uploadDir string) *UploadService {
	return &UploadService{uploadDir: uploadDir}
}

// SaveImage 保存上传图片，返回可访问 URL
func (s *UploadService) SaveImage(ctx context.Context, userID uint, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", util.NewAppError(constants.CodeUploadInvalid, "未选择文件: field=file")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] {
		return "", util.NewAppError(constants.CodeUploadInvalid, fmt.Sprintf(constants.MsgErrUploadInvalid, file.Filename))
	}
	if file.Size > 5<<20 {
		return "", util.NewAppError(constants.CodeUploadInvalid, "图片大小不能超过 5MB: filename="+file.Filename)
	}
	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return "", util.WrapAppError(constants.CodeUploadFailed, fmt.Sprintf(constants.MsgErrUploadFailed, file.Filename), err)
	}
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.NewString()[:8], ext)
	dst := filepath.Join(s.uploadDir, filename)
	src, err := file.Open()
	if err != nil {
		return "", util.WrapAppError(constants.CodeUploadFailed, fmt.Sprintf(constants.MsgErrUploadFailed, file.Filename), err)
	}
	defer src.Close()
	out, err := os.Create(dst)
	if err != nil {
		return "", util.WrapAppError(constants.CodeUploadFailed, fmt.Sprintf(constants.MsgErrUploadFailed, file.Filename), err)
	}
	defer out.Close()
	if _, err := out.ReadFrom(src); err != nil {
		return "", util.WrapAppError(constants.CodeUploadFailed, fmt.Sprintf(constants.MsgErrUploadFailed, file.Filename), err)
	}
	url := "/uploads/" + filename
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUploadImage, userID, file.Filename, url))
	return url, nil
}
