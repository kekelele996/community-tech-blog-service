package dto

// RegisterRequest 邮箱+验证码注册
type RegisterRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Code     string   `json:"code" binding:"required,len=6"`
	Password string   `json:"password" binding:"required,min=6,max=64"`
	Nickname string   `json:"nickname" binding:"required,min=2,max=32"`
	Avatar   string   `json:"avatar" binding:"omitempty,max=512"`
	Bio      string   `json:"bio" binding:"omitempty,max=512"`
	TechTags []string `json:"tech_tags" binding:"omitempty,dive,max=32"`
}

// LoginRequest 邮箱+密码登录
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// LoginCodeRequest 邮箱验证码登录
type LoginCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// SendCodeRequest 发送邮箱验证码
type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// GithubLoginRequest GitHub 第三方登录（演示模式：code 为 GitHub OAuth 授权码）
type GithubLoginRequest struct {
	Code     string `json:"code" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=32"`
}

// SendCodeResponse 发送验证码响应：debug_code 仅演示/测试环境返回
type SendCodeResponse struct {
	Email     string `json:"email"`
	ExpiresIn int    `json:"expires_in"`
	DebugCode string `json:"debug_code,omitempty"`
}

// AuthResponse 登录/注册响应
type AuthResponse struct {
	Token string         `json:"token"`
	User  UserProfileDTO `json:"user"`
}
