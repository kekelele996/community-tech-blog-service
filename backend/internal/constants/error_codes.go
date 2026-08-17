package constants

// 全局错误码集中维护：所有 service/handler 与前端 utils/request.ts 共用
const (
	CodeSuccess          = 0
	CodeBadRequest       = 40000
	CodeUnauthorized     = 40100
	CodeForbidden        = 40300
	CodeNotFound         = 40400
	CodeConflict         = 40900
	CodeValidationFailed = 42200
	CodeRateLimited      = 42900
	CodeInternalError    = 50000

	// 用户/认证错误码
	CodeUserNotFound       = 40401
	CodeUserDisabled       = 40301
	CodeEmailExists        = 40901
	CodeInvalidCredentials = 40101
	CodeInvalidCode        = 42201
	CodeCodeExpired        = 42202
	CodeNicknameExists     = 40902

	// 文章错误码
	CodeArticleNotFound = 40402
	CodeArticleForbidden = 40302

	// 话题错误码
	CodeTopicNotFound = 40403
	CodeTopicExists   = 40903

	// 收藏夹错误码
	CodeCollectionNotFound = 40404
	CodeCollectionForbidden = 40304
	CodeArticleCollected   = 40904

	// 关注错误码
	CodeFollowSelf      = 42203
	CodeAlreadyFollowed = 40905
	CodeNotFollowed     = 40405

	// 点赞错误码
	CodeAlreadyLiked = 40906
	CodeNotLiked     = 40406

	// 评论错误码
	CodeCommentNotFound  = 40407
	CodeCommentForbidden = 40307

	// 通知错误码
	CodeNotificationNotFound = 40408

	// 上传错误码
	CodeUploadInvalid = 42204
	CodeUploadFailed  = 50001
)

// MessageOf 返回错误码对应的默认文案（与 constants/messages.go 联动）
func MessageOf(code int) string {
	switch code {
	case CodeSuccess:
		return "ok"
	case CodeBadRequest:
		return "请求参数错误"
	case CodeUnauthorized:
		return "未认证或登录已过期"
	case CodeForbidden:
		return "无权限执行该操作"
	case CodeNotFound:
		return "资源不存在"
	case CodeConflict:
		return "资源冲突"
	case CodeValidationFailed:
		return "参数校验失败"
	case CodeRateLimited:
		return "请求过于频繁，请稍后再试"
	case CodeInternalError:
		return "服务器内部错误"
	default:
		return "操作失败"
	}
}
