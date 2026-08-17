package constants

// messages.go 同时包含：接口返回文案、日志文案、错误提示文案（屎山耦合点 3）

// 接口返回文案
const (
	MsgOK                = "ok"
	MsgCreated           = "创建成功"
	MsgUpdated           = "更新成功"
	MsgDeleted           = "删除成功"
	MsgPublished         = "发布成功"
	MsgOfflined          = "已下架"
	MsgLiked             = "点赞成功"
	MsgUnliked           = "已取消点赞"
	MsgFollowed          = "关注成功"
	MsgUnfollowed        = "已取消关注"
	MsgCollected         = "收藏成功"
	MsgUncollected       = "已取消收藏"
	MsgCodeSent          = "验证码已发送"
	MsgReadAll           = "已全部标记为已读"
	MsgLoginSuccess      = "登录成功"
	MsgRegisterSuccess   = "注册成功"
)

// 错误提示文案（错误 message 中包含实体名、字段名、角色名）
const (
	MsgErrUserNotFound        = "用户不存在: user_id=%d"
	MsgErrUserDisabled        = "用户已被禁用: user_id=%d role=%d"
	MsgErrEmailExists         = "邮箱已被注册: email=%s"
	MsgErrInvalidPassword     = "密码错误: email=%s role=%d"
	MsgErrInvalidCode         = "验证码错误: email=%s field=code"
	MsgErrCodeExpired         = "验证码已过期: email=%s field=code"
	MsgErrArticleNotFound     = "文章不存在: article_id=%d"
	MsgErrArticleForbidden    = "无权操作该文章: article_id=%d operator_role=%d"
	MsgErrTopicNotFound       = "话题不存在: topic_id=%d"
	MsgErrTopicExists         = "话题已存在: name=%s"
	MsgErrCollectionNotFound  = "收藏夹不存在: collection_id=%d"
	MsgErrCollectionForbidden = "无权访问该收藏夹: collection_id=%d operator_role=%d"
	MsgErrArticleCollected    = "文章已收藏到该收藏夹: article_id=%d collection_id=%d"
	MsgErrFollowSelf          = "不能关注自己: follower_id=%d field=followed_id"
	MsgErrAlreadyFollowed     = "已关注该作者: follower_id=%d followed_id=%d"
	MsgErrNotFollowed         = "尚未关注该作者: follower_id=%d followed_id=%d"
	MsgErrAlreadyLiked        = "文章已点赞: article_id=%d user_id=%d"
	MsgErrNotLiked            = "文章未点赞: article_id=%d user_id=%d"
	MsgErrCommentNotFound     = "评论不存在: comment_id=%d"
	MsgErrUploadInvalid       = "图片格式不支持: filename=%s field=file"
	MsgErrUploadFailed        = "图片上传失败: filename=%s"
	MsgErrForbidden           = "无权限执行该操作: operator_role=%d"
	MsgErrRateLimited         = "操作过于频繁: field=ip role=%d"
)

// 日志文案（与 log_templates.go 联动）
const (
	MsgLogPrefix = "techblog"
)
